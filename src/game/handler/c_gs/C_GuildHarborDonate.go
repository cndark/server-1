package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildHarborDonate(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildHarborDonate)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildHarborDonate_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotAMember
		}

		// conf
		conf := gamedata.ConfGlobalPublic.Query(1)
		if conf == nil {
			return Err.Failed
		}

		// donate mode
		v := int64(conf.HarborDonateValue)
		if req.Which != 1 {
			v *= 10
		}

		// cost
		op := plr.GetBag().NewOp(gconst.ObjFrom_GuildHarborDonate)
		op.Dec(gconst.Gold, v)
		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		// ok
		op.Apply()

		// donate
		gld.Harbor.AddDonate(plr.GetId(), v)

		evtmgr.Fire(gconst.Evt_GuildDonate, plr)

		return Err.OK
	}()

	plr.SendMsg(res)
}
