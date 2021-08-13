package maze

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/math"
	"fw/src/core/sched/loop"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/ranksvc"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
	"math/rand"
	"strings"
	"time"
)

// ============================================================================

const (
	C_Maze_Entrance = 1  // 入口
	C_Maze_Exit     = 2  // 出口
	C_Maze_Item     = 10 // 道具
	C_Maze_Key      = 11 // 钥匙
	C_Maze_Power    = 20 // 体力
	C_Maze_Box      = 30 // 宝箱
	C_Maze_Mine     = 40 // 矿
	C_Maze_Master   = 50 // 小怪
	C_Maze_Boss     = 51 // Boss
	C_Maze_Trade    = 60 // 商人
	C_Maze_Trap     = 70 // 陷阱
)

const (
	C_Task_Daily = 1
	C_Task_Achv  = 2
)

// ============================================================================

var rand_d = rand.New(rand.NewSource(time.Now().Unix()))

// ============================================================================

var actObj = &act_t{}

// ============================================================================

type act_t struct {
	act.BaseAct

	rank_tid *core.Timer
}

type data_svr_t struct {
	TotalScore map[string]int32
}

type data_plr_t struct {
	Lv      int32 // 层数
	Clicked map[int32]bool
	Points  map[int32]int32

	Attain map[int32]*attain_obj_t // 统计进度
	Taken  []int32                 // 任务领取

	Box     map[int32]int32     // 宝箱 pos:seq
	Mine    map[int32]int32     // 矿
	Item    map[int32]int32     // 道具
	Battle  map[int32]*battle_t // 战斗
	BuffIds map[int32]int32     // Buff加成

	TradeLimit map[int32]*trade_t // 商人限购

}

type trade_t struct {
	Shop map[int32]int32
}

type battle_t struct {
	Lv  int32
	Seq int32 // 配置表id
	Tp  int32 // 0: 小怪、1: Boss
}

type attain_obj_t struct {
	Id  int32
	Val float64
}

type clone_info_t struct {
	info  *ranksvc.RankRowInfo
	score int32
}

// ============================================================================

func init() {
	act.RegisterAct(gconst.ActName_Maze, actObj)

	evtmgr.On(gconst.Evt_PlrResetDaily, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		plr_data := actObj.GetPlrData(plr)
		if plr_data == nil {
			return
		}

		// 每日任务重置
		for _, v := range gamedata.ConfActMazeTask.Items() {
			if v.Type == C_Task_Daily {
				for i := 0; i < len(plr_data.Taken); i++ {
					if plr_data.Taken[i] == v.Id {
						plr_data.Taken = append(plr_data.Taken[:i], plr_data.Taken[i+1:]...)
						for _, v := range v.Attain {
							plr_data.Attain[v.AttainId] = &attain_obj_t{}
						}
						i--
					}
				}
			}
		}

	})
}

// ============================================================================

func (self *act_t) NewSvrData() interface{} {
	return &data_svr_t{
		TotalScore: map[string]int32{},
	}
}

func (self *act_t) NewPlrData() interface{} {
	return &data_plr_t{
		Attain:     map[int32]*attain_obj_t{},
		Points:     map[int32]int32{},
		BuffIds:    map[int32]int32{},
		Box:        map[int32]int32{},
		Mine:       map[int32]int32{},
		Item:       map[int32]int32{},
		Battle:     map[int32]*battle_t{},
		TradeLimit: map[int32]*trade_t{},
		Clicked:    map[int32]bool{},
	}
}

func (self *act_t) ResetSvrData() {

}

func (self *act_t) ResetPlrData(iplr interface{}) {
	plr := iplr.(IPlayer)

	// 初始counter
	cnt_init := int64(0)
	conf := gamedata.ConfActivity.Query(gconst.ActName_Maze)
	if conf != nil {
		cnt_init = int64(conf.InitialPower)
	}

	cop := plr.GetCounter().NewOp(gconst.ObjFrom_ActMaze)
	cop.SetCounter(gconst.Cnt_MazeCnt, plr.GetCounter().GetMax(gconst.Cnt_MazeCnt)-cnt_init)
	cop.SetCounter(gconst.Cnt_MazeLife, 0)

	// 清空临时物品
	conf_g := gamedata.ConfActivityPublic.Query(1)
	if conf_g != nil {
		for _, v := range conf_g.MazeItemClear {
			n := plr.GetBag().GetItem(v)
			if n > 0 {
				cop.Dec(v, n)
			}
		}
	}

	cop.Apply()
}

func (self *act_t) GetSvrData() *data_svr_t {
	return self.GetActRawData().(*data_svr_t)
}

