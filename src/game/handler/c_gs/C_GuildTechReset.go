package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildTechReset(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_GuildTechReset)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildTechReset_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotAMember
		}

		// get tech object
		tech_obj := plr.GetGuildPlrData().Tech

		// calc return res
		cost := make(map[int32]int64)

		for id, lv := range tech_obj.Techs {
			// conf
			conf := gamedata.ConfGuildTech.Query(id)
			if conf == nil {
				return Err.Failed
			}

			// nb + g*(n-1)*n/2
			for _, v := range conf.UpCost {
				cost[v.Id] += int64(lv) * v.N
			}
			for _, v := range conf.UpCostGrowth {
				cost[v.Id] += v.N * int64(lv-1) * int64(lv) / 2
			}
		}

		// cost
		cop := plr.GetCounter().NewOp(gconst.ObjFrom_GuildTechReset)
		cop.DecCounter(gconst.Cnt_GuildTechReset, 1)

		if ec := cop.CheckEnough(); ec != Err.OK {
			return ec
		}

		// return res
		for id, v := range cost {
			cop.Inc(id, v)
		}
		res.Rewards = cop.Apply().ToMsg()

		// reset
		tech_obj.Reset()

		return Err.OK
	}()

	plr.SendMsg(res)
}
