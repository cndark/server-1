package gw_gs

import (
	"fw/src/gate/app"
	"fw/src/gate/msg"
)

func GS_Kick(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_Kick)

	sess := app.NetMgr.FindSession(req.Sid)
	if sess != nil {
		sess.Close()
	}
}
