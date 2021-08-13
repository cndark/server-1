package taskgrow

import "fw/src/game/app/comp"

type IPlayer interface {
	comp.IPlayer

	GetTaskGrow() *TaskGrow
}
