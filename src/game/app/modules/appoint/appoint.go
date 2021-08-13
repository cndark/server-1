package appoint

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math/rand"
	"time"
)

// ============================================================================

const (
	C_Appoint_Hero_Star = 1 // 英雄星级,英雄最大星级
	C_Appoint_Hero_Elem = 2 // 英雄阵营,对应元素的英雄个数
	C_Appoint_Hero_Job  = 3 // 英雄职业,对应职业的英雄个数

	C_Appoint_Cnt_Max = 100 // 任务数总上限
)

// ============================================================================

// 酒馆探索派遣
type Appoint struct {
	TotalTake int32     // 总共生成的派遣数量
	Tasks     task_m    // 任务列表
	AddTs     time.Time // 上次新增任务时间

	heroes map[int64]bool
	plr    IPlayer
}

type task_m map[int32]*task_t // key为递增的seq
type task_t struct {
	Id     int32           // 表索引,可重复
	Cond   map[int32]int32 // 条件
	IsLock bool            // 是否锁定
	Ts     time.Time       // 结束时间
	Heroes []int64         // 上阵英雄
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_MOpen, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		mid := args[1].(int32)

		if mid == gconst.ModuleId_Appoint {
			plr.GetAppoint().CheckAdd()
		}
	})
}

func NewAppoint() *Appoint {
	return &Appoint{
		Tasks: make(task_m),
	}
}

// ============================================================================

func (self *Appoint) Init(plr IPlayer) {
	self.plr = plr

	self.heroes = make(map[int64]bool)
	for _, t := range self.Tasks {
		for _, v := range t.Heroes {
			self.heroes[v] = true
		}
	}
}

// star1 star2 为0 表示通随
func (self *Appoint) AddTask(n int32, star1, star2 int32, send_update bool) (addN int32) {
	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil {
		return 0
	}

	ids := self.pick_n(n, star1, star2)
	if len(ids) == 0 {
		return 0
	}

	var taskRes []*msg.AppointTask
	for _, id := range ids {
		conf := gamedata.ConfAppoint.Query(id)
		if conf == nil {
			continue
		}

		t := &task_t{
			Id:   id,
			Cond: make(map[int32]int32),
		}

		for _, v := range conf.Require {
			switch v.Type {
			case C_Appoint_Hero_Star:
				t.Cond[v.Type] = v.N
			case C_Appoint_Hero_Elem:
				L := len(conf.HeroElemRange)
				if L > 0 {
					t.Cond[v.Type] = conf.HeroElemRange[rand.Intn(L)]*100 + v.N
				}
			case C_Appoint_Hero_Job:
				L := len(conf_g.HeroJobRange)
				if L > 0 {
					t.Cond[v.Type] = conf_g.HeroJobRange[rand.Intn(L)]*100 + v.N
				}
			}
		}

		self.TotalTake++
		self.Tasks[self.TotalTake] = t

		taskRes = append(taskRes, &msg.AppointTask{
			Seq:    self.TotalTake,
			Id:     t.Id,
			Cond:   t.Cond,
			Ts:     t.Ts.Unix(),
			IsLock: t.IsLock,
			Heroes: t.Heroes,
		})

		addN++

		if len(self.Tasks) >= C_Appoint_Cnt_Max {
			break
		}
	}

	// notify
	if len(taskRes) > 0 && send_update {
		self.plr.SendMsg(&msg.GS_AppointAddTask{
			Tasks: taskRes,
		})
	}

	return
}

func (self *Appoint) pick_n(n int32, star1, star2 int32) (ids []int32) {
	slt := make(map[int32]int32)

	if star1 == 0 && star2 == 0 {
		for _, v := range gamedata.ConfAppoint.Items() {
			slt[v.Id] += v.Weight
		}
	} else {
		for i := star1; i <= star2; i++ {
			for _, v := range gamedata.ConfAppointM.QueryItems(i) {
				slt[v.Id] += v.Weight
			}
		}
	}

	for i := int32(0); i < n; i++ {
		id := utils.PickWeightedMapId(slt)
		if id != 0 {
			ids = append(ids, id)
		}
	}

	return
}

