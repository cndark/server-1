package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_UserExtInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_UserExtInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_UserExtInfo_R{}

	// bill
	res.GS_BillInfo_R = plr.GetBill().ToMsg_Info()

	// act
	res.GS_ActStateGet_R = &msg.GS_ActStateGet_R{
		Acts: plr.GetAct().ToMsg(),
	}

	// crusade
	res.GS_CrusadeGetInfo_R = plr.GetCrusade().GetInfo()

	// gld full info
	gld := plr.GetGuild()
	if gld != nil {
		res.GS_GuildInfoFull_R = &msg.GS_GuildInfoFull_R{
			Info: gld.ToMsg_InfoFull(plr),
		}
	}

	// gwar
	res.GS_GWarGetSummary_R = plr.GetGWar().GetSummary()

	plr.SendMsg(res)
}
