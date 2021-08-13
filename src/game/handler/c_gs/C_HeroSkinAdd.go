package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroSkinAdd(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroSkinAdd)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroSkinAdd_R{}
	res.ErrorCode = func() int32 {
		ec := plr.GetHeroSkin().AddSkin(req.Skin)
		if ec != Err.OK {
			return ec
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
