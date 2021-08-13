package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ArmorEquipOnekey(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ArmorEquipOnekey)
	plr := ctx.(*app.Player)

	res := &msg.GS_ArmorEquipOnekey_R{}

	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.HeroSeq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		// equip all
		eop := hero.NewEquipOp()

		for _, id := range req.Ids {
			ar := plr.GetBag().PrepareArmor(id)
			if ar == nil {
				return Err.Equip_NotFound
			}

			eop.Equip(ar)
		}

		return eop.Apply()
	}()

	plr.SendMsg(res)
}
