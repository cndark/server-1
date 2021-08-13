package wleveldraw

import (
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math/rand"
)

// ============================================================================

// 推图10连
type WLevelDraw struct {
	Id    int32                 // 当前已经领取id
	Items map[int32][]*msg.Item // 当前随出物品

	plr IPlayer
}

// ============================================================================

func NewWLevelDraw() *WLevelDraw {
	return &WLevelDraw{
		Items: make(map[int32][]*msg.Item),
	}
}

func (self *WLevelDraw) Init(plr IPlayer) {
	self.plr = plr
}

func (self *WLevelDraw) Draw(idx int32) (int32, []*msg.Item) {
	if idx < 0 || idx >= 3 {
		return Err.Failed, nil
	}

	if len(self.Items[idx]) > 0 {
		return Err.WLevel_DrawBefore, nil
	}

	conf := gamedata.ConfWorldLevelDraw.Query(self.Id + 1)
	if conf == nil || idx >= int32(len(conf.SeniorDrop)) {
		return Err.Failed, nil
	}

	var items []*msg.Item
	for _, v := range utils.Drop(self.plr, conf.SeniorDrop[idx]) {
		items = append(items, &msg.Item{
			Id:  v.Id,
			Num: int32(v.N),
		})
	}

	for _, v := range utils.Drop(self.plr, conf.NormalDrop) {
		items = append(items, &msg.Item{
			Id:  v.Id,
			Num: int32(v.N),
		})
	}

	if len(items) > 0 {

		rand.Shuffle(len(items), func(i, j int) {
			items[i], items[j] = items[j], items[i]
		})

		self.Items[idx] = items
	}

	return Err.OK, items
}

func (self *WLevelDraw) Take(idx int32, isAutoDec bool) (int32, *msg.Rewards, []*msg.Item) {
	items := self.Items[idx]
	if len(items) == 0 {
		return Err.WLevel_DrawNull, nil, nil
	}

	ditem := make(map[int32]int32)
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_WLevelDraw)
	for _, v := range items {
		// 自动分解3星以下英雄
		b := false
		if isAutoDec && gconst.IsHero(v.Id) {
			conf_m := gamedata.ConfMonster.Query(v.Id)
			if conf_m != nil && conf_m.Star <= 2 {

				conf_s := gamedata.ConfHeroStarUp.Query(conf_m.Star)
				if conf_s != nil {
					for _, s := range conf_s.Sacrifice {
						if s.Type == 1 {
							ditem[s.Id] += s.N * v.Num
						}
					}
					b = true
				}
			}
		}

		if !b {
			op.Inc(v.Id, v.Num)
		}
	}

	// 自动分解的需要整合
	var autoItem []*msg.Item
	for did, dn := range ditem {
		op.Inc(did, dn)
		autoItem = append(autoItem, &msg.Item{Id: did, Num: dn})
	}

	self.Id++
	self.Items = make(map[int32][]*msg.Item)

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds, autoItem
}

// ============================================================================

func (self *WLevelDraw) ToMsg() *msg.WLevelDrawData {
	ret := &msg.WLevelDrawData{
		Id:        self.Id,
		DrawItems: make(map[int32]*msg.WLevelDrawItem),
	}

	for i, v := range self.Items {
		ret.DrawItems[i] = &msg.WLevelDrawItem{
			Items: v,
		}
	}

	return ret
}
