package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

func C_GuildApply(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildApply)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildApply_R{}

	res.ErrorCode = func() int32 {
		if !plr.IsModuleOpen(gconst.ModuleId_Guild) {
			return Err.Guild_Locked
		}

		// get guild
		gld := guild.GuildMgr.FindGuild(req.GuildId)
		if gld == nil {
			return Err.Guild_NotFound
		}

		conf := gamedata.ConfGlobalPublic.Query(1)
		if conf == nil {
			return Err.Failed
		}

		if time.Now().Before(plr.GetMisc().GldLeaveTs) {
			return Err.Guild_JoinTimeLimited
		}

		// apply
		if ec := gld.ApplyRequst(plr); ec != Err.OK {
			return ec
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
