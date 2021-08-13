"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-07 11:26:46
 * @LastEditors: zyb
 * @LastEditTime: 2021-07-14 11:52:52
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.SkillFunc = void 0;
var Configs_1 = require("../../bundle0/framework/config/Configs");
var BattleConst_1 = require("../../bundle0/battleCommon/BattleConst");
var FindTargets_1 = require("./FindTargets");
var SkillFunc = /** @class */ (function () {
    function SkillFunc(skill, id, inheritUnits) {
        this.id = 0;
        this.skill = null;
        this.funcParams = null;
        this.targets = null;
        this.inheritTargets = null;
        this.isValid = false;
        this.hitJsonName = "";
        this.killed = false;
        this._owner = null;
        this._funcParams = null;
        this.id = id;
        this.skill = skill;
        this._owner = skill.getOwner();
        this.targets = [];
        var conf = Configs_1.Configs.skillFuncConf[id];
        if (!conf) {
            console.warn("invalid skillfunc id: " + id);
        }
        this.hitJsonName = conf.hitJsonName;
        this._funcParams = conf.skillFunc.split("~");
        // 查找目标
        var allTargets = [];
        if (conf.inherit === "inherit_all" || conf.inherit === "inherit_else" || conf.inherit === "inherit_self") {
            allTargets = inheritUnits;
        }
        else {
            allTargets = this._owner.bCtrl.getUnits();
        }
        var isFriendly = conf.isFriendly;
        var targetType = conf.targetType;
        if (this._funcParams[0] === "skill_dmg") {
            // 优先嘲讽
            var isTaunt = false;
            if (this._owner.taunt > 0) {
                var seqs = this._owner.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.TAUNT];
                for (var i = seqs.length - 1; i >= 0; i--) {
                    var seq = seqs[i];
                    var buff = this._owner.bCtrl.buffCtrl.getBuffBySeq(seq);
                    var atker = buff.getAtker();
                    if (!buff.isRemove && atker.isAlive()) {
                        this.targets = [atker];
                        isTaunt = true;
                        break;
                    }
                }
            }
            if (!isTaunt) {
                // 混乱
                if (this._owner.confusion > 0) {
                    isFriendly = (isFriendly + 1) % 2;
                }
                // 锁定后排
                var dType = Number(this._funcParams[1]);
                if (dType !== BattleConst_1.DMG_TYPE.CURE && this._owner.atkTargetBack > 0) {
                    var info = conf.targetType.split("|");
                    var tType = info[0];
                    if (tType === "mixFirst" || tType === "mixAnd" || tType === "mixOr") {
                        tType = info[1];
                    }
                    var ret = tType.split("_");
                    targetType = "back_" + ret[1] + "_1";
                }
                this.targets = FindTargets_1.FindTargets.getTargets(this._owner, isFriendly, targetType, allTargets);
            }
        }
        else {
            this.targets = FindTargets_1.FindTargets.getTargets(this._owner, isFriendly, targetType, allTargets);
        }
        if (conf.inherit === "inherit_all") {
            this.inheritTargets = inheritUnits;
        }
        else if (conf.inherit === "inherit_else") {
            var seqMap = {};
            for (var i = 0; i < this.targets.length; i++) {
                var target = this.targets[i];
                seqMap[target.seq] = true;
            }
            var tmp = [];
            for (var i = 0; i < inheritUnits.length; i++) {
                var target = inheritUnits[i];
                if (!seqMap[target.seq]) {
                    tmp.push(target);
                }
            }
            this.inheritTargets = tmp;
        }
        else {
            this.inheritTargets = this.targets;
        }
        this.isValid = (conf.ratio >= 1 || skill.getOwner().bCtrl.random() < conf.ratio);
    }
    SkillFunc.prototype.execute = function (dmg) {
        if (!this.isValid)
            return 0;
        var funcName = this._funcParams[0];
        var obj = this;
        if (obj[funcName] && typeof obj[funcName] == 'function') {
            if (funcName === "skill_dmg") {
                return obj[funcName](dmg);
            }
            else {
                obj[funcName]();
            }
        }
        else {
            console.warn("SkillFunc invalid funcName", this._funcParams);
        }
        return 0;
    };
    // ============================================================================================
    // 造成伤害
    // 判断暴击。（治疗不计算抗暴）
    SkillFunc.prototype.checkCrit = function (target, dType) {
        // 最终暴击率 = MIN ( MAX ( 攻击方暴击 - 防守方抗暴 , 0 ) , 实际暴击上限 )
        //     attacker.critChance = MIN (  
        //         MAX ( attacker.crit  – defender.calm， 0）
        //    ，[globalBattle].critRatioLimit ）
        if (dType === BattleConst_1.DMG_TYPE.PHYIC) {
            if (target.critNot > 0)
                return false;
            if (this._owner.mustCrit > 0)
                return true;
        }
        var atkerCrit = this._owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.CRIT];
        var critValue = atkerCrit;
        if (dType != BattleConst_1.DMG_TYPE.CURE) {
            critValue -= target.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.CALM];
        }
        var conf = Configs_1.Configs.globalBattleConf[1];
        var critChance = Math.min(Math.max(critValue, 0), conf.critRatioLimit);
        return (this._owner.bCtrl.random() < critChance);
    };
    // 格挡
    SkillFunc.prototype.checkBlock = function (target) {
        // 实际格挡率 = MAX ( 攻击方精准 - 防守方格挡 , 0 )
        // defender.blockChance = MAX ( attacker.hit  – defender.block， 0）
        // 根据实际格挡率defender.blockChance,计算是否格挡
        var blockChance = Math.max(this._owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.HIT] - target.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.BLOCK], 0);
        return (this._owner.bCtrl.random() < blockChance);
    };
    SkillFunc.prototype.isDebug = function () {
        var bType = this._owner.bCtrl.bType;
        if (bType == BattleConst_1.BattleConst.BATTLE_TYPE.DEBUG ||
            bType == BattleConst_1.BattleConst.BATTLE_TYPE.AUTO_FIGHT ||
            bType == BattleConst_1.BattleConst.BATTLE_TYPE.RANDOM_FIGHT)
            return true;
        return false;
    };
    SkillFunc.prototype.skill_dmg = function (inheritDmg) {
        var ret = 0;
        var dType = Number(this._funcParams[1]);
        var coefficient = Number(this._funcParams[2]);
        var eType = Number(this._funcParams[3]);
        var eRatio = Number(this._funcParams[4]);
        var eLimit = Number(this._funcParams[5]);
        var gbConf = Configs_1.Configs.globalBattleConf[1];
        for (var i = 0; i < this.targets.length; i++) {
            var target = this.targets[i];
            var atkAttrs = this._owner.attrs;
            var defAttrs = target.attrs;
            var crit = false;
            var block = false;
            var dmg = 0;
            if (dType == BattleConst_1.DMG_TYPE.CURE && coefficient < 0) {
                //critDamage，判定暴击
                // 如果产生暴击，critDamage = attacker.critDam
                // 如果没有产生暴击，critDamage = 1
                // //治疗量 = 攻击 * 技能倍数 * 暴击倍数 * 治疗加成  * 受疗效果
                // "attacker.healingBonusResult = attacker.atk
                //                                             * coefficient 
                //                                             * critDamage
                //                                             * attacker.cureRatio
                //                                             * defender.healingRatio
                //                                             * -1"
                // 最后治疗值再进行浮动处理，伤害浮动值，从0.9～1.1之间随机取值
                crit = this.checkCrit(target, dType);
                var critDamage = crit ? atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.CRIT_DAM] : 1;
                var healingBonusResult = atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK] *
                    coefficient *
                    critDamage *
                    atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.CURE_RATIO] *
                    defAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.HEALING_RATIO];
                dmg = healingBonusResult;
                if (!this.isDebug()) {
                    dmg = healingBonusResult * (0.9 + 0.2 * this._owner.bCtrl.random());
                }
            }
            else if (dType == BattleConst_1.DMG_TYPE.PHYIC && coefficient > 0) {
                //critDamage，判定暴击
                // 如果产生暴击，critDamage = attacker.critDam
                // 如果没有产生暴击，critDamage = 1
                crit = this.checkCrit(target, dType);
                var critDamage = crit ? this._owner.getFixedCritDam() : 1;
                // //blockDown，判定格挡
                // 如果产生格挡，blockDown =  1 - [globalBattle].blockDamRatio
                // 如果没有产生格挡，blockDown = 1
                block = this.checkBlock(target);
                var blockDown = block ? (1 - gbConf.blockDamRatio) : 1;
                // //jobDam，计算职业增伤
                // 如果防守方为坦克
                // jobDam = MAX（ 1 +  attacker.tanDamUp - defender.tanDamDown  ， 0 ）
                // 如果防守方为战士
                // jobDam = MAX（ 1 +  attacker.zhanDamUp - defender.zhanDamDown  ， 0 ）
                // 如果防守方为法师
                // jobDam = MAX（ 1 +  attacker.faDamUp - defender.faDamDown  ， 0 ）
                // 如果防守方为游侠
                // jobDam = MAX（ 1 +  attacker.xiaDamUp - defender.xiaDamDown  ， 0 ）
                // 如果防守方为辅助
                // jobDam = MAX（ 1 +  attacker.fuDamUp - defender.fuDamDown  ， 0 ）
                var jobDam = 0;
                if (target.job == BattleConst_1.BattleConst.HERO_JOB.TAN) {
                    jobDam = Math.max(1 + atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.TAN_DAM_UP] - defAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.TAN_DAM_DOWN], 0);
                }
                else if (target.job == BattleConst_1.BattleConst.HERO_JOB.ZHAN) {
                    jobDam = Math.max(1 + atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.ZHAN_DAM_UP] - defAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.ZHAN_DAM_DOWN], 0);
                }
                else if (target.job == BattleConst_1.BattleConst.HERO_JOB.FA) {
                    jobDam = Math.max(1 + atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.FA_DAM_UP] - defAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.FA_DAM_DOWN], 0);
                }
                else if (target.job == BattleConst_1.BattleConst.HERO_JOB.XIA) {
                    jobDam = Math.max(1 + atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.XIA_DAM_UP] - defAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.XIA_DAM_DOWN], 0);
                }
                else if (target.job == BattleConst_1.BattleConst.HERO_JOB.FU) {
                    jobDam = Math.max(1 + atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.FU_DAM_UP] - defAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.FU_DAM_DOWN], 0);
                }
                // 保底基础伤害 ： [globalBattle].baseDamRatio
                // //普攻伤害 =  max ( 攻击 - 防御 ，攻击保底)  
                // basicDamage =  MAX (  attacker.atk - defender.def , attacker.atk * [globalBattle].baseDamRatio  ) 
                var basicDamage = Math.max(atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK] - defAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.DEF], atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK] * gbConf.baseDamRatio);
                // //技能伤害 = 普攻伤害  * 技能倍数 * 暴击伤害 *格挡减伤 *  技能增伤 * 职业增伤 * 伤害加深减免
                // "physicsDamage = basicDamage 
                //                         * coefficient  
                //                         * critDamage
                //                         * blockDown
                //                         * max (  1+ attacker.skillDamUp -  defender.skillDamDown ，0 )
                //                         * jobDam
                //                         * max (  1+ attacker.damUp -  defender.damDown ，0 )"
                // 伤害再进行浮动处理，伤害浮动值，从0.9～1.1之间随机取值
                // 最后伤害不能小于1点
                var physicsDamage = basicDamage *
                    coefficient *
                    critDamage *
                    blockDown *
                    Math.max(1 + atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.SKILL_DAM_UP] - defAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.SKILL_DAM_DOWN], 0) *
                    jobDam *
                    Math.max(1 + atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.DAM_UP] - defAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.DAM_DOWN], 0);
                dmg = Math.max(physicsDamage, 1);
                if (!this.isDebug()) {
                    dmg = Math.max(physicsDamage * (0.9 + 0.2 * this._owner.bCtrl.random()), 1);
                }
            }
            else if (dType == BattleConst_1.DMG_TYPE.INHERIT) {
                dmg = inheritDmg * coefficient;
            }
            if (dType !== BattleConst_1.DMG_TYPE.CURE) {
                target.updateMp(gbConf.defendToEnergy);
            }
            var extraDmg = 0;
            if (eType === 1) {
                extraDmg = Math.min(dmg * eRatio, atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK] * eLimit);
            }
            else if (eType === 2) {
                extraDmg = Math.min(defAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE] * eRatio, atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK] * eLimit);
            }
            else if (eType === 3) {
                extraDmg = Math.min(target.hp * eRatio, atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK] * eLimit);
            }
            else if (eType === 4) {
                extraDmg = Math.min(atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK] * eRatio, atkAttrs[BattleConst_1.BattleConst.ATTR_TO_ID.ATK] * eLimit);
            }
            var buffExtraDmg = 0;
            var percent = 1;
            if (dmg > 0) {
                buffExtraDmg = this._owner.bCtrl.buffCtrl.buff_atkExtraDmg(this._owner, target, this.skill.sType, dmg, crit);
                var revise = this._owner.bCtrl.buffCtrl.buff_dmgRevise(this._owner, target, dmg, crit);
                buffExtraDmg += revise.extra;
                percent *= revise.percent;
            }
            dmg = (dmg + extraDmg + buffExtraDmg) * percent;
            ret += dmg;
            var preHp = target.hp;
            var delta = target.updateHp(dmg, this.skill, null, crit);
            if (delta > 0) {
                this._owner.bCtrl.buffCtrl.buff_reboundDmg(target, this._owner, delta);
                this._owner.bCtrl.buffCtrl.buff_hitDmg(target);
            }
            if (preHp > 0 && target.hp <= 0) {
                this.killed = true;
            }
            if (dmg > 0) {
                this._owner.bCtrl.buffCtrl.buff_vampire(this._owner, dmg, crit);
            }
            else if (dmg < 0) {
                this._owner.bCtrl.buffCtrl.buff_cureExtraBuff(target);
            }
            if (dmg > 0) {
                this._owner.bCtrl.buffCtrl.buff_atkEnergy(this._owner, target, dType, crit);
                this._owner.bCtrl.buffCtrl.buff_atkExtraBuff(this._owner, target, this.skill.sType, crit, this.killed);
                if (this.skill.isAct) {
                    this._owner.bCtrl.buffCtrl.buff_hitExtraNormal(target);
                }
            }
            this._owner.bCtrl.callSceneFunc("unitHit", target.seq, this.hitJsonName, i == 0);
        }
        return ret;
    };
    // ============================================================================================
    // 驱散buff
    SkillFunc.prototype.skill_dispel = function () {
        var ratio = Number(this._funcParams[1]);
        for (var i = 0; i < this.targets.length; i++) {
            var target = this.targets[i];
            if (this._owner.bCtrl.random() < ratio) {
                var cnt = Number(this._funcParams[2]);
                var buffs = this._owner.bCtrl.buffCtrl.getBuffs(target);
                var tmp = [];
                var buffType = [];
                for (var j = 3; j < this._funcParams.length; j++) {
                    buffType.push(Number(this._funcParams[j]));
                }
                for (var j = 0; j < buffs.length; j++) {
                    var buff = buffs[j];
                    if (buff.isBuffType(buffType) && buff.dispel) {
                        tmp.push(buff);
                    }
                }
                while (cnt != -1 && tmp.length > cnt) {
                    var idx = Math.floor(this._owner.bCtrl.random() * tmp.length);
                    tmp.splice(idx, 1);
                }
                for (var i_1 = tmp.length - 1; i_1 >= 0; i_1--) {
                    this._owner.bCtrl.buffCtrl.removeBuff(tmp[i_1]);
                }
            }
        }
    };
    // ============================================================================================
    // 附加buff
    SkillFunc.prototype.skill_extraBuff = function () {
        var ratio = Number(this._funcParams[1]);
        var cnt = Number(this._funcParams[2]);
        var buffIds = [];
        for (var j = 3; j < this._funcParams.length; j++) {
            buffIds.push(Number(this._funcParams[j]));
        }
        for (var i = 0; i < this.targets.length; i++) {
            var target = this.targets[i];
            var idxs = [];
            for (var j = 0; j < buffIds.length; j++) {
                idxs.push(j);
            }
            var num = 0;
            while (idxs.length > 0 && num < cnt) {
                var idx = Math.floor(this._owner.bCtrl.random() * idxs.length);
                var buffId = buffIds[idxs[idx]];
                this._owner.bCtrl.buffCtrl.addBuff(target, buffId, this._owner, ratio);
                idxs.splice(idx, 1);
                num++;
            }
        }
    };
    // ============================================================================================
    // 英雄数量加buff
    SkillFunc.prototype.skill_heroCntBuff = function () {
        var cnt = 0;
        var condition = this._funcParams[1].split("_");
        if (condition[0] === "elem") {
            var group = this._owner.group;
            if (condition[1] === "enemy") {
                group += 1;
                if (group > 2)
                    group = 1;
            }
            var elemMap = {};
            for (var i = 2; i < condition.length; i++) {
                elemMap[Number(condition[i])] = true;
            }
            var units = this._owner.bCtrl.getUnits();
            for (var i = 0; i < units.length; i++) {
                var unit = units[i];
                if (unit.isAlive() && unit.group === group && elemMap[unit.elem])
                    cnt++;
            }
        }
        else if (condition[0] === "buffType") {
            var group = this._owner.group;
            if (condition[1] === "enemy") {
                group += 1;
                if (group > 2)
                    group = 1;
            }
            var bType = Number(condition[2]);
            var units = this._owner.bCtrl.getUnits();
            for (var i = 0; i < units.length; i++) {
                var unit = units[i];
                if (unit.isAlive() && unit.group === group && unit.hasBuffType([bType]))
                    cnt++;
            }
        }
        for (var i = 0; i < cnt; i++) {
            for (var j = 0; j < this.targets.length; j++) {
                var target = this.targets[j];
                for (var k = 2; k < this._funcParams.length; k++) {
                    var buffId = Number(this._funcParams[k]);
                    this._owner.bCtrl.buffCtrl.addBuff(target, buffId, this._owner);
                }
            }
        }
    };
    // ============================================================================================
    // 替换技能
    SkillFunc.prototype.skill_replace = function () {
        var sType = this._funcParams[1];
        var skillId = Number(this._funcParams[2]);
        this._owner.replaceSkill(sType, skillId);
    };
    // ============================================================================================
    // 能量增减
    SkillFunc.prototype.skill_energy = function () {
        var ratio = Number(this._funcParams[1]);
        var delta = Number(this._funcParams[2]);
        for (var i = 0; i < this.targets.length; i++) {
            var target = this.targets[i];
            if (this._owner.bCtrl.random() < ratio) {
                target.updateMp(delta);
            }
        }
    };
    // ============================================================================================
    // 额外触发技能
    SkillFunc.prototype.skill_next = function () {
        var ratio = Number(this._funcParams[1]);
        var skillType = this._funcParams[2];
        var nextCnt = Number(this._funcParams[3]);
        if (this._owner.bCtrl.random() < ratio) {
            this._owner.bCtrl.pushExtraSkill(BattleConst_1.BattleConst.SKILL_FUNC.SKILL_NEXT, this._owner, skillType, this.skill.isAct ? nextCnt : 0);
        }
    };
    // ============================================================================================
    // dot伤害立即生效
    SkillFunc.prototype.skill_dotExecute = function () {
        for (var i = 0; i < this.targets.length; i++) {
            var target = this.targets[i];
            var seqs = target.buffFuncs[BattleConst_1.BattleConst.BUFF_FUNC.DOT] || [];
            for (var j = 0; j < seqs.length; j++) {
                var seq = seqs[j];
                var buff = this._owner.bCtrl.buffCtrl.getBuffBySeq(seq);
                var func = buff.getFunc();
                func.buff_dot();
            }
        }
    };
    return SkillFunc;
}());
exports.SkillFunc = SkillFunc;
