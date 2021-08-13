package c_gs

import (
	"fw/src/bot/app"
	"fw/src/bot/msg"
	"fw/src/core/log"
)

func GS_LoginError(message msg.Message, ctx interface{}) {
	req := message.(*msg.GS_LoginError)
	bot := ctx.(*app.Bot)

	log.Error("get uesrinfo failed:", req.ErrorCode, "botid:", bot.Id)
	bot.Close()
}
