package targets

import (
	"context"
)

type Service struct {
	targets *Repository
}

func New(targets *Repository) *Service {
	return &Service{targets: targets}
}

func (s *Service) Create(ctx context.Context, userID string, in CreateTargetInput) (*Target, error) {
	in.normalize()
	if err := in.Validate(); err != nil {
		return nil, err
	}

	return s.targets.Create(ctx, userID, in)
}

func (s *Service) Update(ctx context.Context, userID string, id string, in UpdateTargetInput) (*Target, error) {
	in.normalize()
	if err := in.Validate(); err != nil {
		return nil, err
	}

	return s.targets.Update(ctx, userID, id, in)
}

func (s *Service) Delete(ctx context.Context, userID string, targetID string) error {
	return s.targets.Delete(ctx, userID, targetID)
}

func (s *Service) GetByID(ctx context.Context, userID string, targetID string) (*Target, error) {
	target, err := s.targets.GetByID(ctx, targetID)
	if err != nil {
		return nil, err
	}

	if target.Status == StatusDeleted {
		return nil, ErrNotFound
	}

	if target.UserID != userID && target.Visibility != VisibilityPublic {
		return nil, ErrNotFound
	}

	return target, nil
}

func (s *Service) ListOwned(ctx context.Context, userID string) ([]Target, error) {
	return s.targets.ListByOwner(ctx, userID)
}

func (s *Service) ListPublic(ctx context.Context) ([]Target, error) {
	return s.targets.ListPublic(ctx)
}
