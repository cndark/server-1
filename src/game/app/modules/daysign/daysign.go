package daysign

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================
// 七日之约
type DaySign struct {
	SignDay int32   // 登录天数
	Taken   []int32 // 领取的天数
	Close   bool    // 活动是否结束

	plr IPlayer
}

// ============================================================================

func NewDaySign() *DaySign {
	return &DaySign{}
}

func init() {
	evtmgr.On(gconst.Evt_PlrDailyOnline, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		plr_data := plr.GetDaySign()
		if plr_data == nil {
			return
		}

		if !plr_data.Close {
			plr_data.SignDay++
		}
	})

}

// ============================================================================

func (self *DaySign) Init(plr IPlayer) {
	self.plr = plr
}

func (self *DaySign) Take(id int32) (int32, *msg.Rewards) {

	for _, v := range self.Taken {
		if id == v {
			return Err.DaySign_AlreadSigned, nil
		}
	}

	conf := gamedata.ConfDaySign.Query(id)
	if conf == nil {
		return Err.DaySign_NotFound, nil
	}

	if conf.Day > self.SignDay {
		return Err.DaySign_NotCond, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_DaySign)

	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	rwds := op.Apply().ToMsg()

	self.Taken = append(self.Taken, id)

	if id >= int32(len(gamedata.ConfDaySign.Items())) {
		self.Close = true
	}

	return Err.OK, rwds
}

func (self *DaySign) ToMsg() *msg.DaySignData {
	return &msg.DaySignData{
		Taken:   self.Taken,
		SignDay: self.SignDay,
		Close:   self.Close,
	}
}
