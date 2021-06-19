package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/egsam98/voting/gosuslugi/db/repositories"
	"github.com/egsam98/voting/gosuslugi/db/repositories/usersdb"
	"github.com/egsam98/voting/gosuslugi/internal/dbext"
	"github.com/egsam98/voting/gosuslugi/internal/randext"
	"github.com/egsam98/voting/gosuslugi/services/users"
)

var envs struct {
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
	Count int `envconfig:"COUNT" default:"100"`
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
		return errors.Wrap(err, "failed to parse ENVs to struct")
	}

	dbCfg := dbext.Config{
		User:         envs.DB.User,
		Password:     envs.DB.Password,
		Host:         envs.DB.Host,
		Name:         envs.DB.Name,
		DisableTLS:   envs.DB.DisableTLS,
		Logger:       nil,
		MaxOpenConns: 0,
		MaxIdleConns: 0,
	}
	if envs.DB.Log {
		dbCfg.Logger = dbext.NewZeroLogger(log.Logger)
	}

	db, err := dbext.Open(dbCfg)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to PostgreSQL, config: %#v", dbCfg)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Stack().Err(err).Msg("main: Failed to close db connection")
		}
	}()

	return seed(&users.Service{}, repositories.New(db))
}

func seed(users *users.Service, r *repositories.Repositories) error {
	randUsers := make([]usersdb.CreateParams, envs.Count)
	for i := range randUsers {
		randUsers[i] = usersdb.CreateParams{
			Passport:  gofakeit.Numerify("########"),
			Fullname:  gofakeit.Name(),
			BirthDate: gofakeit.Date(),
			DeathDate: sql.NullTime{Time: gofakeit.Date(), Valid: randext.Bool()},
		}
	}

	if err := users.CreateMany(context.Background(), r, randUsers); err != nil {
		return err
	}

	log.Info().Int("count", envs.Count).Msg("main: Users added")
	return nil
}
