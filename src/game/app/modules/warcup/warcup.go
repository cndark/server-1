package warcup

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/math"
	. "fw/src/core/math"
	"fw/src/core/sched/loop"
	"fw/src/core/wordsfilter"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"sort"
	"time"
)

// ============================================================================

const (
	c_Audition_MaxPlrs            = 256 // 海选赛参数人数
	c_Audition_MaxRound           = 6   // 海选赛最大轮数; 2^(Max+1)要<=(256/2)才满足分组
	c_Audition_Round_GuessTime    = 90  // 海选赛每轮竞猜时长
	c_Audition_Round_FtResultTime = 180 // 海选赛每轮打架结果时长
	c_Audition_Round_RwdTime      = 30  // 海选赛每轮结束时长

	c_Knockout_MaxRound           = 3   // 淘汰赛最大轮数
	c_Knockout_Round_GuessTime    = 90  // 淘汰赛每轮竞猜时长
	c_Knockout_Round_FtResultTime = 180 // 淘汰赛每轮打架结果时长
	c_Knockout_Round_RwdTime      = 30  // 淘汰赛每轮结束时长

	c_Top64_MaxPlrs = 64 // 64强淘汰赛参数人数
	c_Top8_MaxPlrs  = 8  // 8强冠军赛参数人数

	c_Round_Piece_Guess       = 1 // 本轮小阶段:竞猜(战斗也开始跑)
	c_Round_Piece_FightResult = 2 // 本轮小阶段:战斗结果
	c_Round_Piece_Award       = 3 // 本轮小阶段:看比赛发奖

)

var (
	force_close = true // 未经历open阶段本期不开
)

// ============================================================================

var (
	g_stage int32     // 当前阶段
	g_t0    time.Time // 本期版本
	g_t1    time.Time // 当前阶段开始时间
	g_t2    time.Time // 当前阶段结束时间(下个阶段开始时间)

	warcup_plrs     map[string]*warcup_plr_t     // 参赛人员共256
	warcup_plrguess map[string]*plr_guess_data_t // 所有玩家竞猜信息

	warcup_vsseq  int32                // 对战信息索引
	warcup_vsdata map[int32]*vs_data_t // 对战信息集合

	warcup_chat []*chat_ont_t // 聊天数据

	audition *audition_t // 海选赛
	top64    *knockout_t // 淘汰赛(64强)
	top8     *knockout_t // 冠军赛(8强)

)

// 参赛人员
type warcup_plr_t struct {
	T      *msg.BattleTeam // 上阵阵容
	Score  int32           // 总分数
	VsSeqs []int32         // 自己的对战信息
}

// 对战信息
type vs_data_t struct {
	vsseq    int32
	stage    int32
	round    int32
	plrid1   string
	plrid2   string
	addScore int32
	replay   *msg.BattleReplay
}

// 每轮公共竞猜信息
type guess_info_t struct {
	guess_plrs map[string]*guess_one_t // 竞猜输赢[plrid]
	plrid1     string                  // 左边
	plrid2     string                  // 右边
	win1       int32                   // 压左边赢的人数
	replay     *msg.BattleReplay       // 竞猜的战斗结果
}

// 竞猜内容
type guess_one_t struct {
	winner int32 // 输赢状况
	num    int32 // 竞猜数
}

// 玩家竞猜信息
type plr_guess_data_t struct {
	score   int32             // 总积分
	records []*guess_record_t // 竞猜记录
}

// 竞猜记录
type guess_record_t struct {
	stage int32  // 阶段
	round int32  // 轮数
	name  string // 谁赢
	num   int32  // 压的数量
	add   int32  // 赢的数量
	isFin bool   // 是否出结果
}

// 聊天记录
type chat_ont_t struct {
	name    string // 玩家名字
	content string // 内容
}

// ============================================================================
// 海选赛

func new_audition() {
	audition = &audition_t{
		guess: &guess_info_t{
			guess_plrs: make(map[string]*guess_one_t),
		},
	}

	audition.begin()
}

// ============================================================================
// 64强淘汰赛

func new_top64() {
	top64 = &knockout_t{
		rank_grps: &kk8_grp_t{
			round1:   make([]*kk8_fight_t, 32),
			round2:   make([]*kk8_fight_t, 16),
			round3:   make([]*kk8_fight_t, 8),
			champion: make([]string, 8),
		},
		guess: &guess_info_t{
			guess_plrs: make(map[string]*guess_one_t),
		},
	}

	top64_begin()
}

