"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-20 17:52:58
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-05-19 09:52:29
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
exports.BFDot = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var BuffFunc_1 = require("./BuffFunc");
var BFDot = /** @class */ (function (_super) {
    __extends(BFDot, _super);
    function BFDot() {
        var _this = _super !== null && _super.apply(this, arguments) || this;
        _this._dotCnt = 0;
        return _this;
    }
    BFDot.prototype.add = function () {
        this._dotCnt = 0;
        this.buff_dot();
    };
    BFDot.prototype.buff_dot = function () {
        var atker = this._buff.getAtker();
        var owner = this._buff.getOwner();
        var value = 0;
        var dType = Number(this._funcParams[0]);
        var tType = this._funcParams[1];
        if (tType === "selfMaxLife") {
            value = owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
        }
        else if (tType === "selfHp") {
            value = owner.hp;
        }
        else if (tType === "originAtk") {
            value = atker.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK];
        }
        var idx = 3 + this._dotCnt;
        if (this._funcParams.length <= idx)
            idx = this._funcParams.length - 1;
        value *= Number(this._funcParams[idx]);
        if (dType === BattleConst_1.DMG_TYPE.PHYIC) {
            value *= Math.max(1 + atker.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.DOT_DAM_UP] - owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.DOT_DAM_DOWN], 0);
        }
        else {
            value *= Math.max(atker.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.CURE_RATIO], 0) * Math.max(owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.HEALING_RATIO], 0);
        }
        var limitRatio = Number(this._funcParams[2]);
        value = Math.min(value, limitRatio * atker.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK]);
        owner.updateHp(value, null, this._buff, false);
        if (value > 0) {
            owner.bCtrl.buffCtrl.buff_dotExtraBuff(owner);
        }
        else if (value < 0) {
            owner.bCtrl.buffCtrl.buff_cureExtraBuff(owner);
        }
        this._dotCnt++;
    };
    return BFDot;
}(BuffFunc_1.BuffFunc));
exports.BFDot = BFDot;
