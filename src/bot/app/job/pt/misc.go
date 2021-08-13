package pt

import (
	"fw/src/bot/app"
	"fw/src/bot/msg"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gconst"
	"math/rand"
	"time"
)

// ============================================================================

func init() {
	evtmgr.On("userinfo", func(args ...interface{}) {
		bot := args[0].(*app.Bot)

		bot.JobAdd("pt_item_use", pt_item_use)
		bot.JobAdd("pt_draw_tp", pt_draw_tp)
		bot.JobAdd("pt_team_setdef", pt_team_setdef)
		bot.JobAdd("pt_golden_hand", pt_golden_hand)
		bot.JobAdd("pt_appoint_refresh", pt_appoint_refresh)
	})
}

func team_rand_create(bot *app.Bot) *msg.TeamFormation {
	T := &msg.TeamFormation{
		Formation: make(map[int64]int32),
	}

	for i, seq := range hero_rand_seqs(bot, 6) {
		T.Formation[seq] = int32(i)
	}

	return T
}

// ============================================================================

func pt_item_use(bot *app.Bot) {
	// 合成英雄
	bot.SendMsg(&msg.C_ItemUse{
		Id: 21031,
		N:  20,
	})

	time.Sleep(time.Millisecond * time.Duration(core.RandInt(500, 1000)))

	bot.SendMsg(&msg.C_ItemUse{
		Id: 21041,
		N:  30,
	})

	time.Sleep(time.Millisecond * time.Duration(core.RandInt(500, 1000)))

	bot.SendMsg(&msg.C_ItemUse{
		Id: 21051,
		N:  50,
	})
}

func pt_draw_tp(bot *app.Bot) {
	T := []string{"normal", "senior", "water", "thunder", "fire", "blend"}
	L := len(T)

	for i := 0; i < 2; i++ {
		n := 1
		if rand.Float32() < 0.5 {
			n = 10
		}

		bot.SendMsg(&msg.C_DrawTp{
			Tp: T[core.RandInt(0, L-1)],
			N:  int32(n),
		})

		time.Sleep(time.Second * time.Duration(core.RandInt(1, 3)))
	}
}

func pt_team_setdef(bot *app.Bot) {
	bot.SendMsg(&msg.C_SetTeam{
		Tp: gconst.TeamType_Dfd,
		T:  team_rand_create(bot),
	})
}

func pt_golden_hand(bot *app.Bot) {
	for i := 0; i < 3; i++ {
		bot.SendMsg(&msg.C_MiscGoldenHand{})

		time.Sleep(time.Millisecond * time.Duration(300))
	}
}

func pt_appoint_refresh(bot *app.Bot) {
	bot.SendMsg(&msg.C_AppointRefresh{})
}
