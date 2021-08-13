package pt

import (
	"fw/src/bot/app"
	"fw/src/bot/msg"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"time"
)

// ============================================================================

func init() {
	evtmgr.On("userinfo", func(args ...interface{}) {
		bot := args[0].(*app.Bot)

		bot.JobAdd("pt_crusade_fight", pt_crusade_fight)
	})
}

// ============================================================================

func pt_crusade_fight(bot *app.Bot) {
	// get info
	bot.SendMsg(&msg.C_CrusadeGetInfo{})

	// wait
	time.Sleep(time.Second * time.Duration(core.RandInt(1, 1)))

	// fight
	bot.SendMsg(&msg.C_CrusadeFight{
		T: team_rand_create(bot),
	})
}
