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

func C_TokenGet(message msg.Message, ctx interface{}) {
	// !Note: in net-thread

	// req := message.(*msg.C_TokenGet)
	sess := ctx.(*app.Session)

	ec := func() int32 {
		// authenticated ?
		if !sess.IsAuthenticated() {
			return Err.Failed
		}

		// check
		if !sess.CheckTokenGet() {
			return Err.Failed
		}

		// http 请求
		webret := core.HttpPost(
			fmt.Sprintf(
				"http://%s:%d/token/get",
				config.Auth.IP,
				config.Auth.Port,
			),
			url.Values{
				"auth_id": {sess.AuthReq.AuthId},
				"sdk":     {sess.AuthReq.Sdk},
				"devid":   {sess.AuthReq.DevId},
			},
		)

		// checkout ret-values
		retobj := make(map[string]string)
		err := json.Unmarshal([]byte(webret), &retobj)
		if err != nil {
			return Err.Failed
		}

		ret_err := retobj["err"]
		ret_tk := retobj["token"]

		// check ret-err
		if ret_err != "" {
			return Err.Failed
		}

		// ok
		sess.SendMsg(&msg.GW_TokenGet_R{
			ErrorCode: Err.OK,
			Token:     ret_tk,
		})

		return Err.OK
	}()

	if ec != Err.OK {
		sess.SendMsg(&msg.GW_TokenGet_R{
			ErrorCode: ec,
		})
	}
}