// 检查恢复
func (self *Appoint) CheckAdd() {
	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil || conf_g.AppointAddTaskTime == 0 {
		return
	}

	now := time.Now()
	cnt := int64(0)
	for _, v := range self.Tasks {
		if len(v.Heroes) > 0 && v.Ts.Before(now) {
			continue
		}
		cnt++
	}

	max := self.plr.GetCounter().GetMax(gconst.Cnt_AppointTask)
	if cnt >= max {
		return
	}

	sub := int64(time.Since(self.AddTs).Minutes())
	n := sub / int64(conf_g.AppointAddTaskTime)

	if n > max-cnt {
		n = max - cnt
	}

	if n <= 0 {
		return
	}

	self.AddTs = now
	self.AddTask(int32(n), 0, 0, true)
}

// 刷新
func (self *Appoint) Refresh() int32 {
	opSeq := []int32{}
	for seq, v := range self.Tasks {
		if v.IsLock {
			continue
		}

		if len(v.Heroes) > 0 {
			continue
		}
		opSeq = append(opSeq, seq)
	}

	L := int32(len(opSeq))
	if L == 0 {
		return Err.Appoint_AllLock
	}

	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return Err.Failed
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_AppointRefresh)
	for _, v := range conf.AppointRefreshCost {
		op.Dec(v.Id, v.N*L)
	}

	if ec := op.CheckEnough(); ec != Err.OK {
		return ec
	}

	self.AddTask(L, 0, 0, false)

	for _, v := range opSeq {
		delete(self.Tasks, v)
	}

	op.Apply()

	return Err.OK
}

// seq == 0 表示一键操作
func (self *Appoint) Lock(seq int32, isLock bool) {
	t := self.Tasks[seq]
	if t == nil || seq == 0 {
		for _, t := range self.Tasks {
			if len(t.Heroes) == 0 {
				t.IsLock = isLock
			}
		}
	} else if len(t.Heroes) == 0 {
		t.IsLock = isLock
	}
}

// 派遣
func (self *Appoint) Send(seq int32, heroes []int64) (int32, int64) {
	task := self.Tasks[seq]
	if task == nil {
		return Err.Appoint_TaskNotFound, 0
	}

	if len(task.Heroes) > 0 {
		return Err.Appoint_TaskSending, 0
	}

	conf_t := gamedata.ConfAppoint.Query(task.Id)
	if conf_t == nil {
		return Err.Failed, 0
	}

	if int32(len(heroes)) > conf_t.HeroNum {
		return Err.Appoint_TooManyHero, 0
	}

	// check cond
	for tp, val := range task.Cond {
		switch tp {
		case C_Appoint_Hero_Star:
			maxStar := int32(0)
			for _, v := range heroes {
				if self.heroes[v] {
					return Err.Appoint_HeroSending, 0
				}

				hero := self.plr.GetBag().FindHero(v)
				if hero == nil {
					return Err.Hero_NotFound, 0
				}

				if hero.Star > maxStar {
					maxStar = hero.Star
				}
			}

			if maxStar < val {
				return Err.Appoint_HeroCondLimit, 0
			}
		case C_Appoint_Hero_Elem:
			elemN := int32(0)
			for _, v := range heroes {
				if self.heroes[v] {
					return Err.Appoint_HeroSending, 0
				}

				hero := self.plr.GetBag().FindHero(v)
				if hero == nil {
					return Err.Hero_NotFound, 0
				}

				conf_m := gamedata.ConfMonster.Query(hero.Id)
				if conf_m != nil && conf_m.Elem == (val/100) {
					elemN++
				}
			}

			if elemN < val%100 {
				return Err.Appoint_HeroCondLimit, 0
			}
		case C_Appoint_Hero_Job:
			jobN := int32(0)
			for _, v := range heroes {
				if self.heroes[v] {
					return Err.Appoint_HeroSending, 0
				}
				hero := self.plr.GetBag().FindHero(v)
				if hero == nil {
					return Err.Hero_NotFound, 0
				}

				conf_m := gamedata.ConfMonster.Query(hero.Id)
				if conf_m != nil && conf_m.JobId == (val/100) {
					jobN++
				}
			}

			if jobN < val%100 {
				return Err.Appoint_HeroCondLimit, 0
			}
		}
	}

	// cost
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_AppointSend)
	for _, v := range conf_t.Cost {
		op.Dec(v.Id, v.N)
	}

	if ec := op.CheckEnough(); ec != Err.OK {
		return ec, 0
	}

	op.Apply()

	for _, v := range heroes {
		self.heroes[v] = true
	}

	task.Heroes = heroes
	task.IsLock = true
	task.Ts = time.Now().Add(time.Duration(conf_t.Time) * time.Second)

	evtmgr.Fire(gconst.Evt_AppointSend, self.plr, conf_t.Star)

	return Err.OK, task.Ts.Unix()
}

