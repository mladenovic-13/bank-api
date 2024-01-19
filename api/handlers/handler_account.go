package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/internal/database"
	"github.com/mladenovic-13/bank-api/models"
)

func (ctx *HandlerContext) HandleGetAccounts(
	w http.ResponseWriter,
	r *http.Request,
	user models.User,
) {
	accounts, err := ctx.DB.GetAccounts(r.Context(), user.ID)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Failed to get accounts")
		return
	}

	api.RespondWithJSON(w, http.StatusOK, models.ToAccounts(accounts))
}

func (ctx *HandlerContext) HandleCreateAccount(
	w http.ResponseWriter,
	r *http.Request,
	user models.User,
) {
	createAccountRequest := new(api.CreateAccountRequest)

	err := json.NewDecoder(r.Body).Decode(createAccountRequest)
	defer r.Body.Close()

	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	newAccount := *models.NewAccount(
		createAccountRequest.Name,
		createAccountRequest.Currency,
		user.ID,
	)

	account, err := ctx.DB.CreateAccount(
		r.Context(),
		database.CreateAccountParams(newAccount),
	)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Failed to create account")
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, models.ToAccount(account))
}
