-- name: CreateUser :one
INSERT INTO users (
  username, hashed_password, full_name, email, total_expenses
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1;

-- name: GetEmail :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE username = $1
FOR NO KEY UPDATE;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET total_expenses = $2
WHERE username = $1
RETURNING *;

-- name: ResetPassword :exec
UPDATE users
SET hashed_password = $2
WHERE username = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;