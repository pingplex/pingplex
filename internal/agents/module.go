package agents

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"agents",
		logger.WithNamedLogger("agents"),
		fx.Provide(NewMetrics, fx.Private),
		fx.Provide(NewRepository, fx.Private),
		fx.Provide(New),
	)
}
