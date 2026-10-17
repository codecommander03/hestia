-- name: CreateUser :one
INSERT INTO users (id, name, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND status <> 'DELETED';

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND status <> 'DELETED';