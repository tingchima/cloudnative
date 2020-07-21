package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"apigateway/pkg/config"
	"apigateway/pkg/repository"
	pkgHTTP "apigateway/pkg/router/http"
	"apigateway/pkg/service"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// ServerCmd ...
var ServerCmd = &cobra.Command{
	Use: "server",
	Run: run,
}

func run(command *cobra.Command, args []string) {
	defer CmdRecover()
	srv := &http.Server{}
	exitCode := 0

	// fx injection
	app := fx.New(
		config.Module,
		repository.Module,
		service.Module,
		pkgHTTP.Module,
		fx.Populate(&srv),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(exitCode)
		return
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	// The context is used to inform the server it has 30 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Msgf("Server forced to shutdown: %s", err.Error())
	}

	os.Exit(exitCode)
}
