{
    "version": "0.0.1",

    "pprof": {
        "enabled": false,
        "port":    8090
    },

    "pwd": "8YUYSmRmtKACv3TqD6GNdDvfyhEyYPUFMqX3Z6vgmwkX",

    "bots": [
        {
            "from":       {{add (mul (add svr -1) 30000) (mul (add batch -1) 10000) 1}},
            "to":         {{add (mul (add svr -1) 30000) (mul (add batch -1) 10000) 4000}},
            "addr":       "ws://10.101.0.115:2300{{svr}}/gate",
            "svr":        "game{{svr}}",
            "sdk":        "soda.pressure",
            "model":       "",
            "job_prefix": ["pt"],
            "job_itv":    [10, 20],
            "args":       "",
            "grp":        1
        },
        {
            "from":       {{add (mul (add svr -1) 30000) (mul (add batch -1) 10000) 4001}},
            "to":         {{add (mul (add svr -1) 30000) (mul (add batch -1) 10000) 7000}},
            "addr":       "ws://10.101.0.115:2300{{svr}}/gate",
            "svr":        "game{{svr}}",
            "sdk":        "soda.pressure",
            "model":       "",
            "job_prefix": ["pt"],
            "job_itv":    [2000, 4000],
            "args":       "",
            "grp":        1
        },
        {
            "from":       {{add (mul (add svr -1) 30000) (mul (add batch -1) 10000) 7001}},
            "to":         {{add (mul (add svr -1) 30000) (mul (add batch -1) 10000) 10000}},
            "addr":       "ws://10.101.0.115:2300{{svr}}/gate",
            "svr":        "game{{svr}}",
            "sdk":        "soda.pressure",
            "model":       "",
            "job_prefix": ["pt"],
            "job_itv":    [5000, 10000],
            "args":       "",
            "grp":        1
        }
    ],

    "jobctl": {
        "grp_binding": {
        },

        "sel_weight": {
            "pt_item_use": 10,
            "pt_draw_tp": 10,
            "pt_armor_compose_onekey": 10,
            "pt_armor_equip": 20,
            "pt_team_setdef": 20,
            "pt_golden_hand": 10,
            "pt_appoint_refresh": 10,
            "pt_relic_equip": 10,
            "pt_relic_eat": 2,
            "pt_tower_fight": 20,
            "pt_wlevel_fight": 50,
            "pt_wlevel_gjtake": 10,
            "pt_guild_create": 2,
            "pt_guild_join": 20,
            "pt_guild_sign": 10,
            "pt_guild_wish": 20,
            "pt_guild_harbor": 10,
            "pt_guild_order": 10,
            "pt_guild_tech": 20,
            "pt_guild_boss": 40,
            "pt_crusade_fight": 10,
            "pt_arena_fight": 5,
            "pt_chat_send": 2,
            "pt_friend_apply": 10,
            "pt_friend_interact": 10,
            "pt_hero_levelup": 100,
            "pt_hero_starup": 100,
            "pt_hero_reset": 1,
            "pt_hero_decompose": 1,
            "pt_hero_trinket": 100
        },

        "tri_prob": {
            "pt_item_use": 1,
            "pt_draw_tp": 1,
            "pt_armor_compose_onekey": 1,
            "pt_armor_equip": 1,
            "pt_team_setdef": 1,
            "pt_golden_hand": 1,
            "pt_appoint_refresh": 1,
            "pt_relic_equip": 1,
            "pt_relic_eat": 1,
            "pt_tower_fight": 1,
            "pt_wlevel_fight": 1,
            "pt_wlevel_gjtake": 1,
            "pt_guild_create": 1,
            "pt_guild_join": 1,
            "pt_guild_sign": 1,
            "pt_guild_wish": 1,
            "pt_guild_harbor": 1,
            "pt_guild_order": 1,
            "pt_guild_tech": 1,
            "pt_guild_boss": 1,
            "pt_crusade_fight": 1,
            "pt_arena_fight": 1,
            "pt_chat_send": 1,
            "pt_friend_apply": 1,
            "pt_friend_interact": 1,
            "pt_hero_levelup": 1,
            "pt_hero_starup": 1,
            "pt_hero_reset": 1,
            "pt_hero_decompose": 1,
            "pt_hero_trinket": 1
        }
    }
}
