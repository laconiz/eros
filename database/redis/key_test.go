package redis

import (
	"testing"
	"time"
)

func TestKey(t *testing.T) {

	k := r.Key()

	// Remove Exist
	assert(k.Delete(key) == nil)
	ok, err := k.Exist(key)
	assert(err == nil, !ok)
	assert(k.Set(key, 1) == nil)
	ok, err = k.Exist(key)
	assert(err == nil, ok)

	// Set Get
	assert(k.Set(key, 10) == nil)
	var Int int
	exists, err := k.Get(key, &Int)
	assert(err == nil, exists, Int == 10)
	assert(k.Set(key, true) == nil)

	var Bool bool
	exists, err = k.Get(key, &Bool)
	assert(err == nil, exists, Bool)
	assert(k.Set(key, false) == nil)
	exists, err = k.Get(key, &Bool)
	assert(err == nil, exists, !Bool)

	assert(k.Set(key, ComplexPointer) == nil)
	Complex := &Struct{}
	exists, err = k.Get(key, Complex)
	assert(err == nil, exists, Complex.Equal(*ComplexPointer))

	assert(k.Set(key, []int{11, 12}) == nil)
	var slice []int
	exists, err = k.Get(key, &slice)
	assert(err == nil, exists, len(slice) == 2, slice[0] == 11, slice[1] == 12)

	// SetNX
	assert(k.Delete(key) == nil)
	Bool, err = k.SetNX(key, 20)
	assert(err == nil, Bool)
	exists, err = k.Get(key, &Int)
	assert(err == nil, exists, Int == 20)
	Bool, err = k.SetNX(key, 21)
	assert(err == nil, !Bool)
	Int = 0
	exists, err = k.Get(key, &Int)
	assert(err == nil, exists, Int == 20)

	// SetNEX
	assert(k.Delete(key) == nil)
	Bool, err = k.SetNEX(key, 30, 1)
	assert(err == nil, Bool)
	exists, err = k.Get(key, &Int)
	assert(err == nil, exists, Int == 30)

	Bool, err = k.SetNEX(key, 2, 1)
	assert(err == nil, !Bool)
	Int = 0
	exists, err = k.Get(key, &Int)
	assert(err == nil, exists, Int == 30)

	time.Sleep(time.Second)
	Bool, err = k.SetNEX(key, 31, 1)
	assert(err == nil, Bool)
	exists, err = k.Get(key, &Int)
	assert(err == nil, exists, Int == 31)

	// SetEX
	assert(k.Delete(key) == nil)
	assert(k.SetEX(key, 2, 1) == nil)
	exists, err = k.Get(key, &Int)
	assert(err == nil, exists, Int == 2)
	time.Sleep(time.Second)
	exists, err = k.Get(key, &Int)
	assert(err == nil, !exists)

	// Incr
	assert(k.Delete(key) == nil)
	Int64, err := k.Incr(key, 40)
	assert(err == nil, Int64 == 40)
	Int64, err = k.Incr(key, 41)
	assert(err == nil, Int64 == 81)
	Int64 = 0
	exists, err = k.Get(key, &Int64)
	assert(err == nil, exists, Int64 == 81)

	assert(k.Delete(key) == nil)
}
