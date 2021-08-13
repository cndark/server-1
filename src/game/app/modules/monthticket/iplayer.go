package monthticket

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetMonthTicket() *MonthTicket
}
