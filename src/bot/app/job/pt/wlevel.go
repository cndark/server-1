package pt

import (
	"fw/src/bot/app"
	"fw/src/bot/msg"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"math/rand"
	"time"
)

// ============================================================================

func init() {
	evtmgr.On("userinfo", func(args ...interface{}) {
		bot := args[0].(*app.Bot)

		bot.JobAdd("pt_wlevel_fight", pt_wlevel_fight)
		bot.JobAdd("pt_wlevel_gjtake", pt_wlevel_gjtake)
	})
}

// ============================================================================

func pt_wlevel_fight(bot *app.Bot) {
	T := team_rand_create(bot)

	bot.SendMsg(&msg.C_WLevelFight{
		T: T,
	})

	time.Sleep(time.Second * time.Duration(core.RandInt(1, 3)))

	if rand.Float32() < 0.05 {
		bot.SendMsg(&msg.C_WLevelFightOneKey{
			T: T,
		})
	}
}

func pt_wlevel_gjtake(bot *app.Bot) {
	bot.SendMsg(&msg.C_WLevelGJTake{})

	if rand.Float32() < 0.05 {
		bot.SendMsg(&msg.C_WLevelOneKeyGJTake{})
	}
}
