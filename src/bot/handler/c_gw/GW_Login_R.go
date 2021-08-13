package c_gw

import (
	"fw/src/bot/app"
	"fw/src/bot/msg"
)

func GW_Login_R(message msg.Message, ctx interface{}) {
	req := message.(*msg.GW_Login_R)
	bot := ctx.(*app.Bot)

	bot.OnLogin(req)
}
