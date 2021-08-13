package appoint

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetAppoint() *Appoint
}
