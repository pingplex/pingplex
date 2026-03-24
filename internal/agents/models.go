package agents

import (
	"time"

	"github.com/pingplex/pingplex/internal/db"
	"github.com/uptrace/bun"
)

type agent struct {
	bun.BaseModel `bun:"table:agents,alias:a"`
	db.TimedModel

	ID string `bun:"id,pk"`

	UserID      string `bun:"user_id"`
	Name        string `bun:"name"`
	Description string `bun:"description,nullzero"`

	Visibility Visibility `bun:"visibility"`
	Status     Status     `bun:"status"`
}

func newAgent(id string, userID string, name string, description string, visibility Visibility) *agent {
	now := time.Now()

	return &agent{
		BaseModel: bun.BaseModel{},
		TimedModel: db.TimedModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
		ID:          id,
		UserID:      userID,
		Name:        name,
		Description: description,
		Visibility:  visibility,
		Status:      StatusActive,
	}
}

func (a *agent) toDomain() *Agent {
	if a == nil {
		return nil
	}

	return &Agent{
		ID: a.ID,

		UserID:      a.UserID,
		Name:        a.Name,
		Description: a.Description,

		Visibility: a.Visibility,
		Status:     a.Status,

		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
