package crusade

import (
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/mdata"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"time"
)

// ============================================================================

var CrusadeMgr *crusade_mgr_t

// ============================================================================

type crusade_mgr_t struct {
	VerTs time.Time // 当前版本时间(t0)
	Ts1   time.Time // 当前阶段开始时间(t1)
	Ts2   time.Time // 当前阶段结束时间(t2)
	stage string    // 当前stage
}

// ============================================================================

func new_data() interface{} {
	return &crusade_mgr_t{}
}

func data_loaded() {
	CrusadeMgr = mdata.Get(NAME).(*crusade_mgr_t)
}

// ============================================================================

func on_stage(stg string, t0, t1, t2 time.Time, f func(bool)) {
	if CrusadeMgr == nil {
		CrusadeMgr = mdata.Get(NAME).(*crusade_mgr_t)
	}

	if CrusadeMgr.VerTs.IsZero() || !CrusadeMgr.VerTs.Equal(t0) {
		CrusadeMgr.VerTs = t0
	}

	CrusadeMgr.Ts1 = t1
	CrusadeMgr.Ts2 = t2
	CrusadeMgr.stage = stg

	f(true)
}

func stage_start() {
	utils.BroadcastPlayers(&msg.GS_CrusadeStageUpdate{
		Stage: CrusadeMgr.stage,
		Ts1:   CrusadeMgr.Ts1.Unix(),
		Ts2:   CrusadeMgr.Ts2.Unix(),
	})
}

func stage_end() {
	utils.BroadcastPlayers(&msg.GS_CrusadeStageUpdate{
		Stage: CrusadeMgr.stage,
		Ts1:   CrusadeMgr.Ts1.Unix(),
		Ts2:   CrusadeMgr.Ts2.Unix(),
	})

	// box_award()
}

func IsStart() bool {
	return CrusadeMgr.stage == "start"
}

// ============================================================================

// 没有领取的宝箱发邮件
func box_award() {
	ids := []int32{}
	for _, v := range gamedata.ConfCrusade.Items() {
		if len(v.Chest) > 0 {
			ids = append(ids, v.Id)
		}
	}

	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil {
		return
	}

	utils.ForEachLoadedPlayer(func(plr interface{}) {
		mplr := plr.(IPlayer)
		if !mplr.IsModuleOpen(gconst.ModuleId_Crusade) {
			return
		}

		if mplr.GetCrusade().LvNum == 0 {
			return
		}

		m := mail.New(mplr).SetKey(conf_g.CrusadeRewardMailId)

		isSend := false
		for _, id := range ids {
			if mplr.GetCrusade().LvNum >= id && !mplr.GetCrusade().IsBoxTake(id) {
				conf_c := gamedata.ConfCrusade.Query(id)
				if conf_c != nil {
					for _, v := range conf_c.Chest {
						m.AddAttachment(v.Id, float64(v.N))
						isSend = true
					}

					mplr.GetCrusade().BoxTaken = append(mplr.GetCrusade().BoxTaken, id)
				}
			}
		}

		if isSend {
			m.Send()
		}
	})
}
