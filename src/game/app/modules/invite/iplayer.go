package invite

import "fw/src/game/app/comp"

type IPlayer interface {
	comp.IPlayer

	GetInvite() *Invite
}
