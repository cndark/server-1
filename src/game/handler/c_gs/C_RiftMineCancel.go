package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/rift"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_RiftMineCancel(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_RiftMineCancel)
	plr := ctx.(*app.Player)

	res := &msg.GS_RiftMineCancel_R{}
	res.ErrorCode = func() int32 {

		ec := rift.MineMgr.CancelOccupy(plr, req.Seq)
		if ec != Err.OK {
			return ec
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
