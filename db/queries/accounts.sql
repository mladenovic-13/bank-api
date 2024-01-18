-- name: CreateAccount :one
INSERT INTO accounts(
  id, name, number, balance, currency, user_id, created_at, updated_at
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAccounts :many
SELECT * FROM accounts WHERE user_id=$1; 

-- name: GetAccountByID :one
SELECT * FROM accounts WHERE user_id=$1 AND id=$2; 

-- name: DeleteAccount :one
DELETE FROM accounts WHERE id=$1 AND  user_id=$2
RETURNING *;
