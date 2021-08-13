"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-26 15:09:03
 * @LastEditors: chenjie
 * @LastEditTime: 2021-04-26 15:11:02
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
exports.BFAtkTargetBack = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFAtkTargetBack = /** @class */ (function (_super) {
    __extends(BFAtkTargetBack, _super);
    function BFAtkTargetBack() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFAtkTargetBack.prototype.add = function () {
        var owner = this._buff.getOwner();
        owner.atkTargetBack++;
    };
    BFAtkTargetBack.prototype.remove = function () {
        var owner = this._buff.getOwner();
        owner.atkTargetBack--;
    };
    return BFAtkTargetBack;
}(BuffFunc_1.BuffFunc));
exports.BFAtkTargetBack = BFAtkTargetBack;
