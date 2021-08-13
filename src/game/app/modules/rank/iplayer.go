package rank

import (
	"fw/src/game/app/comp"
)

type IPlayer interface {
	comp.IPlayer

	GetWLevelLvNum() int32
	GetTowerLvNum() int32
	GetRankPlay() *RankPlay
}
