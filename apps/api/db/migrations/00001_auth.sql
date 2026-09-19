-- +goose Up
CREATE TABLE users (
    id text PRIMARY KEY,
    email text NOT NULL UNIQUE,
    first_name text NOT NULL,
    last_name text NOT NULL,
    password_hash text NOT NULL,
    role text NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CHECK (char_length(email) BETWEEN 3 AND 254),
    CHECK (char_length(first_name) BETWEEN 1 AND 100),
    CHECK (char_length(last_name) BETWEEN 1 AND 100)
);

CREATE TABLE sessions (
    token_hash bytea PRIMARY KEY CHECK (octet_length(token_hash) = 32),
    user_id text NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL,
    last_seen_at timestamptz NOT NULL,
    idle_expires_at timestamptz NOT NULL,
    absolute_expires_at timestamptz NOT NULL,
    CHECK (last_seen_at >= created_at),
    CHECK (idle_expires_at > last_seen_at),
    CHECK (absolute_expires_at > created_at)
);

CREATE INDEX sessions_user_id_idx ON sessions(user_id);
CREATE INDEX sessions_expiry_idx ON sessions(idle_expires_at, absolute_expires_at);

-- +goose Down
DROP TABLE sessions;
DROP TABLE users;
