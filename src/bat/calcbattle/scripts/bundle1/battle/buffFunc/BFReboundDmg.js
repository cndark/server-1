"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-21 15:18:50
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-05 18:27:33
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
exports.BFReboundDmg = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var BuffFunc_1 = require("./BuffFunc");
var BFReboundDmg = /** @class */ (function (_super) {
    __extends(BFReboundDmg, _super);
    function BFReboundDmg() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFReboundDmg.prototype.buff_reboundDmg = function (atker, delta) {
        var owner = this._buff.getOwner();
        var ratio = Number(this._funcParams[0]);
        if (owner.bCtrl.random() >= ratio)
            return;
        var dmgRatio = Number(this._funcParams[1]);
        delta *= dmgRatio;
        var limitRatio = Number(this._funcParams[2]);
        delta = Math.min(delta, owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK] * limitRatio);
        atker.updateHp(delta, null, this._buff, false);
    };
    return BFReboundDmg;
}(BuffFunc_1.BuffFunc));
exports.BFReboundDmg = BFReboundDmg;
