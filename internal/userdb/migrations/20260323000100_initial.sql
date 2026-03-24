-- +goose Up
-- +goose StatementBegin
CREATE TABLE check_results (
    id TEXT PRIMARY KEY,
    agent_id TEXT NOT NULL,
    target_id TEXT NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('up', 'down', 'degraded')),
    response_time INTEGER,
    details TEXT,
    checked_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_check_results_target_checked_at ON check_results(target_id, checked_at DESC);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE incidents (
    id TEXT PRIMARY KEY,
    target_id TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    details TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_incidents_target_id ON incidents(target_id);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER update_incidents_updated_at
AFTER
UPDATE ON incidents FOR EACH ROW BEGIN
UPDATE incidents
SET updated_at = CURRENT_TIMESTAMP
WHERE id = OLD.id;
END;
-- +goose StatementEnd
---
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_incidents_updated_at;
-- +goose StatementEnd
-- +goose StatementBegin
DROP INDEX idx_incidents_target_id;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE incidents;
-- +goose StatementEnd
-- +goose StatementBegin
DROP INDEX idx_check_results_target_checked_at;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE check_results;
-- +goose StatementEnd