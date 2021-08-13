package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroTrinketTransformGen(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroTrinketTransformGen)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroTrinketTransformGen_R{}

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

		// conf
		conf := gamedata.ConfTrinket.Query(hero.Trinket.Lv)
		if conf == nil {
			return Err.Failed
		}

		// cost
		op := plr.GetBag().NewOp(gconst.ObjFrom_HeroTrinketRefresh)

		for _, v := range conf.TransformCost {
			op.Dec(v.Id, v.N)
		}

		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		// cost ok
		op.Apply()

		// transform gen
		res.Props = hero.TrinketTransformGen()

		return Err.OK
	}()

	plr.SendMsg(res)
}
