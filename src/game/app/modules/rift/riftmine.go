package rift

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math"
	"time"
)

// ============================================================================

// 裂隙--矿
type rift_mine_t struct {
	TotalCnt int64
	Mines    map[int64]*mine_t
	Plrs     map[string]*plr_mine_t
}

// 玩家当前矿区
type plr_mine_t struct {
	Id     int32                 // 当前区矿
	Lv     int32                 // 当前区矿等级
	RMines map[int32]*plr_mine_r // 自己有关系的矿[id]seq
}

type plr_mine_r struct {
	Occ int64 // 自己占的
	Rob int64 // 被抢
}

// 单个矿
type mine_t struct {
	Seq     int64           // 矿seq
	Id      int32           // 矿id
	Lv      int32           // 矿等级
	FinTs   time.Time       // 完成时间
	CurTeam *msg.BattleTeam // 当前玩家阵容信息

	isBattle bool
}

// ============================================================================

func (self *rift_mine_t) data_loaded() {
}

// ============================================================================

// 探索生成矿
func (self *rift_mine_t) ExploreRiftMine(plr IPlayer) {
	plr_data := self.Plrs[plr.GetId()]
	if plr_data == nil {
		self.Plrs[plr.GetId()] = &plr_mine_t{
			RMines: make(map[int32]*plr_mine_r),
		}
	}

	var mine *mine_t
	if rand_rift.Float64() < 0.5 {
		mine = self.rand_empty(plr)
	} else {
		mine = self.rand_plr_mine(plr)
	}

	if mine == nil {
		return
	}

	// res
	plr.SendMsg(&msg.GS_RiftMineNew{
		Mine: mine.ToMsg(),
	})
}

// 没探索到就清掉空矿
func (self *rift_mine_t) CleanEmptyMine(plr IPlayer) {
	plr_data := self.Plrs[plr.GetId()]
	if plr_data == nil {
		plr_data = &plr_mine_t{
			RMines: make(map[int32]*plr_mine_r),
		}

		self.Plrs[plr.GetId()] = plr_data
	}

	plr_data.Id = 0
	plr_data.Lv = 0
}

// 随机一个空矿
func (self *rift_mine_t) rand_empty(plr IPlayer) *mine_t {
	plr_data := self.Plrs[plr.GetId()]

	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil || len(conf_g.RiftMineLv) == 0 {
		return nil
	}

	slt := make(map[int32]int32)
	for _, v := range gamedata.ConfRiftMine.Items() {
		slt[v.Id] += v.Weight
	}

	id := utils.PickWeightedMapId(slt)
	conf_m := gamedata.ConfRiftMine.Query(id)
	if conf_m == nil {
		return nil
	}

	// lv
	min := conf_g.RiftMineLv[0].Min
	max := conf_g.RiftMineLv[0].Max
	r := min + float32(rand_rift.Float64())*(max-min)

	lv := int32(float32(plr.GetLevel()) * r)
	if lv <= 0 || lv > gamedata.ConfLimitM.Query().MaxPlrLv {
		lv = plr.GetLevel()
	}

	plr_data.Id = id
	plr_data.Lv = lv

	return &mine_t{
		Id: id,
		Lv: lv,
	}
}

// 随机一个玩家防守矿
func (self *rift_mine_t) rand_plr_mine(plr IPlayer) *mine_t {
	seqs := []int64{}
	now := time.Now()
	for seq, m := range self.Mines {
		if now.After(m.FinTs) {
			continue
		}

		if m.CurTeam == nil ||
			m.CurTeam.Player.Id == plr.GetId() {
			continue
		}

		for _, v := range plr.GetRift().RecentMines {
			if v == seq {
				continue
			}
		}

		seqs = append(seqs, seq)
	}

	L := len(seqs)
	if L > 0 {
		i := rand_rift.Intn(L)

		conf_g := gamedata.ConfGlobalPublic.Query(1)
		if conf_g == nil {
			return nil
		}

		// 将最近N次随机到的矿放入队列中
		if len(plr.GetRift().RecentMines) > int(conf_g.RiftMineNorepeat) {
			plr.GetRift().RecentMines = append(plr.GetRift().RecentMines[1:conf_g.RiftMineNorepeat], seqs[i])
		} else {
			plr.GetRift().RecentMines = append(plr.GetRift().RecentMines, seqs[i])
		}

		return self.Mines[seqs[i]]
	}

	return nil
}

