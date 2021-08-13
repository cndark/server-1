"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-26 14:13:11
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-19 20:22:09
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
exports.BFDotGamble = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var BuffFunc_1 = require("./BuffFunc");
var BFDotGamble = /** @class */ (function (_super) {
    __extends(BFDotGamble, _super);
    function BFDotGamble() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFDotGamble.prototype.add = function () {
        this.buff_dotGamble();
    };
    BFDotGamble.prototype.buff_dotGamble = function () {
        var atker = this._buff.getAtker();
        var owner = this._buff.getOwner();
        if (!atker.isAlive() || !owner.isAlive())
            return;
        var limit = atker.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE] * Number(this._funcParams[2]);
        var value = Math.min(owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE] * Number(this._funcParams[0]), limit);
        owner.updateHp(value, null, this._buff, false);
        value = Math.min(owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE] * Number(this._funcParams[1]), limit);
        atker.updateHp(value, null, this._buff, false);
        owner.bCtrl.buffCtrl.buff_dotExtraBuff(owner);
    };
    return BFDotGamble;
}(BuffFunc_1.BuffFunc));
exports.BFDotGamble = BFDotGamble;
