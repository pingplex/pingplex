package targets_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pingplex/pingplex/internal/targets"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "modernc.org/sqlite"
)

func TestRepositoryCRUD(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupTestDB(t)
	repo := targets.NewRepository(db)

	created, err := repo.Create(ctx, "user-1", targets.CreateTargetInput{
		Name:       "api-main",
		Endpoint:   "https://example.com/health",
		Visibility: targets.VisibilityPrivate,
		CheckParams: map[string]string{
			targets.ParamTimeoutSeconds:  "8",
			targets.ParamRetries:         "3",
			targets.ParamIntervalSeconds: "30",
			"dns_expected_ip":            "1.1.1.1",
		},
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	fetched, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}

	if fetched.Endpoint != "https://example.com/health" {
		t.Fatalf("unexpected endpoint: %s", fetched.Endpoint)
	}

	if fetched.CheckParams["dns_expected_ip"] != "1.1.1.1" {
		t.Fatalf("expected dns_expected_ip param, got: %q", fetched.CheckParams["dns_expected_ip"])
	}

	updated, err := repo.Update(ctx, "user-1", created.ID, targets.UpdateTargetInput{
		Name:        "api-main-v2",
		Endpoint:    "https://example.com/v2/health",
		Description: "updated",
		Visibility:  targets.VisibilityPublic,
		CheckParams: map[string]string{
			targets.ParamTimeoutSeconds:  "12",
			targets.ParamRetries:         "4",
			targets.ParamIntervalSeconds: "45",
			"tls_expiry_threshold_days":  "21",
		},
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}

	if updated.Name != "api-main-v2" {
		t.Fatalf("expected updated name, got %s", updated.Name)
	}

	if updated.CheckParams["tls_expiry_threshold_days"] != "21" {
		t.Fatalf("expected threshold param to be persisted, got: %q", updated.CheckParams["tls_expiry_threshold_days"])
	}

	refetched, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("refetch after update failed: %v", err)
	}
	if refetched.Endpoint != "https://example.com/v2/health" {
		t.Fatalf("expected persisted endpoint, got %s", refetched.Endpoint)
	}
	if refetched.CheckParams["tls_expiry_threshold_days"] != "21" {
		t.Fatalf("expected persisted threshold param, got: %q", refetched.CheckParams["tls_expiry_threshold_days"])
	}

	list, err := repo.ListByOwner(ctx, "user-1")
	if err != nil {
		t.Fatalf("list by owner failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 target, got %d", len(list))
	}

	if err = repo.Delete(ctx, "user-1", created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = repo.GetByID(ctx, created.ID)
	if !errors.Is(err, targets.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRepositoryUpdateWithWrongOwnerReturnsNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupTestDB(t)
	repo := targets.NewRepository(db)

	created, err := repo.Create(ctx, "user-1", targets.CreateTargetInput{
		Name:       "api-main",
		Endpoint:   "https://example.com/health",
		Visibility: targets.VisibilityPrivate,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, err = repo.Update(ctx, "user-2", created.ID, targets.UpdateTargetInput{
		Name:       "new-name",
		Endpoint:   "https://example.com/new",
		Visibility: targets.VisibilityPublic,
	})
	if !errors.Is(err, targets.ErrNotFound) {
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
	sqldb.SetMaxOpenConns(1)

	t.Cleanup(func() {
		_ = sqldb.Close()
	})

	db := bun.NewDB(sqldb, sqlitedialect.New())
	t.Cleanup(func() {
		_ = db.Close()
	})

	_, err = db.ExecContext(context.Background(), "PRAGMA foreign_keys = ON;")
	if err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	schema := []string{
		`CREATE TABLE users (
			id TEXT PRIMARY KEY,
			status TEXT NOT NULL DEFAULT 'active',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE targets (
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
