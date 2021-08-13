package comp

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

// ============================================================================

// 计数
type Counter struct {
	Cnt    map[int32]int64     // 当前值
	maxCnt map[int32]int64     // 最大值
	Ts     map[int32]time.Time // 上次恢复时间

	plr IPlayer
}

type CounterOp struct {
	cnt    map[int32]int64
	maxCnt map[int32]int64

	c *Counter
	*BagOp
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_PlrResetDaily, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		plr.GetCounter().reset_daily()
	})
}

// ============================================================================

func NewCounter() *Counter {
	return &Counter{
		Cnt:    make(map[int32]int64),
		maxCnt: make(map[int32]int64),
		Ts:     make(map[int32]time.Time),
	}
}

func (self *Counter) NewOp(from int32) *CounterOp {
	ret := &CounterOp{
		cnt:    make(map[int32]int64),
		maxCnt: make(map[int32]int64),

		c: self,
	}

	ret.BagOp = self.plr.GetBag().NewOp(from)

	return ret
}

// ============================================================================

func (self *Counter) Init(plr IPlayer) {
	self.plr = plr
	self.maxCnt = make(map[int32]int64)

	self.init_max()
}

func (self *Counter) init_max() {
	for _, conf := range gamedata.ConfCounter.Items() {
		if conf.MaxValue > 0 {
			self.Add(conf.Id, int64(conf.MaxValue), true)
		}
	}
}

func (self *Counter) reset_daily() {
	now := time.Now()

	for _, conf := range gamedata.ConfCounter.Items() {
		switch conf.Type {
		// 每天重置
		case gconst.CounterType_Daily:
			if self.Cnt[conf.Id] > 0 {
				self.Cnt[conf.Id] = 0
			}

			self.Ts[conf.Id] = now

		// 每周重置
		case gconst.CounterType_Weekly:
			t := self.Ts[conf.Id]

			if t.IsZero() || !core.IsSameWeek(now, t) {
				if self.Cnt[conf.Id] > 0 {
					self.Cnt[conf.Id] = 0
				}

				self.Ts[conf.Id] = now
			}
		}
	}
}

func (self *Counter) mod_ts(id int32) {
	conf := gamedata.ConfCounter.Query(id)
	if conf != nil && conf.Type == gconst.CounterType_Recover {
		self.Ts[id] = time.Now()
	}
}

// 获取本计数值
func (self *Counter) Get(id int32) int64 {
	return self.Cnt[id]
}

// 获取本计数的上限值
func (self *Counter) GetMax(id int32) int64 {
	return self.maxCnt[id]
}

// 获取本计数的剩余次数
func (self *Counter) GetRemain(id int32) int64 {
	return self.maxCnt[id] - self.Cnt[id]
}

// 获取本计数的恢复时间
func (self *Counter) GetTs(id int32) time.Time {
	return self.Ts[id]
}

// 直接增加计数次数
func (self *Counter) Add(id int32, n int64, isMax ...bool) {
	if n == 0 {
		return
	}

	if core.DefFalse(isMax) {
		v := self.maxCnt[id] + n
		if v == 0 {
			delete(self.maxCnt, id)
		} else {
			self.maxCnt[id] = v
		}
	} else {
		v := self.Cnt[id] + n
		if v == 0 {
			delete(self.Cnt, id)
		} else {
			self.Cnt[id] = v

			// <=0 to >0
			if v-n <= 0 && v > 0 {
				self.mod_ts(id)
			}
		}
	}
}

// 检查恢复
func (self *Counter) CheckRecover(id int32) (ec int32, cnt int64, ts int64) {
	conf := gamedata.ConfCounter.Query(id)
	if conf == nil || conf.Type != gconst.CounterType_Recover ||
		len(conf.Recover) == 0 || conf.Recover[0].Sec <= 0 {
		ec = Err.Failed
		return
	}

	// 无需恢复
	if self.Cnt[id] <= 0 {
		ec = Err.Failed
		return
	}

	// calc
	t := self.Ts[id]
	m := int32(time.Since(t).Seconds()) / conf.Recover[0].Sec
	n := m * conf.Recover[0].N
	if n == 0 {
		ec = Err.Failed
		return
	}

	self.Ts[id] = t.Add(time.Duration(m*conf.Recover[0].Sec) * time.Second)

	self.Cnt[id] -= int64(n)
	if self.Cnt[id] < 0 {
		self.Cnt[id] = 0
	}

	cnt, ts = self.Cnt[id], self.Ts[id].Unix()

	// offline counter full
	if !self.plr.IsOnline() && cnt <= 0 && cnt+int64(n) > 0 {
		evtmgr.Fire(gconst.Evt_OfflineCounterFull, self.plr, id, self.maxCnt[id])
	}

	return
}

