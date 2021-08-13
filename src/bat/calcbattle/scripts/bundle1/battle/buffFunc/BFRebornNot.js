"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-27 18:10:32
 * @LastEditors: chenjie
 * @LastEditTime: 2021-04-27 18:13:09
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
exports.BFRebornNot = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFRebornNot = /** @class */ (function (_super) {
    __extends(BFRebornNot, _super);
    function BFRebornNot() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFRebornNot.prototype.add = function () {
        var owner = this._buff.getOwner();
        owner.rebornNot++;
    };
    BFRebornNot.prototype.remove = function () {
        var owner = this._buff.getOwner();
        owner.rebornNot--;
    };
    return BFRebornNot;
}(BuffFunc_1.BuffFunc));
exports.BFRebornNot = BFRebornNot;
