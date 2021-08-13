package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_FriendGetBlackList(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_FriendGetBlackList)
	plr := ctx.(*app.Player)
	res := &msg.GS_FriendGetBlackList_R{}

	for plrid := range plr.GetFriend().BlackList {
		fplr := app.PlayerMgr.FindPlayerById(plrid)
		if fplr == nil {
			continue
		}

		res.BlackList = append(res.BlackList, plr.ToMsg_SimpleInfo())
	}

	plr.SendMsg(res)
}
