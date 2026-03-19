package users_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/pingplex/pingplex/internal/users"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "modernc.org/sqlite"
)

func TestRepositoryRegisterOrLoginIsIdempotent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupTestDB(t)
	repo := users.NewRepository(db)

	ident := users.Identity{
		Provider:     "telegram",
		ProviderID:   "123456789",
		ProviderData: `{"username":"johndoe"}`,
	}

	firstUser, err := repo.RegisterOrLogin(ctx, ident)
	if err != nil {
		t.Fatalf("first register or login failed: %v", err)
	}

	secondUser, err := repo.RegisterOrLogin(ctx, ident)
	if err != nil {
		t.Fatalf("second register or login failed: %v", err)
	}

	if firstUser.ID == "" {
		t.Fatalf("expected non-empty user id")
	}

	if firstUser.ID != secondUser.ID {
		t.Fatalf("expected same user id, got %q and %q", firstUser.ID, secondUser.ID)
	}

	usersCount := countRows(t, db, "users")
	if usersCount != 1 {
		t.Fatalf("expected one user row, got %d", usersCount)
	}

	identitiesCount := countRows(t, db, "user_identities")
	if identitiesCount != 1 {
		t.Fatalf("expected one identity row, got %d", identitiesCount)
	}
}

func setupTestDB(t *testing.T) *bun.DB {
	t.Helper()

	sqldb, err := sql.Open("sqlite", "file::memory:?cache=shared")
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
		`CREATE TABLE user_identities (
			id INTEGER PRIMARY KEY,
			user_id TEXT NOT NULL,
			provider TEXT NOT NULL,
			provider_id TEXT NOT NULL,
			provider_data TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(provider, provider_id),
			UNIQUE(user_id, provider)
		);`,
	}

	for _, query := range schema {
		if _, qErr := db.ExecContext(context.Background(), query); qErr != nil {
			t.Fatalf("failed to create schema: %v", qErr)
		}
	}

	return db
}

func countRows(t *testing.T, db *bun.DB, table string) int {
	t.Helper()

	var count int
	query := "SELECT COUNT(*) FROM " + table
	if err := db.QueryRowContext(context.Background(), query).Scan(&count); err != nil {
		t.Fatalf("failed to count rows in %s: %v", table, err)
	}

	return count
}