func (self *act_t) GetPlrData(plr IPlayer) *data_plr_t {
	return plr.GetActRawData(self.GetName()).(*data_plr_t)
}

// ============================================================================

func (self *act_t) Started() bool {
	return self.GetStage() == "start"
}

func (self *act_t) Ended() bool {
	return self.GetStage() == "end"
}

func (self *act_t) Closed() bool {
	return self.GetStage() == "close"
}

// ============================================================================

func (self *act_t) OnInit() {

}

func (self *act_t) OnQuit() {
}

func (self *act_t) OnStage() {
	if self.Started() {

		self.stop_timer_push_rank()
		self.push_rank()
		self.start_timer_push_rank()
	}

	if self.Ended() {
		self.stop_timer_push_rank()
		loop.SetTimeout(time.Now().Add(time.Duration(2)*time.Minute), func() {
			self.push_rank()
		})
		loop.SetTimeout(time.Now().Add(time.Duration(5)*time.Minute), func() {
			self.check_box_reward_mail()
		})
	}
}

func (self *act_t) check_box_reward_mail() {
	ranksvc.Get(ranksvc.RankType_Local, config.CurGame.Id, gconst.RankId_ActMaze_Score, func(rows []*ranksvc.RankRow) {

		conf_g := gamedata.ConfActivityPublic.Query(1)
		if conf_g == nil {
			return
		}

		for idx, v := range rows {
			plr := load_player(v.Info.Id())

			svr_data := self.GetSvrData()

			m := mail.New(plr).SetKey(conf_g.MazeMail)

			isAward := false
			zero_item := make(map[int32]int64)

			for _, item := range gamedata.ConfActMazeRankM.QueryItems(actObj.GetConfGrp()) {
				if item.Rank[0].Low == 0 {
					for _, v := range item.Reward {
						zero_item[v.Id] = v.N
					}
				}

				if int32(idx+1) >= item.Rank[0].Low && int32(idx+1) <= item.Rank[0].High {
					for _, v := range item.Reward {
						m.AddAttachment(v.Id, float64(v.N))
					}
					isAward = true
					break
				}
			}

			if !isAward {
				for id, n := range zero_item {
					m.AddAttachment(id, float64(n))
				}
			}

			m.AddDictInt32("score", svr_data.TotalScore[plr.GetId()])
			m.AddDictInt32("rank", int32(idx+1))
			m.Send()
		}
	})
}

func (self *act_t) start_timer_push_rank() {
	self.rank_tid = loop.SetTimeout(time.Now().Add(time.Duration(5)*time.Minute), func() {
		if self.Started() {
			self.push_rank()
		}
		self.start_timer_push_rank()
	})
}

func (self *act_t) stop_timer_push_rank() {
	if self.rank_tid != nil {
		loop.CancelTimer(self.rank_tid)
		self.rank_tid = nil
	}
}

func (self *act_t) push_rank() {

	svr_data := self.GetSvrData()
	if svr_data == nil {
		return
	}

	cp := make([]*clone_info_t, 0, len(svr_data.TotalScore))
	for k, v := range svr_data.TotalScore {
		plr := find_player(k)

		info := &ranksvc.RankRowInfo{
			Plr: plr.ToMsg_SimpleInfo(),
		}

		cp = append(cp, &clone_info_t{
			info:  info,
			score: v,
		})
	}

	core.Go(func() {
		raw := &ranksvc.RankRaw{
			Id:        gconst.RankId_ActMaze_Score,
			SortLevel: ranksvc.SortLevel_Local,
			A:         cp,
			ScoreFunc: func(i int) float64 { return float64(cp[i].score) },
			InfoFunc:  func(i int) *ranksvc.RankRowInfo { return cp[i].info },
		}
		ranksvc.Push(raw)
	})
}

