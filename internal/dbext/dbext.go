package dbext

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/jmoiron/sqlx"
	sqldblogger "github.com/simukti/sqldb-logger"

	"github.com/egsam98/voting/gosuslugi/internal/sqalx"
)

// Config is the required properties to use the database.
type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
	CertPath   string
	Logger     sqldblogger.Logger
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
	// MaxIdleConns, then MaxIdleConns will be reduced to match the new
	// MaxOpenConns limit.
	//
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	MaxOpenConns int
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
	// MaxIdleConns, then MaxIdleConns will be reduced to match the new
	// MaxOpenConns limit.
	//
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	MaxIdleConns int
}

// URL returns database config in URL presentation
func (c Config) URL() *url.URL {
	sslMode := "verify-full"
	if c.DisableTLS {
		sslMode = "disable"
	}
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	if !c.DisableTLS {
		q.Set("sslrootcert", c.CertPath)
	}
	q.Set("timezone", "utc")

	return &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     c.Host,
		Path:     c.Name,
		RawQuery: q.Encode(),
	}
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config) (sqalx.Node, error) {
	dsn := cfg.URL().String()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if cfg.Logger != nil {
		db = sqldblogger.OpenDriver(
			dsn,
			db.Driver(),
			cfg.Logger,
			sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
			sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
			sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug),
		)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)

	if cfg.MaxIdleConns != 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	sqlxDb := sqlx.NewDb(db, "postgres")
	node, err := sqalx.New(sqlxDb)
	if err != nil {
		return nil, err
	}

	return node, StatusCheck(context.Background(), node)
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db sqalx.Node) error {
	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	_, err := db.ExecContext(ctx, `SELECT true`)
	return err
}
