package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/internal/database"
	"github.com/mladenovic-13/bank-api/models"
	"github.com/mladenovic-13/bank-api/utils"
)

func (ctx *HandlerContext) HandleGetAccounts(
	w http.ResponseWriter,
	r *http.Request,
	user models.User,
) {
	accounts, err := ctx.Queries.GetAccounts(r.Context(), user.ID)

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

	account, err := ctx.Queries.CreateAccount(
		r.Context(),
		database.CreateAccountParams(newAccount),
	)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Failed to create account")
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, models.ToAccount(account))
}

func (ctx *HandlerContext) HandleGetAccount(w http.ResponseWriter, r *http.Request, user models.User) {
	accountID := chi.URLParam(r, "id")

	if accountID == "" {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	accountUUID, err := uuid.Parse(accountID)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	account, err := ctx.Queries.GetAccountByID(
		r.Context(),
		database.GetAccountByIDParams{
			ID:     accountUUID,
			UserID: user.ID,
		},
	)

	if err != nil {
		api.RespondWithError(w, http.StatusNotFound, "Account not found")
		return
	}

	api.RespondWithJSON(w, http.StatusOK, models.ToAccount(account))
}

func (ctx *HandlerContext) HandleDeleteAccount(
	w http.ResponseWriter,
	r *http.Request,
	user models.User,
) {
	accountID := chi.URLParam(r, "id")

	if accountID == "" {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	accountUUID, err := uuid.Parse(accountID)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	account, err := ctx.Queries.DeleteAccount(
		r.Context(),
		database.DeleteAccountParams{
			ID:     accountUUID,
			UserID: user.ID,
		},
	)

	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Failed to delete account")
	}

	api.RespondWithJSON(w, http.StatusOK, models.ToAccount(account))
}

func (ctx *HandlerContext) HandleSend(w http.ResponseWriter, r *http.Request, user models.User) {
	senderNumber := chi.URLParam(r, "number")

	if senderNumber == "" {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid account number")
		return
	}

	senderNumberUUID, err := uuid.Parse(senderNumber)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid account number")
		return
	}

	sendRequest := new(api.SendRequest)

	err = json.NewDecoder(r.Body).Decode(sendRequest)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid request data")
		return
	}
	defer r.Body.Close()

	senderAccount, err := ctx.Queries.GetAccountByNumber(r.Context(), senderNumberUUID)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Internal server error")
		return
	}
	receiverAccount, err := ctx.Queries.GetAccountByNumber(r.Context(), sendRequest.ToAccountNumber)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Internal server error")
		return
	}

	if senderAccount.Currency != sendRequest.Currency {
		api.RespondWithError(w, http.StatusBadRequest, "Can not send this currency from account")
		return
	}

	if receiverAccount.Currency != sendRequest.Currency {
		api.RespondWithError(w, http.StatusBadRequest, "Can not receiver this currency to account")
		return
	}

	senderBalance, err := strconv.ParseFloat(senderAccount.Balance, 32)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Internal server error")
		return
	}

	if senderBalance < sendRequest.Amount {
		api.RespondWithError(w, http.StatusBadRequest, "You don't have enough funds in your account")
		return
	}

	transactionResult, err := utils.ExecSendTransaction(
		r.Context(),
		ctx.DB,
		ctx.Queries,
		&senderAccount,
		&receiverAccount,
		sendRequest,
	)

	if err != nil {
		log.Printf("Transaction: %+v\n", err)
		api.RespondWithError(w, http.StatusInternalServerError, "Failed to execute transaction")
		return
	}

	api.RespondWithJSON(w, http.StatusOK, transactionResult)
}
