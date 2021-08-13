package pt

import (
	"fmt"
	"fw/src/bot/app"
	"fw/src/bot/msg"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"math/rand"
	"sync"
	"time"
)

// ============================================================================

var (
	v_gldlist  []*msg.GuildRow
	v_wishlist = make(map[string][]*msg.GuildWish) // [gid]

	v_wishlist_lock sync.Mutex
)

// ============================================================================

type data_guild_t struct {
	orders []*msg.GuildOrderRec
}

// ============================================================================

func init() {
	evtmgr.On("userinfo", func(args ...interface{}) {
		bot := args[0].(*app.Bot)

		bot.SetData("pt_guild", &data_guild_t{})

		bot.JobAdd("pt_guild_create", pt_guild_create)
		bot.JobAdd("pt_guild_join", pt_guild_join)
		bot.JobAdd("pt_guild_sign", pt_guild_sign)
		bot.JobAdd("pt_guild_wish", pt_guild_wish)
		bot.JobAdd("pt_guild_harbor", pt_guild_harbor)
		bot.JobAdd("pt_guild_order", pt_guild_order)
		bot.JobAdd("pt_guild_tech", pt_guild_tech)
		bot.JobAdd("pt_guild_boss", pt_guild_boss)
	})

	evtmgr.On(app.MsgEvt(&msg.GS_GuildList_R{}), func(args ...interface{}) {
		res := args[1].(*msg.GS_GuildList_R)

		if len(res.Rows) > 0 {
			v_gldlist = res.Rows
		}
	})

	evtmgr.On(app.MsgEvt(&msg.GS_Guild_Join{}), func(args ...interface{}) {
		bot := args[0].(*app.Bot)
		res := args[1].(*msg.GS_Guild_Join)

		if res.Mb.Plr.Id == plr_get_data(bot).id {
			plr_get_data(bot).gid = res.GuildId
		}
	})

	evtmgr.On(app.MsgEvt(&msg.GS_GuildWishList_R{}), func(args ...interface{}) {
		bot := args[0].(*app.Bot)
		res := args[1].(*msg.GS_GuildWishList_R)

		if len(res.Wishes) > 0 {
			v_wishlist_lock.Lock()
			v_wishlist[plr_get_data(bot).gid] = res.Wishes
			v_wishlist_lock.Unlock()
		}
	})

	evtmgr.On(app.MsgEvt(&msg.GS_GuildOrderList_R{}), func(args ...interface{}) {
		bot := args[0].(*app.Bot)
		res := args[1].(*msg.GS_GuildOrderList_R)

		d := bot.GetData("pt_guild").(*data_guild_t)

		d.orders = res.Records
	})
}

// ============================================================================

func pt_guild_create(bot *app.Bot) {
	bot.SendMsg(&msg.C_GuildCreate{
		Name:   fmt.Sprintf("大团团-%d", rand.Int31()),
		Icon:   1101,
		Notice: fmt.Sprintf("带你飞-%d", rand.Int31()),
	})
}

func pt_guild_join(bot *app.Bot) {
	// pull list
	bot.SendMsg(&msg.C_GuildList{
		Page: core.RandInt32(1, 50),
	})

	// get list
	time.Sleep(time.Second * time.Duration(core.RandInt(1, 2)))
	lst := v_gldlist
	L := len(lst)

	// pick one
	if L == 0 {
		return
	}
	gld := lst[rand.Intn(L)]

	bot.SendMsg(&msg.C_GuildApply{
		GuildId: gld.Id,
	})
}

func pt_guild_sign(bot *app.Bot) {
	bot.SendMsg(&msg.C_GuildSign{})
}

func pt_guild_wish(bot *app.Bot) {
	// pull list
	bot.SendMsg(&msg.C_GuildWishList{})

	// get list
	time.Sleep(time.Second * time.Duration(core.RandInt(1, 2)))
	v_wishlist_lock.Lock()
	lst := v_wishlist[plr_get_data(bot).gid]
	v_wishlist_lock.Unlock()
	L := len(lst)

	// do
	p := rand.Float32()

	if p < 0.3 { // wish
		conf := gamedata.ConfGlobalPublic.Query(1)
		if conf != nil {
			bot.SendMsg(&msg.C_GuildWishItem{
				Num: int32(rand.Intn(len(conf.GuildWish))),
			})
		}
	} else if p < 0.7 { // close
		for _, v := range lst {
			if v.PlrId == plr_get_data(bot).id && v.Helps >= 3 {
				bot.SendMsg(&msg.C_GuildWishClose{
					Seq: v.Seq,
				})
				break
			}
		}
	} else { // help
		// pick one
		if L == 0 {
			return
		}
		w := lst[rand.Intn(L)]

		bot.SendMsg(&msg.C_GuildWishHelp{
			Seq: w.Seq,
		})
	}
}

func pt_guild_harbor(bot *app.Bot) {
	// get list
	bot.SendMsg(&msg.C_GuildHarborDonateList{})

	// wait
	time.Sleep(time.Second * time.Duration(core.RandInt(1, 2)))

	// donate
	bot.SendMsg(&msg.C_GuildHarborDonate{
		Which: core.RandInt32(1, 10),
	})
}

func pt_guild_order(bot *app.Bot) {
	// get orders
	bot.SendMsg(&msg.C_GuildOrderGet{})

	// wait
	time.Sleep(time.Second * time.Duration(core.RandInt(1, 2)))

	// pull list
	bot.SendMsg(&msg.C_GuildOrderList{})

	// wait
	time.Sleep(time.Second * time.Duration(core.RandInt(1, 2)))

	// get list
	d := bot.GetData("pt_guild").(*data_guild_t)
	lst := d.orders

	// pick one
	L := len(lst)
	if L == 0 {
		return
	}
	o := lst[rand.Intn(L)]

	// do
	p := rand.Float32()

	if p < 0.4 { // star up
		bot.SendMsg(&msg.C_GuildOrderStarup{
			Seq: o.Seq,
		})
	} else if p < 0.7 { // start
		bot.SendMsg(&msg.C_GuildOrderStart{
			Seq: o.Seq,
		})
	} else { // close
		bot.SendMsg(&msg.C_GuildOrderClose{
			Seq: o.Seq,
		})
	}
}

func pt_guild_tech(bot *app.Bot) {
	// pick elem
	e := rand.Int31n(5)

	start := e*8 + 1
	end := e*8 + 8

	for i := start; i <= end; i++ {
		bot.SendMsg(&msg.C_GuildTechLevelup{
			Id: i,
		})

		time.Sleep(time.Millisecond * time.Duration(core.RandInt(300, 500)))
	}

	// reset
	if rand.Float32() < 0.0001 {
		bot.SendMsg(&msg.C_GuildTechReset{})
	}
}

func pt_guild_boss(bot *app.Bot) {
	// get history
	bot.SendMsg(&msg.C_GuildBossGetHistory{
		Num: int32(core.RandInt(1, 10)),
	})

	// fight
	bot.SendMsg(&msg.C_GuildBossFight{
		Team: team_rand_create(bot),
	})
}
