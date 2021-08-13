package tower

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetTower() *Tower
}
