package wlevelfund

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetWLevelFund() *WLevelFund
	GetWLevelLvNum() int32
}
