package daysign

import "fw/src/game/app/comp"

type IPlayer interface {
	comp.IPlayer

	GetDaySign() *DaySign
}
