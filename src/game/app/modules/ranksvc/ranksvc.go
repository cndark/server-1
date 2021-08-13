package ranksvc

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	. "fw/src/core/math"
	"fw/src/core/sched/loop"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/svrgrp"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"fw/src/shared/config"
	"math"
	"reflect"
	"sort"
	"time"
)

// ============================================================================

const (
	RankType_Local = 1
	RankType_Cross = 2
)

const (
	SortLevel_Local = 1
	SortLevel_Cross = 2
	SortLevel_All   = 3
)

// ============================================================================
/*
	rank table cols:
		_id, type, sgid, rankid, data, ts

	type:	1 位十进制. (1: local; 2: cross)
	sgid:	4 位十进制. (svrid if local; grpid if cross)
	rankid:	4 位十进制

	_id (full-rankid):
		type * 10^8 + sgid * 10^4 + rankid

	grpid = 0 means:
		cross all servers
*/

var (
	r_ranks = make(map[int32]*rank_t) // [full-rankid]

	G *svrgrp.Group
)

// ============================================================================

type RankRaw struct {
	Id        int32                  // rankid
	A         interface{}            // 要排序的数组
	ScoreFunc func(int) float64      // 分数函数
	InfoFunc  func(int) *RankRowInfo // 信息函数
	SortLevel int                    // 排序等级
}

type RankRowInfo struct {
	Plr *msg.PlayerSimpleInfo
	Gld *msg.GuildRow
}

type rank_t struct {
	data         []*RankRow         // data
	need_updates int                // need update count
	req          []func([]*RankRow) // pending requests
}

type RankRow struct {
	Score float64
	Info  *RankRowInfo
}

// ============================================================================

func init() {
	evtmgr.On(gconst.Evt_SvrGrpReady, func(...interface{}) {
		G = svrgrp.GetGroupCross()
	})

	evtmgr.On(gconst.Evt_GsPush_RankCacheExpire, func(args ...interface{}) {
		sarg := args[0].([]string)

		full_rkid := core.Atoi32(sarg[0])
		rk := r_ranks[full_rkid]
		if rk != nil {
			rk.expire()
		}
	})
}

// ============================================================================

func (self *RankRowInfo) Id() string {
	if self.Plr != nil {
		return self.Plr.Id
	} else if self.Gld != nil {
		return self.Gld.Id
	} else {
		return ""
	}
}

func (self *RankRowInfo) AtkPwr() int32 {
	if self.Plr != nil {
		return self.Plr.AtkPwr
	} else if self.Gld != nil {
		return self.Gld.AtkPwr
	} else {
		return 0
	}
}

func (self *rank_t) expire() {
	if self.need_updates < 2 {
		self.need_updates++
	}
}

// ============================================================================

func Push(raw *RankRaw) {
	core.Go(func() {
		// decorate score func with additional info-speed
		sf := raw.ScoreFunc
		raw.ScoreFunc = func(i int) float64 {
			ap := raw.InfoFunc(i).AtkPwr() + 1
			if ap > 1 {
				return sf(i)*1000 + math.Log(float64(ap))
			} else {
				return sf(i) * 1000
			}
		}

		// local sort
		local_sort(raw)

		// cross sort
		time.Sleep(time.Second * 30)
		if raw.SortLevel >= SortLevel_Cross && G.IsMaster() {
			cross_sort(raw)
		}

		// all sort
		time.Sleep(time.Second * 30)
		if raw.SortLevel >= SortLevel_All && config.CurGame.Id == 1 {
			all_sort(raw)
		}

		// expire cross cache (local cache has already updated just now)
		time.Sleep(time.Second * 30)
		if raw.SortLevel >= SortLevel_Cross {
			loop.Push(func() {
				max_grpid := gamedata.ConfLimitM.Query().MaxSvrGrpId
				for i := int32(0); i <= max_grpid; i++ {
					rk := r_ranks[make_full_rankid(RankType_Cross, i, raw.Id)]
					if rk != nil {
						rk.expire()
					}
				}
			})
		}
	})
}

func local_sort(raw *RankRaw) {
	// sort
	sort.Slice(raw.A, func(i, j int) bool {
		return raw.ScoreFunc(i) > raw.ScoreFunc(j)
	})

	// create rank data
	L := MinInt(gconst.RankMaxRows(raw.Id), reflect.ValueOf(raw.A).Len())
	D := make([]*RankRow, 0, L)
	for i := 0; i < L; i++ {
		score := math.Floor(raw.ScoreFunc(i) / 1000)
		if score > 0 {
			D = append(D, &RankRow{
				Score: score,
				Info:  raw.InfoFunc(i),
			})
		}
	}

	// update db
	full_rkid := make_full_rankid(RankType_Local, config.CurGame.Id, raw.Id)
	err := dbmgr.DBCross.Replace(
		dbmgr.C_tabname_rank,
		full_rkid,
		db.M{
			"type":   RankType_Local,
			"sgid":   config.CurGame.Id,
			"rankid": raw.Id,
			"d":      D,
			"ts":     time.Now(),
		},
	)
	if err != nil {
		log.Error("update local rank failed:", config.CurGame.Id, raw.Id, err)
	}

	// update rank map
	loop.Push(func() {
		r_ranks[full_rkid] = &rank_t{
			data: D,
		}
	})
}

