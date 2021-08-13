package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ArmorUnequip(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ArmorUnequip)
	plr := ctx.(*app.Player)

	res := &msg.GS_ArmorUnequip_R{}

	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.HeroSeq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		// unequip
		return hero.NewEquipOp().Unequip(req.Slot).Apply()
	}()

	plr.SendMsg(res)
}
