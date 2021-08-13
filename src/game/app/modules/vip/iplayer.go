package vip

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetVip() *Vip
}
