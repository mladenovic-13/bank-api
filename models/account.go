package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/internal/database"
)

type Account struct {
	ID        uuid.UUID         `json:"id"`
	Name      string            `json:"name"`
	Number    uuid.UUID         `json:"number"`
	Balance   string            `json:"balance"`
	Currency  database.Currency `json:"currency"`
	UserID    uuid.UUID         `json:"userId"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

func NewAccount(name string, currency database.Currency, userId uuid.UUID) *Account {
	return &Account{
		ID:        uuid.New(),
		Name:      name,
		Number:    uuid.New(),
		Balance:   "0.00",
		Currency:  currency,
		UserID:    userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func ToAccount(a database.Account) *Account {
	return &Account{
		ID:        a.ID,
		Name:      a.Name,
		Number:    a.Number,
		Balance:   a.Balance,
		Currency:  a.Currency,
		UserID:    a.UserID,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func ToAccounts(a []database.Account) []*Account {
	accounts := []*Account{}
	for _, dbAccount := range a {
		accounts = append(accounts, ToAccount(dbAccount))
	}

	return accounts
}
