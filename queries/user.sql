-- name: GetUsers :many
SELECT *
FROM users;
-- name: CreateUser :execresult
INSERT INTO users (email, hashed_password, uuid, created_at)
VALUES (?, ?, ?, now());
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?;
-- name: GetUserById :one
SELECT *
FROM users
WHERE id = ?;
