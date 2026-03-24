package agents_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pingplex/pingplex/internal/agents"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "modernc.org/sqlite"
)

func TestRepositoryCRUD(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupTestDB(t)
	repo := agents.NewRepository(db)

	created, err := repo.Create(ctx, "user-1", agents.CreateAgentInput{
		Name:        "agent-1",
		Description: "private agent",
		Visibility:  agents.VisibilityPrivate,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if created.ID == "" {
		t.Fatalf("expected non-empty id")
	}

	fetched, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}

	if fetched.Name != "agent-1" {
		t.Fatalf("expected name agent-1, got %s", fetched.Name)
	}

	updated, err := repo.Update(ctx, "user-1", created.ID, agents.UpdateAgentInput{
		Name:        "agent-1-renamed",
		Description: "updated",
		Visibility:  agents.VisibilityPublic,
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}

	if updated.Name != "agent-1-renamed" {
		t.Fatalf("expected updated name, got %s", updated.Name)
	}

	list, err := repo.ListByOwner(ctx, "user-1")
	if err != nil {
		t.Fatalf("list by owner failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected list length 1, got %d", len(list))
	}

	if err = repo.Delete(ctx, "user-1", created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	deleted, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get deleted by id failed: %v", err)
	}

	if deleted.Status != agents.StatusDeleted {
		t.Fatalf("expected deleted status, got %s", deleted.Status)
	}
}

func TestRepositoryUpdateWithWrongOwnerReturnsNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupTestDB(t)
	repo := agents.NewRepository(db)

	created, err := repo.Create(ctx, "user-1", agents.CreateAgentInput{
		Name:       "agent-1",
		Visibility: agents.VisibilityPrivate,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, err = repo.Update(ctx, "user-2", created.ID, agents.UpdateAgentInput{
		Name:       "new-name",
		Visibility: agents.VisibilityPublic,
	})
	if !errors.Is(err, agents.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func setupTestDB(t *testing.T) *bun.DB {
	t.Helper()

	dbName := strings.ReplaceAll(t.Name(), "/", "_")
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", dbName)

	sqldb, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	t.Cleanup(func() {
		_ = sqldb.Close()
	})

	db := bun.NewDB(sqldb, sqlitedialect.New())
	t.Cleanup(func() {
		_ = db.Close()
	})

	schema := []string{
		`CREATE TABLE users (
			id TEXT PRIMARY KEY,
			status TEXT NOT NULL DEFAULT 'active',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE agents (
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
		);`,
		`INSERT INTO users (id) VALUES ('user-1'), ('user-2');`,
	}

	for _, query := range schema {
		if _, qErr := db.ExecContext(context.Background(), query); qErr != nil {
			t.Fatalf("failed to create schema: %v", qErr)
		}
	}

	return db
}