// 64强淘汰赛开始
func top64_begin() {
	if len(audition.rank_grps) < c_Top64_MaxPlrs {
		force_close = true
		log.Warning("warcup knockout player not enough", len(audition.rank_grps))
		return
	}
	audition.rank_sort()

	// pick 64
	for i := 0; i < 8; i++ {
		grp := [8]string{}
		for j := 0; j < 8; j++ {
			grp[j] = audition.rank_grps[i+j*8].plrid
		}

		top64.rank_grps.round1[i*4] = &kk8_fight_t{tt: [2]string{grp[0], grp[7]}}
		top64.rank_grps.round1[i*4+1] = &kk8_fight_t{tt: [2]string{grp[2], grp[5]}}
		top64.rank_grps.round1[i*4+2] = &kk8_fight_t{tt: [2]string{grp[1], grp[6]}}
		top64.rank_grps.round1[i*4+3] = &kk8_fight_t{tt: [2]string{grp[3], grp[4]}}
	}

	top64.round_timer()
}

// ============================================================================
// 8强冠军赛

func new_top8() {
	top8 = &knockout_t{
		rank_grps: &kk8_grp_t{
			round1:   make([]*kk8_fight_t, 4),
			round2:   make([]*kk8_fight_t, 2),
			round3:   make([]*kk8_fight_t, 1),
			champion: make([]string, 1),
		},
		guess: &guess_info_t{
			guess_plrs: make(map[string]*guess_one_t),
		},
	}

	top8_begin()
}

// 8强冠军赛开始
func top8_begin() {
	// 添加参赛赛人员
	arr := []*score_grp_t{}
	for _, plrid := range top64.rank_grps.champion {
		plr := warcup_plrs[plrid]
		if plr == nil {
			continue
		}
		arr = append(arr, &score_grp_t{
			plrid: plrid,
			ap:    plr.T.Player.AtkPwr,
		})
	}

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].ap > arr[j].ap
	})

	if len(arr) < c_Top8_MaxPlrs {
		force_close = true
		log.Warning("warcup final player not enough", len(arr))
		return
	}

	top8.rank_grps.round1[0] = &kk8_fight_t{tt: [2]string{arr[0].plrid, arr[7].plrid}}
	top8.rank_grps.round1[1] = &kk8_fight_t{tt: [2]string{arr[2].plrid, arr[5].plrid}}
	top8.rank_grps.round1[2] = &kk8_fight_t{tt: [2]string{arr[1].plrid, arr[6].plrid}}
	top8.rank_grps.round1[3] = &kk8_fight_t{tt: [2]string{arr[3].plrid, arr[4].plrid}}

	top8.round_timer()
}

// ============================================================================

// 定时发送开始提醒邮件
func send_open_mail() {
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil || len(conf.WarCupForcastMail) == 0 {
		return
	}

	mid := conf.WarCupForcastMail[0].MailId
	conf_m := gamedata.ConfMail.Query(mid)
	if conf_m == nil {
		return
	}

	ts := g_t1.Add(time.Duration(conf.WarCupForcastMail[0].Sec) * time.Second)
	loop.SetTimeout(ts, func() {
		utils.ForEachLoadedPlayer(func(plr interface{}) {
			m := mail.New(plr).SetKey(mid)
			m.Send()
		})
	})
}

// 计算海选战斗积分
func calc_audition_fight_score(plr IWarCupPlayer, atk_sf, atk_rl int32) int32 {
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return 0
	}

	a := float64(conf.WarCupWinBaseScore) * math.MaxFloat64((float64(atk_rl)/float64(atk_sf)),
		float64(conf.WarCupWinFloorsRatio))

	b := math.MinFloat64(float64(plr.GetArenaScore())/float64(conf.WarCupWinScoreP1),
		float64(conf.WarCupWinScoreP2))

	return int32(a * (1 + b))
}

// 添加对战信息
func add_vs_data(stage, round int32, plrid1, plrid2 string, plr1, plr2 *warcup_plr_t) int32 {
	warcup_vsseq++
	warcup_vsdata[warcup_vsseq] = &vs_data_t{
		vsseq:  warcup_vsseq,
		stage:  stage,
		round:  round,
		plrid1: plrid1,
		plrid2: plrid2,
	}

	plr1.VsSeqs = append(plr1.VsSeqs, warcup_vsseq)
	plr2.VsSeqs = append(plr2.VsSeqs, warcup_vsseq)

	return warcup_vsseq
}

