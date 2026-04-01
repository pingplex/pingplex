-- +goose Up
-- +goose StatementBegin
CREATE TABLE targets (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    description TEXT,
    visibility TEXT NOT NULL CHECK (visibility IN ('private', 'public', 'internal')),
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'deleted')),
    check_params TEXT NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (user_id, name)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_targets_owner_visibility ON targets(user_id, visibility);
-- +goose StatementEnd
---
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_targets_owner_visibility;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS targets;
-- +goose StatementEnd