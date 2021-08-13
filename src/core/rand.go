package core

import "math/rand"

// ============================================================================

// [a, b]
func RandInt(a, b int) int {
	return rand.Intn(b-a+1) + a
}

// [a, b]
func RandInt32(a, b int32) int32 {
	return rand.Int31n(b-a+1) + a
}

// [a, b)
func RandFloat32(a, b float32) float32 {
	return rand.Float32()*(b-a) + a
}

// [a, b)
func RandFloat64(a, b float64) float64 {
	return rand.Float64()*(b-a) + a
}
