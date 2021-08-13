package targetdays

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetTargetDays() *TargetDays
}
