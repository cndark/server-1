package guild

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

const (
	RK_Owner  = 1
	RK_Vice   = 2
	RK_Member = 3
)

// ============================================================================

type Member struct {
	Id   string
	Rank int32

	gld *Guild
}

// ============================================================================
// 核心

func (self *Member) daily_reset() {
}

func (self *Member) GetName() string {
	plr := load_player(self.Id)
	if plr == nil {
		return ""
	} else {
		return plr.GetName()
	}
}

func (self *Member) GetLevel() int32 {
	plr := load_player(self.Id)
	if plr == nil {
		return 0
	} else {
		return plr.GetLevel()
	}
}

func (self *Member) SetRank(rk int32) int32 {
	if rk == self.Rank {
		return Err.OK
	}

	// check range
	if rk < RK_Owner || rk > RK_Member {
		return Err.Guild_InvalidRank
	}

	// vice full ?
	if ec := self.gld.IsFullForRank(rk); ec != Err.OK {
		return ec
	}

	// set
	self.Rank = rk

	// bind new owner
	if rk == RK_Owner {
		// demote old owner to member
		if self.gld.owner != nil {
			self.gld.owner.SetRank(RK_Member)
		}

		// bind
		self.gld.owner = self
	}

	// load player
	plr := load_player(self.Id)
	if plr != nil {
		// notify
		self.gld.Broadcast(&msg.GS_Guild_MbRank{
			PId:   self.Id,
			PName: plr.GetName(),
			Rank:  rk,
		})

		// fire
		evtmgr.Fire(gconst.Evt_GuildMemberRank, self.gld, plr, rk)
		evtmgr.Fire(gconst.Evt_GuildUserChange, self.gld, plr)

		// log
		self.gld.AddLog(C_Guild_Log_SetRank, map[string]string{
			"player":  plr.GetName(),
			"gldRank": core.I32toa(rk),
		})
	}

	return Err.OK
}

func (self *Member) Leave(reason int32) int32 {
	// check rank
	if self.Rank == RK_Owner {
		return Err.Guild_OwnerCantLeave
	}

	if self.gld.LockMemberLeave() {
		return Err.Guild_LeaveTimeLimited
	}

	// leave
	plr := load_player(self.Id)
	if plr == nil {
		return Err.Failed
	}

	self.gld.RemoveMember(plr, reason)

	return Err.OK
}

func (self *Member) ToMsg_Info() *msg.GuildMemberInfo {
	plr := load_player(self.Id)

	if plr == nil {
		return &msg.GuildMemberInfo{
			Plr: &msg.PlayerSimpleInfo{
				Id: self.Id,
			},
			Rank: self.Rank,
		}
	} else {
		return &msg.GuildMemberInfo{
			Plr:       plr.ToMsg_SimpleInfo(),
			Rank:      self.Rank,
			AtkPwr:    plr.GetAtkPwr(),
			Online:    plr.IsOnline(),
			OfflineTs: plr.GetOfflineTs().Unix(),
			SignTs:    plr.GetGuildPlrData().SignTs.Unix(),
		}
	}
}

// ============================================================================
