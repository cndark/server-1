"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-15 17:01:35
 * @LastEditors: chenjie
 * @LastEditTime: 2021-04-15 18:16:39
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
exports.BFSteal = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFSteal = /** @class */ (function (_super) {
    __extends(BFSteal, _super);
    function BFSteal(buff, func, params) {
        var _this = _super.call(this, buff, func, params) || this;
        _this._target = null;
        _this._value = 0;
        _this._target = null;
        _this._value = 0;
        var owner = _this._buff.getOwner();
        if (params[0] === "pos") {
            var units = owner.bCtrl.getUnits();
            for (var i = 0; i < units.length; i++) {
                var unit = units[i];
                if (unit.isAlive() && unit.group !== owner.group && unit.order === owner.order) {
                    _this._target = unit;
                    break;
                }
            }
        }
        else {
            _this._target = _this._buff.getAtker();
        }
        if (_this._target) {
            var attrId = Number(_this._funcParams[1]);
            var ratio = Number(_this._funcParams[2]);
            var limitRatio = Number(_this._funcParams[3]);
            var val = owner.attrs[attrId] * ratio;
            var limit = _this._target.attrs[attrId] * limitRatio;
            _this._value = Math.min(val, limit);
        }
        return _this;
    }
    BFSteal.prototype.add = function () {
        if (this._target) {
            var owner = this._buff.getOwner();
            var attrId = Number(this._funcParams[1]);
            owner.updateBuffAttr(attrId, -1 * this._value);
            this._target.updateBuffAttr(attrId, this._value);
        }
    };
    BFSteal.prototype.remove = function () {
        if (this._target) {
            var owner = this._buff.getOwner();
            var attrId = Number(this._funcParams[1]);
            owner.updateBuffAttr(attrId, this._value);
            this._target.updateBuffAttr(attrId, -1 * this._value);
        }
    };
    return BFSteal;
}(BuffFunc_1.BuffFunc));
exports.BFSteal = BFSteal;
