package refund

import (
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
)

// ============================================================================

var RefundCode = make(map[string]*refund_code_t)

// ============================================================================

type refund_code_t struct {
	BillCnt int32 `bson:"bill"`  // 充值金额
	Taken   bool  `bson:"taken"` // 是否领取
}

// ============================================================================

func Open() {
	if config.CurGame.Id != 1 {
		return
	}

	var obj []*struct {
		Id    string `bson:"_id"`
		Bill  int32  `bson:"bill"`
		Taken bool   `bson:"taken"`
	}

	// load
	err := dbmgr.DBBill.GetAllObjects(
		"refund_take",
		&obj,
	)

	if err != nil {
		log.Warning("get refund code failed:", err)
	}

	for _, v := range obj {
		RefundCode[v.Id] = &refund_code_t{
			BillCnt: v.Bill,
			Taken:   v.Taken,
		}
	}
}

func TakeRefund(plr IPlayer, code string) (ec int32, rwd *msg.Rewards) {
	rc := RefundCode[code]
	if rc == nil {
		ec = Err.Bill_RefundCodeError
		return
	}

	if rc.Taken {
		ec = Err.Bill_RefundCodeTaken
		return
	}

	cnt := (rc.BillCnt / 100) * 30 // 人民币元*10*3倍

	op := plr.GetBag().NewOp(gconst.ObjFrom_RefundCode)
	op.Inc(gconst.Diamond, cnt)

	// flush
	err := dbmgr.DBBill.Update(
		"refund_take",
		code,
		db.M{"$set": db.M{"taken": true}},
	)

	if err != nil {
		log.Warning("flush refund_take failed:", code, err)
		ec = Err.Bill_RefundUpdateFailed
		return
	}

	// ok
	rc.Taken = true

	ec = Err.OK
	rwd = op.Apply().ToMsg()

	return
}