func (self *act_t) ToMsg(plr IPlayer) *msg.ActMazeData {
	ret := &msg.ActMazeData{}

	plr_data := self.GetPlrData(plr)
	if plr_data == nil {
		return nil
	}

	svr_data := actObj.GetSvrData()
	if svr_data == nil {
		return nil
	}

	if len(plr_data.Points) == 0 && plr_data.Lv == 0 {
		reset_maze(plr_data)
		gen_maze(plr_data, 16)
	}

	for _, v := range plr_data.Attain {
		ret.Attain = append(ret.Attain, &msg.ActAttainObj{
			OId: v.Id,
			Val: v.Val,
		})
	}

	ret.Lv = plr_data.Lv
	ret.Clicked = plr_data.Clicked
	ret.Score = svr_data.TotalScore[plr.GetId()]
	ret.Taken = plr_data.Taken
	ret.Points = plr_data.Points
	ret.Box = plr_data.Box
	ret.Mine = plr_data.Mine
	ret.Item = plr_data.Item
	ret.Lv = plr_data.Lv + 1
	ret.TradeLimit = make(map[int32]*msg.TradeData)
	ret.Battle = make(map[int32]*msg.BattleData)

	for idx, v := range plr_data.Battle {
		max_lv := int32(len(gamedata.ConfMonster.Items()))
		if v.Lv > max_lv {
			plr_data.Battle[idx].Lv = max_lv
		}
	}
	for pos, v := range plr_data.Battle {
		ret.Battle[pos] = &msg.BattleData{
			Lv:  v.Lv,
			Seq: v.Seq,
			Tp:  v.Tp,
		}
	}

	for pos, v := range plr_data.TradeLimit {
		ret.TradeLimit[pos] = &msg.TradeData{
			Data: v.Shop,
		}
	}

	return ret
}

// ============================================================================

func (self *data_plr_t) get_attain_obj(oid int32) *attain_obj_t {
	obj := self.Attain[oid]
	if obj == nil {
		self.Attain[oid] = &attain_obj_t{Id: oid}
		return self.Attain[oid]
	}
	return obj
}

// ============================================================================
// implements ICondObj interface

func (self *attain_obj_t) GetVal() float64 {
	return self.Val
}

func (self *attain_obj_t) SetVal(v float64) {
	self.Val = v
}

func (self *attain_obj_t) AddVal(v float64) {
	self.Val += v
}

func (self *attain_obj_t) Done(body interface{}, confid int32, isChange bool) {
	if !actObj.Started() || !isChange {
		return
	}

	body.(IPlayer).SendMsg(&msg.GS_ActMazeObjValueChanged{
		OId: self.Id,
		Val: self.Val,
	})
}

// ============================================================================

func seq_to_lv(seq int32) int32 {
	return seq%(actObj.GetConfGrp()*1000) - 1
}

func lv_to_seq(lv int32) int32 {
	return (actObj.GetConfGrp() * 1000) + lv + 1
}

func reset_maze(plr_data *data_plr_t) {

	plr_data.Clicked = map[int32]bool{}
	plr_data.Battle = map[int32]*battle_t{}
	plr_data.Box = map[int32]int32{}
	plr_data.Item = map[int32]int32{}
	plr_data.Mine = map[int32]int32{}
	plr_data.Points = map[int32]int32{}
	plr_data.TradeLimit = map[int32]*trade_t{}

	return
}

// 迷宫数据生成
func gen_maze(plr_data *data_plr_t, enterPos int32) {
	conf := gamedata.ConfActMazeLv.Query(lv_to_seq(plr_data.Lv))
	if conf == nil || len(conf.RandEventGrp) == 0 {
		return
	}

	rand_grp := conf.RandEventGrp[rand_d.Intn(len(conf.RandEventGrp))]

	conf_e := gamedata.ConfActMazeEventGrp.Query(rand_grp)
	if conf_e == nil {
		return
	}

	plr_data.Points[enterPos] = C_Maze_Entrance

	plr_data.Points[pick_exit(plr_data)] = C_Maze_Exit

	for _, v := range conf_e.EventNum {

		if v.EventId == C_Maze_Entrance || v.EventId == C_Maze_Exit {
			continue
		}

		for i := int32(0); i < v.N; i++ {
			pos := pick_pos(plr_data)

			plr_data.Points[pos] = v.EventId
		}
	}

	return
}

func pick_exit(plr_data *data_plr_t) int32 {
	conf := gamedata.ConfActMazeLv.Query(lv_to_seq(plr_data.Lv))

	exit_ids := []int32{}

	for i := 0; i < len(conf.ExitPos); i++ {
		pos := conf.ExitPos[i] - 1
		if plr_data.Points[pos] == 0 {
			exit_ids = append(exit_ids, pos)
		}
	}

	if len(exit_ids) == 0 {
		return 0
	}

	idx := rand_d.Intn(len(exit_ids))

	return exit_ids[idx]
}

func pick_pos(plr_data *data_plr_t) int32 {

	ids := []int32{}

	for i := int32(0); i < 25; i++ {
		if plr_data.Points[i] == 0 {
			ids = append(ids, i)
		}
	}

	if len(ids) == 0 {
		return 0
	}

	return ids[rand_d.Intn(len(ids))]
}

