"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-21 17:40:49
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-05-19 09:52:18
 */
var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (Object.prototype.hasOwnProperty.call(b, p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        if (typeof b !== "function" && b !== null)
            throw new TypeError("Class extends value " + String(b) + " is not a constructor or null");
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.BFDeathDmg = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var FindTargets_1 = require("../FindTargets");
var BuffFunc_1 = require("./BuffFunc");
var BFDeathDmg = /** @class */ (function (_super) {
    __extends(BFDeathDmg, _super);
    function BFDeathDmg() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFDeathDmg.prototype.buff_deathDmg = function (death) {
        var owner = this._buff.getOwner();
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "self") {
                if (owner.seq !== death.seq)
                    return;
            }
        }
        var ratio = Number(this._funcParams[3]);
        if (owner.bCtrl.random() >= ratio)
            return false;
        var targetType = this._funcParams[2];
        var info = targetType.split("_");
        var isFriendly = Number(info[0]);
        targetType = targetType.substr(2);
        var targets = FindTargets_1.FindTargets.getTargets(owner, isFriendly, targetType, owner.bCtrl.getUnits());
        if (targets.length <= 0)
            return;
        var value = 0;
        var dmgRatio = this._funcParams[4];
        if (dmgRatio === "maxLife") {
            value = owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
        }
        else if (dmgRatio === "selfAtk") {
            value = owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK];
        }
        value *= Number(this._funcParams[5]);
        var dType = this._funcParams[0];
        var isCure = (dType === "cure");
        if (isCure) {
            value *= Math.max(owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.CURE_RATIO], 0);
        }
        for (var i = 0; i < targets.length; i++) {
            var target = targets[i];
            if (isCure) {
                value *= Math.max(target.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.HEALING_RATIO], 0);
            }
            target.updateHp(value, null, this._buff, false);
        }
    };
    return BFDeathDmg;
}(BuffFunc_1.BuffFunc));
exports.BFDeathDmg = BFDeathDmg;
