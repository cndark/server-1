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

		bot.JobAdd("pt_armor_compose_onekey", pt_armor_compose_onekey)
		bot.JobAdd("pt_armor_equip", pt_armor_equip)
	})
}

// ============================================================================

func pt_armor_compose_onekey(bot *app.Bot) {
	T := []int32{23101, 23201, 23301, 23401}
	L := len(T)

	bot.SendMsg(&msg.C_ArmorComposeOnekey{
		SrcId: T[core.RandInt(0, L-1)],
	})
}

func pt_armor_equip(bot *app.Bot) {
	// get hero
	hero := hero_rand_get(bot)
	if hero == nil {
		return
	}

	// equip
	for _, a := range []int32{23101, 23201, 23301, 23401} {
		bot.SendMsg(&msg.C_ArmorEquip{
			HeroSeq: hero.seq,
			Id:      core.RandInt32(a, a+17),
		})

		time.Sleep(time.Millisecond * 500)
	}
}
