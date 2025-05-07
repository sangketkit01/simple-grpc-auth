-- name: CreateUser :one
INSERT INTO users 
    (username, hashed_password, full_name, email)
VALUES
    ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT username, full_name, email, created_at FROM users 
WHERE username = $1;

-- name: LoginUser :one
SELECT username, hashed_password
FROM users WHERE username = $1;