package pt

import (
	"fw/src/bot/app"
	"fw/src/bot/msg"
	"fw/src/core/evtmgr"
)

// ============================================================================

type data_plr_t struct {
	id  string
	lv  int32
	gid string
}

// ============================================================================

func init() {
	evtmgr.On("userinfo", func(args ...interface{}) {
		bot := args[0].(*app.Bot)
		info := args[1].(*msg.GS_UserInfo)

		bot.SetData("pt_plr", &data_plr_t{
			id:  info.UserId,
			lv:  info.Lv,
			gid: info.GuildId,
		})

	})

	evtmgr.On(app.MsgEvt(&msg.GS_PlayerUpdateLv{}), func(args ...interface{}) {
		bot := args[0].(*app.Bot)
		res := args[1].(*msg.GS_PlayerUpdateLv)

		d := bot.GetData("pt_plr").(*data_plr_t)

		if res.Level > 0 {
			d.lv = res.Level
		}
	})
}

func plr_get_data(bot *app.Bot) *data_plr_t {
	return bot.GetData("pt_plr").(*data_plr_t)
}

// ============================================================================
