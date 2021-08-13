package gw_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func GW_LogoutPlayer(message msg.Message, ctx interface{}) {
	req := message.(*msg.GW_LogoutPlayer)

	plr := app.PlayerMgr.FindPlayerBySid(req.Sid)
	if plr != nil {
		app.PlayerMgr.SetOffline(plr, false)
	}
}
