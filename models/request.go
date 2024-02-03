package models

import (
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type RequestType string

const (
	OPEN_ACCOUNT  RequestType = "OPEN_ACCOUNT"
	CLOSE_ACCOUNT RequestType = "CLOSE_ACCOUNT"
	DEPOSIT       RequestType = "DEPOSIT"
	WITHDRAW      RequestType = "WITHDRAW"
)

type Request struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID   `gorm:"not null" json:"userId"`
	AccountID   uuid.UUID   `gorm:"not null" json:"accountId"`
	Type        RequestType `gorm:"not null" json:"type"`
	Amount      float32     `gorm:"type:decimal(10,2)" json:"amount"`
	Currency    Currency    `json:"currency"`
	IsProcessed bool        `gorm:"not null" json:"isProcessed"`
	CreateAt    time.Time   `gorm:"not null" json:"createdAt"`
}

type RequestProps struct {
	AccountID uuid.UUID `json:"accountId"`
	Amount    float32   `json:"amount"`
	Currency  Currency  `json:"currency"`
}

func NewRequest(
	reqType RequestType,
	userId uuid.UUID,
	props RequestProps,
) *Request {
	return &Request{
		ID: uuid.New(),
		// BUG
		AccountID:   uuid.New(),
		UserID:      userId,
		Type:        reqType,
		Amount:      props.Amount,
		Currency:    props.Currency,
		IsProcessed: false,
		CreateAt:    time.Now(),
	}
}
