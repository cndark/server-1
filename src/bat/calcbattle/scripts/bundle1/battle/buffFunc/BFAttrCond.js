"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-26 11:55:00
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-05-19 09:51:59
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
exports.BFAttrCond = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var BuffFunc_1 = require("./BuffFunc");
var BFAttrCond = /** @class */ (function (_super) {
    __extends(BFAttrCond, _super);
    function BFAttrCond(buff, func, params) {
        var _this = _super.call(this, buff, func, params) || this;
        _this._bfAttrs = [];
        _this._bfValid = false;
        _this._bfAttrs = [];
        var info = _this._funcParams[1].split("|");
        for (var i = info.length - 1; i >= 0; i--) {
            var data = info[i].split("_");
            _this._bfAttrs.push({ id: Number(data[0]), val: Number(data[1]) });
        }
        return _this;
    }
    Object.defineProperty(BFAttrCond.prototype, "bfValid", {
        get: function () {
            return this._bfValid;
        },
        set: function (valid) {
            if (this._bfValid != valid) {
                this._bfValid = valid;
                if (valid)
                    this.addAttr();
                else
                    this.removeAttr();
            }
        },
        enumerable: false,
        configurable: true
    });
    BFAttrCond.prototype.add = function () {
        this.buff_attrCond();
    };
    BFAttrCond.prototype.remove = function () {
        this.bfValid = false;
    };
    BFAttrCond.prototype.addAttr = function () {
        var owner = this._buff.getOwner();
        for (var i = this._bfAttrs.length - 1; i >= 0; i--) {
            var attrData = this._bfAttrs[i];
            owner.updateBuffAttr(attrData.id, attrData.val);
        }
    };
    BFAttrCond.prototype.removeAttr = function () {
        var owner = this._buff.getOwner();
        for (var i = this._bfAttrs.length - 1; i >= 0; i--) {
            var attrData = this._bfAttrs[i];
            owner.updateBuffAttr(attrData.id, -1 * attrData.val);
        }
    };
    BFAttrCond.prototype.buff_attrCond = function () {
        var owner = this._buff.getOwner();
        var atker = this._buff.getAtker();
        var valid = true;
        var conditions = this._funcParams[0].split("|");
        for (var i = 0; i < conditions.length; i++) {
            var condition = conditions[i].split("_");
            if (condition[0] === "hpPoint") {
                var hpRatio = owner.hp / owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
                var limit = Number(condition[2]);
                if ((condition[1] === "lt" && hpRatio >= limit) ||
                    (condition[1] === "gt" && hpRatio <= limit)) {
                    valid = false;
                    break;
                }
            }
            else if (condition[0] === "originLive") {
                if (!atker.isAlive()) {
                    valid = false;
                    break;
                }
            }
        }
        this.bfValid = valid;
    };
    return BFAttrCond;
}(BuffFunc_1.BuffFunc));
exports.BFAttrCond = BFAttrCond;
