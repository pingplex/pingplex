package db

import (
	"github.com/go-core-fx/goosefx"
	"github.com/go-core-fx/logger"
	"github.com/pingplex/pingplex/internal/db/migrations"
	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/schema"
	"go.uber.org/fx"

	_ "modernc.org/sqlite" // Pure-Go SQLite driver
)

func Module() fx.Option {
	return fx.Module(
		"db",
		logger.WithNamedLogger("db"),
		fx.Provide(func() schema.Dialect {
			return sqlitedialect.New()
		}),
		fx.Provide(func() goose.Dialect {
			return goose.DialectSQLite3
		}),
		fx.Provide(func() goosefx.Storage {
			return goosefx.Storage(migrations.FS)
		}),
	)
}
