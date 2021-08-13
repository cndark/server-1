package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ArmorEquip(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ArmorEquip)
	plr := ctx.(*app.Player)

	res := &msg.GS_ArmorEquip_R{}

	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.HeroSeq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		// equip
		ar := plr.GetBag().PrepareArmor(req.Id)
		if ar == nil {
			return Err.Equip_NotFound
		}

		return hero.NewEquipOp().Equip(ar).Apply()
	}()

	plr.SendMsg(res)
}