func pick_item(plr_data *data_plr_t) int32 {
	seq := int32(0)
	cnt := int32(0)
	slt := map[int32]int32{}
	for _, v := range gamedata.ConfActMazeEventItem.Items() {
		if cnt == 0 {
			seq = v.Seq
			cnt++
		}
		if v.ShowLvCond <= (plr_data.Lv + 1) {
			slt[v.Seq] = v.Weight
		}
	}

	if idx := utils.PickWeightedMapId(slt); idx != 0 {
		seq = idx
	}

	return seq
}

func pick_box(plr_data *data_plr_t) int32 {
	seq := int32(0)
	cnt := int32(0)
	slt := map[int32]int32{}
	for _, v := range gamedata.ConfActMazeEventBox.Items() {
		if cnt == 0 {
			seq = v.Seq
			cnt++
		}
		if v.ShowLvCond <= (plr_data.Lv + 1) {
			slt[v.Seq] = v.Weight
		}
	}
	if idx := utils.PickWeightedMapId(slt); idx != 0 {
		seq = idx
	}
	return seq
}

func pick_mine(plr_data *data_plr_t) int32 {
	seq := int32(0)
	cnt := int32(0)
	slt := map[int32]int32{}
	for _, v := range gamedata.ConfActMazeEventMine.Items() {
		if cnt == 0 {
			seq = v.Seq
			cnt++
		}
		if v.ShowLvCond <= (plr_data.Lv + 1) {
			slt[v.Seq] = v.Weight
		}
	}
	if idx := utils.PickWeightedMapId(slt); idx != 0 {
		seq = idx
	}
	return seq
}

func pick_master(plr_data *data_plr_t, tp int32) int32 {
	seq := int32(0)
	cnt := int32(0)
	slt := map[int32]int32{}
	for _, v := range gamedata.ConfActMazeEventMonsterM.QueryItems(tp) {
		if cnt == 0 {
			seq = v.Seq
			cnt++
		}
		if v.ShowLvCond <= (plr_data.Lv + 1) {
			slt[v.Seq] = v.Weight
		}
	}
	if idx := utils.PickWeightedMapId(slt); idx != 0 {
		seq = idx
	}
	return seq
}

func pick_buff(plr_data *data_plr_t) int32 {
	seq := int32(0)
	cnt := int32(0)
	slt := map[int32]int32{}
	for _, v := range gamedata.ConfActMazeBuff.Items() {
		if cnt == 0 {
			seq = v.BuffId
			cnt++
		}
		if v.ShowLvCond <= (plr_data.Lv + 1) {
			slt[v.BuffId] = v.Weight
		}
	}
	if idx := utils.PickWeightedMapId(slt); idx != 0 {
		seq = idx
	}
	return seq
}

func wlevel_aver(plr IPlayer) int32 {
	conf_w := gamedata.ConfWorldLevelM.Query(plr.GetWLevelLvNum())
	if conf_w == nil {
		return plr.GetLevel() + 1
	}

	plr_data := actObj.GetPlrData(plr)

	conf_lv := gamedata.ConfActMazeLv.Query(lv_to_seq(plr_data.Lv))
	if conf_lv == nil {
		return plr.GetLevel() + 1
	}

	sum := int32(0)
	for _, v := range conf_w.Monster {
		sum += v.Lv
	}

	lv := int32(float32(sum/int32(len(conf_w.Monster))) * conf_lv.MonsterLvRatio)

	if lv <= 0 {
		lv = 1
	}

	limit_lv := int32(len(gamedata.ConfMonster.Items()))
	if lv > limit_lv {
		lv = int32(limit_lv)
	}

	return lv
}

// ============================================================================

// 任务领取（每日和成就）
func TakeTask(plr IPlayer, id int32) (int32, *msg.Rewards) {
	if !actObj.Started() && !actObj.Ended() {
		return Err.Act_ActClosed, nil
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		return Err.Act_ActPlrDataNotFound, nil
	}

	conf := gamedata.ConfActMazeTask.Query(id)
	if conf == nil {
		return Err.Failed, nil
	}

	if conf.ConfGrp != actObj.GetConfGrp() {
		return Err.Act_ConfGrp, nil
	}

	op := plr.GetBag().NewOp(gconst.ObjFrom_ActMaze)

	for _, v := range plr_data.Taken {
		if id == v {
			return Err.Plr_TakenBefore, nil
		}
	}

	for _, v := range conf.Attain {
		if plr_data.Attain[v.AttainId].Val < v.P2 {
			return Err.Activity_CondLimited, nil
		} else {
			for _, v := range conf.Reward {
				op.Inc(v.Id, v.N)
			}

			plr_data.Taken = append(plr_data.Taken, id)
		}
	}

	rwds := op.Apply().ToMsg()
	plr_data.Taken = append(plr_data.Taken, id)

	return Err.OK, rwds
}

