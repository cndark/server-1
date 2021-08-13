package invite

import (
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

type Invite struct {
	Taken []int32
	plr   IPlayer
}

func NewInvite() *Invite {
	return &Invite{}
}

func (self *Invite) Init(plr IPlayer) {
	self.plr = plr
}

func (self *Invite) Take(id int32) (int32, *msg.Rewards) {
	conf := gamedata.ConfInviteReward.Query(id)
	if conf == nil {
		return Err.Failed, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_Invite)

	for _, v := range conf.Reward {
		op.Inc(v.Id, v.N)
	}

	rwds := op.Apply().ToMsg()

	self.Taken = append(self.Taken, id)

	return Err.OK, rwds
}

func (self *Invite) ToMsg() *msg.InviteData {
	return &msg.InviteData{
		Taken: self.Taken,
	}
}
