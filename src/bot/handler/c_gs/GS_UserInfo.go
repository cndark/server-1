package c_gs

import (
	"fw/src/bot/app"
	"fw/src/bot/msg"
)

func GS_UserInfo(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_UserInfo)
	bot := ctx.(*app.Bot)

	bot.OnUserInfo(req)
}
