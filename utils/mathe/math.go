package mathe

import "math"

func MaxInt8(a, b int8) int8 {
	if a > b {
		return a
	}
	return b
}

func MaxInt16(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}

func MaxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MaxUint8(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

func MaxUint16(a, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}

func MaxUint32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func MaxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func MaxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

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
