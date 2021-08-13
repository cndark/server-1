package crusade

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/robot"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"sort"
	"time"
)

// ============================================================================

// 远征
type Crusade struct {
	VerTs      time.Time         // 版本时间
	LvNum      int32             // 通过层数
	HpLoss     map[int64]float64 // 损失血量百分比[seq]hp
	Enemies    []*enemy_t        // N波怪
	BoxTaken   []int32           // 奖励领取
	Difficulty int32             // 难度

	plr IPlayer
}

type enemy_t struct {
	T      *msg.BattleTeam   // 战队
	HpLoss map[int32]float64 // 损失血量百分比[pos]hp
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_MOpen, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		mid := args[1].(int32)

		if mid == gconst.ModuleId_Crusade {
			plr.GetCrusade().CheckVerTs()
		}
	})
}

func NewCrusade() *Crusade {
	return &Crusade{
		HpLoss: make(map[int64]float64),
	}
}

// ============================================================================

func (self *Crusade) Init(plr IPlayer) {
	self.plr = plr
}

// 确认版本信息
func (self *Crusade) CheckVerTs() {

	// 历史最大战力为0，就忽略
	if self.plr.GetAttainObjVal(gconst.AttainId_MaxAtkPwr) <= 0 {
		return
	}

	// reset
	if self.VerTs.IsZero() || !self.VerTs.Equal(CrusadeMgr.VerTs) {
		self.reset()
	}

	// check null
	if len(self.Enemies) == 0 {
		self.gen_enemies()
	}
}

func (self *Crusade) reset() {
	self.HpLoss = make(map[int64]float64)

	if self.LvNum < int32(len(gamedata.ConfCrusade.Items())) {
		self.Difficulty = 0
	}

	self.LvNum = 0
	self.BoxTaken = []int32{}

	self.gen_enemies()
	self.VerTs = CrusadeMgr.VerTs
}

// 重新生成对手
func (self *Crusade) gen_enemies() {
	N := len(gamedata.ConfCrusade.Items())
	self.Enemies = make([]*enemy_t, N)
	lv := self.plr.GetLevel()
	dup := map[string]bool{self.plr.GetId(): true}
	max_atkpwr := self.plr.GetAttainObjVal(gconst.AttainId_MaxAtkPwr)

	for _, conf := range gamedata.ConfCrusade.Items() {
		if conf.Id > int32(N) {
			continue
		}

		idx := self.Difficulty
		L := int32(len(conf.PowerCorrect))
		if self.Difficulty >= L {
			idx = L - 1
		}

		rate := conf.PowerCorrect[idx]
		low := float64(rate.Low) * float64(max_atkpwr)
		high := float64(rate.High) * float64(max_atkpwr)

		// player
		plrid := ""
		utils.ForEachLoadedPlayerBreakable(func(plr interface{}) bool {
			fplr := plr.(ICrusadePlayer)
			if dup[fplr.GetId()] {
				return true
			}

			if fplr.GetAtkPwr() > 0 && (fplr.GetAtkPwr() >= int32(low)) &&
				(fplr.GetAtkPwr() <= int32(high)) {
				plrid = fplr.GetId()
				return false
			}

			return true
		})

		// robot
		if plrid == "" {
			for _, bot := range robot.RobotMgr.Index {
				if dup[bot.GetId()] {
					continue
				}

				if bot.GetAtkPwr() >= int32(low) &&
					bot.GetAtkPwr() <= int32(high) {
					plrid = bot.GetId()
					break
				}

				if (bot.Lv-8) <= lv &&
					(bot.Lv+8) >= lv {

					plrid = bot.GetId()
					break
				}
			}
		}

		// 随机一个
		if plrid == "" {
			for _, bot := range robot.RobotMgr.Index {
				if dup[bot.GetId()] {
					continue
				}

				plrid = bot.GetId()
				break
			}
		}

		eplr := FindCrusadePlayer(plrid)
		if eplr != nil {
			t := eplr.ToMsg_BattleTeam(eplr.GetTeam(gconst.TeamType_Dfd), true)

			enemy := &enemy_t{
				HpLoss: make(map[int32]float64),
				T:      t,
			}

			self.Enemies[conf.Id-1] = enemy
			dup[plrid] = true
		}
	}

	// atkpwr sort
	sort.Slice(self.Enemies, func(i, j int) bool {
		ei, ej := self.Enemies[i], self.Enemies[j]
		if ei == nil || ej == nil {
			return false
		}

		return self.Enemies[i].T.Player.AtkPwr <= self.Enemies[j].T.Player.AtkPwr
	})
}

func (self *Crusade) IsBoxTake(id int32) bool {
	for _, v := range self.BoxTaken {
		if v == id {
			return true
		}
	}

	return false
}

