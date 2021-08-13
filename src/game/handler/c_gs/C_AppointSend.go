package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_AppointSend(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_AppointSend)
	plr := ctx.(*app.Player)

	res := &msg.GS_AppointSend_R{}
	res.ErrorCode = func() int32 {

		if !plr.IsModuleOpen(gconst.ModuleId_Appoint) {
			return Err.Plr_ModuleLocked
		}

		ec, ts := plr.GetAppoint().Send(req.Seq, req.Heroes)
		if ec != Err.OK {
			return ec
		}

		res.Seq = req.Seq
		res.Heroes = req.Heroes
		res.Ts = ts

		return Err.OK
	}()

	plr.SendMsg(res)
}
