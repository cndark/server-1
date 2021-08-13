package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_RelicEquip(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_RelicEquip)
	plr := ctx.(*app.Player)

	res := &msg.GS_RelicEquip_R{}

	res.ErrorCode = func() int32 {
		// find hero
		hero := plr.GetBag().FindHero(req.HeroSeq)
		if hero == nil {
			return Err.Hero_NotFound
		}

		if hero.Lv < 100 {
			return Err.Hero_LowLevel
		}

		// equip
		rlc := plr.GetBag().FindRelic(req.Seq)
		if rlc == nil {
			return Err.Equip_NotFound
		}

		return hero.NewEquipOp().Equip(rlc).Apply()
	}()

	plr.SendMsg(res)
}
