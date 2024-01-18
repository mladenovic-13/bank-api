package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mladenovic-13/bank-api/models"
	"github.com/mladenovic-13/bank-api/utils"
)

func (ctx *RouterCtx) HandleLogin(w http.ResponseWriter, r *http.Request) {
	credentials := new(LoginRequest)

	err := json.NewDecoder(r.Body).Decode(credentials)
	defer r.Body.Close()

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Missing credentials")
		return
	}

	user, err := ctx.DB.GetUserByUsername(r.Context(), credentials.Username)

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "User does not exist")
		return
	}

	if utils.CheckPasswordHash(credentials.Password, user.Password) {
		token, err := utils.CreateJWT(models.ToUser(user))

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(24 * time.Hour),
		})

		RespondWithJSON(
			w, http.StatusOK,
			map[string]string{"message": "login success"},
		)
	} else {
		RespondWithError(w, http.StatusBadRequest, "Wrong password")
		return
	}
}
