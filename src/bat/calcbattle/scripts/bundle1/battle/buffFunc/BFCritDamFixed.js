"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-17 16:51:49
 * @LastEditors: chenjie
 * @LastEditTime: 2021-04-17 12:57:20
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
exports.BFCritDamFixed = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFCritDamFixed = /** @class */ (function (_super) {
    __extends(BFCritDamFixed, _super);
    function BFCritDamFixed(buff, func, params) {
        var _this = _super.call(this, buff, func, params) || this;
        _this._critDmg = 0;
        _this._critDmg = Number(_this._funcParams[0]);
        return _this;
    }
    BFCritDamFixed.prototype.buff_critDamFixed = function () {
        return this._critDmg;
    };
    return BFCritDamFixed;
}(BuffFunc_1.BuffFunc));
exports.BFCritDamFixed = BFCritDamFixed;
