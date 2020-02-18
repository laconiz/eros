package session

import "sync/atomic"

var incremental uint64

func Increment() ID {
	return ID(atomic.AddUint64(&incremental, 1))
}
