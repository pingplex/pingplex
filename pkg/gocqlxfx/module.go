package gocqlxfx

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"gocqlxfx",
		logger.WithNamedLogger("gocqlxfx"),
		fx.Provide(New),
		fx.Provide(fx.Annotate(
			NewMigrator,
			fx.ParamTags("", `optional:"true"`),
		)),
		fx.Invoke(func(lc fx.Lifecycle, migrator *Migrator) {
			lc.Append(
				fx.StartHook(migrator.Migrate),
			)
		}),
	)
}
