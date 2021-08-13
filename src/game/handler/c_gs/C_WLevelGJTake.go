package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WLevelGJTake(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_WLevelGJTake)
	plr := ctx.(*app.Player)

	res := &msg.GS_WLevelGJTake_R{}
	res.ErrorCode = func() int32 {
		if ec := plr.GetBag().CheckFull(); ec != Err.OK {
			return ec
		}

		plr.GetWLevel().GJLoot()

		ec, rwds := plr.GetWLevel().GJLootTake()
		if ec != Err.OK {
			return ec
		}

		res.Rewards = rwds
		res.GJTs = plr.GetWLevel().GJTs.Unix()

		return Err.OK
	}()

	plr.SendMsg(res)
}
