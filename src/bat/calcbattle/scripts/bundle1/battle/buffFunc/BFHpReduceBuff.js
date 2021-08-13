"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-21 12:11:56
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-28 15:17:07
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
exports.BFHpReduceBuff = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var BuffFunc_1 = require("./BuffFunc");
var BFHpReduceBuff = /** @class */ (function (_super) {
    __extends(BFHpReduceBuff, _super);
    function BFHpReduceBuff() {
        var _this = _super !== null && _super.apply(this, arguments) || this;
        _this._reduceHp = 0;
        _this._reduceRatio = 0;
        return _this;
    }
    BFHpReduceBuff.prototype.add = function () {
        this._reduceHp = 0;
        this._reduceRatio = Number(this._funcParams[0]);
    };
    BFHpReduceBuff.prototype.buff_hpReduceBuff = function (reduceHp) {
        this._reduceHp += reduceHp;
        var owner = this._buff.getOwner();
        var value = owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
        var delta = value * this._reduceRatio;
        var cnt = Math.floor(this._reduceHp / delta);
        for (var i = 0; i < cnt; i++) {
            this._reduceHp -= delta;
            var targetType = this._funcParams[2];
            var targets = [];
            if (targetType === "self") {
                targets = [owner];
            }
            if (targets.length <= 0)
                continue;
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
        }
    };
    return BFHpReduceBuff;
}(BuffFunc_1.BuffFunc));
exports.BFHpReduceBuff = BFHpReduceBuff;
