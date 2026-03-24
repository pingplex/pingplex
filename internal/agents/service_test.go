package agents_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pingplex/pingplex/internal/agents"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func TestServiceGetAgentByID_AccessRules(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupTestDB(t)
	repo := agents.NewRepository(db)
	svc := agents.New(repo, agents.NewMetricsWithRegisterer(prometheus.NewRegistry()), zap.NewNop())

	privateAgent, err := svc.Create(ctx, "user-1", agents.CreateAgentInput{
		Name:       "private-agent",
		Visibility: agents.VisibilityPrivate,
	})
	if err != nil {
		t.Fatalf("create private agent failed: %v", err)
	}

	_, err = svc.GetByID(ctx, "user-2", privateAgent.ID)
	if !errors.Is(err, agents.ErrForbidden) {
		t.Fatalf("expected ErrForbidden for private agent, got %v", err)
	}

	publicAgent, err := svc.Create(ctx, "user-1", agents.CreateAgentInput{
		Name:       "public-agent",
		Visibility: agents.VisibilityPublic,
	})
	if err != nil {
		t.Fatalf("create public agent failed: %v", err)
	}

	got, err := svc.GetByID(ctx, "user-2", publicAgent.ID)
	if err != nil {
		t.Fatalf("expected public access to succeed: %v", err)
	}

	if got.ID != publicAgent.ID {
		t.Fatalf("expected agent %s, got %s", publicAgent.ID, got.ID)
	}
}

func TestServiceListOwnedAgents_SkipsDeleted(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupTestDB(t)
	repo := agents.NewRepository(db)
	svc := agents.New(repo, agents.NewMetricsWithRegisterer(prometheus.NewRegistry()), zap.NewNop())

	agent, err := svc.Create(ctx, "user-1", agents.CreateAgentInput{
		Name:       "to-delete",
		Visibility: agents.VisibilityPrivate,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if err = svc.Delete(ctx, "user-1", agent.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	list, err := svc.ListOwned(ctx, "user-1")
	if err != nil {
		t.Fatalf("list owned failed: %v", err)
	}

	if len(list) != 0 {
		t.Fatalf("expected no non-deleted agents, got %d", len(list))
	}
}
