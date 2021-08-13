package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_FriendGetApplyList(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_FriendGetApplyList)
	plr := ctx.(*app.Player)
	res := &msg.GS_FriendGetApplyList_R{}

	for plrid := range plr.GetFriend().ApplyList {
		fplr := app.PlayerMgr.FindPlayerById(plrid)
		if fplr == nil {
			continue
		}

		res.ApplyList = append(res.ApplyList, fplr.ToMsg_SimpleInfo())
	}

	plr.SendMsg(res)
}
