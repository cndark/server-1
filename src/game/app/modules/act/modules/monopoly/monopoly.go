package monopoly

import (
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math/rand"
	"time"
)

// ============================================================================

var rand_d = rand.New(rand.NewSource(time.Now().Unix()))

// ============================================================================

const (
	C_Monopoly_LevelMax = 60
)

// 奇遇类型

const (
	C_Monopoly_AdvBox       = 1 // 宝箱奇遇
	C_Monopoly_AdvProblem   = 2 // 答题奇遇
	C_Monopoly_AdvShop      = 3 // 商店折扣奇遇
	C_Monopoly_AdvNormalBat = 4 // 普通战斗奇遇
	C_Monopoly_AdvHeroBat   = 5 // 英雄战斗奇遇
)

// 战斗类型

const (
	C_Battle_Normal = 1
	C_Battle_Hero   = 2
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
	Attain    map[int32]*attain_obj_t // 统计进度
	TaskTaken []int32                 // 任务领取
	CurLv     int32                   // 已通关
	Pos       int32                   // 点位
	RepeatCnt int32                   // 最大关卡后重复次数

	PosInfo map[int32]int32 // 本关奇遇点位的信息

	NormalBat item_l    // 普通战斗
	Problems  item_l    // 答题
	HeroBat   []*bat_t  // 英雄挑战
	Shops     []*shop_t // 折扣商店
}

type shop_t struct {
	Id         int32     // 配置id
	ExpireTs   time.Time // 到期时间
	AlreadyBuy []int32   // 已经购买的商品下标
	DisCnt     int32     // 商店折扣
}

type bat_t struct {
	Id        int32     // 配置id
	ExpireTs  time.Time // 到期时间
	HeroLv    int32
	MonsterLv int32
}

type item_l []*item_t
type item_t struct {
	Id       int32     // 配置id
	ExpireTs time.Time // 到期时间
}

type attain_obj_t struct {
	Id  int32
	Val float64 // progress value
}

// ============================================================================

func init() {
	act.RegisterAct(gconst.ActName_Monopoly, actObj)
}

// ============================================================================

func (self *act_t) NewSvrData() interface{} {
	return new(data_svr_t)
}

