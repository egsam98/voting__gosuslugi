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
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/egsam98/voting/gosuslugi/cmd/server/handlers/rest"
	"github.com/egsam98/voting/gosuslugi/internal/dbext"
	"github.com/egsam98/voting/gosuslugi/services/users"
)

const serviceName = "Gosuslugi"

var envs struct {
	Web struct {
		Addr            string        `envconfig:"WEB_ADDR" default:":3000"`
		ShutdownTimeout time.Duration `envconfig:"WEB_SHUTDOWN_TIMEOUT" default:"5s"`
	}
	DB struct {
		// User is the username for the database
		User string `envconfig:"DB_USER" default:"postgres"`
		// Password is the password for the database
		Password string `envconfig:"DB_PASSWORD" default:"postgres"`
		// Host is the address of database
		Host string `envconfig:"DB_HOST" default:"db"`
		// Name is the database name to connect to
		Name string `envconfig:"DB_NAME" default:"postgres"`
		// DisableTLS is flag indicating to disable TLS
		DisableTLS bool `envconfig:"DB_DISABLE_TLS" default:"true"`
		// Log is flag to enable logging of SQL queries
		Log bool `envconfig:"DB_LOG" default:"true"`
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

	log.Info().
		Interface("envs", envs).
		Msg("main: ENVs")

	dbCfg := dbext.Config{
		User:       envs.DB.User,
		Password:   envs.DB.Password,
		Host:       envs.DB.Host,
		Name:       envs.DB.Name,
		DisableTLS: envs.DB.DisableTLS,
	}
	if envs.DB.Log {
		dbCfg.Logger = dbext.NewZeroLogger(log.Logger)
	}

	db, err := dbext.Open(dbCfg)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to PostgreSQL, config: %#v", dbCfg)
	}

	defer func() {
		log.Info().Msg("main: Closing database")
		if err := db.Close(); err != nil {
			log.Error().Stack().Err(err).Msg("main: Failed to close db connection")
		}
	}()

	srv := http.Server{
		Addr:    envs.Web.Addr,
		Handler: rest.API(&users.Service{}, db),
	}

	apiErr := make(chan error)
	go func() {
		log.Info().Msgf("main: %q REST service is listening on %q", serviceName, envs.Web.Addr)
		apiErr <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)

	select {
	case err := <-apiErr:
		return errors.WithStack(err)
	case sig := <-shutdown:
		log.Info().Msgf("main: Received signal %q", sig)

		ctx, cancel := context.WithTimeout(context.Background(), envs.Web.ShutdownTimeout)
		defer cancel()

		log.Info().Msg("main: Doing server shutdown")
		if err := srv.Shutdown(ctx); err != nil {
			_ = srv.Close()
			return errors.Wrapf(err, "failed to shutdown server")
		}
	}

	return nil
}
