package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/warcup"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WarCupTop8Info(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WarCupTop8Info)
	plr := ctx.(*app.Player)

	res := &msg.GS_WarCupTop8Info_R{}
	res.ErrorCode = func() int32 {

		res.VsData = warcup.WarCupTop8Info()

		return Err.OK
	}()

	plr.SendMsg(res)
}
