package repository

import (
	"apigateway/pkg/database"

	"github.com/jinzhu/gorm"
	"go.uber.org/fx"
)

type repository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
	// microservice
}

// NewRepository ...
func NewRepository(conn *database.RdbmsConn) IRepository {
	return &repository{
		readDB:  conn.ReadDB,
		writeDB: conn.WriteDB,
	}
}

// IRepository ...
type IRepository interface {
	BookRepository
}

// Module Export repository module
var Module = fx.Options(
	fx.Provide(NewRepository),
)
