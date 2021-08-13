package taskachv

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetTaskAchv() *TaskAchv
}
