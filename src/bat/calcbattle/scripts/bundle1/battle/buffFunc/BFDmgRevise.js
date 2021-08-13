"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-19 20:18:12
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-20 16:34:09
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
exports.BFDmgRevise = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFDmgRevise = /** @class */ (function (_super) {
    __extends(BFDmgRevise, _super);
    function BFDmgRevise() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFDmgRevise.prototype.buff_dmgRevise = function (atker, defer, baseDmg, crit) {
        var ret = { extra: 0, percent: 1 };
        var owner = this._buff.getOwner();
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "job") {
                var t = (condition[1] === "atker") ? atker : defer;
                var valid = false;
                for (var j = 2; j < condition.length; j++) {
                    if (t.job == Number(condition[j])) {
                        valid = true;
                        break;
                    }
                }
                if (!valid)
                    return ret;
            }
            else if (condition[0] === "elem") {
                var t = (condition[1] === "atker") ? atker : defer;
                if (t.elem !== Number(condition[2]))
                    return ret;
            }
            else if (condition[0] === "atker") {
                if (atker.seq !== owner.seq)
                    return ret;
            }
            else if (condition[0] === "defenser") {
                if (defer.seq !== owner.seq)
                    return ret;
            }
            else if (condition[0] === "crit") {
                if (!crit)
                    return ret;
            }
            else if (condition[0] === "buffType") {
                var t = (condition[1] === "atker") ? atker : defer;
                var exist = t.hasBuffType([Number(condition[2])]);
                if (!exist)
                    return ret;
            }
        }
        var ratio = Number(this._funcParams[1]);
        if (owner.bCtrl.random() >= ratio)
            return ret;
        var dType = this._funcParams[2];
        var dmgRatio = Number(this._funcParams[3]);
        if (dType === "skillDmg") {
            ret.extra = dmgRatio * baseDmg;
        }
        else {
            ret.percent = dmgRatio;
        }
        var isDel = (this._funcParams[4] === "1");
        if (isDel) {
            owner.bCtrl.buffCtrl.removeBuff(this._buff);
        }
        return ret;
    };
    return BFDmgRevise;
}(BuffFunc_1.BuffFunc));
exports.BFDmgRevise = BFDmgRevise;
