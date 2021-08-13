"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-15 11:30:51
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-06-05 15:39:12
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.startBattle = void 0;
var BattleConst_1 = require("./scripts/bundle0/battleCommon/BattleConst");
var Configs_1 = require("./scripts/bundle0/framework/config/Configs");
var BattleCtrl_1 = require("./scripts/bundle1/battle/BattleCtrl");
Configs_1.Configs.loadCsv();
function startBattle(input) {
    "";
    var bCtrl = new BattleCtrl_1.BattleCtrl();
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
    var rt = {
        Winner: (bCtrl.atkCnt > 0 && bCtrl.defCnt <= 0) ? 1 : 2,
        Args: {},
    };
    switch (bCtrl.bType) {
        case BattleConst_1.BattleConst.BATTLE_TYPE.CRUSADE: //英灵试炼返回单位损失血量
            //args格式："hp_loss.group.pos": "decHp"
            bCtrl.getUnits().forEach(function (unit) {
                if (unit.group == rt.Winner) {
                    var decHp = Math.max(0, 1 - unit.hp / unit.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE]);
                    rt.Args["hp_loss." + unit.group + "." + unit.order] = String(decHp);
                }
                else { //失败方全部变成1
                    rt.Args["hp_loss." + unit.group + "." + unit.order] = "1";
                }
            });
            break;
        case BattleConst_1.BattleConst.BATTLE_TYPE.GUILD_BOSS: //公会boss返回boss损失血量
            bCtrl.getUnits().forEach(function (unit) {
                //只记录boss血量
                if (unit.group == BattleConst_1.BattleConst.UNIT_GROUP.DEFENDER) {
                    var decHp = Math.max(0, 1 - unit.hp / unit.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE]);
                    //公会boss损失血量不能比初始血量少
                    var initDecHp = input.Args["init_hp." + unit.group + "." + unit.order] || '0';
                    rt.Args["hp_loss." + unit.group + "." + unit.order] = String(Math.max(Number(initDecHp), decHp));
                }
            });
            break;
        default:
            break;
    }
    //统计数据
    var countData = bCtrl.countCtrl.getCountData();
    for (var key in countData.dmg) {
        var dmg = countData.dmg[key];
        rt.Args["dmg." + key] = String(Math.floor(dmg));
    }
    for (var key in countData.cure) {
        var cure = countData.cure[key];
        rt.Args["cure." + key] = String(Math.floor(cure));
    }
    for (var key in countData.resist) {
        var resist = countData.resist[key];
        rt.Args["resist." + key] = String(Math.floor(resist));
    }
    for (var key in countData.totalDmg) {
        var totalDmg = countData.totalDmg[key];
        rt.Args["dmg_total." + key] = String(Math.floor(totalDmg));
    }
    return rt;
}
exports.startBattle = startBattle;
