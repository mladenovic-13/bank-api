package handlers

import (
	"net/http"
	"time"

	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/models"
)

func (ctx *HandlerContext) HandleLogout(w http.ResponseWriter, r *http.Request, user models.User) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})

	api.RespondWithJSON(
		w, http.StatusOK,
		map[string]string{"message": "login success"},
	)
}
