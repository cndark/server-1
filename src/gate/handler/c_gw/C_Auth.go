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

func C_Auth(message msg.Message, ctx interface{}) {
	// !Note: in net-thread

	req := message.(*msg.C_Auth)
	sess := ctx.(*app.Session)

	ec := func() int32 {
		// 开始认证
		if !sess.AuthBegin() {
			return Err.Auth_Failed
		}

		// 延后结束认证
		defer sess.AuthEnd()

		// 检查版本号
		if req.VerMajor != config.Common.Ver_Major ||
			req.VerMinor != config.Common.Ver_Minor {
			return Err.Auth_InvalidVersion
		}

		// 参数检查
		if req.AuthId == "" {
			return Err.Auth_Failed
		}

		// http 请求
		webret := core.HttpPost(
			fmt.Sprintf(
				"http://%s:%d/sdk/auth",
				config.Auth.IP,
				config.Auth.Port,
			),
			url.Values{
				"auth_id":    {req.AuthId},
				"auth_token": {req.AuthToken},
				"sdk":        {req.Sdk},
				"devid":      {req.DevId},
				"ip":         {sess.GetIP()},
				"j":          {req.J},
			},
		)

		// checkout ret-values
		var retobj struct {
			Err    string
			AuthId string
			S      map[string]string
			C      map[string]string
		}

		err := json.Unmarshal([]byte(webret), &retobj)
		if err != nil {
			return Err.Auth_Failed
		}

		// check ret-err
		if retobj.Err != "" {
			return Err.Auth_Failed
		}

		// set ret-authid
		if retobj.AuthId != "" {
			req.AuthId = retobj.AuthId
		}

		// set auth-info
		sess.AuthReq = req
		sess.AuthRet = retobj.S

		// ok
		sess.SendMsg(&msg.GW_Auth_R{
			ErrorCode: Err.OK,
			AuthId:    req.AuthId,
			C:         retobj.C,
		})

		return Err.OK
	}()

	if ec != Err.OK {
		sess.SendMsg(&msg.GW_Auth_R{
			ErrorCode: ec,
		})
	}
}
