package queriesdb

import (
	"github.com/pkg/errors"

	"github.com/egsam98/voting/gosuslugi/db/queriesdb/usersdb"
	"github.com/egsam98/voting/gosuslugi/internal/sqalx"
)

type Queries struct {
	dbtx  sqalx.Node
	Users usersdb.Querier
}

func New(db sqalx.Node) *Queries {
	return &Queries{
		dbtx:  db,
		Users: usersdb.New(db),
	}
}

// ExecuteTx runs function with new database tx
func (q *Queries) ExecuteTx(f func(*Queries) error) error {
	tx, err := q.dbtx.Beginx()
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
