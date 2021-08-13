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

		bot.JobAdd("pt_tower_fight", pt_tower_fight)
	})
}

// ============================================================================

func pt_tower_fight(bot *app.Bot) {
	bot.SendMsg(&msg.C_TowerFight{
		T: team_rand_create(bot),
	})

	time.Sleep(time.Second * time.Duration(core.RandInt(1, 3)))

	if rand.Float32() < 0.05 {
		bot.SendMsg(&msg.C_TowerRaid{})
	}
}
