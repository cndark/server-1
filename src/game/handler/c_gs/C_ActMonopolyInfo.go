package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/monopoly"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActMonopolyInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_ActMonopolyInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_ActMonoPolyInfo_R{}
	res.ErrorCode = func() int32 {
		// find act
		a := act.FindAct(gconst.ActName_Monopoly)
		if a == nil {
			return Err.Act_ActNotFound
		}

		res.Data = monopoly.ActMonopolyInfo(plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
