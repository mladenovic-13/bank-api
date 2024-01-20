package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/utils"
)

func (ctx *HandlerContext) HandleDeposit(w http.ResponseWriter, r *http.Request) {
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

	account, err := ctx.Queries.GetAccountByNumber(r.Context(), accountNumberUUID)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Failed to get account")
		return
	}

	if account.Currency != depositRequest.Currency {
		api.RespondWithError(w, http.StatusBadRequest, "Money needs to be in the same currency")
		return
	}

	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	transactionResult, err := utils.ExecDepositTransaction(
		r.Context(),
		ctx.DB,
		ctx.Queries,
		&account,
		depositRequest,
	)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Failed to deposit money")
		return
	}

	api.RespondWithJSON(w, http.StatusOK, transactionResult)
}

func (ctx *HandlerContext) HandleWithdraw(w http.ResponseWriter, r *http.Request) {
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

	account, err := ctx.Queries.GetAccountByNumber(r.Context(), accountNumberUUID)

	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Failed to get account")
		return
	}

	balanceNumber, err := strconv.ParseFloat(account.Balance, 32)

	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if balanceNumber < withdrawRequest.Amount {
		api.RespondWithError(w, http.StatusBadRequest, "Do not have enough founds")
		return
	}

	transactionResult, err := utils.ExecWithdrawTransaction(
		r.Context(),
		ctx.DB,
		ctx.Queries,
		&account,
		withdrawRequest,
	)

	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Failed to execute transaction")
		return
	}

	api.RespondWithJSON(w, http.StatusOK, transactionResult)
}
