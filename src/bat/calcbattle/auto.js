"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-05-10 14:44:23
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-06-03 16:06:43
 */
Object.defineProperty(exports, "__esModule", { value: true });
var fs = require("fs");
var BattleConst_1 = require("./scripts/bundle0/battleCommon/BattleConst");
var SvrTeamCtrl_1 = require("./scripts/bundle0/battleCommon/SvrTeamCtrl");
var Configs_1 = require("./scripts/bundle0/framework/config/Configs");
var BattleCtrl_1 = require("./scripts/bundle1/battle/BattleCtrl");
var RecordCtrl_1 = require("./scripts/bundle1/battle/RecordCtrl");
Configs_1.Configs.loadCsv();
function startFight(input, speed) {
    var bCtrl = new BattleCtrl_1.BattleCtrl();
    bCtrl.recordCtrl = new RecordCtrl_1.RecordCtrl(bCtrl);
    bCtrl.init(input, null);
    bCtrl.stateStart();
    var dt = 0.016;
    var frame = 0;
    while (true) {
        bCtrl.update(dt * speed);
        frame++;
        if (bCtrl.state === BattleConst_1.BATTLE_STATE.COMPLETE)
            break;
    }
    console.log("complete", frame, bCtrl.atkCnt, bCtrl.defCnt);
    return bCtrl;
}
function checkRecord(r1, r2) {
    console.log(r1.length, r2.length);
    if (r1.length != r2.length) {
        return false;
    }
    for (var i = 0, j = r1.length; i < j; i++) {
        if (r1[i] != r2[i]) {
            console.error(r1[i], r2[i]);
            return false;
        }
    }
    return true;
}
function auto() {
    var input = {
        T1: { Fighters: [] },
        T2: { Fighters: [] },
        Args: {
            RoundType: "0",
            LvNum: "0",
            Module: BattleConst_1.BattleConst.BATTLE_TYPE.DEBUG,
            seed: String(Math.floor(Math.random() * 1000000000) + 1),
        }, // battle args. keys: [ "RoundType", "LvNum", "Module":["wlevel"] ]
    };
    var keys = Object.keys(Configs_1.Configs.monsterConf);
    var info = [];
    for (var i = 0; i < 6; i++) {
        var idx = keys[Math.floor(Math.random() * keys.length)];
        var mConf = Configs_1.Configs.monsterConf[idx];
        info.push({ lv: 1, id: mConf.id });
    }
    input.T1.Fighters = SvrTeamCtrl_1.SvrTeamCtrl.getFighterByConf(info);
    info = [];
    for (var i = 6; i < 12; i++) {
        var idx = keys[Math.floor(Math.random() * keys.length)];
        var mConf = Configs_1.Configs.monsterConf[idx];
        info.push({ lv: 1, id: mConf.id });
    }
    input.T2.Fighters = SvrTeamCtrl_1.SvrTeamCtrl.getFighterByConf(info);
    var b1 = startFight(input, 1);
    var b2 = startFight(input, 3);
    if (!checkRecord(b1.recordCtrl.content, b2.recordCtrl.content)) {
        fs.writeFileSync('input.txt', JSON.stringify(input));
        fs.writeFileSync('record1.txt', b1.recordCtrl.content.join("\n"));
        fs.writeFileSync('record2.txt', b2.recordCtrl.content.join("\n"));
        console.error("fatal! record is not equal!!!!");
    }
    else {
        setTimeout(auto, 100);
    }
}
auto();
