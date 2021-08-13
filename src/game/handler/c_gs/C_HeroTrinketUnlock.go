package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroTrinketUnlock(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroTrinketUnlock)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroTrinketUnlock_R{}

	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.Seq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		// unlock
		return hero.TrinketUnlock()
	}()

	plr.SendMsg(res)
}
