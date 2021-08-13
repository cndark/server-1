package signdaily

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

// 日常签到
type SignDaily struct {
	CanDays int32 // 签到未领取天数
	Day     int32 // 最近领取的是第几条
	Round   int32 // 当前领取的奖励的轮数

	LastTs time.Time // 最近领取时间
	plr    IPlayer
}

// ============================================================================

func NewSignDaily() *SignDaily {
	return &SignDaily{
		Day:   0,
		Round: 1,
	}
}

func (self *SignDaily) Init(plr IPlayer) {
	self.plr = plr
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_PlrDailyOnline, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		plr_data := plr.GetSignDaily()
		if plr_data == nil {
			return
		}

		plr_data.CanDays++

		if plr_data.Day >= int32(len(gamedata.ConfSignDaily.Items())) {
			plr_data.Day = 0
			plr_data.Round += 1
		}

	})

}

// ============================================================================

func (self *SignDaily) Sign() (int32, *msg.Rewards) {

	now := time.Now()
	if core.IsSameDay(now, self.LastTs) {
		return Err.Plr_TakenBefore, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_SignDaily)

	tDay := self.Day
	tCanDays := self.CanDays

	for i := int32(1); i <= self.CanDays; i++ {
		conf := gamedata.ConfSignDaily.Query(i + self.Day)
		if conf == nil {
			continue
		}

		L := len(conf.Rewards)
		if L > 1 {
			idx := int(self.Round) % L
			if idx == 0 {
				idx = L
			}

			op.Inc(conf.Rewards[idx-1].Id, conf.Rewards[idx-1].N)
		} else {
			op.Inc(conf.Rewards[0].Id, conf.Rewards[0].N)
		}

		tDay = conf.Day
		tCanDays--
	}

	self.Day = tDay
	self.CanDays = tCanDays
	self.LastTs = now

	rwds := op.Apply().ToMsg()

	// fire
	evtmgr.Fire(gconst.Evt_SignDaily, self.plr, self.Day)

	return Err.OK, rwds
}

// ============================================================================

func (self *SignDaily) ToMsg() *msg.SignDailyData {
	return &msg.SignDailyData{
		CanDays: self.CanDays,
		Day:     self.Day,
		Round:   self.Round,
	}
}
