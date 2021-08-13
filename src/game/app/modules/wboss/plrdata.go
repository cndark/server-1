package wboss

import (
	"fw/src/core"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"sort"
	"strings"
	"time"
)

// ============================================================================

type WBoss struct {
	Ver time.Time

	RwdMaxDmgTaken map[int32]bool
	BuffIds        []int32

	plr IPlayer
}

// ============================================================================

func NewWBoss() *WBoss {
	return &WBoss{
		RwdMaxDmgTaken: make(map[int32]bool),
	}
}

// ============================================================================

func (self *WBoss) Init(plr IPlayer) {
	self.plr = plr
}

func (self *WBoss) check_data() {
	if !g_t0.IsZero() && !self.Ver.Equal(g_t0) {
		self.Ver = g_t0

		self.RwdMaxDmgTaken = make(map[int32]bool)
		self.BuffIds = nil
	}
}

func (self *WBoss) Fight(tf *msg.TeamFormation, f func(*msg.GS_WBossFight_R)) int {
	fail := func(ec int32) int {
		f(&msg.GS_WBossFight_R{
			ErrorCode: ec,
		})
		return 0
	}

	// check stage
	if g_stage != c_Stage_Start {
		return fail(Err.Failed)
	}

	// can play ?
	if !self.plr.IsModuleOpen(gconst.ModuleId_WBoss) {
		return fail(Err.Plr_ModuleLocked)
	}

	// check team
	if !self.plr.IsTeamFormationValid(tf) {
		return fail(Err.Plr_TeamInvalid)
	}

	// check data
	self.check_data()

	// check counter
	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_WBossFight)
	cop.DecCounter(gconst.Cnt_WBossFight, 1)
	if ec := cop.CheckEnough(); ec != Err.OK {
		return fail(ec)
	}

	// boss conf
	bossid := current_boss_id()
	if bossid == 0 {
		return fail(Err.Failed)
	}

	conf_boss := gamedata.ConfWorldBoss.Query(bossid)
	if conf_boss == nil {
		return fail(Err.Failed)
	}

	// prepare buff ids
	var buff_list string

	L := len(self.BuffIds)
	if L > 0 {
		var sb strings.Builder
		sb.Grow(L * 10)

		for _, v := range self.BuffIds {
			sb.WriteString(",")
			sb.WriteString(core.I32toa(v))
		}
		buff_list = sb.String()[1:]
	}

	// fight
	bi := &msg.BattleInput{
		T1: self.plr.ToMsg_BattleTeam(tf),
		T2: battle.NewMonsterTeam().AddMonster(conf_boss.MonsterId, conf_boss.Lv, 6).ToMsg_BattleTeam(),
		Args: map[string]string{
			"Module":    "WORLD_BOSS",
			"RoundType": core.I32toa(conf_boss.RoundType),
			"buffs.1":   buff_list,
		},
	}

	battle.Fight(bi, func(r *msg.BattleResult) {
		// check
		if r == nil {
			fail(Err.Failed)
			return
		}

		// ok
		replay := &msg.BattleReplay{Bi: bi, Br: r}
		dmg := core.Atof64(r.Args["dmg_total.1"])

		// calc rewards
		for _, v := range conf_boss.DamageReward {
			if dmg >= v.N {
				conf_rwd := gamedata.ConfWorldBossRewardList.Query(v.Id)
				if conf_rwd != nil {
					for _, v2 := range conf_rwd.Reward {
						cop.Inc(v2.Id, v2.N)
					}
				}
			}
		}

		// apply rewards
		rwd := cop.Apply().ToMsg()

		// calc jf
		var jfadd int32
		conf_g := gamedata.ConfGlobalPublic.Query(1)
		if conf_g != nil {
			jfadd = conf_g.WorldBossPointFloors + int32(dmg*float64(conf_g.WorldBossPointRatio))
			g_sd.update_rank(self.plr.GetId(), jfadd)
		}

		// max dmg stats
		if dmg > g_sd.MaxDmg.Dmg {
			g_sd.MaxDmg.Id = self.plr.GetId()
			g_sd.MaxDmg.Dmg = dmg
		}

		// callback
		f(&msg.GS_WBossFight_R{
			ErrorCode: Err.OK,
			Replay:    replay,
			Rewards:   rwd,
			JfAdd:     jfadd,
		})
	})

	return 0
}

