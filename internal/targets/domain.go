package targets

import (
	"maps"
	"strings"
	"time"
)

type Visibility string

type Status string

const (
	VisibilityPrivate  Visibility = "private"  // Visible to the private agents only
	VisibilityInternal Visibility = "internal" // Visible to the private and service-provided agents
	VisibilityPublic   Visibility = "public"   // Visible to the any agents
)

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusDeleted  Status = "deleted"
)

const (
	ParamTimeoutSeconds  = "timeout_seconds"
	ParamRetries         = "retries"
	ParamIntervalSeconds = "interval_seconds"
)

func defaultCheckParams() map[string]string {
	return map[string]string{
		ParamTimeoutSeconds:  "5",
		ParamRetries:         "1",
		ParamIntervalSeconds: "60",
	}
}

type Target struct {
	ID string

	UserID      string
	Name        string
	Endpoint    string
	Description string

	Visibility Visibility
	Status     Status

	CheckParams map[string]string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateTargetInput struct {
	Name        string
	Endpoint    string
	Description string
	Visibility  Visibility

	CheckParams map[string]string
}

func (in *CreateTargetInput) normalize() {
	params := defaultCheckParams()

	maps.Copy(params, in.CheckParams)

	in.CheckParams = params
}

func (in *CreateTargetInput) Validate() error {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Endpoint) == "" {
		return ErrValidation
	}

	switch in.Visibility {
	case VisibilityPrivate, VisibilityInternal, VisibilityPublic:
	default:
		return ErrValidation
	}

	for key, value := range in.CheckParams {
		if strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			return ErrValidation
		}
	}

	return nil
}

type UpdateTargetInput = CreateTargetInput
