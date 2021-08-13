package gw_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
)

func GW_UserOnline(message msg.Message, ctx interface{}) {
	req := message.(*msg.GW_UserOnline)

	// * try load player
	// * if not found, create new
	plr, ec := app.PlayerMgr.LoadPlayerWithError(req.UserId)

	if ec == Err.Plr_NotFound {
		// create new user
		plr = app.PlayerMgr.CreatePlayer(req.UserId, func(user *app.User) {
			user.AuthId = req.AuthId
			user.Svr0 = req.Svr0
			user.Svr = config.CurGame.Name
			user.Sdk = req.Sdk

			user.Name = utils.GenRandName(req.Language)
		})
		if plr == nil {
			app.NetMgr.Send2Player(req.Sid, &msg.GS_LoginError{ErrorCode: Err.Login_GetUserInfo})
			return
		}
	} else if ec != Err.OK {
		// critical error
		app.NetMgr.Send2Player(req.Sid, &msg.GS_LoginError{ErrorCode: Err.Login_GetUserInfo})
		return
	} else {
		// ok. user exists. check double login
		if plr.IsOnline() {
			if plr.Sid() == req.Sid {
				// same session: just offline
				app.PlayerMgr.SetOffline(plr, false)
			} else {
				// different session: logout previous one
				plr.Logout()
			}
		}
	}

	// update info
	plr.User().Model = req.Model
	plr.User().DevId = req.DevId
	plr.User().Os = req.Os
	plr.User().OsVer = req.OsVer
	plr.User().LoginIP = req.LoginIP
	plr.User().AuthRet = req.AuthRet

	// set online
	app.PlayerMgr.SetOnline(plr, req.Sid)
}
