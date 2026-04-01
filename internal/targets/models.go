package targets

import (
	"encoding/json"
	"time"

	"github.com/pingplex/pingplex/internal/db"
	"github.com/uptrace/bun"
)

type target struct {
	bun.BaseModel `bun:"table:targets,alias:t"`
	db.TimedModel

	ID string `bun:"id,pk"`

	UserID      string `bun:"user_id"`
	Name        string `bun:"name"`
	Endpoint    string `bun:"endpoint"`
	Description string `bun:"description,nullzero"`

	Visibility Visibility `bun:"visibility"`
	Status     Status     `bun:"status"`

	CheckParams string `bun:"check_params"`
}

func newTarget(id string, userID string, in CreateTargetInput) *target {
	now := time.Now()
	paramsJSON := mustMarshalParams(in.CheckParams)

	return &target{
		BaseModel: bun.BaseModel{},
		TimedModel: db.TimedModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
		ID:          id,
		UserID:      userID,
		Name:        in.Name,
		Endpoint:    in.Endpoint,
		Description: in.Description,
		Visibility:  in.Visibility,
		Status:      StatusActive,
		CheckParams: paramsJSON,
	}
}

func (t *target) toDomain() *Target {
	if t == nil {
		return nil
	}

	params := parseParamsOrDefault(t.CheckParams)

	return &Target{
		ID: t.ID,

		UserID:      t.UserID,
		Name:        t.Name,
		Endpoint:    t.Endpoint,
		Description: t.Description,

		Visibility: t.Visibility,
		Status:     t.Status,

		CheckParams: params,

		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func mustMarshalParams(params map[string]string) string {
	if params == nil {
		return ""
	}

	encoded, err := json.Marshal(params)
	if err != nil {
		return "{}"
	}

	return string(encoded)
}

func parseParamsOrDefault(raw string) map[string]string {
	if raw == "" {
		return map[string]string{}
	}

	out := map[string]string{}
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return map[string]string{}
	}

	return out
}
