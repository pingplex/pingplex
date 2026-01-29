package config

import (
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/redisfx"
	"github.com/go-core-fx/sqlfx"
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
		fx.Provide(
			func(cfg Config) redisfx.Config {
				return redisfx.Config{
					URL: cfg.Redis.URL,
				}
			},
			func(cfg Config) sqlfx.Config {
				return sqlfx.Config{
					URL: cfg.Database.URL,

					ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
					ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
					MaxOpenConns:    cfg.Database.MaxOpenConns,
					MaxIdleConns:    cfg.Database.MaxIdleConns,
				}
			},
		),
	)
}
