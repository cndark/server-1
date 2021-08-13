package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroSetSkin(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroSetSkin)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroSetSkin_R{}
	res.ErrorCode = func() int32 {

		hero := plr.GetBag().FindHero(req.Seq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		if req.Skin != 0 {
			conf := gamedata.ConfHeroSkin.Query(req.Skin)
			if conf == nil || conf.Hero != hero.Id {
				return Err.Failed
			}

			if !plr.GetHeroSkin().IsExist(req.Skin) {
				return Err.Hero_SkinNotExist
			}
		}

		hero.SetSkin(req.Skin)

		return Err.OK
	}()

	plr.SendMsg(res)
}
