-- name: CreateTransaction :one
INSERT INTO transactions(id, sender_number, receiver_number, amount, currency, transaction_type, created_at)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

