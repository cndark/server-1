package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/warcup"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WarCupTop1Info(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WarCupTop1Info)
	plr := ctx.(*app.Player)

	res := &msg.GS_WarCupTop1Info_R{}
	res.ErrorCode = func() int32 {

		res.Plr, res.HeroId, res.HeroStar, res.HeroSkin = warcup.WarCupTop1Info()

		return Err.OK
	}()
	plr.SendMsg(res)
}
