package rest

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/egsam98/voting/gosuslugi/cmd/server/handlers/rest/requests/users"
	"github.com/egsam98/voting/gosuslugi/cmd/server/handlers/rest/responses/users"
	"github.com/egsam98/voting/gosuslugi/db/queriesdb"
	"github.com/egsam98/voting/gosuslugi/internal/web"
	"github.com/egsam98/voting/gosuslugi/services/users"
)

type usersController struct {
	users *users.Service
	q     *queriesdb.Queries
}

func newUsersController(users *users.Service, q *queriesdb.Queries) *usersController {
	return &usersController{users: users, q: q}
}

func (uc *usersController) FindByPassport(w http.ResponseWriter, r *http.Request) {
	var req requests.FindByPassport
	if err := render.Bind(r, &req); err != nil {
		web.RespondError(w, r, web.WrapWithError(users.ErrInvalidInput, err))
		return
	}

	user, err := uc.users.FindByPassport(r.Context(), uc.q, req.Passport)
	if err != nil {
		web.RespondError(w, r, err)
		return
	}

	render.JSON(w, r, responses.NewUser(*user))
}
