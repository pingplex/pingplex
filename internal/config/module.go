package config

import (
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/redisfx"
	"github.com/pingplex/pingplex/pkg/gocqlfx"
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
		fx.Provide(func(cfg Config) gocqlfx.Config {
			return gocqlfx.Config{
				Hosts:    cfg.Database.Hosts,
				Keyspace: cfg.Database.Keyspace,
				Username: cfg.Database.Username,
				Password: cfg.Database.Password,
			}
		}),
		fx.Provide(func(cfg Config) redisfx.Config {
			return redisfx.Config{
				URL: cfg.Redis.URL,
			}
		}),
	)
}
