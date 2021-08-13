package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	"time"
)

func C_TimeSync(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_TimeSync)
	plr := ctx.(*app.Player)

	// update heart-beat
	plr.HeartbeatUpdate()

	now := time.Now()

	plr.SendMsg(&msg.GS_TimeSync_R{
		Ts: now.Unix(),
	})
}
