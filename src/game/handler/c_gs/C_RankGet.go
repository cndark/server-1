package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/modules/rank"
	"fw/src/game/app/modules/ranksvc"
	"fw/src/game/msg"
)

func C_RankGet(message msg.Message, ctx interface{}) {
	req := message.(*msg.C_RankGet)
	plr := ctx.(*app.Player)

	res := &msg.GS_RankGet_R{}

	ranksvc.Get(req.Type, req.SgId, req.RkId, func(rows []*ranksvc.RankRow) {

		res.Rows = make([]*msg.RankRow, 0, len(rows))
		for _, v := range rows {
			// be liked count
			n := int32(0)

			c := rank.V_LikableRanks[req.RkId]
			if c != 0 && v.Info.Plr != nil {
				plr2 := app.PlayerMgr.LoadPlayer(v.Info.Plr.Id)
				if plr2 != nil {
					n = plr2.GetRankPlay().BeLikedCnt[req.RkId]
				}
			}

			res.Rows = append(res.Rows, &msg.RankRow{
				Score:   v.Score,
				Info:    (*msg.RankRowInfo)(v.Info),
				BeLiked: n,
			})
		}

		plr.SendMsg(res)
	})
}
