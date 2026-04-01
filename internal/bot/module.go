package bot

import (
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	"github.com/pingplex/pingplex/internal/bot/handler"
	"github.com/pingplex/pingplex/internal/bot/handlers/agents"
	"github.com/pingplex/pingplex/internal/bot/handlers/start"
	"github.com/pingplex/pingplex/internal/bot/middlewares/userauth"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"bot",
		logger.WithNamedLogger("bot"),
		fx.Provide(func() []telego.BotOption {
			return []telego.BotOption{
				telego.WithFastHTTPClient(&fasthttp.Client{Dial: fasthttpproxy.FasthttpProxyHTTPDialer()}),
			}
		}),
		fx.Provide(
			fx.Annotate(userauth.New, fx.ResultTags(`group:"middlewares"`)),

			fx.Annotate(start.New, fx.ResultTags(`group:"handlers"`)),
			fx.Annotate(agents.New, fx.ResultTags(`group:"handlers"`)),
		),
		fx.Invoke(
			fx.Annotate(
				func(handlers []handler.Handler, middlewares []th.Handler, r *telegofx.Router) {
					for _, m := range middlewares {
						r.Use(m)
					}

					for _, h := range handlers {
						h.Register(r)
					}
				},
				fx.ParamTags(`group:"handlers"`, `group:"middlewares"`),
			),
		),
	)
}
