package rest

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/egsam98/voting/gosuslugi/db/repositories"
	"github.com/egsam98/voting/gosuslugi/services/users"
)

func API(users *users.Service, r *repositories.Repositories) http.Handler {
	uc := newUsersController(users, r)

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