// 打怪
func ClickBattle(plr IPlayer, pos int32, T *msg.TeamFormation, cb func(ec, score int32, replay *msg.BattleReplay, rwds *msg.Rewards)) {
	if !actObj.Started() {
		cb(Err.Act_ActClosed, 0, nil, nil)
		return
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		cb(Err.Act_ActPlrDataNotFound, 0, nil, nil)
		return
	}

	svr_data := actObj.GetSvrData()
	if svr_data == nil {
		cb(Err.Act_ActSvrDataNotFound, 0, nil, nil)
		return
	}

	conf_lv := gamedata.ConfActMazeLv.Query(lv_to_seq(plr_data.Lv))
	if conf_lv == nil {
		cb(Err.Failed, 0, nil, nil)
		return
	}

	if conf_lv.ConfGrp != actObj.GetConfGrp() {
		cb(Err.Act_ConfGrp, 0, nil, nil)
		return
	}

	conf_g := gamedata.ConfActivityPublic.Query(1)
	if conf_g == nil {
		cb(Err.Failed, 0, nil, nil)
		return
	}

	point_tp := plr_data.Points[pos]
	if (point_tp != C_Maze_Master) && (point_tp != C_Maze_Boss) || !plr_data.Clicked[pos] {
		cb(Err.Activity_MazeNotClicked, 0, nil, nil)
		return
	}

	b1 := plr_data.Battle[pos]
	if b1 == nil {
		cb(Err.Activity_MazeNoBattle, 0, nil, nil)
		return
	}

	conf_m := gamedata.ConfActMazeEventMonster.Query(b1.Seq)
	if conf_m == nil {
		cb(Err.Failed, 0, nil, nil)
		return
	}

	conf_e := gamedata.ConfActMazeEvent.Query(point_tp)
	if conf_e == nil {
		cb(Err.Failed, 0, nil, nil)
		return
	}

	var buffids string
	for buff, cnt := range plr_data.BuffIds {
		buffids = buffids + strings.Repeat(core.I32toa(buff)+",", int(cnt))
	}

	if len(buffids) != 0 {
		buffids = buffids[:len(buffids)-1]
	}

	T2 := battle.NewMonsterTeam()

	Args := map[string]string{}

	for idx, v := range conf_m.Monster {
		lv := b1.Lv

		if b1.Tp == 1 && int32(idx) == conf_m.BossPos-1 {
			lv += 1
		}

		T2.AddMonster(v, lv, int32(idx))
	}

	if len(buffids) != 0 {
		Args["buffs.1"] = buffids
	}

	ratio := conf_lv.PowerRatio
	if plr.GetLevel() > conf_lv.PowerSwitch {
		ratio = math.MaxFloat32(0.25,
			conf_lv.PowerRatio*(1-float32(plr.GetLevel()-conf_lv.PowerSwitch)*0.05))
	}
	T2.ModifyProps(ratio)

	Args["Module"] = "MAZE_BATTLE"

	input := &msg.BattleInput{
		T1:   plr.ToMsg_BattleTeam(T),
		T2:   T2.ToMsg_BattleTeam(),
		Args: Args,
	}

	cop := plr.GetCounter().NewOp(gconst.ObjFrom_ActMaze)

	// 提前计算需要消耗的生命
	cn := int64(conf_g.MazeLoseBlood)
	r := plr.GetCounter().GetRemain(gconst.Cnt_MazeLife)
	if r <= 0 {
		cb(Err.NotEnoughObject(gconst.Cnt_MazeLife), 0, nil, nil)
		return
	}

	if cn > r {
		cn = r
	}

	battle.Fight(input, func(r *msg.BattleResult) {
		if r == nil {
			return
		}

		if r.Winner == 1 {
			for _, v := range conf_m.Reward {
				cop.Inc(v.Id, v.N)
			}

			for i := 0; i < int(conf_m.BuffNum); i++ {
				plr_data.BuffIds[pick_buff(plr_data)] += 1
			}

			delete(plr_data.Battle, pos)
			plr_data.Points[pos] = 0
			svr_data.TotalScore[plr.GetId()] += conf_e.Score

			evtmgr.Fire(gconst.Evt_ActMaze_Score, plr, conf_e.Score)
			evtmgr.Fire(gconst.Evt_ActMaze_Event, plr)
		} else {
			cop.DecCounter(gconst.Cnt_MazeLife, int64(cn))
			if ec := cop.CheckEnough(); ec != Err.OK {
				cb(ec, 0, nil, nil)
				return
			}
		}

		rwds := cop.Apply().ToMsg()

		reply := &msg.BattleReplay{
			Ts: time.Now().Unix(),
			Bi: input,
			Br: r,
		}

		cb(Err.OK, svr_data.TotalScore[plr.GetId()], reply, rwds)
		return
	})
}

