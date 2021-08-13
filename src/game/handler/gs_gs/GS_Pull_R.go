package gs_gs

import (
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
)

func GS_Pull_R(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_Pull_R)

	utils.GSPull_R(req.Tk, req.ErrorCode, req.Result)
}
