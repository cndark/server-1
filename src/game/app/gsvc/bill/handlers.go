package bill

import (
	"fw/src/core/db"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/game/app"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	"fw/src/shared/config"
	"net/http"
)

// ============================================================================

var (
	done_orders = make(map[string]bool)
)

// ============================================================================

func handle_give_items(req *http.Request) (r string, err error) {
	order_key := get_string(req, "order_key")
	userid := get_string(req, "userid")
	prod_id := get_int32(req, "prod_id")
	csext := get_string(req, "csext")
	amount := get_int32(req, "amount")
	ccy := get_string(req, "ccy")
	orderid := get_string(req, "orderid")
	cp_orderid := get_string(req, "cp_orderid")

	// check done_orders
	if done_orders[order_key] {
		r = "ok"
		return
	}

	// load player
	plr := app.PlayerMgr.LoadPlayer(userid)
	if plr == nil {
		r = "e_user"
		return
	}

	// give items
	if !plr.GetBill().GiveItems(prod_id, csext, amount, ccy) {
		r = "e_cfg"
		return
	}

	// notify
	if orderid != "" {
		plr.SendMsg(&msg.GS_BillOrder{
			OrderId: orderid,
			Amount:  amount,
		})
	}

	// isfirst
	isFirst := false
	if !plr.GetBill().IsRealFirst {
		isFirst = true
		plr.GetBill().IsRealFirst = true
	}

	// stats
	evtmgr.Fire(gconst.Evt_BillStats, plr, amount, ccy, cp_orderid, prod_id, csext)

	// update done-orders
	done_orders[order_key] = true

	// svr & lv
	svr := config.CurGame.Name
	lv := plr.GetLevel()

	// update order
	async.PushQ(
		gconst.AQ_Bill,
		func() {
			err := dbmgr.DBBill.Update(
				"order",
				order_key,
				db.M{
					"$set": db.M{
						"status": "ok",
						"svr":    svr,
						"lv":     lv,
						"first":  isFirst,
					},
				},
			)
			if err != nil {
				log.Warning("update order status failed:", order_key, err)
			}
		},
	)

	r = "ok"
	return
}
