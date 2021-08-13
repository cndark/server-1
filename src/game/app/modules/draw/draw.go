package draw

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math/rand"
	"time"
)

// ============================================================================

var rand_d = rand.New(rand.NewSource(time.Now().Unix()))

// ============================================================================
// 祈愿
type Draw struct {
	Tps   map[string]*draw_types_t
	Score int32

	plr IPlayer
}

type draw_types_t struct {
	Tp          string    // tp
	Cnt         int32     // 抽卡次数
	TotalCnt    int32     // 总抽卡次数
	LastTs      time.Time // 上次抽卡时间
	InitialTake int32     // 初始保底第几次生效
}

// ============================================================================
func NewDraw() *Draw {
	return &Draw{
		Tps: make(map[string]*draw_types_t),
	}
}

// ============================================================================
func (self *Draw) Init(plr IPlayer) {
	self.plr = plr

	for tp := range gamedata.ConfDraw.Items() {
		if self.Tps[tp] == nil {
			self.Tps[tp] = &draw_types_t{
				Tp: tp,
			}
		}
	}
}

func (self *Draw) GetDrawTp(tp string) *draw_types_t {
	return self.Tps[tp]
}

// 抽卡
func (self *Draw) DrawOne(tp string) (id int32) {
	drw := self.Tps[tp]
	if drw == nil {
		return
	}

	conf := gamedata.ConfDraw.Query(drw.Tp)
	if conf == nil {
		return
	}

	// 掉落组 && 保底是否清零
	drop_grp, isZeroExtra, isInitialGroup := drw.Tp, false, false
	if conf.InitialTimes != 0 && drw.TotalCnt < conf.InitialTimes { // 前N次保底
		isZeroExtra = true
		drop_grp = conf.InitialBaseGroup
		// 最后一次必出
		if drw.InitialTake == 0 {
			if (drw.TotalCnt == conf.InitialTimes-1) || (rand_d.Float32() < conf.InitialOdds) {
				drop_grp, isInitialGroup = conf.InitialGroup, true
			}
		}

	} else { // 循环保底

		// 固定组
		inFixed := false
		for _, v := range conf.FixedTimes {
			if drw.TotalCnt == v-1 {
				drop_grp, isZeroExtra = conf.FixedGroup, true
				inFixed = true
				break
			}
		}

		// 循环组
		if !inFixed {
			if conf.ExtraTimes != 0 && drw.Cnt >= conf.ExtraTimes-1 { // 循环保底次数
				drop_grp, isZeroExtra = conf.ExtraGroup, true
			}
		}
	}

	if drop_grp == "" {
		drop_grp = drw.Tp
	}

	id = self.pick_id(drw, drop_grp)
	if id == 0 {
		return
	}

	drw.TotalCnt++
	if isZeroExtra {
		drw.Cnt = 0
	} else {
		drw.Cnt++
	}

	if isInitialGroup {
		drw.InitialTake = drw.TotalCnt
	}

	self.Score += conf.Score

	return
}

func (self *Draw) pick_id(drw *draw_types_t, dgroup string) int32 {
	// select arr
	slt := make(map[int32]int32)

	lv := self.plr.GetLevel()
	for _, conf := range gamedata.ConfDrawOddsM.Items(dgroup) {
		if conf.IsDelivery != 0 {
			continue
		}

		odd := conf.BasicOdds

		// 附加权重
		// 等级
		for _, l := range conf.LevelOdds {
			if lv >= l.Lv {
				odd += l.Odd
				break
			}
		}

		// 充值等值钻石
		tb := self.plr.GetBillTotalBaseCcy()
		for _, b := range conf.PayOdds {
			if tb >= int64(b.Diam) {
				odd += b.Odd
				break
			}
		}

		// 动态召唤权重
		for _, v := range conf.DrawCntOdds {
			if drw.TotalCnt >= v.Cnt {
				odd += v.Odd
			}
		}

		// ##TODO
		//
		// ##TODO
		slt[conf.Id] = odd
	}

	return utils.PickWeightedMapId(slt)
}

func (self *Draw) Tutorial(tp string) *msg.GS_DrawTp_R {
	res := &msg.GS_DrawTp_R{}
	cnt := int32(1)

	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil {
		return nil
	}

	conf_d := gamedata.ConfDraw.Query(tp)
	if conf_d == nil {
		return nil
	}

	d := self.GetDrawTp(tp)
	if d == nil {
		return nil
	}

	//cost
	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_Draw)

	costid := int32(gconst.DrawHigh)
	itemid := conf_g.TutorialDrawSenior

	if tp == "normal" {
		costid = 0
		itemid = conf_g.TutorialDrawNormal
		cop.DecCounter(gconst.Cnt_DrawNormal_Free, 1)
	}

	cur := self.plr.GetBag().GetItem(costid)
	if cur < cnt {
		cop.Dec(costid, cur)
	} else {
		cop.Dec(costid, cnt)
	}

	cop.Inc(itemid, 1)
	res.Items = append(res.Items, &msg.Item{Id: itemid, Num: 1})

	d.TotalCnt += cnt
	self.Score += conf_d.Score * cnt

	rwds := cop.Apply().ToMsg()
	res.Score = self.Score
	res.Rewards = rwds
	res.DrawTp = d.ToMsg()

	evtmgr.Fire(gconst.Evt_Draw, self.plr, tp, conf_d.ModuleId, cnt, res.Items)

	return res
}

// 领取积分宝箱
func (self *Draw) ScoreBoxTake() (int32, *msg.Rewards) {
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil || len(conf.DrawScore) == 0 {
		return Err.Failed, nil
	}

	if self.plr.GetVipLevel() < conf.DrawScore[0].Vip {
		return Err.Plr_LowVipLevel, nil
	}

	if self.Score < conf.DrawScore[0].N {
		return Err.Draw_ScoreBoxScoreNotEnough, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_DrawBox)
	for _, v := range utils.Drop(self.plr, conf.DrawScore[0].Id) {
		op.Inc(v.Id, int64(v.N))
	}

	self.Score -= conf.DrawScore[0].N

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

// ============================================================================

func (self *Draw) ToMsg() *msg.DrawData {
	ret := &msg.DrawData{}
	for _, v := range self.Tps {
		ret.Tps = append(ret.Tps, &msg.DrawTpData{
			Tp:          v.Tp,
			Cnt:         v.Cnt,
			TotalCnt:    v.TotalCnt,
			InitialTake: v.InitialTake,
			LastTs:      v.LastTs.Unix(),
		})
	}

	ret.Score = self.Score

	return ret
}

func (self *draw_types_t) ToMsg() *msg.DrawTpData {
	return &msg.DrawTpData{
		Tp:          self.Tp,
		Cnt:         self.Cnt,
		TotalCnt:    self.TotalCnt,
		InitialTake: self.InitialTake,
		LastTs:      self.LastTs.Unix(),
	}
}
