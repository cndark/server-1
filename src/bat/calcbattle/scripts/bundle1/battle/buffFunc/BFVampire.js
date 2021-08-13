"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-16 16:58:47
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-13 15:22:44
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
exports.BFVampire = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var BuffFunc_1 = require("./BuffFunc");
var BFVampire = /** @class */ (function (_super) {
    __extends(BFVampire, _super);
    function BFVampire() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFVampire.prototype.buff_vampire = function (dmg, crit) {
        var owner = this._buff.getOwner();
        var hp = 0;
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "crit") {
                if (!crit)
                    return;
            }
        }
        var dType = this._funcParams[1];
        if (dType === "skillDmg") {
            hp = dmg;
        }
        else if (dType === "maxLife_self") {
            hp = owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
        }
        var ratio = Number(this._funcParams[2]);
        hp = hp *
            ratio *
            Math.max(owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.CURE_RATIO], 0) *
            Math.max(owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.HEALING_RATIO], 0) *
            -1;
        if (hp < 0)
            owner.updateHp(hp, null, this._buff, false);
    };
    return BFVampire;
}(BuffFunc_1.BuffFunc));
exports.BFVampire = BFVampire;
