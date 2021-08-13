package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/warcup"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WarCupGetReplay(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_WarCupGetReplay)
	plr := ctx.(*app.Player)

	res := &msg.GS_WarCupGetReplay_R{}
	res.ErrorCode = func() int32 {

		res.Replay = warcup.WarCupGetReplay(req.VsSeq)
		if res.Replay == nil {
			return Err.Common_BattleNotFound
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
