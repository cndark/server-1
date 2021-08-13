"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-26 14:30:17
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-20 16:29:48
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
exports.BFRoundChangeHp = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var BuffFunc_1 = require("./BuffFunc");
var BFRoundChangeHp = /** @class */ (function (_super) {
    __extends(BFRoundChangeHp, _super);
    function BFRoundChangeHp() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFRoundChangeHp.prototype.buff_roundChangeHp = function () {
        var owner = this._buff.getOwner();
        if (owner.bCtrl.atkCnt <= 0 || owner.bCtrl.defCnt <= 0)
            return;
        var units = owner.bCtrl.getUnits();
        var atkMinHp = null;
        var atkMaxHp = null;
        var defMinHp = null;
        var defMaxHp = null;
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive() || unit.order > 5)
                continue;
            if (unit.group === owner.group) {
                if (!atkMinHp || atkMinHp.hp > unit.hp)
                    atkMinHp = unit;
                if (!atkMaxHp || atkMaxHp.hp < unit.hp)
                    atkMaxHp = unit;
            }
            else {
                if (!defMinHp || defMinHp.hp > unit.hp)
                    defMinHp = unit;
                if (!defMaxHp || defMaxHp.hp < unit.hp)
                    defMaxHp = unit;
            }
        }
        if (!atkMinHp || !defMinHp || atkMinHp.hp > defMinHp.hp)
            return;
        var limit = Number(this._funcParams[0]);
        if (defMaxHp.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE] / atkMinHp.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE] > limit)
            return;
        var ratio = Number(this._funcParams[1]);
        if (owner.bCtrl.random() >= ratio)
            return;
        var minHp = atkMinHp.hp;
        var maxHp = defMaxHp.hp;
        atkMinHp.forceUpdateHp(maxHp);
        defMaxHp.forceUpdateHp(minHp);
    };
    return BFRoundChangeHp;
}(BuffFunc_1.BuffFunc));
exports.BFRoundChangeHp = BFRoundChangeHp;
