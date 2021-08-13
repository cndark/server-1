"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-16 17:53:58
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-05-19 09:51:42
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
exports.BFAtkEnergy = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var BuffFunc_1 = require("./BuffFunc");
var BFAtkEnergy = /** @class */ (function (_super) {
    __extends(BFAtkEnergy, _super);
    function BFAtkEnergy() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFAtkEnergy.prototype.buff_atkEnergy = function (atker, defer, dType, crit) {
        var owner = this._buff.getOwner();
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "atker") {
                if (condition[1] === "self" && atker.seq !== owner.seq)
                    return;
                if (condition[1] === "teammate" && (atker.seq === owner.seq || atker.group !== owner.group))
                    return;
            }
            else if (condition[0] === "defenser") {
                if (condition[1] === "self" && defer.seq !== owner.seq)
                    return;
                if (condition[1] === "teammate" && (defer.seq === owner.seq || defer.group !== owner.group))
                    return;
            }
            else if (condition[0] === "dmgType") {
                if (condition[1] === "dmg" && dType === BattleConst_1.DMG_TYPE.CURE)
                    return;
                if (condition[1] === "cure" && dType !== BattleConst_1.DMG_TYPE.CURE)
                    return;
            }
            else if (condition[0] === "crit") {
                if (!crit)
                    return;
            }
        }
        var ratio = Number(this._funcParams[1]);
        if (owner.bCtrl.random() >= ratio)
            return;
        var t = this._funcParams[2];
        var target = null;
        if (t === "self")
            target = owner;
        else if (t === "target")
            target = defer;
        else if (t === "randomTeam") {
            var units = owner.bCtrl.getUnits();
            var tmp = [];
            for (var j = 0; j < units.length; j++) {
                var unit = units[j];
                if (!unit.isAlive())
                    continue;
                if (unit.seq !== owner.seq && unit.group === owner.group) {
                    tmp.push(unit);
                }
            }
            if (tmp.length > 0) {
                var idx = Math.floor(owner.bCtrl.random() * tmp.length);
                target = tmp[idx];
            }
        }
        if (target) {
            var value = Number(this._funcParams[3]);
            target.updateMp(value);
        }
    };
    return BFAtkEnergy;
}(BuffFunc_1.BuffFunc));
exports.BFAtkEnergy = BFAtkEnergy;
