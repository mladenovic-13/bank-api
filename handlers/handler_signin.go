package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/internal/database"
	"github.com/mladenovic-13/bank-api/models"
	"github.com/mladenovic-13/bank-api/utils"
)

func (ctx *RouterCtx) HandleSignin(w http.ResponseWriter, r *http.Request) {
	credentials := new(SigninRequest)

	err := json.NewDecoder(r.Body).Decode(credentials)
	defer r.Body.Close()

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Missing credentials")
		return
	}

	user, err := ctx.DB.GetUserByUsername(r.Context(), credentials.Username)

	if err == nil && user.Username == credentials.Username {
		RespondWithError(w, http.StatusBadRequest, "User already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(credentials.Password)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
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
		RespondWithError(w, http.StatusBadRequest, "Failed to create user")
		return
	}

	RespondWithJSON(w, http.StatusCreated, models.ToUser(newUser))
}
