package marvelroll

import "fw/src/game/app/comp"

type IPlayer interface {
	comp.IPlayer

	GetMarvelRoll() *MarvelRoll
}
