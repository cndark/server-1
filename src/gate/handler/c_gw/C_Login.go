package c_gw

import (
	"fw/src/gate/app"
	"fw/src/gate/app/dbmgr"
	"fw/src/gate/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
	"time"
)

func C_Login(message msg.Message, ctx interface{}) {
	// !Note: in net-thread

	req := message.(*msg.C_Login)
	sess := ctx.(*app.Session)

	ec := func() int32 {
		// 开始登录
		if !sess.LoginBegin() {
			return Err.Login_Failed
		}

		// 延后结束登录
		defer sess.LoginEnd()

		// authenticated ?
		if !sess.IsAuthenticated() {
			return Err.Login_Failed
		}

		// authid & sdk
		auth_id := sess.AuthReq.AuthId
		sdk := sess.AuthReq.Sdk

		// get merged info
		svr, closereg := dbmgr.Center_GetMergedSvrInfo(req.Svr0)
		if svr == "" {
			return Err.Login_Failed
		}

		// find real gs
		gs := config.Games[svr]
		if gs == nil {
			return Err.Login_Failed
		}

		// 查找玩家信息
		ec, user := dbmgr.Center_GetUserInfo(auth_id, sdk, req.Svr0, svr, closereg, func(newinfo *dbmgr.UserInfo) {
			newinfo.Model = sess.AuthReq.Model
			newinfo.DevId = sess.AuthReq.DevId
			newinfo.Os = sess.AuthReq.Os
			newinfo.OsVer = sess.AuthReq.OsVer
			newinfo.IP = sess.GetIP()

			dbmgr.Center_UpdateAcctAuthRet(auth_id, sdk, sess.AuthRet)
		})

		if ec != Err.OK {
			return ec
		}

		// 检查是否封号
		if user.BanTs.After(time.Now()) {
			return Err.Login_UserBanned
		}

		// ok. update acct last svr
		if req.ChgSvr {
			dbmgr.Center_UpdateAcctLastSvr(auth_id, sdk, req.Svr0)
		}

		// login player: 通知到 GS
		b := sess.LoginPlayer(gs.Id, &msg.GW_UserOnline{
			Sid:      sess.GetId(),
			UserId:   user.UserId,
			AuthId:   user.AuthId,
			Svr0:     user.Svr0,
			Sdk:      user.Sdk,
			Model:    sess.AuthReq.Model, // use newer one
			DevId:    sess.AuthReq.DevId, // use newer one
			Os:       sess.AuthReq.Os,    // use newer one
			OsVer:    sess.AuthReq.OsVer, // use newer one
			LoginIP:  sess.GetIP(),
			Language: req.Language,
			AuthRet:  sess.AuthRet,
		})
		if b {
			return Err.OK
		} else {
			return Err.Login_Failed
		}
	}()

	// send result
	sess.SendMsg(&msg.GW_Login_R{
		ErrorCode: ec,
	})
}
