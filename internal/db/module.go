package db

import (
	"github.com/go-core-fx/logger"
	"github.com/pingplex/pingplex/internal/db/migrations"
	"github.com/pingplex/pingplex/pkg/gocqlxfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"db",
		logger.WithNamedLogger("db"),
		fx.Provide(func() gocqlxfx.Storage {
			return migrations.Files
		}),
	)
}
