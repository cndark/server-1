package shop

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetShop() *Shop
}
