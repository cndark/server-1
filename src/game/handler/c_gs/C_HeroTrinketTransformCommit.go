package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroTrinketTransformCommit(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroTrinketTransformCommit)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroTrinketTransformCommit_R{}

	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.Seq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		// check if unlocked
		if hero.Trinket.Lv == 0 {
			return Err.Hero_TrinketNotUnlocked
		}

		// transform commit
		return hero.TrinketTransformCommit()
	}()

	plr.SendMsg(res)
}
