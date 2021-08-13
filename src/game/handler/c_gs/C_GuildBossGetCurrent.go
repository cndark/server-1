package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_GuildBossGetCurrent(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_GuildBossGetCurrent)
	plr := ctx.(*app.Player)

	res := &msg.GS_GuildBossGetCurrent_R{}

	res.ErrorCode = func() int32 {
		// get guild
		gld := plr.GetGuild()
		if gld == nil {
			return Err.Guild_NotAMember
		}

		// get
		if req.OnlyNum {
			res.Info = &msg.GuildBossHistory{
				Num: gld.Boss.Cur.Num,
			}
		} else {
			res.Info = gld.Boss.Cur.ToMsg()
			res.HpLeft = gld.Boss.HpLeft
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
