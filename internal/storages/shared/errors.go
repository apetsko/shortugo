package shared

import "errors"

// ErrNotFound is returned when a requested resource is not found.
var ErrNotFound = errors.New("not found")

// ErrGone is returned when a requested resource is permanently deleted.
var ErrGone = errors.New("Gone")
