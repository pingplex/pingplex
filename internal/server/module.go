package server

import (
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/fiberfx/health"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-core-fx/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module(
		"server",
		logger.WithNamedLogger("server"),

		fx.Provide(func(log *zap.Logger) fiberfx.Options {
			opts := fiberfx.Options{}
			opts.WithErrorHandler(fiberfx.NewJSONErrorHandler(log))
			opts.WithMetrics()
			return opts
		}),

		fx.Provide(
			fx.Annotate(health.NewHandler, fx.ResultTags(`group:"handlers"`)), fx.Private,
			// fx.Annotate(stacks.NewHandler, fx.ResultTags(`group:"handlers"`)), fx.Private,
		),

		fx.Invoke(
			fx.Annotate(
				func(handlers []handler.Handler, app *fiber.App) {
					// Swagger documentation
					// app.Get("/swagger/*", fiberSwagger.WrapHandler)

					// Version 1 API group
					v1 := app.Group("/api/v1")
					v1.Use(validation.Middleware)

					for _, h := range handlers {
						h.Register(v1)
					}
				},
				fx.ParamTags(`group:"handlers"`),
			),
		),
	)
}
