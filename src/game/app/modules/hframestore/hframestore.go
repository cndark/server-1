package hframestore

import (
	"fw/src/core/sched/loop"
	"fw/src/game/app/gamedata"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

// ============================================================================

// 头像框
type HFrameStore struct {
	HFrames map[int32]*hframe_t

	plr IPlayer
}

type hframe_t struct {
	ExpireTs time.Time // 过期时间
}

// ============================================================================

func NewHFrameStore() *HFrameStore {
	ret := &HFrameStore{
		HFrames: make(map[int32]*hframe_t),
	}

	ret.HFrames[0] = &hframe_t{}

	return ret
}

// ============================================================================

func (self *HFrameStore) Init(plr IPlayer) {
	self.plr = plr

	// check expire
	now := time.Now()
	for id, hf := range self.HFrames {
		conf := gamedata.ConfPlayerHFrame.Query(id)
		if conf != nil && conf.Time != 0 {
			if now.After(hf.ExpireTs) {
				if self.plr.GetHFrame() == id {
					self.plr.SetHFrame(0)
				}

				delete(self.HFrames, id)
			} else {
				loop.SetTimeout(hf.ExpireTs, func() {
					if self.plr.GetHFrame() == id {
						self.plr.SetHFrame(0)
					}

					delete(self.HFrames, id)
				})
			}
		}
	}
}

func (self *HFrameStore) IsValid(id int32) int32 {
	hf := self.HFrames[id]
	if hf == nil {
		return Err.Plr_HFrameError
	}

	conf := gamedata.ConfPlayerHFrame.Query(id)
	if conf == nil {
		return Err.Failed
	}

	if conf.Time != 0 && time.Now().After(hf.ExpireTs) {
		return Err.Plr_HFrameError
	}

	return Err.OK
}

func (self *HFrameStore) Add(id int32) int32 {
	conf := gamedata.ConfPlayerHFrame.Query(id)
	if conf == nil {
		return Err.Failed
	}

	hf := self.HFrames[id]
	if hf == nil {
		hf = &hframe_t{ExpireTs: time.Unix(0, 0)}
		self.HFrames[id] = hf
	}

	now := time.Now()
	if conf.Time != 0 {
		dur := time.Duration(conf.Time*24) * time.Hour

		if now.After(hf.ExpireTs) {
			hf.ExpireTs = now.Add(dur)
		} else {
			hf.ExpireTs = hf.ExpireTs.Add(dur)
		}

		loop.SetTimeout(hf.ExpireTs, func() {
			self.plr.SetHFrame(0)
			delete(self.HFrames, id)
		})
	}

	self.plr.SendMsg(&msg.GS_PlayerUpdateHFrame{
		One: &msg.HFrameStoreOne{
			Id:       id,
			ExpireTs: hf.ExpireTs.Unix(),
		},
	})

	return Err.OK
}

// ============================================================================

func (self *HFrameStore) ToMsg() *msg.HFrameStoreData {
	ret := &msg.HFrameStoreData{}
	for id, hf := range self.HFrames {
		ret.HFrames = append(ret.HFrames, &msg.HFrameStoreOne{
			Id:       id,
			ExpireTs: hf.ExpireTs.Unix(),
		})
	}

	return ret
}
