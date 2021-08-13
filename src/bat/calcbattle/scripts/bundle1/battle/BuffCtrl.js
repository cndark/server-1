"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-25 17:52:19
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-07-23 17:45:25
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.BuffCtrl = void 0;
var Configs_1 = require("../../bundle0/framework/config/Configs");
var BattleConst_1 = require("../../bundle0/battleCommon/BattleConst");
var UnitBuff_1 = require("./UnitBuff");
var BTYPE_TO_ATTR = {
    102: [BattleConst_1.BattleConst.ATTR_TO_ID.DOT_HIT, BattleConst_1.BattleConst.ATTR_TO_ID.FIRE_DOT_RES, BattleConst_1.BattleConst.ATTR_TO_ID.DOT_RES],
    103: [BattleConst_1.BattleConst.ATTR_TO_ID.DOT_HIT, BattleConst_1.BattleConst.ATTR_TO_ID.BLOOD_DOT_RES, BattleConst_1.BattleConst.ATTR_TO_ID.DOT_RES],
    104: [BattleConst_1.BattleConst.ATTR_TO_ID.DOT_HIT, BattleConst_1.BattleConst.ATTR_TO_ID.DU_DOT_RES, BattleConst_1.BattleConst.ATTR_TO_ID.DOT_RES],
    105: [BattleConst_1.BattleConst.ATTR_TO_ID.DOT_HIT, BattleConst_1.BattleConst.ATTR_TO_ID.ICE_DOT_RES, BattleConst_1.BattleConst.ATTR_TO_ID.DOT_RES],
    201: [BattleConst_1.BattleConst.ATTR_TO_ID.KONG_HIT, BattleConst_1.BattleConst.ATTR_TO_ID.YUN_RES, BattleConst_1.BattleConst.ATTR_TO_ID.KONG_RES],
    202: [BattleConst_1.BattleConst.ATTR_TO_ID.KONG_HIT, BattleConst_1.BattleConst.ATTR_TO_ID.MO_RES, BattleConst_1.BattleConst.ATTR_TO_ID.KONG_RES],
    206: [BattleConst_1.BattleConst.ATTR_TO_ID.KONG_HIT, BattleConst_1.BattleConst.ATTR_TO_ID.NU_RES, BattleConst_1.BattleConst.ATTR_TO_ID.KONG_RES],
    207: [BattleConst_1.BattleConst.ATTR_TO_ID.KONG_HIT, BattleConst_1.BattleConst.ATTR_TO_ID.BING_RES, BattleConst_1.BattleConst.ATTR_TO_ID.KONG_RES],
    208: [BattleConst_1.BattleConst.ATTR_TO_ID.KONG_HIT, BattleConst_1.BattleConst.ATTR_TO_ID.SHI_RES, BattleConst_1.BattleConst.ATTR_TO_ID.KONG_RES],
    209: [BattleConst_1.BattleConst.ATTR_TO_ID.KONG_HIT, BattleConst_1.BattleConst.ATTR_TO_ID.MA_RES, BattleConst_1.BattleConst.ATTR_TO_ID.KONG_RES], // 麻痹
};
var BuffCtrl = /** @class */ (function () {
    function BuffCtrl(bCtrl) {
        this.bCtrl = null;
        this.reboundLimit = {}; // 反弹buff死循环处理{ "buffId_atkSeq_ownSeq": true }
        this._buffs = {};
        this.bCtrl = bCtrl;
        this._buffs = {};
        this.reboundLimit = {};
    }
    BuffCtrl.prototype._getBuffExtraRatio = function (id, owner, atker, ratio) {
        var ret = ratio;
        var conf = Configs_1.Configs.buffConf[id];
        for (var i = 0; i < conf.buffType.length; i++) {
            var bType = conf.buffType[i];
            var attrs = BTYPE_TO_ATTR[bType];
            if (attrs) {
                ret = Math.max(ratio + atker.attrs[attrs[0]] - owner.attrs[attrs[1]] - owner.attrs[attrs[2]], 0);
                break;
            }
        }
        return ret;
    };
    BuffCtrl.prototype.addBuff = function (owner, buffId, atker, ratio, showTips) {
        if (ratio === void 0) { ratio = 1; }
        if (showTips === void 0) { showTips = true; }
        if (!owner.isAlive())
            return;
        // 免疫buffType
        var conf = Configs_1.Configs.buffConf[buffId];
        for (var i = 0; i < conf.buffType.length; i++) {
            var bType = conf.buffType[i];
            if (owner.immuneBuffs[bType])
                return;
        }
        // buff抗性
        var newRatio = this._getBuffExtraRatio(buffId, owner, atker, ratio);
        if (this.bCtrl.random() >= newRatio)
            return;
        // 反弹buff
        if (this.buff_reboundBuff(owner, buffId, atker))
            return;
        if (!owner.buffGrpCnt[conf.group]) {
            owner.buffGrpCnt[conf.group] = 0;
        }
        var cnt = owner.buffGrpCnt[conf.group];
        var isOverLap = false;
        // buff组叠加上限，超出是按照优先级剔除
        if (cnt == conf.overlap) {
            isOverLap = true;
            var buffSeqs = owner.buffSeqs;
            var preBuff = null;
            for (var i = buffSeqs.length - 1; i >= 0; i--) {
                var seq = buffSeqs[i];
                var b = this._buffs[seq];
                if (!b.isRemove && b.group === conf.group && b.priority < conf.priority) {
                    if (!preBuff ||
                        (b.priority < preBuff.priority) ||
                        (b.priority === preBuff.priority && b.round <= preBuff.round)) {
                        preBuff = b;
                    }
                }
            }
            if (!preBuff)
                return;
            this.removeBuff(preBuff);
        }
        var buff = new UnitBuff_1.UnitBuff(owner, buffId, atker);
        this._buffs[buff.seq] = buff;
        // 更新BattleUnit相关buff数据
        owner.buffGrpCnt[buff.group]++;
        for (var i = conf.buffType.length - 1; i >= 0; i--) {
            var bType = conf.buffType[i];
            var cnt_1 = owner.buffTypeCnt[bType] || 0;
            owner.buffTypeCnt[bType] = cnt_1 + 1;
        }
        var funcName = conf.func.split("~")[0];
        var funcs = owner.buffFuncs[funcName] || [];
        funcs.push(buff.seq);
        owner.buffFuncs[funcName] = funcs;
        owner.buffSeqs.push(buff.seq);
        if (showTips) {
            this.bCtrl.callSceneFunc("addBuff", buff);
        }
        // buff生效
        buff.add();
        this.buff_buffTypeExtraBuff(buff);
        if (!isOverLap) {
            this.buff_buffGrpBuff(owner, buff.group);
        }
    };
    BuffCtrl.prototype.removeBuff = function (buff) {
        var conf = Configs_1.Configs.buffConf[buff.id];
        var owner = buff.getOwner();
        owner.buffGrpCnt[buff.group]--;
        for (var i = conf.buffType.length - 1; i >= 0; i--) {
            var bType = conf.buffType[i];
            owner.buffTypeCnt[bType]--;
        }
        this.bCtrl.callSceneFunc("removeBuff", buff);
        buff.remove();
    };
    BuffCtrl.prototype.removeBuffByBTypes = function (owner, bTypes) {
        var buffSeqs = owner.buffSeqs;
        for (var i = 0; i < buffSeqs.length; i++) {
            var seq = buffSeqs[i];
            var buff = this._buffs[seq];
            if (!buff.isRemove && buff.isBuffType(bTypes)) {
                this.removeBuff(buff);
            }
        }
    };
    BuffCtrl.prototype.removeBuffByGrp = function (owner, group) {
        var buffSeqs = owner.buffSeqs;
        for (var i = 0; i < buffSeqs.length; i++) {
            var seq = buffSeqs[i];
            var buff = this._buffs[seq];
            if (!buff.isRemove && buff.group === group) {
                this.removeBuff(buff);
            }
        }
    };
    BuffCtrl.prototype.clearRemoveBuff = function () {
        var keys = Object.keys(this._buffs);
        for (var i = keys.length - 1; i >= 0; i--) {
            var seq = Number(keys[i]);
            var buff = this._buffs[seq];
            if (buff.isRemove) {
                delete this._buffs[seq];
            }
        }
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            var buffSeqs = unit.buffSeqs;
            for (var j = buffSeqs.length - 1; j >= 0; j--) {
                var seq = buffSeqs[j];
                var buff = this._buffs[seq];
                if (!buff) {
                    buffSeqs.splice(j, 1);
                }
            }
            var buffFuncs = unit.buffFuncs;
            for (var bfunc in buffFuncs) {
                var seqs = buffFuncs[bfunc];
                for (var j = seqs.length; j >= 0; j--) {
                    var seq = seqs[j];
                    var buff = this._buffs[seq];
                    if (!buff) {
                        seqs.splice(j, 1);
                    }
                }
            }
        }
    };
    BuffCtrl.prototype.updateRound = function (owner) {
        var buffSeqs = owner.buffSeqs;
        for (var i = buffSeqs.length - 1; i >= 0; i--) {
            var seq = buffSeqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            if (buff.updateRound()) {
                this.removeBuff(buff);
            }
        }
    };
    BuffCtrl.prototype.deathDropBuff = function (owner) {
        var buffSeqs = owner.buffSeqs;
        for (var i = buffSeqs.length - 1; i >= 0; i--) {
            var seq = buffSeqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            if (buff.deathDrop == 1) {
                this.removeBuff(buff);
            }
        }
    };
    BuffCtrl.prototype.getBuffBySeq = function (seq) {
        return this._buffs[seq];
    };
    BuffCtrl.prototype.getBuffs = function (owner) {
        var buffs = [];
        var buffSeqs = owner.buffSeqs;
        for (var i = buffSeqs.length - 1; i >= 0; i--) {
            var seq = buffSeqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            buffs.push(buff);
        }
        return buffs;
    };
    // ============================================================================================
    // buff func相关
    BuffCtrl.prototype.buff_reboundBuff = function (owner, buffId, atker) {
        var key = buffId + "_" + atker.seq + "_" + owner.seq;
        if (this.reboundLimit[key])
            return false;
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.REBOUND_BUFF] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            if (func.buff_reboundBuff(buffId, atker)) {
                this.reboundLimit[key] = true;
                return true;
            }
        }
        return false;
    };
    BuffCtrl.prototype.buff_buffTypeExtraBuff = function (buff) {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.BUFF_TYPE_EXTRA_BUFF] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_buffTypeExtraBuff(buff);
            }
        }
    };
    BuffCtrl.prototype.buff_attrToShield = function (owner, dmg) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ATTR_TO_SHIELD] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            dmg = func.buff_attrToShield(dmg);
            if (dmg <= 0)
                break;
        }
        return dmg;
    };
    BuffCtrl.prototype.buff_atkExtraDmg = function (owner, target, sType, baseDmg, crit) {
        var ret = 0;
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ATK_EXTRA_DMG] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            ret += func.buff_atkExtraDmg(target, sType, baseDmg, crit);
        }
        return ret;
    };
    BuffCtrl.prototype.buff_dmgRevise = function (atker, defer, baseDmg, crit) {
        var ret = { extra: 0, percent: 1 };
        var units = [atker, defer];
        for (var i = units.length - 1; i >= 0; i--) {
            var owner = units[i];
            var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.DMG_REVISE] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var buff = this._buffs[seq];
                if (buff.isRemove)
                    continue;
                var func = buff.getFunc();
                var revise = func.buff_dmgRevise(atker, defer, baseDmg, crit);
                ret.extra += revise.extra;
                ret.percent *= revise.percent;
            }
        }
        return ret;
    };
    BuffCtrl.prototype.buff_vampire = function (owner, dmg, crit) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.VAMPIRE] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_vampire(dmg, crit);
        }
    };
    BuffCtrl.prototype.buff_atkEnergy = function (atker, defer, dType, crit) {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ATK_ENERGY] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_atkEnergy(atker, defer, dType, crit);
            }
        }
    };
    BuffCtrl.prototype.buff_atkExtraBuff = function (atker, defer, sType, crit, kill) {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ATK_EXTRA_BUFF] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_atkExtraBuff(atker, defer, sType, crit, kill);
            }
        }
    };
    BuffCtrl.prototype.buff_deathExtraBuff = function (death) {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (unit.seq != death.seq && !unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.DEATH_EXTRA_BUFF] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_deathExtraBuff(death);
            }
        }
    };
    BuffCtrl.prototype.buff_hpReduceBuff = function (owner, reduceHp) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.HP_REDUCE_BUFF] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_hpReduceBuff(reduceHp);
        }
    };
    BuffCtrl.prototype.buff_hpPointBuff = function (owner, reduceHp) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.HP_POINT_BUFF] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_hpPointBuff(reduceHp);
        }
    };
    BuffCtrl.prototype.buff_deathEnergy = function (death) {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (unit.seq != death.seq && !unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.DEATH_ENERGY] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_deathEnergy(death);
            }
        }
    };
    BuffCtrl.prototype.buff_cureEnergy = function (owner) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.CURE_ENERGY] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_cureEnergy();
        }
    };
    BuffCtrl.prototype.buff_actEnergy = function (owner) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ACT_ENERGY] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_actEnergy();
        }
    };
    BuffCtrl.prototype.buff_actStealEnergy = function (owner) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ACT_STEAL_ENERGY] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_actStealEnergy();
        }
    };
    BuffCtrl.prototype.buff_dot = function (owner) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.DOT] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_dot();
        }
    };
    BuffCtrl.prototype.buff_dotGamble = function (owner) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.DOT_GAMBLE] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_dotGamble();
        }
    };
    BuffCtrl.prototype.buff_dotExtraBuff = function (owner) {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.DOT_EXTRA_BUFF] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_dotExtraBuff(owner);
            }
        }
    };
    BuffCtrl.prototype.buff_cureExtraBuff = function (owner) {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.CURE_EXTRA_BUFF] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_cureExtraBuff(owner);
            }
        }
    };
    BuffCtrl.prototype.buff_reboundDmg = function (owner, atker, delta) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.REBOUND_DMG] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_reboundDmg(atker, delta);
        }
    };
    BuffCtrl.prototype.buff_hitDmg = function (owner) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.HIT_DMG] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_hitDmg();
        }
    };
    BuffCtrl.prototype.buff_deathDmg = function (death) {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (unit.seq != death.seq && !unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.DEATH_DMG] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_deathDmg(death);
            }
        }
    };
    BuffCtrl.prototype.buff_actStartDmgBuff = function (owner) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ACT_START_DMG_BUFF] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_actStartDmgBuff();
        }
    };
    BuffCtrl.prototype.buff_atkDeath = function (owner, sType, hpRatio) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ATK_DEATH] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            if (func.buff_atkDeath(sType, hpRatio))
                return true;
        }
        return false;
    };
    BuffCtrl.prototype.buff_roundChangeHp = function () {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ROUND_CHANGE_HP] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_roundChangeHp();
            }
        }
    };
    BuffCtrl.prototype.buff_roundStartBuff = function () {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ROUND_START_BUFF] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_roundStartBuff();
            }
        }
    };
    BuffCtrl.prototype.buff_roundStartDmgBuff = function () {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ROUND_START_DMG_BUFF] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_roundStartDmgBuff();
            }
        }
    };
    BuffCtrl.prototype.buff_roundEndDmgBuff = function () {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ROUND_END_DMG_BUFF] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_roundEndDmgBuff();
            }
        }
    };
    BuffCtrl.prototype.buff_roundEndBuff = function () {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ROUND_END_BUFF] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_roundEndBuff();
            }
        }
    };
    BuffCtrl.prototype.buff_actExtraBuff = function (owner) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ACT_EXTRA_BUFF] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_actExtraBuff();
        }
    };
    BuffCtrl.prototype.buff_critDamFixed = function (owner) {
        var ret = -1;
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.CRIT_DAM_FIXED] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            var value = func.buff_critDamFixed();
            if (ret === -1 || ret > value)
                ret = value;
        }
        return ret;
    };
    BuffCtrl.prototype.buff_crit = function (owner) {
        var ret = 0;
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.CRIT] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            ret += func.buff_crit();
        }
        return ret;
    };
    BuffCtrl.prototype.buff_attrCond = function () {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ATTR_COND] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                func.buff_attrCond();
            }
        }
    };
    BuffCtrl.prototype.buff_buffGrpBuff = function (owner, group) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.BUFF_GRP_BUFF] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            func.buff_buffGrpBuff(group);
        }
    };
    BuffCtrl.prototype.buff_actNewSkill = function (owner, sType) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ACT_NEW_SKILL] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            var sId = func.buff_actNewSkill(sType);
            if (sId !== 0)
                return sId;
        }
        return 0;
    };
    BuffCtrl.prototype.buff_actExtraSpecial = function (owner, skill) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ACT_EXTRA_SPECIAL] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            if (func.buff_actExtraSpecial(skill))
                return true;
        }
        return false;
    };
    BuffCtrl.prototype.buff_hitExtraNormal = function (defer) {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.HIT_EXTRA_NORMAL] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                if (func.buff_hitExtraNormal(defer)) {
                    this.bCtrl.pushExtraSkill(BattleConst_1.BattleConst.BUFF_FUNC.HIT_EXTRA_NORMAL, b.getOwner(), BattleConst_1.BattleConst.SKILL_TYPE.NORMAL);
                    return;
                }
            }
        }
    };
    BuffCtrl.prototype.buff_deathExtraSpecial = function (death) {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            if (!unit.isAlive())
                continue;
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.DEATH_EXTRA_SPECIAL] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                var ret = func.buff_deathExtraSpecial(death);
                if (ret) {
                    this.bCtrl.pushExtraSkill(BattleConst_1.BattleConst.BUFF_FUNC.DEATH_EXTRA_SPECIAL, unit, BattleConst_1.BattleConst.SKILL_TYPE.SPECIAL);
                }
            }
        }
    };
    BuffCtrl.prototype.buff_actExtraAct = function (owner) {
        var seqs = owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.ACT_EXTRA_ACT] || [];
        for (var i = seqs.length - 1; i >= 0; i--) {
            var seq = seqs[i];
            var buff = this._buffs[seq];
            if (buff.isRemove)
                continue;
            var func = buff.getFunc();
            if (func.buff_actExtraAct())
                return true;
        }
        return false;
    };
    BuffCtrl.prototype.buff_reborn = function (death) {
        var units = this.bCtrl.getUnits();
        for (var i = units.length - 1; i >= 0; i--) {
            var unit = units[i];
            var seqs = unit.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.REBORN] || [];
            for (var j = seqs.length - 1; j >= 0; j--) {
                var seq = seqs[j];
                var b = this._buffs[seq];
                if (b.isRemove)
                    continue;
                var func = b.getFunc();
                var ret = func.buff_reborn(death);
                if (ret > 0) {
                    return ret;
                }
            }
        }
        return 0;
    };
    return BuffCtrl;
}());
exports.BuffCtrl = BuffCtrl;
