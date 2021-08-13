package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_MiscSkipTutorial(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_MiscSkipTutorial)
	plr := ctx.(*app.Player)

	res := &msg.GS_MiscSkipTutorial_R{}

	func() int {
		if plr.GetLevel() > 1 {
			return Err.Failed
		}

		conf := gamedata.ConfInitial.Query(2)
		if conf == nil {
			return Err.Failed
		}

		// add stuffs
		op := plr.GetBag().NewOp(gconst.ObjFrom_Init)

		// ccy
		for _, v := range conf.InitialCurrency {
			op.Inc(v.Id, int64(v.N))
		}

		// heroes
		for _, v := range conf.InitialHero {
			op.Inc(v, 1)
		}

		// items
		for _, v := range conf.InitialItem {
			op.Inc(v.Id, int64(v.N))
		}

		op.Apply()

		return Err.OK
	}()

	plr.SendMsg(res)
}
