package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroLock(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroLock)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroLock_R{}

	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.Seq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		// lock
		hero.Locked = req.Lock

		// ok
		return Err.OK
	}()

	plr.SendMsg(res)
}
