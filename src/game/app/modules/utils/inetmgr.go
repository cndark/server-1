package utils

import (
	"fw/src/core"
	"fw/src/game/msg"
	"fw/src/shared/config"
)

// ============================================================================

var inetmgr INetMgr

// ============================================================================

type INetMgr interface {
	Send2Game(svrid int32, message msg.Message)
	Send2CrossPlayer(svrid int32, plrid string, message msg.Message)
}

// ============================================================================

func Send2Game(svrid int32, message msg.Message) {
	inetmgr.Send2Game(svrid, message)
}

func Send2CrossPlayer(svrid int32, plrid string, message msg.Message) {
	inetmgr.Send2CrossPlayer(svrid, plrid, message)
}

func BroadcastGames(message msg.Message, except_self ...bool) {
	m := config.Games
	b := core.DefFalse(except_self)
	curid := config.CurGame.Id

	for _, v := range m {
		if b && v.Id == curid {
			continue
		}

		inetmgr.Send2Game(v.Id, message)
	}
}
