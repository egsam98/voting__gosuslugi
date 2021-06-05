package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/egsam98/voting/gosuslugi/handlers/rest"
)

const serviceName = "Gosuslugi"

var envs struct {
	Web struct {
		Addr            string        `envconfig:"WEB_ADDR"`
		ShutdownTimeout time.Duration `envconfig:"WEB_SHUTDOWN_TIMEOUT" default:"5s"`
	}
}

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := run(); err != nil {
		log.Fatal().Stack().Err(err).Msg("main: Fatal error")
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Warn().Err(err).Msg("main: Read ENVs from .env file")
	}
	if err := envconfig.Process("", &envs); err != nil {
		return err
	}

	http.HandleFunc("/api/validate", rest.ValidateVote)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)

	srv := http.Server{
		Addr:    envs.Web.Addr,
		Handler: http.DefaultServeMux,
	}

	go func() {
		log.Info().Msgf("main: %q REST service is listening on %q", serviceName, envs.Web.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Stack().Err(err).Msg("main: Failed to start HTTP server")
		}
	}()

	sig := <-sigint

	ctx, cancel := context.WithTimeout(context.Background(), envs.Web.ShutdownTimeout)
	defer cancel()

	log.Info().Msg("main: Shutdown server")
	if err := srv.Shutdown(ctx); err != nil {
		return errors.Wrapf(err, "failed to shutdown server")
	}

	log.Info().Msgf("main: Terminated via signal %q", sig)
	return nil
}
