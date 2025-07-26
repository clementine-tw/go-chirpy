-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens
(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES
($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users
WHERE id = (
	SELECT refresh_tokens.user_id
	FROM refresh_tokens
	WHERE refresh_tokens.token = $1
	AND refresh_tokens.revoked_at IS NULL
	AND refresh_tokens.expires_at > NOW()
);

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE token = $1
RETURNING *;
