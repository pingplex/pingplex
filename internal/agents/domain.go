package agents

import (
	"strings"
	"time"
)

type Visibility string

type Status string

const (
	VisibilityPrivate Visibility = "private"
	VisibilityPublic  Visibility = "public"
)

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusDeleted  Status = "deleted"
)

type Agent struct {
	ID string

	UserID      string
	Name        string
	Description string

	Visibility Visibility
	Status     Status

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateAgentInput struct {
	Name        string
	Description string
	Visibility  Visibility
}

func (in CreateAgentInput) Validate() error {
	if strings.TrimSpace(in.Name) == "" {
		return ErrValidation
	}

	if in.Visibility != VisibilityPrivate && in.Visibility != VisibilityPublic {
		return ErrValidation
	}

	return nil
}

type UpdateAgentInput = CreateAgentInput
