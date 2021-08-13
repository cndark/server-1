package robot

import (
	"fw/src/game/app/modules/battle"
	"fw/src/game/msg"
	"fw/src/shared/config"
)

// ============================================================================

// 机器人
type Robot struct {
	Id         string
	Name       string
	Head       string
	HFrame     int32
	Lv         int32
	AtkPwr     int32
	ArenaScore int32
	Team       *msg.TeamFormation // 布阵信息
	HeroLv     map[int32]int32    // 英雄等级[pos]lv
}

// ============================================================================

func (self *Robot) GetId() string {
	return self.Id
}

func (self *Robot) GetName() string {
	return self.Name
}

func (self *Robot) GetHead() string {
	return self.Head
}

func (self *Robot) GetHFrame() int32 {
	return self.HFrame
}

func (self *Robot) GetLevel() int32 {
	return self.Lv
}

func (self *Robot) GetAtkPwr() int32 {
	return self.AtkPwr
}

func (self *Robot) GetShowHero() int32 {
	for id := range self.Team.Formation {
		return int32(id)
	}

	return 0
}

func (self *Robot) GetArenaScore() int32 {
	return self.ArenaScore
}

func (self *Robot) AddArenaScore(v int32) {

}

func (self *Robot) AddArenaReplay(rid string, sscore, escore, ascore int32, plrInfo *msg.PlayerSimpleInfo) {

}

func (self *Robot) GetTeam(tp int32) *msg.TeamFormation {
	return self.Team
}

func (self *Robot) ToMsg_SimpleInfo(atkpwr ...int32) *msg.PlayerSimpleInfo {
	ret := &msg.PlayerSimpleInfo{
		Id:     self.GetId(),
		Name:   self.GetName(),
		SvrId:  config.CurGame.Id,
		Head:   self.GetHead(),
		HFrame: self.GetHFrame(),
		Lv:     self.GetLevel(),
		AtkPwr: self.GetAtkPwr(),
	}

	return ret
}

func (self *Robot) ToMsg_BattleTeam(tf *msg.TeamFormation, teamAtkPwr ...bool) *msg.BattleTeam {
	t := battle.NewMonsterTeam()
	for id, pos := range tf.Formation {
		lv := self.HeroLv[pos]
		if lv == 0 {
			lv = self.Lv
		}

		t.AddMonster(int32(id), lv, pos)
	}

	ret := t.ToMsg_BattleTeam()
	ret.Player = self.ToMsg_SimpleInfo()

	return ret
}
