"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-05-05 18:30:41
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-20 18:07:18
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.CountCtrl = void 0;
var BattleConst_1 = require("../../bundle0/battleCommon/BattleConst");
var CountCtrl = /** @class */ (function () {
    function CountCtrl(bCtrl) {
        this.bCtrl = null;
        this._dmg = {}; // 输出
        this._cure = {}; // 治疗
        this._resist = {}; // 承伤
        this.bCtrl = bCtrl;
        this._dmg = {};
        this._cure = {};
        this._resist = {};
    }
    CountCtrl.prototype.countDmg = function (atker, defer, value) {
        var key = atker.order + "." + atker.group;
        if (!this._dmg[key])
            this._dmg[key] = 0;
        this._dmg[key] += value;
        key = defer.order + "." + defer.group;
        if (!this._resist[key])
            this._resist[key] = 0;
        this._resist[key] += value;
    };
    CountCtrl.prototype.countCure = function (defer, value) {
        var key = defer.order + "." + defer.group;
        if (!this._cure[key])
            this._cure[key] = 0;
        this._cure[key] += value;
    };
    CountCtrl.prototype.getCountData = function () {
        var _a;
        //总伤害
        var totalDmg = (_a = {},
            _a[BattleConst_1.BattleConst.UNIT_GROUP.ATTACKER] = 0,
            _a[BattleConst_1.BattleConst.UNIT_GROUP.DEFENDER] = 0,
            _a);
        for (var key in this._dmg) {
            var val = Math.floor(this._dmg[key]);
            var group = Number(key.split(".")[1]);
            totalDmg[group] += val;
        }
        return {
            dmg: this._dmg,
            cure: this._cure,
            resist: this._resist,
            totalDmg: totalDmg,
        };
    };
    return CountCtrl;
}());
exports.CountCtrl = CountCtrl;
