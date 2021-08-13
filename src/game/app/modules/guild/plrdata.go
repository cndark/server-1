package guild

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/battle"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

// ============================================================================

// 玩家身上家族数据
type GuildPlrData struct {
	// 签到
	Signed bool
	SignTs time.Time

	// 许愿
	WishTs  time.Time
	WishCnt map[int32]int32 // [item-num]cnt

	// 订单
	Order *order_t

	// 科技
	Tech *tech_t

	plr IPlayer
}

// ============================================================================

func init() {
	// reset daily
	evtmgr.On(gconst.Evt_PlrResetDaily, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		plr.GetGuildPlrData().daily_reset()
	})
}

func NewGuildPlrData() *GuildPlrData {
	return &GuildPlrData{
		WishTs:  time.Unix(0, 0),
		WishCnt: make(map[int32]int32),

		Order: new_order(),
		Tech:  new_tech(),
	}
}

// ============================================================================

func (self *GuildPlrData) Init(plr IPlayer) {
	self.plr = plr

	self.Order.init(plr)
	self.Tech.init(plr)
}

func (self *GuildPlrData) daily_reset() {
	self.Signed = false
	self.WishCnt = make(map[int32]int32)
}

// ============================================================================
// 签到

func (self *GuildPlrData) Sign() (ec int32, rwd *msg.Rewards) {
	// check
	gld := self.plr.GetGuild()
	if gld == nil {
		return Err.Guild_NotAMember, nil
	}

	if self.Signed {
		return Err.Guild_SignAlready, nil
	}

	conf := gamedata.ConfGuild.Query(gld.Lv)
	if conf == nil {
		return Err.Failed, nil
	}

	// sign
	op := self.plr.GetBag().NewOp(gconst.ObjFrom_GuildSign)
	for _, v := range conf.CheckinReward {
		op.Inc(v.Id, v.N)
	}
	rwd = op.Apply().ToMsg()

	gld.AddExp(self.plr, int64(conf.CheckinExp))

	self.Signed = true
	self.SignTs = time.Now()

	evtmgr.Fire(gconst.Evt_GuildSign, self.plr)

	return Err.OK, rwd
}

// ============================================================================
// 许愿

func (self *GuildPlrData) WishItem(num int32) (ec int32, seq int64) {
	// check
	gld := self.plr.GetGuild()
	if gld == nil {
		return Err.Guild_NotAMember, 0
	}

	// conf
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return Err.Failed, 0
	}

	// cd
	now := time.Now()
	if now.Sub(self.WishTs).Minutes() < float64(conf.GuildWishCd) {
		return Err.Guild_WishCD, 0
	}

	// item num
	L := int32(len(conf.GuildWish))
	if num < 0 || num >= L {
		return Err.Failed, 0
	}

	if self.WishCnt[num] > 0 {
		return Err.Guild_WishItemCnt, 0
	}

	// wish count
	if gld.WishCount(self.plr.GetId()) >= conf.GuildWishLimit {
		return Err.Guild_WishLimit, 0
	}

	// ok. update wish data
	self.WishTs = now
	self.WishCnt[num]++

	// add wish
	seq = gld.WishAdd(self.plr.GetId(), num)

	// broadcast
	gld.Broadcast(&msg.GS_GuildWishNew{
		Name: self.plr.GetName(),
		Num:  num,
		Seq:  seq,
	})

	evtmgr.Fire(gconst.Evt_GuildWish, self.plr)

	return Err.OK, seq
}

func (self *GuildPlrData) WishHelp(seq int64) (ec int32, rwd *msg.Rewards) {
	// check
	gld := self.plr.GetGuild()
	if gld == nil {
		return Err.Guild_NotAMember, nil
	}

	// find wish
	w := gld.Wish[seq]
	if w == nil {
		return Err.Guild_WishNotFound, nil
	}

	// already full help ?
	if w.Helps >= 3 {
		return Err.Guild_WishFullHelp, nil
	}

	// check counter
	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_GuildWishHelp)
	cop.DecCounter(gconst.Cnt_GuildWishHelp, 1)
	if ec := cop.CheckEnough(); ec != Err.OK {
		return ec, nil
	}

	// help rewards
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf != nil {
		n := int(self.plr.GetCounter().Get(gconst.Cnt_GuildWishHelp))
		L := len(conf.GuildHelpReward)

		if n < 0 {
			n = 0
		} else if n >= L {
			n = L - 1
		}

		cop.Inc(conf.GuildHelpReward[n].Id, conf.GuildHelpReward[n].N)
	}

	rwd = cop.Apply().ToMsg()

	// ok
	w.Helps++
	if w.Helps >= 3 {
		iplr := utils.FindPlayer(w.PlrId)
		if iplr != nil {
			iplr.(IPlayer).SendMsg(&msg.GS_GuildWishFullHelp{})
		}
	}

	evtmgr.Fire(gconst.Evt_GuildHelp, self.plr)

	return Err.OK, rwd
}