// 跳到下一关
func ClickNext(plr IPlayer) (int32, *msg.ActMazeData) {

	if !actObj.Started() {
		return Err.Act_ActClosed, nil
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		return Err.Act_ActPlrDataNotFound, nil
	}

	svr_data := actObj.GetSvrData()
	if svr_data == nil {
		return Err.Act_ActSvrDataNotFound, nil
	}

	conf_lv := gamedata.ConfActMazeLv.Query(lv_to_seq(plr_data.Lv))
	if conf_lv == nil {
		return Err.Failed, nil
	}

	if conf_lv.ConfGrp != actObj.GetConfGrp() {
		return Err.Act_ConfGrp, nil
	}

	conf_g := gamedata.ConfActivityPublic.Query(1)
	if conf_g == nil {
		return Err.Failed, nil
	}

	conf_e := gamedata.ConfActMazeEvent.Query(C_Maze_Exit)
	if conf_e == nil {
		return Err.Failed, nil
	}

	pos := int32(0)
	for idx, v := range plr_data.Points {
		if v == C_Maze_Exit {
			pos = idx
			break
		}
	}

	plr_data.Lv += 1

	reset_maze(plr_data)
	gen_maze(plr_data, pos)

	svr_data.TotalScore[plr.GetId()] += conf_e.Score

	evtmgr.Fire(gconst.Evt_ActMaze_Score, plr, conf_e.Score)

	return Err.OK, actObj.ToMsg(plr)
}

