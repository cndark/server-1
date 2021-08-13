"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-26 21:16:25
 * @LastEditors: chenjie
 * @LastEditTime: 2021-06-26 18:47:26
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
exports.BFActExtraSpecial = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFActExtraSpecial = /** @class */ (function (_super) {
    __extends(BFActExtraSpecial, _super);
    function BFActExtraSpecial() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFActExtraSpecial.prototype.buff_actExtraSpecial = function (skill) {
        var owner = this._buff.getOwner();
        var conditions = this._funcParams[1].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "skillType") {
                if (skill.sType != condition[1])
                    return false;
            }
            else if (condition[0] === "buffType") {
                var t = (condition[1] === "self") ? [owner] : skill.getTargets();
                var bTypes = condition.slice(2).map(function (a) { return Number(a); });
                var isValid = false;
                for (var j = t.length - 1; j >= 0; j--) {
                    if (t[j].hasBuffType(bTypes)) {
                        isValid = true;
                        break;
                    }
                }
                if (!isValid)
                    return false;
            }
        }
        var ratio = Number(this._funcParams[0]);
        if (owner.bCtrl.random() >= ratio)
            return false;
        return true;
    };
    return BFActExtraSpecial;
}(BuffFunc_1.BuffFunc));
exports.BFActExtraSpecial = BFActExtraSpecial;
