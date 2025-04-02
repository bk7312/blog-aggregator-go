-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE users.name = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE users.id = $1;

-- name: ResetUsers :exec
DELETE from users;

-- name: GetUsers :many
SELECT * FROM users;