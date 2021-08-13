package msummon

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

// ============================================================================

var actObj = &act_t{}

// ============================================================================

type act_t struct {
	act.BaseAct
}

type data_svr_t struct {
}

type data_plr_t struct {
	NormalCnt int32 // 常规保底计数
	MagicCnt  int32 // 魔法保底计数
	DimCnt    int32 // 钻石消耗次数

	LastTs  time.Time       // 上次免费召唤时间
	ActGift map[int32]int32 // 活动礼包购买次数

}

// ============================================================================

func init() {
	act.RegisterAct(gconst.ActName_MagicSummon, actObj)

	// 更新购买活动礼包次数
	evtmgr.On(gconst.Evt_ActGift, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		actName := args[1].(string)
		id := args[2].(int32)

		if actObj.GetName() != actName {
			return
		}

		plr_data := actObj.GetPlrData(plr)
		if plr_data == nil {
			return
		}

		plr_data.ActGift[id]++
	})
}

// ============================================================================

func (self *act_t) NewSvrData() interface{} {
	return new(data_svr_t)
}

func (self *act_t) NewPlrData() interface{} {
	return &data_plr_t{
		ActGift: map[int32]int32{},
	}
}

func (self *act_t) ResetSvrData() {
}

func (self *act_t) ResetPlrData(iplr interface{}) {
}

func (self *act_t) GetSvrData() *data_svr_t {
	return self.GetActRawData().(*data_svr_t)
}

func (self *act_t) GetPlrData(plr IPlayer) *data_plr_t {
	return plr.GetActRawData(self.GetName()).(*data_plr_t)
}

// ============================================================================

func (self *act_t) Started() bool {
	return self.GetStage() == "start"
}

func (self *act_t) Ended() bool {
	return self.GetStage() == "end"
}

func (self *act_t) Closed() bool {
	return self.GetStage() == "close"
}

// ============================================================================

func (self *act_t) OnInit() {
}

func (self *act_t) OnQuit() {
}

func (self *act_t) OnStage() {
}

func (self *act_t) ToMsg(plr IPlayer) *msg.ActMSummonData {
	plr_data := self.GetPlrData(plr)
	if plr_data == nil {
		return nil
	}

	ret := &msg.ActMSummonData{
		MagicCnt:  plr_data.MagicCnt,
		NormalCnt: plr_data.NormalCnt,
		DiamCnt:   plr_data.DimCnt,
		LastTs:    plr_data.LastTs.Unix(),
		ActGift:   plr_data.ActGift,
	}

	return ret
}

// ============================================================================

func Draw(plr IPlayer, dim bool, cnt int32) {

	res := &msg.GS_ActMSummonDraw_R{}
	res.ErrorCode = func() int32 {
		if !actObj.Started() {
			return Err.Act_ActClosed
		}

		plr_data := actObj.GetPlrData(plr)
		if plr_data == nil {
			return Err.Act_ActPlrDataNotFound
		}

		conf := gamedata.ConfActMagicSummon.Query(actObj.GetConfGrp())
		if conf == nil {
			return Err.Act_ConfGrp
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_ActMagicSummon)

		if cnt != 10 {
			cnt = 1
		}

		now := time.Now()
		if (cnt == 1) && now.After(plr_data.LastTs.Add(time.Duration(conf.CostFreeTime)*time.Second)) {
			plr_data.LastTs = now
		} else {
			if dim {
				if plr_data.DimCnt+cnt > conf.DiamCnt {
					return Err.Activity_SummonDiamLimit
				}

				if cnt == 10 {
					op.Dec(gconst.Diamond, int32(float32(conf.DiamCost*cnt)*conf.Discount))
				} else {
					op.Dec(gconst.Diamond, conf.DiamCost*cnt)
				}

				if ec := op.CheckEnough(); ec != Err.OK {
					return ec
				}

				plr_data.DimCnt += cnt
			} else {
				for _, v := range conf.Cost {
					if cnt == 10 {
						op.Dec(v.Id, int32(float32(v.N*cnt)*conf.Discount))
					} else {
						op.Dec(v.Id, v.N*cnt)
					}
				}

				if ec := op.CheckEnough(); ec != Err.OK {
					return ec
				}
			}
		}

		for i := 0; i < int(cnt); i++ {

			zeroNormal, zeroMagic := false, false
			if plr_data.MagicCnt+1 >= conf.MagicExtraTimes {
				hid := actObj.pick_hero(plr)
				if hid == 0 {
					continue
				}

				op.Inc(hid, 1)

				zeroNormal = true
				zeroMagic = true
			} else {
				drop_grp := conf.SummonGrp
				if plr_data.NormalCnt+1 >= conf.ExtraTimes {
					drop_grp = conf.ExtraGroup
					zeroNormal = true
				}

				seq := pick_id(drop_grp, plr_data.MagicCnt)
				conf_o := gamedata.ConfActMagicSummonOdds.Query(seq)
				if conf_o == nil {
					continue
				}

				op.Inc(conf_o.RewardId, conf_o.RewardNum)

				for _, v := range conf.MagicHero {
					if v.Hero == conf_o.RewardId {
						zeroMagic = true
						break
					}
				}
			}

			if zeroNormal {
				plr_data.NormalCnt = 0
			} else {
				plr_data.NormalCnt++
			}

			if zeroMagic {
				plr_data.MagicCnt = 0
			} else {
				plr_data.MagicCnt++
			}
		}

		res.DiamCnt = plr_data.DimCnt
		res.LastTs = plr_data.LastTs.Unix()
		res.MagicCnt = plr_data.MagicCnt
		res.NormalCnt = plr_data.NormalCnt
		res.Rewards = op.Apply().ToMsg()

		return Err.OK
	}()

	plr.SendMsg(res)
}

func (self *act_t) pick_hero(plr IPlayer) int32 {
	conf := gamedata.ConfActMagicSummon.Query(self.GetConfGrp())
	if conf == nil {
		return 0
	}

	slt := map[int32]int32{}
	for _, v := range conf.MagicHero {
		slt[v.Hero] = v.Weight
	}

	return utils.PickWeightedMapId(slt)
}

func pick_id(group string, mcnt int32) int32 {
	slt := map[int32]int32{}
	for _, v := range gamedata.ConfActMagicSummonOddsM.Items(group) {
		if v.MagicTimesCond > 0 && mcnt < v.MagicTimesCond {
			continue
		}

		slt[v.Seq] = v.Weight
	}

	return utils.PickWeightedMapId(slt)
}

// ============================================================================

func ActMSummonInfo(plr IPlayer) *msg.ActMSummonData {
	return actObj.ToMsg(plr)
}
