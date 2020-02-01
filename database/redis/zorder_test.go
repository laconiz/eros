package redis

import (
	"testing"
)

func TestZOrder(t *testing.T) {

	z := r.ZOrder(key)

	assert(r.Key().Delete(key) == nil)

	// Incr
	assert(z.Incr(1, 2) == nil)
	assert(z.Incr(3, 4) == nil)

	// Range
	m := map[int8]int64{}
	assert(z.Range(0, 1, &m) == nil,
		len(m) == 2, m[1] == 2, m[3] == 4)

	// Score
	Int64, exists, err := z.Score(1)
	assert(err == nil, exists, Int64 == 2)
	Int64, exists, err = z.Score(2)
	assert(err == nil, !exists)

	// Rank
	Int64, exists, err = z.Rank(1)
	assert(err == nil, exists, Int64 == 1)
	Int64, exists, err = z.Rank(2)
	assert(err == nil, !exists)

	assert(r.Key().Delete(key) == nil)
	assert(z.Range(0, 1, &m) == nil, len(m) == 0)
}
