"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-26 18:20:03
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-05-19 09:51:15
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
exports.BFActNewSkill = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var BuffFunc_1 = require("./BuffFunc");
var BFActNewSkill = /** @class */ (function (_super) {
    __extends(BFActNewSkill, _super);
    function BFActNewSkill() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFActNewSkill.prototype.buff_actNewSkill = function (sType) {
        var owner = this._buff.getOwner();
        var skillId = 0;
        var ratio = Number(this._funcParams[1]);
        if (owner.bCtrl.random() >= ratio)
            return skillId;
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "normal") {
                if (sType !== BattleConst_1.BattleConst.SKILL_TYPE.NORMAL)
                    return skillId;
            }
            else if (condition[0] === "special") {
                if (sType !== BattleConst_1.BattleConst.SKILL_TYPE.SPECIAL)
                    return skillId;
            }
            else if (condition[0] === "buffGrp") {
                var grpId = Number(condition[1]);
                var cnt = owner.buffGrpCnt[grpId] || 0;
                var grpCnt = Number(condition[2]);
                if (cnt === 0 || cnt < grpCnt)
                    return skillId;
                var isDel = (condition[3] == "1");
                owner.bCtrl.buffCtrl.removeBuffByGrp(owner, grpId);
            }
        }
        skillId = Number(this._funcParams[2]);
        return skillId;
    };
    return BFActNewSkill;
}(BuffFunc_1.BuffFunc));
exports.BFActNewSkill = BFActNewSkill;
