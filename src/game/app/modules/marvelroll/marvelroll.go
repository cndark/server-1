package marvelroll

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

// ============================================================================

// 奇迹转盘
type MarvelRoll struct {
	Groups map[string]*group_t

	plr IPlayer
}

type group_t struct {
	EnforceTs time.Time  // 强制刷新
	FreeTs    time.Time  // 上次免费刷新时间
	Blanks    []*blank_t // 格子
}

type blank_t struct {
	Id      int32
	AttainN int32 // 获得次数
}

// ============================================================================

func NewMarvelRoll() *MarvelRoll {
	return &MarvelRoll{
		Groups: make(map[string]*group_t),
	}
}

// ============================================================================

func (self *MarvelRoll) Init(plr IPlayer) {
	self.plr = plr

	for _, conf := range gamedata.ConfMarvelRoll.Items() {
		group := self.Groups[conf.Group]
		if group == nil {
			self.Groups[conf.Group] = &group_t{
				EnforceTs: time.Unix(0, 0),
				FreeTs:    time.Unix(0, 0),
			}
		}
	}
}

func (self *MarvelRoll) refresh(grp string) {
	conf := gamedata.ConfMarvelRoll.Query(grp)
	if conf == nil {
		return
	}

	group := self.Groups[grp]
	if group == nil {
		return
	}

	blanks := []*blank_t{}
	items := gamedata.ConfMarvelRollOddsM.QueryGrp(grp)
	for i := 1; i <= int(conf.BlankNum); i++ {
		b_items := items[int32(i)]

		slt := make(map[int32]int32)
		for _, v := range b_items {
			lv := self.plr.GetLevel()
			if len(v.LevelRange) > 0 && (lv < v.LevelRange[0].Low || lv > v.LevelRange[0].High) {
				continue
			}

			slt[v.Id] += v.RefOdds
		}

		id := utils.PickWeightedMapId(slt)
		if id == 0 {
			continue
		}

		blanks = append(blanks, &blank_t{Id: id})
	}

	group.Blanks = blanks
}

// 检查强制刷新
func (self *MarvelRoll) CheckRefresh() {
	now := time.Now()
	for grp, group := range self.Groups {
		if !core.IsSameDay(group.EnforceTs, now) {
			self.refresh(grp)
			group.EnforceTs = now
		}
	}
}

// 刷新
func (self *MarvelRoll) Refresh(grp string) int32 {
	conf := gamedata.ConfMarvelRoll.Query(grp)
	if conf == nil {
		return Err.Failed
	}

	group := self.Groups[grp]
	if group == nil {
		return Err.MarvelRoll_GroupNotFound
	}

	now := time.Now()
	// free or cost
	if core.IsSameDay(group.FreeTs, now) {
		op := self.plr.GetBag().NewOp(gconst.ObjFrom_MarvelRollRefresh)
		for _, v := range conf.RefreshCost {
			op.Dec(v.Id, v.N)
		}

		if ec := op.CheckEnough(); ec != Err.OK {
			return ec
		}

		op.Apply()
	}

	self.refresh(grp)

	group.FreeTs = now

	return Err.OK
}

func (self *MarvelRoll) Take(grp string, isTen bool) (int32, []int32) {
	conf := gamedata.ConfMarvelRoll.Query(grp)
	if conf == nil {
		return Err.Failed, nil
	}

	if ec := utils.IsStatusTabArrayConform(self.plr, conf.OpenStatus); ec != Err.OK {
		return ec, nil
	}

	group := self.Groups[grp]
	if group == nil {
		return Err.MarvelRoll_GroupNotFound, nil
	}

	// check cost
	costR, cnt := float32(1), 1
	if isTen {
		costR = conf.Discount * 10
		cnt = 10
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_MarvelRollTake)
	for _, v := range conf.Cost {
		op.Dec(v.Id, int64(float32(v.N)*costR))
	}

	if ec := op.CheckEnough(); ec != Err.OK {
		return ec, nil
	}

	// roll
	lastId := []int32{}
	for i := 0; i < cnt; i++ {
		id := self.roll_one(group)
		conf_o := gamedata.ConfMarvelRollOdds.Query(id)
		if conf_o == nil {
			continue
		}

		for _, v := range conf_o.Item {
			op.Inc(v.Id, v.N)
		}

		for _, v := range conf.ExtReward {
			op.Inc(v.Id, v.N)
		}

		for _, v := range group.Blanks {
			if v.Id == id {
				v.AttainN++
			}
		}

		// 最后有效id返回给客户端
		lastId = append(lastId, id)
	}

	if len(lastId) == 0 {
		return Err.MarvelRoll_RollNothing, nil
	}

	op.Apply()

	evtmgr.Fire(gconst.Evt_MarvelRoll, self.plr, grp, conf.ModuleId, int32(cnt))

	return Err.OK, lastId
}

func (self *MarvelRoll) roll_one(group *group_t) int32 {
	slt := make(map[int32]int32)
	for _, v := range group.Blanks {
		conf := gamedata.ConfMarvelRollOdds.Query(v.Id)
		if conf == nil {
			continue
		}

		if conf.Upper > 0 && v.AttainN >= conf.Upper {
			continue
		}

		slt[v.Id] += conf.RollOdds
	}

	return utils.PickWeightedMapId(slt)
}

// ============================================================================

func (self *MarvelRoll) ToMsg() *msg.MarvelRollData {
	ret := &msg.MarvelRollData{}
	for grp, group := range self.Groups {
		ret.Groups = append(ret.Groups, group.ToMsg_Group(grp))
	}

	return ret
}

func (self *group_t) ToMsg_Group(grp string) *msg.MarvelRollGroup {
	now := time.Now()

	ret := &msg.MarvelRollGroup{
		Grp:       grp,
		IsEnforce: !core.IsSameDay(self.EnforceTs, now),
		IsFree:    !core.IsSameDay(self.FreeTs, now),
	}

	for _, v := range self.Blanks {
		ret.Blanks = append(ret.Blanks, &msg.MarvelRollBlank{
			Id:      v.Id,
			AttainN: v.AttainN,
		})
	}

	return ret
}
