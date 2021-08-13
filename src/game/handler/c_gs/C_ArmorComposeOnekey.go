package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_ArmorComposeOnekey(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ArmorComposeOnekey)
	plr := ctx.(*app.Player)

	res := &msg.GS_ArmorComposeOnekey_R{}

	res.ErrorCode = func() int32 {
		bag := plr.GetBag()

		// collect costs and final armors
		op := bag.NewOp(gconst.ObjFrom_ArmorCompose)

		// accumulated ccy cost
		cost := make(map[int32]int32)

		// helper func
		f := func(srcid int32) (add int32, nextid int32) {
			id := srcid

			for {
				// composable ?
				conf := gamedata.ConfItem.Query(id)
				if conf == nil || len(conf.ArmorCompose) == 0 {
					nextid = 0
					break
				}

				// calc
				n1 := bag.GetItem(id) + add
				n2 := n1 / 3
				n1 = n2 * 3

				if n2 <= 0 {
					nextid = id + 1
					break
				}

				// cost: ccy
				for _, v := range conf.ArmorCompose {
					cost[v.Id] += v.N
				}

				// check ccy cost
				b := true
				for k, v := range cost {
					if plr.GetBag().GetCcy(k) < int64(v) {
						nextid = 0
						b = false
						break
					}
				}
				if !b {
					break
				}

				// yes, ccy is enough. put them in 'op'
				for _, v := range conf.ArmorCompose {
					op.Dec(v.Id, v.N)
				}

				// cost: id-item. !Note: MUST use op.Add()
				m := add - n1
				op.Add(id, m, 0)

				// next
				id++
				add = n2
			}

			// add
			if id > srcid {
				op.Inc(id, add)
			}

			// return nextid
			return
		}

		cnt := int32(0)
		// collect
		for id := req.SrcId; id > 0; {
			add, v := f(id)
			id = v

			if add > 0 {
				cnt += add
			}
		}

		// final check: make sure our calculation is correct.
		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		// apply
		res.Rewards = op.Apply().ToMsg()

		// fire
		evtmgr.Fire(gconst.Evt_ArmorCompose, plr, cnt)

		return Err.OK
	}()

	plr.SendMsg(res)
}
