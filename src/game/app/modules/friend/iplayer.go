package friend

import (
	"fw/src/game/app/comp"
	"fw/src/game/app/modules/utils"
)

type IPlayer interface {
	comp.IPlayer

	GetFriend() *Friend
}

func load_player(uid string) IPlayer {
	plr := utils.LoadPlayer(uid)
	if plr == nil {
		return nil
	} else {
		return plr.(IPlayer)
	}
}

func find_player(uid string) IPlayer {
	plr := utils.FindPlayer(uid)
	if plr == nil {
		return nil
	} else {
		return plr.(IPlayer)
	}
}
