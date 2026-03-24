package userdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-core-fx/goosefx"
	"github.com/go-core-fx/sqlfx"
	"github.com/pingplex/pingplex/internal/userdb/migrations"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type Service struct {
	config Config

	logger *zap.Logger

	conns map[string]*sql.DB
	mu    sync.RWMutex
}

func New(config Config, logger *zap.Logger) *Service {
	return &Service{
		config: config,

		logger: logger,

		conns: map[string]*sql.DB{},
		mu:    sync.RWMutex{},
	}
}

func (s *Service) Get(ctx context.Context, userID string) (*sql.DB, error) {
	s.mu.RLock()
	conn, ok := s.conns[userID]
	s.mu.RUnlock()

	if ok {
		return conn, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Re-check after acquiring write lock
	if conn, ok = s.conns[userID]; ok {
		return conn, nil
	}

	db, err := sqlfx.New(sqlfx.Config{
		URL:             strings.Replace(s.config.URLTemplate, "{user_id}", userID, 1),
		ConnMaxIdleTime: s.config.ConnMaxIdleTime,
		ConnMaxLifetime: s.config.ConnMaxLifetime,
		MaxOpenConns:    s.config.MaxOpenConns,
		MaxIdleConns:    s.config.MaxIdleConns,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to user db: %w", err)
	}

	migrator, err := goosefx.NewProvider(db, goosefx.Storage(migrations.FS), goose.DialectSQLite3)
	if err != nil {
		err = errors.Join(err, db.Close())
		return nil, fmt.Errorf("failed to create user db migrator: %w", err)
	}

	up, err := migrator.Up(ctx)
	if err != nil {
		err = errors.Join(err, db.Close())
		return nil, fmt.Errorf("failed to migrate user db: %w", err)
	}

	if len(up) > 0 {
		s.logger.Info("migrated user db", zap.Int("migrations", len(up)))
	}

	s.conns[userID] = db

	return db, nil
}

func (s *Service) Close() error {
	var err error

	s.mu.Lock()
	for _, conn := range s.conns {
		err = errors.Join(err, conn.Close())
	}
	s.conns = map[string]*sql.DB{}
	s.mu.Unlock()

	return err
}
