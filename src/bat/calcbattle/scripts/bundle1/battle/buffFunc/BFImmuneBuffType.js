"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-21 10:29:49
 * @LastEditors: chenjie
 * @LastEditTime: 2021-04-25 18:26:11
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
exports.BFImmuneBuffType = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFImmuneBuffType = /** @class */ (function (_super) {
    __extends(BFImmuneBuffType, _super);
    function BFImmuneBuffType(buff, func, params) {
        var _this = _super.call(this, buff, func, params) || this;
        _this._bTypes = [];
        _this._bTypes = [];
        for (var i = 0; i < _this._funcParams.length; i++) {
            var bType = Number(_this._funcParams[i]);
            _this._bTypes.push(bType);
        }
        return _this;
    }
    BFImmuneBuffType.prototype.add = function () {
        var owner = this._buff.getOwner();
        owner.bCtrl.buffCtrl.removeBuffByBTypes(owner, this._bTypes);
        for (var i = 0; i < this._bTypes.length; i++) {
            var bType = this._bTypes[i];
            if (!owner.immuneBuffs[bType])
                owner.immuneBuffs[bType] = 0;
            owner.immuneBuffs[bType]++;
        }
    };
    BFImmuneBuffType.prototype.remove = function () {
        var owner = this._buff.getOwner();
        for (var i = 0; i < this._bTypes.length; i++) {
            var bType = this._bTypes[i];
            owner.immuneBuffs[bType]--;
        }
    };
    return BFImmuneBuffType;
}(BuffFunc_1.BuffFunc));
exports.BFImmuneBuffType = BFImmuneBuffType;
