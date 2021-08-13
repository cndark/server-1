"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-15 11:30:51
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-06-03 16:06:24
 */
Object.defineProperty(exports, "__esModule", { value: true });
var _1 = require(".");
var BattleConst_1 = require("./scripts/bundle0/battleCommon/BattleConst");
var SvrTeamCtrl_1 = require("./scripts/bundle0/battleCommon/SvrTeamCtrl");
var Configs_1 = require("./scripts/bundle0/framework/config/Configs");
function test() {
    var input = {
        T1: { Fighters: [] },
        T2: { Fighters: [] },
        Args: {
            RoundType: "0",
            LvNum: "0",
            Module: BattleConst_1.BattleConst.BATTLE_TYPE.DEBUG,
        }, // battle args. keys: [ "RoundType", "LvNum", "Module":["wlevel"] ]
    };
    // let keys = Object.keys(Configs.monsterConf);
    var keys = [30133, 30178, 30141, 30064, 30136, 30060];
    var info = [];
    for (var i = 0; i < 6; i++) {
        var idx = keys[i];
        var mConf = Configs_1.Configs.monsterConf[idx];
        info.push({ lv: 1, id: mConf.id });
    }
    input.T1.Fighters = SvrTeamCtrl_1.SvrTeamCtrl.getFighterByConf(info);
    info = [];
    for (var i = 0; i < 6; i++) {
        var idx = keys[i];
        var mConf = Configs_1.Configs.monsterConf[idx];
        info.push({ lv: 1, id: mConf.id });
    }
    input.T2.Fighters = SvrTeamCtrl_1.SvrTeamCtrl.getFighterByConf(info);
    _1.startBattle(input);
}
test();
