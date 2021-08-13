package gwar

import (
	. "fw/src/core/math"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

type GWar struct {
	plr IPlayer
}

// ============================================================================

func NewGWar() *GWar {
	return &GWar{}
}

func (self *GWar) Init(plr IPlayer) {
	self.plr = plr
}

func (self *GWar) Fight(tarid string, tf *msg.TeamFormation, f func(*msg.GS_GWarFight_R)) int32 {
	fail := func(ec int32) int32 {
		f(&msg.GS_GWarFight_R{
			ErrorCode: ec,
		})
		return ec
	}

	// check conf
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return fail(Err.Failed)
	}

	// check stage
	if g_stage != c_Stage_Match {
		return fail(Err.Failed)
	}

	// get guild
	gld := self.plr.GetGuild()
	if gld == nil {
		return fail(Err.Guild_NotAMember)
	}

	// enrolled ?
	gd := localdata.Glds[gld.Id]
	if gd == nil {
		return fail(Err.GWar_NotEnrolled)
	}

	// arena module open?
	if !self.plr.IsModuleOpen(gconst.ModuleId_Arena) {
		return fail(Err.GWar_NotEnrolled)
	}

	// check if matched
	if gd.G2 == nil {
		return fail(Err.GWar_NoMatch)
	}

	// find target
	T2 := gd.G2.Mbs[tarid]
	if T2 == nil {
		return fail(Err.GWar_TargetNotFound)
	}

	// check target
	info2 := gd.get_g2_plrinfo(tarid)
	if info2.Done {
		return fail(Err.GWar_TargetDead)
	}
	if info2.locked {
		return fail(Err.GWar_TargetLocked)
	}

	// check team
	if !self.plr.IsTeamFormationValid(tf) {
		return fail(Err.Plr_TeamInvalid)
	}

	// check counter
	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_GWarFight)
	cop.DecCounter(gconst.Cnt_GWarFight, 1)
	if ec := cop.CheckEnough(); ec != Err.OK {
		return fail(ec)
	}

	// ok. fight
	info2.locked = true

	bi := &msg.BattleInput{
		T1: self.plr.ToMsg_BattleTeam(tf),
		T2: T2,
		Args: map[string]string{
			"Module":    "GUILD_WAR",
			"RoundType": "0",
		},
	}

	battle.Fight(bi, func(r *msg.BattleResult) {
		defer func() { info2.locked = false }()

		// check
		if r == nil {
			fail(Err.Failed)
			return
		}

		// ok
		cop.Apply()
		replay := &msg.BattleReplay{Bi: bi, Br: r}

		// calc jf
		jf_add := int32(0)
		if r.Winner == 1 {
			jf_add = int32(float64(T2.Player.AtkPwr)/10000+float64(bi.T1.Player.AtkPwr)/30000) + 1
		}

		// update plr1 info
		info1 := gd.get_g1_plrinfo(self.plr.GetId())
		info1.Cnt++
		info1.Jf += jf_add

		// update plr2 info
		info2.Val += int32(
			MaxFloat64(float64(bi.T1.Player.AtkPwr)/float64(T2.Player.AtkPwr), 0.1) * float64(conf.GuildWarKillScore),
		)

		if r.Winner == 1 || info2.Val >= conf.GuildWarKillShield {
			info2.Done = true
			gd.G2DeadCount++
		}

		// callback
		f(&msg.GS_GWarFight_R{
			ErrorCode: Err.OK,
			Replay:    replay,
			JfAdd:     jf_add,
			Val:       info2.Val,
			Done:      info2.Done,
		})

		// check if g2 is defeated
		if gd.G2DeadCount >= len(gd.G2.Mbs) {
			on_gld_win(gd)
		}
	})

	return 0
}

func (self *GWar) GetSummary() (ret *msg.GS_GWarGetSummary_R) {
	ret = &msg.GS_GWarGetSummary_R{
		Stage: int32(g_stage),
		Ts2:   g_t2.Unix(),
	}

	// get guild
	gld := self.plr.GetGuild()
	if gld == nil {
		return
	}

	// enrolled ?
	gd := localdata.Glds[gld.Id]
	if gd == nil {
		return
	}

	// ok
	if gd.G2 != nil {
		ret.G2 = gd.G2.Base.tomsg()
	}

	e := enroll_index[gld.Id]
	if e != nil {
		ret.G1Jf = e.Jf
	}

	return
}

func (self *GWar) GetG2Members() (ec int32, ret []*msg.GWarGuildMember) {
	// check stage
	if g_stage < c_Stage_Match {
		return Err.Failed, nil
	}

	// get guild
	gld := self.plr.GetGuild()
	if gld == nil {
		return Err.Guild_NotAMember, nil
	}

	// enrolled ?
	gd := localdata.Glds[gld.Id]
	if gd == nil {
		return Err.GWar_NotEnrolled, nil
	}

	// ok
	if gd.G2 != nil {
		ret = make([]*msg.GWarGuildMember, 0, len(gd.G2.Mbs))
		for _, v := range gd.G2.Mbs {
			info2 := gd.get_g2_plrinfo(v.Player.Id)

			ret = append(ret, &msg.GWarGuildMember{
				Plr:  v.Player,
				Val:  info2.Val,
				Done: info2.Done,
			})
		}
	}

	return Err.OK, ret
}

func (self *GWar) GetGuildRank() []*msg.GWarGuildRankRow {
	// rank: top 20
	L := len(enroll_rank)
	if L > 20 {
		L = 20
	}

	ret := make([]*msg.GWarGuildRankRow, 0, L)
	for _, v := range enroll_rank[:L] {
		ret = append(ret, &msg.GWarGuildRankRow{
			Id:    v.Id,
			Name:  v.Name,
			Icon:  v.Icon,
			Lv:    v.Lv,
			SvrId: v.SvrId,
			Jf:    v.Jf,
		})
	}

	return ret
}

func (self *GWar) GetPlrRank() (ec int32, ret []*msg.GWarPlrRankRow) {
	// get guild
	gld := self.plr.GetGuild()
	if gld == nil {
		return Err.Guild_NotAMember, nil
	}

	// enrolled ?
	gd := localdata.Glds[gld.Id]
	if gd == nil {
		return Err.GWar_NotEnrolled, nil
	}

	// ok
	ret = make([]*msg.GWarPlrRankRow, 0, len(gd.G1Plrs))
	for uid, v := range gd.G1Plrs {
		iplr := utils.LoadPlayer(uid)
		if iplr == nil {
			continue
		}

		plr := iplr.(IPlayer)
		ap := plr.GetTeamAtkPwr(plr.GetTeam(gconst.TeamType_Dfd))

		ret = append(ret, &msg.GWarPlrRankRow{
			Plr: plr.ToMsg_SimpleInfo(ap),
			Cnt: v.Cnt,
			Jf:  v.Jf,
		})
	}

	return Err.OK, ret
}

// ============================================================================

func (self *GWar) ToMsg() *msg.GWarData {
	return &msg.GWarData{}
}