// 布防空矿
func (self *rift_mine_t) Occupy(plr IPlayer, seq int64, T *msg.TeamFormation,
	cb func(ec int32, mine *msg.RiftMine, replay *msg.BattleReplay)) {
	if !plr.IsTeamFormationValid(T) {
		cb(Err.Plr_TeamInvalid, nil, nil)
		return
	}

	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil || len(conf_g.RiftCost) < 4 {
		cb(Err.Failed, nil, nil)
		return
	}

	cop := plr.GetCounter().NewOp(gconst.ObjFrom_RiftMine)
	cop.DecCounter(gconst.Cnt_PlayerStrength, int64(conf_g.RiftCost[C_RiftType_Monster]))
	if ec := cop.CheckEnough(); ec != Err.OK {
		cb(ec, nil, nil)
		return
	}

	// 占领自己的空矿
	if seq == 0 {
		ec, m := self.occupy_self(plr, T)
		if ec == Err.OK {
			cop.Apply()
		}

		cb(ec, m, nil)
		return
	}

	// 抢别人的矿
	plr_data := self.Plrs[plr.GetId()]
	if plr_data == nil {
		cb(Err.Rift_MineExploreBefore, nil, nil)
		return
	}

	mine := self.Mines[seq]
	if mine == nil || mine.CurTeam == nil {
		cb(Err.OK, nil, nil) // 这个客户端错误不要被拦截
		return
	}

	if mine.isBattle {
		cb(Err.Rift_MineIsBattle, nil, nil)
		return
	}

	if mine.FinTs.Before(time.Now()) {
		cb(Err.Rift_MineFin, nil, nil)
		return
	}

	if mine.CurTeam.Player.Id == plr.GetId() {
		cb(Err.Rift_MineIsSelf, nil, nil)
		return
	}

	if self.has_same_mine(plr.GetId(), plr_data, mine.Id) {
		cb(Err.Rift_MineOccupyNoMore, nil, nil)
		return
	}

	if mine.isBattle {
		cb(Err.Rift_BoxIsBattle, nil, nil)
		return
	}

	// 进攻
	eplr := load_player(mine.CurTeam.Player.Id)
	if eplr == nil {
		cb(Err.Plr_NotFound, nil, nil)
		return
	}

	input := &msg.BattleInput{
		T1: plr.ToMsg_BattleTeam(T),
		T2: mine.CurTeam,
		Args: map[string]string{
			"Module":    "RIFT_MINE",
			"RoundType": "2",
		},
	}

	mine.isBattle = true

	// fight
	battle.Fight(input, func(r *msg.BattleResult) {
		mine.isBattle = false

		if r == nil {
			cb(Err.Common_BattleResError, nil, nil)
			return
		}

		cop := plr.GetCounter().NewOp(gconst.ObjFrom_RiftMine)
		cop.DecCounter(gconst.Cnt_PlayerStrength, int64(conf_g.RiftCost[C_RiftType_Monster]))
		if ec := cop.CheckEnough(); ec != Err.OK {
			cb(ec, nil, nil)
			return
		}

		// 成功
		if r.Winner == 1 {
			mine.CurTeam = input.T1
			mine.FinTs = self.get_mine_fints(plr, mine.Id)

			plr_data.update_rmine(true, mine.Id, mine.Seq)

			eplr_data := self.Plrs[eplr.GetId()]
			if eplr_data != nil {
				eplr_data.update_rmine(false, mine.Id, mine.Seq)
			}

			// evt
			evtmgr.Fire(gconst.Evt_RiftMineOccupy, plr, mine.Id)

			plr.GetRift().AddMineCnt(mine.Id, 1)
			eplr.GetRift().AddMineCnt(mine.Id, -1)
		}

		cop.Apply()

		replay := &msg.BattleReplay{
			Ts: time.Now().Unix(),
			Bi: input,
			Br: r,
		}

		m := mine.ToMsg()
		cb(Err.OK, m, replay)

		if r.Winner == 1 {
			eplr.SendMsg(&msg.GS_RiftMineOccupied{
				Mine: m,
			})

			// mail
			conf_m := gamedata.ConfMail.Query(conf_g.RiftMineRobMail)
			if conf_m != nil {
				m := mail.New(eplr).SetKey(conf_g.RiftMineRobMail)
				m.AddDict("player", plr.GetName())
				m.Send()
			}
		}

	})
}

