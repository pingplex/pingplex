package agents

import "errors"

var (
	ErrNotFound     = errors.New("agent not found")
	ErrValidation   = errors.New("validation failed")
	ErrForbidden    = errors.New("forbidden")
	ErrAgentDeleted = errors.New("agent deleted")
)
