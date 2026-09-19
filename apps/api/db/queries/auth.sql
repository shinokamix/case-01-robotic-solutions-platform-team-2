-- name: CreateUser :one
INSERT INTO users (id, email, first_name, last_name, password_hash, role)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, email, first_name, last_name, password_hash, role, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, email, first_name, last_name, password_hash, role, created_at, updated_at
FROM users
WHERE email = $1;

-- name: CreateSession :exec
INSERT INTO sessions (
    token_hash,
    user_id,
    created_at,
    last_seen_at,
    idle_expires_at,
    absolute_expires_at
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetSessionUser :one
SELECT
    s.token_hash,
    s.created_at AS session_created_at,
    s.last_seen_at,
    s.idle_expires_at,
    s.absolute_expires_at,
    u.id,
    u.email,
    u.first_name,
    u.last_name,
    u.password_hash,
    u.role,
    u.created_at AS user_created_at,
    u.updated_at AS user_updated_at
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE s.token_hash = $1;

-- name: TouchSession :execrows
UPDATE sessions
SET last_seen_at = GREATEST(last_seen_at, $2),
    idle_expires_at = GREATEST(idle_expires_at, $3)
WHERE token_hash = $1
  AND idle_expires_at > $2
  AND absolute_expires_at > $2;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token_hash = $1;
