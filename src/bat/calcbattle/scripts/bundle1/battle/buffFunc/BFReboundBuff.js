"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-21 14:58:32
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-28 15:04:11
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
exports.BFReboundBuff = void 0;
var Configs_1 = require("../../../bundle0/framework/config/Configs");
var BuffFunc_1 = require("./BuffFunc");
var BFReboundBuff = /** @class */ (function (_super) {
    __extends(BFReboundBuff, _super);
    function BFReboundBuff() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFReboundBuff.prototype.buff_reboundBuff = function (buffId, atker) {
        var bType = Number(this._funcParams[1]);
        var valid = false;
        var conf = Configs_1.Configs.buffConf[buffId];
        for (var i = 0; i < conf.buffType.length; i++) {
            if (bType === conf.buffType[i]) {
                valid = true;
                break;
            }
        }
        if (!valid)
            return false;
        var owner = this._buff.getOwner();
        var ratio = Number(this._funcParams[0]);
        if (owner.bCtrl.random() >= ratio)
            return false;
        owner.bCtrl.buffCtrl.addBuff(atker, buffId, atker);
        return true;
    };
    return BFReboundBuff;
}(BuffFunc_1.BuffFunc));
exports.BFReboundBuff = BFReboundBuff;
