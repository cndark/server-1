"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-21 11:38:30
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-05-19 09:52:09
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
exports.BFAttrToShield = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFAttrToShield = /** @class */ (function (_super) {
    __extends(BFAttrToShield, _super);
    function BFAttrToShield() {
        var _this = _super !== null && _super.apply(this, arguments) || this;
        _this._shieldHp = 0;
        return _this;
    }
    BFAttrToShield.prototype.add = function () {
        this._shieldHp = 0;
        var owner = this._buff.getOwner();
        var atker = this._buff.getAtker();
        var target = null;
        var targetType = this._funcParams[0];
        if (targetType === "self") {
            target = owner;
        }
        else if (targetType === "origin") {
            target = atker;
        }
        if (!target)
            return;
        var attrId = Number(this._funcParams[1]);
        var ratio = Number(this._funcParams[2]);
        this._shieldHp = target.attrs[attrId] * ratio;
        owner.shieldHp += this._shieldHp;
        owner.bCtrl.callSceneFunc("updateHp", owner, this._shieldHp, null, this._buff, false);
    };
    BFAttrToShield.prototype.remove = function () {
        var owner = this._buff.getOwner();
        owner.shieldHp -= this._shieldHp;
    };
    BFAttrToShield.prototype.buff_attrToShield = function (dmg) {
        if (dmg <= 0)
            return dmg;
        var owner = this._buff.getOwner();
        if (dmg >= this._shieldHp) {
            owner.shieldHp -= this._shieldHp;
            dmg -= this._shieldHp;
            this._shieldHp = 0;
            owner.bCtrl.buffCtrl.removeBuff(this._buff);
        }
        else {
            owner.shieldHp -= dmg;
            this._shieldHp -= dmg;
            dmg = 0;
        }
        return dmg;
    };
    return BFAttrToShield;
}(BuffFunc_1.BuffFunc));
exports.BFAttrToShield = BFAttrToShield;
