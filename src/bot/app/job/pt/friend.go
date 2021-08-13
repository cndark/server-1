package pt

import (
	"fw/src/bot/app"
	"fw/src/bot/msg"
	"fw/src/core/evtmgr"
	"math/rand"
	"time"
)

// ============================================================================

type data_friend_t struct {
	frds []string
}

// ============================================================================

func init() {
	evtmgr.On("userinfo", func(args ...interface{}) {
		bot := args[0].(*app.Bot)

		bot.SetData("pt_friend", &data_friend_t{})

		bot.JobAdd("pt_friend_apply", pt_friend_apply)
		bot.JobAdd("pt_friend_interact", pt_friend_interact)
	})

	evtmgr.On(app.MsgEvt(&msg.GS_FriendNewApplied{}), func(args ...interface{}) {
		bot := args[0].(*app.Bot)
		res := args[1].(*msg.GS_FriendNewApplied)

		bot.SendMsg(&msg.C_FriendAccept{
			IsAccept: true,
			PlrIds:   []string{res.PlrId},
		})
	})

	evtmgr.On(app.MsgEvt(&msg.GS_FriendGetFrds_R{}), func(args ...interface{}) {
		bot := args[0].(*app.Bot)
		res := args[1].(*msg.GS_FriendGetFrds_R)

		arr := make([]string, 0, len(res.Friends))
		for _, v := range res.Friends {
			arr = append(arr, v.Info.Id)
		}

		d := bot.GetData("pt_friend").(*data_friend_t)

		d.frds = arr
	})
}

// ============================================================================

func pt_friend_apply(bot *app.Bot) {
	id := global_rand_plrid()
	if id == "" {
		return
	}

	bot.SendMsg(&msg.C_FriendApply{
		PlrId: id,
	})
}

func pt_friend_interact(bot *app.Bot) {
	// pull list
	bot.SendMsg(&msg.C_FriendGetFrds{})

	// wait
	time.Sleep(time.Millisecond * 1500)

	// get list
	d := bot.GetData("pt_friend").(*data_friend_t)
	lst := d.frds

	// pick one
	L := len(lst)
	if L == 0 {
		return
	}

	f := lst[rand.Intn(L)]

	bot.SendMsg(&msg.C_FriendGiveAndRecv{
		PlrIds: []string{f},
	})
}
