package itemuse

import (
	"fw/src/core"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================
// 类型

const (
	use_method_drop         = "drop"
	use_method_counter      = "counter"
	use_method_appoint      = "appoint"
	use_method_hero_compose = "hero"
	use_method_vipLv        = "vipLv"
	use_method_bill         = "bill"
)

// ============================================================================
// 新的处理函数添加在下面

var use_handlers = map[string]use_func_t{

	// 掉落
	use_method_drop: func(plr *app.Player, id, n int32) (res *msg.GS_ItemUse_R) {
		res = &msg.GS_ItemUse_R{ErrorCode: Err.OK}

		conf := gamedata.ConfItem.Query(id)
		if conf == nil || len(conf.FuncParameter) < 2 {
			res.ErrorCode = Err.Failed
			return
		}

		if conf.UseLimit > 0 && conf.UseLimit < n {
			n = conf.UseLimit
		}

		if ec := plr.GetBag().CheckFull(); ec != Err.OK {
			res.ErrorCode = ec
			return
		}

		// modn
		cnt := n / conf.FuncParameter[0]
		if cnt == 0 {
			res.ErrorCode = Err.NotEnoughObject(id)
			return
		}
		n = cnt * conf.FuncParameter[0]

		op := plr.GetBag().NewOp(gconst.ObjFrom_ItemUse)
		op.Dec(id, n)

		if ec := op.CheckEnough(); ec != Err.OK {
			res.ErrorCode = ec
			return
		}

		for i := int32(0); i < cnt; i++ {
			for _, v := range utils.Drop(plr, conf.FuncParameter[1]) {
				op.Inc(v.Id, v.N)
			}
		}

		res.Rewards = op.Apply().ToMsg()

		return
	},

	// counter
	use_method_counter: func(plr *app.Player, id, n int32) (res *msg.GS_ItemUse_R) {
		res = &msg.GS_ItemUse_R{ErrorCode: Err.OK}

		conf := gamedata.ConfItem.Query(id)
		if conf == nil || len(conf.FuncParameter) < 2 {
			res.ErrorCode = Err.Failed
			return
		}

		if conf.UseLimit > 0 && conf.UseLimit < n {
			n = conf.UseLimit
		}

		conf_c := gamedata.ConfCounter.Query(conf.FuncParameter[0])
		if conf_c == nil {
			res.ErrorCode = Err.Failed
			return
		}

		cop := plr.GetCounter().NewOp(gconst.ObjFrom_ItemUse)

		cop.Dec(id, int64(n))
		if ec := cop.CheckEnough(); ec != Err.OK {
			res.ErrorCode = ec
			return
		}

		cop.IncCounter(conf.FuncParameter[0], int64(n*conf.FuncParameter[1]))

		cop.Apply()

		return
	},

	// appoint add task
	use_method_appoint: func(plr *app.Player, id, n int32) (res *msg.GS_ItemUse_R) {
		res = &msg.GS_ItemUse_R{ErrorCode: Err.OK}

		conf := gamedata.ConfItem.Query(id)
		if conf == nil || len(conf.FuncParameter) < 2 {
			res.ErrorCode = Err.Failed
			return
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_ItemUse)

		op.Dec(id, 1)
		if ec := op.CheckEnough(); ec != Err.OK {
			res.ErrorCode = ec
			return
		}

		addN := plr.GetAppoint().AddTask(1, conf.FuncParameter[0], conf.FuncParameter[1], true)
		if addN == 0 {
			res.ErrorCode = Err.Appoint_TaskCntLimit
			return
		}

		op.Apply()

		return
	},

	use_method_hero_compose: func(plr *app.Player, id, n int32) (res *msg.GS_ItemUse_R) {
		res = &msg.GS_ItemUse_R{ErrorCode: Err.OK}

		conf := gamedata.ConfItem.Query(id)
		if conf == nil || len(conf.FuncParameter) < 2 || conf.FuncParameter[0] <= 0 {
			res.ErrorCode = Err.Failed
			return
		}

		if conf.UseLimit > 0 && conf.UseLimit < n {
			n = conf.UseLimit
		}

		// modn
		cnt := n / conf.FuncParameter[0]
		if cnt == 0 {
			res.ErrorCode = Err.NotEnoughObject(id)
			return
		}
		n = cnt * conf.FuncParameter[0]

		op := plr.GetBag().NewOp(gconst.ObjFrom_ItemUse)
		op.Dec(id, n)

		if ec := op.CheckEnough(); ec != Err.OK {
			res.ErrorCode = ec
			return
		}

		op.Inc(conf.FuncParameter[1], cnt)

		res.Rewards = op.Apply().ToMsg()

		return
	},

	use_method_vipLv: func(plr *app.Player, id, n int32) (res *msg.GS_ItemUse_R) {
		res = &msg.GS_ItemUse_R{ErrorCode: Err.OK}

		conf := gamedata.ConfItem.Query(id)
		if conf == nil {
			res.ErrorCode = Err.Failed
			return
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_ItemUse)
		op.Dec(id, 1)
		if ec := op.CheckEnough(); ec != Err.OK {
			res.ErrorCode = ec
			return
		}

		if ec := plr.GetVip().LevelUp(); ec != Err.OK {
			res.ErrorCode = ec
			return
		}

		op.Apply()

		return
	},

	use_method_bill: func(plr *app.Player, id, n int32) (res *msg.GS_ItemUse_R) {
		res = &msg.GS_ItemUse_R{ErrorCode: Err.OK}

		conf := gamedata.ConfItem.Query(id)
		if conf == nil || len(conf.FuncParameter) < 3 {
			res.ErrorCode = Err.Failed
			return
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_ItemUse)
		op.Dec(id, 1)
		if ec := op.CheckEnough(); ec != Err.OK {
			res.ErrorCode = ec
			return
		}

		pid := int32(0)
		for _, prod := range gamedata.ConfBillProduct.Items() {
			if prod.PayId == conf.FuncParameter[0] && plr.GetSdk() == prod.Sdk {
				pid = prod.Id
				break
			}
		}

		conf_p := gamedata.ConfBillProduct.Query(pid)
		if conf_p == nil {
			res.ErrorCode = Err.Failed
			return
		}

		csext := ""
		tp := core.I32toa(conf.FuncParameter[1])
		if tp == gconst.Bill_CsExt_Type_GiftShop ||
			tp == gconst.Bill_CsExt_Type_PushGift ||
			tp == gconst.Bill_CsExt_Type_ActGift {
			csext = tp + "_" + core.I32toa(conf.FuncParameter[2])
		}

		op.Apply()

		plr.GetBill().GiveItems(pid, csext, conf_p.Price, conf_p.Ccy)

		return
	},

	// #todo:

}

// ============================================================================
// api
// ============================================================================

type use_func_t func(plr *app.Player, id, n int32) *msg.GS_ItemUse_R

// ============================================================================

func GetUseHandler(id int32) use_func_t {
	conf := gamedata.ConfItem.Query(id)
	if conf == nil {
		return nil
	}

	f, ok := use_handlers[conf.Func]
	if !ok {
		return nil
	}

	return f
}

// ============================================================================
