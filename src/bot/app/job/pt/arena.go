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

type data_arena_t struct {
	enemies []*msg.ArenaEnemy
}

// ============================================================================

func init() {
	evtmgr.On("userinfo", func(args ...interface{}) {
		bot := args[0].(*app.Bot)

		bot.SetData("pt_arena", &data_arena_t{})

		bot.JobAdd("pt_arena_fight", pt_arena_fight)
	})

	evtmgr.On(app.MsgEvt(&msg.GS_ArenaUpdateEnemy_R{}), func(args ...interface{}) {
		bot := args[0].(*app.Bot)
		res := args[1].(*msg.GS_ArenaUpdateEnemy_R)

		d := bot.GetData("pt_arena").(*data_arena_t)

		if len(res.Enemies) > 0 {
			d.enemies = res.Enemies
		}
	})
}

// ============================================================================

func pt_arena_fight(bot *app.Bot) {
	// update enemy
	bot.SendMsg(&msg.C_ArenaUpdateEnemy{})

	// wait
	time.Sleep(time.Second * time.Duration(core.RandInt(1, 2)))

	// fight
	d := bot.GetData("pt_arena").(*data_arena_t)
	lst := d.enemies
	L := len(lst)
	if L == 0 {
		return
	}
	e := lst[rand.Intn(L)]

	bot.SendMsg(&msg.C_ArenaFight{
		PlrId:      e.Plr.Id,
		IsOneKey:   false,
		RevengeIdx: -1,
		T:          team_rand_create(bot),
	})

	// pull replay
	p := rand.Float32()
	if p < 0.3 {
		bot.SendMsg(&msg.C_ArenaRecordInfo{})
	}

	// pull rank
	p = rand.Float32()
	if p < 0.3 {
		bot.SendMsg(&msg.C_ArenaRank{
			Top: 20,
		})
	}
}
