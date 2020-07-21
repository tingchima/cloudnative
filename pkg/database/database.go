package database

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"

	"github.com/cenk/backoff"
)

// DBType 類型
type DBType string

const (
	// MySQL ...
	MySQL DBType = "mysql"
	// Postgres ...
	Postgres DBType = "postgres"
)

// RdbmsConn ...
type RdbmsConn struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

// RdbmsConfig for db config
type RdbmsConfig struct {
	Read       *Rdbms
	Write      *Rdbms
	Secrets    string `yaml:"secrets"`
	WithColor  bool   `yaml:"withColor"`
	WithCaller bool   `yaml:"withCaller"`
}

// Rdbms ...
type Rdbms struct {
	Type           DBType
	Debug          bool
	Host           string
	Port           int
	Username       string
	Password       string
	DBName         string
	MaxIdleConns   int
	MaxOpenConns   int
	MaxLifetimeSec int
	ReadTimeout    string `yaml:"read_timeout"`
	WriteTimeout   string `yaml:"write_timeout"`
	SearchPath     string `yaml:"search_path" mapstructure:"search_path"`
}

// InitRdbms init and return write and read DB objects
func InitRdbms(cfg *RdbmsConfig) (*RdbmsConn, error) {
	var err error
	var conn RdbmsConn

	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}

	conn.ReadDB, err = cfg.Read.Open()
	if err != nil {
		return nil, err
	}

	conn.WriteDB, err = cfg.Write.Open()
	if err != nil {
		return nil, err
	}

	if conn.ReadDB == nil {
		return nil, fmt.Errorf("read db initialization was failed, name: %s", cfg.Read.DBName)
	} else if err = conn.ReadDB.DB().Ping(); err != nil {
		return nil, fmt.Errorf("read db ping was failed")
	}

	if conn.WriteDB == nil {
		return nil, fmt.Errorf("write db initialization was failed, name: %s", cfg.Write.DBName)
	} else if err = conn.WriteDB.DB().Ping(); err != nil {
		return nil, fmt.Errorf("write db ping was failed")
	}

	if cfg.Read.Debug {
		conn.ReadDB = conn.ReadDB.LogMode(true)
		readLogger := log.With().Str("db_log", "read").Logger()
		conn.ReadDB.SetLogger(&GormLogger{logger: readLogger, WithColor: cfg.WithColor, WithCaller: cfg.WithCaller})
	}

	if cfg.Write.Debug {
		conn.WriteDB = conn.WriteDB.LogMode(true)
		writeLogger := log.With().Str("db_log", "write").Logger()
		conn.WriteDB.SetLogger(&GormLogger{logger: writeLogger, WithColor: cfg.WithColor, WithCaller: cfg.WithCaller})
	}

	return &conn, nil
}

// Open create database connection with dbType
func (r *Rdbms) Open() (*gorm.DB, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	var db *gorm.DB
	var err error

	switch r.Type {
	case MySQL:
		db, err = OpenMySQL(r)
		if err != nil {
			return db, err
		}
	case Postgres:
		db, err = OpenPostgres(r)
		if err != nil {
			return db, err
		}
	default:
		return nil, errors.New("Not support sql driver")
	}

	log.Info().Msgf("database ping success")

	if r.WriteTimeout == "" {
		r.WriteTimeout = "10s"
	}
	if r.ReadTimeout == "" {
		r.ReadTimeout = "10s"
	}

	if r.MaxIdleConns != 0 {
		db.DB().SetMaxIdleConns(r.MaxIdleConns)
	} else {
		db.DB().SetMaxIdleConns(2)
	}

	if r.MaxOpenConns != 0 {
		db.DB().SetMaxOpenConns(r.MaxOpenConns)
	} else {
		db.DB().SetMaxOpenConns(5)
	}

	if r.MaxLifetimeSec != 0 {
		db.DB().SetConnMaxLifetime(time.Duration(r.MaxLifetimeSec) * time.Second)
	} else {
		db.DB().SetConnMaxLifetime(14400 * time.Second)
	}

	return db, nil
}

// OpenMySQL ...
func OpenMySQL(database *Rdbms) (*gorm.DB, error) {
	var sqlDriver = "mysql"

	connectionStr := fmt.Sprintf(
		`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true&readTimeout=%s&writeTimeout=%s`,
		database.Username,
		database.Password,
		database.Host+":"+strconv.Itoa(database.Port),
		database.DBName,
		database.ReadTimeout,
		database.WriteTimeout,
	)

	log.Debug().Msgf("main: database connection string: %s", connectionStr)

	var db *gorm.DB
	var err error

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	err = backoff.Retry(func() error {
		db, err = gorm.Open(sqlDriver, connectionStr)
		if err != nil {
			log.Error().Msgf("main: %s open failed: %v", sqlDriver, err)
			return err
		}
		err = db.DB().Ping()
		if err != nil {
			log.Error().Msgf("main: %s ping error: %v", sqlDriver, err)
			return err
		}
		return nil
	}, bo)

	if err != nil {
		log.Error().Msgf("main: mysql connect err: %s", err.Error())
		return nil, err
	}

	return db, nil
}

// OpenPostgres ...
func OpenPostgres(rdbms *Rdbms) (*gorm.DB, error) {
	// sqlDriver = "postgres"

	// connectionString = fmt.Sprintf(`user=%s password=%s host=%s port=%d dbname=%s sslmode=disable `, database.Username, database.Password, database.Host, database.Port, database.DBName)
	// if strings.TrimSpace(database.SearchPath) != "" {
	// 	connectionString = fmt.Sprintf("%s search_path=%s", connectionString, database.SearchPath)
	// }

	return nil, nil
}
