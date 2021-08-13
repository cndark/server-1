package wlevel

import (
	"fw/src/game/app/comp"
	"fw/src/game/msg"
)

type IPlayer interface {
	comp.IPlayer

	GetWLevel() *WLevel

	SetTeam(tp int32, T *msg.TeamFormation)
}
