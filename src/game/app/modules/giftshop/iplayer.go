package giftshop

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetGiftShop() *GiftShop
}
