package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildKick(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildKick)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildKick_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotFound
		}

		// find operator member
		m_oper := gld.FindMember(plr.GetId())
		if m_oper == nil {
			return Err.Guild_NotAMember
		}

		// check operator rank
		if m_oper.Rank > guild.RK_Vice {
			return Err.Guild_LowPriv
		}

		// find plr member
		m_plr := gld.FindMember(req.PId)
		if m_plr == nil {
			return Err.Guild_NotAMember
		}

		// check plr rank
		if m_oper.Rank >= m_plr.Rank {
			return Err.Guild_LowPriv
		}

		// kick
		if ec := m_plr.Leave(guild.LR_Kick); ec != Err.OK {
			return ec
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
