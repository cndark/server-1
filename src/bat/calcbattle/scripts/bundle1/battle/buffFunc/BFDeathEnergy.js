"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-21 10:58:22
 * @LastEditors: chenjie
 * @LastEditTime: 2021-04-21 11:03:12
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
exports.BFDeathEnergy = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFDeathEnergy = /** @class */ (function (_super) {
    __extends(BFDeathEnergy, _super);
    function BFDeathEnergy() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFDeathEnergy.prototype.buff_deathEnergy = function (death) {
        var owner = this._buff.getOwner();
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "enemy") {
                if (death.group === owner.group)
                    return;
            }
        }
        var ratio = Number(this._funcParams[1]);
        if (owner.bCtrl.random() >= ratio)
            return;
        var targets = [];
        var targetType = this._funcParams[2];
        if (targetType === "self") {
            targets = [owner];
        }
        if (targets.length <= 0)
            return;
        var value = Number(this._funcParams[3]);
        for (var i = 0; i < targets.length; i++) {
            var target = targets[i];
            target.updateMp(value);
        }
    };
    return BFDeathEnergy;
}(BuffFunc_1.BuffFunc));
exports.BFDeathEnergy = BFDeathEnergy;
