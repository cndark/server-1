package gs_rt

import (
	"fw/src/router/app"
	"fw/src/router/msg"
)

func GS_RegisterGame(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_RegisterGame)
	c := ctx.(*app.SocketC)

	b := app.NetMgr.RegisterC(c, req.Id)
	c.SendMsg(&msg.RT_RegisterGame_R{
		Success: b,
	})
}
