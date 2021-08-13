"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-03-28 12:06:21
 * @LastEditors: zyb
 * @LastEditTime: 2021-07-13 11:00:49
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.BattleUnit = void 0;
var Configs_1 = require("../../bundle0/framework/config/Configs");
var BattleConst_1 = require("../../bundle0/battleCommon/BattleConst");
var UnitSkill_1 = require("./UnitSkill");
var BattleUnit = /** @class */ (function () {
    function BattleUnit(ctrl, fighter, group) {
        var _a, _b, _c, _d, _e, _f;
        this.bCtrl = null;
        this.attrs = {};
        // buff功能相关
        this.noAttack = 0; // 不能普攻
        this.noSkill = 0; // 不能释放技能
        this.invincible = 0; // 无敌
        this.critNot = 0; // 必定不暴击
        this.mustCrit = 0; // 必定暴击
        this.firstAct = 0; // 优先出手
        this.confusion = 0; // 混乱
        this.atkTargetBack = 0; // 锁定后排
        this.taunt = 0; // 嘲讽
        this.immuneDeath = 0; // 免疫致命伤害
        this.immuneDeathCnt = 0; // 免疫致命伤害次数
        this.rebornNot = 0; // 必定不能复活
        this.rebornCnt = 0; // 复活次数
        this.actCnt = 0; // 一回合行动次数
        this.shieldHp = 0; // 护盾总血量
        this.buffGrpCnt = {};
        this.buffTypeCnt = {};
        this.immuneBuffs = {};
        this.buffFuncs = {};
        this.buffSeqs = [];
        this._normalSkill = null;
        this._specialSkill = null;
        this._passiveSkills = [];
        this._originAttr = {};
        this._state = BattleConst_1.UNIT_STATE.IDLE;
        this._curSkill = null;
        this._rebornDt = 0;
        this.hp = 0; // 血量
        this.mp = 0; // 能量
        this.maxMp = 0; // 能量上限
        this.seq = 0; // 唯一标识
        this.group = 0; // 阵营
        this.id = 0; // monsterConf中id
        this.star = 0; // 星级
        this.order = 0; // 站位
        this.spine = ""; // 使用的spine名
        this.job = 0; // 职业
        this.elem = 0; // 元素
        this.distance = 0; // 近战远程
        this.level = 0; // 等级
        this._transform = (_a = {},
            _a[BattleConst_1.UNIT_STATE.IDLE] = (_b = {},
                _b[BattleConst_1.UNIT_STATE.SKILL] = true,
                _b[BattleConst_1.UNIT_STATE.DIE] = true,
                _b[BattleConst_1.UNIT_STATE.REBORN_DIE] = true,
                _b),
            _a[BattleConst_1.UNIT_STATE.SKILL] = (_c = {},
                _c[BattleConst_1.UNIT_STATE.IDLE] = true,
                _c[BattleConst_1.UNIT_STATE.DIE] = true,
                _c[BattleConst_1.UNIT_STATE.REBORN_DIE] = true,
                _c),
            _a[BattleConst_1.UNIT_STATE.REBORN_DIE] = (_d = {},
                _d[BattleConst_1.UNIT_STATE.WAIT_REBORN] = true,
                _d),
            _a[BattleConst_1.UNIT_STATE.WAIT_REBORN] = (_e = {},
                _e[BattleConst_1.UNIT_STATE.REBORN] = true,
                _e),
            _a[BattleConst_1.UNIT_STATE.REBORN] = (_f = {},
                _f[BattleConst_1.UNIT_STATE.IDLE] = true,
                _f),
            _a);
        this.bCtrl = ctrl;
        this.group = group;
        this.id = fighter.Id;
        this.star = fighter.Star || 1;
        this.order = fighter.Pos || 0;
        this.level = fighter.Lv;
        this.seq = BattleUnit.unique;
        this.noAttack = 0;
        this.noSkill = 0;
        this.invincible = 0;
        this.critNot = 0;
        this.mustCrit = 0;
        this.firstAct = 0;
        this.confusion = 0;
        this.atkTargetBack = 0;
        this.taunt = 0;
        this.immuneDeath = 0;
        this.immuneDeathCnt = 0;
        this.rebornNot = 0;
        this.actCnt = 0;
        this.buffGrpCnt = {};
        this.immuneBuffs = {};
        this.buffSeqs = [];
        var monsterConf = Configs_1.Configs.monsterConf[this.id];
        var role = monsterConf.role[this.star - 1] || monsterConf.role[monsterConf.role.length - 1];
        if (fighter.Skin) {
            var sConf = Configs_1.Configs.heroSkinConf[fighter.Skin];
            if (sConf) {
                role = sConf.role;
            }
        }
        var roleConf = Configs_1.Configs.roleConf[role];
        this.spine = roleConf.spine;
        this.job = monsterConf.jobId;
        this.elem = monsterConf.elem;
        this.distance = monsterConf.distance;
        // 初始化属性
        this.attrs = {};
        for (var attrId in Configs_1.Configs.attributeConf) {
            var conf_1 = Configs_1.Configs.attributeConf[attrId];
            this.attrs[conf_1.attributeId] = 0;
        }
        for (var attrId in fighter.Props) {
            this.attrs[attrId] = fighter.Props[attrId];
        }
        for (var attrId in Configs_1.Configs.attributeConf) {
            var conf_2 = Configs_1.Configs.attributeConf[attrId];
            this._originAttr[conf_2.attributeId] = { a: this.attrs[conf_2.attributeId], d: 0, e: 0 };
        }
        var conf = Configs_1.Configs.globalBattleConf[1];
        this.mp = conf.initailEnergy;
        this.maxMp = conf.energyLimit;
        this.hp = this.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
        // 初始化技能
        var lvConf = Configs_1.Configs.heroUpConf[this.level];
        var starConf = Configs_1.Configs.heroStarUpConf[this.star];
        this._passiveSkills = [];
        this._normalSkill = new UnitSkill_1.UnitSkill(this, monsterConf.skillNormal);
        for (var i = 0; i < monsterConf.skills.length; i++) {
            var skillId = monsterConf.skills[i];
            var lv = 0;
            //通过等级星级计算技能id
            if (lvConf && starConf) {
                lv = (lvConf.skillUnlock >= i + 1) ? (starConf.skillLv[i] || 0) : 0;
                //取可以找到的最大技能id
                for (var j = lv; j >= 1; j--) {
                    skillId = monsterConf.skills[i] + j - 1;
                    if (Configs_1.Configs.skillConf[skillId]) {
                        break;
                    }
                }
            }
            if (lv == 0)
                continue;
            var skill = new UnitSkill_1.UnitSkill(this, skillId);
            if (skill.sType === BattleConst_1.BattleConst.SKILL_TYPE.SPECIAL) {
                this._specialSkill = skill;
            }
            else {
                this._passiveSkills.push(skill);
            }
        }
    }
    Object.defineProperty(BattleUnit, "unique", {
        get: function () {
            var limit = Math.pow(10, 9);
            BattleUnit._unique++;
            if (BattleUnit._unique > limit) {
                BattleUnit._unique -= limit;
            }
            return BattleUnit._unique;
        },
        enumerable: false,
        configurable: true
    });
    BattleUnit.prototype.update = function (dt) {
        if (this._state === BattleConst_1.UNIT_STATE.SKILL) {
            if (this._curSkill.isComplete) {
                this.switchState(BattleConst_1.UNIT_STATE.IDLE);
                this.bCtrl.callSceneFunc("hideShade");
                this.bCtrl.callSceneFunc("unitSkillComplete", this._curSkill);
                if (this._curSkill.isAct) {
                    this._curSkill.isAct = false;
                    if (this.actCnt === 1) {
                        if (this.bCtrl.buffCtrl.buff_actExtraAct(this)) {
                            this.bCtrl.pushExtraSkill(BattleConst_1.BattleConst.BUFF_FUNC.ACT_EXTRA_ACT, this, "");
                        }
                    }
                }
                this._curSkill = null;
            }
            else {
                this._curSkill.update(dt);
            }
        }
        else if (this._state === BattleConst_1.UNIT_STATE.REBORN_DIE) {
            this._rebornDt = 0;
            this.switchState(BattleConst_1.UNIT_STATE.WAIT_REBORN);
        }
        else if (this._state === BattleConst_1.UNIT_STATE.REBORN) {
            this._rebornDt += dt;
            if (this._rebornDt > 1.9) { //reborn动画播放时间1.84
                this._rebornDt = 0;
                this.switchState(BattleConst_1.UNIT_STATE.IDLE);
            }
        }
    };
    BattleUnit.prototype.stateReborn = function (ratio) {
        this.forceUpdateHp(ratio * this.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE]);
        this.switchState(BattleConst_1.UNIT_STATE.REBORN);
        this.bCtrl.callSceneFunc("unitReborn", this);
    };
    BattleUnit.prototype.switchState = function (newState) {
        if (this._transform[this._state] && this._transform[this._state][newState]) {
            this.bCtrl.callSceneFunc("switchState", this, this._state, newState);
            this._state = newState;
            return true;
        }
        return false;
    };
    BattleUnit.prototype.updateControlType = function (noAttack, noSkill) {
        this.noAttack += noAttack;
        this.noSkill += noSkill;
    };
    BattleUnit.prototype.actStart = function () {
        this.actCnt++;
        if (this._normalSkill) {
            this._normalSkill.nextCnt = -1;
        }
        if (this._specialSkill) {
            this._specialSkill.nextCnt = -1;
        }
        var skills = [this._specialSkill, this._normalSkill];
        for (var i = 0; i < skills.length; i++) {
            var skill = skills[i];
            if (skill && skill.valid()) {
                if (this.switchState(BattleConst_1.UNIT_STATE.SKILL)) {
                    var sId = this.bCtrl.buffCtrl.buff_actNewSkill(this, skill.sType);
                    if (sId !== 0) {
                        this._curSkill = new UnitSkill_1.UnitSkill(this, sId);
                    }
                    else {
                        this._curSkill = skill;
                    }
                    this._curSkill.isAct = true;
                    this._curSkill.start();
                    return this._curSkill;
                }
            }
        }
    };
    BattleUnit.prototype.useSpecialSkill = function (nextCnt) {
        var skill = this._specialSkill;
        if (nextCnt && nextCnt !== 0) {
            skill.nextCnt = nextCnt;
        }
        if (skill.nextCnt === 0)
            return;
        if (skill && this.noSkill <= 0) {
            if (this.switchState(BattleConst_1.UNIT_STATE.SKILL)) {
                this._curSkill = skill;
                skill.start();
                skill.nextCnt = skill.nextCnt == -1 ? -1 : skill.nextCnt - 1;
                return skill;
            }
        }
    };
    BattleUnit.prototype.useNormalSkill = function (nextCnt) {
        var skill = this._normalSkill;
        if (nextCnt && nextCnt !== 0) {
            skill.nextCnt = nextCnt;
        }
        if (skill.nextCnt === 0)
            return;
        if (skill && this.noAttack <= 0) {
            if (this.switchState(BattleConst_1.UNIT_STATE.SKILL)) {
                this._curSkill = skill;
                skill.start();
                skill.nextCnt = skill.nextCnt == -1 ? -1 : skill.nextCnt - 1;
                return skill;
            }
        }
    };
    BattleUnit.prototype.usePassiveSkill = function () {
        for (var i = 0; i < this._passiveSkills.length; i++) {
            var skill = this._passiveSkills[i];
            if (skill.isComplete)
                continue;
            if (this.switchState(BattleConst_1.UNIT_STATE.SKILL)) {
                this._curSkill = skill;
                skill.start();
                return skill;
            }
        }
        return undefined;
    };
    BattleUnit.prototype.useSkill = function (id, nextCnt) {
        var skill = new UnitSkill_1.UnitSkill(this, id);
        skill.nextCnt = nextCnt;
        if ((skill.sType === BattleConst_1.BattleConst.SKILL_TYPE.SPECIAL && this.noSkill > 0) ||
            (skill.sType === BattleConst_1.BattleConst.SKILL_TYPE.NORMAL && this.noAttack > 0)) {
            return;
        }
        if (this.switchState(BattleConst_1.UNIT_STATE.SKILL)) {
            this._curSkill = skill;
            skill.start();
            return skill;
        }
    };
    BattleUnit.prototype.replaceSkill = function (sType, skillId) {
        if (sType === BattleConst_1.BattleConst.SKILL_TYPE.NORMAL) {
            this._normalSkill = new UnitSkill_1.UnitSkill(this, skillId);
        }
        else if (sType === BattleConst_1.BattleConst.SKILL_TYPE.SPECIAL) {
            this._specialSkill = new UnitSkill_1.UnitSkill(this, skillId);
        }
    };
    BattleUnit.prototype.forceUpdateHp = function (value) {
        this.hp = Math.min(value, this.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE]);
        this.bCtrl.callSceneFunc("updateHp", this, 0, null, null, false);
    };
    BattleUnit.prototype.updateHp = function (delta, skill, buff, crit) {
        var totalDmg = delta;
        var totalReduceHp = 0;
        if (!this.isAlive())
            return totalReduceHp;
        // 无敌
        if (this.invincible > 0)
            return totalReduceHp;
        var atker = skill ? skill.getOwner() : buff.getAtker();
        var atkDeath = false;
        if (delta > 0 && skill && this.order < 6) { //触发条件：伤害、技能伤害、非boss
            atkDeath = this.bCtrl.buffCtrl.buff_atkDeath(skill.getOwner(), skill.sType, this.hp / this.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE]);
        }
        // 斩杀
        if (delta > 0 && atkDeath) {
            if (this.immuneDeath > 0 && this.immuneDeathCnt === 0) {
                this.immuneDeathCnt++;
                return totalReduceHp;
            }
            else {
                totalDmg = Math.max(this.hp, totalDmg);
                totalReduceHp = this.hp;
                delta = totalDmg;
            }
        }
        else {
            // 护盾
            if (delta > 0) {
                var newDmg = this.bCtrl.buffCtrl.buff_attrToShield(this, delta);
                // 统计护盾伤害
                if (newDmg < delta) {
                    this.bCtrl.countCtrl.countDmg(atker, this, delta - newDmg);
                }
                if (newDmg <= 0) {
                    this.bCtrl.callSceneFunc("updateHp", this, totalDmg, skill, buff, crit);
                    return totalReduceHp;
                }
                else {
                    delta = newDmg;
                }
            }
            else {
                totalReduceHp = delta;
            }
        }
        var hp = this.hp - delta;
        hp = this.bCtrl.unitHpRefresh(this, hp);
        if (hp < this.hp) {
            totalReduceHp += (this.hp - Math.max(hp, 0));
            this.bCtrl.countCtrl.countDmg(atker, this, totalReduceHp);
        }
        else if (hp > this.hp) {
            var cureHp = Math.min(hp, this.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE]) - this.hp;
            if (cureHp > 0)
                this.bCtrl.countCtrl.countCure(atker, cureHp);
        }
        //世界boss打不死
        if (this.bCtrl.bType == BattleConst_1.BattleConst.BATTLE_TYPE.WORLD_BOSS && this.group == BattleConst_1.BattleConst.UNIT_GROUP.DEFENDER) {
            hp = this.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
        }
        if (hp <= 0) {
            this.hp = 0;
            // 复活
            var canReborn = false;
            if (this.rebornNot <= 0 && this.rebornCnt <= 0) {
                var ratio = this.bCtrl.buffCtrl.buff_reborn(this);
                canReborn = (ratio > 0);
                if (canReborn) {
                    this.switchState(BattleConst_1.UNIT_STATE.REBORN_DIE);
                    this.rebornCnt++;
                    this.bCtrl.pushExtraSkill(BattleConst_1.BattleConst.BUFF_FUNC.REBORN, this, String(ratio));
                    this.bCtrl.buffCtrl.deathDropBuff(this);
                }
            }
            if (!canReborn) {
                this.switchState(BattleConst_1.UNIT_STATE.DIE);
                this.bCtrl.unitDie(this);
            }
            this.bCtrl.callSceneFunc("unitDie", this);
        }
        else {
            this.hp = Math.min(hp, this.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE]);
        }
        var isVampire = (buff && buff.getFunc().funcName == BattleConst_1.BattleConst.BUFF_FUNC.VAMPIRE);
        if (totalDmg < 0 && !isVampire) {
            this.bCtrl.buffCtrl.buff_cureEnergy(this);
        }
        if (totalReduceHp > 0) {
            this.bCtrl.buffCtrl.buff_hpReduceBuff(this, totalReduceHp);
            this.bCtrl.buffCtrl.buff_hpPointBuff(this, totalReduceHp);
        }
        this.bCtrl.buffCtrl.buff_attrCond();
        this.bCtrl.callSceneFunc("updateHp", this, totalDmg, skill, buff, crit);
        return totalReduceHp;
    };
    BattleUnit.prototype.updateMp = function (delta) {
        if (!this.isAlive())
            return;
        if (delta > 0) {
            delta *= (1 + this.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.ENERTY_RATIO]);
        }
        this.mp += delta;
        this.mp = Math.max(this.mp, 0);
        this.mp = Math.min(this.maxMp, this.mp);
        this.bCtrl.callSceneFunc("updateMp", this);
    };
    BattleUnit.prototype.updateBuffAttr = function (attrId, val) {
        var d = 0;
        var e = 0;
        if (attrId >= 1000) {
            d += val;
            attrId = attrId / 100;
        }
        else {
            e += val;
        }
        var attr = this._originAttr[attrId];
        var d1 = attr.d + d;
        var e1 = attr.e + e;
        var value = attr.a * (1 + d1) + e1;
        // 当前血量根据血量上限发生变化
        if (attrId == BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE && this.hp > 0) {
            this.hp = value * (this.hp / this.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE]);
        }
        this._originAttr[attrId].d = d1;
        this._originAttr[attrId].e = e1;
        this.attrs[attrId] = value;
    };
    BattleUnit.prototype.isAlive = function () {
        return (this.hp > 0 && this._state != BattleConst_1.UNIT_STATE.DIE);
    };
    BattleUnit.prototype.getState = function () {
        return this._state;
    };
    BattleUnit.prototype.hasBuffType = function (bTypes) {
        for (var i = bTypes.length - 1; i >= 0; i--) {
            var bType = bTypes[i];
            if (this.buffTypeCnt[bType])
                return true;
        }
        return false;
    };
    BattleUnit.prototype.getFixedCritDam = function () {
        var ret = -1;
        // 暴击伤害锁定
        ret = this.bCtrl.buffCtrl.buff_critDamFixed(this);
        if (ret === -1) {
            ret = this.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.CRIT_DAM] + this.bCtrl.buffCtrl.buff_crit(this);
        }
        return ret;
    };
    BattleUnit.prototype.getSkills = function () {
        var skillIds = [];
        if (this._normalSkill)
            skillIds.push(this._normalSkill.id);
        if (this._specialSkill)
            skillIds.push(this._specialSkill.id);
        for (var i = 0; i < this._passiveSkills.length; i++) {
            skillIds.push(this._passiveSkills[i].id);
        }
        return skillIds;
    };
    BattleUnit._unique = 1;
    return BattleUnit;
}());
exports.BattleUnit = BattleUnit;
