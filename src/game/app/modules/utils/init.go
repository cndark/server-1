package utils

// ============================================================================

func Init(plrmgr IPlayerMgr, netmgr INetMgr, statustab IStatusTab) {
	iplrmgr = plrmgr
	inetmgr = netmgr
	istatustab = statustab
}
