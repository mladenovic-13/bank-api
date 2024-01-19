package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/internal/database"
	"github.com/mladenovic-13/bank-api/models"
)

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

	if account.Balance >= int32(withdrawRequest.Amount) &&
		account.Currency == withdrawRequest.Currency {
		updatedAccount, err := ctx.DB.UpdateAccountBalance(r.Context(),
			database.UpdateAccountBalanceParams{
				ID:      accountNumberUUID,
				Balance: account.Balance - int32(withdrawRequest.Amount),
			},
		)

		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, "Failed to withdraw money")
			return
		}

		api.RespondWithJSON(w, http.StatusOK, models.ToAccount(updatedAccount))
	}
}
