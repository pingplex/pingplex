package targets

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module(
		"targets",
		fx.Provide(NewRepository, fx.Private),
		fx.Provide(New),
	)
}
