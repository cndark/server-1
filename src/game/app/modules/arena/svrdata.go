package arena

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/sched/loop"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/mdata"
	"fw/src/game/app/modules/robot"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"math/rand"
	"sort"
	"time"
)

// ============================================================================

const (
	C_Enemy_Cnt  = 2    // 敌方人数最后会多加一个
	C_Award_Time = 22   // 发奖时间
	C_Row_Max    = 100  // 排行榜最多人数
	C_BeginScore = 1000 // 起始分数

)

var (
	rand_d   = rand.New(rand.NewSource(time.Now().Unix()))
	ArenaMgr *arena_mgr_t
)

// ============================================================================

type arena_mgr_t struct {
	VerTs time.Time // 当前版本时间(t0)
	Ts1   time.Time // 当前阶段开始时间(t1)
	Ts2   time.Time // 当前阶段结束时间(t2)
	stage string    // 当前stage

	Rows         []*arena_row_t // 玩家竞技场分数
	DailyAwardTs time.Time      // 每天奖励发放时间
	IsWeekAward  bool           // 赛季奖励是否发放
}

// 积分列表
type arena_row_t struct {
	PlrId string
	Score int32
}

// ============================================================================

func new_data() interface{} {
	return &arena_mgr_t{}
}

func data_loaded() {
	ArenaMgr = mdata.Get(NAME).(*arena_mgr_t)
	ArenaMgr.init()
}

// ============================================================================

func (self *arena_mgr_t) init() {
	self.start_daily_award_timer()
	next_ts := self.daily_award_next_time()

	if !core.IsSameDay(time.Now(), next_ts) {
		loop.Push(func() {
			self.daily_award()
		})
	}
}

// 赛季调度
func (self *arena_mgr_t) reset() {
	self.IsWeekAward = false
	self.Rows = []*arena_row_t{}
}

func on_stage(stg string, t0, t1, t2 time.Time, f func(bool)) {
	if ArenaMgr == nil {
		ArenaMgr = mdata.Get(NAME).(*arena_mgr_t)
		ArenaMgr.init()
	}

	if ArenaMgr.VerTs.IsZero() || !ArenaMgr.VerTs.Equal(t0) {
		ArenaMgr.VerTs = t0
		ArenaMgr.reset()
	}

	ArenaMgr.Ts1 = t1
	ArenaMgr.Ts2 = t2
	ArenaMgr.stage = stg

	f(true)
}

func stage_start() {
	ArenaMgr.stage = "start"

	utils.BroadcastPlayers(&msg.GS_ArenaStageUpdate{
		Stage: ArenaMgr.stage,
		Ts1:   ArenaMgr.Ts1.Unix(),
		Ts2:   ArenaMgr.Ts2.Unix(),
	})
}

func stage_end() {
	ArenaMgr.stage = "end"

	ArenaMgr.week_award()

	utils.BroadcastPlayers(&msg.GS_ArenaStageUpdate{
		Stage: ArenaMgr.stage,
		Ts1:   ArenaMgr.Ts1.Unix(),
		Ts2:   ArenaMgr.Ts2.Unix(),
	})
}

func stage_close() {
	ArenaMgr.stage = "close"

	utils.BroadcastPlayers(&msg.GS_ArenaStageUpdate{
		Stage: ArenaMgr.stage,
		Ts1:   ArenaMgr.Ts1.Unix(),
		Ts2:   ArenaMgr.Ts2.Unix(),
	})
}

// ============================================================================
// timer

// 开启每日奖励调度
func (self *arena_mgr_t) start_daily_award_timer() {
	loop.SetTimeout(self.daily_award_next_time(), func() {
		self.daily_award()
		self.start_daily_award_timer()
	})
}

func (self *arena_mgr_t) daily_award_next_time() time.Time {
	now := time.Now()
	zero := core.StartOfDay(now)
	key := zero.Add(time.Duration(C_Award_Time) * time.Hour)

	if now.After(key) {
		key = key.Add(time.Duration(24) * time.Hour)
	}

	return key
}

