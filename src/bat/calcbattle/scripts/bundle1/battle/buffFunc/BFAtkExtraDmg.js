"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-15 20:45:11
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-19 20:53:10
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
exports.BFAtkExtraDmg = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var BuffFunc_1 = require("./BuffFunc");
var BFAtkExtraDmg = /** @class */ (function (_super) {
    __extends(BFAtkExtraDmg, _super);
    function BFAtkExtraDmg() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFAtkExtraDmg.prototype.buff_atkExtraDmg = function (target, sType, baseDmg, crit) {
        var owner = this._buff.getOwner();
        var ratio = Number(this._funcParams[0]);
        if (owner.bCtrl.random() >= ratio)
            return 0;
        var conditions = this._funcParams[1].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "job") {
                var t = (condition[1] === "self") ? owner : target;
                if (t.job !== Number(condition[2]))
                    return 0;
            }
            else if (condition[0] === "elem") {
                var t = (condition[1] === "self") ? owner : target;
                if (t.elem !== Number(condition[2]))
                    return 0;
            }
            else if (condition[0] === "buffType") {
                var t = (condition[1] === "self") ? owner : target;
                var exist = t.hasBuffType([Number(condition[2])]);
                if (!exist)
                    return 0;
            }
            else if (condition[0] === "hpPoint") {
                var t = (condition[1] === "self") ? owner : target;
                var hpRatio = t.hp / t.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
                var limit = Number(condition[2]);
                if (condition[2] === "lt" && hpRatio >= limit)
                    return 0;
                if (condition[2] === "gt" && hpRatio < limit)
                    return 0;
            }
            else if (condition[0] === "maxBuffGrp") {
                var t = (condition[1] === "self") ? owner : target;
                var group = Number(condition[2]);
                var cnt = t.buffGrpCnt[group] || 0;
                var limit = Number(condition[3]);
                if (cnt === 0 || cnt < limit)
                    return 0;
                var isDel = (condition[4] === "1");
                if (isDel)
                    owner.bCtrl.buffCtrl.removeBuffByGrp(t, group);
            }
            else if (condition[0] === "skillType") {
                if (sType !== condition[1])
                    return 0;
            }
            else if (condition[0] === "originLive") {
                var atker = this._buff.getAtker();
                if (!atker.isAlive())
                    return 0;
            }
            else if (condition[0] === "crit") {
                if (!crit)
                    return 0;
            }
        }
        var dmg = 0;
        var dType = this._funcParams[2];
        if (dType === "targetMaxLife")
            dmg = target.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
        else if (dType === "targetNowHp")
            dmg = target.hp;
        else if (dType === "selfMaxLife")
            dmg = owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
        else if (dType === "skillDmg")
            dmg = baseDmg;
        else if (dType === "attrAtk")
            dmg = owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK];
        var dmgRatio = Number(this._funcParams[3]);
        var atkRatio = Number(this._funcParams[4]);
        dmg = Math.min(dmg * dmgRatio, owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK] * atkRatio);
        return dmg;
    };
    return BFAtkExtraDmg;
}(BuffFunc_1.BuffFunc));
exports.BFAtkExtraDmg = BFAtkExtraDmg;
