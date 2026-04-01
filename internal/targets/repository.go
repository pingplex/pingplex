package targets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, userID string, in CreateTargetInput) (*Target, error) {
	if userID == "" {
		return nil, ErrValidation
	}

	item := newTarget(uuid.NewString(), userID, in)

	if _, err := r.db.NewInsert().Model(item).Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("duplicate target name for owner: %w", ErrValidation)
		}

		return nil, fmt.Errorf("create target: %w", err)
	}

	return item.toDomain(), nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Target, error) {
	if id == "" {
		return nil, ErrValidation
	}

	var existing target
	if err := r.db.NewSelect().Model(&existing).Where("id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("get target by id: %w", err)
	}

	return existing.toDomain(), nil
}

func (r *Repository) ListByOwner(ctx context.Context, ownerUserID string) ([]Target, error) {
	if ownerUserID == "" {
		return nil, ErrValidation
	}

	var items []target
	if err := r.db.NewSelect().
		Model(&items).
		Where("user_id = ?", ownerUserID).
		Where("status != ?", StatusDeleted).
		Order("name").
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("list targets by owner: %w", err)
	}

	out := make([]Target, 0, len(items))
	for _, item := range items {
		out = append(out, *item.toDomain())
	}

	return out, nil
}

func (r *Repository) ListPublic(ctx context.Context) ([]Target, error) {
	var items []target
	if err := r.db.NewSelect().
		Model(&items).
		Where("visibility = ?", VisibilityPublic).
		Where("status = ?", StatusActive).
		Order("created_at DESC").
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("list public targets: %w", err)
	}

	out := make([]Target, 0, len(items))
	for _, item := range items {
		out = append(out, *item.toDomain())
	}

	return out, nil
}

func (r *Repository) Update(ctx context.Context, userID string, id string, in UpdateTargetInput) (*Target, error) {
	if userID == "" || id == "" {
		return nil, ErrValidation
	}

	item := newTarget(id, userID, in)

	res, err := r.db.NewUpdate().
		OmitZero().
		Model(item).
		ExcludeColumn("created_at").
		Where("id = ?", id).
		Where("user_id = ?", userID).
		Where("status != ?", StatusDeleted).
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("duplicate target name for owner: %w", ErrValidation)
		}

		return nil, fmt.Errorf("update target: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("update target rows affected: %w", err)
	}

	if affected == 0 {
		return nil, ErrNotFound
	}

	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, userID string, targetID string) error {
	if userID == "" || targetID == "" {
		return ErrValidation
	}

	res, err := r.db.NewDelete().
		Model((*target)(nil)).
		Where("id = ?", targetID).
		Where("user_id = ?", userID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete target: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete target rows affected: %w", err)
	}

	if affected == 0 {
		return ErrNotFound
	}

	return nil
}
