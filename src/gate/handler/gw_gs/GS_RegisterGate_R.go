package gw_gs

import (
	"fw/src/core/log"
	"fw/src/gate/app"
	"fw/src/gate/msg"
)

func GS_RegisterGate_R(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_RegisterGate_R)
	gs := ctx.(*app.SocketGS)

	if req.Success {
		log.Notice("register to gs OK:", gs.Id)
	} else {
		log.Error("register to gs failed:", gs.Id)
		gs.Close()
	}
}
