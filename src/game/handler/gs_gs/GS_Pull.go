package gs_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func GS_Pull(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_Pull)

	ret := func(ec int32, r interface{}) {
		// marshal result
		var data []byte
		var err error

		if ec == Err.OK && r != nil {
			data, err = utils.MarshalArg(r)
			if err != nil {
				ec = Err.Failed
			}
		}

		// send msg
		app.NetMgr.Send2Game(req.SrcSvrId, &msg.GS_Pull_R{
			Tk:        req.Tk,
			ErrorCode: ec,
			Result:    data,
		})
	}

	evtmgr.Fire(req.EvtName, req.SArg, req.OArg, ret)
}