// 加速
func (self *Appoint) Acc(seq int32) int32 {
	task := self.Tasks[seq]
	if task == nil {
		return Err.Appoint_TaskNotFound
	}

	if len(task.Heroes) == 0 {
		return Err.Appoint_TaskNotSend
	}

	conf_t := gamedata.ConfAppoint.Query(task.Id)
	if conf_t == nil {
		return Err.Failed
	}

	if task.Ts.Before(time.Now()) {
		return Err.Appoint_TaskFinished
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_AppointAcc)
	for _, v := range conf_t.AccCost {
		op.Dec(v.Id, v.N)
	}

	if ec := op.CheckEnough(); ec != Err.OK {
		return ec
	}

	op.Apply()

	task.Ts = time.Unix(0, 0)

	return Err.OK
}

// 领奖
func (self *Appoint) Take(seq int32) (int32, *msg.Rewards) {
	task := self.Tasks[seq]
	if task == nil {
		return Err.Appoint_TaskNotFound, nil
	}

	if len(task.Heroes) == 0 {
		return Err.Appoint_TaskNotSend, nil
	}

	conf_t := gamedata.ConfAppoint.Query(task.Id)
	if conf_t == nil {
		return Err.Failed, nil
	}

	if task.Ts.After(time.Now()) {
		return Err.Appoint_TaskNotFinished, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_AppointTake)
	for _, v := range conf_t.Reward {
		op.Inc(v.Id, v.N)
	}

	rwds := op.Apply().ToMsg()

	for _, v := range task.Heroes {
		delete(self.heroes, v)
	}

	delete(self.Tasks, seq)

	maxN := self.plr.GetCounter().GetMax(gconst.Cnt_AppointTask)
	if int64(len(self.Tasks)) == maxN-1 {
		self.AddTs = time.Now()
	}

	evtmgr.Fire(gconst.Evt_AppointFin, self.plr, conf_t.Star)

	return Err.OK, rwds
}

// 取消
func (self *Appoint) Cancel(seq int32) int32 {
	task := self.Tasks[seq]
	if task == nil {
		return Err.Appoint_TaskNotFound
	}

	if len(task.Heroes) == 0 {
		return Err.Appoint_TaskNotSend
	}

	if task.Ts.Before(time.Now()) {
		return Err.Appoint_TaskFinished
	}

	for _, seq := range task.Heroes {
		delete(self.heroes, seq)
	}

	task.Heroes = []int64{}
	task.Ts = time.Unix(0, 0)
	task.IsLock = false

	return Err.OK
}

// ============================================================================

func (self *Appoint) ToMsg() *msg.AppointData {
	ret := &msg.AppointData{
		AddTs: self.AddTs.Unix(),
	}

	for seq, v := range self.Tasks {
		ret.Tasks = append(ret.Tasks, &msg.AppointTask{
			Seq:    seq,
			Id:     v.Id,
			Cond:   v.Cond,
			Ts:     v.Ts.Unix(),
			IsLock: v.IsLock,
			Heroes: v.Heroes,
		})
	}

	return ret
}
