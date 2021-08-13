package c_gw

import (
	"encoding/json"
	"fmt"
	"fw/src/core"
	"fw/src/gate/app"
	"fw/src/gate/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
	"net/url"
)

func C_TokenAuth(message msg.Message, ctx interface{}) {
	// !Note: in net-thread

	req := message.(*msg.C_TokenAuth)
	sess := ctx.(*app.Session)

	ec := func() int32 {
		// 开始认证
		if !sess.AuthBegin() {
			return Err.Auth_Failed
		}

		// 延后结束认证
		defer sess.AuthEnd()

		// http 请求
		webret := core.HttpPost(
			fmt.Sprintf(
				"http://%s:%d/token/auth",
				config.Auth.IP,
				config.Auth.Port,
			),
			url.Values{
				"token":   {req.Token},
				"auth_id": {req.AuthId},
				"sdk":     {req.Sdk},
				"devid":   {req.DevId},
			},
		)

		// checkout ret-values
		retobj := make(map[string]string)
		err := json.Unmarshal([]byte(webret), &retobj)
		if err != nil {
			return Err.Auth_Failed
		}

		ret_err := retobj["err"]

		// check ret-err
		if ret_err != "" {
			return Err.Auth_Failed
		}

		// set auth-info
		sess.AuthReq = &msg.C_Auth{
			AuthId: req.AuthId,
			Sdk:    req.Sdk,
			DevId:  req.DevId,
		}

		// ok
		sess.SendMsg(&msg.GW_TokenAuth_R{
			ErrorCode: Err.OK,
		})

		return Err.OK
	}()

	if ec != Err.OK {
		sess.SendMsg(&msg.GW_TokenAuth_R{
			ErrorCode: ec,
		})
	}
}
