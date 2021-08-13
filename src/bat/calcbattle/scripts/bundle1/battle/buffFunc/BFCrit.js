"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-19 16:51:49
 * @LastEditors: chenjie
 * @LastEditTime: 2021-04-26 11:08:02
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
exports.BFCrit = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFCrit = /** @class */ (function (_super) {
    __extends(BFCrit, _super);
    function BFCrit() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFCrit.prototype.add = function () {
        var owner = this._buff.getOwner();
        owner.mustCrit++;
    };
    BFCrit.prototype.remove = function () {
        var owner = this._buff.getOwner();
        owner.mustCrit--;
    };
    BFCrit.prototype.buff_crit = function () {
        return Number(this._funcParams[0]);
    };
    return BFCrit;
}(BuffFunc_1.BuffFunc));
exports.BFCrit = BFCrit;
