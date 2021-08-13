package heroskin

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/act"
	"fw/src/game/msg"
)

// ============================================================================

var actObj = &act_t{}

// ============================================================================

type act_t struct {
	act.BaseAct
}

type data_svr_t struct {
}

type data_plr_t struct {
	ActGift map[int32]int32 // 活动礼包购买次数
}

// ============================================================================

func init() {
	act.RegisterAct(gconst.ActName_HeroSkin, actObj)

	// 更新购买活动礼包次数
	evtmgr.On(gconst.Evt_ActGift, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		actName := args[1].(string)
		id := args[2].(int32)

		if actObj.GetName() != actName {
			return
		}

		plr_data := actObj.GetPlrData(plr)
		if plr_data == nil {
			return
		}

		plr_data.ActGift[id]++
	})
}

// ============================================================================

func (self *act_t) NewSvrData() interface{} {
	return new(data_svr_t)
}

func (self *act_t) NewPlrData() interface{} {
	return &data_plr_t{
		ActGift: make(map[int32]int32),
	}
}

func (self *act_t) ResetSvrData() {
}

func (self *act_t) ResetPlrData(iplr interface{}) {
}

func (self *act_t) GetSvrData() *data_svr_t {
	return self.GetActRawData().(*data_svr_t)
}

func (self *act_t) GetPlrData(plr IPlayer) *data_plr_t {
	return plr.GetActRawData(self.GetName()).(*data_plr_t)
}

// ============================================================================

func (self *act_t) Started() bool {
	return self.GetStage() == "start"
}

func (self *act_t) Ended() bool {
	return self.GetStage() == "end"
}

func (self *act_t) Closed() bool {
	return self.GetStage() == "close"
}

// ============================================================================

func (self *act_t) OnInit() {
}

func (self *act_t) OnQuit() {
}

func (self *act_t) OnStage() {
}

func (self *act_t) ToMsg(plr IPlayer) *msg.ActHeroSkinData {

	plr_data := self.GetPlrData(plr)
	if plr_data == nil {
		return nil
	}

	ret := &msg.ActHeroSkinData{
		ActGift: plr_data.ActGift,
	}

	return ret
}

// ============================================================================

func ActHeroSkinInfo(plr IPlayer) *msg.ActHeroSkinData {
	return actObj.ToMsg(plr)
}
