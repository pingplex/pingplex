package agents

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, userID string, in CreateAgentInput) (*Agent, error) {
	if userID == "" {
		return nil, ErrValidation
	}

	newAgent := newAgent(
		uuid.NewString(),
		userID,
		in.Name,
		in.Description,
		in.Visibility,
	)

	if _, err := r.db.NewInsert().Model(newAgent).Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("duplicate agent name for owner: %w", ErrValidation)
		}

		return nil, fmt.Errorf("create agent: %w", err)
	}

	return newAgent.toDomain(), nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Agent, error) {
	if id == "" {
		return nil, ErrValidation
	}

	var existing agent
	if err := r.db.NewSelect().Model(&existing).Where("id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("get agent by id: %w", err)
	}

	return existing.toDomain(), nil
}

func (r *Repository) ListByOwner(ctx context.Context, ownerUserID string) ([]Agent, error) {
	if ownerUserID == "" {
		return nil, ErrValidation
	}

	var items []agent
	if err := r.db.NewSelect().
		Model(&items).
		Where("user_id = ?", ownerUserID).
		Where("status != ?", StatusDeleted).
		Order("name").
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("list agents by owner: %w", err)
	}

	out := make([]Agent, 0, len(items))
	for _, item := range items {
		out = append(out, *item.toDomain())
	}

	return out, nil
}

func (r *Repository) ListPublic(ctx context.Context) ([]Agent, error) {
	var items []agent
	if err := r.db.NewSelect().
		Model(&items).
		Where("visibility = ?", VisibilityPublic).
		Where("status = ?", StatusActive).
		Order("created_at DESC").
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("list public agents: %w", err)
	}

	out := make([]Agent, 0, len(items))
	for _, item := range items {
		out = append(out, *item.toDomain())
	}

	return out, nil
}

func (r *Repository) Update(ctx context.Context, userID string, id string, in UpdateAgentInput) (*Agent, error) {
	if userID == "" || id == "" {
		return nil, ErrValidation
	}

	res, err := r.db.NewUpdate().
		Model((*agent)(nil)).
		Set("name = ?", in.Name).
		Set("description = ?", in.Description).
		Set("visibility = ?", in.Visibility).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Where("user_id = ?", userID).
		Where("status != ?", StatusDeleted).
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("duplicate agent name for owner: %w", ErrValidation)
		}

		return nil, fmt.Errorf("update agent: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("update agent rows affected: %w", err)
	}

	if affected == 0 {
		return nil, ErrNotFound
	}

	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, userID string, agentID string) error {
	if userID == "" || agentID == "" {
		return ErrValidation
	}

	res, err := r.db.NewUpdate().
		Model((*agent)(nil)).
		Set("status = ?", StatusDeleted).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", agentID).
		Where("user_id = ?", userID).
		Where("status != ?", StatusDeleted).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("soft delete agent: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("soft delete agent rows affected: %w", err)
	}

	if affected == 0 {
		return ErrNotFound
	}

	return nil
}
