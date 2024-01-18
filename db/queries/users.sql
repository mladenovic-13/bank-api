-- name: CreateUser :one
INSERT INTO users(id, username, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id=$1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username=$1;


-- name: DeleteUserByID :one
DELETE FROM users WHERE id=$1
RETURNING *;