func cross_sort(raw *RankRaw) {
	// sort
	var obj []*struct {
		D []*RankRow
	}

	ids := G.Svrs[:G.AvailLen]
	err := dbmgr.DBCross.Aggregate(
		dbmgr.C_tabname_rank,
		[]db.M{
			{"$match": db.M{
				"type":   RankType_Local,
				"rankid": raw.Id,
				"sgid":   db.M{"$in": ids},
			}},
			{"$project": db.M{
				"d": 1,
			}},
			{"$unwind": "$d"},
			{"$sort": db.M{"d.score": -1}},
			{"$limit": gconst.RankMaxRows(raw.Id)},
			{"$group": db.M{
				"_id": nil,
				"d":   db.M{"$push": "$d"},
			}},
		},
		&obj,
	)
	if err != nil {
		log.Error("rank cross sort failed:", G.Id, raw.Id, err)
		return
	}
	if len(obj) == 0 {
		return
	}

	// update db
	now := time.Now()
	err = dbmgr.DBCross.Replace(
		dbmgr.C_tabname_rank,
		make_full_rankid(RankType_Cross, G.Id, raw.Id),
		db.M{
			"type":   RankType_Cross,
			"sgid":   G.Id,
			"rankid": raw.Id,
			"d":      obj[0].D,
			"ts":     now,
		},
	)
	if err != nil {
		log.Error("update cross rank failed:", G.Id, raw.Id, err)
	}
}

func all_sort(raw *RankRaw) {
	// sort
	var obj []*struct {
		D []*RankRow
	}

	err := dbmgr.DBCross.Aggregate(
		dbmgr.C_tabname_rank,
		[]db.M{
			{"$match": db.M{
				"type":   RankType_Cross,
				"rankid": raw.Id,
				"sgid":   db.M{"$gt": 0},
			}},
			{"$project": db.M{
				"d": 1,
			}},
			{"$unwind": "$d"},
			{"$sort": db.M{"d.score": -1}},
			{"$limit": gconst.RankMaxRows(raw.Id)},
			{"$group": db.M{
				"_id": nil,
				"d":   db.M{"$push": "$d"},
			}},
		},
		&obj,
	)
	if err != nil {
		log.Error("rank all sort failed:", raw.Id, err)
		return
	}
	if len(obj) == 0 {
		return
	}

	// update db
	now := time.Now()
	err = dbmgr.DBCross.Replace(
		dbmgr.C_tabname_rank,
		make_full_rankid(RankType_Cross, 0, raw.Id),
		db.M{
			"type":   RankType_Cross,
			"sgid":   0,
			"rankid": raw.Id,
			"d":      obj[0].D,
			"ts":     now,
		},
	)
	if err != nil {
		log.Error("update all rank failed:", raw.Id, err)
	}
}

func SetCrossOnly(sgid, rkid int32, data []*RankRow) {
	core.Go(func() {
		// update db
		full_rkid := make_full_rankid(RankType_Cross, sgid, rkid)
		err := dbmgr.DBCross.Replace(
			dbmgr.C_tabname_rank,
			full_rkid,
			db.M{
				"type":   RankType_Cross,
				"sgid":   sgid,
				"rankid": rkid,
				"d":      data,
				"ts":     time.Now(),
			},
		)
		if err != nil {
			log.Error("set cross-only rank failed:", sgid, rkid, err)
		}

		// tell ALL servers to invalidate the cache of this item
		utils.GsPushAll(gconst.Evt_GsPush_RankCacheExpire, []string{core.I32toa(full_rkid)}, nil)
	})
}

func Get(tp, sgid, rkid int32, f func([]*RankRow)) {
	// find from cache
	fullid := make_full_rankid(tp, sgid, rkid)
	rk := r_ranks[fullid]
	if rk == nil {
		rk = &rank_t{
			need_updates: 1,
		}
		r_ranks[fullid] = rk
	}

	// check cache expiration
	if rk.need_updates <= 0 {
		f(rk.data)
		return
	}

	// cache expired. fetch from db
	rk.req = append(rk.req, f)
	if len(rk.req) > 1 {
		return
	}

	core.Go(func() {
		var obj struct {
			D  []*RankRow
			Ts time.Time
		}

		err := dbmgr.DBCross.GetObject(
			dbmgr.C_tabname_rank,
			fullid,
			&obj,
		)
		if err != nil && !db.IsNotFound(err) {
			log.Warning("fetch rank from db failed:", tp, sgid, rkid, err)
		}

		// ok
		loop.Push(func() {
			rk := r_ranks[fullid] // MUST re-fetch it. because it may be replaced by updaters
			if rk.need_updates > 0 {
				rk.data = obj.D
				rk.need_updates--
			}

			for _, f := range rk.req {
				f(rk.data)
			}
			rk.req = nil
		})
	})
}

func TryGet(tp, sgid, rkid int32) (ret []*RankRow) {
	// find from cache only
	fullid := make_full_rankid(tp, sgid, rkid)
	rk := r_ranks[fullid]
	if rk == nil {
		return nil
	} else {
		return rk.data
	}
}

func GetRowPos(rk []*RankRow, id string) int32 {
	for i, v := range rk {
		if v.Info.Id() == id {
			return int32(i) + 1
		}
	}

	return 0
}

// ============================================================================

func make_full_rankid(tp int32, sgid int32, rkid int32) int32 {
	return tp*1_0000_0000 + sgid*1_0000 + rkid
}
