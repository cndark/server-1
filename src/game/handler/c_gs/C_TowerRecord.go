package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/tower"
	"fw/src/game/msg"
)

func C_TowerRecord(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_TowerRecord)
	plr := ctx.(*app.Player)

	res := &msg.GS_TowerRecord_R{}

	rec := tower.TowerMgr.Records[req.LvNum]
	if rec != nil {
		res.First = rec.First
		res.MinPower = rec.MinPower
	}

	plr.SendMsg(res)
}
