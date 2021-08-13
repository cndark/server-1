package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ArenaFight(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ArenaFight)
	plr := ctx.(*app.Player)

	res := &msg.GS_ArenaFight_R{}
	res.ErrorCode = func() int32 {
		if !plr.IsModuleOpen(gconst.ModuleId_Arena) {
			return Err.Plr_ModuleLocked
		}

		if !plr.GetTeamMgr().IsSetTeam(gconst.TeamType_Dfd) {
			return Err.Common_NotSetTeam
		}

		if req.PlrId == plr.GetId() {
			return Err.Plr_IsSelf
		}

		// fight
		plr.GetArena().Fight(req.PlrId, req.T, req.RevengeIdx,
			func(ec int32, replay *msg.BattleReplay, rwds *msg.Rewards, addScore1, addScore2 int32) {
				res.ErrorCode = ec
				res.IsOneKey = req.IsOneKey
				res.Replay = replay
				res.Rewards = rwds
				res.AddScore1 = addScore1
				res.AddScore2 = addScore2

				plr.SendMsg(res)
			})

		return Err.OK
	}()

	if res.ErrorCode != Err.OK {
		plr.SendMsg(res)
	}
}