// 阶段本轮小阶段和它的结束时间
func round_piece() (round int32, rp int32, endts int64) {
	now := time.Now()
	round, t1, t2, t3 := int32(0), 0, 0, 0

	switch g_stage {
	case c_Stage_Audition:
		if audition == nil {
			return
		}

		round = audition.round
		t1 = c_Audition_Round_GuessTime
		t2 = c_Audition_Round_FtResultTime
		t3 = c_Audition_Round_RwdTime

	case c_Stage_Top64:
		if top64 == nil {
			return
		}

		round = top64.round
		t1 = c_Knockout_Round_GuessTime
		t2 = c_Knockout_Round_FtResultTime
		t3 = c_Knockout_Round_RwdTime
	case c_Stage_Top8:
		if top8 == nil {
			return
		}

		round = top8.round
		t1 = c_Knockout_Round_GuessTime
		t2 = c_Knockout_Round_FtResultTime
		t3 = c_Knockout_Round_RwdTime
	}

	if round == 0 {
		return
	}

	fts := g_t1.Add(time.Duration((round-1)*int32(t1+t2+t3)) * time.Second)

	p1_endts := fts.Add(time.Duration(t1) * time.Second)
	p2_endts := fts.Add(time.Duration(t1+t2) * time.Second)
	p3_endts := fts.Add(time.Duration(t1+t2+t3) * time.Second)

	if now.Before(p1_endts) {
		rp = c_Round_Piece_Guess
		endts = p1_endts.Unix()
	} else if now.Before(p2_endts) {
		rp = c_Round_Piece_FightResult
		endts = p2_endts.Unix()
	} else {
		rp = c_Round_Piece_Award
		endts = p3_endts.Unix()
	}

	return
}

// 阶段更新竞猜结果
func round_guess_result(guess *guess_info_t, replay *msg.BattleReplay) {
	if guess == nil {
		return
	}

	if (guess.plrid1 == replay.Bi.T1.Player.Id) &&
		(guess.plrid2 == replay.Bi.T2.Player.Id) {
		guess.replay = replay
	}
}

// 阶段竞猜赔率
func calc_guess_ratio(guess *guess_info_t) (r1, r2 float32) {
	r1, r2 = float32(1), float32(1)
	if guess == nil {
		return
	}

	plr1 := warcup_plrs[guess.plrid1]
	plr2 := warcup_plrs[guess.plrid2]
	if plr1 == nil || plr2 == nil {
		return
	}

	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return
	}

	ap1, ap2 := plr1.T.Player.AtkPwr, plr2.T.Player.AtkPwr
	a := float32(ap1) / float32(ap2)
	if ap1 < ap2 {
		a = float32(ap2) / float32(ap1)
	}

	b := MinFloat32(a, conf.WarCupGuessOddsLimit)

	L := int32(len(guess.guess_plrs))
	if L > 0 {
		r1 += b * float32(L-guess.win1) / float32(L)
		r2 += b * float32(guess.win1) / float32(L)
	}

	r1 = math.MaxFloat32(r1, conf.WarCupGuessFloors)
	r2 = math.MaxFloat32(r2, conf.WarCupGuessFloors)

	return
}

// 阶段结算
func round_award(guess *guess_info_t) {
	// 竞猜
	if guess == nil ||
		guess.replay == nil ||
		guess.replay.Br == nil {
		return
	}

	plr1 := warcup_plrs[guess.plrid1]
	plr2 := warcup_plrs[guess.plrid2]
	if plr1 == nil || plr2 == nil {
		return
	}

	round, _, _ := round_piece()
	r1, r2 := calc_guess_ratio(guess)

	winner := guess.replay.Br.Winner
	for plrid, v := range guess.guess_plrs {
		var add int32
		if winner == v.winner {
			if winner == 1 {
				add = int32(float32(v.num) * r1)
			} else {
				add = int32(float32(v.num) * r2)
			}

			plr := find_player(plrid)
			if plr != nil {
				evtmgr.Fire(gconst.Evt_WarCupGuessWin, plr)
			}
		}

		plr_guess_result_update(plrid, round, add)
	}
}

// 玩家竞猜结果更新
func plr_guess_result_update(plrid string, round int32, add int32) {
	plr_data := warcup_plrguess[plrid]
	if plr_data == nil {
		return
	}

	plr_data.score += add

	for _, v := range plr_data.records {
		if v.stage == g_stage && v.round == round {
			v.add = add
			v.isFin = true
			break
		}
	}
}

