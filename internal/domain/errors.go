package domain

import "errors"

// Sentinel errors used by the application and delivery layers to map to HTTP status codes.
var (
	ErrNotFound    = errors.New("not found")
	ErrConflict    = errors.New("conflict")
	ErrForbidden   = errors.New("forbidden")
	ErrValidation  = errors.New("validation error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInternal    = errors.New("internal server error")
)
