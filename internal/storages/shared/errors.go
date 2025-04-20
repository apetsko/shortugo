// Package shared contains shared utilities and error definitions used across the application.
// It includes common error types and other reusable components.
package shared

import "errors"

// ErrNotFound is returned when a requested resource is not found.
var ErrNotFound = errors.New("not found")

// ErrGone is returned when a requested resource is permanently deleted.
var ErrGone = errors.New("Gone")