// 日常奖励是否发放
func (self *arena_mgr_t) is_daily_award() bool {
	now := time.Now()
	zero := core.StartOfDay(now)
	t_key := zero.Add(time.Duration(C_Award_Time) * time.Hour)

	if now.After(t_key) {
		return self.DailyAwardTs.After(t_key)
	} else {
		y_key := t_key.Add(time.Duration(-24) * time.Hour)
		return self.DailyAwardTs.After(y_key)
	}
}

// 当日奖励
func (self *arena_mgr_t) daily_award() {
	if self.is_daily_award() {
		return
	}

	self.award_mail(true)

	self.DailyAwardTs = time.Now()
	log.Warning("arena daily award end ")
}

// 发奖励邮件
func (self *arena_mgr_t) award_mail(isDaily bool) {
	for idx, row := range self.Rows {
		idx := idx + 1
		for _, conf := range gamedata.ConfArena.Items() {
			if len(conf.RankRange) > 0 && conf.RankRange[0].Up <= int32(idx) &&
				conf.RankRange[0].Down >= int32(idx) {
				mid := conf.WeekMail
				if isDaily {
					mid = conf.DailyMail
				}

				conf_m := gamedata.ConfMail.Query(mid)
				if conf_m == nil {
					continue
				}

				plr := find_player(row.PlrId)
				if plr == nil {
					continue
				}

				m := mail.New(plr).SetKey(mid)
				for _, v := range conf_m.MailItem {
					m.AddAttachment(v.Id, float64(v.N))
				}
				m.Send()
				break
			}
		}
	}
}

// 赛季奖励
func (self *arena_mgr_t) week_award() {
	if self.IsWeekAward {
		return
	}

	self.award_mail(false)

	self.IsWeekAward = true
}

// ============================================================================

// 战斗结束后排序
func (self *arena_mgr_t) end_battle_sort_rank(plrid, eplrid string) {
	// old rank idx
	idxs := []int32{}
	idxs = append(idxs, self.find_top_idx(plrid))
	if !IsRobot(eplrid) {
		idxs = append(idxs, self.find_top_idx(eplrid))
	}

	// rank
	self.sort_score()

	// new rank idx
	idxs = append(idxs, self.find_top_idx(plrid))
	if !IsRobot(eplrid) {
		idxs = append(idxs, self.find_top_idx(eplrid))
	}

	min, max := int32(0), int32(0)
	for _, v := range idxs {
		if v > 0 {
			if min == 0 || min > v {
				min = v
			}

			if max == 0 || max < v {
				max = v
			}
		}
	}

	self.evt_fire(min, max)
}

// 积分排序
func (self *arena_mgr_t) sort_score() {
	sort.Slice(self.Rows, func(i, j int) bool {
		if self.Rows[i].Score > self.Rows[j].Score {
			return true
		} else if self.Rows[i].Score == self.Rows[j].Score {
			plr_i := find_player(self.Rows[i].PlrId)
			if plr_i == nil {
				return false
			}

			plr_j := find_player(self.Rows[i].PlrId)
			if plr_j == nil {
				return true
			}

			return plr_i.GetAtkPwr() >= plr_j.GetAtkPwr()
		}
		return false
	})

	return
}

// 获取100内玩家名次
func (self *arena_mgr_t) find_top_idx(plrid string) int32 {
	for i, v := range self.Rows {
		if i >= C_Row_Max {
			return 0
		}

		if v.PlrId == plrid {
			return int32(i) + 1
		}
	}

	return 0
}

// 排名事件
func (self *arena_mgr_t) evt_fire(min, max int32) {
	L := int32(len(self.Rows))
	for i := min; i <= max; i++ {
		if i > 0 && i <= L {
			plr := find_player(self.Rows[i-1].PlrId)
			if plr == nil {
				continue
			}

			evtmgr.Fire(gconst.Evt_ArenaRank, plr, -i)
		}
	}
}

// ============================================================================

// 获取玩家积分
func (self *arena_mgr_t) GetScore(plrid string) int32 {
	for _, v := range self.Rows {
		if plrid == v.PlrId {
			return v.Score
		}
	}

	return 0
}