// 玩家计算首次参与竞猜,可得免费积分
func plr_first_guess_free_score(plrid string) (freeScore int32) {
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return
	}

	for st := int32(c_Stage_Audition); st <= g_stage; st++ {
		switch st {
		case c_Stage_Audition:
			L := int32(len(conf.WarCupAFreePoints))
			for rd := int32(1); rd <= audition.round; rd++ {
				if L >= rd {
					freeScore += conf.WarCupAFreePoints[rd-1]
				}
			}

		case c_Stage_Top64:
			L := int32(len(conf.WarCupKFreePoints))
			for rd := int32(1); rd <= top64.round; rd++ {
				if L >= rd {
					freeScore += conf.WarCupKFreePoints[rd-1]
				}
			}

		case c_Stage_Top8:
			L := int32(len(conf.WarCupFFreePoints))
			for rd := int32(1); rd <= top8.round; rd++ {
				if L >= rd {
					freeScore += conf.WarCupFFreePoints[rd-1]
				}
			}
		default:
		}
	}

	return
}

// 玩家每轮开始免费给竞猜积分(至少竞猜过一轮才给)
func plr_round_guess_free_score() {
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return
	}

	var freeScore int32
	switch g_stage {
	case c_Stage_Audition:
		if len(conf.WarCupAFreePoints) >= int(audition.round) && (audition.round > 0) {
			freeScore = conf.WarCupAFreePoints[audition.round-1]
		}
	case c_Stage_Top64:
		if len(conf.WarCupAFreePoints) >= int(top64.round) && (top64.round > 0) {
			freeScore = conf.WarCupAFreePoints[top64.round-1]
		}
	case c_Stage_Top8:
		if len(conf.WarCupAFreePoints) >= int(top8.round) && (top8.round > 0) {
			freeScore = conf.WarCupAFreePoints[top8.round-1]
		}
	default:
	}

	for _, v := range warcup_plrguess {
		if len(v.records) == 0 {
			continue
		}

		v.score += freeScore
	}
}

// 赛季发奖
func end_award() {
	end_award_guess()
	end_award_rank()
}

// 竞猜奖励
func end_award_guess() {
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return
	}

	// 竞猜奖励
	mid := int32(conf.WarCupGuessMailId)
	conf_m := gamedata.ConfMail.Query(mid)
	if conf_m == nil {
		return
	}

	for plrid, v := range warcup_plrguess {
		plr := load_player(plrid)
		if v.score > 0 && plr != nil {
			m := mail.New(plr).SetKey(mid)
			m.AddAttachment(gconst.ArenaCoin, float64(float32(v.score)*conf.WarCupCoinRatio))
			m.AddDict("num", core.I32toa(int32(v.score)))
			m.Send()
		}
	}
}

