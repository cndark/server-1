package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/act/modules/heroskin"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ActHeroSkinInfo(message msg.Message, ctx interface{}) {
	plr := ctx.(*app.Player)

	res := &msg.GS_ActHeroSkinInfo_R{}
	res.ErrorCode = func() int32 {
		// find act
		a := act.FindAct(gconst.ActName_HeroSkin)
		if a == nil {
			return Err.Act_ActNotFound
		}

		res.Data = heroskin.ActHeroSkinInfo(plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
