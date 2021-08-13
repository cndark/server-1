package math

// ============================================================================

func MinInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func MinInt32(a, b int32) int32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt32(a, b int32) int32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func MinFloat32(a, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxFloat32(a, b float32) float32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func MinFloat64(a, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxFloat64(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func AbsFloat32(a, b float32) float32 {
	v := a - b
	if v >= 0 {
		return v
	} else {
		return -v
	}
}

func AbsInt32(a, b int32) int32 {
	v := a - b
	if v >= 0 {
		return v
	} else {
		return -v
	}
}

func SignInt32(v int32) int32 {
	if v > 0 {
		return 1
	} else if v < 0 {
		return -1
	} else {
		return 0
	}
}

func SignInt64(v int64) int32 {
	if v > 0 {
		return 1
	} else if v < 0 {
		return -1
	} else {
		return 0
	}
}
