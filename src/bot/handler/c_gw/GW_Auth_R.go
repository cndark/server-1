package c_gw

import (
	"fw/src/bot/app"
	"fw/src/bot/msg"
)

func GW_Auth_R(message msg.Message, ctx interface{}) {
	req := message.(*msg.GW_Auth_R)
	bot := ctx.(*app.Bot)

	bot.OnAuth(req)
}
