package act

import (
	"fw/src/game/app/gamedata"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"time"
)

// ============================================================================

type IAct interface {
	GetName() string
	set_name(v string)

	get_ver() time.Time
	set_ver(v time.Time)

	GetNum() int32
	set_num(v int32)

	GetStage() string
	set_stage(v string)

	GetT1() time.Time
	set_t1(v time.Time)

	GetT2() time.Time
	set_t2(v time.Time)

	GetActRawData() interface{}
	get_actdata() *act_data_t
	set_actdata(v *act_data_t)

	// --------------------------------

	send_state_change()
	tomsg_state() *msg.ActState

	// --------------------------------

	GetConfGrp() int32

	// --------------------------------

	NewSvrData() interface{}
	NewPlrData() interface{}

	ResetSvrData() // 只能对 svrdata 本身做操作
	ResetPlrData(iplr interface{})

	// --------------------------------

	OnInit()  // 初始化
	OnQuit()  // 停服退出
	OnStage() // 服务器启动后, 除了close阶段, 都会有当前阶段通知过来
}

// ============================================================================

type BaseAct struct {
	name    string      // act name
	ver     time.Time   // current ver
	num     int32       // the nth number
	stage   string      // current stage
	t1      time.Time   // stage start time
	t2      time.Time   // stage end time
	actdata *act_data_t // act data
}

// ============================================================================

func (self *BaseAct) GetName() string {
	return self.name
}

func (self *BaseAct) set_name(name string) {
	self.name = name
}

func (self *BaseAct) get_ver() time.Time {
	return self.ver
}

func (self *BaseAct) set_ver(v time.Time) {
	if !self.ver.Equal(v) {
		self.ver = v
		self.num++

		self.ResetSvrData()
	}
}

func (self *BaseAct) GetNum() int32 {
	return self.num
}

func (self *BaseAct) set_num(v int32) {
	self.num = v
}

func (self *BaseAct) GetStage() string {
	return self.stage
}

func (self *BaseAct) set_stage(v string) {
	self.stage = v
}

func (self *BaseAct) GetT1() time.Time {
	return self.t1
}

func (self *BaseAct) set_t1(v time.Time) {
	self.t1 = v
}

func (self *BaseAct) GetT2() time.Time {
	return self.t2
}

func (self *BaseAct) set_t2(v time.Time) {
	self.t2 = v
}

func (self *BaseAct) GetActRawData() interface{} {
	// reset data if ver differs
	if !self.actdata.Ver.Equal(self.ver) {
		a := FindAct(self.name)
		if a != nil {
			self.actdata.Ver = self.ver
			self.actdata.Data = a.NewSvrData()
		}
	}

	// return
	return self.actdata.Data
}

func (self *BaseAct) get_actdata() *act_data_t {
	return self.actdata
}

func (self *BaseAct) set_actdata(v *act_data_t) {
	self.actdata = v
}

func (self *BaseAct) send_state_change() {
	utils.BroadcastPlayers(&msg.GS_ActStateChange{
		Act: self.tomsg_state(),
	})
}

func (self *BaseAct) tomsg_state() *msg.ActState {
	return &msg.ActState{
		Name:    self.name,
		Stage:   self.stage,
		T1:      self.t1.Unix(),
		T2:      self.t2.Unix(),
		ConfGrp: self.GetConfGrp(),
	}
}

func (self *BaseAct) GetConfGrp() int32 {
	conf := gamedata.ConfActivityM.Query(self.name)
	if conf == nil {
		return 0
	}

	n := self.num
	if n < 1 {
		return 0
	}

	for _, v := range conf.ConfCycles {
		L := int32(len(v.Grps))

		if v.Cnt > 0 {
			m := L * v.Cnt
			if m < n {
				n -= m
				continue
			}
		}

		return v.Grps[(n-1)%L]
	}

	return 0
}

func (self *BaseAct) ResetSvrData()                 {}
func (self *BaseAct) ResetPlrData(iplr interface{}) {}

func (self *BaseAct) OnInit()  {}
func (self *BaseAct) OnQuit()  {}
func (self *BaseAct) OnStage() {}
