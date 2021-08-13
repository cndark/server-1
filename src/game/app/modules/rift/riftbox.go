package rift

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/sched/loop"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/mail"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

// ============================================================================

// 宝箱
type rift_box_t struct {
	Boxes []*box_t // 宝箱
}

type box_t struct {
	Id     int32     // 宝箱id
	CurPlr string    // 当前占领玩家
	FinTs  time.Time // 完成时间

	isBattle  bool
	tid_award *core.Timer // 发奖timer
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_GlobalResetDaily, func(args ...interface{}) {
		BoxMgr.daily_reset()
	})
}

// ============================================================================

func (self *rift_box_t) data_loaded() {
	for _, box := range self.Boxes {
		if box.CurPlr == "" {
			continue
		}

		box := box
		box.tid_award = loop.SetTimeout(box.FinTs, func() {
			plr := load_player(box.CurPlr)
			if plr != nil {
				self.award_box(plr, box.Id)
			}
		})
	}
}

func (self *rift_box_t) daily_reset() {
	for _, box := range self.Boxes {
		box.cancel_timer()
	}

	cnt := 30
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf != nil {
		cnt = int(conf.RiftBoxNum)
	}

	self.Boxes = []*box_t{}
	for i := 0; i < cnt; i++ {
		self.Boxes = append(self.Boxes, &box_t{Id: int32(i + 1)})
	}
}

func (self *rift_box_t) find_box(id int32) *box_t {
	for _, box := range self.Boxes {
		if box.Id == id {
			return box
		}
	}

	return nil
}

// ============================================================================

// 探索裂隙宝箱
func (self *rift_box_t) ExploreRiftBox(plr IPlayer) {
	if !self.IsStart() {
		return
	}

	L := len(self.Boxes)
	if L == 0 {
		return
	}

	idx := rand_rift.Intn(L)
	box := self.Boxes[idx]

	plr.SendMsg(&msg.GS_RiftBoxNew{
		Box: box.ToMsg(),
	})
}

// 是否是开始状态
func (self *rift_box_t) IsStart() bool {
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return false
	}

	now := time.Now()
	zero := core.StartOfDay(now)
	begin := zero.Add(time.Duration(conf.RiftBoxRefresh) * time.Hour)
	end := begin.Add(time.Duration(conf.RiftBoxLifeMinute) * time.Minute)

	return now.After(begin) && now.Before(end)
}

