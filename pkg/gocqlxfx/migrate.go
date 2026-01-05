package gocqlxfx

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/migrate"
)

type Storage fs.FS

type Migrator struct {
	db gocqlx.Session

	storage Storage
}

func NewMigrator(db gocqlx.Session, storage Storage) *Migrator {
	if storage == nil {
		return nil
	}

	return &Migrator{
		db: db,

		storage: storage,
	}
}

func (m *Migrator) Migrate(ctx context.Context) error {
	if m == nil {
		return nil
	}

	if err := migrate.FromFS(ctx, m.db, m.storage); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
