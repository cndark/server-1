package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/summon"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActSummonPick(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ActSummonPick)
	plr := ctx.(*app.Player)

	res := &msg.GS_ActSummonPick_R{}
	res.ErrorCode = func() int32 {
		// find act
		a := act.FindAct(gconst.ActName_Summon)
		if a == nil {
			return Err.Act_ActNotFound
		}

		ec := summon.Pick(plr, req.HeroPos)
		if ec != Err.OK {
			return ec
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
