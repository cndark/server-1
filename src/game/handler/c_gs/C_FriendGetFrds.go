package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

func C_FriendGetFrds(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_FriendGetFrds)
	plr := ctx.(*app.Player)

	res := &msg.GS_FriendGetFrds_R{}
	frd := plr.GetFriend()

	for plrid := range frd.Friends {
		fplr := app.PlayerMgr.LoadPlayer(plrid)
		if fplr != nil {
			res.Friends = append(res.Friends, plr.GetFriend().FriendInfo_ToMsg(plrid))
		}
	}

	plr.SendMsg(res)
}
