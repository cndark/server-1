package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildTechLevelup(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildTechLevelup)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildTechLevelup_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotAMember
		}

		// conf
		conf := gamedata.ConfGuildTech.Query(req.Id)
		if conf == nil {
			return Err.Failed
		}

		// get tech object
		tech_obj := plr.GetGuildPlrData().Tech

		// check pre cond
		if len(conf.PreId) > 0 {
			pre_lv := tech_obj.Techs[conf.PreId[0].Id]
			if pre_lv < conf.PreId[0].Level {
				return Err.Guild_TechPreCond
			}
		}

		// full level ?
		lv := tech_obj.Techs[req.Id]
		if lv >= conf.LevelMax {
			return Err.Guild_TechFullLevel
		}

		// cost
		op := plr.GetBag().NewOp(gconst.ObjFrom_GuildTechLevelup)

		// init cost
		for _, v := range conf.UpCost {
			op.Dec(v.Id, v.N)
		}

		// growth cost
		for _, v := range conf.UpCostGrowth {
			op.Dec(v.Id, v.N*int64(lv))
		}

		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		// ok
		op.Apply()

		// set level
		tech_obj.SetLevel(req.Id, lv+1)

		return Err.OK
	}()

	plr.SendMsg(res)
}