// 占领自己的空矿
func (self *rift_mine_t) occupy_self(plr IPlayer, T *msg.TeamFormation) (int32, *msg.RiftMine) {
	plr_data := self.Plrs[plr.GetId()]
	if plr_data == nil {
		return Err.Rift_MineExploreBefore, nil
	}

	if plr_data.Id == 0 {
		return Err.Rift_MineNotFound, nil

	}

	if self.has_same_mine(plr.GetId(), plr_data, plr_data.Id) {
		return Err.Rift_MineOccupyNoMore, nil
	}

	mine := &mine_t{
		Seq:     self.TotalCnt + 1,
		Id:      plr_data.Id,
		Lv:      plr_data.Lv,
		FinTs:   self.get_mine_fints(plr, plr_data.Id),
		CurTeam: plr.ToMsg_BattleTeam(T, true),
	}

	plr_data.Id = 0
	plr_data.Lv = 0
	plr_data.update_rmine(true, mine.Id, mine.Seq)

	self.Mines[mine.Seq] = mine
	self.TotalCnt++

	plr.GetRift().AddMineCnt(mine.Id, 1)

	evtmgr.Fire(gconst.Evt_RiftMineOccupy, plr, mine.Id)

	return Err.OK, mine.ToMsg()
}

// 获取完成时间
func (self *rift_mine_t) get_mine_fints(plr IPlayer, id int32) time.Time {
	cnt := plr.GetRift().GetMineCnt(id)

	min := int32(30)
	conf := gamedata.ConfRiftMine.Query(id)
	if conf != nil {
		L := int32(len(conf.TimeMin))
		if L > 0 {
			if cnt >= L {
				cnt = L - 1
			}

			min = conf.TimeMin[cnt]
		}
	}

	return time.Now().Add(time.Duration(min) * time.Minute)
}

// 是否有同类型的矿
func (self *rift_mine_t) has_same_mine(plrid string, plr_data *plr_mine_t, id int32) bool {
	r := plr_data.RMines[id]
	if r != nil && r.Occ != 0 {
		mine := self.Mines[r.Occ]
		if (mine != nil) &&
			(mine.CurTeam != nil) &&
			(mine.CurTeam.Player.Id == plrid) {
			return true
		}
	}

	return false
}

// 放弃矿
func (self *rift_mine_t) CancelOccupy(plr IPlayer, seq int64) int32 {
	plr_data := self.Plrs[plr.GetId()]
	if plr_data == nil {
		return Err.Rift_MineExploreBefore
	}

	mine := self.Mines[seq]
	if mine == nil || mine.CurTeam == nil {
		return Err.Rift_MineNotFound
	}

	if mine.CurTeam.Player.Id != plr.GetId() {
		return Err.Rift_MineNotSelf
	}

	delete(plr_data.RMines, mine.Id)
	delete(self.Mines, seq)

	plr.GetRift().AddMineCnt(mine.Id, -1)

	return Err.OK
}