// 占领宝箱
func (self *rift_box_t) Occupy(plr IPlayer, id int32,
	cb func(ec int32, finTs int64, replay *msg.BattleReplay)) {

	if !self.IsStart() {
		cb(Err.Rift_BoxIsEnd, 0, nil)
		return
	}

	if !plr.IsSetTeam(gconst.TeamType_Dfd) {
		cb(Err.Common_NotSetTeam, 0, nil)
		return
	}

	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil {
		cb(Err.Failed, 0, nil)
		return
	}

	box := self.find_box(id)
	if box == nil || box.isBattle {
		cb(Err.Rift_BoxNotFound, 0, nil)
		return
	}

	if box.CurPlr == plr.GetId() {
		cb(Err.Rift_BoxIsSelf, 0, nil)
		return
	}

	if box.isBattle {
		cb(Err.Rift_BoxIsBattle, 0, nil)
		return
	}

	// 空的
	if box.CurPlr == "" {
		box.CurPlr = plr.GetId()
		box.FinTs = time.Now().Add(time.Duration(conf_g.RiftBoxOpenMinute) * time.Minute)
		box.tid_award = loop.SetTimeout(box.FinTs, func() {
			self.award_box(plr, id)
		})

		cb(Err.OK, box.FinTs.Unix(), nil)
		return
	}

	// 抢别人的
	cop := plr.GetCounter().NewOp(gconst.ObjFrom_RiftBox)
	cop.DecCounter(gconst.Cnt_PlayerStrength, int64(conf_g.RiftCost[C_RiftType_Box]))
	if ec := cop.CheckEnough(); ec != Err.OK {
		cb(ec, 0, nil)
		return
	}

	// 进攻
	eplr := load_player(box.CurPlr)
	if eplr == nil {
		cb(Err.Plr_NotFound, 0, nil)
		return
	}

	input := &msg.BattleInput{
		T1: plr.ToMsg_BattleTeam(plr.GetTeam(gconst.TeamType_Dfd)),
		T2: eplr.ToMsg_BattleTeam(eplr.GetTeam(gconst.TeamType_Dfd)),
		Args: map[string]string{
			"Module":    "RIFT_BOX",
			"RoundType": "2",
		},
	}

	box.isBattle = true

	// fight
	battle.Fight(input, func(r *msg.BattleResult) {
		box.isBattle = false

		if r == nil {
			cb(Err.Common_BattleResError, 0, nil)
			return
		}

		cop := plr.GetCounter().NewOp(gconst.ObjFrom_RiftBox)
		cop.DecCounter(gconst.Cnt_PlayerStrength, int64(conf_g.RiftCost[C_RiftType_Box]))
		if ec := cop.CheckEnough(); ec != Err.OK {
			cb(ec, 0, nil)
			return
		}

		if r.Winner == 1 {
			box.cancel_timer()
			box.CurPlr = plr.GetId()
			box.FinTs = time.Now().Add(time.Duration(conf_g.RiftBoxOpenMinute) * time.Minute)
			box.tid_award = loop.SetTimeout(box.FinTs, func() {
				self.award_box(plr, id)
			})

			// notify
			eplr.SendMsg(&msg.GS_RiftBoxOccupied{
				Box: box.ToMsg(),
			})
		}

		cop.Apply()

		// res
		replay := &msg.BattleReplay{
			Ts: time.Now().Unix(),
			Bi: input,
			Br: r,
		}

		cb(Err.OK, box.FinTs.Unix(), replay)
	})
}

func (self *rift_box_t) award_box(plr IPlayer, id int32) {
	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil {
		return
	}

	conf_lv := gamedata.ConfPlayerUp.Query(plr.GetLevel())
	if conf_lv == nil {
		return
	}

	// mail
	conf_m := gamedata.ConfMail.Query(conf_g.RiftBoxMail)
	if conf_m == nil {
		return
	}

	m := mail.New(plr).SetKey(conf_g.RiftBoxMail)
	for _, v := range conf_g.RiftBoxReward {
		if rand_rift.Float32() < v.Odds {
			a := float32(1)
			for _, vv := range conf_lv.RiftBoxRatio {
				if vv.Id == v.Id {
					a = vv.N
				}
			}
			n := int64(float32(v.N) * a)
			if n > 0 {
				m.AddAttachment(v.Id, float64(n))
			}
		}
	}

	m.Send()

	plr.SendMsg(&msg.GS_RiftBoxRewards{
		Id: id,
	})

	// evt
	evtmgr.Fire(gconst.Evt_RiftBoxTake, plr)

	// delete
	for i, v := range self.Boxes {
		if v.Id == id {
			self.Boxes = append(self.Boxes[:i], self.Boxes[i+1:]...)
			break
		}
	}
}

func (self *rift_box_t) ToMsg(plr IPlayer) (box *msg.RiftBox, num int32) {
	num = int32(len(self.Boxes))
	for _, v := range self.Boxes {
		if v.CurPlr == plr.GetId() {
			box = v.ToMsg()
		}
	}

	return
}

// ============================================================================

func (self *box_t) cancel_timer() {
	if self.tid_award != nil {
		loop.CancelTimer(self.tid_award)
		self.tid_award = nil
	}
}

func (self *box_t) ToMsg() *msg.RiftBox {
	ret := &msg.RiftBox{
		Id:    self.Id,
		FinTs: self.FinTs.Unix(),
	}

	plr := load_player(self.CurPlr)
	if plr != nil {
		ret.CurPlr = plr.ToMsg_SimpleInfo()

		T := plr.GetTeam(gconst.TeamType_Dfd)
		if T != nil {
			ret.Fighters = plr.ToMsg_BattleTeam(T).Fighters
		}
	}

	return ret
}
