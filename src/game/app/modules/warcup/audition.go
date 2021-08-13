package warcup

import (
	"fw/src/core/log"
	"fw/src/core/sched/loop"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/arena"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/robot"
	"fw/src/game/msg"
	"math"
	"sort"
	"time"
)

// 海选赛
type audition_t struct {
	round     int32          // 第几轮
	rank_grps []*score_grp_t // 海选赛256分组排序
	guess     *guess_info_t  // 竞猜信息
}

// 海选积分赛组
type score_grp_t struct {
	plrid string
	sc    int32 // 海选排名分数
	ap    int32 // 阵容战力
}

// ============================================================================
// 海选赛

// 海选赛--开始
func (self *audition_t) begin() {
	self.add_plrs()
	if len(self.rank_grps) < c_Audition_MaxPlrs {
		force_close = true
		log.Warning("warcup audition player not enough", len(self.rank_grps))
		return
	}

	self.round_timer()
}

// 海选赛--添加参数人员(竞技场+机器人=256人)
func (self *audition_t) add_plrs() {
	rlv := int32(99999)

	// 竞技场玩家
	for _, row := range arena.ArenaMgr.Rows {
		plr := FindWarCupPlayer(row.PlrId)
		if plr == nil {
			continue
		}

		if rlv > plr.GetLevel() {
			rlv = plr.GetLevel()
		}

		T := plr.ToMsg_BattleTeam(plr.GetTeam(gconst.TeamType_Dfd), true)
		warcup_plrs[row.PlrId] = &warcup_plr_t{
			T: T,
		}

		self.rank_grps = append(self.rank_grps, &score_grp_t{
			plrid: plr.GetId(),
			ap:    T.Player.AtkPwr,
		})

		if len(warcup_plrs) >= c_Audition_MaxPlrs {
			break
		}
	}

	if rlv == 99999 {
		rlv = 30
	}

	// 机器人
	if len(warcup_plrs) < c_Audition_MaxPlrs {
		for {
			for _, v := range robot.RobotMgr.GetByLevel(rlv, 300) {
				T := v.ToMsg_BattleTeam(v.GetTeam(gconst.TeamType_Dfd), true)
				warcup_plrs[v.GetId()] = &warcup_plr_t{
					T: T,
				}

				self.rank_grps = append(self.rank_grps, &score_grp_t{
					plrid: v.GetId(),
					ap:    T.Player.AtkPwr,
				})

				if len(warcup_plrs) >= c_Audition_MaxPlrs {
					goto FULLPLAYER
				}
			}

			rlv--
			if rlv <= 0 {
				break
			}
		}
	FULLPLAYER:
	}

	// sort
	sort.Slice(self.rank_grps, func(i, j int) bool {
		return self.rank_grps[i].ap > self.rank_grps[j].ap
	})
}

// 海选赛--每轮timer
func (self *audition_t) round_timer() {
	for i := 0; i < c_Audition_MaxRound; i++ {
		fts := g_t1.Add(time.Duration(i*
			(c_Audition_Round_GuessTime+c_Audition_Round_FtResultTime+c_Audition_Round_RwdTime)) * time.Second)

		// guess & fight
		loop.SetTimeout(fts, func() {
			self.round_fight()

			plr_round_guess_free_score()
			// set stage
			broadcast_stage()
		})

		// fight result
		loop.SetTimeout(fts.Add(time.Duration(c_Audition_Round_GuessTime)*time.Second), func() {

			// set stage
			broadcast_stage()
		})

		// reward
		loop.SetTimeout(fts.Add(time.Duration(c_Audition_Round_GuessTime+c_Audition_Round_FtResultTime)*time.Second), func() {
			self.update_rank_score()
			round_award(self.guess)

			// set stage
			broadcast_stage()
		})
	}
}

// 海选赛--当前轮打架
func (self *audition_t) round_fight() {
	self.round++

	vs_grp := self.round_group()
	self.round_pick_guess(vs_grp)

	// 每秒一场战斗
	loop.TimeSlice(len(vs_grp), 1, 400, func(i int) {
		gplr1 := self.rank_grps[vs_grp[i][0]]
		gplr2 := self.rank_grps[vs_grp[i][1]]

		plr1 := warcup_plrs[gplr1.plrid]
		plr2 := warcup_plrs[gplr2.plrid]
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
			var addScore int32
			if r.Winner == 1 {
				plrw := FindWarCupPlayer(gplr1.plrid)
				if plrw != nil {
					addScore = calc_audition_fight_score(plrw, gplr1.ap, gplr2.ap)
				}

				plr1.Score += addScore
			} else {
				plrw := FindWarCupPlayer(gplr2.plrid)
				if plrw != nil {
					addScore = calc_audition_fight_score(plrw, gplr2.ap, gplr1.ap)
				}

				plr2.Score += addScore
			}

			replay := &msg.BattleReplay{
				Ts: time.Now().Unix(),
				Bi: input,
				Br: r,
			}

			// 更新对战信息
			for _, v := range warcup_vsdata {
				if v.stage == g_stage && v.round == self.round &&
					v.plrid1 == gplr1.plrid && v.plrid2 == gplr2.plrid {
					v.replay = replay
					v.addScore = addScore
				}
			}

			// 更新竞猜结果
			round_guess_result(self.guess, replay)
		})
	})
}

// 海选赛--当前轮分组
func (self *audition_t) round_group() [][2]int {
	vs_grp := [][2]int{}

	L := int(math.Exp2(float64(self.round + 1)))
	for i := 0; i < c_Audition_MaxPlrs; i += (2 * L) {
		for j := 0; j < L; j++ {
			vs_grp = append(vs_grp, [2]int{i + j, i + (2*L - 1 - j)})
		}
	}

	return vs_grp
}

// 海选赛--当前轮竞猜队伍
func (self *audition_t) round_pick_guess(bat_grp [][2]int) {
	self.guess = &guess_info_t{
		guess_plrs: make(map[string]*guess_one_t),
	}

	min := float64(99999999)
	for _, v := range bat_grp {
		plrid1 := self.rank_grps[v[0]].plrid
		plrid2 := self.rank_grps[v[1]].plrid

		plr1 := warcup_plrs[plrid1]
		plr2 := warcup_plrs[plrid2]
		if plr1 == nil || plr2 == nil {
			continue
		}

		// 生成对战信息
		add_vs_data(g_stage, self.round, plrid1, plrid2, plr1, plr2)

		// 战力相差最小
		n := math.Abs(float64(plr1.T.Player.AtkPwr - plr2.T.Player.AtkPwr))
		if min > n {
			min = n
			self.guess.plrid1 = plrid1
			self.guess.plrid2 = plrid2
		}
	}
}

func (self *audition_t) update_rank_score() {
	for _, v := range self.rank_grps {
		plr_data := warcup_plrs[v.plrid]
		if plr_data == nil {
			continue
		}

		v.sc = plr_data.Score
	}
}

// 海选赛--排序
func (self *audition_t) rank_sort() {
	if len(audition.rank_grps) == 0 {
		return
	}

	sort.Slice(audition.rank_grps, func(i, j int) bool {
		if audition.rank_grps[i].sc == audition.rank_grps[j].sc {
			return audition.rank_grps[i].ap > audition.rank_grps[j].ap
		} else {
			return audition.rank_grps[i].sc > audition.rank_grps[j].sc
		}
	})
}
