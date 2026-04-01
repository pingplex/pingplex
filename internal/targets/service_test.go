package targets_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pingplex/pingplex/internal/targets"
)

func TestServiceGetByID_AccessRules(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupTestDB(t)
	repo := targets.NewRepository(db)
	svc := targets.New(repo)

	privateTarget, err := svc.Create(ctx, "user-1", targets.CreateTargetInput{
		Name:       "private-api",
		Endpoint:   "https://private.example/health",
		Visibility: targets.VisibilityPrivate,
	})
	if err != nil {
		t.Fatalf("create private target failed: %v", err)
	}

	_, err = svc.GetByID(ctx, "user-2", privateTarget.ID)
	if !errors.Is(err, targets.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for private target, got %v", err)
	}

	publicTarget, err := svc.Create(ctx, "user-1", targets.CreateTargetInput{
		Name:       "public-api",
		Endpoint:   "https://public.example/health",
		Visibility: targets.VisibilityPublic,
	})
	if err != nil {
		t.Fatalf("create public target failed: %v", err)
	}

	got, err := svc.GetByID(ctx, "user-2", publicTarget.ID)
	if err != nil {
		t.Fatalf("expected public access to succeed: %v", err)
	}
	if got.ID != publicTarget.ID {
		t.Fatalf("expected target %s, got %s", publicTarget.ID, got.ID)
	}
}

func TestServiceCreate_DefaultsCheckParams(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupTestDB(t)
	svc := targets.New(targets.NewRepository(db))

	created, err := svc.Create(ctx, "user-1", targets.CreateTargetInput{
		Name:       "defaulted-target",
		Endpoint:   "https://example.com/live",
		Visibility: targets.VisibilityPrivate,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if created.CheckParams[targets.ParamTimeoutSeconds] != "5" || created.CheckParams[targets.ParamRetries] != "1" ||
		created.CheckParams[targets.ParamIntervalSeconds] != "60" {
		t.Fatalf("unexpected defaults params: %#v", created.CheckParams)
	}
}
