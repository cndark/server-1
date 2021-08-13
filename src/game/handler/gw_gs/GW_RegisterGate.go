package gw_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func GW_RegisterGate(message msg.Message, ctx interface{}) {
	req := message.(*msg.GW_RegisterGate)
	gw := ctx.(*app.SocketGW)

	b := app.NetMgr.RegisterGate(gw, req.Id)
	gw.SendMsg(&msg.GS_RegisterGate_R{
		Success: b,
	})
}
