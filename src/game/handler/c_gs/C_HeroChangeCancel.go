package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroChangeCancel(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroChangeCancel)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroChangeCancel_R{}
	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.Seq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		hero.SetChangeId(0)

		// ok
		return Err.OK
	}()

	plr.SendMsg(res)
}
