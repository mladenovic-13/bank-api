package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

	account, err := ctx.DB.GetAccountByID(
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

	account, err := ctx.DB.DeleteAccount(
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

func (ctx *HandlerContext) HandleDeposit(w http.ResponseWriter, r *http.Request, user models.User) {
	if !user.IsAdmin {
		api.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountNumber := chi.URLParam(r, "number")

	if accountNumber == "" {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid account number")
		return
	}

	accountNumberUUID, err := uuid.Parse(accountNumber)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Account number is not valid UUID")
		return
	}

	depositRequest := new(api.DepositRequest)

	err = json.NewDecoder(r.Body).Decode(depositRequest)
	defer r.Body.Close()

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	// TODO: Create transaction

	account, err := ctx.DB.GetAccountByNumber(r.Context(), accountNumberUUID)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Failed to get account")
		return
	}

	if account.Currency != depositRequest.Currency {
		api.RespondWithError(w, http.StatusBadRequest, "Money needs to be in the same currency")
		return
	}

	balanceNumber, err := strconv.ParseFloat(account.Balance, 32)

	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	newBalance := balanceNumber + float64(depositRequest.Amount)

	updatedAccount, err := ctx.DB.UpdateAccountBalance(
		r.Context(),
		database.UpdateAccountBalanceParams{
			ID:        account.ID,
			Balance:   fmt.Sprintf("%0.2f", newBalance),
			UpdatedAt: time.Now(),
		},
	)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Failed to deposit money")
		return
	}

	api.RespondWithJSON(w, http.StatusOK, models.ToAccount(updatedAccount))
}

func (ctx *HandlerContext) HandleWithdraw(w http.ResponseWriter, r *http.Request, user models.User) {
	if !user.IsAdmin {
		api.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountNumber := chi.URLParam(r, "number")

	if accountNumber == "" {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid account number")
		return
	}

	accountNumberUUID, err := uuid.Parse(accountNumber)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Account number is not valid UUID")
		return
	}
	withdrawRequest := new(api.DepositRequest)

	err = json.NewDecoder(r.Body).Decode(withdrawRequest)
	defer r.Body.Close()

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	account, err := ctx.DB.GetAccountByNumber(r.Context(), accountNumberUUID)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Failed to get account")
		return
	}

	balanceNumber, err := strconv.ParseFloat(account.Balance, 32)

	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if balanceNumber >= float64(withdrawRequest.Amount) &&
		account.Currency == withdrawRequest.Currency {
		updatedAccount, err := ctx.DB.UpdateAccountBalance(r.Context(),
			database.UpdateAccountBalanceParams{
				ID:      accountNumberUUID,
				Balance: fmt.Sprintf("%0.2f", balanceNumber-float64(withdrawRequest.Amount)),
			},
		)

		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, "Failed to withdraw money")
			return
		}

		api.RespondWithJSON(w, http.StatusOK, models.ToAccount(updatedAccount))
	}
}

func (ctx *HandlerContext) HandleSend(w http.ResponseWriter, r *http.Request, user models.User) {
	accountNumber := chi.URLParam(r, "number")

	if accountNumber == "" {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid account number")
		return
	}

	// accountNumberUUID, err := uuid.Parse(accountNumber)

	// if err != nil {
	// 	api.RespondWithError(w, http.StatusBadRequest, "Account number is not valid UUID")
	// 	return
	// }

	// sendRequest := new(api.SendRequest)

	// err = json.NewDecoder(r.Body).Decode(sendRequest)
	// defer r.Body.Close()

	// if err != nil {
	// 	api.RespondWithError(w, http.StatusBadRequest, "Invalid request data")
	// 	return
	// }

	// account, err := ctx.DB.GetAccountByNumber(r.Context(), accountNumberUUID)

	// if err != nil {
	// 	api.RespondWithError(w, http.StatusBadRequest, "Failed to get account")
	// 	return
	// }
	api.RespondWithError(w, http.StatusOK, "Not implemented yet")
}
