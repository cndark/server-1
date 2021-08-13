package mopen

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetMOpen() *MOpen
}
