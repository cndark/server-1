package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/rift"
	"fw/src/game/msg"
)

func C_RiftMineOccupy(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_RiftMineOccupy)
	plr := ctx.(*app.Player)

	rift.MineMgr.Occupy(plr, req.Seq, req.T, func(ec int32, mine *msg.RiftMine, replay *msg.BattleReplay) {
		res := &msg.GS_RiftMineOccupy_R{}
		res.ErrorCode = ec
		res.Mine = mine
		res.Replay = replay

		plr.SendMsg(res)
	})
}
