package app

// ============================================================================

func sid_to_gateid(sid uint64) int32 {
	return int32(sid >> 41)
}
