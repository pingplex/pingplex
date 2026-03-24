package userdb

import (
	"context"

	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"userdb",
		logger.WithNamedLogger("userdb"),
		fx.Provide(New),
		fx.Invoke(func(lc fx.Lifecycle, service *Service) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					return nil
				},
				OnStop: func(_ context.Context) error {
					return service.Close()
				},
			})
		}),
	)
}
