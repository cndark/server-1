package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

func C_PushGiftSetCreateTs(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_PushGiftSetCreateTs)
	plr := ctx.(*app.Player)

	res := &msg.GS_PushGiftSetCreateTs_R{}
	res.ErrorCode = func() int32 {

		gift := plr.GetPushGift().Get(req.Id)
		gift.BuyCnt = 0
		gift.CreateTs = time.Now()

		res.Id = req.Id
		res.CreateTs = gift.CreateTs.Unix()

		return Err.OK
	}()

	plr.SendMsg(res)
}
