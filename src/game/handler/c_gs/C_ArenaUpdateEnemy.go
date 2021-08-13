package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/arena"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ArenaUpdateEnemy(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_ArenaUpdateEnemy)
	plr := ctx.(*app.Player)

	res := &msg.GS_ArenaUpdateEnemy_R{}
	res.ErrorCode = func() int32 {

		res.Enemies = arena.ArenaMgr.PickEnemies(plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
