package api

import (
	"net/http"

	"github.com/mladenovic-13/bank-api/internal/database"
	"github.com/mladenovic-13/bank-api/models"
)

type ProtectedHandler func(
	http.ResponseWriter,
	*http.Request,
	models.User,
)

type SigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateAccountRequest struct {
	Name     string            `json:"name"`
	Currency database.Currency `json:"currency"`
}

type DepositRequest struct {
	Amount   int               `json:"amount"`
	Currency database.Currency `json:"currency"`
}
