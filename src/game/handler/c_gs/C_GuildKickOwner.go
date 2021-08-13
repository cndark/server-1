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

func C_GuildKickOwner(message msg.Message, ctx interface{}) {
	//req := message.(*msg.C_GuildKickOwner)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildKickOwner_R{}

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

		// find Owner
		m_owner := gld.Owner()
		if m_owner == nil {
			return Err.Guild_NotAMember
		}

		// can NOT operate on self
		if m_owner == m_oper {
			return Err.Guild_OperateOnSelf
		}

		own_plr := app.PlayerMgr.LoadPlayer(m_owner.Id)
		if own_plr == nil {
			return Err.Plr_NotLoad
		}

		conf := gamedata.ConfGlobalPublic.Query(1)
		if conf == nil {
			return Err.Failed
		}

		lts := own_plr.GetOfflineTs().Add(time.Duration(conf.GuildKickHour) * time.Hour)
		if lts.After(time.Now()) {
			return Err.Guild_KickOwnerTimeError
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_GuildKick)
		for _, v := range conf.GuildKickCost {
			op.Dec(v.Id, int64(v.N))
		}

		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		// set rank
		if ec := m_oper.SetRank(guild.RK_Owner); ec != Err.OK {
			return ec
		}

		op.Apply()

		gld.AddLog(guild.C_Guild_Log_KickOwner, map[string]string{
			"player": plr.GetName(),
		})

		return Err.OK
	}()

	plr.SendMsg(res)
}
