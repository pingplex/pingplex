package bot

import (
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	"github.com/pingplex/pingplex/internal/bot/handler"
	"github.com/pingplex/pingplex/internal/bot/handlers/start"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"bot",
		logger.WithNamedLogger("bot"),
		fx.Provide(func() []telego.BotOption {
			return nil
		}),
		fx.Provide(
			fx.Annotate(start.New, fx.ResultTags(`group:"handlers"`)),
		),
		fx.Invoke(
			fx.Annotate(
				func(handlers []handler.Handler, r *telegofx.Router) {
					for _, h := range handlers {
						h.Register(r)
					}
				},
				fx.ParamTags(`group:"handlers"`),
			),
		),
	)
}
