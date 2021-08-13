package tutorial

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
)

// ============================================================================

type Tutorial struct {
	Steps map[string]int32 // 引导[type]step
	Keys  map[string]bool  // used keys

	plr IPlayer
}

// ============================================================================

func NewTutorial() *Tutorial {
	return &Tutorial{
		Steps: make(map[string]int32),
		Keys:  make(map[string]bool),
	}
}

func (self *Tutorial) Init(plr IPlayer) {
	self.plr = plr
}

func (self *Tutorial) HasKey(key string) bool {
	return self.Keys[key]
}

func (self *Tutorial) Set(tp string, step int32, key string) {
	if tp == "" || len(tp) > 32 {
		return
	}

	self.Steps[tp] = step

	if key != "" {
		self.Keys[key] = true
	}

	// fire
	evtmgr.Fire(gconst.Evt_Tutorial, self.plr, tp, step)
}

func (self *Tutorial) ToMsg() *msg.TutorialData {
	return &msg.TutorialData{
		Steps: self.Steps,
	}
}
