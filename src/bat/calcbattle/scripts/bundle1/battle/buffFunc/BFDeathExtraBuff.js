"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-20 15:53:47
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-28 15:15:16
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
exports.BFDeathExtraBuff = void 0;
var FindTargets_1 = require("../FindTargets");
var BuffFunc_1 = require("./BuffFunc");
var BFDeathExtraBuff = /** @class */ (function (_super) {
    __extends(BFDeathExtraBuff, _super);
    function BFDeathExtraBuff() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFDeathExtraBuff.prototype.buff_deathExtraBuff = function (death) {
        var owner = this._buff.getOwner();
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i];
            if (condition === "self") {
                if (owner.seq !== death.seq)
                    return;
            }
            else if (condition === "teammate") {
                if (owner.group !== death.group || owner.seq == death.seq)
                    return;
            }
            else if (condition === "teamAll") {
                if (owner.group != death.group)
                    return;
            }
            else if (condition === "enemy") {
                if (owner.group == death.group)
                    return;
            }
        }
        var targetType = this._funcParams[2];
        var info = targetType.split("_");
        var isFriendly = Number(info[0]);
        targetType = targetType.substr(2);
        var targets = FindTargets_1.FindTargets.getTargets(owner, isFriendly, targetType, owner.bCtrl.getUnits());
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
    return BFDeathExtraBuff;
}(BuffFunc_1.BuffFunc));
exports.BFDeathExtraBuff = BFDeathExtraBuff;
