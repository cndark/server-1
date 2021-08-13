package actgift

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/guild"
	"fw/src/game/msg"
	"strings"
)

// 所有活动礼包在这里监听
func init() {
	on_bill_done()
}

func on_bill_done() {
	evtmgr.On(gconst.Evt_BillDone, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		prod_id := args[1].(int32)
		csext := args[2].(string)

		prod := gamedata.ConfBillProduct.Query(prod_id)
		if prod == nil || prod.TypeId != gconst.Bill_Gift {
			return
		}

		arr := strings.Split(csext, "_")
		if len(arr) < 2 || arr[0] != gconst.Bill_CsExt_Type_ActGift {
			return
		}

		id := core.Atoi32(arr[1])

		conf := gamedata.ConfActGift.Query(id)
		if conf == nil || conf.PayId != prod.PayId {
			return
		}

		actObj := act.FindAct(conf.ActName)
		if actObj == nil || actObj.GetStage() != "start" {
			return
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_ActGift)
		for _, v := range conf.Reward {
			op.Inc(v.Id, v.N)
		}

		rwds := op.Apply().ToMsg()

		//res
		plr.SendMsg(&msg.GS_ActGiftNew{
			Id:      id,
			Rewards: rwds,
		})

		//evt
		evtmgr.Fire(gconst.Evt_ActGift, plr, conf.ActName, id)

		// guild gift
		if conf.Type == gconst.C_ActGift_Gld {
			gld := guild.GuildMgr.FindGuild(plr.GetGuildId())
			if gld == nil {
				return
			}

			for plrid := range gld.Members {
				mplr := load_player(plrid)
				if mplr == nil {
					continue
				}

				mplr.GetMisc().GldActGiftAdd(id, plr.GetName())
			}
		}
	})
}
