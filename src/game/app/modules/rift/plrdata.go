package rift

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

// ============================================================================

// 裂隙
type Rift struct {
	Monster     *monster_t      // 怪物
	MineCnt     map[int32]int32 // 每天占领矿的次数[id]cnt
	RecentMines []int64
	IsFirst     bool // 是否第一次探索

	plr IPlayer
}

type monster_t struct {
	Id int32 // 当前怪id
	Lv int32 // 怪等级
}

// ============================================================================

func init() {
	// reset daily
	evtmgr.On(gconst.Evt_PlrResetDaily, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		plr.GetRift().reset_daily()
	})
}

func NewRift() *Rift {
	return &Rift{
		Monster: &monster_t{},
		MineCnt: make(map[int32]int32),
		IsFirst: true,
	}
}

// ============================================================================

func (self *Rift) Init(plr IPlayer) {
	self.plr = plr
}

func (self *Rift) reset_daily() {
	self.MineCnt = make(map[int32]int32)
}

func (self *Rift) GetMineCnt(id int32) int32 {
	return self.MineCnt[id]
}

func (self *Rift) AddMineCnt(id int32, v int32) {
	self.MineCnt[id] += v

	if self.MineCnt[id] < 0 {
		self.MineCnt[id] = 0
	}
}

func (self *Rift) GetFirst() bool {
	return self.IsFirst
}

func (self *Rift) SetFirst() {
	self.IsFirst = false
}

func (self *Rift) SetMonster(id int32, lv int32) {
	self.Monster.Id = id
	self.Monster.Lv = lv
}

// 探索生成怪物
func (self *Rift) ExploreMonster() {
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil || len(conf.RiftMonsterLv) == 0 {
		return
	}

	slt := make(map[int32]int32)
	for _, v := range gamedata.ConfRiftMonster.Items() {
		slt[v.Id] = v.Weight
	}

	id := utils.PickWeightedMapId(slt)
	if id == 0 {
		return
	}

	self.Monster.Id = id

	min := conf.RiftMonsterLv[0].Min
	max := conf.RiftMonsterLv[0].Max
	r := min + float32(rand_rift.Float64())*(max-min)

	self.Monster.Lv = int32(float32(self.plr.GetLevel()) * r)
	if self.Monster.Lv <= 0 || self.Monster.Lv > gamedata.ConfLimitM.Query().MaxPlrLv {
		self.Monster.Lv = self.plr.GetLevel()
	}

	// res
	self.plr.SendMsg(&msg.GS_RiftMonsterNew{
		Monster: &msg.RiftMonster{
			Id: self.Monster.Id,
			Lv: self.Monster.Lv,
		},
	})

}

// 没探索到就清掉这个
func (self *Rift) CleanMonster() {
	self.Monster.Id = 0
	self.Monster.Lv = 0
}

func (self *Rift) FightMonster(T *msg.TeamFormation, cb func(ec int32, replay *msg.BattleReplay, rwds *msg.Rewards)) {
	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil || len(conf_g.RiftCost) < 4 {
		cb(Err.Failed, nil, nil)
		return
	}

	conf_m := gamedata.ConfRiftMonster.Query(self.Monster.Id)
	if conf_m == nil {
		cb(Err.Rift_MonsterNotFound, nil, nil)
		return
	}

	if !self.plr.IsTeamFormationValid(T) {
		cb(Err.Plr_TeamInvalid, nil, nil)
		return
	}

	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_RiftMonster)
	cop.DecCounter(gconst.Cnt_PlayerStrength, int64(conf_g.RiftCost[C_RiftType_Monster]))
	if ec := cop.CheckEnough(); ec != Err.OK {
		cb(ec, nil, nil)
		return
	}

	// gen input
	T2 := battle.NewMonsterTeam()
	for i, v := range conf_m.Monster {
		T2.AddMonster(v, self.Monster.Lv, int32(i))
	}

	input := &msg.BattleInput{
		T1: self.plr.ToMsg_BattleTeam(T),
		T2: T2.ToMsg_BattleTeam(),
		Args: map[string]string{
			"Module": "RIFT_MONSTER",
			"Round":  core.I32toa(conf_m.Round),
		},
	}

	// fight
	battle.Fight(input, func(r *msg.BattleResult) {
		if r == nil {
			cb(Err.Common_BattleResError, nil, nil)
			return
		}

		cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_RiftMonster)
		cop.DecCounter(gconst.Cnt_PlayerStrength, int64(conf_g.RiftCost[C_RiftType_Monster]))
		if ec := cop.CheckEnough(); ec != Err.OK {
			cb(ec, nil, nil)
			return
		}

		conf_lv := gamedata.ConfPlayerUp.Query(self.Monster.Lv)
		if conf_lv == nil {
			cb(Err.Failed, nil, nil)
			return
		}

		if r.Winner == 1 {
			for _, v := range conf_m.Reward {
				if rand_rift.Float32() < v.Odds {
					a := float32(1)
					for _, vv := range conf_lv.RiftMonsterRatio {
						if vv.Id == v.Id {
							a = vv.N
						}
					}

					cop.Inc(v.Id, int32(float32(v.N)*a))
				}
			}

			evtmgr.Fire(gconst.Evt_RiftMonsterFight, self.plr, self.Monster.Id)

			self.Monster.Id = 0
			self.Monster.Lv = 0
		}

		rwds := cop.Apply().ToMsg()

		replay := &msg.BattleReplay{
			Ts: time.Now().Unix(),
			Bi: input,
			Br: r,
		}

		cb(Err.OK, replay, rwds)
	})
}

// ============================================================================

func (self *Rift) ToMsg() *msg.RiftMonster {
	return &msg.RiftMonster{
		Id: self.Monster.Id,
		Lv: self.Monster.Lv,
	}
}
