package rest

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/egsam98/voting/gosuslugi/db/queriesdb"
	"github.com/egsam98/voting/gosuslugi/services/users"
)

func API(users *users.Service, q *queriesdb.Queries) http.Handler {
	uc := newUsersController(users, q)

	mux := chi.NewMux()

	api := chi.NewRouter()
	mux.Mount("/api", api)
	api.Use(
		middleware.Recoverer,
		middleware.RequestLogger(&middleware.DefaultLogFormatter{
			Logger: &log.Logger,
		}),
	)

	api.Post("/users/passport", uc.FindByPassport)

	return mux
}
