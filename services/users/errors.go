package users

import (
	"github.com/egsam98/voting/gosuslugi/internal/web"
)

var (
	ErrInvalidInput = &web.ClientError{
		Err:  "invalid input",
		Code: 1,
	}
	ErrNotFound = &web.ClientError{
		Err:  "user is not found",
		Code: 2,
	}
)
