package models

import (
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Currency string

const (
	RSD Currency = "RSD"
	USD Currency = "USD"
	EUR Currency = "USD"
)

type Account struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Balance   float32   `gorm:"type:decimal(10,2);not null" json:"balance"`
	Currency  Currency  `gorm:"not null" json:"currency"`
	UserID    uuid.UUID `gorm:"not null" json:"userId"`
	CreateAt  time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt"`
}

func NewAccount(
	accountId uuid.UUID,
	name string,
	currency Currency,
	userId uuid.UUID,
) *Account {
	return &Account{
		ID:        accountId,
		Name:      name,
		Balance:   0,
		Currency:  currency,
		UserID:    userId,
		CreateAt:  time.Now(),
		UpdatedAt: time.Now(),
	}
}
