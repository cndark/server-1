package misc

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/math"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/worlddata"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
	"math/rand"
	"strings"
	"time"
)

// ============================================================================

type Misc struct {
	FreeRename bool // 是否免费改名

	GldLeaveTs  time.Time     // 玩家离开家族,可再次加入时间
	GldLeaveCnt int32         // 离开家族次数
	GldActGift  []*act_gift_t // 家族他人分享活动礼包

	// 点金手
	GoldenHandCrit float32

	OnlineBoxId  int32 // 上次领取宝箱ID
	OnlineBoxDur int64 // 上次领取到现在的时长

	plr IPlayer
}

// 分享活动礼包
type act_gift_t struct {
	Id   int32
	Name string
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_GuildLeave, func(args ...interface{}) {
		plr := args[1].(IPlayer)

		conf := gamedata.ConfGlobalPublic.Query(1)
		if conf == nil {
			return
		}

		L := len(conf.GuildJoinCD)
		if L == 0 {
			return
		}

		data := plr.GetMisc()

		ts := int32(0)
		ct := int(data.GldLeaveCnt)
		if L <= ct {
			ts = conf.GuildJoinCD[L-1]
		} else {
			ts = conf.GuildJoinCD[ct]
		}

		data.GldLeaveTs = time.Now().Add(time.Duration(ts) * time.Minute)
		data.GldLeaveCnt++

		plr.SendMsg(&msg.GS_GuildPlrLeaveTs{
			LeaveTs:  data.GldLeaveTs.Unix(),
			LeaveCnt: data.GldLeaveCnt,
		})
	})

	evtmgr.On(gconst.Evt_PlrResetDaily, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		plr.GetMisc().reset_daily()
	})
}

func (self *Misc) reset_daily() {
	self.GoldenHandCrit = 0
}

// ============================================================================

func NewMisc() *Misc {
	return &Misc{
		FreeRename: true,

		GldLeaveTs:  time.Unix(0, 0),
		GldLeaveCnt: 0,
	}
}

func (self *Misc) Init(plr IPlayer) {
	self.plr = plr
}

func (self *Misc) GiftExchange(code string) (ec int32, rwd *msg.Rewards) {
	// get info
	info := gift_get_info(code)
	if info == nil {
		return Err.Gift_NoCode, nil
	}

	// check area
	if info.Area != -1 && info.Area != config.Common.Area.Id {
		return Err.Gift_AreaLimit, nil
	}

	// check expire
	if time.Now().After(info.Expire) {
		return Err.Gift_CodeExpired, nil
	}

	// check use
	if gift_code_used(info, code, self.plr.GetId()) {
		return Err.Gift_CodeUsed, nil
	}

	// ok. give rewards
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_GiftCode)
	for _, v := range info.Rewards {
		p := strings.Split(v.ResK, " - ")
		if len(p) != 2 {
			continue
		}

		id := core.Atoi32(p[1])
		n := v.ResV

		if id == 0 || n == 0 {
			continue
		}

		op.Inc(id, n)
	}
	rwd = op.Apply().ToMsg()

	// update use
	gift_update_use(info, code, self.plr.GetId())

	// return
	return Err.OK, rwd
}

// 获得家族活动礼包
func (self *Misc) GldActGiftAdd(id int32, name string) {
	self.GldActGift = append(self.GldActGift, &act_gift_t{
		Id:   id,
		Name: name,
	})

	// res
	self.plr.SendMsg(&msg.GS_MiscGldActGift{
		One: &msg.MiscGldActGiftOne{
			Id:   id,
			Name: name,
		},
	})
}

// 领取家族活动礼包
func (self *Misc) GldActGiftTake(idx int32) (int32, *msg.Rewards) {
	if idx >= int32(len(self.GldActGift)) || idx < 0 {
		return Err.Plr_TakenBefore, nil
	}

	conf := gamedata.ConfActGift.Query(self.GldActGift[idx].Id)
	if conf == nil || conf.Type != gconst.C_ActGift_Gld {
		return Err.Failed, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_ActGift)
	for _, v := range conf.GuildSharedReward {
		op.Inc(v.Id, v.N)
	}

	self.GldActGift = append(self.GldActGift[:idx], self.GldActGift[idx+1:]...)

	rwds := op.Apply().ToMsg()

	return Err.OK, rwds
}

func (self *Misc) ToMsg() *msg.MiscData {
	ret := &msg.MiscData{
		FreeRename:     self.FreeRename,
		SvrOpenTs:      worlddata.GetSvrCreateTs().Unix(),
		GldLeaveTs:     self.GldLeaveTs.Unix(),
		GldLeaveCnt:    self.GldLeaveCnt,
		GoldenHandCrit: self.GoldenHandCrit,
		OnlineBoxId:    self.OnlineBoxId,
		OnlineBoxDur:   self.OnlineBoxDur,
	}

	for _, v := range self.GldActGift {
		ret.GldActGift = append(ret.GldActGift, &msg.MiscGldActGiftOne{
			Id:   v.Id,
			Name: v.Name,
		})
	}

	return ret
}

// ============================================================================
// 点金手

func (self *Misc) GoldenHandClick() (ec int32, rwd *msg.Rewards) {
	// conf
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return Err.Failed, nil
	}

	// check counter
	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_GoldenHand)
	cop.DecCounter(gconst.Cnt_GoldenHand, 1)
	if ec := cop.CheckEnough(); ec != Err.OK {
		return ec, nil
	}

	// ok. add gold
	b := false
	d := conf.BuyGoldBasic + (self.plr.GetLevel()-1)*conf.BuyGoldLevelRatio
	v := math.MinInt32(d, conf.BuyGoldLimit)
	if rand.Float32() < self.GoldenHandCrit {
		v *= 2
		b = true
	}

	cop.Inc(gconst.Gold, v)

	// apply
	rwd = cop.Apply().ToMsg()

	// update next crit
	if b {
		self.GoldenHandCrit = 0
	} else {
		var idx int
		if self.GoldenHandCrit == 0 {
			idx = 0
		} else {
			idx = 1
		}

		add_crit := core.RandFloat32(conf.BuyGoldCritAdd[idx].Min, conf.BuyGoldCritAdd[idx].Max)
		self.GoldenHandCrit += add_crit
	}

	evtmgr.Fire(gconst.Evt_GoldHand, self.plr)

	return Err.OK, rwd
}

// ============================================================================
// 在线宝箱领奖

func (self *Misc) TakeOnlineBox(id int32) (int32, *msg.Rewards) {

	if id <= self.OnlineBoxId {
		return Err.OnlineBox_AlreadRewarded, nil
	}

	conf := gamedata.ConfOnlineBox.Query(id)
	if conf == nil {
		return Err.Failed, nil
	}

	d1 := int64(self.plr.AccOnlineDur())

	if int64(d1-self.OnlineBoxDur) < int64(conf.OnlineTime) {
		return Err.OnlineBox_NotCompleted, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_OnlineBox)

	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	rwds := op.Apply().ToMsg()

	self.OnlineBoxId = id
	self.OnlineBoxDur = d1

	return Err.OK, rwds
}