func (self *WBoss) TakeMaxDmgRwd(n int32) (ec int32, rwd *msg.Rewards) {
	// check stage
	if g_stage < c_Stage_Start {
		return Err.Failed, nil
	}

	// can play ?
	if !self.plr.IsModuleOpen(gconst.ModuleId_WBoss) {
		return Err.Plr_ModuleLocked, nil
	}

	// check data
	self.check_data()

	// already taken ?
	if self.RwdMaxDmgTaken[n] {
		return Err.WBoss_MaxDmgRwdTaken, nil
	}

	// boss conf
	bossid := current_boss_id()
	if bossid == 0 {
		return Err.Failed, nil
	}

	conf_boss := gamedata.ConfWorldBoss.Query(bossid)
	if conf_boss == nil {
		return Err.Failed, nil
	}

	// can take ?
	if n < 1 || n > int32(len(conf_boss.MaxDamageReward)) || g_sd.MaxDmg.Dmg < conf_boss.MaxDamageReward[n-1].N {
		return Err.WBoss_MaxDmgRwdNotAllowed, nil
	}

	// take it
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_WBossMaxDmgRwd)
	conf_rwd := gamedata.ConfWorldBossRewardList.Query(conf_boss.MaxDamageReward[n-1].Id)
	if conf_rwd != nil {
		for _, v := range conf_rwd.Reward {
			op.Inc(v.Id, v.N)
		}
		self.BuffIds = append(self.BuffIds, conf_rwd.BuffId...)
	}
	rwd = op.Apply().ToMsg()

	// mark taken
	self.RwdMaxDmgTaken[n] = true

	// ok
	return Err.OK, rwd
}

func (self *WBoss) GetSelfRank() (rk int32, jf int32) {
	if g_stage >= c_Stage_Start {
		e := g_sd.rank_idx[self.plr.GetId()]
		if e != nil {
			L := len(g_sd.Rank)
			start := sort.Search(L, func(i int) bool {
				return g_sd.Rank[i].Jf <= e.Jf
			})
			for i := start; i < L && g_sd.Rank[i].Jf == e.Jf; i++ {
				if g_sd.Rank[i].Id == self.plr.GetId() {
					return int32(i) + 1, e.Jf
				}
			}
		}
	}

	return 0, 0
}

func (self *WBoss) GetRank() (ret []*msg.WBossRankRow) {
	if g_stage < c_Stage_Start {
		return
	}

	L := len(g_sd.Rank)
	if L > 100 {
		L = 100
	}

	ret = make([]*msg.WBossRankRow, 0, L)
	for _, v := range g_sd.Rank[:L] {
		var pi *msg.PlayerSimpleInfo
		iplr := utils.LoadPlayer(v.Id)
		if iplr != nil {
			pi = iplr.(IPlayer).ToMsg_SimpleInfo()
		}

		ret = append(ret, &msg.WBossRankRow{
			Player: pi,
			Jf:     v.Jf,
		})
	}

	return
}

// ============================================================================

func (self *WBoss) ToMsg() *msg.WBossData {
	return &msg.WBossData{
		Stage:   int32(g_stage),
		Ts2:     g_t2.Unix(),
		Summary: self.ToMsg_Summary(),
	}
}

func (self *WBoss) ToMsg_Summary() *msg.GS_WBossGetSummary_R {
	return &msg.GS_WBossGetSummary_R{
		BossId:         current_boss_id(),
		RwdMaxDmgTaken: self.RwdMaxDmgTaken,
	}
}
