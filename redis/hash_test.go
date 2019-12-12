package redis

import (
	"testing"
)

func TestHash(t *testing.T) {

	k := r.Key()
	h := r.Hash(key)

	assert(k.Delete(key, key2) == nil)

	String := ""

	// Set & Get
	assert(h.Set(key, key2) == nil)
	exists, err := h.Get(key, &String)
	assert(err == nil, exists, String == key2)

	exists, err = k.Get(key2, &String)
	assert(err == nil, !exists)

	assert(h.Set(key2, ComplexValue) == nil)
	value := &Struct{}
	exists, err = h.Get(key2, value)
	assert(err == nil, exists, value.Equal(ComplexValue))

	// 获取一组数据并反序列化到slice中
	assert(h.Set(key, ComplexValue) == nil)
	assert(h.Set(key2, ComplexPointer) == nil)

	var values []Struct
	assert(h.Gets(&values, key, key2) == nil,
		len(values) == 2,
		values[0].Equal(ComplexValue),
		values[1].Equal(*ComplexPointer),
	)

	var points []*Struct
	assert(h.Gets(&points, key, key2) == nil,
		len(points) == 2,
		points[0].Equal(ComplexValue),
		points[1].Equal(*ComplexPointer),
	)

	// 获取所有数据并反序列化到map中
	var mValues map[string]Struct
	assert(h.GetAll(&mValues) == nil,
		len(mValues) == 2,
		mValues[key].Equal(ComplexValue),
		mValues[key2].Equal(*ComplexPointer),
	)

	var pmValues map[string]*Struct
	assert(h.GetAll(&pmValues) == nil,
		len(pmValues) == 2,
		pmValues[key].Equal(ComplexValue),
		pmValues[key2].Equal(*ComplexPointer),
	)

	// 获取一组数据并反序列化到struct中
	String = "hello world"
	assert(h.Set(key, String) == nil)
	complex := &struct {
		A string
		B *Struct
	}{}
	assert(h.Gets(complex, key, key2) == nil,
		complex.A == String,
		complex.B.Equal(*ComplexPointer),
	)

	// 增量
	Int64, err := h.Incr(key, 1)
	assert(err != nil)
	assert(h.Delete(key, key2) == nil)
	Int64, err = h.Incr(key, 3)
	assert(err == nil, Int64 == 3)
	Int64, err = h.Incr(key, -4)
	assert(err == nil, Int64 == -1)

	// 无符号增量
	Int64, success, err := h.UnsignedIncr(key2, 5)
	assert(err == nil, success, Int64 == 5)
	Int64, success, err = h.UnsignedIncr(key2, -6)
	assert(err == nil, !success)
	Int64, success, err = h.UnsignedIncr(key2, -3)
	assert(err == nil, success, Int64 == 2)

	assert(k.Delete(key, key2) == nil)
}
