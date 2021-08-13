package gs_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/msg"
)

func GS_Push(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_Push)

	evtmgr.Fire(req.EvtName, req.SArg, req.OArg)
}
