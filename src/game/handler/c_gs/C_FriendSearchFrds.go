package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"

	Err "fw/src/proto/errorcode"
)

func C_FriendSearchFrds(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_FriendSearchFrds)
	plr := ctx.(*app.Player)

	res := &msg.GS_FriendSearchFrds_R{}
	func() int32 {
		if len(req.Name) > 0 {
			if req.Name == plr.GetName() {
				return Err.Plr_IsSelf
			}

			fplr := app.PlayerMgr.FindPlayerByName(req.Name)
			if fplr == nil {
				return Err.Plr_NotFound
			}

			res.Plrs = append(res.Plrs, fplr.ToMsg_SimpleInfo())
		} else {
			i := 0
			app.PlayerMgr.ForEachOnlineIPlayerBreakable(func(fplr interface{}) bool {
				if i >= 4 {
					return false
				}

				rplr := fplr.(*app.Player)
				if plr.GetFriend().IsSearch(rplr.GetId()) {
					res.Plrs = append(res.Plrs, rplr.ToMsg_SimpleInfo())
					i++
				}
				return true
			})
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