func (self *act_t) NewPlrData() interface{} {
	return &data_plr_t{
		Attain:  map[int32]*attain_obj_t{},
		PosInfo: map[int32]int32{},
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

func (self *act_t) ToMsg(plr IPlayer) *msg.ActMonopolyData {
	ret := &msg.ActMonopolyData{}

	plr_data := self.GetPlrData(plr)
	if plr_data == nil {
		return nil
	}

	for _, v := range plr_data.Attain {
		ret.Attain = append(ret.Attain, &msg.ActAttainObj{
			OId: v.Id,
			Val: v.Val,
		})
	}

	for i := 0; i < len(plr_data.NormalBat); i++ {
		v := plr_data.NormalBat[i]
		if v != nil {
			if v.ExpireTs.After(time.Now()) {
				ret.NormalBat = append(ret.NormalBat, &msg.ActMonopolyAdvItem{
					Id:       v.Id,
					ExpireTs: v.ExpireTs.Unix(),
				})
			} else {
				plr_data.NormalBat = append(plr_data.NormalBat[0:i], plr_data.NormalBat[i+1:]...)
				i--
			}
		}
	}

	for i := 0; i < len(plr_data.HeroBat); i++ {
		v := plr_data.HeroBat[i]
		if v != nil {
			if v.ExpireTs.After(time.Now()) {
				ret.HeroBat = append(ret.HeroBat, &msg.ActMonopolyBatItem{
					Id:        v.Id,
					ExpireTs:  v.ExpireTs.Unix(),
					HeroLv:    v.HeroLv,
					MonsterLv: v.MonsterLv,
				})
			} else {
				plr_data.HeroBat = append(plr_data.HeroBat[0:i], plr_data.HeroBat[i+1:]...)
				i--
			}
		}
	}

	for i := 0; i < len(plr_data.Problems); i++ {
		v := plr_data.Problems[i]
		if v != nil {
			if v.ExpireTs.After(time.Now()) {
				ret.Problems = append(ret.Problems, &msg.ActMonopolyAdvItem{
					Id:       v.Id,
					ExpireTs: v.ExpireTs.Unix(),
				})
			} else {
				plr_data.Problems = append(plr_data.Problems[0:i], plr_data.Problems[i+1:]...)
				i--
			}
		}
	}

	for i := 0; i < len(plr_data.Shops); i++ {
		v := plr_data.Shops[i]
		if v != nil {
			if v.ExpireTs.After(time.Now()) {
				ret.Shops = append(ret.Shops, &msg.ActMonolyShopItem{
					Id:         v.Id,
					ExpireTs:   v.ExpireTs.Unix(),
					AlreadyBuy: v.AlreadyBuy,
					DisCnt:     v.DisCnt,
				})
			} else {
				plr_data.Shops = append(plr_data.Shops[0:i], plr_data.Shops[i+1:]...)
				i--
			}
		}
	}

	if len(plr_data.PosInfo) == 0 {
		plr_data.PosInfo = map[int32]int32{}
		rand_pos(plr)
	}

	ret.AdvPosInfo = plr_data.PosInfo
	ret.TaskTaken = plr_data.TaskTaken
	ret.RepeatCnt = plr_data.RepeatCnt
	ret.Pos = plr_data.Pos

	ret.Seq = lv_to_seq(plr_data.CurLv)

	return ret
}

// ============================================================================

func (self *data_plr_t) get_attain_obj(oid int32) *attain_obj_t {
	obj := self.Attain[oid]
	if obj == nil {
		self.Attain[oid] = &attain_obj_t{Id: oid}
		return self.Attain[oid]
	}
	return obj
}

// ============================================================================
// implements ICondObj interface

func (self *attain_obj_t) GetVal() float64 {
	return self.Val
}

func (self *attain_obj_t) SetVal(v float64) {
	self.Val = v
}

func (self *attain_obj_t) AddVal(v float64) {
	self.Val += v
}

func (self *attain_obj_t) Done(body interface{}, confid int32, isChange bool) {
	if !actObj.Started() || !isChange {
		return
	}

	body.(IPlayer).SendMsg(&msg.GS_ActMonopolyObjValueChanged{
		OId: self.Id,
		Val: self.Val,
	})
}

// ============================================================================
// 活动领奖
func TaskTake(plr IPlayer, id int32) (int32, *msg.Rewards) {
	if !actObj.Started() {
		return Err.Act_ActClosed, nil
	}

	conf := gamedata.ConfActMonopolyTask.Query(id)
	if conf == nil {
		return Err.Failed, nil
	}

	if conf.ConfGrp != actObj.GetConfGrp() {
		return Err.Act_ConfGrp, nil
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		return Err.Act_ActPlrDataNotFound, nil
	}

	for _, v := range plr_data.TaskTaken {
		if v == id {
			return Err.Plr_TakenBefore, nil
		}
	}

	op := plr.GetBag().NewOp(gconst.ObjFrom_ActMonopoly)

	for _, v := range conf.Attain {
		if plr_data.Attain[v.AttainId].Val < v.P2 {
			return Err.Activity_CondLimited, nil
		} else {
			for _, v := range conf.Reward {
				op.Inc(v.Id, v.N)
			}
			plr_data.TaskTaken = append(plr_data.TaskTaken, id)
		}
	}

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

// 色子前进N步
func NextNPos(plr IPlayer, step int32) {

	res := &msg.GS_ActMonopolyNPos_R{}

	res.ErrorCode, res.Rewards, res.Step = func() (int32, *msg.Rewards, int32) {
		if !actObj.Started() {
			return Err.Act_ActClosed, nil, step
		}

		plr_data := actObj.GetPlrData(plr)
		if plr_data == nil {
			return Err.Act_ActPlrDataNotFound, nil, step
		}

		conf_lv := gamedata.ConfActMonopolyLevel.Query(lv_to_seq(plr_data.CurLv))
		if conf_lv == nil {
			return Err.Failed, nil, step
		}

		if conf_lv.ConfGrp != actObj.GetConfGrp() {
			return Err.Act_ConfGrp, nil, step
		}

		conf_g := gamedata.ConfActivityPublic.Query(1)
		if conf_g == nil {
			return Err.Failed, nil, step
		}

		op := plr.GetBag().NewOp(gconst.ObjFrom_ActMonopoly)

		if step < 1 || step > 6 {
			for _, v := range conf_g.MonopolyConsRandom {
				if v.Cg == actObj.GetConfGrp() {
					op.Dec(v.Id, v.N)
				}
			}
			step = rand_d.Int31n(6) + 1
		} else {
			for _, v := range conf_g.MonopolyConsNormal {
				if v.Cg == actObj.GetConfGrp() {
					op.Dec(v.Id, v.N)
				}
			}
		}

		if ec := op.CheckEnough(); ec != Err.OK {
			return ec, nil, step
		}

		next_p := plr_data.Pos + step
		if next_p >= (C_Monopoly_LevelMax - 1) {

			for _, v := range conf_lv.PassReward {
				op.Inc(v.Id, v.N)
			}

			rand_pos(plr)

			// 关卡达到最大值后，可继续重走最后一关，记录重复次数
			conf_nlv := gamedata.ConfActMonopolyLevel.Query(lv_to_seq(plr_data.CurLv + 1))
			if conf_nlv != nil {
				plr_data.CurLv = plr_data.CurLv + 1
			} else {
				plr_data.RepeatCnt += 1
			}

			plr_data.Pos = 0

			plr.SendMsg(&msg.GS_ActMonoPolyNextLv{
				Data: actObj.ToMsg(plr),
			})

			rwds := op.Apply().ToMsg()

			return Err.OK, rwds, step
		}

		next_tp := plr_data.PosInfo[next_p]
		if next_tp == 0 {
			plr_data.Pos = next_p

			slt := map[int32]int32{}
			for i, v := range conf_lv.BaseReward {
				slt[int32(i)] = v.Weight
			}

			idx := utils.PickWeightedMapId(slt)

			if len(conf_lv.BaseReward) > 0 {

				op.Inc(conf_lv.BaseReward[idx].Id, conf_lv.BaseReward[idx].N)
			}

			return Err.OK, op.Apply().ToMsg(), step
		}

		conf_tp := gamedata.ConfActMonopolyType.Query(next_tp)
		if conf_tp == nil {
			plr_data.Pos = next_p
			return Err.Activity_MonopolyTpNotFound, nil, step
		}

		p_id := int32(0)

		pick_id(plr_data, next_tp, func(id int32) {
			p_id = id
		})

		item := &item_t{
			Id:       p_id,
			ExpireTs: time.Now().Add(time.Duration(conf_tp.Time) * time.Second),
		}

		switch next_tp {
		case C_Monopoly_AdvBox: // 宝箱奖励单独推送（先推送，后续一起显示）
			conf := gamedata.ConfActMonopolyBox.Query(p_id)
			if conf != nil {
				boxOp := plr.GetBag().NewOp(gconst.ObjFrom_ActMonopoly)
				for _, v := range conf.BoxReward {
					boxOp.Inc(v.Id, v.N)
				}
				plr.SendMsg(&msg.GS_ActMonopolyBoxReward{
					Rewards: boxOp.Apply().ToMsg(),
				})
			}
			delete(plr_data.PosInfo, next_p)

		case C_Monopoly_AdvShop: // 折扣商店
			plr_data.Shops = append(plr_data.Shops, &shop_t{
				Id:       item.Id,
				ExpireTs: item.ExpireTs,
				DisCnt:   pick_shop_weight(item.Id),
			})
			res.Shops = shop_items_tomsg(plr_data.Shops)
			delete(plr_data.PosInfo, next_p)

		case C_Monopoly_AdvProblem: // 答题奇遇
			plr_data.Problems = append(plr_data.Problems, item)
			res.Problems = items_tomsg(plr_data.Problems)
			delete(plr_data.PosInfo, next_p)

		case C_Monopoly_AdvNormalBat: // 普通战斗
			plr_data.NormalBat = append(plr_data.NormalBat, item)
			res.NormalBat = items_tomsg(plr_data.NormalBat)
			delete(plr_data.PosInfo, next_p)

		case C_Monopoly_AdvHeroBat: // 英雄战斗
			plr_data.HeroBat = append(plr_data.HeroBat, &bat_t{
				Id:        item.Id,
				ExpireTs:  item.ExpireTs,
				MonsterLv: wlevel_aver(plr),
				HeroLv:    wlevel_aver(plr) + 1,
			})
			res.HeroBat = bat_items_tomsg(plr_data.HeroBat)
			delete(plr_data.PosInfo, next_p)
		}

		plr_data.Pos = next_p

		rwds := op.Apply().ToMsg()

		return Err.OK, rwds, step
	}()

	plr.SendMsg(res)

}

// 答题
func Answer(plr IPlayer, idx int32, ansIdx int32) (int32, *msg.Rewards) {
	if !actObj.Started() {
		return Err.Act_ActClosed, nil
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		return Err.Act_ActPlrDataNotFound, nil
	}

	conf_lv := gamedata.ConfActMonopolyLevel.Query(lv_to_seq(plr_data.CurLv))
	if conf_lv == nil {
		return Err.Failed, nil
	}

	if conf_lv.ConfGrp != actObj.GetConfGrp() {
		return Err.Act_ConfGrp, nil
	}

	op := plr.GetBag().NewOp(gconst.ObjFrom_ActMonopoly)

	p1 := plr_data.Problems[idx]
	if p1 == nil {
		return Err.Failed, nil
	}

	if p1.ExpireTs.Before(time.Now()) {
		return Err.Activity_MonopolyExpired, nil
	}

	conf_p := gamedata.ConfActMonopolyProblem.Query(p1.Id)
	if conf_p == nil {
		return Err.Activity_MonopolyPNotFound, nil
	}

	if (ansIdx + 1) != conf_p.Right {
		for _, v := range conf_p.ErrorReward {
			op.Inc(v.Id, v.N)
		}
	} else {
		for _, v := range conf_p.RightReward {
			op.Inc(v.Id, v.N)
		}
	}

	plr_data.Problems[idx].ExpireTs = time.Now()

	plr_data.Problems = append(plr_data.Problems[:idx], plr_data.Problems[idx+1:]...)

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

// 折扣商店购买，shopid为数组下标，itemidx为商品下标
func Buy(plr IPlayer, shopIdx int32, itemIdx int32) (int32, *msg.Rewards) {
	if !actObj.Started() {
		return Err.Act_ActClosed, nil
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		return Err.Act_ActPlrDataNotFound, nil
	}

	conf_lv := gamedata.ConfActMonopolyLevel.Query(lv_to_seq(plr_data.CurLv))
	if conf_lv == nil {
		return Err.Failed, nil
	}

	if conf_lv.ConfGrp != actObj.GetConfGrp() {
		return Err.Act_ConfGrp, nil
	}

	op := plr.GetBag().NewOp(gconst.ObjFrom_ActMonopoly)

	s1 := plr_data.Shops[shopIdx]
	if s1 == nil {
		return Err.Failed, nil
	}

	if s1.ExpireTs.Before(time.Now()) {
		return Err.Activity_MonopolyExpired, nil
	}

	conf_item := gamedata.ConfActMonopolyShop.Query(s1.Id)
	if conf_item == nil {
		return Err.Failed, nil
	}

	for idx := 0; idx < len(s1.AlreadyBuy); idx++ {
		if s1.AlreadyBuy[idx] == itemIdx {
			return Err.Activity_TakeBefore, nil
		}
	}

	for idx, v := range conf_item.ShopItem {

		if int32(idx) == itemIdx {

			conf_price := gamedata.ConfBasePrice.Query(conf_item.ShopItem[idx].Id)
			if conf_price == nil {
				return Err.Failed, nil
			}

			// TODO: float64问题，
			for _, v1 := range conf_price.Price {
				if v1.Id == conf_item.Ccy {
					op.Dec(v1.Id, int32(v1.N*float64(v.N)*float64(plr_data.Shops[shopIdx].DisCnt)/10000.0))
					break
				}
			}

			if ec := op.CheckEnough(); ec != Err.OK {
				return ec, nil
			}

			op.Inc(v.Id, v.N)

			plr_data.Shops[shopIdx].AlreadyBuy = append(plr_data.Shops[shopIdx].AlreadyBuy, itemIdx)

			if len(plr_data.Shops[shopIdx].AlreadyBuy) >= len(conf_item.ShopItem) {
				plr_data.Shops = append(plr_data.Shops[:shopIdx], plr_data.Shops[shopIdx+1:]...)
			}

			break
		}
	}

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

// 战斗，tp为类型（1：普通，2：英雄），tid为对应列表下标
func Battle(plr IPlayer, tp int32, tid int32, T *msg.TeamFormation, cb func(ec int32, replay *msg.BattleReplay, rwds *msg.Rewards)) {

	if !actObj.Started() {
		cb(Err.Act_ActClosed, nil, nil)
		return
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		cb(Err.Act_ActPlrDataNotFound, nil, nil)
		return
	}

	conf_lv := gamedata.ConfActMonopolyLevel.Query(lv_to_seq(plr_data.CurLv))
	if conf_lv == nil {
		cb(Err.Failed, nil, nil)
		return
	}

	conf_g := gamedata.ConfActivityPublic.Query(1)
	if conf_g == nil {
		cb(Err.Failed, nil, nil)
		return
	}

	if conf_lv.ConfGrp != actObj.GetConfGrp() {
		cb(Err.Act_ConfGrp, nil, nil)
		return
	}

	conf_w := gamedata.ConfWorldLevelM.Query(plr.GetWLevelLvNum() + 1)
	if conf_w == nil {
		cb(Err.Failed, nil, nil)
		return
	}

	op := plr.GetBag().NewOp(gconst.ObjFrom_ActMonopoly)

	bat1 := &item_t{}
	bat2 := &bat_t{}

	T2 := battle.NewMonsterTeam()
	Args := map[string]string{}

	switch tp {
	case C_Battle_Normal:

		bat1 = plr_data.NormalBat[tid]

		if len(plr_data.NormalBat) <= int(tid) {
			cb(Err.Failed, nil, nil)
			return
		}

		if bat1.ExpireTs.Before(time.Now()) {
			cb(Err.Activity_MonopolyExpired, nil, nil)
			return
		}
		conf := gamedata.ConfActMonopolyBattle.Query(bat1.Id)
		if conf == nil {
			cb(Err.Failed, nil, nil)
			return
		}

		lv := wlevel_aver(plr)

		for i, v := range conf.Monster {
			T2.AddMonster(v, lv, int32(i))
		}
		Args["Module"] = "MONOPOLY_MONSTER"

	case C_Battle_Hero:
		bat2 = plr_data.HeroBat[tid]

		if len(plr_data.HeroBat) <= int(tid) {
			cb(Err.Failed, nil, nil)
			return
		}

		if bat2.ExpireTs.Before(time.Now()) {
			cb(Err.Activity_MonopolyExpired, nil, nil)
			return
		}

		conf := gamedata.ConfActMonopolyHero.Query(bat2.Id)
		if conf == nil {
			cb(Err.Failed, nil, nil)
			return
		}

		for i, v := range conf.Monster {
			if i == int(conf.BossPos-1) {
				T2.AddMonster(v, bat2.HeroLv, int32(i))
			} else {
				T2.AddMonster(v, bat2.MonsterLv, int32(i))
			}
		}
		Args["Module"] = "MONOPOLY_HERO"
	default:
		cb(Err.Failed, nil, nil)
		return
	}

	input := &msg.BattleInput{
		T1:   plr.ToMsg_BattleTeam(T),
		T2:   T2.ToMsg_BattleTeam(),
		Args: Args,
	}

	battle.Fight(input, func(r *msg.BattleResult) {
		if r == nil {
			cb(Err.Common_BattleResError, nil, nil)
			return
		}

		if r.Winner == 1 {
			switch tp {

			case C_Battle_Normal:

				for _, v := range gamedata.ConfActMonopolyBattle.Query(bat1.Id).Reward {
					op.Inc(v.Id, v.N)
				}

				plr_data.NormalBat = append(plr_data.NormalBat[:tid], plr_data.NormalBat[tid+1:]...)

			case C_Battle_Hero:

				for _, v := range gamedata.ConfActMonopolyHero.Query(bat2.Id).Reward {
					op.Inc(v.Id, v.N)
				}

				plr_data.HeroBat = append(plr_data.HeroBat[:tid], plr_data.HeroBat[tid+1:]...)
			}
		}

		rwds := op.Apply().ToMsg()

		reply := &msg.BattleReplay{
			Ts: time.Now().Unix(),
			Bi: input,
			Br: r,
		}

		cb(Err.OK, reply, rwds)
	})

}

func seq_to_lv(seq int32) int32 {
	return seq%(actObj.GetConfGrp()*1000) - 1
}

// 关卡等级映射Seq
func lv_to_seq(lv int32) int32 {
	return (actObj.GetConfGrp() * 1000) + lv + 1
}

// ============================================================================
// 重新随机本关的60个点位
func rand_pos(plr IPlayer) {

	if !actObj.Started() {
		return
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		return
	}

	plr_data.PosInfo = map[int32]int32{}

	max_lv := len(gamedata.ConfActMonopolyLevelM.QueryItems(actObj.GetConfGrp()))

	conf_next := gamedata.ConfActMonopolyLevel.Query(lv_to_seq(plr_data.CurLv + 1))
	if conf_next == nil {
		// 如果到最大关卡后，随机数据为最后一关
		if (plr_data.CurLv) >= int32(max_lv) {
			conf_next = gamedata.ConfActMonopolyLevel.Query(lv_to_seq(plr_data.CurLv))
		} else {
			return
		}
	}

	conf_g := gamedata.ConfActivityPublic.Query(1)
	if conf_g == nil {
		return
	}

	slt_type := map[int32]int32{}
	for _, v := range conf_next.EventWeight {
		slt_type[v.Tp] = v.Weight
	}

	tp_cnt := map[int32]int32{}
	for i := int32(0); i < C_Monopoly_LevelMax; i += conf_next.CreateEvent {

		tp_id := utils.PickWeightedMapId(slt_type)

		if tp_id != 0 && (tp_cnt[tp_id]+1) > conf_next.EventMax {
			i -= conf_next.CreateEvent
			continue
		}

		pos := rand_d.Int31n(conf_next.CreateEvent) + i

		flag := false

		for k := int32(conf_g.MonopolyAdvSplitNum); k > 0; k-- {
			if plr_data.PosInfo[pos-k] != 0 {
				flag = true
				break
			}
		}

		if flag || pos == 0 || pos >= (C_Monopoly_LevelMax-1) {
			i -= conf_next.CreateEvent
			continue
		}

		plr_data.PosInfo[pos] = tp_id
		tp_cnt[tp_id]++

	}

	return
}

// 商店的折扣随机
func pick_shop_weight(sid int32) (dis int32) {

	dis = 10000

	conf := gamedata.ConfActMonopolyShop.Query(sid)
	if conf == nil {
		return dis
	}

	if len(conf.DiscountAndOdds) > 0 {
		dis = conf.DiscountAndOdds[0].Discount
	}

	slt := map[int32]int32{}
	for i, v := range conf.DiscountAndOdds {
		slt[int32(i)] = v.Weight
	}

	if aid := utils.PickWeightedMapId(slt); aid != 0 {
		dis = conf.DiscountAndOdds[aid].Discount
	}

	return
}

// 随机指定的类型的奇遇的id
func pick_id(plr_data *data_plr_t, tp int32, cb func(id int32)) {

	conf_lv := gamedata.ConfActMonopolyLevel.Query(lv_to_seq(plr_data.CurLv))

	for _, v := range conf_lv.EventWeight {
		if v.Tp == tp {
			switch v.Tp {
			case C_Monopoly_AdvBox:
				cb(pick_box(v.Grp))
				return

			case C_Monopoly_AdvHeroBat:
				cb(pick_hero(v.Grp))
				return

			case C_Monopoly_AdvNormalBat:
				cb(pick_bat(v.Grp))
				return

			case C_Monopoly_AdvProblem:
				cb(pick_problem(v.Grp))
				return

			case C_Monopoly_AdvShop:
				cb(pick_shop(v.Grp))
				return
			}
			break
		}
	}

}

// 宝箱随机
func pick_box(grp int32) int32 {

	id := int32(0)

	slt := map[int32]int32{}
	for i, v := range gamedata.ConfActMonopolyBoxM.QueryItems(grp) {
		if i == 0 {
			id = v.Seq
		}
		slt[v.Seq] = v.Weight

	}

	r_id := utils.PickWeightedMapId(slt)
	if r_id != 0 {
		id = r_id
	}

	if id == 0 {
		conf := gamedata.ConfActMonopolyBox.Query(1)
		if conf != nil {
			return conf.Seq
		}
	}

	return id
}

// 英雄随机
func pick_hero(grp int32) int32 {

	id := int32(0)

	slt := map[int32]int32{}
	for i, v := range gamedata.ConfActMonopolyHeroM.QueryItems(grp) {
		if i == 0 {
			id = v.Seq
		}
		slt[v.Seq] = v.Weight
	}

	r_id := utils.PickWeightedMapId(slt)
	if r_id != 0 {
		id = r_id
	}

	if id == 0 {
		conf := gamedata.ConfActMonopolyHero.Query(1)
		if conf != nil {
			return conf.Seq
		}
	}

	return id
}

// 战斗随机
func pick_bat(grp int32) int32 {

	id := int32(0)

	slt := map[int32]int32{}
	for i, v := range gamedata.ConfActMonopolyBattleM.QueryItems(grp) {
		if i == 0 {
			id = v.Seq
		}
		slt[v.Seq] = v.Weight
	}

	r_id := utils.PickWeightedMapId(slt)
	if r_id != 0 {
		id = r_id
	}

	if id == 0 {
		conf := gamedata.ConfActMonopolyBattle.Query(1)
		if conf != nil {
			return conf.Seq
		}
	}

	return id
}

// 答题随机
func pick_problem(grp int32) int32 {

	id := int32(0)

	slt := map[int32]int32{}
	for i, v := range gamedata.ConfActMonopolyProblemM.QueryItems(grp) {
		if i == 0 {
			id = v.Seq
		}
		slt[v.Seq] = v.Weight
	}

	r_id := utils.PickWeightedMapId(slt)
	if r_id != 0 {
		id = r_id
	}

	if id == 0 {
		conf := gamedata.ConfActMonopolyProblem.Query(1)
		if conf != nil {
			return conf.Seq
		}
	}

	return id
}

// 商店随机
func pick_shop(grp int32) int32 {

	id := int32(0)

	slt := map[int32]int32{}
	for i, v := range gamedata.ConfActMonopolyShopM.QueryItems(grp) {
		if i == 0 {
			id = v.Seq
		}
		slt[v.Seq] = v.Weight

	}

	r_id := utils.PickWeightedMapId(slt)
	if r_id != 0 {
		id = r_id
	}

	if id == 0 {
		conf := gamedata.ConfActMonopolyShop.Query(1)
		if conf != nil {
			return conf.Seq
		}
	}

	return id
}

func wlevel_aver(plr IPlayer) int32 {
	conf_w := gamedata.ConfWorldLevelM.Query(plr.GetWLevelLvNum())
	if conf_w == nil {
		return plr.GetLevel()
	}

	conf_g := gamedata.ConfActivityPublic.Query(1)
	if conf_g == nil {
		return plr.GetLevel()
	}

	sum := int32(0)
	for _, v := range conf_w.Monster {
		sum += v.Lv
	}

	return int32(float32(sum/int32(len(conf_w.Monster))) * conf_g.MonopolyBattleRatio)
}

func bat_items_tomsg(params []*bat_t) []*msg.ActMonopolyBatItem {
	ret := []*msg.ActMonopolyBatItem{}
	for i := 0; i < len(params); i++ {
		v := params[i]
		if v != nil {
			if v.ExpireTs.After(time.Now()) {
				ret = append(ret, &msg.ActMonopolyBatItem{
					Id:        v.Id,
					ExpireTs:  v.ExpireTs.Unix(),
					HeroLv:    v.HeroLv,
					MonsterLv: v.MonsterLv,
				})
			}
		}
	}
	return ret
}

// 返回符合的奇遇给客户端
func items_tomsg(params item_l) []*msg.ActMonopolyAdvItem {
	ret := []*msg.ActMonopolyAdvItem{}

	for i := 0; i < len(params); i++ {
		v := params[i]
		if v != nil {
			if v.ExpireTs.After(time.Now()) {
				ret = append(ret, &msg.ActMonopolyAdvItem{
					Id:       v.Id,
					ExpireTs: v.ExpireTs.Unix(),
				})
			}
		}
	}
	return ret
}

// 返回符合的商店奇遇给客户端
func shop_items_tomsg(params []*shop_t) []*msg.ActMonolyShopItem {

	ret := []*msg.ActMonolyShopItem{}

	for i := 0; i < len(params); i++ {
		v := params[i]
		if v != nil {
			if v.ExpireTs.After(time.Now()) {
				ret = append(ret, &msg.ActMonolyShopItem{
					DisCnt:     v.DisCnt,
					Id:         v.Id,
					ExpireTs:   v.ExpireTs.Unix(),
					AlreadyBuy: v.AlreadyBuy,
				})
			}
		}
	}
	return ret

}

// ============================================================================

func ActMonopolyInfo(plr IPlayer) *msg.ActMonopolyData {
	return actObj.ToMsg(plr)
}
