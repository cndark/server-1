package comp

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math"
)

// ============================================================================

type ccy_map_t = map[int32]int64  // 货币容器
type item_map_t = map[int32]int32 // 道具容器
type hero_map_t = map[int64]*Hero // 英雄容器
type armor_map_t map[int64]*Armor // 装备容器: 仅用于穿戴。 装备实质上是 Item
type relic_map_t map[int64]*Relic // 神器容器

type Bag struct {
	Ccy    ccy_map_t
	Items  item_map_t
	Heroes hero_map_t
	Armors armor_map_t // 仅用于穿戴。 装备实质上是 Item
	Relics relic_map_t

	plr IPlayer
}

type BagOp struct {
	Ccy       ccy_map_t
	Items     item_map_t
	Heroes    hero_map_t
	HeroesDel []int64
	Relics    relic_map_t
	RelicsDel []int64

	rewards *rewards_t

	bag  *Bag
	from int32
}

type rewards_t struct {
	Ccy      ccy_map_t  //
	Items    item_map_t //
	HeroIds  []int64    // id -> seq
	RelicIds []int64    // id -> seq
}

// ============================================================================

func NewBag() *Bag {
	bag := &Bag{
		Ccy:    make(ccy_map_t),
		Items:  make(item_map_t),
		Heroes: make(hero_map_t),
		Armors: make(armor_map_t),
		Relics: make(relic_map_t),
	}

	return bag
}

func (self *Bag) Init(plr IPlayer) {
	self.plr = plr

	// init heroes
	for _, hero := range self.Heroes {
		hero.init(self)
	}

	// init armors
	for _, ar := range self.Armors {
		ar.init(self)
	}

	// init relics
	for _, rlc := range self.Relics {
		rlc.init(self)
	}
}

func (self *Bag) ApplyHeroMods() {
	for _, hero := range self.Heroes {
		hero.apply_mods(false)
	}
}

func (self *Bag) NewOp(from int32) *BagOp {
	return &BagOp{
		Ccy:       make(ccy_map_t),
		Items:     make(item_map_t),
		Heroes:    make(hero_map_t),
		HeroesDel: nil,
		Relics:    make(relic_map_t),
		RelicsDel: nil,

		rewards: &rewards_t{
			Ccy:      make(ccy_map_t),
			Items:    make(item_map_t),
			HeroIds:  nil,
			RelicIds: nil,
		},

		bag:  self,
		from: from,
	}
}

func (self *Bag) CheckFull(n ...int32) int32 {
	return Err.OK
}

func (self *Bag) ToMsg() *msg.BagData {
	return &msg.BagData{
		Currency: self.ToMsg_CcyArray(),
		Items:    self.ToMsg_ItemArray(),
		Heroes:   self.ToMsg_HeroArray(),
		Armors:   self.ToMsg_ArmorArray(),
		Relics:   self.ToMsg_RelicArray(),
	}
}

// ============================================================================

func (self *BagOp) assert_n_val(n interface{}, sign int32) (int64, float64) {
	var ival int64
	var fval float64
	var isfloat bool

	switch n.(type) {
	case int:
		ival = int64(n.(int))
		fval = float64(n.(int))
	case int32:
		ival = int64(n.(int32))
		fval = float64(n.(int32))
	case int64:
		ival = n.(int64)
		fval = float64(n.(int64))

	case float32:
		ival = int64(n.(float32))
		fval = float64(n.(float32))
		isfloat = true
	case float64:
		ival = int64(n.(float64))
		fval = n.(float64)
		isfloat = true

	default:
		return 0, 0
	}

	// invalid floating number check
	if isfloat && (math.IsInf(fval, 0) || math.IsNaN(fval)) {
		return 0, 0
	}

	// normal checks
	if fval == 0 {
		return 0, 0
	}

	if sign > 0 {
		if fval < 0 {
			return 0, 0
		}
	} else if sign < 0 {
		if fval < 0 {
			return 0, 0
		}
		ival, fval = -ival, -fval
	} else {
		// sign == 0: as is
	}

	return ival, fval
}

