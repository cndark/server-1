package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/rift"
	"fw/src/game/msg"
)

func C_RiftBoxInfo(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_RiftBoxInfo)
	plr := ctx.(*app.Player)

	res := &msg.GS_RiftBoxInfo_R{}

	res.BoxNum = int32(len(rift.BoxMgr.Boxes))

	plr.SendMsg(res)
}
