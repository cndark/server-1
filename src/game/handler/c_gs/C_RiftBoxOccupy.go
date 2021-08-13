package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/rift"
	"fw/src/game/msg"
)

func C_RiftBoxOccupy(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_RiftBoxOccupy)
	plr := ctx.(*app.Player)

	rift.BoxMgr.Occupy(plr, req.Id, func(ec int32, finTs int64, replay *msg.BattleReplay) {
		res := &msg.GS_RiftBoxOccupy_R{}
		res.ErrorCode = ec
		res.FinTs = finTs
		res.Replay = replay

		plr.SendMsg(res)
	})
}
