// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: transactions.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions(id, sender_number, receiver_number, amount, currency, transaction_type, created_at)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING id, sender_number, receiver_number, amount, currency, transaction_type, created_at
`

type CreateTransactionParams struct {
	ID              uuid.UUID
	SenderNumber    uuid.UUID
	ReceiverNumber  uuid.UUID
	Amount          string
	Currency        Currency
	TransactionType TransactionType
	CreatedAt       time.Time
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.ID,
		arg.SenderNumber,
		arg.ReceiverNumber,
		arg.Amount,
		arg.Currency,
		arg.TransactionType,
		arg.CreatedAt,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.SenderNumber,
		&i.ReceiverNumber,
		&i.Amount,
		&i.Currency,
		&i.TransactionType,
		&i.CreatedAt,
	)
	return i, err
}
