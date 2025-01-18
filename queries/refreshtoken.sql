-- name: CreateRefreshToken :execresult
INSERT INTO refresh_tokens (user_id, hashed_token, expires_at)
VALUES (?, ?, ?);
-- name: GetRefreshTokenByID :one
SELECT *
FROM refresh_tokens
WHERE id = ?;
-- name: DeleteUserTokens :execresult
Delete FROM refresh_tokens
WHERE user_id = ?;
-- name: GetRefreshTokenByUserIDAndToken :one
SELECT *
FROM refresh_tokens
WHERE user_id = ?
    AND hashed_token = ?;
