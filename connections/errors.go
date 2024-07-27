package connections

import "errors"

// Per-function error wrappers
var (
	ErrGetSchedules            = errors.New("failed to get schedules")
	ErrLoginFail               = errors.New("login has failed")
	ErrListenerRead            = errors.New("could not read message from listener")
	ErrSignPacket              = errors.New("could not sign packet")
	ErrTryNewBidirectionalConn = errors.New("could not create a new bi-directional connection")
)
