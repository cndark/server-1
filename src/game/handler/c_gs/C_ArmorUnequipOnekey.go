package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ArmorUnequipOnekey(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ArmorUnequipOnekey)
	plr := ctx.(*app.Player)

	res := &msg.GS_ArmorUnequipOnekey_R{}

	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.HeroSeq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		// unequip all
		return hero.NewEquipOp().UnequipAll(gconst.EquipGroup_Armor).Apply()
	}()

	plr.SendMsg(res)
}
