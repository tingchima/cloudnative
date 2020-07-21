package config

import (
	"apigateway/pkg/database"
	"apigateway/pkg/router/http"

	"fmt"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Config the structure for application
type Config struct {
	fx.Out

	HTTP  *http.Config
	MySQL *database.RdbmsConfig
}

// ProvideConfig ...
func ProvideConfig() Config {
	cfg, err := CreateConfig()
	if err != nil {
		log.Fatal().Msg("Error create configuration failed")
	}

	return cfg
}

// CreateConfig 讀取App 啟動程式設定檔
func CreateConfig(path ...string) (Config, error) {
	viper.AutomaticEnv()

	configPath := resolvePath(path)

	configName := viper.GetString("CONFIG_NAME")
	if configName == "" {
		configName = "app"
	}

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Msgf("error reading config file, %s", err)
		return cfg, err
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Error().Msgf("unable to decode into struct, %v", err)
		return cfg, err
	}

	return cfg, nil
}

func resolvePath(paths []string) string {
	switch len(paths) {
	case 0:
		configPath := viper.GetString("CONFIGPATH")

		_, b, _, _ := runtime.Caller(0)
		if configPath == "" {
			configPath = fmt.Sprintf(`%s/../../env/`, filepath.Dir(b))
		}
		return configPath
	case 1:
		return paths[0]
	default:
		panic("too many parameters")
	}
}

// Module provide related configuration
var Module = fx.Options(
	fx.Provide(
		ProvideConfig,
		database.InitRdbms,
	),
)
