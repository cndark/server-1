package signdaily

import "fw/src/game/app/comp"

type IPlayer interface {
	comp.IPlayer
	GetSignDaily() *SignDaily
}
