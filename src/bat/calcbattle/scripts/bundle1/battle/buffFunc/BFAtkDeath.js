"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-21 15:56:25
 * @LastEditors: chenjie
 * @LastEditTime: 2021-04-21 16:03:19
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
exports.BFAtkDeath = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFAtkDeath = /** @class */ (function (_super) {
    __extends(BFAtkDeath, _super);
    function BFAtkDeath() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFAtkDeath.prototype.buff_atkDeath = function (sType, hpRatio) {
        var owner = this._buff.getOwner();
        var limitRatio = Number(this._funcParams[2]);
        if (hpRatio >= limitRatio)
            return false;
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "skillType") {
                if (sType !== condition[1])
                    return false;
            }
        }
        var ratio = Number(this._funcParams[1]);
        if (owner.bCtrl.random() >= ratio)
            return false;
        return true;
    };
    return BFAtkDeath;
}(BuffFunc_1.BuffFunc));
exports.BFAtkDeath = BFAtkDeath;
