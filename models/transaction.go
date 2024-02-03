package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeTRANSFER   TransactionType = "TRANSFER"
	TransactionTypePAYMENT    TransactionType = "PAYMENT"
	TransactionTypeWITHDRAWAL TransactionType = "WITHDRAWAL"
	TransactionTypeDEPOSIT    TransactionType = "DEPOSIT"
)

type Transaction struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	Sender          uuid.UUID       `gorm:"type:uuid;not null" json:"sender"`
	Receiver        uuid.UUID       `gorm:"type:uuid;not null" json:"receiver"`
	Amount          float32         `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency        Currency        `gorm:"not null" json:"currency"`
	TransactionType TransactionType `gorm:"not null" json:"transactionType"`
	CreatedAt       time.Time       `gorm:"not null" json:"createdAt"`
}

func NewTransaction(
	sender, receiver uuid.UUID,
	amount float32,
	currency Currency,
	transactionType TransactionType,
) *Transaction {
	return &Transaction{
		ID:              uuid.New(),
		Sender:          sender,
		Receiver:        receiver,
		Amount:          amount,
		Currency:        currency,
		TransactionType: transactionType,
		CreatedAt:       time.Now(),
	}
}