// 购买
func (self *Counter) Buy(id int32) int32 {
	conf := gamedata.ConfCounter.Query(id)
	if conf == nil {
		return Err.Counter_NotFound
	}

	if conf.BuyId == 0 {
		return Err.Counter_CanNotBuy
	}

	// from := fmt.Sprintf("%s%d", gconst.ObjFrom_CounterBuy, id)
	cop := self.NewOp(id)

	cop.DecCounter(conf.BuyId, 1)
	if ec := cop.CheckEnough(); ec != Err.OK {
		return ec
	}

	// 一次补满
	if conf.BuyN == 0 {
		if self.GetRemain(id) > 0 {
			return Err.Counter_Enough
		}

		cop.SetCounter(id, 0)
	} else {
		cop.IncCounter(id, conf.BuyN)
	}

	cop.Apply()

	return Err.OK
}

func (self *Counter) ToMsg() *msg.CounterData {
	ret := &msg.CounterData{
		Cnt:    self.Cnt,
		MaxCnt: self.maxCnt,
		Ts:     make(map[int32]int64),
	}

	for id, v := range self.Ts {
		ret.Ts[id] = v.Unix()
	}

	return ret
}

// ============================================================================
// cop

func (self *CounterOp) AddCounter(id int32, n int64, isMax ...bool) {
	if n == 0 {
		return
	}

	conf := gamedata.ConfCounter.Query(id)
	if conf == nil {
		return
	}

	if core.DefFalse(isMax) {
		self.maxCnt[id] += n
	} else {
		// opCost
		L := int64(len(conf.OpCost))
		if L > 0 {
			for i := int64(0); i < n; i++ {
				idx := self.c.Get(id) + i
				if idx >= L {
					idx = L - 1
				}

				self.BagOp.Dec(conf.OpCost[idx].Id, conf.OpCost[idx].N)
			}
		}

		self.cnt[id] += n
	}
}

// 增加计数值(远离max)
func (self *CounterOp) IncCounter(id int32, n int64) {
	if n <= 0 {
		return
	}

	self.AddCounter(id, -n)
}

// 增加计数值Max
func (self *CounterOp) IncCounterMax(id int32, n int64) {
	if n <= 0 {
		return
	}

	self.AddCounter(id, n, true)
}

// 减少计数值(向max靠近)
func (self *CounterOp) DecCounter(id int32, n int64) {
	if n <= 0 {
		return
	}

	self.AddCounter(id, n)
}

// 减少计数值Max
func (self *CounterOp) DecCounterMax(id int32, n int64) {
	if n <= 0 {
		return
	}

	self.AddCounter(id, -n, true)
}

// 设置值
func (self *CounterOp) SetCounter(id int32, n int64, isMax ...bool) {
	if core.DefFalse(isMax) {
		self.maxCnt[id] = 0
		self.AddCounter(id, -self.c.GetMax(id)+n, true)
	} else {
		self.cnt[id] = 0
		self.AddCounter(id, -self.c.Get(id)+n)
	}
}

// 回退计数值
func (self *CounterOp) RetCounter(id int32, n int64) {
	if n <= 0 {
		return
	}

	conf := gamedata.ConfCounter.Query(id)
	if conf == nil {
		return
	}

	// opCost
	L := int64(len(conf.OpCost))
	if L > 0 {
		for i := int64(0); i < n; i++ {
			idx := self.c.Get(id) + i
			if idx >= L {
				idx = L - 1
			}

			self.BagOp.Ret(conf.OpCost[idx].Id, conf.OpCost[idx].N)
		}
	}

	self.cnt[id] -= n
}

func (self *CounterOp) CheckEnough() int32 {
	for id, v := range self.cnt {
		if v < 0 {
			continue
		}

		conf := gamedata.ConfCounter.Query(id)
		if conf == nil {
			return Err.Counter_NotFound
		}

		if conf.Type != gconst.CounterType_Unlimit &&
			self.c.Get(id)+v > self.c.GetMax(id)+self.maxCnt[id] {

			return Err.NotEnoughObject(id)
		}
	}

	if ec := self.BagOp.CheckEnough(); ec != Err.OK {
		return ec
	}

	return Err.OK
}

func (self *CounterOp) Apply() (rwds *rewards_t) {
	rwds = self.BagOp.Apply()

	// counter apply
	if len(self.cnt) == 0 && len(self.maxCnt) == 0 {
		return
	}

	// 增量
	res := &msg.GS_CounterOpUpdate{
		Cnt:    make(map[int32]int64),
		MaxCnt: make(map[int32]int64),
		Ts:     make(map[int32]int64),
	}

	for id, v := range self.cnt {
		self.c.Add(id, v)

		res.Cnt[id] += v
		t := self.c.GetTs(id)
		if !t.IsZero() {
			res.Ts[id] = t.Unix()
		}
	}

	for id, v := range self.maxCnt {
		self.c.Add(id, v, true)

		res.MaxCnt[id] += v
	}

	self.c.plr.SendMsg(res)

	return
}
