package usersdb

import (
	"context"
)

type Querier interface {
	Create(ctx context.Context, arg CreateParams) error
	FindByPassport(ctx context.Context, passport string) (User, error)
}

var _ Querier = (*Queries)(nil)
