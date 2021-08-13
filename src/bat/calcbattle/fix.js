"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-05-10 15:54:18
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-05-19 10:54:11
 */
Object.defineProperty(exports, "__esModule", { value: true });
var fs = require("fs");
var BattleConst_1 = require("./scripts/bundle0/battleCommon/BattleConst");
var Configs_1 = require("./scripts/bundle0/framework/config/Configs");
var BattleCtrl_1 = require("./scripts/bundle1/battle/BattleCtrl");
var RecordCtrl_1 = require("./scripts/bundle1/battle/RecordCtrl");
Configs_1.Configs.loadCsv();
function fix() {
    var path = "./input.txt";
    var content = fs.readFileSync(path);
    var input = JSON.parse(content.toString());
    var bCtrl = new BattleCtrl_1.BattleCtrl();
    bCtrl.recordCtrl = new RecordCtrl_1.RecordCtrl(bCtrl);
    bCtrl.init(input, null);
    bCtrl.stateStart();
    var dt = 0.016;
    var frame = 0;
    while (true) {
        bCtrl.update(dt);
        frame++;
        if (bCtrl.state === BattleConst_1.BATTLE_STATE.COMPLETE)
            break;
    }
    console.log("complete", frame, bCtrl.atkCnt, bCtrl.defCnt);
    fs.writeFileSync('fix.txt', bCtrl.recordCtrl.content.join("\n"));
}
fix();
