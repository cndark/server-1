package warcup

import (
	"fw/src/core/sched/loop"
	"fw/src/game/app/modules/battle"
	"fw/src/game/msg"
	"math"
	"time"
)

// 淘汰赛
type knockout_t struct {
	round     int32         // 第几轮
	rank_grps *kk8_grp_t    // 淘汰赛64强分组打8人赛
	guess     *guess_info_t // 竞猜信息
}

// 8人赛淘汰组
type kk8_grp_t struct {
	round1   []*kk8_fight_t // 前8组
	round2   []*kk8_fight_t // 前4组
	round3   []*kk8_fight_t // 前2组
	champion []string       // 冠军
}

// 8人赛淘汰各轮战斗队伍
type kk8_fight_t struct {
	tt    [2]string // 参战人员
	vsseq int32     // 对战索引
}

// ============================================================================
// 8人一组淘汰制

// 开启timer
func (self *knockout_t) round_timer() {
	for i := 0; i < c_Knockout_MaxRound; i++ {
		fts := g_t1.Add(time.Duration(i*
			(c_Knockout_Round_GuessTime+c_Knockout_Round_FtResultTime+c_Knockout_Round_RwdTime)) * time.Second)

		// guess & fight
		loop.SetTimeout(fts, func() {
			self.round_fight()

			plr_round_guess_free_score()
			// set stage
			broadcast_stage()
		})

		// fight result
		loop.SetTimeout(fts.Add(time.Duration(c_Knockout_Round_GuessTime)*time.Second), func() {

			// set stage
			broadcast_stage()
		})

		// reward
		loop.SetTimeout(fts.Add(time.Duration(c_Knockout_Round_GuessTime+c_Knockout_Round_FtResultTime)*time.Second), func() {
			round_award(self.guess)

			// set stage
			broadcast_stage()
		})
	}
}

// 当前轮打架
func (self *knockout_t) round_fight() {
	self.round++
	self.round_pickguess()

	switch self.round {
	case 1:
		self.round_fight_1()
	case 2:
		self.round_fight_2()
	case 3:
		self.round_fight_3()
	default:
	}
}

// 当前轮竞猜队伍
func (self *knockout_t) round_pickguess() {
	self.guess = &guess_info_t{
		guess_plrs: make(map[string]*guess_one_t),
	}

	min := float64(99999999)
	switch self.round {
	case 1:
		for _, bat := range self.rank_grps.round1 {
			plr1 := warcup_plrs[bat.tt[0]]
			plr2 := warcup_plrs[bat.tt[1]]
			if plr1 == nil || plr2 == nil {
				continue
			}

			// 生成对战信息
			bat.vsseq = add_vs_data(g_stage, self.round, bat.tt[0], bat.tt[1], plr1, plr2)

			// 战力相差最小
			n := math.Abs(float64(plr1.T.Player.AtkPwr - plr2.T.Player.AtkPwr))
			if min > n {
				min = n
				self.guess.plrid1 = bat.tt[0]
				self.guess.plrid2 = bat.tt[1]
			}
		}
	case 2:
		for _, bat := range self.rank_grps.round2 {
			plr1 := warcup_plrs[bat.tt[0]]
			plr2 := warcup_plrs[bat.tt[1]]
			if plr1 == nil || plr2 == nil {
				continue
			}

			// 生成对战信息
			bat.vsseq = add_vs_data(g_stage, self.round, bat.tt[0], bat.tt[1], plr1, plr2)

			// 战力相差最小
			n := math.Abs(float64(plr1.T.Player.AtkPwr - plr2.T.Player.AtkPwr))
			if min > n {
				min = n
				self.guess.plrid1 = bat.tt[0]
				self.guess.plrid2 = bat.tt[1]
			}
		}
	case 3:
		for _, bat := range self.rank_grps.round3 {
			plr1 := warcup_plrs[bat.tt[0]]
			plr2 := warcup_plrs[bat.tt[1]]
			if plr1 == nil || plr2 == nil {
				continue
			}

			// 生成对战信息
			bat.vsseq = add_vs_data(g_stage, self.round, bat.tt[0], bat.tt[1], plr1, plr2)

			// 战力相差最小
			n := math.Abs(float64(plr1.T.Player.AtkPwr - plr2.T.Player.AtkPwr))
			if min > n {
				min = n
				self.guess.plrid1 = bat.tt[0]
				self.guess.plrid2 = bat.tt[1]
			}
		}

	default:
	}
}

