package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/summon"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActSummonInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_ActSummonInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_ActSummonInfo_R{}
	res.ErrorCode = func() int32 {
		// find act
		a := act.FindAct(gconst.ActName_Summon)
		if a == nil {
			return Err.Act_ActNotFound
		}

		res.Data = summon.ActSummonInfo(plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