func (self *BagOp) Add(id int32, n interface{}, sign int32, isRet ...bool) *BagOp {
	ival, fval := self.assert_n_val(n, sign)
	if fval == 0 {
		return self
	}

	// add
	switch gconst.ObjectType(id) {
	case gconst.ObjType_Currency:
		if ival > 0 {
			if core.DefFalse(isRet) {
				self.Ccy[id] += ival
			} else {
				self.rewards.Ccy[id] += ival
			}
		} else {
			self.Ccy[id] += ival
		}

	case gconst.ObjType_Item:
		if ival > 0 {
			if core.DefFalse(isRet) {
				self.Items[id] += int32(ival)
			} else {
				self.rewards.Items[id] += int32(ival)
			}
		} else {
			self.Items[id] += int32(ival)
		}

	case gconst.ObjType_Hero:
		for i := int64(0); i < ival; i++ {
			// !important: limit hero number per id
			if ival > 20 {
				ival = 20
			}

			self.rewards.HeroIds = append(self.rewards.HeroIds, int64(id))
		}

	case gconst.ObjType_Relic:
		for i := int64(0); i < ival; i++ {
			// !important: limit relic number per id
			if ival > 50 {
				ival = 50
			}

			self.rewards.RelicIds = append(self.rewards.RelicIds, int64(id))
		}

	default:
		log.Warning("invalid object id:", id)
	}

	return self
}

func (self *BagOp) Inc(id int32, n interface{}) *BagOp {
	return self.Add(id, n, 1)
}

func (self *BagOp) Dec(id int32, n interface{}) *BagOp {
	return self.Add(id, n, -1)
}

// just return item and ccy, n > 0
func (self *BagOp) Ret(id int32, n interface{}) *BagOp {
	if gconst.IsItem(id) || gconst.IsCurrency(id) {
		return self.Add(id, n, 1, true)
	}

	return self
}

func (self *BagOp) DelHero(seq int64) *BagOp {
	self.HeroesDel = append(self.HeroesDel, seq)
	return self
}

func (self *BagOp) DelRelic(seq int64) *BagOp {
	self.RelicsDel = append(self.RelicsDel, seq)
	return self
}

func (self *BagOp) CheckEnough() int32 {
	// check currency
	for id, val := range self.Ccy {
		if val >= 0 {
			continue
		}

		if -val > self.bag.Ccy[id] {
			return Err.NotEnoughObject(id)
		}
	}

	// check items
	for id, num := range self.Items {
		if num >= 0 {
			continue
		}

		if -num > self.bag.Items[id] {
			return Err.NotEnoughObject(id)
		}
	}

	// check heroes
	for _, seq := range self.HeroesDel {
		if self.bag.Heroes[seq] == nil {
			return Err.Bag_NotEnoughHero
		}
	}

	// check relics
	for _, seq := range self.RelicsDel {
		if self.bag.Relics[seq] == nil {
			return Err.Bag_NotEnoughRelic
		}
	}

	// ok
	return Err.OK
}

func (self *BagOp) Apply() *rewards_t {

	// convert rewards
	self.convert_rewards()

	// check if there's any change
	if self.is_empty() {
		return self.rewards
	}

	// glog: MUST be called before applying (reason: before object-deletion)
	self.glog_bag_change()

	// ----- apply -----

	// ccy
	for id, val := range self.Ccy {
		self.bag.add_ccy(id, val)
	}

	// items
	for id, num := range self.Items {
		self.bag.add_item(id, num)
	}

	// hero
	for _, hero := range self.Heroes {
		self.bag.add_hero(hero)
	}

	// heroes-del
	for _, seq := range self.HeroesDel {
		self.bag.del_hero(seq)
	}

	// relic
	for _, rlc := range self.Relics {
		self.bag.add_relic(rlc)
	}

	// relics-del
	for _, seq := range self.RelicsDel {
		self.bag.del_relic(seq)
	}

	// push: bag update
	self.send_update()

	// return rewards
	return self.rewards
}

