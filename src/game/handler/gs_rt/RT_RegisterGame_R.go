package gs_rt

import (
	"fw/src/core/log"
	"fw/src/game/app"
	"fw/src/game/msg"
)

func RT_RegisterGame_R(message msg.Message, ctx interface{}) {
	req := message.(*msg.RT_RegisterGame_R)
	rt := ctx.(*app.SocketRt)

	if req.Success {
		log.Notice("register to router OK", rt.Id)
	} else {
		log.Error("register to router failed", rt.Id)
		rt.Close()
	}
}
