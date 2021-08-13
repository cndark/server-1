package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/modules/battle"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ArenaReplayGet(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ArenaReplayGet)
	plr := ctx.(*app.Player)

	battle.ReplayGet(dbmgr.DBGame, dbmgr.C_tabname_replayarena, req.ReplayId, func(r *msg.BattleReplay) {
		res := &msg.GS_ArenaReplayGet_R{}
		if r == nil {
			res.ErrorCode = Err.Common_BattleNotFound
		}

		res.Replay = r
		plr.SendMsg(res)
	})
}
