package rest

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/egsam98/voting/gosuslugi/internal/dbext"
	"github.com/egsam98/voting/gosuslugi/internal/sqalx"
)

type healthController struct {
	db sqalx.Node
}

func newHealthController(db sqalx.Node) *healthController {
	return &healthController{db: db}
}

type status struct {
	Status string `json:"status"`
}

func fail(w http.ResponseWriter, r *http.Request, serviceName string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	render.JSON(w, r, status{Status: serviceName + ": " + err.Error()})
}

func ok(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, status{Status: "ok"})
}

func (hc *healthController) Readiness(w http.ResponseWriter, r *http.Request) {
	if err := dbext.StatusCheck(r.Context(), hc.db); err != nil {
		fail(w, r, "postgresql", err)
		return
	}

	ok(w, r)
}
