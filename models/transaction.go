package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/internal/database"
)

type Transaction struct {
	ID              uuid.UUID                `json:"id"`
	SenderNumber    uuid.UUID                `json:"senderNumber"`
	ReceiverNumber  uuid.UUID                `json:"receiverNumber"`
	Amount          string                   `json:"amount"`
	Currency        database.Currency        `json:"currency"`
	TransactionType database.TransactionType `json:"transactionType"`
	CreatedAt       time.Time                `json:"createdAt"`
}

func ToTransaction(t database.Transaction) *Transaction {
	return &Transaction{
		ID:              t.ID,
		SenderNumber:    t.SenderNumber,
		ReceiverNumber:  t.ReceiverNumber,
		Amount:          t.Amount,
		Currency:        t.Currency,
		TransactionType: t.TransactionType,
		CreatedAt:       t.CreatedAt,
	}
}

func (t *Transaction) ToString() string {
	return fmt.Sprintf(
		"%s\n From: %s\n To: %s\n Amount: %s, Currency: %s\n",
		t.TransactionType,
		t.SenderNumber.String(),
		t.ReceiverNumber.String(),
		t.Amount,
		t.Currency,
	)
}

func ToTransactions(t []database.Transaction) []*Transaction {
	accounts := []*Transaction{}
	for _, dbTransaction := range t {
		accounts = append(accounts, ToTransaction(dbTransaction))
	}

	return accounts
}
