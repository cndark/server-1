package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/rift"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_RiftMineInfo(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_RiftMineInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_RiftMineInfo_R{}
	res.ErrorCode = func() int32 {
		mine := rift.MineMgr.Mines[req.Seq]
		if mine != nil {
			res.Mine = mine.ToMsg()
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
