package act0

import (
	"fw/src/game/app/comp"
	"fw/src/game/app/modules/act"
)

// ============================================================================

type IPlayer interface {
	comp.IPlayer
}

// ============================================================================

var act0 = &act_t{}

// ============================================================================

type act_t struct {
	act.BaseAct
}

type data_svr_t struct {
}

type data_plr_t struct {
}

// ============================================================================

func init() {
	act.RegisterAct("act0", act0)
}

// ============================================================================

func (self *act_t) NewSvrData() interface{} {
	return new(data_svr_t)
}

func (self *act_t) NewPlrData() interface{} {
	return &data_plr_t{}
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

// ============================================================================
