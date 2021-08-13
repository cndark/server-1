package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_HeroTrinketUp(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_HeroTrinketUp)
	plr := ctx.(*app.Player)

	res := &msg.GS_HeroTrinketUp_R{}

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
		if conf == nil || len(conf.UpCost) == 0 {
			return Err.Hero_TrinketFullLevel
		}

		// cost
		op := plr.GetBag().NewOp(gconst.ObjFrom_HeroTrinketUp)

		for _, v := range conf.UpCost {
			op.Dec(v.Id, v.N)
		}

		if req.Lock {
			for _, v := range conf.LockCost {
				op.Dec(v.Id, v.N)
			}
		}

		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		// const ok
		op.Apply()

		// level up
		hero.TrinketSetLevel(hero.Trinket.Lv+1, req.Lock)

		return Err.OK
	}()

	plr.SendMsg(res)
}
