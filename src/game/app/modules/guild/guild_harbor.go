package guild

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
)

// ============================================================================

type harbor_t struct {
	Lv      int32
	Xp      int64
	Donates map[string]int64

	gld *Guild
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_GuildLeave, func(args ...interface{}) {
		gld := args[0].(*Guild)
		plr := args[1].(IPlayer)

		gld.Harbor.remove_donate(plr.GetId())
	})
}

// ============================================================================

func new_harbor() *harbor_t {
	return &harbor_t{
		Lv:      1,
		Xp:      0,
		Donates: make(map[string]int64),
	}
}

func (self *harbor_t) init(g *Guild) {
	self.gld = g
}

func (self *harbor_t) add_xp(v int64) {
	if v <= 0 {
		return
	}

	new_lv := self.Lv
	self.Xp += v

	// consume exp
	for {
		conf := gamedata.ConfGuildHarbor.Query(new_lv)
		if conf == nil || conf.LevelExp == 0 {
			// 满级
			self.Xp = 0
			break
		}

		if self.Xp < conf.LevelExp {
			break
		}

		self.Xp -= conf.LevelExp
		new_lv++
	}

	// set new lv
	self.Lv = new_lv

	// notify
	self.gld.Broadcast(&msg.GS_GuildHarborXpChange{
		Lv: self.Lv,
		Xp: self.Xp,
	})
}

func (self *harbor_t) AddDonate(plrid string, v int64) {
	// add xp
	self.add_xp(v)

	// acc record
	self.Donates[plrid] += v
}

func (self *harbor_t) remove_donate(plrid string) {
	delete(self.Donates, plrid)
}

func (self *harbor_t) ToMsg_DonateList() (ret []*msg.GuildHarborDonateRec) {
	ret = make([]*msg.GuildHarborDonateRec, 0, len(self.Donates))
	for k, v := range self.Donates {
		plr := load_player(k)
		if plr == nil {
			continue
		}

		ret = append(ret, &msg.GuildHarborDonateRec{
			Plr: plr.ToMsg_SimpleInfo(),
			Val: v,
		})
	}
	return
}
