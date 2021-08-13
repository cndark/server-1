"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-19 20:57:22
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-28 15:13:00
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
exports.BFAtkExtraBuff = void 0;
var FindTargets_1 = require("../FindTargets");
var BuffFunc_1 = require("./BuffFunc");
var BFAtkExtraBuff = /** @class */ (function (_super) {
    __extends(BFAtkExtraBuff, _super);
    function BFAtkExtraBuff() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFAtkExtraBuff.prototype.buff_atkExtraBuff = function (atker, defer, sType, crit, kill) {
        var owner = this._buff.getOwner();
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "buffType") {
                var t = (condition[1] === "atker") ? atker : defer;
                var exist = t.hasBuffType([Number(condition[2])]);
                if (!exist)
                    return;
            }
            else if (condition[0] === "buffGrp") {
                var t = (condition[1] === "atker") ? atker : defer;
                var group = Number(condition[2]);
                var cnt = t.buffGrpCnt[group] || 0;
                if (cnt <= 0)
                    return;
            }
            else if (condition[0] === "job") {
                var t = (condition[1] === "atker") ? atker : defer;
                if (t.job !== Number(condition[2]))
                    return 0;
            }
            else if (condition[0] === "skillType") {
                if (sType !== condition[1])
                    return 0;
            }
            else if (condition[0] === "atker") {
                if (condition[1] === "self" && atker.seq !== owner.seq)
                    return;
                if (condition[1] === "teammate" && (atker.seq === owner.seq || atker.group !== owner.group))
                    return;
            }
            else if (condition[0] === "defenser") {
                if (condition[1] === "self" && defer.seq !== owner.seq)
                    return;
                if (condition[1] === "teammate" && (defer.seq === owner.seq || defer.group !== owner.group))
                    return;
            }
            else if (condition[0] === "crit") {
                if (!crit)
                    return;
            }
            else if (condition[0] === "kill") {
                if (!kill)
                    return;
            }
        }
        var targetType = this._funcParams[2];
        var targets = [];
        if (targetType === "atker") {
            targets = [atker];
        }
        else if (targetType === "defenser") {
            targets = [defer];
        }
        else {
            var info = targetType.split("_");
            var isFriendly = Number(info[0]);
            targetType = targetType.substr(2);
            targets = FindTargets_1.FindTargets.getTargets(owner, isFriendly, targetType, owner.bCtrl.getUnits());
        }
        if (targets.length <= 0)
            return;
        var num = Number(this._funcParams[3]);
        var buffIds = [];
        for (var i = 4; i < this._funcParams.length; i++) {
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
    return BFAtkExtraBuff;
}(BuffFunc_1.BuffFunc));
exports.BFAtkExtraBuff = BFAtkExtraBuff;