// 购买商人
func ClickTrade(plr IPlayer, pos int32, seq int32) (int32, *msg.Rewards) {
	if !actObj.Started() {
		return Err.Act_ActClosed, nil
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		return Err.Act_ActPlrDataNotFound, nil
	}

	svr_data := actObj.GetSvrData()
	if svr_data == nil {
		return Err.Act_ActSvrDataNotFound, nil
	}

	conf_lv := gamedata.ConfActMazeLv.Query(lv_to_seq(plr_data.Lv))
	if conf_lv == nil {
		return Err.Failed, nil
	}

	if conf_lv.ConfGrp != actObj.GetConfGrp() {
		return Err.Act_ConfGrp, nil
	}

	conf_g := gamedata.ConfActivityPublic.Query(1)
	if conf_g == nil {
		return Err.Failed, nil
	}

	conf_e := gamedata.ConfActMazeEvent.Query(C_Maze_Trade)
	if conf_e == nil {
		return Err.Failed, nil
	}

	point_tp := plr_data.Points[pos]
	if point_tp != C_Maze_Trade || !plr_data.Clicked[pos] {
		return Err.Failed, nil
	}

	op := plr.GetBag().NewOp(gconst.ObjFrom_ActMaze)

	shop := gamedata.ConfActMazeEventShop.Query(seq)
	if shop == nil {
		return Err.Failed, nil
	}

	tdata := plr_data.TradeLimit[pos]
	if tdata == nil {
		tdata = &trade_t{
			Shop: map[int32]int32{},
		}
	}

	if tdata.Shop[seq] >= shop.ExchangeLimit {
		return Err.Activity_MazeBuyCntLimit, nil
	}

	for _, c := range shop.Cost {
		op.Dec(c.Id, c.N)
	}

	if ec := op.CheckEnough(); ec != Err.OK {
		return ec, nil
	}

	for _, v := range shop.Id {
		op.Inc(v.Id, v.N)
	}

	tdata.Shop[seq] += 1

	plr_data.TradeLimit[pos] = tdata

	svr_data.TotalScore[plr.GetId()] += conf_e.Score

	evtmgr.Fire(gconst.Evt_ActMaze_Score, plr, conf_e.Score)

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

// 领取物品（道具、宝箱、宝矿、体力、钥匙）
func ClickThing(plr IPlayer, pos int32, item_id int32) (int32, int32, *msg.Rewards) {
	if !actObj.Started() {
		return Err.Act_ActClosed, 0, nil
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		return Err.Act_ActPlrDataNotFound, 0, nil
	}

	svr_data := actObj.GetSvrData()
	if svr_data == nil {
		return Err.Act_ActSvrDataNotFound, 0, nil
	}

	conf_lv := gamedata.ConfActMazeLv.Query(lv_to_seq(plr_data.Lv))
	if conf_lv == nil {
		return Err.Failed, 0, nil
	}

	if conf_lv.ConfGrp != actObj.GetConfGrp() {
		return Err.Act_ConfGrp, 0, nil
	}

	conf_g := gamedata.ConfActivityPublic.Query(1)
	if conf_g == nil {
		return Err.Failed, 0, nil
	}

	point_tp := plr_data.Points[pos]
	if (point_tp != C_Maze_Power) && (point_tp != C_Maze_Key) &&
		(point_tp != C_Maze_Item) && (point_tp != C_Maze_Box) &&
		(point_tp != C_Maze_Mine) {
		return Err.Failed, 0, nil
	}

	conf_e := gamedata.ConfActMazeEvent.Query(point_tp)
	if conf_e == nil {
		return Err.Failed, 0, nil
	}

	cop := plr.GetCounter().NewOp(gconst.ObjFrom_ActMaze)

	switch point_tp {
	case C_Maze_Power:
		cop.IncCounter(gconst.Cnt_MazeCnt, int64(conf_g.MazePowerFree))

	case C_Maze_Key:
		for _, v := range conf_g.MazeKeyReward {
			cop.Inc(v.Id, v.N)
		}

	case C_Maze_Item:
		item_seq := plr_data.Item[pos]
		if item_seq == 0 {
			return Err.Failed, 0, nil
		}

		conf_i := gamedata.ConfActMazeEventItem.Query(item_seq)
		if conf_i == nil {
			return Err.Failed, 0, nil
		}

		for _, v := range conf_i.Reward {
			cop.Inc(v.Id, v.N)
		}

		delete(plr_data.Item, pos)

	case C_Maze_Box:
		box_seq := plr_data.Box[pos]
		if box_seq == 0 {
			return Err.Failed, 0, nil
		}

		conf_b := gamedata.ConfActMazeEventBox.Query(box_seq)
		if conf_b == nil {
			return Err.Failed, 0, nil
		}

		flag := false
		for _, v := range conf_b.Drop {
			if v.Item == item_id {

				cop.Dec(item_id, 1)

				if ec := cop.CheckEnough(); ec != Err.OK {
					return ec, 0, nil
				}

				for _, v2 := range utils.Drop(plr, v.DropId) {
					cop.Inc(v2.Id, v2.N)
				}

				flag = true
				break
			}
		}

		if !flag {
			return Err.Failed, 0, nil
		}

		delete(plr_data.Box, pos)

	case C_Maze_Mine:
		mine_seq := plr_data.Mine[pos]
		if mine_seq == 0 {
			return Err.Failed, 0, nil
		}

		conf_m := gamedata.ConfActMazeEventMine.Query(mine_seq)
		if conf_m == nil {
			return Err.Failed, 0, nil
		}

		flag := false
		for _, v := range conf_m.Drop {
			if v.Item == item_id {

				cop.Dec(item_id, 1)

				if ec := cop.CheckEnough(); ec != Err.OK {
					return ec, 0, nil
				}

				for _, v2 := range utils.Drop(plr, v.DropId) {
					cop.Inc(v2.Id, v2.N)
				}

				flag = true
				break
			}
		}

		if !flag {
			return Err.Failed, 0, nil
		}

		delete(plr_data.Mine, pos)
	}

	svr_data.TotalScore[plr.GetId()] += conf_e.Score

	evtmgr.Fire(gconst.Evt_ActMaze_Event, plr)
	evtmgr.Fire(gconst.Evt_ActMaze_Score, plr, conf_e.Score)

	plr_data.Points[pos] = 0

	rwds := cop.Apply().ToMsg()

	return Err.OK, svr_data.TotalScore[plr.GetId()], rwds
}

// 翻格子
func Click(plr IPlayer, pos int32, cb func(int32, int32, int32, *msg.BattleData, *msg.Rewards)) {
	var r_seq int32
	var b1 *msg.BattleData

	if !actObj.Started() {
		cb(Err.Act_ActClosed, 0, 0, nil, nil)
		return
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		cb(Err.Plr_NotFound, 0, 0, nil, nil)
		return
	}

	svr_data := actObj.GetSvrData()
	if svr_data == nil {
		cb(Err.Act_ActSvrDataNotFound, 0, 0, nil, nil)
		return
	}

	conf_lv := gamedata.ConfActMazeLv.Query(lv_to_seq(plr_data.Lv))
	if conf_lv == nil {
		cb(Err.Failed, 0, 0, nil, nil)
		return
	}

	if conf_lv.ConfGrp != actObj.GetConfGrp() {
		cb(Err.Act_ConfGrp, 0, 0, nil, nil)
		return
	}

	conf_g := gamedata.ConfActivityPublic.Query(1)
	if conf_g == nil {
		cb(Err.Failed, 0, 0, nil, nil)
		return
	}

	point_tp := plr_data.Points[pos]
	if plr_data.Clicked[pos] {
		cb(Err.Activity_MazeAlreadyClicked, 0, 0, nil, nil)
		return
	}

	cop := plr.GetCounter().NewOp(gconst.ObjFrom_ActMaze)
	cop.DecCounter(gconst.Cnt_MazeCnt, 1)
	if ec := cop.CheckEnough(); ec != Err.OK {
		cb(ec, 0, 0, nil, nil)
		return
	}

	switch {
	case point_tp == C_Maze_Item: // 道具
		r_seq = pick_item(plr_data)
		plr_data.Item[pos] = r_seq

	case point_tp == C_Maze_Box: // 宝箱
		r_seq = pick_box(plr_data)
		plr_data.Box[pos] = r_seq

	case point_tp == C_Maze_Mine: // 矿
		r_seq = pick_mine(plr_data)
		plr_data.Mine[pos] = r_seq

	case point_tp == C_Maze_Boss: // Boss

		boss_seq := pick_master(plr_data, 1)
		boss_lv := wlevel_aver(plr)

		b := &battle_t{
			Tp:  1,
			Lv:  boss_lv,
			Seq: boss_seq,
		}

		r_seq = b.Seq

		b1 = &msg.BattleData{
			Tp:  1,
			Lv:  boss_lv,
			Seq: boss_seq,
		}

		plr_data.Battle[pos] = b

	case point_tp == C_Maze_Master: // Master
		monster_seq := pick_master(plr_data, 0)
		monster_lv := wlevel_aver(plr)

		b := &battle_t{
			Tp:  0,
			Lv:  monster_lv,
			Seq: monster_seq,
		}

		r_seq = b.Seq

		b1 = &msg.BattleData{
			Tp:  0,
			Lv:  monster_seq,
			Seq: monster_seq,
		}
		plr_data.Battle[pos] = b

	case point_tp == C_Maze_Trap: // 陷阱: 直接扣血
		cn := int64(conf_g.MazeTrapBlood)
		r := plr.GetCounter().GetRemain(gconst.Cnt_MazeLife)
		if r <= 0 {
			cb(Err.NotEnoughObject(gconst.Cnt_MazeLife), 0, 0, nil, nil)
			return
		}

		if cn > r {
			cn = r
		}

		cop.DecCounter(gconst.Cnt_MazeLife, int64(cn))
		if ec := cop.CheckEnough(); ec != Err.OK {
			cb(ec, 0, 0, nil, nil)
			return
		}

		plr_data.Points[pos] = 0
		evtmgr.Fire(gconst.Evt_ActMaze_Event, plr)
	}

	plr_data.Clicked[pos] = true
	svr_data.TotalScore[plr.GetId()] += 1

	evtmgr.Fire(gconst.Evt_ActMaze_Click, plr)
	evtmgr.Fire(gconst.Evt_ActMaze_Score, plr, int32(1))

	rwds := cop.Apply().ToMsg()

	cb(Err.OK, svr_data.TotalScore[plr.GetId()], r_seq, b1, rwds)
	return
}

// 重置
func Reset(plr IPlayer) (int32, *msg.ActMazeData) {
	if !actObj.Started() {
		return Err.Act_ActClosed, nil
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		return Err.Act_ActPlrDataNotFound, nil
	}

	cop := plr.GetCounter().NewOp(gconst.ObjFrom_ActMaze)

	if plr.GetCounter().Get(gconst.Cnt_MazeLife) > 0 {
		cop.SetCounter(gconst.Cnt_MazeLife, 0)
	}

	plr_data.Lv = 0

	reset_maze(plr_data)
	gen_maze(plr_data, 16)

	cop.Apply()

	return Err.OK, actObj.ToMsg(plr)
}

func GetBuffIds(plr IPlayer) (int32, map[int32]int32) {
	if !actObj.Started() {
		return Err.Act_ActClosed, nil
	}

	plr_data := actObj.GetPlrData(plr)
	if plr_data == nil {
		return Err.Act_ActPlrDataNotFound, nil
	}

	return Err.OK, plr_data.BuffIds
}

func ActMazeInfo(plr IPlayer) *msg.ActMazeData {
	return actObj.ToMsg(plr)
}
