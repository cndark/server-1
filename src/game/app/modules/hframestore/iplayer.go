package hframestore

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetHFrameStore() *HFrameStore
}
