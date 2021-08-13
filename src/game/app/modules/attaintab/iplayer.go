package attaintab

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetAttainTab() *AttainTab
}
