package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroChangeApply(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroChangeApply)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroChangeApply_R{}
	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.Seq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		if hero.Locked {
			return Err.Hero_Locked
		}

		conf_h := gamedata.ConfMonster.Query(hero.ChangeId)
		if conf_h == nil {
			return Err.Failed
		}

		hero.ChangeApply()

		res.Hero = hero.ToMsg()

		return Err.OK
	}()

	plr.SendMsg(res)
}
