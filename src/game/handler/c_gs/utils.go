package c_gs

// ============================================================================

func battleid_to_batid(battleid int64) int32 {
	return int32(battleid >> 41)
}
