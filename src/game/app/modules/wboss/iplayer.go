package wboss

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetWBoss() *WBoss
}
