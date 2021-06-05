package rest

import (
	"math/rand"
	"net/http"
)

var randStatusCode = map[int]int{
	0: http.StatusOK,
	1: http.StatusBadRequest,
}

func ValidateVote(w http.ResponseWriter, _ *http.Request) {
	key := rand.Intn(2)
	w.WriteHeader(randStatusCode[key])
}
