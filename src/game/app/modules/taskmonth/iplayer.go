package taskmonth

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetTaskMonth() *TaskMonth
}
