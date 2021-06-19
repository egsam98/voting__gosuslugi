package repositories

import (
	"github.com/pkg/errors"

	"github.com/egsam98/voting/gosuslugi/db/repositories/usersdb"
	"github.com/egsam98/voting/gosuslugi/internal/sqalx"
)

type Repositories struct {
	dbtx  sqalx.Node
	Users usersdb.Querier
}

func New(db sqalx.Node) *Repositories {
	return &Repositories{
		dbtx:  db,
		Users: usersdb.New(db),
	}
}

// ExecuteTx runs function with new database tx
func (r *Repositories) ExecuteTx(f func(r *Repositories) error) error {
	tx, err := r.dbtx.Beginx()
	if err != nil {
		return errors.Wrap(err, "failed to begin database tx")
	}

	if err := f(New(tx)); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrapf(err, "tx rollback error: %v", rollbackErr)
		}
		return err
	}
	return errors.Wrap(tx.Commit(), "failed to commit database tx")
}
