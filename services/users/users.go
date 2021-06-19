package users

import (
	"context"

	"github.com/pkg/errors"

	"github.com/egsam98/voting/gosuslugi/db/repositories"
	"github.com/egsam98/voting/gosuslugi/db/repositories/usersdb"
)

type Service struct{}

func (s *Service) CreateMany(ctx context.Context, r *repositories.Repositories, newUsers []usersdb.CreateParams) error {
	return r.ExecuteTx(func(r *repositories.Repositories) error {
		for _, params := range newUsers {
			if err := r.Users.Create(ctx, params); err != nil {
				return errors.Wrapf(err, "failed to create user, params: %#v", params)
			}
		}
		return nil
	})
}
