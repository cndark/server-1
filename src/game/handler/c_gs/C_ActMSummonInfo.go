package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/msummon"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActMSummonInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_ActMSummonInfo)

	plr := ctx.(*app.Player)

	res := &msg.GS_ActMSummonInfo_R{}
	res.ErrorCode = func() int32 {

		a := act.FindAct(gconst.ActName_Summon)
		if a == nil {
			return Err.Act_ActNotFound
		}

		res.Data = msummon.ActMSummonInfo(plr)

		return Err.OK
	}()

	plr.SendMsg(res)

}
