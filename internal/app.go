package internal

import (
	"context"

	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/healthfx"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/redisfx"
	"github.com/pingplex/pingplex/internal/config"
	"github.com/pingplex/pingplex/internal/db"
	"github.com/pingplex/pingplex/internal/server"
	"github.com/pingplex/pingplex/pkg/gocqlfx"
	"github.com/pingplex/pingplex/pkg/gocqlxfx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(version healthfx.Version) {
	fx.New(
		// CORE MODULES
		logger.Module(),
		logger.WithFxDefaultLogger(),
		fiberfx.Module(),
		gocqlfx.Module(),
		gocqlxfx.Module(),
		redisfx.Module(),
		healthfx.Module(),
		//
		// APP MODULES
		config.Module(),
		db.Module(),
		server.Module(),
		// bot.Module(),
		//
		// BUSINESS MODULES
		// example.Module(),
		//
		fx.Supply(version),
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					logger.Info("app started")
					return nil
				},
				OnStop: func(_ context.Context) error {
					logger.Info("app stopped")
					return nil
				},
			})
		}),
	).Run()
}
