-- +goose Up
-- +goose StatementBegin
CREATE TABLE agents (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    visibility TEXT NOT NULL CHECK (visibility IN ('private', 'public')),
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'deleted')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (user_id, name)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_agents_owner_visibility ON agents(user_id, visibility);
-- +goose StatementEnd
---
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_agents_owner_visibility;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS agents;
-- +goose StatementEnd