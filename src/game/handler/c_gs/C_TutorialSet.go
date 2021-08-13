package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_TutorialSet(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_TutorialSet)
	plr := ctx.(*app.Player)

	res := &msg.GS_TutorialSet_R{}

	res.ErrorCode = func() int32 {
		tut := plr.GetTutorial()

		// use key
		if req.Key != "" {
			// check key
			if tut.HasKey(req.Key) {
				return Err.Tut_KeyAlreadyUsed
			}

			// perform
			if !tut_perform(plr, req.Key, res) {
				return Err.Tut_KeyNotFound
			}
		}

		// set
		if req.Step > tut.Steps[req.Tp] {
			tut.Set(req.Tp, req.Step, req.Key)
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}

// ============================================================================

func tut_perform(plr *app.Player, key string, res *msg.GS_TutorialSet_R) bool {
	switch key {

	// 首次
	case "draw-normal":
		ret := plr.GetDraw().Tutorial("normal")
		if ret == nil {
			return false
		}

		plr.SendMsg(ret)

		res.Rewards = ret.Rewards

	case "draw-senior":
		ret := plr.GetDraw().Tutorial("senior")
		if ret == nil {
			return false
		}

		plr.SendMsg(ret)

		res.Rewards = ret.Rewards

	case "guaji":
		conf := gamedata.ConfGlobalPublic.Query(1)
		if conf == nil {
			return false
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_WLevelGJ)

		for _, v := range conf.TutorialGj {
			op.Inc(v.Id, v.N)
		}

		res.Rewards = op.Apply().ToMsg()

		evtmgr.Fire(gconst.Evt_WLevelGj, plr)

	default:
		// return false
	}

	return true
}