// 排行奖励
func end_award_rank() {
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return
	}

	mid := int32(conf.WarCupRankMailId)
	conf_m := gamedata.ConfMail.Query(mid)
	if conf_m == nil {
		return
	}

	dup := make(map[string]bool)
	for i := int32(1); i <= 8; i++ {
		conf := gamedata.ConfWarCupRankReward.Query(i)
		if conf == nil {
			continue
		}

		switch conf.Id {
		case 1: //冠军
			if top8 != nil &&
				top8.rank_grps != nil &&
				len(top8.rank_grps.champion) > 0 {

				plr := find_player(top8.rank_grps.champion[0])
				if plr != nil && !dup[plr.GetId()] {
					m := mail.New(plr).SetKey(mid)
					for _, v := range conf.Reward {
						m.AddAttachment(v.Id, float64(v.N))
					}
					m.AddDict("cupWarRank", conf.WarCupRank)
					m.Send()
					dup[plr.GetId()] = true
				}
			}
		case 2: //亚军
			if top8 != nil &&
				top8.rank_grps != nil &&
				len(top8.rank_grps.round3) > 0 {
				for _, plrid := range top8.rank_grps.round3[0].tt {
					if dup[plrid] {
						continue
					}

					plr := find_player(plrid)
					if plr != nil {
						m := mail.New(plr).SetKey(mid)
						for _, v := range conf.Reward {
							m.AddAttachment(v.Id, float64(v.N))
						}
						m.AddDict("cupWarRank", conf.WarCupRank)
						m.Send()
						dup[plrid] = true
					}
				}
			}
		case 3: //四强
			if top8 != nil &&
				top8.rank_grps != nil {
				for _, v4 := range top8.rank_grps.round2 {
					for _, plrid := range v4.tt {
						if dup[plrid] {
							continue
						}

						plr := find_player(plrid)
						if plr != nil {
							m := mail.New(plr).SetKey(mid)
							for _, v := range conf.Reward {
								m.AddAttachment(v.Id, float64(v.N))
							}
							m.AddDict("cupWarRank", conf.WarCupRank)
							m.Send()
							dup[plrid] = true
						}
					}
				}
			}
		case 4: //八强
			if top8 != nil &&
				top8.rank_grps != nil &&
				len(top8.rank_grps.round1) > 0 {
				for _, v8 := range top8.rank_grps.round1 {
					for _, plrid := range v8.tt {
						if dup[plrid] {
							continue
						}

						plr := find_player(plrid)
						if plr != nil {
							m := mail.New(plr).SetKey(mid)
							for _, v := range conf.Reward {
								m.AddAttachment(v.Id, float64(v.N))
							}
							m.AddDict("cupWarRank", conf.WarCupRank)
							m.Send()
							dup[plrid] = true
						}
					}
				}
			}
		case 5: //十六强
			if top64 != nil &&
				top64.rank_grps != nil {
				for _, v16 := range top64.rank_grps.round3 {
					for _, plrid := range v16.tt {
						if dup[plrid] {
							continue
						}

						plr := find_player(plrid)
						if plr != nil {
							m := mail.New(plr).SetKey(mid)
							for _, v := range conf.Reward {
								m.AddAttachment(v.Id, float64(v.N))
							}
							m.AddDict("cupWarRank", conf.WarCupRank)
							m.Send()
							dup[plrid] = true
						}
					}
				}
			}
		case 6: //三十二强
			if top64 != nil &&
				top64.rank_grps != nil {
				for _, v32 := range top64.rank_grps.round2 {
					for _, plrid := range v32.tt {
						if dup[plrid] {
							continue
						}

						plr := find_player(plrid)
						if plr != nil {
							m := mail.New(plr).SetKey(mid)
							for _, v := range conf.Reward {
								m.AddAttachment(v.Id, float64(v.N))
							}
							m.AddDict("cupWarRank", conf.WarCupRank)
							m.Send()
							dup[plrid] = true
						}
					}
				}
			}
		case 7: //六十四强
			if top64 != nil &&
				top64.rank_grps != nil {
				for _, v64 := range top64.rank_grps.round1 {
					for _, plrid := range v64.tt {
						if dup[plrid] {
							continue
						}

						plr := find_player(plrid)
						if plr != nil {
							m := mail.New(plr).SetKey(mid)
							for _, v := range conf.Reward {
								m.AddAttachment(v.Id, float64(v.N))
							}
							m.AddDict("cupWarRank", conf.WarCupRank)
							m.Send()
							dup[plrid] = true
						}
					}
				}
			}
		case 8: //参与奖
			if audition != nil &&
				audition.rank_grps != nil {
				for _, v256 := range audition.rank_grps {
					if dup[v256.plrid] {
						continue
					}

					plr := find_player(v256.plrid)
					if plr != nil {
						m := mail.New(plr).SetKey(mid)
						for _, v := range conf.Reward {
							m.AddAttachment(v.Id, float64(v.N))
						}
						m.AddDict("cupWarRank", conf.WarCupRank)
						m.Send()
						dup[v256.plrid] = true
					}
				}
			}
		}
	}
}

// ============================================================================

func (self *vs_data_t) to_msg(isWhole bool) *msg.WarCupVsData {
	ret := &msg.WarCupVsData{
		VsSeq: self.vsseq,
	}

	if isWhole {
		plr1 := warcup_plrs[self.plrid1]
		plr2 := warcup_plrs[self.plrid2]
		if plr1 != nil && plr2 != nil {
			ret.Plr1 = plr1.T.Player
			ret.Plr2 = plr2.T.Player
		}

		ret.Stage = self.stage
		ret.Round = self.round
		ret.AddScore = self.addScore
	}

	if self.replay != nil && self.replay.Br != nil {
		ret.Winner = self.replay.Br.Winner
	}

	return ret
}

// ============================================================================
// req

// 是否开启
func IsOpen() bool {
	return !force_close && g_stage != c_Stage_Close
}

// 杯赛信息
func WarCup_ToMsg(plr IPlayer) *msg.WarCupData {
	ret := &msg.WarCupData{}
	if !IsOpen() {
		return ret
	}

	ret.Stage = g_stage
	ret.Ts2 = g_t2.Unix()
	ret.Round, ret.RoundPiece, ret.PieceEndTs = round_piece()

	for _, v := range warcup_chat {
		ret.Chat = append(ret.Chat, &msg.WarCupChatOne{
			Name:    v.name,
			Content: v.content,
		})
	}

	return ret
}