// 第1轮8进4
func (self *knockout_t) round_fight_1() {
	loop.TimeSlice(len(self.rank_grps.round1), 2, 1000, func(i int) {

		kk8_ft := self.rank_grps.round1[i]
		plrid1 := kk8_ft.tt[0]
		plrid2 := kk8_ft.tt[1]

		plr1 := warcup_plrs[plrid1]
		plr2 := warcup_plrs[plrid2]
		if plr1 == nil || plr2 == nil {
			return
		}

		input := &msg.BattleInput{
			T1: plr1.T,
			T2: plr2.T,
			Args: map[string]string{
				"Module":    "WAR_CUP",
				"RoundType": "3",
			},
		}

		battle.Fight(input, func(r *msg.BattleResult) {
			if r == nil {
				return
			}

			// 加积分
			winner := plrid1
			if r.Winner != 1 {
				winner = plrid2
			}

			replay := &msg.BattleReplay{
				Ts: time.Now().Unix(),
				Bi: input,
				Br: r,
			}

			// 更新对战信息
			vsdata := warcup_vsdata[kk8_ft.vsseq]
			if vsdata != nil {
				vsdata.replay = replay
			}

			// 生成下一轮选手
			idx0 := i / 2 // 新战斗队顺序
			idx1 := i % 2 // 新战斗左右位置

			// 新战斗人员
			next_kk := self.rank_grps.round2[idx0]
			if next_kk == nil {
				next_kk = &kk8_fight_t{}
				self.rank_grps.round2[idx0] = next_kk
			}

			if idx1 == 0 {
				next_kk.tt[0] = winner
			} else {
				next_kk.tt[1] = winner
			}

			// 更新竞猜结果
			round_guess_result(self.guess, replay)
		})
	})
}

// 第2轮4进2
func (self *knockout_t) round_fight_2() {
	loop.TimeSlice(len(self.rank_grps.round2), 2, 1000, func(i int) {
		kk8_ft := self.rank_grps.round2[i]
		plrid1 := kk8_ft.tt[0]
		plrid2 := kk8_ft.tt[1]

		plr1 := warcup_plrs[plrid1]
		plr2 := warcup_plrs[plrid2]
		if plr1 == nil || plr2 == nil {
			return
		}

		input := &msg.BattleInput{
			T1: plr1.T,
			T2: plr2.T,
			Args: map[string]string{
				"Module":    "WAR_CUP",
				"RoundType": "3",
			},
		}

		battle.Fight(input, func(r *msg.BattleResult) {
			if r == nil {
				return
			}

			// 加积分
			winner := plrid1
			if r.Winner != 1 {
				winner = plrid2
			}

			replay := &msg.BattleReplay{
				Ts: time.Now().Unix(),
				Bi: input,
				Br: r,
			}

			// 更新对战信息
			vsdata := warcup_vsdata[kk8_ft.vsseq]
			if vsdata != nil {
				vsdata.replay = replay
			}

			// 生成下一轮选手
			idx0 := i / 2 // 新战斗队顺序
			idx1 := i % 2 // 新战斗左右位置

			// 新战斗人员
			next_kk := self.rank_grps.round3[idx0]
			if next_kk == nil {
				next_kk = &kk8_fight_t{}
				self.rank_grps.round3[idx0] = next_kk
			}

			if idx1 == 0 {
				next_kk.tt[0] = winner
			} else {
				next_kk.tt[1] = winner
			}

			// 更新竞猜结果
			round_guess_result(self.guess, replay)
		})
	})
}

// 第3轮2进1
func (self *knockout_t) round_fight_3() {
	loop.TimeSlice(len(self.rank_grps.round3), 1, 1000, func(i int) {
		kk8_ft := self.rank_grps.round3[i]
		plrid1 := kk8_ft.tt[0]
		plrid2 := kk8_ft.tt[1]

		plr1 := warcup_plrs[plrid1]
		plr2 := warcup_plrs[plrid2]
		if plr1 == nil || plr2 == nil {
			return
		}

		input := &msg.BattleInput{
			T1: plr1.T,
			T2: plr2.T,
			Args: map[string]string{
				"Module":    "WAR_CUP",
				"RoundType": "3",
			},
		}

		battle.Fight(input, func(r *msg.BattleResult) {
			if r == nil {
				return
			}

			// 加积分
			winner := plrid1
			if r.Winner != 1 {
				winner = plrid2
			}

			replay := &msg.BattleReplay{
				Ts: time.Now().Unix(),
				Bi: input,
				Br: r,
			}

			// 更新对战信息
			vsdata := warcup_vsdata[kk8_ft.vsseq]
			if vsdata != nil {
				vsdata.replay = replay
			}

			// 生成冠军
			self.rank_grps.champion[i] = winner

			// 更新竞猜结果
			round_guess_result(self.guess, replay)
		})
	})
}
