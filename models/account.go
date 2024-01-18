package models

import (
	"math/rand"
	"time"
)

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	UserID    int       `json:"userId"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewAccount(firstName, lastName string, userID int) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		UserID:    userID,
		Number:    int64(rand.Intn(1000000)),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
