package statustab

import (
	"fw/src/core"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/worlddata"
	Err "fw/src/proto/errorcode"
	"time"
)

// ============================================================================
// 类型

const (
	type_svr_open_day  = "svrOpenDay" // 开服时间达到
	type_plr_level     = "level"      // 等级
	type_plr_counter   = "counter"    // counter
	type_plr_attainTab = "attainTab"  // 统计数据是否到达
	type_plr_ccyEqual  = "ccyEqual"   // 当前货币数
	type_plr_wLevel    = "wLevel"     // 主线关卡
	type_plr_vip       = "vip"        // VIP达到
	type_actState      = "actState"   // 活动状态
	type_actConfGrp    = "actConfGrp" // 活动Grp
	type_priv_card     = "privCard"   // 激活特权卡
)

// ============================================================================

// 新的处理函数添加在下面
var StatusTab = statustab_t{

	type_svr_open_day: func(body interface{}, flag int32, val float64) int32 {
		t0 := core.StartOfDay(worlddata.GetSvrCreateTs())
		now := time.Now()
		day := int32(now.Sub(t0).Hours())/24 + 1

		// 1=天数达到，2=天数之内，3=指定第多少天
		switch flag {
		case 1:
			if float64(day) >= val {
				return Err.OK
			} else {
				return Err.Common_TimeNotUp
			}
		case 2:
			if float64(day) <= val {
				return Err.OK
			} else {
				return Err.Common_TimeNotUp
			}
		case 3:
			if float64(day) == val {
				return Err.OK
			} else {
				return Err.Common_TimeNotUp
			}

		}

		return Err.Failed
	},

	type_plr_level: func(body interface{}, flag int32, val float64) int32 {
		if body == nil {
			return Err.Plr_NotFound
		}

		plr := body.(IPlayer)

		if float64(plr.GetLevel()) < val {
			return Err.Plr_LowLevel
		}

		return Err.OK
	},

	type_plr_counter: func(body interface{}, flag int32, val float64) int32 {
		if body == nil {
			return Err.Plr_NotFound
		}

		plr := body.(*app.Player)

		n := plr.GetCounter().GetRemain(flag)

		if float64(n) < val {
			return Err.NotEnoughObject(flag)
		}

		return Err.OK
	},

	type_plr_attainTab: func(body interface{}, flag int32, val float64) int32 {
		if body == nil {
			return Err.Plr_NotFound
		}

		plr := body.(*app.Player)

		if float64(plr.GetAttainObjVal(flag)) < val {
			return Err.Plr_AttainTabNotEnough
		}

		return Err.OK
	},

	type_actState: func(body interface{}, flag int32, val float64) int32 {

		a := act.FindAct(core.I32toa(int32(val)))
		if a == nil {
			return Err.Act_ActNotFound
		}

		switch a.GetStage() {
		case "open":
			if flag != 1 {
				return Err.Act_StageError
			}
		case "start":
			if flag != 2 {
				return Err.Act_StageError
			}
		case "end":
			if flag != 3 {
				return Err.Act_StageError
			}
		case "close":
			if flag != 4 {
				return Err.Act_StageError
			}
		default:
			return Err.Act_StageError
		}

		return Err.OK
	},

	type_actConfGrp: func(body interface{}, flag int32, val float64) int32 {
		if body == nil {
			return Err.Plr_NotFound
		}

		a := act.FindAct(core.I32toa(flag))
		if a == nil {
			return Err.Act_ActNotFound
		}

		if a.GetConfGrp() != int32(val) {
			return Err.Act_ConfGrp
		}

		return Err.OK
	},

	type_plr_ccyEqual: func(body interface{}, flag int32, val float64) int32 {
		if body == nil {
			return Err.Plr_NotFound
		}

		plr := body.(*app.Player)

		ba := plr.GetBag().GetCcy(flag)

		if float64(ba) < val {
			return Err.NotEnoughObject(flag)
		}

		return Err.OK
	},

	type_plr_wLevel: func(body interface{}, flag int32, val float64) int32 {
		if body == nil {
			return Err.Plr_NotFound
		}

		plr := body.(*app.Player)

		if plr.GetWLevelLvNum() < int32(val) {
			return Err.Plr_WLevelLow
		}

		return Err.OK
	},

	type_plr_vip: func(body interface{}, flag int32, val float64) int32 {
		if body == nil {
			return Err.Plr_NotFound
		}

		plr := body.(*app.Player)

		if plr.GetVipLevel() < int32(val) {
			return Err.Plr_LowVipLevel
		}

		return Err.OK

	},

	type_priv_card: func(body interface{}, flag int32, val float64) int32 {
		if body == nil {
			return Err.Plr_NotFound
		}

		plr := body.(*app.Player)

		if !plr.GetPrivCard().IsPrivCardValid(int32(val)) {
			return Err.Plr_PrivCardInValid
		}

		return Err.OK

	},
	// #todo:

}

// ============================================================================
// api
// ============================================================================

type statustab_handler_t func(body interface{}, flag int32, val float64) int32

type statustab_t map[string]statustab_handler_t

// ============================================================================

// 是否满足条件
func (self statustab_t) IsStatusTabConform(body interface{}, id int32, val float64) int32 {
	conf := gamedata.ConfStatusTab.Query(id)
	if conf == nil {
		return Err.Failed
	}

	h, ok := self[conf.Type]
	if !ok {
		return Err.Common_NoStatusHandler
	}

	return h(body, conf.Flag, val)
}

// ============================================================================
