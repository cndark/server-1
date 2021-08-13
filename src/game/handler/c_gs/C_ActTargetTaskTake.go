package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/targettask"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActTargetTaskTake(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActTargetTaskTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_ActTargetTaskTake_R{}
	res.ErrorCode = func() int32 {
		// find act
		a := act.FindAct(gconst.ActName_TargetTask)
		if a == nil {
			return Err.Act_ActNotFound
		}

		ec, rwds := targettask.TakeRewards(plr, req.Id)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