// 竞猜信息
func WarCupGuessInfo(plr IPlayer) *msg.WarCupGuessData {
	ret := &msg.WarCupGuessData{}

	var guess *guess_info_t
	switch g_stage {
	case c_Stage_Audition:
		guess = audition.guess

	case c_Stage_Top64:
		guess = top64.guess

	case c_Stage_Top8:
		guess = top8.guess

	}

	if guess != nil {
		_, rp, _ := round_piece()
		if guess.replay == nil || rp == c_Round_Piece_Guess {
			plr1 := warcup_plrs[guess.plrid1]
			plr2 := warcup_plrs[guess.plrid2]
			if plr1 != nil && plr2 != nil {
				ret.Replay = &msg.BattleReplay{
					Bi: &msg.BattleInput{
						T1: plr1.T,
						T2: plr2.T,
						Args: map[string]string{
							"Module":    "WAR_CUP",
							"RoundType": "3",
						},
					},
				}
			}
		} else {
			ret.Replay = guess.replay
		}

		n1, n2 := guess.win1, int32(len(guess.guess_plrs))
		if n2 > 0 {
			ret.GuessRatio = float32(n1) / float32(n2)
		} else {
			ret.GuessRatio = 0.5
		}

		ret.GuessWinRatio1, ret.GuessWinRatio2 = calc_guess_ratio(guess)

		plr_guess := warcup_plrguess[plr.GetId()]
		if plr_guess != nil {
			ret.GuessScore = plr_guess.score
			if len(plr_guess.records) > 0 {
				ret.GuessHas = true
			}
		}

		guess_plr := guess.guess_plrs[plr.GetId()]
		if guess_plr != nil {
			ret.GuessWin = guess_plr.winner
			ret.GuessNum = guess_plr.num
		}

	}

	return ret
}

// 竞猜
func WarCupGuess(plr IPlayer, winner int32, num int32) (int32, int32) {
	if !IsOpen() {
		return Err.Common_TimeNotUp, 0
	}

	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil || num <= 0 {
		return Err.Failed, 0
	}

	if winner != 1 && winner != 2 {
		return Err.Failed, 0
	}

	if num > conf.WarCupGuessLimit {
		num = conf.WarCupGuessLimit
	}

	var guess *guess_info_t
	switch g_stage {
	case c_Stage_Audition:
		guess = audition.guess
	case c_Stage_Top64:
		guess = top64.guess
	case c_Stage_Top8:
		guess = top8.guess
	default:
	}

	if guess == nil {
		return Err.WarCup_NotInGuess, 0
	}

	round, rp, _ := round_piece()
	if rp != c_Round_Piece_Guess {
		return Err.WarCup_NotInGuess, 0
	}

	if guess.guess_plrs[plr.GetId()] != nil {
		return Err.WarCup_GuessBefore, 0
	}

	// add score
	var freeScore int32
	plr_guess := warcup_plrguess[plr.GetId()]
	if plr_guess == nil {
		freeScore = plr_first_guess_free_score(plr.GetId())

		plr_guess = &plr_guess_data_t{}
		warcup_plrguess[plr.GetId()] = plr_guess
	}

	plr_score := plr_guess.score + freeScore
	if num > plr_score {
		num = plr_score
	}
	plr_guess.score = plr_score - num

	// add record
	rec := &guess_record_t{
		stage: g_stage,
		round: round,
		num:   num,
	}

	if winner == 1 {
		guess.win1++

		plr1 := warcup_plrs[guess.plrid1]
		if plr1 != nil {
			rec.name = plr1.T.Player.Name
		}
	} else {
		plr2 := warcup_plrs[guess.plrid2]
		if plr2 != nil {
			rec.name = plr2.T.Player.Name
		}
	}
	plr_guess.records = append(plr_guess.records, rec)

	guess.guess_plrs[plr.GetId()] = &guess_one_t{
		winner: winner,
		num:    num,
	}

	r1, r2 := calc_guess_ratio(guess)

	// res
	utils.BroadcastPlayers(&msg.GS_WarCupGuessRatio{
		GuessRatio:     float32(guess.win1) / float32(len(guess.guess_plrs)),
		GuessWinRatio1: r1,
		GuessWinRatio2: r2,
	})

	//evt
	evtmgr.Fire(gconst.Evt_WarCupGuess, plr)

	return Err.OK, plr_guess.score
}

