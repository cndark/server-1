package act

import (
	"fw/src/game/app/comp"
	"fw/src/game/msg"
	"time"

	"fw/src/core/log"

	"go.mongodb.org/mongo-driver/bson"
)

// ============================================================================

type IPlayer interface {
	comp.IPlayer

	GetAct() *Act
}

// ============================================================================

type Act struct {
	Acts map[string]*act_plrdata_t // [act-name]actdata

	plr IPlayer
}

type act_plrdata_t struct {
	Ver  time.Time   // data key
	Data interface{} // real data
}

// ============================================================================

func NewAct() *Act {
	return &Act{
		Acts: make(map[string]*act_plrdata_t),
	}
}

func (self *Act) Init(plr IPlayer) {
	self.plr = plr

	unix_t0 := time.Unix(0, 0)

	EachAct(func(a IAct) {
		name := a.GetName()

		d := self.Acts[name]
		if d == nil {
			// no data. create it
			d = &act_plrdata_t{
				Ver:  unix_t0,
				Data: nil,
			}
			self.Acts[name] = d
		}

		d.Data = alloc_plr_data(d.Data, a)
	})
}

func alloc_plr_data(m interface{}, a IAct) (out interface{}) {
	out = a.NewPlrData()

	if m == nil {
		return
	}

	data, err := bson.Marshal(m)
	if err != nil {
		log.Error("alloc act data failed:", err)
		return
	}

	err = bson.Unmarshal(data, out)
	if err != nil {
		log.Error("alloc act data failed:", err)
		return
	}

	return
}

func (self *Act) GetActRawData(name string) interface{} {
	d := self.Acts[name]
	if d == nil {
		return nil
	} else {
		// reset data if ver differs
		a := FindAct(name)
		if a != nil {
			if !d.Ver.Equal(a.get_ver()) {
				d.Ver = a.get_ver()
				d.Data = a.NewPlrData()
				a.ResetPlrData(self.plr)
			}
		}

		// return
		return d.Data
	}
}

// ============================================================================

func (self *Act) ToMsg() (ret []*msg.ActState) {
	EachAct(func(a IAct) {
		if a.GetStage() != "close" {
			ret = append(ret, a.tomsg_state())
		}
	})

	return
}
