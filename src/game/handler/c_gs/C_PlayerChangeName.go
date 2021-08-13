package c_gs

import (
	"fw/src/core/wordsfilter"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"strings"
)

const (
	C_player_name_maxlen = 24
)

func C_PlayerChangeName(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_PlayerChangeName)
	plr := ctx.(*app.Player)

	res := &msg.GS_PlayerChangeName_R{}

	res.ErrorCode = func() int32 {
		// check name
		req.Name = strings.TrimSpace(req.Name)

		if req.Name == "" {
			return Err.Plr_NameLen
		} else if req.Name == plr.GetName() {
			return Err.Plr_SameName
		} else if len(req.Name) > C_player_name_maxlen {
			return Err.Plr_NameLen
		}

		if wordsfilter.IsSensitive(req.Name) {
			return Err.Plr_InvalidName
		}

		// check cost
		op := plr.GetBag().NewOp(gconst.ObjFrom_PlrChangeName)

		// free for the 1st time
		is_free := plr.User().Misc.FreeRename
		if !is_free {
			conf := gamedata.ConfGlobalPublic.Query(1)
			for _, v := range conf.PlayerNameChange {
				op.Dec(v.Id, int64(v.N))
			}

			if ec := op.CheckEnough(); ec != Err.OK {
				return ec
			}
		}

		// change name
		plr.ChangeName(req.Name, func(b bool) {
			if b {
				if is_free {
					// NOT free any more
					plr.User().Misc.FreeRename = false
				} else {
					op.Apply()
				}

				res.ErrorCode = Err.OK
				res.Name = req.Name
			} else {
				res.ErrorCode = Err.Plr_DupName
			}

			plr.SendMsg(res)
		})

		return Err.OK
	}()

	if res.ErrorCode != Err.OK {
		plr.SendMsg(res)
	}
}
