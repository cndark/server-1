package guild

import (
	"fw/src/game/app/gamedata"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"sort"
)

// ============================================================================

type boss_t struct {
	Cur    *boss_hist_t // 当前信息
	HpLeft float64      // 当前 boss 剩余血量百分比

	Hist map[int32]*boss_hist_t // 历史. [num]

	gld *Guild
}

type boss_hist_t struct {
	Num  int32              // 层数
	Dmgs map[string]float64 // 伤害榜单. [plrid]dmg
}

// ============================================================================

func new_boss() *boss_t {
	return &boss_t{
		Cur: &boss_hist_t{
			Num:  1,
			Dmgs: make(map[string]float64),
		},
		HpLeft: 1,
		Hist:   make(map[int32]*boss_hist_t),
	}
}

func (self *boss_t) init(g *Guild) {
	self.gld = g
}

func (self *boss_t) update(plrid string, dmg float64, hpleft float64) {
	// add player dmg
	self.Cur.Dmgs[plrid] += dmg

	// check boss death
	if hpleft > 0 { // update hp
		self.HpLeft = hpleft
	} else { // dead
		// send rank rewards mails
		self.send_rank_rewards_mails()

		// update history
		self.Hist[self.Cur.Num] = self.Cur

		// next boss
		self.next_boss()
	}
}

func (self *boss_t) send_rank_rewards_mails() {
	// conf
	conf := gamedata.ConfGuildBoss.Query(self.Cur.Num)
	if conf == nil {
		return
	}

	// sort
	type dmg_t struct {
		id  string
		dmg float64
	}
	arr := make([]*dmg_t, 0, len(self.Cur.Dmgs))
	for id, dmg := range self.Cur.Dmgs {
		arr = append(arr, &dmg_t{id, dmg})
	}

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].dmg > arr[j].dmg ||
			arr[i].dmg == arr[j].dmg && arr[i].id < arr[j].id
	})

	// send mails
	for i, v := range arr {
		plr := utils.LoadPlayer(v.id)
		if plr == nil {
			continue
		}

		// new mail
		m := mail.New(plr)
		m.SetKey(205)

		// attachments
		rk := int32(i) + 1
		for _, w := range conf.KillReward {
			if w.A <= rk && rk <= w.B {
				m.AddAttachment(w.Id, float64(w.N))
			}
		}

		// dict
		m.AddDictInt32("num", self.Cur.Num)
		m.AddDictInt32("rank", rk)

		// send
		m.Send()
	}
}

func (self *boss_t) next_boss() {
	// next num
	num := self.Cur.Num
	conf := gamedata.ConfGuildBoss.Query(num + 1)
	if conf != nil {
		num++
	}

	// new current
	self.Cur = &boss_hist_t{
		Num:  num,
		Dmgs: make(map[string]float64),
	}
	self.HpLeft = 1
}

func (self *boss_t) ToMsg_History(num int32) (ec int32, info *msg.GuildBossHistory) {
	// find
	hist := self.Hist[num]
	if hist == nil {
		return Err.Guild_BossHistNotFound, nil
	}

	return Err.OK, hist.ToMsg()
}

// ============================================================================

func (self *boss_hist_t) ToMsg() *msg.GuildBossHistory {
	arr := make([]*msg.GuildBossDmg, 0, len(self.Dmgs))
	for id, dmg := range self.Dmgs {
		iplr := utils.LoadPlayer(id)
		if iplr == nil {
			continue
		}

		arr = append(arr, &msg.GuildBossDmg{
			Plr: iplr.(IPlayer).ToMsg_SimpleInfo(),
			Dmg: dmg,
		})
	}

	return &msg.GuildBossHistory{
		Num:  self.Num,
		Dmgs: arr,
	}
}
