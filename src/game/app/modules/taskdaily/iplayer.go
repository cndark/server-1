package taskdaily

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetTaskDaily() *TaskDaily
}