// 领奖
func (self *rift_mine_t) TakeRewards(plr IPlayer, seq int64) (int32, *msg.Rewards) {
	plr_data := self.Plrs[plr.GetId()]
	if plr_data == nil {
		return Err.Rift_MineExploreBefore, nil
	}

	mine := self.Mines[seq]
	if mine == nil || mine.CurTeam == nil {
		return Err.Rift_MineNotFound, nil
	}

	if mine.CurTeam.Player.Id != plr.GetId() {
		return Err.Rift_MineNotSelf, nil
	}

	id := mine.Id
	conf_m := gamedata.ConfRiftMine.Query(id)
	if conf_m == nil {
		return Err.Failed, nil
	}

	conf_lv := gamedata.ConfPlayerUp.Query(mine.Lv)
	if conf_lv == nil {
		return Err.Failed, nil
	}

	// reward
	op := plr.GetBag().NewOp(gconst.ObjFrom_RiftMine)

	// fast take
	now := time.Now()
	if mine.FinTs.After(time.Now()) {
		sub := mine.FinTs.Sub(now).Minutes()
		n := int64(math.Ceil(sub))

		conf_g := gamedata.ConfGlobalPublic.Query(1)
		if conf_g != nil {
			n *= int64(conf_g.RiftMineFastCost)
		}

		if n <= 0 {
			return Err.Common_TimeNotUp, nil
		}

		op.Dec(gconst.Diamond, n)
		if ec := op.CheckEnough(); ec != Err.OK {
			return ec, nil
		}
	}

	for _, v := range conf_m.Reward {
		if rand_rift.Float32() < v.Odds {
			a := float32(1)
			for _, vv := range conf_lv.RiftMineRatio {
				if vv.Id == v.Id {
					a = vv.N
				}
			}

			op.Inc(v.Id, int32(float32(v.N)*a))
		}
	}

	delete(plr_data.RMines, id)
	delete(self.Mines, seq)

	rwds := op.Apply().ToMsg()

	// evt
	evtmgr.Fire(gconst.Evt_RiftMineTake, plr, id)

	return Err.OK, rwds
}

func (self *rift_mine_t) ToMsg(plr IPlayer) *msg.RiftPlrMine {
	plr_data := self.Plrs[plr.GetId()]
	if plr_data == nil {
		return nil
	}

	ret := &msg.RiftPlrMine{
		Id:     plr_data.Id,
		Lv:     plr_data.Lv,
		RMines: make(map[int32]*msg.RiftMine),
	}

	for id, v := range plr_data.RMines {
		if v.Occ != 0 {
			mine := self.Mines[v.Occ]
			if mine != nil {
				ret.RMines[id] = mine.ToMsg()
			}
		} else if v.Rob != 0 {
			mine := self.Mines[v.Rob]
			if mine != nil {
				ret.RMines[id] = mine.ToMsg()
			}
		}
	}

	return ret
}

// ============================================================================

// 更新自己相关矿信息
func (self *plr_mine_t) update_rmine(isOcc bool, id int32, seq int64) {
	r := self.RMines[id]
	if r == nil {
		r = &plr_mine_r{}
		self.RMines[id] = r
	}

	if isOcc {
		r.Occ = seq
		r.Rob = 0
	} else {
		r.Occ = 0
		r.Rob = seq
	}
}

func (self *mine_t) ToMsg() *msg.RiftMine {
	ret := &msg.RiftMine{
		Seq:   self.Seq,
		Id:    self.Id,
		Lv:    self.Lv,
		FinTs: self.FinTs.Unix(),
	}

	if self.CurTeam != nil {
		plr := load_player(self.CurTeam.Player.Id)
		if plr != nil {
			ret.CurPlr = plr.ToMsg_SimpleInfo()
			ret.Fighters = self.CurTeam.Fighters
		}
	}

	return ret
}
