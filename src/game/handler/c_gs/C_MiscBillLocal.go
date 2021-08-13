package c_gs

import (
	"encoding/json"
	"fmt"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
)

type bill_local_t struct {
	Op      string `json:"op"`
	OrderId string `json:"orderid"`
	PId     int32  `json:"prod_id"`
	CsExt   string `json:"csext"`
	Ts      int64  `json:"ts"`
}

func C_MiscBillLocal(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_MiscBillLocal)
	plr := ctx.(*app.Player)

	res := &msg.GS_MiscBillLocal_R{}

	addr := fmt.Sprintf("http://%s:%d/local", config.Bill.IP, config.Bill.Port)

	go func() {
		res.ErrorCode = Err.OK
		res.R = core.HttpPostJson(addr, req.J)

		plr.SendMsg(res)

		// Unmarshal genbill
		var obj *bill_local_t
		err := json.Unmarshal([]byte(res.R), &obj)
		if err != nil {
			return
		}

		if obj != nil && obj.Op == "genorder" {
			evtmgr.Fire(gconst.Evt_BillGen, plr, obj.PId, obj.CsExt, obj.OrderId, obj.Ts)
		}
	}()
}
