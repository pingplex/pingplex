package config

import (
	"github.com/go-core-fx/fiberfx"
	"github.com/pingplex/pingplex/internal/example"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"config",
		fx.Provide(New),
		fx.Provide(func(cfg Config) fiberfx.Config {
			return fiberfx.Config{
				Address:     cfg.HTTP.Address,
				ProxyHeader: cfg.HTTP.ProxyHeader,
				Proxies:     cfg.HTTP.Proxies,
			}
		}),
		fx.Provide(func(cfg Config) example.Config {
			return example.Config{
				Example: cfg.Example.Example,
			}
		}),
	)
}
