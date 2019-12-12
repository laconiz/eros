package redis

import (
	"testing"
)

func TestSet(t *testing.T) {

	s := r.Set(key)

	assert(r.Key().Delete(key) == nil)

	assert(s.Add(key, key2) == nil)

	var keys []string
	assert(s.Keys(&keys) == nil, len(keys) == 2)

	assert(s.Remove(key2) == nil)
	assert(s.Keys(&keys) == nil, len(keys) == 1, keys[0] == key)

	assert(r.Key().Delete(key) == nil)
}
