"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-27 17:10:06
 * @LastEditors: chenjie
 * @LastEditTime: 2021-04-27 17:14:36
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
exports.BFHitExtraNormal = void 0;
var BuffFunc_1 = require("./BuffFunc");
var BFHitExtraNormal = /** @class */ (function (_super) {
    __extends(BFHitExtraNormal, _super);
    function BFHitExtraNormal() {
        return _super !== null && _super.apply(this, arguments) || this;
    }
    BFHitExtraNormal.prototype.buff_hitExtraNormal = function (defer) {
        var owner = this._buff.getOwner();
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i];
            if (condition === "self") {
                if (defer.seq !== owner.seq)
                    return false;
            }
            else if (condition === "teammate") {
                if (defer.seq === owner.seq || defer.group !== owner.group)
                    return false;
            }
        }
        var ratio = Number(this._funcParams[1]);
        if (owner.bCtrl.random() >= ratio)
            return false;
        return true;
    };
    return BFHitExtraNormal;
}(BuffFunc_1.BuffFunc));
exports.BFHitExtraNormal = BFHitExtraNormal;
