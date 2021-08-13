package shop

import (
	"fw/src/game/app/gamedata"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math/rand"
	"time"
)

// ============================================================================

var rand_shop = rand.New(rand.NewSource(time.Now().Unix()))

// ============================================================================

type Shop struct {
	Shops shop_m //所有商店

	plr IPlayer
}

type shop_m map[int32]*shop_t
type shop_t struct {
	Goods  []*goods_t //拥有的商品
	LastTs time.Time  //上次自动刷新时间
}

//单个商品结构
type goods_t struct {
	Id     int32 //商品表Id
	BuyCnt int32 //已经购买次数
}

// ============================================================================

func NewShop() *Shop {
	return &Shop{
		Shops: make(shop_m),
	}
}

func (self *Shop) Init(plr IPlayer) {
	self.plr = plr

	self.init_shop()
}

// 商店初始化
func (self *Shop) init_shop() {
	for _, v := range gamedata.ConfShop.Items() {
		_, ok := self.Shops[v.Id]
		if !ok {
			self.Shops[v.Id] = &shop_t{}
			self.Shops[v.Id].new_shop(self.plr, v.Id)
		}
	}

}

func (self *Shop) IsShopOpen(shopid int32) int32 {
	conf := gamedata.ConfShop.Query(shopid)
	if conf == nil {
		return Err.Failed
	}

	return utils.IsStatusTabArrayConform(self.plr, conf.OpenStatus, true)
}

func (self *Shop) Update(shopid int32) {
	shop := self.Shops[shopid]
	if shop == nil {
		return
	}

	if self.is_need_update(shopid) {
		shop.update_shop(self.plr, shopid)
	} else {
		self.check_goods_conf(shopid)
	}
}

// 检查商品更新
func (self *Shop) is_need_update(shopid int32) bool {
	conf := gamedata.ConfShop.Query(shopid)
	if conf == nil || len(conf.RefreshTime) == 0 {
		return false
	}

	shop := self.Shops[shopid]
	if conf.IsRandom == 1 && len(shop.Goods) == 0 {
		return true
	}

	now := time.Now()
	y, M, d := now.Date()

	for _, v := range conf.RefreshTime {
		if v >= 24 || v < -1 {
			return false
		}

		var key time.Time
		if v >= 0 {
			key = time.Date(y, M, d, int(v), 0, 0, 0, time.Local)
		} else {
			key = time.Date(y, M, -int(v), 0, 0, 0, 0, time.Local)
		}

		if shop.LastTs.Before(key) && now.After(key) {
			return true
		}
	}

	return false
}

// 检查更新商品
func (self *Shop) check_goods_conf(shopid int32) {
	shop := self.Shops[shopid]
	if shop == nil {
		return
	}

	for i := 0; i < len(shop.Goods); i++ {
		conf := gamedata.ConfShopItem.Query(shop.Goods[i].Id)
		if conf == nil ||
			(len(conf.ItemClose) != 0 &&
				utils.IsStatusTabArrayConform(self.plr, conf.ItemClose) == Err.OK) {

			shop.Goods = append(shop.Goods[:i], shop.Goods[i+1:]...)
			i--
			continue
		}
	}
}

// 刷新
func (self *Shop) Refresh(shopid int32) int32 {
	shop := self.Shops[shopid]
	if shop == nil {
		return Err.Shop_NotFound
	}

	shop.refresh_shop(self.plr, shopid)

	return Err.OK
}

// 购买商品
func (self *Shop) CheckBuy(shopid int32, itemid int32) int32 {
	conf := gamedata.ConfShopItem.Query(itemid)
	if conf == nil {
		return Err.Failed
	}

	return utils.IsStatusTabArrayConform(self.plr, conf.BlankOpen)
}

func (self *Shop) BuyGoods(shopid int32, itemid, num int32) {
	i := self.find_goods_idx(shopid, itemid)
	if i < 0 {
		self.Shops[shopid].Goods = append(self.Shops[shopid].Goods, &goods_t{itemid, num})
		return
	}

	self.Shops[shopid].Goods[i].BuyCnt += num
}

func (self *Shop) GetGoodBuyCnt(shopid int32, itemid int32) int32 {
	i := self.find_goods_idx(shopid, itemid)
	if i < 0 {
		return 0
	}

	return self.Shops[shopid].Goods[i].BuyCnt
}

func (self *Shop) ModifyGoodBuyCnt(shopid int32, itemid int32, n int32) int32 {
	conf := gamedata.ConfShopItem.Query(itemid)
	if conf == nil {
		return 0
	}

	if conf.BuyCount == -1 {
		return n
	}

	r := conf.BuyCount - self.GetGoodBuyCnt(shopid, itemid)
	if n > r {
		return r
	}

	return n
}

func (self *Shop) find_goods_idx(shopid, itemid int32) int32 {
	for i, v := range self.Shops[shopid].Goods {
		if itemid == v.Id {
			return int32(i)
		}
	}
	return -1
}

// ============================================================================

func (self *Shop) ToMsg(shopid int32) *msg.ShopData {
	shop := self.Shops[shopid]
	if shop == nil {
		return nil
	}

	ret := &msg.ShopData{Id: shopid}
	for _, v := range shop.Goods {
		ret.Goods = append(ret.Goods, &msg.Good{
			Id:     v.Id,
			BuyCnt: v.BuyCnt,
		})
	}

	return ret
}

// ============================================================================

func (self *shop_t) new_shop(plr IPlayer, shopid int32) {
	self.Goods = gen_goods(plr, shopid)
	self.LastTs = time.Now()
}

func (self *shop_t) update_shop(plr IPlayer, shopid int32) {
	self.Goods = gen_goods(plr, shopid)
	self.LastTs = time.Now()
}

func (self *shop_t) refresh_shop(plr IPlayer, shopid int32) {
	self.Goods = gen_goods(plr, shopid)
}

func gen_goods(plr IPlayer, shopid int32) (ret []*goods_t) {
	ret = []*goods_t{}

	conf_s := gamedata.ConfShop.Query(shopid)
	if conf_s == nil || conf_s.IsRandom != 1 {
		return
	}

	conf := gamedata.ConfShopItemM.QueryS(shopid)
	if len(conf) == 0 {
		return
	}

	//商品不符合开启条件的，不进入随机数组
	choose := gamedata.ConfShopItemM.NewShopItemToBlank()
	for k, v := range conf {
		for i, g := range v {
			if len(g.ItemClose) != 0 && utils.IsStatusTabArrayConform(plr, g.ItemClose) == Err.OK {
				continue
			}

			if len(g.ItemOpen) != 0 && utils.IsStatusTabArrayConform(plr, g.ItemOpen) != Err.OK {
				continue
			}

			choose[k] = append(choose[k], v[i])
		}
	}

	//数组里面随机
	for _, b := range choose {
		if len(b) == 0 {
			continue
		}

		slt := make(map[int32]int32)
		for _, v := range b { //total weight
			slt[v.Id] += v.Odds
		}

		id := utils.PickWeightedMapId(slt)
		if id == 0 {
			continue
		}

		ret = append(ret, &goods_t{Id: id})
	}

	return
}
