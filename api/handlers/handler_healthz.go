package handlers

import (
	"net/http"

	"github.com/mladenovic-13/bank-api/api"
)

func (ctx *HandlerContext) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	api.RespondWithJSON(w, http.StatusOK, "ok")
}
