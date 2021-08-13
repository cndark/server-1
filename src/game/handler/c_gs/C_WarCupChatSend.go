package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/chat"
	"fw/src/game/app/modules/warcup"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_WarCupChatSend(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_WarCupChatSend)
	plr := ctx.(*app.Player)

	res := &msg.GS_WarCupChatSend_R{}
	res.ErrorCode = func() int32 {
		if !warcup.IsOpen() {
			return Err.Common_TimeNotUp
		}

		L := len(req.Content)
		if L >= chat.C_Content_Max_Len || L <= 0 {
			return Err.Chat_ContentLenLimited
		}

		warcup.WarCupChat(plr, req.Content)

		return Err.OK
	}()

	plr.SendMsg(res)
}
