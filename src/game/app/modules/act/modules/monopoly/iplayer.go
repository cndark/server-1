package monopoly

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer
	GetWLevelLvNum() int32
}
