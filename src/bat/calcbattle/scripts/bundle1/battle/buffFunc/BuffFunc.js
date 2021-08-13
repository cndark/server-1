"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-14 20:36:37
 * @LastEditors: chenjie
 * @LastEditTime: 2021-06-26 18:38:21
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.BuffFunc = void 0;
var BuffFunc = /** @class */ (function () {
    function BuffFunc(buff, func, params) {
        this.funcName = "";
        this._buff = null;
        this._funcParams = [];
        this.funcName = func;
        this._funcParams = params;
        this._buff = buff;
    }
    BuffFunc.prototype.add = function () { };
    BuffFunc.prototype.remove = function () { };
    BuffFunc.prototype.buff_reboundBuff = function (buffId, atker) { return false; };
    ;
    BuffFunc.prototype.buff_buffTypeExtraBuff = function (buff) { };
    ;
    BuffFunc.prototype.buff_attrToShield = function (dmg) { return dmg; };
    ;
    BuffFunc.prototype.buff_deathEnergy = function (death) { };
    ;
    BuffFunc.prototype.buff_cureEnergy = function () { };
    ;
    BuffFunc.prototype.buff_actEnergy = function () { };
    ;
    BuffFunc.prototype.buff_actStealEnergy = function () { };
    ;
    BuffFunc.prototype.buff_vampire = function (dmg, crit) { };
    BuffFunc.prototype.buff_atkExtraDmg = function (target, sType, baseDmg, crit) { return 0; };
    BuffFunc.prototype.buff_dmgRevise = function (atker, defer, baseDmg, crit) { return { extra: 0, percent: 1 }; };
    ;
    BuffFunc.prototype.buff_dot = function () { };
    ;
    BuffFunc.prototype.buff_dotGamble = function () { };
    ;
    BuffFunc.prototype.buff_reboundDmg = function (atker, delta) { };
    ;
    BuffFunc.prototype.buff_hpReduceBuff = function (reduceHp) { };
    ;
    BuffFunc.prototype.buff_hpPointBuff = function (reduceHp) { };
    ;
    BuffFunc.prototype.buff_hitDmg = function () { };
    ;
    BuffFunc.prototype.buff_deathDmg = function (death) { };
    ;
    BuffFunc.prototype.buff_actStartDmgBuff = function () { };
    ;
    BuffFunc.prototype.buff_atkDeath = function (sType, hpRatio) { return false; };
    ;
    BuffFunc.prototype.buff_roundStartBuff = function () { };
    ;
    BuffFunc.prototype.buff_roundStartDmgBuff = function () { };
    ;
    BuffFunc.prototype.buff_roundEndDmgBuff = function () { };
    ;
    BuffFunc.prototype.buff_roundEndBuff = function () { };
    ;
    BuffFunc.prototype.buff_roundChangeHp = function () { };
    ;
    BuffFunc.prototype.buff_actExtraBuff = function () { };
    ;
    BuffFunc.prototype.buff_critDamFixed = function () { return 0; };
    ;
    BuffFunc.prototype.buff_crit = function () { return 0; };
    ;
    BuffFunc.prototype.buff_atkEnergy = function (atker, defer, dType, crit) { };
    BuffFunc.prototype.buff_atkExtraBuff = function (atker, defer, sType, crit, kill) { };
    ;
    BuffFunc.prototype.buff_deathExtraBuff = function (death) { };
    ;
    BuffFunc.prototype.buff_dotExtraBuff = function (dotUnit) { };
    ;
    BuffFunc.prototype.buff_cureExtraBuff = function (cureUnit) { };
    ;
    BuffFunc.prototype.buff_attrCond = function () { };
    ;
    BuffFunc.prototype.buff_buffGrpBuff = function (group) { };
    ;
    BuffFunc.prototype.buff_actNewSkill = function (sType) { return 0; };
    ;
    BuffFunc.prototype.buff_actExtraSpecial = function (skill) { return false; };
    ;
    BuffFunc.prototype.buff_hitExtraNormal = function (defer) { return false; };
    ;
    BuffFunc.prototype.buff_deathExtraSpecial = function (death) { return false; };
    ;
    BuffFunc.prototype.buff_actExtraAct = function () { return false; };
    ;
    BuffFunc.prototype.buff_reborn = function (death) { return 0; };
    ;
    return BuffFunc;
}());
exports.BuffFunc = BuffFunc;
