package rank

import (
	"fmt"
	"fw/src/core"
	"fw/src/core/sched/loop"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/ranksvc"
	"fw/src/game/app/modules/utils"
	"fw/src/game/app/modules/worlddata"
	"fw/src/shared/config"
	"time"
)

// ============================================================================

type clone_plr_t struct {
	atkPwr   float64 // 最大战力
	wlevelLv float64 // 推图
	towerLv  float64 // 爬塔

	info *ranksvc.RankRowInfo
}

// ============================================================================

func Init() {
	// preload some ranks
	preload()

	// sched
	sched()
}

func preload() {
	f := func([]*ranksvc.RankRow) {}

	ranksvc.Get(ranksvc.RankType_Local, config.CurGame.Id, gconst.RankId_AtkPwr, f)
	ranksvc.Get(ranksvc.RankType_Local, config.CurGame.Id, gconst.RankId_WLevelLv, f)
	ranksvc.Get(ranksvc.RankType_Local, config.CurGame.Id, gconst.RankId_TowerLv, f)
}

func sched() {
	// every hour
	sched_every_hour(0)
	sched_every_hour(15)
	sched_every_hour(30)
	sched_every_hour(45)

	// first few times for server open
	Ms := []int{5, 10}

	t0 := worlddata.GetSvrCreateTs()
	now := time.Now()

	for _, M := range Ms {
		t := t0.Add(time.Minute * time.Duration(M))
		if t.After(now) {
			loop.SetTimeout(t, push)
		}
	}
}

func sched_every_hour(M int) {
	u, t := core.ParseRepeatTime(fmt.Sprintf("H/%d", M))

	var f func()
	f = func() {
		t = core.AddTimeByUnit(t, u, 1)
		loop.SetTimeout(t, func() {
			push()
			f()
		})
	}
	f()
}

func push() {
	// clone data
	cp := make([]*clone_plr_t, 0, utils.PlayerNumLoaded())

	// first, clone player pointers
	plrs := make([]IPlayer, 0, utils.PlayerNumLoaded())
	utils.ForEachLoadedPlayer(func(plr interface{}) {
		p := plr.(IPlayer)
		if !p.IsBan() {
			plrs = append(plrs, p)
		}
	})

	// then, clone with time distributed
	loop.TimeSlice(
		len(plrs), 500, 200,
		func(i int) {
			plr := plrs[i]

			// info
			info := &ranksvc.RankRowInfo{
				Plr: plr.ToMsg_SimpleInfo(),
			}

			// clone plr
			cp = append(cp, &clone_plr_t{
				atkPwr:   plr.GetAttainObjVal(gconst.AttainId_MaxAtkPwr),
				wlevelLv: float64(plr.GetWLevelLvNum()),
				towerLv:  float64(plr.GetTowerLvNum()),

				info: info,
			})

		},
		func() {
			core.Go(func() {
				// make rank-raws
				raws := make([]*ranksvc.RankRaw, 0, 3)

				// atk pwr
				raws = append(raws, &ranksvc.RankRaw{
					Id:        gconst.RankId_AtkPwr,
					A:         cp,
					ScoreFunc: func(i int) float64 { return cp[i].atkPwr },
					InfoFunc:  func(i int) *ranksvc.RankRowInfo { return cp[i].info },
					SortLevel: ranksvc.SortLevel_Local,
				})

				//  wlevel
				raws = append(raws, &ranksvc.RankRaw{
					Id:        gconst.RankId_WLevelLv,
					A:         cp,
					ScoreFunc: func(i int) float64 { return cp[i].wlevelLv },
					InfoFunc:  func(i int) *ranksvc.RankRowInfo { return cp[i].info },
					SortLevel: ranksvc.SortLevel_Local,
				})

				// tower
				raws = append(raws, &ranksvc.RankRaw{
					Id:        gconst.RankId_TowerLv,
					A:         cp,
					ScoreFunc: func(i int) float64 { return cp[i].towerLv },
					InfoFunc:  func(i int) *ranksvc.RankRowInfo { return cp[i].info },
					SortLevel: ranksvc.SortLevel_Local,
				})

				// push
				for _, raw := range raws {
					ranksvc.Push(raw)
					time.Sleep(time.Second * 5)
				}
			})
		},
	)

}
