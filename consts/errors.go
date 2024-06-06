package consts

import "errors"

var (
	ErrInvalidIP   = errors.New("invalid IP address")
	ErrInvalidPort = errors.New("invalid port (must be between 100-65000)")
)
