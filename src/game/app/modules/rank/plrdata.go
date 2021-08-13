package rank

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/ranksvc"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"fw/src/shared/config"
)

// ============================================================================

var (
	V_LikableRanks = map[int32]int32{ // [rankid]like-counter
		gconst.RankId_WLevelLv: gconst.Cnt_RankLikeWLevel,
		gconst.RankId_TowerLv:  gconst.Cnt_RankLikeTower,
	}
)

// ============================================================================

type RankPlay struct {
	LikeRecords []*like_rec_t   // like records every day
	BeLikedCnt  map[int32]int32 // [rankid]cnt

	plr IPlayer
}

type like_rec_t struct {
	rkid  int32
	plrid string
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_PlrResetDaily, func(args ...interface{}) {
		plr := args[0].(IPlayer)

		plr.GetRankPlay().daily_reset()
	})
}

// ============================================================================

func NewRankPlay() *RankPlay {
	return &RankPlay{
		BeLikedCnt: make(map[int32]int32),
	}
}

func (self *RankPlay) Init(plr IPlayer) {
	self.plr = plr
}

func (self *RankPlay) Like(rkid int32, plrid string) (ec int32, rwd *msg.Rewards) {
	// check rankid
	c := V_LikableRanks[rkid]
	if c == 0 {
		return Err.Failed, nil
	}

	// check plrid
	rk := ranksvc.TryGet(ranksvc.RankType_Local, config.CurGame.Id, rkid)
	if rk == nil {
		return Err.RankPlay_RankNotReady, nil
	}

	n := ranksvc.GetRowPos(rk, plrid)
	if n == 0 {
		return Err.RankPlay_PlayerNotOnList, nil
	}

	// check cost
	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_RankPlayLike)
	cop.DecCounter(c, 1)
	if ec := cop.CheckEnough(); ec != Err.OK {
		return ec, nil
	}

	// check if we should get the rewards
	found := false
	for _, v := range self.LikeRecords {
		if v.rkid == rkid && v.plrid == plrid {
			found = true
			break
		}
	}

	if !found {
		conf := gamedata.ConfGlobalPublic.Query(1)
		if conf != nil {
			for _, v := range conf.RankLikeReward {
				cop.Inc(v.Id, v.N)
			}
		}
	}

	// apply
	rwd = cop.Apply().ToMsg()

	// add like rec
	self.LikeRecords = append(self.LikeRecords, &like_rec_t{
		rkid:  rkid,
		plrid: plrid,
	})

	// inc be liked count
	iplr := utils.LoadPlayer(plrid)
	if iplr != nil {
		plr := iplr.(IPlayer)
		plr.GetRankPlay().BeLikedCnt[rkid]++
	}

	//evt
	evtmgr.Fire(gconst.Evt_RankLike, self.plr, rkid)

	return Err.OK, rwd
}

func (self *RankPlay) ToMsg() *msg.RankPlayData {
	return nil
}

// ============================================================================

func (self *RankPlay) daily_reset() {
	self.LikeRecords = nil
}
