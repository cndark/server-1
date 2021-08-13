package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/warcup"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WarCupGuess(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_WarCupGuess)
	plr := ctx.(*app.Player)

	res := &msg.GS_WarCupGuess_R{}
	res.ErrorCode = func() int32 {

		ec, gscore := warcup.WarCupGuess(plr, req.GuessWin, req.GuessNum)
		if ec != Err.OK {
			return ec
		}

		res.GuessScore = gscore

		return Err.OK
	}()

	plr.SendMsg(res)
}
