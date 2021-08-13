package c_gs

import (
	"fw/src/core/evtmgr"
	"fw/src/core/wordsfilter"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/chat"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

func C_ChatSend(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_ChatSend)
	plr := ctx.(*app.Player)

	res := &msg.GS_ChatSend_R{}

	res.ErrorCode = func() int32 {
		L := len(req.Content)
		if L >= chat.C_Content_Max_Len || L <= 0 {
			return Err.Chat_ContentLenLimited
		}

		if req.Tp == gconst.C_ChatType_Cross {
			conf := gamedata.ConfGlobalPublic.Query(1)
			if conf == nil {
				return Err.Failed
			}

			op := plr.GetBag().NewOp(gconst.ObjFrom_Chat)
			for _, v := range conf.ChatCrsSvrCost {
				op.Dec(v.Id, v.N)
			}

			if ec := op.CheckEnough(); ec != Err.OK {
				return ec
			}

			op.Apply()
		}

		// filter
		content := wordsfilter.Filter(req.Content)

		ci := &chat.ChatOne{
			Tp:      req.Tp,
			FromId:  plr.GetId(),
			ToId:    req.ToPlrId,
			Content: content,
			Ts:      time.Now(),
		}

		toName := ""
		if ci.Tp == gconst.C_ChatType_Friend {
			if !plr.IsFriend(ci.ToId) {
				return Err.Friend_IsNotFriend
			}

			tplr := app.PlayerMgr.LoadPlayer(ci.ToId)
			if tplr == nil {
				return Err.Plr_NotLoad
			}

			toName = tplr.GetName()
			// save
			plr.GetFriend().AddChat(ci)
			tplr.GetFriend().AddChat(ci)

		} else {
			ec := chat.Add(plr, ci)
			if ec != Err.OK {
				return ec
			}
		}

		evtmgr.Fire(gconst.Evt_SendChat, plr, req.Tp, content, req.ToPlrId, toName)

		return Err.OK
	}()

	plr.SendMsg(res)
}