func (self *GuildPlrData) WishClose(seq int64) (ec int32, rwd *msg.Rewards) {
	// check
	gld := self.plr.GetGuild()
	if gld == nil {
		return Err.Guild_NotAMember, nil
	}

	// find wish
	w := gld.Wish[seq]
	if w == nil {
		return Err.Guild_WishNotFound, nil
	}

	// is yours ?
	if w.PlrId != self.plr.GetId() {
		return Err.Guild_WishNotYours, nil
	}

	// full help ?
	if w.Helps < 3 {
		return Err.Guild_WishNotFullHelp, nil
	}

	// ok. remove wish
	delete(gld.Wish, seq)

	// take rewards
	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return Err.Failed, nil
	}

	op := self.plr.GetBag().NewOp(gconst.ObjFrom_GuildWishRewards)
	if int(w.ItemNum) < len(conf.GuildWish) {
		op.Inc(conf.GuildWish[w.ItemNum].Id, conf.GuildWish[w.ItemNum].N)
	}
	rwd = op.Apply().ToMsg()

	return Err.OK, rwd
}

func (self *GuildPlrData) WishList() (ret []*msg.GuildWish) {
	// check
	gld := self.plr.GetGuild()
	if gld == nil {
		return
	}

	// list
	ret = make([]*msg.GuildWish, 0, len(gld.Wish))
	for _, v := range gld.Wish {
		ret = append(ret, &msg.GuildWish{
			Seq:     v.Seq,
			PlrId:   v.PlrId,
			ItemNum: v.ItemNum,
			Helps:   v.Helps,
		})
	}

	return
}

// ============================================================================
// 副本 boss

func (self *GuildPlrData) BossFight(tf *msg.TeamFormation, f func(*msg.GS_GuildBossFight_R)) {
	fail := func(ec int32) {
		f(&msg.GS_GuildBossFight_R{
			ErrorCode: ec,
		})
	}

	// get guild
	gld := self.plr.GetGuild()
	if gld == nil {
		fail(Err.Guild_NotAMember)
		return
	}

	// check team
	if !self.plr.IsTeamFormationValid(tf) {
		fail(Err.Plr_TeamInvalid)
		return
	}

	// cost
	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_GuildBossFight)
	cop.DecCounter(gconst.Cnt_GuildBossFight, 1)

	if ec := cop.CheckEnough(); ec != Err.OK {
		fail(ec)
		return
	}

	// ok. add fight rewards
	conf := gamedata.ConfGuildBoss.Query(gld.Boss.Cur.Num)
	if conf == nil {
		fail(Err.Failed)
		return
	}

	for _, v := range conf.ChallengeReward {
		cop.Inc(v.Id, v.N)
	}

	// fight
	bi := &msg.BattleInput{
		T1: self.plr.ToMsg_BattleTeam(tf),
		T2: battle.NewMonsterTeam().AddMonster(conf.Monster[0].Id, conf.Monster[0].Lv, 6).ToMsg_BattleTeam(),
		Args: map[string]string{
			"Module":      "GUILD_BOSS",
			"init_hp.2.6": core.F64toa(1 - gld.Boss.HpLeft),
			"RoundType":   core.I32toa(conf.RoundType),
		},
	}

	battle.Fight(bi, func(r *msg.BattleResult) {
		// check
		if r == nil {
			fail(Err.Failed)
			return
		}

		// ok
		replay := &msg.BattleReplay{Bi: bi, Br: r}
		rwd := cop.Apply().ToMsg()
		hpleft := 1 - core.Atof64(r.Args["hp_loss.2.6"])
		dmg := core.Atof64(r.Args["dmg_total.1"])

		f(&msg.GS_GuildBossFight_R{
			ErrorCode: Err.OK,
			Replay:    replay,
			Rewards:   rwd,
			Dmg:       dmg,
			HpLeft:    hpleft,
		})

		// update boss
		gld.Boss.update(self.plr.GetId(), dmg, hpleft)

		// evt
		evtmgr.Fire(gconst.Evt_GuildBossFight, self.plr, gld, dmg)
	})
}

// ============================================================================

func (self *GuildPlrData) ToMsg() *msg.GuildPlrData {
	ret := &msg.GuildPlrData{
		Signed:  self.Signed,
		WishTs:  self.WishTs.Unix(),
		WishCnt: self.WishCnt,
	}

	return ret
}
