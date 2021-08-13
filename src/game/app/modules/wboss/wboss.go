package wboss

import (
	. "fw/src/core/math"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/mdata"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"sort"
	"time"
)

// ============================================================================

const (
	c_Stage_Close  = 0
	c_Stage_Init   = 1
	c_Stage_Open   = 2
	c_Stage_Start  = 3
	c_Stage_End    = 4
	c_Stage_Reward = 5
)

// ============================================================================

var (
	g_stage int       // 当前阶段
	g_t0    time.Time // 本期版本
	g_t2    time.Time // 当前阶段结束时间 (下个阶段开始时间)

	g_sd *wboss_t
)

// ============================================================================

type wboss_t struct {
	T0           time.Time               // 版本
	Num          int32                   // 期数
	MaxDmg       *wboss_maxdmg_plr_t     // 最大伤害信息
	Rank         []*wboss_plr_t          // 排行榜
	rank_idx     map[string]*wboss_plr_t // 排行榜索引
	RankMailSent bool                    // 排行榜奖励邮件是否发送
}

type wboss_maxdmg_plr_t struct {
	Id  string
	Dmg float64
}

type wboss_plr_t struct {
	Id string
	Jf int32
}

// ============================================================================

func new_svrdata() interface{} {
	return &wboss_t{
		MaxDmg: &wboss_maxdmg_plr_t{},
	}
}

func on_svrdata_loaded() {
	g_sd = mdata.Get(NAME).(*wboss_t)
	g_sd.init()
}

func (self *wboss_t) init() {
	// build rank index
	self.rank_idx = make(map[string]*wboss_plr_t)
	for _, v := range self.Rank {
		self.rank_idx[v.Id] = v
	}
}

func (self *wboss_t) reset() {
	self.T0 = g_t0
	self.Num++
	self.MaxDmg = &wboss_maxdmg_plr_t{}
	self.Rank = nil
	self.rank_idx = make(map[string]*wboss_plr_t)
	self.RankMailSent = false
}

func (self *wboss_t) update_rank(id string, jfadd int32) {
	// find
	e := self.rank_idx[id]
	if e == nil {
		e = &wboss_plr_t{
			Id: id,
		}
		self.Rank = append(self.Rank, e)
		self.rank_idx[e.Id] = e
	}

	// update
	e.Jf += jfadd

	// sort
	sort.Slice(self.Rank, func(i, j int) bool {
		return self.Rank[i].Jf > self.Rank[j].Jf
	})
}

// ============================================================================

func on_stage(stg string, t0, t1, t2 time.Time, f func(bool)) {
	// set ts
	g_t0 = t0
	g_t2 = t2

	if stg == "init" {
		g_stage = c_Stage_Init
		g_sd.reset()
	} else {
		f(g_sd.T0.Equal(t0))
	}
}

// ============================================================================

func stage_open() {
	set_stage(c_Stage_Open)
}

func stage_start() {
	set_stage(c_Stage_Start)
}

func stage_end() {
	set_stage(c_Stage_End)
}

func stage_reward() {
	set_stage(c_Stage_Reward)
	send_rank_reward_mails()
}

func stage_close() {
	set_stage(c_Stage_Close)
}

// ============================================================================

func set_stage(v int) {
	g_stage = v
	utils.BroadcastPlayers(&msg.GS_WBossStageChange{
		Stage: int32(v),
		Ts2:   g_t2.Unix(),
	})
}

func current_boss_id() int32 {
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return 0
	}

	L := len(conf.WorldBossArrange)
	if L == 0 {
		return 0
	}

	i := int(g_sd.Num)
	if i < 1 {
		return 0
	} else if i > L {
		i = L
	}

	return conf.WorldBossArrange[i-1]
}

func send_rank_reward_mails() {
	// already sent ?
	if g_sd.RankMailSent {
		return
	}

	// conf
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return
	}

	// send
	L := int32(len(g_sd.Rank))

	for _, v := range gamedata.ConfWorldBossRankReward.Items() {
		a := MaxInt32(1, v.Rank[0].Low)
		b := MinInt32(L, v.Rank[0].High)

		for i := a; i <= b; i++ {
			e := g_sd.Rank[i-1]

			// find player
			iplr := utils.LoadPlayer(e.Id)
			if iplr == nil {
				continue
			}

			// send
			ml := mail.New(iplr).SetKey(conf.WorldBossRankMailId)
			for _, v2 := range v.Reward {
				ml.AddAttachment(v2.Id, float64(v2.N))
			}
			ml.AddDictInt32("rank", i)
			ml.Send()
		}
	}

	// mark as sent
	g_sd.RankMailSent = true
}

func ToMsg_MaxDmgInfo() *msg.WBossMaxDmgInfo {
	ret := &msg.WBossMaxDmgInfo{}

	if g_sd.MaxDmg.Id != "" {
		iplr := utils.LoadPlayer(g_sd.MaxDmg.Id)
		if iplr != nil {
			ret.Player = iplr.(IPlayer).ToMsg_SimpleInfo()
		}
	}
	ret.Dmg = g_sd.MaxDmg.Dmg

	return ret
}