// 玩家竞猜记录
func WarCupGuessRecords(plr IPlayer) (ret []*msg.WarCupPlrGuessOne) {
	plr_guess := warcup_plrguess[plr.GetId()]
	if plr_guess == nil {
		return
	}

	for _, v := range plr_guess.records {
		ret = append(ret, &msg.WarCupPlrGuessOne{
			Stage: v.stage,
			Round: v.round,
			Name:  v.name,
			Num:   v.num,
			Add:   v.add,
			IsFin: v.isFin,
		})
	}

	return
}

// 自己的对战信息
func WarCupSelfVsInfo(plr IPlayer) (ret []*msg.WarCupVsData, curReplay *msg.BattleReplay) {
	if !IsOpen() {
		return
	}

	plr_data := warcup_plrs[plr.GetId()]
	if plr_data == nil {
		return
	}

	for _, v := range plr_data.VsSeqs {
		vsdata := warcup_vsdata[v]
		if vsdata == nil {
			continue
		}

		round, rp, _ := round_piece()
		m := vsdata.to_msg(true)

		// 当前轮轮数竞猜结束前结果不返回
		if vsdata.stage == g_stage &&
			vsdata.round == round {

			plr1 := warcup_plrs[vsdata.plrid1]
			plr2 := warcup_plrs[vsdata.plrid2]
			if plr1 != nil && plr2 != nil {
				curReplay = &msg.BattleReplay{
					Bi: &msg.BattleInput{
						T1: plr1.T,
						T2: plr2.T,
						Args: map[string]string{
							"Module":    "WAR_CUP",
							"RoundType": "3",
						}},
				}
			}

			// 竞猜阶段其他信息不返
			if rp == c_Round_Piece_Guess {
				continue
			} else {
				if rp == c_Round_Piece_FightResult {
					m.Winner = 0
				}

				if vsdata.replay != nil {
					curReplay.Br = vsdata.replay.Br
				}
			}
		}

		ret = append(ret, m)
	}

	return
}

// 获取64强淘汰赛信息
func WarCupTop64Info(grp int32) (ret []*msg.WarCupVsData) {
	if !IsOpen() {
		return
	}

	if grp < 1 || grp > 8 {
		grp = 1
	}

	if top64 == nil || top64.rank_grps == nil {
		return
	}

	round, rp, _ := round_piece()

	if len(top64.rank_grps.round1) >= 32 {
		for i := 0; i < 4; i++ {
			v := top64.rank_grps.round1[int(grp-1)*4+i]
			if v == nil {
				continue
			}

			vsdata := warcup_vsdata[v.vsseq]
			if vsdata == nil {
				continue
			}

			m := vsdata.to_msg(true)
			// 当前轮轮数竞猜结束前结果不返回
			if vsdata.stage == g_stage &&
				vsdata.round == round {

				if rp == c_Round_Piece_Guess {
					m.VsSeq = 0
					m.Winner = 0
				} else if rp == c_Round_Piece_FightResult {
					m.Winner = 0
				}
			}
			ret = append(ret, m)
		}
	}

	if len(top64.rank_grps.round2) >= 16 {
		for i := 0; i < 2; i++ {
			v := top64.rank_grps.round2[int(grp-1)*2+i]
			if v == nil {
				continue
			}

			vsdata := warcup_vsdata[v.vsseq]
			if vsdata == nil {
				continue
			}

			m := vsdata.to_msg(false)
			// 当前轮轮数竞猜结束前结果不返回
			if vsdata.stage >= g_stage &&
				vsdata.round == round {

				if rp == c_Round_Piece_Guess {
					m.VsSeq = 0
					m.Winner = 0
				} else if rp == c_Round_Piece_FightResult {
					m.Winner = 0
				}
			}
			ret = append(ret, m)
		}
	}

	if len(top64.rank_grps.round3) >= 8 {
		v := top64.rank_grps.round3[int(grp-1)]
		if v != nil {
			vsdata := warcup_vsdata[v.vsseq]
			if vsdata != nil {
				m := vsdata.to_msg(false)
				// 当前轮轮数竞猜结束前结果不返回
				if vsdata.stage >= g_stage &&
					vsdata.round == round {

					if rp == c_Round_Piece_Guess {
						m.VsSeq = 0
						m.Winner = 0
					} else if rp == c_Round_Piece_FightResult {
						m.Winner = 0
					}
				}
				ret = append(ret, m)
			}
		}
	}

	return
}

