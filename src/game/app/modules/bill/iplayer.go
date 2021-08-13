package bill

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetBill() *Bill
}
