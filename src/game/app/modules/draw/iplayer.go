package draw

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetDraw() *Draw
}