// 获取8强淘汰赛信息
func WarCupTop8Info() (ret []*msg.WarCupVsData) {
	if !IsOpen() {
		return
	}

	if top8 == nil || top8.rank_grps == nil {
		return
	}

	round, rp, _ := round_piece()

	for _, v := range top8.rank_grps.round1 {
		if v == nil {
			continue
		}

		vsdata := warcup_vsdata[v.vsseq]
		if vsdata == nil {
			continue
		}

		m := vsdata.to_msg(true)
		// 当前轮轮数竞猜结束前结果不返回
		if vsdata.stage == g_stage &&
			vsdata.round == round {

			if rp == c_Round_Piece_Guess {
				m.VsSeq = 0
				m.Winner = 0
			} else if rp == c_Round_Piece_FightResult {
				m.Winner = 0
			}
		}
		ret = append(ret, m)
	}

	for _, v := range top8.rank_grps.round2 {
		if v == nil {
			continue
		}

		vsdata := warcup_vsdata[v.vsseq]
		if vsdata == nil {
			continue
		}

		m := vsdata.to_msg(false)
		// 当前轮轮数竞猜结束前结果不返回
		if vsdata.stage == g_stage &&
			vsdata.round == round {

			if rp == c_Round_Piece_Guess {
				m.VsSeq = 0
				m.Winner = 0
			} else if rp == c_Round_Piece_FightResult {
				m.Winner = 0
			}
		}
		ret = append(ret, m)
	}

	for _, v := range top8.rank_grps.round3 {
		if v == nil {
			continue
		}

		vsdata := warcup_vsdata[v.vsseq]
		if vsdata == nil {
			continue
		}

		m := vsdata.to_msg(false)
		// 当前轮轮数竞猜结束前结果不返回
		if vsdata.stage == g_stage &&
			vsdata.round == round {

			if rp == c_Round_Piece_Guess {
				m.VsSeq = 0
				m.Winner = 0
			} else if rp == c_Round_Piece_FightResult {
				m.Winner = 0
			}
		}
		ret = append(ret, m)
	}

	return
}

// 获取冠军
func WarCupTop1Info() (info *msg.PlayerSimpleInfo, heroid, star, skin int32) {
	if !IsOpen() {
		return
	}

	if top8 == nil || top8.rank_grps == nil {
		return
	}

	if len(top8.rank_grps.champion) == 0 {
		return
	}

	plrid := top8.rank_grps.champion[0]

	plr := warcup_plrs[plrid]
	if plr == nil || plr.T == nil {
		return
	}

	info = plr.T.Player

	max := int32(0)
	for _, v := range plr.T.Fighters {
		if v == nil {
			continue
		}

		if v.Id > 0 && heroid == 0 { // 默认第一个
			heroid = v.Id
			star = v.Star
			skin = v.Skin
		}

		if max < v.AtkPwr { // 再看战力
			max = v.AtkPwr
			heroid = v.Id
			star = v.Star
			skin = v.Skin
		}
	}

	return
}

// 拉取战报
func WarCupGetReplay(vsseq int32) *msg.BattleReplay {
	vsdata := warcup_vsdata[vsseq]
	if vsdata == nil {
		return nil
	}

	return vsdata.replay
}

// 聊天
func WarCupChat(plr IPlayer, content string) {
	// filter
	content = wordsfilter.Filter(content)

	one := &chat_ont_t{
		name:    plr.GetName(),
		content: content,
	}
	warcup_chat = append(warcup_chat, one)

	if L := len(warcup_chat); L > 20 {
		warcup_chat = warcup_chat[L-20:]
	}

	utils.BroadcastPlayers(&msg.GS_WarCupChat{
		Name:    plr.GetName(),
		Content: content,
	})

	//evt
	evtmgr.Fire(gconst.Evt_WarCupChat, plr)
}

// 获取海选排行榜
func WarCupAuditonRank(top, N int32) (ret []*msg.RankRow) {
	if !IsOpen() {
		return
	}

	if audition == nil {
		return
	}

	if N <= 0 || N > c_Audition_MaxPlrs {
		N = 50
	}

	to := top + N
	if top <= 0 || top > c_Audition_MaxPlrs {
		top = 1
	}

	L := int32(len(audition.rank_grps))
	for i := top; i < to && i <= L; i++ {
		v := audition.rank_grps[i-1]

		plr := FindWarCupPlayer(v.plrid)
		if plr == nil {
			continue
		}

		one := &msg.RankRow{
			Score: float64(v.sc),
			Info: &msg.RankRowInfo{
				Plr: plr.ToMsg_SimpleInfo(v.ap),
			},
		}

		ret = append(ret, one)
	}

	return
}
