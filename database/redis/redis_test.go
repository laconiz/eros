package redis

import "fmt"

func assert(assertions ...bool) {
	for i, a := range assertions {
		if !a {
			panic(fmt.Errorf("assert index %d failed", i))
		}
	}
}

const (
	key  = "redis.test.key"
	key2 = "redis.test.key2"
)

var r *Redis

type Struct struct {
	A int64
	B string
	C []float32
}

func (s Struct) Equal(v Struct) bool {
	return s.A == v.A && s.B == v.B
}

var ComplexValue = Struct{A: 100, B: "complex value", C: []float32{1.1, 2.2, 3.3}}
var ComplexPointer = &Struct{A: 200, B: "complex pointer", C: []float32{4.4, 5.5}}

func init() {

	conf := Config{
		Network:   "tcp",
		Address:   "127.0.0.1:6379",
		Password:  "redis",
		Database:  15,
		MaxIdle:   5,
		MaxActive: 50,
		// Log:       true,
	}

	var err error
	if r, err = New(conf); err != nil {
		panic(err)
	}
}
