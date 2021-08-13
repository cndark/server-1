package c_gs

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
)

// ============================================================================
func init() {
	evtmgr.On(gconst.Evt_GsPull_PlrInfo, func(args ...interface{}) {
		sarg := args[0].([]string)
		ret := args[2].(func(int32, interface{}))

		if len(sarg) < 2 {
			ret(Err.Failed, nil)
			return
		}

		data := ToMsg_PlayerInfo_R(sarg[0], core.Atoi32(sarg[1]))

		ret(Err.OK, data)
	})
}

// ============================================================================

// 玩家数据(本服或跨服)
func ToMsg_PlayerInfo_R(plrid string, teamType int32) *msg.GS_PlayerInfo_R {
	ret := &msg.GS_PlayerInfo_R{
		DfdTeamInfo: make(map[int32]*msg.HeroDetail),
	}

	plr := app.PlayerMgr.FindPlayerById(plrid)
	if plr == nil {
		ret.ErrorCode = Err.Plr_NotLoad
		return ret
	}

	ret.Info = plr.ToMsg_SimpleInfo()
	ret.GuildId = plr.GetGuildId()
	if ret.GuildId != "" {
		ret.GuildName = plr.GetGuildName()
		ret.GuildIcon = plr.GetGuildIcon()
	}

	t := plr.GetTeam(teamType)
	if t != nil {
		for seq, pos := range t.Formation {
			hero := plr.GetBag().FindHero(seq)
			if hero != nil {
				ret.DfdTeamInfo[pos] = hero.ToMsg_Detail()
			}
		}
	}

	return ret
}

// ============================================================================

func C_PlayerInfo(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_PlayerInfo)
	plr := ctx.(*app.Player)

	if req.SvrId == config.CurGame.Id {
		res := ToMsg_PlayerInfo_R(req.PlrId, req.TeamType)
		plr.SendMsg(res)
	} else {
		utils.GsPull(req.SvrId, gconst.Evt_GsPull_PlrInfo, []string{req.PlrId, core.I32toa(req.TeamType)},
			nil, &msg.GS_PlayerInfo_R{},
			func(ec int32, r interface{}) {
				if ec != Err.OK {
					plr.SendMsg(&msg.GS_PlayerInfo_R{
						ErrorCode: ec,
					})
				} else {
					plr.SendMsg(r.(*msg.GS_PlayerInfo_R))
				}
			})
	}
}
