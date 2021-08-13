package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_RelicUnequip(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_RelicUnequip)
	plr := ctx.(*app.Player)

	res := &msg.GS_RelicUnequip_R{}

	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.HeroSeq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		// unequip
		return hero.NewEquipOp().Unequip(gconst.EquipGroup_Relic_SlotStart).Apply()
	}()

	plr.SendMsg(res)
}
