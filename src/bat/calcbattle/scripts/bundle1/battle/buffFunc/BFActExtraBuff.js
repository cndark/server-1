"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-20 16:52:21
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-28 15:10:55
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
exports.BFActExtraBuff = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var BuffFunc_1 = require("./BuffFunc");
var BFActExtraBuff = /** @class */ (function (_super) {
    __extends(BFActExtraBuff, _super);
    function BFActExtraBuff() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFActExtraBuff.prototype.buff_actExtraBuff = function () {
        var owner = this._buff.getOwner();
        var cnt = 1;
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "buffType") {
                var bTypes = [Number(condition[1])];
                if (!owner.hasBuffType(bTypes))
                    return;
            }
            else if (condition[0] === "hpPoint") {
                var hpRatio = owner.hp / owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
                var limit = Number(condition[2]);
                if (condition[1] === "lt" && hpRatio >= limit)
                    return;
                if (condition[1] === "gt" && hpRatio < limit)
                    return;
            }
            else if (condition[0] === "elemLive") {
                var elem = Number(condition[1]);
                var units = owner.bCtrl.getUnits();
                cnt = 0;
                for (var j = 0; j < units.length; j++) {
                    var unit = units[j];
                    if (unit.isAlive() && unit.elem === elem && unit.group === owner.group)
                        cnt++;
                }
            }
        }
        var targetType = this._funcParams[2];
        var targets = [];
        if (targetType === "self") {
            targets = [owner];
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
        for (var i = 0; i < cnt; i++) {
            for (var j = 0; j < targets.length; j++) {
                var target = targets[j];
                for (var k = 0; k < buffIds.length; k++) {
                    owner.bCtrl.buffCtrl.addBuff(target, buffIds[k], owner, ratio);
                }
            }
        }
    };
    return BFActExtraBuff;
}(BuffFunc_1.BuffFunc));
exports.BFActExtraBuff = BFActExtraBuff;
