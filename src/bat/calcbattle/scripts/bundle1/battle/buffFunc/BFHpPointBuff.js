"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-21 14:13:46
 * @LastEditors: chenjie
 * @LastEditTime: 2021-06-28 14:26:48
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
exports.BFHpPointBuff = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var FindTargets_1 = require("../FindTargets");
var BuffFunc_1 = require("./BuffFunc");
var BFHpPointBuff = /** @class */ (function (_super) {
    __extends(BFHpPointBuff, _super);
    function BFHpPointBuff() {
        var _this = _super !== null && _super.apply(this, arguments) || this;
        _this._hpRatio = 0;
        return _this;
    }
    BFHpPointBuff.prototype.add = function () {
        this._hpRatio = Number(this._funcParams[0]);
    };
    BFHpPointBuff.prototype.buff_hpPointBuff = function (reduceHp) {
        var owner = this._buff.getOwner();
        var preHp = owner.hp + reduceHp;
        var value = owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
        if (preHp / value <= this._hpRatio || owner.hp / value > this._hpRatio)
            return;
        var targetType = this._funcParams[2];
        var targets = [];
        if (targetType === "self") {
            targets = [owner];
        }
        else {
            var info = targetType.split("_");
            var isFriendly = Number(info[0]);
            targetType = targetType.substr(2);
            targets = FindTargets_1.FindTargets.getTargets(owner, isFriendly, targetType, owner.bCtrl.getUnits());
        }
        if (targets.length <= 0)
            return;
        var num = Number(this._funcParams[3]);
        var buffIds = [];
        for (var j = 4; j < this._funcParams.length; j++) {
            buffIds.push(Number(this._funcParams[j]));
        }
        while (num < buffIds.length) {
            var idx = Math.floor(owner.bCtrl.random() * buffIds.length);
            buffIds.splice(idx, 1);
        }
        var ratio = Number(this._funcParams[1]);
        for (var j = 0; j < targets.length; j++) {
            var target = targets[j];
            for (var k = 0; k < buffIds.length; k++) {
                owner.bCtrl.buffCtrl.addBuff(target, buffIds[k], owner, ratio);
            }
        }
    };
    return BFHpPointBuff;
}(BuffFunc_1.BuffFunc));
exports.BFHpPointBuff = BFHpPointBuff;
