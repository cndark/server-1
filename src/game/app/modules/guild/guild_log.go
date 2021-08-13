package guild

import (
	"time"
)

// ============================================================================

const (
	C_Guild_Log_Join      = 1 // 加入公会
	C_Guild_Log_Leave     = 2 // 离开公会
	C_Guild_Log_SetRank   = 3 // 任命
	C_Guild_Log_KickOwner = 4 // 踢会长
	C_Guild_Log_Create    = 5 // 创建

)

// ============================================================================

// 日志
type log_t struct {
	Id    int32
	Param map[string]string
	Ts    time.Time
}

// ============================================================================

func (self *Guild) AddLog(id int32, param map[string]string) {
	L := len(self.Log)
	if L >= 100 {
		self.Log = self.Log[L-100:]
	}

	self.Log = append(self.Log, &log_t{
		Id:    id,
		Param: param,
		Ts:    time.Now(),
	})
}