func (self *Crusade) TakeBox(id int32) (int32, *msg.Rewards) {
	conf := gamedata.ConfCrusade.Query(id)
	if conf == nil || len(conf.Chest) == 0 {
		return Err.Failed, nil
	}

	if id > self.LvNum {
		return Err.Crusade_NotPass, nil
	}

	if self.IsBoxTake(id) {
		return Err.Plr_TakenBefore, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_CrusadeBox)
	for _, v := range conf.Chest {
		op.Inc(v.Id, v.N)
	}

	self.BoxTaken = append(self.BoxTaken, id)

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

func (self *Crusade) Fight(T *msg.TeamFormation,
	cb func(ec int32, replay *msg.BattleReplay, rwds *msg.Rewards, hploss map[int64]float64)) {

	if !self.plr.IsModuleOpen(gconst.ModuleId_Crusade) {
		cb(Err.Plr_ModuleLocked, nil, nil, nil)
		return
	}

	if !IsStart() {
		cb(Err.Crusade_End, nil, nil, nil)
		return
	}

	conf := gamedata.ConfCrusade.Query(self.LvNum + 1)
	if conf == nil {
		cb(Err.Failed, nil, nil, nil)
		return
	}

	if self.LvNum >= int32(len(self.Enemies)) {
		cb(Err.Crusade_PassAll, nil, nil, nil)
		return
	}

	if !self.plr.IsTeamFormationValid(T) {
		cb(Err.Plr_TeamInvalid, nil, nil, nil)
		return
	}

	// args
	args := map[string]string{
		"Module":    "CRUSADE",
		"RoundType": core.I32toa(conf.RoundType),
	}

	allDeath := true
	for seq, pos := range T.Formation {
		v := self.HpLoss[seq]
		s := `init_hp.1.` + core.I32toa(pos)
		args[s] = core.F64toa(v)

		if v < 1 {
			allDeath = false
		}
	}

	if allDeath {
		cb(Err.Plr_TeamAllDeath, nil, nil, nil)
		return
	}

	enemy := self.Enemies[self.LvNum]
	for pos, v := range enemy.HpLoss {
		s := `init_hp.2.` + core.I32toa(int32(pos))
		args[s] = core.F64toa(v)
	}

	// gen input
	input := &msg.BattleInput{
		T1:   self.plr.ToMsg_BattleTeam(T),
		T2:   enemy.T,
		Args: args,
	}

	// fight
	battle.Fight(input, func(r *msg.BattleResult) {
		if r == nil {
			cb(Err.Common_BattleResError, nil, nil, nil)
			return
		}

		var rwds *msg.Rewards
		if r.Winner == 1 {
			op := self.plr.GetBag().NewOp(gconst.ObjFrom_CrusadeFight)

			for _, v := range conf.Reward {
				lv := self.plr.GetLevel()
				if lv >= v.Low && lv <= v.High {
					op.Inc(v.Id, v.N)
				}
			}

			rwds = op.Apply().ToMsg()
			self.LvNum++
		}

		// update hp
		hploss := make(map[int64]float64)
		t1, t2 := battle.BattleResultHpLoss(r.Args)
		for pos1, v := range t1 {
			for seq, pos2 := range T.Formation {
				if pos1 == pos2 {
					self.HpLoss[seq] = v
					hploss[seq] = v
					break
				}
			}
		}

		for pos1, v := range t2 {
			enemy.HpLoss[pos1] = v
		}

		// res
		cb(Err.OK, &msg.BattleReplay{
			Ts: time.Now().Unix(),
			Bi: input,
			Br: r,
		}, rwds, hploss)

		// evt
		evtmgr.Fire(gconst.Evt_CrusadeFight, self.plr, r.Winner == 1)
	})
}

// ============================================================================

func (self *Crusade) GetInfo() *msg.GS_CrusadeGetInfo_R {
	self.CheckVerTs()

	return &msg.GS_CrusadeGetInfo_R{
		Data: self.ToMsg(),
	}
}

func (self *Crusade) ToMsg() *msg.CrusadeData {
	ret := &msg.CrusadeData{
		Stage:    CrusadeMgr.stage,
		Ts1:      CrusadeMgr.Ts1.Unix(),
		Ts2:      CrusadeMgr.Ts2.Unix(),
		LvNum:    self.LvNum,
		BoxTaken: self.BoxTaken,
		HpLoss:   self.HpLoss,
	}

	for _, v := range self.Enemies {
		if v == nil || v.T == nil {
			continue
		}

		ret.Enemies = append(ret.Enemies, &msg.CrusadeEnemy{
			Plr: v.T.Player,
		})
	}

	return ret
}
