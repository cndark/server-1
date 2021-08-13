package teammgr

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetTeamMgr() *TeamMgr
}
