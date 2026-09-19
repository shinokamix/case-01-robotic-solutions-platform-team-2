-- +goose Up
CREATE UNIQUE INDEX users_single_admin_idx ON users ((role)) WHERE role = 'admin';

-- +goose Down
DROP INDEX users_single_admin_idx;
