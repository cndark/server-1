package arena

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math"
	"time"
)

// ============================================================================

// 竞技场玩家信息
type Arena struct {
	Records []*record_t

	enemies []string // 选中的敌人
	plr     IPlayer
}

// 战报记录
type record_t struct {
	ReplayId   string                // 战报id
	Revenge    int32                 // 复仇次数
	SelfScore  int32                 // 自己分数
	EnemyScore int32                 // 对手分数
	AddScore   int32                 // 分数变化
	Ts         int64                 // 时间
	Enemy      *msg.PlayerSimpleInfo // 对手玩家简要信息
}

// ============================================================================

func NewArena() *Arena {
	return &Arena{}
}

// ============================================================================

func (self *Arena) Init(plr IPlayer) {
	self.plr = plr
}

func (self *Arena) GetScore() int32 {
	return ArenaMgr.GetScore(self.plr.GetId())
}

func (self *Arena) AddScore(v int32) {
	ArenaMgr.AddScore(self.plr.GetId(), v)
}

func (self *Arena) AddArenaReplay(rid string, sscore, escore, ascore int32, eInfo *msg.PlayerSimpleInfo) {
	self.Records = append(self.Records, &record_t{
		ReplayId:   rid,
		SelfScore:  sscore,
		EnemyScore: escore,
		AddScore:   ascore,
		Ts:         time.Now().Unix(),
		Enemy:      eInfo,
	})

	if l := len(self.Records); l > 8 {
		for _, v := range self.Records[:l-8] {
			battle.ReplayDel(dbmgr.DBGame, dbmgr.C_tabname_replayarena, v.ReplayId)
		}

		self.Records = self.Records[l-8:]
	}
}

func (self *Arena) CheckEnemy(uid string) bool {
	for _, v := range self.enemies {
		if uid == v {
			return true
		}
	}

	return false
}

func (self *Arena) CheckRevenge(ridx int32) int32 {
	if ridx >= int32(len(self.Records)) || ridx < 0 {
		return Err.Arena_EnemyNotFound
	}

	if self.Records[ridx].Revenge > 0 {
		return Err.Arena_RevengeError
	}

	return Err.OK
}

func (self Arena) UpdateRevenge(ridx int32) {
	if ridx >= 0 && ridx < int32(len(self.Records)) {
		self.Records[ridx].Revenge++
	}
}

func (self *Arena) SetEnemies(ids []string) {
	self.enemies = ids
}

func (self *Arena) Fight(eid string, T *msg.TeamFormation, revengeIdx int32,
	cb func(ec int32, replay *msg.BattleReplay, rwds *msg.Rewards, addScore1, addScore2 int32)) {

	if revengeIdx < 0 {
		if !self.CheckEnemy(eid) {
			cb(Err.Arena_EnemyNotFound, nil, nil, 0, 0)
			return
		}
	} else {
		if ec := self.CheckRevenge(revengeIdx); ec != Err.OK {
			cb(ec, nil, nil, 0, 0)
			return
		}
	}

	eplr := FindArenaPlayer(eid)
	if eplr == nil {
		cb(Err.Arena_PlayerNotFound, nil, nil, 0, 0)
		return
	}

	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil {
		cb(Err.Failed, nil, nil, 0, 0)
		return
	}

	if !self.plr.IsTeamFormationValid(T) {
		cb(Err.Plr_TeamInvalid, nil, nil, 0, 0)
		return
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_ArenaFight)
	op.Dec(gconst.ArenaCard, 1)
	if ec := op.CheckEnough(); ec != Err.OK {
		cb(ec, nil, nil, 0, 0)
		return
	}

	// gen input
	input := &msg.BattleInput{
		T1: self.plr.ToMsg_BattleTeam(T),
		T2: eplr.ToMsg_BattleTeam(eplr.GetTeam(gconst.TeamType_Dfd)),
		Args: map[string]string{
			"Module":    "AP_ARENA",
			"RoundType": "2",
		},
	}

	// fight
	battle.Fight(input, func(r *msg.BattleResult) {
		if r == nil {
			cb(Err.Common_BattleResError, nil, nil, 0, 0)
			return
		}

		op := self.plr.GetBag().NewOp(gconst.ObjFrom_ArenaFight)
		op.Dec(gconst.ArenaCard, 1)
		if ec := op.CheckEnough(); ec != Err.OK {
			cb(ec, nil, nil, 0, 0)
			return
		}

		addScore := 50 * (1 - 1/(1+math.Pow10(int((eplr.GetArenaScore()-self.plr.GetArenaScore()))/400)))
		if math.IsInf(addScore, 0) || math.IsNaN(addScore) {
			addScore = 1
		}

		// award
		addScore1 := int32(addScore)
		if r.Winner != 1 {
			addScore1 = -addScore1
		}
		addScore2 := -int32(addScore1)

		item := utils.BattleReward(self.plr.GetLevel())
		for id, n := range item {
			op.Inc(id, n)
		}

		self.plr.AddArenaScore(int32(addScore1))
		eplr.AddArenaScore(addScore2)
		self.SetEnemies([]string{})

		rwds := op.Apply().ToMsg()

		self.UpdateRevenge(revengeIdx)

		// replay
		replay := &msg.BattleReplay{
			Bi: input,
			Br: r,
		}

		rid := battle.ReplaySave(dbmgr.DBGame, dbmgr.C_tabname_replayarena, replay)

		eplr.AddArenaReplay(rid, eplr.GetArenaScore(), self.plr.GetArenaScore(),
			addScore2, self.plr.ToMsg_SimpleInfo())

		if !IsRobot(eplr.GetId()) {
			eplr.(IPlayer).SendMsg(&msg.GS_ArenaFighted{
				AddScore: addScore2,
			})
		}

		cb(Err.OK, replay, rwds, addScore1, addScore2)

		// evt
		evtmgr.Fire(gconst.Evt_ArenaFight, self.plr, r.Winner == 1,
			IsRobot(eid), self.plr.GetAtkPwr()-eplr.GetAtkPwr(), addScore1)

		if r.Winner != 1 {
			if !IsRobot(eplr.GetId()) {
				evtmgr.Fire(gconst.Evt_ArenaScore, eplr, eplr.GetArenaScore(), addScore2)
			}
		} else {
			evtmgr.Fire(gconst.Evt_ArenaScore, self.plr, self.plr.GetArenaScore(), addScore1)
		}

		// rank
		ArenaMgr.end_battle_sort_rank(self.plr.GetId(), eplr.GetId())
	})

}

// ============================================================================

func (self *Arena) ToMsg() *msg.ArenaData {
	ret := &msg.ArenaData{
		Score: self.GetScore(),
		Stage: ArenaMgr.stage,
		Ts1:   ArenaMgr.Ts1.Unix(),
		Ts2:   ArenaMgr.Ts2.Unix(),
	}

	return ret
}

func (self *Arena) Records_ToMsg() (ret []*msg.ArenaRecord) {
	for _, v := range self.Records {
		ret = append(ret, &msg.ArenaRecord{
			ReplayId:   v.ReplayId,
			Revenge:    v.Revenge,
			SelfScore:  v.SelfScore,
			EnemyScore: v.EnemyScore,
			AddScore:   v.AddScore,
			Ts:         v.Ts,
			Enemy:      v.Enemy,
		})
	}

	return
}
