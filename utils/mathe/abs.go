package mathe

import "math"

func AbsInt8(n int8) int8 {
	return n &^ 1 << 7
}

func AbsInt16(n int16) int16 {
	return n &^ 1 << 15
}

func AbsInt32(n int32) int32 {
	return n &^ 1 << 31
}

func AbsInt64(n int64) int64 {
	return n &^ 1 << 63
}

func AbsFloat32(f float32) float32 {
	return math.Float32frombits(math.Float32bits(f) &^ 1 << 31)
}
