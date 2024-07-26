package consts

import "errors"

// Common errors
var (
	ErrNotImplemeted = errors.New("not implemented")
	ErrInvalidIP     = errors.New("invalid IP address")
	ErrInvalidPort   = errors.New("invalid port (must be between 100-65000)")
)
