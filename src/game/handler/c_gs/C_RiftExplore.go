package c_gs

import (
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/modules/rift"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
)

func C_RiftExplore(message msg.Message, ctx interface{}) {
	// req := message.(*msg.C_RiftExplore)
	plr := ctx.(*app.Player)

	res := &msg.GS_RiftExplore_R{}
	res.ErrorCode = func() int32 {
		conf := gamedata.ConfGlobalPublic.Query(1)
		if conf == nil {
			return Err.Failed
		}

		// 首次探索指定怪物
		if plr.GetRift().GetFirst() {

			m1 := gamedata.ConfRiftMonster.Query(conf.RiftFirstMonster)
			if m1 == nil {
				return Err.Failed
			}

			plr.GetRift().SetMonster(m1.Id, plr.GetLevel())

			plr.SendMsg(&msg.GS_RiftMonsterNew{
				Monster: &msg.RiftMonster{
					Id: m1.Id,
					Lv: plr.GetLevel(),
				},
			})
			plr.GetRift().SetFirst()
			return Err.OK
		}

		slt := make(map[int32]int32)
		for i, v := range conf.RiftEvent {
			if i == rift.C_RiftType_Box && !rift.BoxMgr.IsStart() {
				continue
			}

			slt[int32(i)] += v
		}

		plr.GetRift().CleanMonster()
		rift.MineMgr.CleanEmptyMine(plr)

		id := utils.PickWeightedMapId(slt)
		switch id {
		case rift.C_RiftType_Mine:
			rift.MineMgr.ExploreRiftMine(plr)

		case rift.C_RiftType_Monster:
			plr.GetRift().ExploreMonster()

		case rift.C_RiftType_Box:
			rift.BoxMgr.ExploreRiftBox(plr)

		default:
		}

		return Err.OK
	}()

	plr.SendMsg(res)
}
