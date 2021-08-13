package lamp

import (
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/game/app/comp"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/guild"
)

// ============================================================================
// 跑马灯 Id

const (
	lamp_guild_create = 1001 // 家族成立
	lamp_guild_lv     = 1002 // 家族升级
	lamp_hero_star    = 1004 // 英雄升星
	lamp_hero_add     = 3001 // 获得英雄
)

// ============================================================================
// 新的处理函数添加在下面

var lamp_handlers = map[int32]*lamp_handler_t{

	lamp_guild_create: {

		gconst.Evt_GuildCreate,

		func(args []interface{}) {
			conf := gamedata.ConfLamp.Query(lamp_guild_create)
			if conf == nil {
				return
			}

			gld := args[0].(*guild.Guild)

			Add(conf.Type, conf.Id, map[string]string{
				"gldName": gld.GetName(),
			})
		},
	},

	lamp_guild_lv: {

		gconst.Evt_GuildLv,

		func(args []interface{}) {
			conf := gamedata.ConfLamp.Query(lamp_guild_lv)
			if conf == nil {
				return
			}

			gld := args[0].(*guild.Guild)
			lv := args[1].(int32)

			Add(conf.Type, conf.Id, map[string]string{
				"gldName": gld.GetName(),
				"gldLv":   core.I32toa(lv),
			})
		},
	},

	lamp_hero_star: {

		gconst.Evt_HeroStar,

		func(args []interface{}) {
			conf := gamedata.ConfLamp.Query(lamp_hero_star)
			if conf == nil {
				return
			}

			plr := args[0].(IPlayer)
			heroid := args[1].(int32)
			star := args[2].(int32)

			isSend := false
			for _, v := range conf.P2 {
				if conf.P1 == 1 {
					if v == star {
						isSend = true
						break
					}
				} else if conf.P1 == 2 {
					if star >= v {
						isSend = true
						break
					}
				} else if conf.P1 == 3 {
					if star <= v {
						isSend = true
						break
					}
				}
			}

			if isSend {
				Add(conf.Type, conf.Id, map[string]string{
					"player": plr.GetName(),
					"heroId": core.I32toa(heroid),
					"star":   core.I32toa(star),
				})
			}
		},
	},

	lamp_hero_add: {

		gconst.Evt_HeroAdd,

		func(args []interface{}) {
			conf := gamedata.ConfLamp.Query(lamp_hero_add)
			if conf == nil {
				return
			}

			plr := args[0].(IPlayer)

			heroes := args[1].((map[int64]*comp.Hero))

			for _, hId := range conf.P2 {

				for _, v := range heroes {

					if hId == v.Id {
						Add(conf.Type, conf.Id, map[string]string{
							"player": plr.GetName(),
							"heroId": core.I32toa(hId),
						})
					}

				}

			}
		},
	},

	// #todo:

}

// ============================================================================
// api
// ============================================================================

func handler_evt_on() {
	for _, item := range gamedata.ConfLamp.Items() {
		item := item

		h := lamp_handlers[item.Id]
		if h == nil {
			continue
		}

		evtmgr.On(h.Evt, func(args ...interface{}) {
			h.F(args)
		})
	}
}

type lamp_handler_t struct {
	Evt string
	F   func(args []interface{})
}

// ============================================================================
