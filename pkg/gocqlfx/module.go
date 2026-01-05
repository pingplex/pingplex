package gocqlfx

import (
	"context"

	"github.com/go-core-fx/logger"
	"github.com/gocql/gocql"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"gocqlfx",
		logger.WithNamedLogger("gocqlfx"),
		fx.Provide(New),
		fx.Invoke(func(lc fx.Lifecycle, session *gocql.Session) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					return nil
				},
				OnStop: func(_ context.Context) error {
					session.Close()
					return nil
				},
			})
		}),
	)
}
