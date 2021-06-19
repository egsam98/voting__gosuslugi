package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var randStatusCode = map[int]int{
	0: http.StatusOK,
	1: http.StatusBadRequest,
}

func ValidateVote(ctx *gin.Context) {
	//key := rand.Intn(2)
	ctx.Status(randStatusCode[1])
}
