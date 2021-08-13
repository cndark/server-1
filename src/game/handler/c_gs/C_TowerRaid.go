package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_TowerRaid(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_TowerRaid)
	plr := ctx.(*app.Player)

	res := &msg.GS_TowerRaid_R{}
	res.ErrorCode = func() int32 {

		ec, rwds := plr.GetTower().Raid()
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds
		res.LastTs = plr.GetTower().LastTs.Unix()

		return Err.OK
	}()

	plr.SendMsg(res)
}
