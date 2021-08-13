"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-10 14:54:22
 * @LastEditors: zyb
 * @LastEditTime: 2021-07-14 11:53:13
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.UnitBuff = void 0;
var Configs_1 = require("../../bundle0/framework/config/Configs");
var BFFactory_1 = require("./buffFunc/BFFactory");
var UnitBuff = /** @class */ (function () {
    function UnitBuff(owner, id, atker) {
        this._owner = null;
        this._atker = null;
        this._func = null;
        this._noAttack = 0;
        this._noSkill = 0;
        this.round = 0;
        this.maxRound = 0;
        this.isRemove = false;
        this.seq = 0; // 唯一标识
        this.id = 0;
        this.group = 0;
        this.priority = 0;
        this.deathDrop = 0;
        this.overlap = 0;
        this.buffType = {};
        this.persistJson = "";
        this.buffPos = "";
        this.dispel = true;
        this._owner = owner;
        this._atker = atker;
        this.id = id;
        this.seq = UnitBuff.unique;
        var conf = Configs_1.Configs.buffConf[id];
        this.group = conf.group;
        this.round = conf.round;
        this.maxRound = conf.round;
        this.priority = conf.priority;
        this.deathDrop = conf.deathDrop;
        this.overlap = conf.overlap;
        this.buffType = {};
        for (var i = 0; i < conf.buffType.length; i++) {
            this.buffType[conf.buffType[i]] = true;
        }
        this.persistJson = conf.persistJson;
        this.buffPos = conf.position;
        this.dispel = conf.dispel == 1; //1可驱散
        this._noAttack = 0;
        this._noSkill = 0;
        for (var i = 0; i < conf.controlType.length; i++) {
            var value = conf.controlType[i];
            if (value == 1) {
                this._noAttack++;
            }
            else if (value == 2) {
                this._noSkill++;
            }
        }
        this._func = BFFactory_1.BFFactory.createBuffFunc(this, conf.func);
    }
    Object.defineProperty(UnitBuff, "unique", {
        get: function () {
            var limit = Math.pow(10, 9);
            UnitBuff._unique++;
            if (UnitBuff._unique > limit) {
                UnitBuff._unique -= limit;
            }
            return UnitBuff._unique;
        },
        enumerable: false,
        configurable: true
    });
    UnitBuff.prototype.updateRound = function () {
        this.round--;
        return (this.round <= 0);
    };
    UnitBuff.prototype.isBuffType = function (bTypes) {
        for (var i = 0; i < bTypes.length; i++) {
            var bType = bTypes[i];
            if (this.buffType[bType] == true)
                return true;
        }
        return false;
    };
    UnitBuff.prototype.add = function () {
        this.isRemove = false;
        this.addAttr();
        this.addControl();
        this._func.add();
    };
    UnitBuff.prototype.addAttr = function () {
        var conf = Configs_1.Configs.buffConf[this.id];
        for (var i = 0; i < conf.attribute.length; i++) {
            var attrData = conf.attribute[i];
            this._owner.updateBuffAttr(attrData.id, attrData.val);
        }
    };
    UnitBuff.prototype.addControl = function () {
        this._owner.updateControlType(this._noAttack, this._noSkill);
    };
    UnitBuff.prototype.remove = function () {
        this.isRemove = true;
        this.removeAttr();
        this.removeControl();
        this._func.remove();
    };
    UnitBuff.prototype.removeAttr = function () {
        var conf = Configs_1.Configs.buffConf[this.id];
        for (var i = 0; i < conf.attribute.length; i++) {
            var attrData = conf.attribute[i];
            this._owner.updateBuffAttr(attrData.id, -1 * attrData.val);
        }
    };
    UnitBuff.prototype.removeControl = function () {
        this._owner.updateControlType(-1 * this._noAttack, -1 * this._noSkill);
    };
    UnitBuff.prototype.getOwner = function () {
        return this._owner;
    };
    UnitBuff.prototype.getAtker = function () {
        return this._atker;
    };
    UnitBuff.prototype.getFunc = function () {
        return this._func;
    };
    UnitBuff._unique = 1;
    return UnitBuff;
}());
exports.UnitBuff = UnitBuff;
