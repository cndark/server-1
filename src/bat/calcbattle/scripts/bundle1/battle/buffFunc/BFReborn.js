"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-27 17:57:20
 * @LastEditors: zyb
 * @LastEditTime: 2021-07-13 11:40:26
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
exports.BFReborn = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFReborn = /** @class */ (function (_super) {
    __extends(BFReborn, _super);
    function BFReborn() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFReborn.prototype.buff_reborn = function (death) {
        var owner = this._buff.getOwner();
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i];
            if (condition === "self") {
                if (owner.seq !== death.seq)
                    return 0;
            }
            else if (condition === "teammate") {
                if (!owner.isAlive() || owner.seq === death.seq || owner.group !== death.group)
                    return 0;
            }
        }
        var ratio = Number(this._funcParams[1]);
        if (owner.bCtrl.random() >= ratio)
            return 0;
        return Number(this._funcParams[2]);
    };
    return BFReborn;
}(BuffFunc_1.BuffFunc));
exports.BFReborn = BFReborn;
