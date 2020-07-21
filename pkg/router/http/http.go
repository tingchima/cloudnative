package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

// Module Export http server module
var Module = fx.Options(
	fx.Provide(
		NewHandler,
		NewServer,
	),
	fx.Invoke(RunServer),
)

// Config the structure for HTTP
type Config struct {
	Mode           string        `json:"mode"`
	Address        string        `json:"address"`
	AppID          string        `yaml:"app_id" mapstructure:"app_id"`
	ReadTimeout    time.Duration `json:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout"`
	MaxHeaderBytes int           `json:"max_header_bytes"`
}

// NewServer ...
func NewServer(cfg *Config) *http.Server {
	router := gin.Default()

	// use middleware

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	// g.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	// g.Use(gin.Recovery())

	RegisteRouter(router)

	// create server to run
	return &http.Server{
		Addr:           cfg.Address,
		Handler:        router,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}
}

// RunServer ...
func RunServer(cfg *Config, srv *http.Server) error {

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("listen: %s\n", err)
		}
	}()

	return nil
}
