package users

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/egsam98/voting/gosuslugi/db/repositories"
	"github.com/egsam98/voting/gosuslugi/db/repositories/usersdb"
)

// Service manipulates users
type Service struct{}

// CreateMany creates many new users wrapped by database tx
func (*Service) CreateMany(ctx context.Context, r *repositories.Repositories, newUsers []usersdb.CreateParams) error {
	return r.ExecuteTx(func(r *repositories.Repositories) error {
		for _, params := range newUsers {
			if err := r.Users.Create(ctx, params); err != nil {
				return errors.Wrapf(err, "failed to create user, params: %#v", params)
			}
		}
		return nil
	})
}

// FindByPassport returns user found by specific passport
func (*Service) FindByPassport(ctx context.Context, r *repositories.Repositories, passport string) (*usersdb.User, error) {
	user, err := r.Users.FindByPassport(ctx, passport)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrapf(ErrNotFound, "passport=%s", passport)
		}
		return nil, errors.Wrapf(err, "failed to find user by passport=%s", passport)
	}
	return &user, nil
}
