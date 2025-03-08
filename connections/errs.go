package connections

import (
	"errors"
)

// Per-function error wrappers
var (
	ErrGetSchedules            = errors.New("failed to get schedules")
	ErrLoginFail               = errors.New("login has failed")
	ErrSignPacket              = errors.New("could not sign packet")
	ErrResponseTooShort        = errors.New("response too short")
	ErrWrongRemote             = errors.New("message did not originate from a Switcher device")
	ErrListenerRead            = errors.New("could not read message from listener")
	ErrTryNewListener          = errors.New("could not create a listener")
	ErrListenerClose           = errors.New("there was an error when closing the listener")
	ErrTryNewBidirectionalConn = errors.New("could not create a new bi-directional connection")
	ErrBiDirConnClose          = errors.New(
		"there was an error when closing the bi-directional connection",
	)
)
