package agents

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	agents  *Repository
	metrics *Metrics
	logger  *zap.Logger
}

func New(agents *Repository, metrics *Metrics, logger *zap.Logger) *Service {
	return &Service{agents: agents, metrics: metrics, logger: logger}
}

func (s *Service) Create(ctx context.Context, userID string, in CreateAgentInput) (*Agent, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	startedAt := time.Now()

	agent, err := s.agents.Create(ctx, userID, in)
	s.metrics.Observe("create_agent", startedAt, err)
	if err != nil {
		s.logger.Warn("create agent failed", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	s.logger.Debug("agent created", zap.String("agent_id", agent.ID), zap.String("user_id", userID))
	return agent, nil
}

func (s *Service) Update(ctx context.Context, userID string, id string, in UpdateAgentInput) (*Agent, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	startedAt := time.Now()

	agent, err := s.agents.Update(ctx, userID, id, in)
	s.metrics.Observe("update_agent", startedAt, err)
	if err != nil {
		s.logger.Warn("update agent failed", zap.String("user_id", userID), zap.String("agent_id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Debug("agent updated", zap.String("agent_id", id), zap.String("user_id", userID))
	return agent, nil
}

func (s *Service) Delete(ctx context.Context, userID string, agentID string) error {
	startedAt := time.Now()
	err := s.agents.Delete(ctx, userID, agentID)
	s.metrics.Observe("delete_agent", startedAt, err)
	if err != nil {
		s.logger.Warn(
			"delete agent failed",
			zap.String("user_id", userID),
			zap.String("agent_id", agentID),
			zap.Error(err),
		)
		return err
	}

	s.logger.Debug("agent deleted", zap.String("agent_id", agentID), zap.String("user_id", userID))
	return nil
}

func (s *Service) GetByID(ctx context.Context, userID string, agentID string) (*Agent, error) {
	startedAt := time.Now()

	agent, err := s.agents.GetByID(ctx, agentID)
	if err == nil {
		if agent.Status == StatusDeleted {
			err = ErrNotFound
		} else if agent.UserID != userID && agent.Visibility != VisibilityPublic {
			err = ErrForbidden
		}
	}

	s.metrics.Observe("get_agent_by_id", startedAt, err)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func (s *Service) ListOwned(ctx context.Context, userID string) ([]Agent, error) {
	startedAt := time.Now()

	items, err := s.agents.ListByOwner(ctx, userID)

	s.metrics.Observe("list_owned_agents", startedAt, err)
	return items, err
}

func (s *Service) ListPublic(ctx context.Context) ([]Agent, error) {
	startedAt := time.Now()
	items, err := s.agents.ListPublic(ctx)
	s.metrics.Observe("list_public_agents", startedAt, err)
	return items, err
}
