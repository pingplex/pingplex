package internal

import (
	"context"

	"github.com/go-core-fx/bunfx"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/goosefx"
	"github.com/go-core-fx/healthfx"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/sqlfx"
	"github.com/go-core-fx/telegofx"
	"github.com/pingplex/pingplex/internal/agents"
	"github.com/pingplex/pingplex/internal/bot"
	"github.com/pingplex/pingplex/internal/config"
	"github.com/pingplex/pingplex/internal/db"
	"github.com/pingplex/pingplex/internal/server"
	"github.com/pingplex/pingplex/internal/userdb"
	"github.com/pingplex/pingplex/internal/users"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(version healthfx.Version) {
	fx.New(
		// CORE MODULES
		logger.Module(),
		logger.WithFxDefaultLogger(),
		fiberfx.Module(),
		sqlfx.Module(),
		goosefx.Module(),
		telegofx.Module(true),
		bunfx.Module(),
		// redisfx.Module(),
		healthfx.Module(),
		//
		// APP MODULES
		config.Module(),
		db.Module(),
		server.Module(),
		bot.Module(),
		userdb.Module(),
		//
		// BUSINESS MODULES
		agents.Module(),
		users.Module(),
		//
		fx.Supply(version),
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					logger.Info("🚀 app started")
					return nil
				},
				OnStop: func(_ context.Context) error {
					logger.Info("👋 app stopped")
					return nil
				},
			})
		}),
	).Run()
}
