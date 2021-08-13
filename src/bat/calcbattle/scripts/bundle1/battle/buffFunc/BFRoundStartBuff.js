"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-20 16:12:55
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-28 15:18:58
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
exports.BFRoundStartBuff = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var FindTargets_1 = require("../FindTargets");
var BuffFunc_1 = require("./BuffFunc");
var BFRoundStartBuff = /** @class */ (function (_super) {
    __extends(BFRoundStartBuff, _super);
    function BFRoundStartBuff() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFRoundStartBuff.prototype.buff_roundStartBuff = function () {
        var owner = this._buff.getOwner();
        var validRound = Number(this._funcParams[2]);
        var curRound = this._buff.maxRound - this._buff.round + 1;
        if (curRound % validRound !== 0)
            return;
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "hpPoint") {
                var hpRatio = owner.hp / owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
                var limit = Number(condition[2]);
                if (condition[1] === "lt" && hpRatio >= limit)
                    return;
                if (condition[1] === "gt" && hpRatio < limit)
                    return;
            }
        }
        var targetType = this._funcParams[3];
        var info = targetType.split("_");
        var isFriendly = Number(info[0]);
        targetType = targetType.substr(2);
        var targets = FindTargets_1.FindTargets.getTargets(owner, isFriendly, targetType, owner.bCtrl.getUnits());
        if (targets.length <= 0)
            return;
        var num = Number(this._funcParams[4]);
        var buffIds = [];
        for (var i = 5; i < this._funcParams.length; i++) {
            buffIds.push(Number(this._funcParams[i]));
        }
        while (num < buffIds.length) {
            var idx = Math.floor(owner.bCtrl.random() * buffIds.length);
            buffIds.splice(idx, 1);
        }
        var ratio = Number(this._funcParams[1]);
        for (var i = 0; i < targets.length; i++) {
            var target = targets[i];
            for (var j = 0; j < buffIds.length; j++) {
                owner.bCtrl.buffCtrl.addBuff(target, buffIds[j], owner, ratio);
            }
        }
    };
    return BFRoundStartBuff;
}(BuffFunc_1.BuffFunc));
exports.BFRoundStartBuff = BFRoundStartBuff;
