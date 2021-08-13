"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-26 15:56:56
 * @LastEditors: chenjie
 * @LastEditTime: 2021-04-26 16:05:08
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
exports.BFBuffGrpBuff = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFBuffGrpBuff = /** @class */ (function (_super) {
    __extends(BFBuffGrpBuff, _super);
    function BFBuffGrpBuff() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFBuffGrpBuff.prototype.buff_buffGrpBuff = function (group) {
        var grpId = Number(this._funcParams[0]);
        if (grpId !== group)
            return;
        var owner = this._buff.getOwner();
        var cnt = owner.buffGrpCnt[grpId] || 0;
        var grpCnt = Number(this._funcParams[1]);
        if (grpCnt !== cnt)
            return;
        for (var i = this._funcParams.length - 1; i >= 2; i--) {
            owner.bCtrl.buffCtrl.addBuff(owner, Number(this._funcParams[i]), owner);
        }
    };
    return BFBuffGrpBuff;
}(BuffFunc_1.BuffFunc));
exports.BFBuffGrpBuff = BFBuffGrpBuff;
