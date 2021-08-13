package cloud

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetCloud() *Cloud
}
