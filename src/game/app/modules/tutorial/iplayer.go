package tutorial

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetTutorial() *Tutorial
}
