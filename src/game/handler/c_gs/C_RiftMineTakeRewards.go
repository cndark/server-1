package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/rift"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_RiftMineTakeRewards(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_RiftMineTakeRewards)
	plr := ctx.(*app.Player)

	res := &msg.GS_RiftMineTakeRewards_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := rift.MineMgr.TakeRewards(plr, req.Seq)
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds

		return Err.OK
	}()

	plr.SendMsg(res)
}
