package warcup

import (
	"fw/src/game/app/comp"
	"fw/src/game/app/modules/robot"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"strings"
)

type IPlayer interface {
	comp.IPlayer

	GetArenaScore() int32
	GetWarCup() *WarCup
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

// 杯赛对象接口(玩家和机器人)
type IWarCupPlayer interface {
	GetId() string
	GetLevel() int32
	GetAtkPwr() int32
	GetArenaScore() int32

	ToMsg_SimpleInfo(atkpwr ...int32) *msg.PlayerSimpleInfo
	GetTeam(tp int32) *msg.TeamFormation
	ToMsg_BattleTeam(tf *msg.TeamFormation, teamAtkPwr ...bool) *msg.BattleTeam
}

func IsRobot(id string) bool {
	return strings.HasPrefix(id, "bot")
}

func FindWarCupPlayer(id string) IWarCupPlayer {
	if IsRobot(id) {
		return robot.RobotMgr.FindRobot(id)
	}

	return find_player(id)
}
