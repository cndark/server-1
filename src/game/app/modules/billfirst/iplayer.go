package billfirst

import "fw/src/game/app/comp"

type IPlayer interface {
	comp.IPlayer

	GetBillFirst() *BillFirst
}
