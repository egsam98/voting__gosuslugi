package requests

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

var _ render.Binder = (*FindByPassport)(nil)

type FindByPassport struct {
	Passport string `json:"passport"`
}

func (v *FindByPassport) Bind(*http.Request) error {
	if v.Passport == "" {
		return errors.New("\"passport\" must be non empty")
	}
	return nil
}