func (self *BagOp) convert_rewards() {
	// fire
	if len(self.Ccy) > 0 {
		evtmgr.Fire(gconst.Evt_CcyDel, self.bag.plr, (map[int32]int64)(self.Ccy))
	}

	if len(self.Items) > 0 {
		evtmgr.Fire(gconst.Evt_ItemDel, self.bag.plr, (map[int32]int32)(self.Items))
	}

	// ccy
	for id, val := range self.rewards.Ccy {
		self.Ccy[id] += val
	}

	// items
	for id, num := range self.rewards.Items {
		self.Items[id] += num
	}

	// heroes
	for i, id := range self.rewards.HeroIds {
		id := int32(id)

		cur := int64(i + len(self.bag.Heroes))
		max := self.bag.plr.GetCounter().GetMax(gconst.Cnt_HeroBag) +
			self.bag.plr.GetCounter().Get(gconst.Cnt_HeroBagBuy)*5
		if cur < max {
			hero := new_hero(id)
			if hero != nil {
				self.Heroes[hero.Seq] = hero
				self.rewards.HeroIds[i] = hero.Seq
			}
		} else {
			conf := gamedata.ConfMonster.Query(id)
			if conf != nil {
				for _, v := range conf.Fragment {
					self.rewards.Items[v.Id] += v.N
					self.Items[v.Id] += v.N
				}
			}
		}
	}

	// relics
	for i, id := range self.rewards.RelicIds {
		id := int32(id)

		rlc := new_relic(id)
		if rlc != nil {
			self.Relics[rlc.Seq] = rlc
			self.rewards.RelicIds[i] = rlc.Seq
		}
	}

	// fire
	if len(self.rewards.Ccy) > 0 {
		evtmgr.Fire(gconst.Evt_CcyAdd, self.bag.plr, (map[int32]int64)(self.rewards.Ccy), self.from)
	}

	if len(self.rewards.Items) > 0 {
		evtmgr.Fire(gconst.Evt_ItemAdd, self.bag.plr, (map[int32]int32)(self.rewards.Items), self.from)
	}

	if len(self.Heroes) > 0 {
		evtmgr.Fire(gconst.Evt_HeroAdd, self.bag.plr, (map[int64]*Hero)(self.Heroes), self.from)
	}

	if len(self.HeroesDel) > 0 {
		evtmgr.Fire(gconst.Evt_HeroDel, self.bag.plr, ([]int64)(self.HeroesDel))
	}

	if len(self.Relics) > 0 {
		evtmgr.Fire(gconst.Evt_RelicAdd, self.bag.plr, (map[int64]*Relic)(self.Relics), self.from)
	}

}

func (self *BagOp) is_empty() bool {
	return len(self.Ccy) == 0 && len(self.Items) == 0 &&
		len(self.Heroes) == 0 && len(self.HeroesDel) == 0 &&
		len(self.Relics) == 0 && len(self.RelicsDel) == 0
}

func (self *BagOp) glog_bag_change() {
	chg := make(map[int32]int64)

	// ccy
	for id, n := range self.Ccy {
		chg[id] += n
	}

	// items
	for id, n := range self.Items {
		chg[id] += int64(n)
	}

	// heroes
	for _, hero := range self.Heroes {
		chg[hero.Id]++
	}

	// heroes-del
	for _, seq := range self.HeroesDel {
		hero := self.bag.FindHero(seq)
		if hero != nil {
			chg[hero.Id]--
		}
	}

	// relics
	for _, rlc := range self.Relics {
		chg[rlc.Id]++
	}

	// relics-del
	for _, seq := range self.RelicsDel {
		rlc := self.bag.FindRelic(seq)
		if rlc != nil {
			chg[rlc.Id]--
		}
	}

	// fire
	if len(chg) > 0 {
		evtmgr.Fire(gconst.Evt_BagChg, self.bag.plr, chg, self.from)
	}
}

func (self *BagOp) send_update() {
	bag := self.bag
	res := &msg.GS_BagUpdate{}

	// ccy
	for id := range self.Ccy {
		res.Currency = append(res.Currency, &msg.Ccy{
			Id:  id,
			Val: bag.Ccy[id],
		})
	}

	// items
	for id := range self.Items {
		res.Items = append(res.Items, &msg.Item{
			Id:  id,
			Num: bag.Items[id],
		})
	}

	// heroes
	for _, hero := range self.Heroes {
		res.Heroes = append(res.Heroes, hero.ToMsg())
	}
	res.HeroesDel = self.HeroesDel

	// relics
	for _, rlc := range self.Relics {
		res.Relics = append(res.Relics, rlc.ToMsg())
	}
	res.RelicsDel = self.RelicsDel

	// send
	self.bag.plr.SendMsg(res)
}

// ============================================================================

func (self *rewards_t) ToMsg() *msg.Rewards {

	ret := &msg.Rewards{}

	// ccy
	for id, val := range self.Ccy {
		ret.Ccy = append(ret.Ccy, &msg.Ccy{
			Id:  id,
			Val: val,
		})
	}

	// items
	for id, num := range self.Items {
		ret.Items = append(ret.Items, &msg.Item{
			Id:  id,
			Num: num,
		})
	}

	// heroes
	ret.Heroes = self.HeroIds

	// relics
	ret.Relics = self.RelicIds

	return ret
}
