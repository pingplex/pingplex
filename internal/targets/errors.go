package targets

import "errors"

var (
	ErrNotFound   = errors.New("target not found")
	ErrValidation = errors.New("validation error")
	ErrForbidden  = errors.New("forbidden")
)
