package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/warcup"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WarCupGuessInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WarCupGuessInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_WarCupGuessInfo_R{}
	res.ErrorCode = func() int32 {
		if !warcup.IsOpen() {
			return Err.Common_TimeNotUp
		}

		res.Guess = warcup.WarCupGuessInfo(plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
