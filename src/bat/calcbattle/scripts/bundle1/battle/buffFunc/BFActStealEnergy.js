"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-21 11:26:23
 * @LastEditors: chenjie
 * @LastEditTime: 2021-04-21 11:37:18
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
exports.BFActStealEnergy = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFActStealEnergy = /** @class */ (function (_super) {
    __extends(BFActStealEnergy, _super);
    function BFActStealEnergy() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFActStealEnergy.prototype.buff_actStealEnergy = function () {
        var owner = this._buff.getOwner();
        var units = owner.bCtrl.getUnits();
        var targets = [];
        for (var i = 0; i < units.length; i++) {
            var unit = units[i];
            if (unit.isAlive() && unit.group !== owner.group && unit.mp >= unit.maxMp) {
                targets.push(unit);
            }
        }
        if (targets.length <= 0)
            return;
        var ratio = Number(this._funcParams[0]);
        if (owner.bCtrl.random() >= ratio)
            return;
        var idx = Math.floor(owner.bCtrl.random() * targets.length);
        var target = targets[idx];
        var value = Number(this._funcParams[1]);
        target.updateMp(-1 * value);
        owner.updateMp(value);
    };
    return BFActStealEnergy;
}(BuffFunc_1.BuffFunc));
exports.BFActStealEnergy = BFActStealEnergy;
