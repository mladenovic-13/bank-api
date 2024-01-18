package handlers

import (
	"net/http"

	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/models"
)

func (ctx *HandlerContext) HandleGetAccounts(
	w http.ResponseWriter,
	r *http.Request,
	user models.User,
) {
	api.RespondWithJSON(w, http.StatusOK, user)
}