// 加积分
func (self *arena_mgr_t) AddScore(plrid string, v int32) {
	for i, row := range self.Rows {
		if plrid == row.PlrId {
			self.Rows[i].Score += v

			if self.Rows[i].Score <= 0 {
				self.Rows[i].Score = 1
			}

			return
		}
	}

	self.Rows = append(self.Rows, &arena_row_t{
		PlrId: plrid,
		Score: C_BeginScore + v,
	})
}

// 匹配对手
func (self *arena_mgr_t) PickEnemies(plr IPlayer) (ret []*msg.ArenaEnemy) {
	ids := []string{}
	L := len(self.Rows)

	// rank pick player
	if L > 0 {
		dup := make(map[string]bool)

		selfIdx := self.find_top_idx(plr.GetId())
		selfScore := plr.GetArenaScore()
		if selfScore == 0 {
			selfScore = C_BeginScore
		}

		conf := gamedata.ConfGlobalPublic.Query(1)
		if conf != nil {
			for i, v := range conf.ArenaMatching {
				for _, row := range self.Rows {
					if float32(selfScore)*v.A <= float32(row.Score) &&
						float32(selfScore)*v.B >= float32(row.Score) {

						if row.PlrId != plr.GetId() && !dup[row.PlrId] {
							eplr := load_player(row.PlrId)
							if eplr != nil {
								ids = append(ids, row.PlrId)
								dup[row.PlrId] = true
								ret = append(ret, &msg.ArenaEnemy{
									Score: eplr.GetArenaScore(),
									Plr:   eplr.ToMsg_SimpleInfo(),
								})
								goto LOOP_END
							}
						}
					}
				}

				// 没找到就在排行榜找,自己前后一名
				if selfIdx > 0 {
					plrid := ""
					if i == 0 { // 强的
						if selfIdx > 1 {
							plrid = self.Rows[selfIdx-2].PlrId
						}
					} else if i == 1 { // 低一名
						if selfIdx < int32(len(self.Rows)) {
							plrid = self.Rows[selfIdx].PlrId
						}
					}

					if plrid != "" && !dup[plrid] {
						eplr := load_player(plrid)
						if eplr != nil {
							ids = append(ids, plrid)
							dup[plrid] = true
							ret = append(ret, &msg.ArenaEnemy{
								Score: eplr.GetArenaScore(),
								Plr:   eplr.ToMsg_SimpleInfo(),
							})
							goto LOOP_END
						}
					}
				}

			LOOP_END:
			}
		}
	}

	// 人数不够凑机器人
	rdN := C_Enemy_Cnt - len(ids)
	if rdN > 0 {
		bots := robot.RobotMgr.ArenaEnemies(plr.GetLevel())
		for _, id := range bots {
			if len(ids) >= C_Enemy_Cnt {
				break
			}

			bot := FindArenaPlayer(id)
			if bot != nil {
				ret = append(ret, &msg.ArenaEnemy{
					Score: bot.GetArenaScore(),
					Plr:   bot.ToMsg_SimpleInfo(),
				})

				ids = append(ids, id)
			}
		}
	}

	// 最后一个是和自己同等级的机器人
	bots := robot.RobotMgr.GetByLevel(plr.GetLevel(), 1)
	if len(bots) > 0 {
		ret = append(ret, &msg.ArenaEnemy{
			Score: bots[0].GetArenaScore(),
			Plr:   bots[0].ToMsg_SimpleInfo(),
		})

		ids = append(ids, bots[0].Id)
	}

	plr.GetArena().SetEnemies(ids)

	return
}

// 获取排行
func (self *arena_mgr_t) ToMsg_Rank(plr IPlayer, top int32) (selfRank int32, ret []*msg.RankRow) {
	if top >= C_Row_Max || top <= 0 {
		top = C_Row_Max
	}

	for i, row := range self.Rows {
		if row.PlrId == plr.GetId() {
			selfRank = int32(i) + 1
		}

		if i >= C_Row_Max {
			continue
		}

		if i <= int(top) {
			rplr := load_player(row.PlrId)
			if rplr != nil {
				ret = append(ret, &msg.RankRow{
					Score: float64(row.Score),
					Info: &msg.RankRowInfo{
						Plr: rplr.ToMsg_SimpleInfo(),
					},
				})
			}
		}
	}

	return
}
