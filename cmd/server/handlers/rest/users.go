package rest

import (
	"net/http"

	"github.com/go-chi/render"

	users2 "github.com/egsam98/voting/gosuslugi/cmd/server/handlers/rest/requests/users"
	responses "github.com/egsam98/voting/gosuslugi/cmd/server/handlers/rest/responses/users"
	"github.com/egsam98/voting/gosuslugi/db/repositories"
	"github.com/egsam98/voting/gosuslugi/internal/web"
	"github.com/egsam98/voting/gosuslugi/services/users"
)

type usersController struct {
	users *users.Service
	r     *repositories.Repositories
}

func newUsersController(users *users.Service, r *repositories.Repositories) *usersController {
	return &usersController{users: users, r: r}
}

func (uc *usersController) FindByPassport(w http.ResponseWriter, r *http.Request) {
	var req users2.FindByPassport
	if err := render.Bind(r, &req); err != nil {
		web.RespondError(w, r, web.WrapWithError(users.ErrInvalidInput, err))
		return
	}

	user, err := uc.users.FindByPassport(r.Context(), uc.r, req.Passport)
	if err != nil {
		web.RespondError(w, r, err)
		return
	}

	render.JSON(w, r, responses.NewUser(*user))
}