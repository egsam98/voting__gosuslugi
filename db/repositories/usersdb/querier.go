package usersdb

import (
	"context"
)

type Querier interface {
	Create(ctx context.Context, arg CreateParams) error
	FindByID(ctx context.Context, id int64) (User, error)
}

var _ Querier = (*Queries)(nil)
