package gwar

import (
	"fw/src/game/app/comp"
	"fw/src/game/app/modules/guild"
)

type IPlayer interface {
	comp.IPlayer

	GetGWar() *GWar
	GetGuild() *guild.Guild
}
