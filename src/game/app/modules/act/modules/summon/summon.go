package summon

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

const (
	C_ActSummon_Pick_Select = 1 // 自选
	C_ActSummon_Pick_Fixed  = 2 // 固定
)

const (
	C_DesirePos = 1 // 心愿英雄位置
	C_UPPos1    = 2
	C_UPPos2    = 3
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
	NormalCnt int32           // 常规计数
	DesireCnt int32           // 心愿计数
	UpCnt     int32           // Up计数
	DiamCnt   int32           // 钻石召唤次数
	LastTs    time.Time       // 上次免费召唤时间
	HeroPos   map[int32]int32 // 英雄的站位【1：心愿，2、3：UP】
	OddCnt    map[int32]int32 // 心愿英雄随机到的次数【限定随机】
	ActGift   map[int32]int32 // 活动礼包购买次数
}

// ============================================================================

func init() {
	act.RegisterAct(gconst.ActName_Summon, actObj)

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
		ActGift: make(map[int32]int32),
		HeroPos: make(map[int32]int32),
		OddCnt:  make(map[int32]int32),
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

func (self *act_t) ToMsg(plr IPlayer) *msg.ActSummonData {
	plr_data := self.GetPlrData(plr)
	if plr_data == nil {
		return nil
	}

	ret := &msg.ActSummonData{
		NormalCnt: plr_data.NormalCnt,
		DesireCnt: plr_data.DesireCnt,
		DiamCnt:   plr_data.DiamCnt,
		UpCnt:     plr_data.UpCnt,
		LastTs:    plr_data.LastTs.Unix(),
		HeroPos:   plr_data.HeroPos,
		ActGift:   plr_data.ActGift,
	}

	return ret
}

// ============================================================================

func ActSummonInfo(plr IPlayer) *msg.ActSummonData {
	return actObj.ToMsg(plr)
}

// 选1个心愿英雄和2个UP英雄
func Pick(plr IPlayer, heroPos map[int32]int32) int32 {
	if !actObj.Started() {
		return Err.Act_ActClosed
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		return Err.Act_ActPlrDataNotFound
	}

	conf := gamedata.ConfActSummon.Query(actObj.GetConfGrp())
	if conf == nil {
		return Err.Failed
	}

	for idx, v := range heroPos {

		if idx == C_DesirePos && conf.Type == C_ActSummon_Pick_Fixed {
			if len(conf.DesireHero) > 0 {
				plr_data.HeroPos[C_DesirePos] = conf.DesireHero[0].Hero
				return Err.OK
			} else {
				return Err.Failed
			}
		}

		plr_data.HeroPos[idx] = v
	}

	return Err.OK
}

// 召唤
func Draw(plr IPlayer, isDiam bool, cnt int32) {
	res := &msg.GS_ActSummonDraw_R{}
	res.ErrorCode = func() int32 {
		if !actObj.Started() {
			return Err.Act_ActClosed
		}

		plr_data := actObj.GetPlrData(plr)
		if plr_data == nil {
			return Err.Act_ActPlrDataNotFound
		}

		conf := gamedata.ConfActSummon.Query(actObj.GetConfGrp())
		if conf == nil {
			return Err.Failed
		}

		if len(plr_data.HeroPos) < 3 {
			return Err.Activity_SummonHeroNotPick
		}

		conf_m := gamedata.ConfMonster.Query(plr_data.HeroPos[C_DesirePos])
		if conf_m == nil {
			return Err.Failed
		}

		if cnt != 10 {
			cnt = 1
		}

		// check cost
		op := plr.GetBag().NewOp(gconst.ObjFrom_ActSummon)
		now := time.Now()
		if (cnt == 1) && now.After(plr_data.LastTs.Add(time.Duration(conf.CostFreeTime)*time.Second)) {
			plr_data.LastTs = now
		} else {
			// 消耗钻石
			if isDiam {

				if plr_data.DiamCnt+cnt > conf.DiamCnt {
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

				plr_data.DiamCnt += cnt
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

		// draw
		for i := int32(0); i < cnt; i++ {

			item := &msg.Item{}

			zeroNormal, zeroDesire, zeroUp := false, false, false

			if conf_m.ActSummonCnt != 0 && plr_data.DesireCnt+1 >= conf_m.ActSummonCnt {
				// 心愿英雄
				op.Inc(plr_data.HeroPos[C_DesirePos], 1)

				item = &msg.Item{
					Id:  plr_data.HeroPos[C_DesirePos],
					Num: 1,
				}

				zeroNormal = true
				zeroDesire = true
				zeroUp = true
			} else {
				if conf.UpTimes != 0 && plr_data.UpCnt+1 >= conf.UpTimes {

					// UP随机
					upId := pick_up_desire(plr_data)

					if upId == plr_data.HeroPos[C_DesirePos] {
						zeroDesire = true
					}

					op.Inc(upId, 1)

					item = &msg.Item{
						Id:  upId,
						Num: 1,
					}

					zeroUp = true
					zeroNormal = true

				} else {

					// 保底组
					drop_grp := conf.SummonGrp
					if plr_data.NormalCnt+1 == conf.ExtraTimes {
						drop_grp = conf.ExtraGroup
						zeroNormal = true
					}

					pick_dropGrp(plr_data, drop_grp, func(zD, zN, zU bool, Id int32, N int32) {
						zeroDesire = zD
						zeroNormal = zN
						zeroUp = zU

						op.Inc(Id, N)

						item = &msg.Item{
							Id:  Id,
							Num: 1,
						}

					})

				}

			}

			if zeroUp {
				plr_data.UpCnt = 0
			} else {
				plr_data.UpCnt++
			}

			if zeroNormal {
				plr_data.NormalCnt = 0
			} else {
				plr_data.NormalCnt++
			}

			if zeroDesire {
				plr_data.DesireCnt = 0
			} else {
				plr_data.DesireCnt++
			}

			res.Items = append(res.Items, item)
		}

		res.LastTs = plr_data.LastTs.Unix()
		res.NormalCnt = plr_data.NormalCnt
		res.DesireCnt = plr_data.DesireCnt
		res.UpCnt = plr_data.UpCnt
		res.DiamCnt = plr_data.DiamCnt
		res.Rewards = op.Apply().ToMsg()

		return Err.OK
	}()

	plr.SendMsg(res)
}

func pick_dropGrp(plr_data *data_plr_t, drop_grp string, cbf func(zD bool, zN bool, zU bool, Id int32, N int32)) {

	zeroNormal, zeroDesire, zeroUp := false, false, false

	conf := gamedata.ConfActSummon.Query(actObj.GetConfGrp())

	seq := pick_id(plr_data, drop_grp)

	conf_o := gamedata.ConfActSummonOdds.Query(seq)
	if conf_o == nil {
		cbf(zeroDesire, zeroNormal, zeroUp, 0, 0)
	}

	for _, v := range conf_o.Item {

		for _, v1 := range conf.DesireHero {

			if v1.Hero == v.Id && plr_data.HeroPos[C_DesirePos] == v.Id {

				if plr_data.OddCnt[v.Id] == 0 {
					plr_data.OddCnt[v.Id] = conf_o.GetCntLimit + conf_o.AddCntLimit
				} else {
					plr_data.OddCnt[v.Id] = plr_data.OddCnt[v.Id] + conf_o.AddCntLimit
				}

				zeroNormal = true
				zeroDesire = true
				zeroUp = true

				break
			}
		}

		if v.Id == plr_data.HeroPos[C_UPPos1] || v.Id == plr_data.HeroPos[C_UPPos2] {
			zeroNormal = true
			zeroUp = true
		}

		cbf(zeroDesire, zeroNormal, zeroUp, v.Id, v.N)
	}

}

func pick_up_desire(plr_data *data_plr_t) (id int32) {
	slt := make(map[int32]int32)

	conf := gamedata.ConfActSummon.Query(actObj.GetConfGrp())

	// 三个英雄随机

	for _, v := range conf.DesireHero {
		if v.Hero == plr_data.HeroPos[C_DesirePos] {
			slt[v.Hero] = v.Weight
			break
		}
	}

	for _, v := range conf.UpHero {
		if v.Hero == plr_data.HeroPos[C_UPPos1] || v.Hero == plr_data.HeroPos[C_UPPos2] {
			slt[v.Hero] = v.Weight
		}
	}

	return utils.PickWeightedMapId(slt)
}

func pick_id(plr_data *data_plr_t, key string) int32 {

	conf := gamedata.ConfActSummon.Query(actObj.GetConfGrp())
	if conf == nil {
		return Err.Failed
	}

	slt := make(map[int32]int32)
	for _, v := range gamedata.ConfActSumonOddsM.Items(key) {

		// 如果是心愿或UP，且条件达到，添加到随机
		flag := false

		for _, v1 := range v.Item {

			if v1.Id == plr_data.HeroPos[C_DesirePos] {

				nflag := false
				for _, v2 := range conf.DesireHero {

					if (v1.Id == v2.Hero) && plr_data.OddCnt[v2.Hero] >= plr_data.DesireCnt {
						nflag = true
						break
					}
				}

				if nflag {
					flag = true
					break
				}
			}

			if (v1.Id == plr_data.HeroPos[C_UPPos1]) || (v1.Id == plr_data.HeroPos[C_UPPos1]) && (conf.UpTimes < plr_data.DesireCnt) {
				flag = true
				break
			}
		}

		if flag {
			continue
		}

		slt[v.Seq] += v.Weight
	}

	return utils.PickWeightedMapId(slt)
}
