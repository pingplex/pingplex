package internal

import (
	"context"

	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/logger"
	"github.com/pingplex/pingplex/internal/config"
	"github.com/pingplex/pingplex/internal/db"
	"github.com/pingplex/pingplex/internal/server"
	"github.com/pingplex/pingplex/pkg/gocqlfx"
	"github.com/pingplex/pingplex/pkg/gocqlxfx"
	"github.com/scylladb/gocqlx/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run() {
	fx.New(
		// CORE MODULES
		logger.Module(),
		logger.WithFxDefaultLogger(),
		fiberfx.Module(),
		gocqlfx.Module(),
		gocqlxfx.Module(),
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
		fx.Invoke(func(lc fx.Lifecycle, _ gocqlx.Session, logger *zap.Logger) {
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
