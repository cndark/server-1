package misc

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetMisc() *Misc
}
