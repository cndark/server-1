package crusade

import (
	"fw/src/game/app/comp"
	"fw/src/game/app/modules/robot"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"strings"
)

type IPlayer interface {
	comp.IPlayer

	GetCrusade() *Crusade
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

// ============================================================================

// 远征对象接口(玩家和机器人)
type ICrusadePlayer interface {
	GetId() string
	GetName() string
	GetHead() string
	GetHFrame() int32
	GetLevel() int32
	GetAtkPwr() int32

	GetTeam(tp int32) *msg.TeamFormation
	ToMsg_SimpleInfo(atkpwr ...int32) *msg.PlayerSimpleInfo
	ToMsg_BattleTeam(tf *msg.TeamFormation, teamAtkPwr ...bool) *msg.BattleTeam
}

func IsRobot(id string) bool {
	return strings.HasPrefix(id, "bot")
}

func FindCrusadePlayer(id string) ICrusadePlayer {
	if IsRobot(id) {
		return robot.RobotMgr.FindRobot(id)
	}

	return find_player(id)
}
