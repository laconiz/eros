package oceanus

import "errors"

var (
	ErrTimeout       = errors.New("timeout")
	ErrNodeNotFound  = errors.New("node not found")
	ErrInvalidSwitch = errors.New("invalid switch")
	ErrDisconnected  = errors.New("connection disconnected")
)
