package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/internal/database"
	"github.com/mladenovic-13/bank-api/models"
	"github.com/mladenovic-13/bank-api/utils"
)

func (ctx *HandlerContext) HandleSignin(w http.ResponseWriter, r *http.Request) {
	credentials := new(api.SigninRequest)

	err := json.NewDecoder(r.Body).Decode(credentials)
	defer r.Body.Close()

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Missing credentials")
		return
	}

	user, err := ctx.DB.GetUserByUsername(r.Context(), credentials.Username)

	if err == nil && user.Username == credentials.Username {
		api.RespondWithError(w, http.StatusBadRequest, "User already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(credentials.Password)

	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	newUser, err := ctx.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Username:  credentials.Username,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Failed to create user")
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, models.ToUser(newUser))
}

func (ctx *HandlerContext) HandleLogin(w http.ResponseWriter, r *http.Request) {
	credentials := new(api.LoginRequest)

	err := json.NewDecoder(r.Body).Decode(credentials)
	defer r.Body.Close()

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Missing credentials")
		return
	}

	user, err := ctx.DB.GetUserByUsername(r.Context(), credentials.Username)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "User does not exist")
		return
	}

	if utils.CheckPasswordHash(credentials.Password, user.Password) {
		token, err := utils.CreateJWT(models.ToUser(user))

		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(24 * time.Hour),
		})

		api.RespondWithJSON(
			w, http.StatusOK,
			map[string]string{"message": "log in success"},
		)
	} else {
		api.RespondWithError(w, http.StatusBadRequest, "Wrong password")
		return
	}
}
func (ctx *HandlerContext) HandleLogout(w http.ResponseWriter, r *http.Request, user models.User) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})

	api.RespondWithJSON(
		w, http.StatusOK,
		map[string]string{"message": "log out success"},
	)
}
