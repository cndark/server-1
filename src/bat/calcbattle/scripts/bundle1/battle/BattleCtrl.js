"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-03-31 10:33:00
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-07-23 17:45:53
 */
var __spreadArray = (this && this.__spreadArray) || function (to, from) {
    for (var i = 0, il = from.length, j = to.length; i < il; i++, j++)
        to[j] = from[i];
    return to;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.BattleCtrl = void 0;
var Configs_1 = require("../../bundle0/framework/config/Configs");
var BattleConst_1 = require("../../bundle0/battleCommon/BattleConst");
var BattleRandom_1 = require("./BattleRandom");
var BattleUnit_1 = require("./BattleUnit");
var BuffCtrl_1 = require("./BuffCtrl");
var CountCtrl_1 = require("./CountCtrl");
var BattleCtrl = /** @class */ (function () {
    function BattleCtrl() {
        this._data = null;
        this._scene = null;
        this._units = null;
        this._random = null;
        this._readyIdxs = [];
        this._round = 1;
        this._maxRound = 0;
        this._actSkills = {};
        this._speed = 1;
        this._bType = "";
        this._atkCnt = 0;
        this._defCnt = 0;
        this._state = BattleConst_1.BATTLE_STATE.NONE;
        this.countCtrl = null;
        this.buffCtrl = null;
        this.pause = false;
        this.recordCtrl = null;
    }
    Object.defineProperty(BattleCtrl.prototype, "speed", {
        get: function () {
            return this._speed;
        },
        set: function (val) {
            if (this._speed !== val) {
                this._speed = val;
                this.callSceneFunc("updateSpeed", this._speed);
            }
        },
        enumerable: false,
        configurable: true
    });
    BattleCtrl.prototype.changeSpeed = function () {
        var value = this.speed + 1;
        if (value > 2)
            value = 1;
        this.speed = value;
    };
    Object.defineProperty(BattleCtrl.prototype, "bType", {
        get: function () {
            return this._bType;
        },
        enumerable: false,
        configurable: true
    });
    Object.defineProperty(BattleCtrl.prototype, "atkCnt", {
        get: function () {
            return this._atkCnt;
        },
        enumerable: false,
        configurable: true
    });
    Object.defineProperty(BattleCtrl.prototype, "defCnt", {
        get: function () {
            return this._defCnt;
        },
        enumerable: false,
        configurable: true
    });
    Object.defineProperty(BattleCtrl.prototype, "state", {
        get: function () {
            return this._state;
        },
        enumerable: false,
        configurable: true
    });
    BattleCtrl.prototype._reset = function (seed) {
        this._units = [];
        this._random = new BattleRandom_1.BattleRandom(seed);
        this._state = BattleConst_1.BATTLE_STATE.NONE;
        this._atkCnt = 0;
        this._defCnt = 0;
        this._actSkills = {};
        this.pause = false;
    };
    BattleCtrl.prototype.init = function (data, scene) {
        var seed = Number(data.Args["seed"] || "1");
        this._reset(seed);
        this._data = data;
        this._scene = scene;
        this._bType = data.Args.Module;
        this.buffCtrl = new BuffCtrl_1.BuffCtrl(this);
        this.countCtrl = new CountCtrl_1.CountCtrl(this);
        this.initUnits();
        this.initElementAdd(); //光环属性加成
    };
    BattleCtrl.prototype.initUnits = function () {
        // this._data.Args中，key以init_hp开头("init_hp.group.pos": "decHp")，需要设置单位初始血量百分比，损失为1的单位不加入战斗
        // 例如："init_hp.2.5": "0.82" 表示 敌方 位置5 单位损失血量0.82
        for (var i = 0; i < this._data.T1.Fighters.length; i++) {
            var fighter = this._data.T1.Fighters[i];
            //损失血量
            var decHpKey = "init_hp." + BattleConst_1.BattleConst.UNIT_GROUP.ATTACKER + "." + fighter.Pos;
            var decHp = Number(this._data.Args[decHpKey]) || 0;
            //血量损失不加入战斗单位
            if (decHp >= 1)
                continue;
            var unit = new BattleUnit_1.BattleUnit(this, fighter, BattleConst_1.BattleConst.UNIT_GROUP.ATTACKER);
            unit.hp *= (1 - (decHp || 0));
            this.addUnit(unit);
        }
        for (var i = 0; i < this._data.T2.Fighters.length; i++) {
            var fighter = this._data.T2.Fighters[i];
            //损失血量
            var decHpKey = "init_hp." + BattleConst_1.BattleConst.UNIT_GROUP.DEFENDER + "." + fighter.Pos;
            var decHp = Number(this._data.Args[decHpKey]) || 0;
            //血量全部损失不加入战斗单位
            if (decHp >= 1)
                continue;
            var unit = new BattleUnit_1.BattleUnit(this, fighter, BattleConst_1.BattleConst.UNIT_GROUP.DEFENDER);
            unit.hp *= (1 - (decHp || 0));
            this.addUnit(unit);
        }
        var gbConf = Configs_1.Configs.globalBattleConf[1];
        this._maxRound = gbConf.battleRound[Number(this._data.Args.RoundType)] || gbConf.battleRound[0];
    };
    BattleCtrl.prototype.initElementAdd = function () {
        var groups = [BattleConst_1.BattleConst.UNIT_GROUP.ATTACKER, BattleConst_1.BattleConst.UNIT_GROUP.DEFENDER];
        for (var i = 0; i < groups.length; i++) {
            var group = groups[i];
            var elemMap = {};
            for (var j = 0; j < this._units.length; j++) {
                var unit = this._units[j];
                if (unit.group === group) {
                    var elem = unit.elem;
                    if (!elemMap[elem]) {
                        elemMap[elem] = 1;
                    }
                    else {
                        elemMap[elem]++;
                    }
                }
            }
            var totalProps = {};
            for (var key in Configs_1.Configs.elementAddConf) {
                var conf = Configs_1.Configs.elementAddConf[key];
                var validCnt = 0;
                for (var k = 0; k < conf.element.length; k++) {
                    var element = conf.element[k];
                    if (elemMap[element] && elemMap[element] == conf.num) {
                        validCnt++;
                    }
                }
                if (validCnt > 0) {
                    for (var k = 0; k < conf.prop.length; k++) {
                        var prop = conf.prop[k];
                        if (!totalProps[prop.id])
                            totalProps[prop.id] = 0;
                        totalProps[prop.id] += prop.val * validCnt;
                    }
                }
            }
            for (var propId in totalProps) {
                var val = totalProps[propId];
                var attrId = Number(propId);
                for (var j = 0; j < this._units.length; j++) {
                    var unit = this._units[j];
                    if (unit.group === group) {
                        unit.updateBuffAttr(attrId, val);
                    }
                }
            }
        }
    };
    BattleCtrl.prototype.initStartBuffs = function () {
        // this._data.Args中，key以buffs开头("buffs.group.pos": "buffId")，需要添加buff
        // 例如："buffs.1.3": "1000" 表示 我方 位置3 添加buff 1000
        //       "buffs.2" : "1000" 表示 敌方 全体 添加buff 1000
        var atkGroupKey = "buffs." + BattleConst_1.BattleConst.UNIT_GROUP.ATTACKER;
        var atkGroupBuffStr = this._data.Args[atkGroupKey];
        var atkGroupBuffs = [];
        if (atkGroupBuffStr) {
            atkGroupBuffs = atkGroupBuffStr.split(',');
        }
        var dfdGroupKey = "buffs." + BattleConst_1.BattleConst.UNIT_GROUP.DEFENDER;
        var dfdGroupBuffStr = this._data.Args[dfdGroupKey];
        var dfdGroupBuffs = [];
        if (dfdGroupBuffStr) {
            dfdGroupBuffs = dfdGroupBuffStr.split(',');
        }
        for (var i = this._units.length - 1; i >= 0; i--) {
            var unit = this._units[i];
            if (unit.isAlive()) {
                var key = "buffs." + unit.group + "." + unit.order;
                var buffStr = this._data.Args[key];
                var orderBuffs = [];
                if (buffStr) {
                    orderBuffs = buffStr.split(',');
                }
                var groupBuff = unit.group == BattleConst_1.BattleConst.UNIT_GROUP.ATTACKER ? atkGroupBuffs : dfdGroupBuffs;
                for (var j = 0; j < groupBuff.length; j++) {
                    var buffId = Number(groupBuff[j]);
                    this.buffCtrl.addBuff(unit, buffId, unit, 1, false);
                }
                for (var j = 0; j < orderBuffs.length; j++) {
                    var buffId = Number(orderBuffs[j]);
                    this.buffCtrl.addBuff(unit, buffId, unit, 1, false);
                }
            }
        }
    };
    BattleCtrl.prototype.logicStart = function () {
        this.callSceneFunc("recordKeyWord", "logicStart begin");
        var finish = true;
        while (this._readyIdxs.length > 0) {
            var idx = this._readyIdxs[this._readyIdxs.length - 1];
            var unit = this._units[idx];
            if (unit.isAlive()) {
                var skill = unit.usePassiveSkill();
                if (skill) {
                    // 【需求】被动技能一帧内完成
                    // 防止死循环，增加一帧内最多执行10次
                    for (var i = 0; i < 10; i++) {
                        if (unit.getState() == BattleConst_1.UNIT_STATE.SKILL) {
                            unit.update(0.016);
                        }
                        else {
                            break;
                        }
                    }
                    finish = false;
                    break;
                }
            }
            this._readyIdxs.pop();
        }
        if (finish) {
            this._state = BattleConst_1.BATTLE_STATE.ROUND_START;
        }
        this.callSceneFunc("recordKeyWord", "logicStart end");
    };
    BattleCtrl.prototype.logicRoundStart = function () {
        this.callSceneFunc("recordKeyWord", "logicRoundStart begin");
        var baseRound = BattleConst_1.BattleConst.BATTLE_BASE_ROUND;
        if (this._round > baseRound) { //超过15回合之后 增加额外buff
            for (var key in Configs_1.Configs.battleRoundBuffConf) {
                var conf = Configs_1.Configs.battleRoundBuffConf[key];
                if (conf.model == this.bType) {
                    var extRound = this._round - baseRound;
                    var buffStr = conf.buff[extRound - 1];
                    if (buffStr !== "") {
                        var buffs = buffStr.split('~');
                        for (var m = 0; m < conf.type.length; m++) {
                            for (var i = this._units.length - 1; i >= 0; i--) {
                                var unit = this._units[i];
                                if (unit.isAlive() && unit.group == conf.type[m]) {
                                    for (var b = 0; b < buffs.length; b++) {
                                        this.buffCtrl.addBuff(unit, Number(buffs[b]), unit);
                                    }
                                }
                            }
                        }
                    }
                    break;
                }
            }
        }
        // 回合开始换血
        this.buffCtrl.buff_roundChangeHp();
        // dot生效
        for (var i = this._units.length - 1; i >= 0; i--) {
            var unit = this._units[i];
            if (unit.isAlive()) {
                this.buffCtrl.buff_dot(unit);
                this.buffCtrl.buff_dotGamble(unit);
            }
        }
        this.buffCtrl.buff_roundStartDmgBuff();
        // 回合开始加buff
        this.buffCtrl.buff_roundStartBuff();
        this._state = BattleConst_1.BATTLE_STATE.SORT_ACT;
        this.callSceneFunc("recordKeyWord", "logicRoundStart end");
    };
    BattleCtrl.prototype.logicSortAct = function () {
        var _this = this;
        this.callSceneFunc("recordKeyWord", "logicSortAct begin");
        // 排序
        this._readyIdxs = [];
        for (var i = this._units.length - 1; i >= 0; i--) {
            var unit = this._units[i];
            unit.actCnt = 0;
            this._readyIdxs.push(i);
        }
        this._readyIdxs.sort(function (a, b) {
            var u1 = _this._units[a];
            var u2 = _this._units[b];
            if (u1.group !== u2.group) {
                if (_this._bType == BattleConst_1.BattleConst.BATTLE_TYPE.WAR_CUP && u1.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.SPEED] == u2.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.SPEED]) {
                    return _this.random() - 0.5;
                }
                else {
                    return u1.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.SPEED] - u2.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.SPEED];
                }
            }
            else {
                if (u1.firstAct > 0 && u2.firstAct > 0) {
                    return u1.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.SPEED] - u2.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.SPEED];
                }
                else if (u1.firstAct) {
                    return 1;
                }
                else if (u2.firstAct) {
                    return -1;
                }
                else {
                    return u1.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.SPEED] - u2.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.SPEED];
                }
            }
        });
        this._state = BattleConst_1.BATTLE_STATE.FIGHT;
        this.callSceneFunc("recordKeyWord", "logicSortAct end");
    };
    BattleCtrl.prototype.logicFight = function (dt) {
        var finish = true;
        while (this._readyIdxs.length > 0) {
            var idx = this._readyIdxs[this._readyIdxs.length - 1];
            var unit = this._units[idx];
            this._readyIdxs.pop();
            if (unit.isAlive()) {
                // 行动开始时加buff
                this.buffCtrl.buff_actExtraBuff(unit);
                this.buffCtrl.buff_actStartDmgBuff(unit);
                this.buffCtrl.buff_actEnergy(unit);
                this.buffCtrl.buff_actStealEnergy(unit);
                unit.actStart();
                finish = false;
                break;
            }
        }
        if (finish) {
            this._state = BattleConst_1.BATTLE_STATE.ROUND_END;
        }
    };
    BattleCtrl.prototype.logicRoundEnd = function () {
        this.callSceneFunc("recordKeyWord", "logicRoundEnd start" + this._round);
        this.buffCtrl.buff_roundEndDmgBuff();
        // 回合结束加buff
        this.buffCtrl.buff_roundEndBuff();
        // 清除生命周期buff
        for (var i = this._units.length - 1; i >= 0; i--) {
            var unit = this._units[i];
            this.buffCtrl.updateRound(unit);
        }
        this.buffCtrl.clearRemoveBuff();
        // 回合数增加
        this._round++;
        if (this._round <= this._maxRound) {
            this._state = BattleConst_1.BATTLE_STATE.ROUND_START;
            this.callSceneFunc("updateRound", this._round);
        }
        this.callSceneFunc("recordKeyWord", "logicRoundEnd end" + this._round);
    };
    BattleCtrl.prototype.logicUpdate = function (dt) {
        if (this.pause)
            return;
        if (this.checkExtraActSkill())
            return;
        if (this._state == BattleConst_1.BATTLE_STATE.START) {
            this.logicStart();
        }
        else if (this._state == BattleConst_1.BATTLE_STATE.ROUND_START) {
            this.logicRoundStart();
        }
        else if (this._state == BattleConst_1.BATTLE_STATE.SORT_ACT) {
            this.logicSortAct();
        }
        else if (this._state == BattleConst_1.BATTLE_STATE.FIGHT) {
            this.logicFight(dt);
        }
        else if (this._state == BattleConst_1.BATTLE_STATE.ROUND_END) {
            this.logicRoundEnd();
        }
    };
    BattleCtrl.prototype.update = function (dt) {
        dt *= this.speed;
        this.buffCtrl.reboundLimit = {};
        var isBusy = false;
        var unitIdxs = [];
        for (var i = 0; i < this._units.length; i++) {
            unitIdxs.push(i);
        }
        //优先刷新使用技能的unit
        for (var i = unitIdxs.length - 1; i >= 0; i--) {
            var unit = this._units[unitIdxs[i]];
            if (unit.getState() === BattleConst_1.UNIT_STATE.SKILL) {
                isBusy = true;
                unit.update(dt);
                unitIdxs.splice(i, 1);
            }
        }
        //再刷新复活的unit
        for (var i = unitIdxs.length - 1; i >= 0; i--) {
            var unit = this._units[unitIdxs[i]];
            if (unit.getState() === BattleConst_1.UNIT_STATE.REBORN) {
                isBusy = true;
                unit.update(dt);
                unitIdxs.splice(i, 1);
            }
        }
        //最后刷新其他unit
        for (var i = unitIdxs.length - 1; i >= 0; i--) {
            var unit = this._units[unitIdxs[i]];
            unit.update(dt);
        }
        // 正在释放技能
        if (isBusy)
            return;
        // 战斗结束
        if (this._state === BattleConst_1.BATTLE_STATE.COMPLETE)
            return;
        if (this._round > this._maxRound || this._atkCnt <= 0 || this._defCnt <= 0) {
            this._state = BattleConst_1.BATTLE_STATE.COMPLETE;
            this.callSceneFunc("fightComplete", this._atkCnt, this._defCnt, false);
            return;
        }
        this.logicUpdate(dt);
    };
    BattleCtrl.prototype.skipBattle = function () {
        if (this._state == BattleConst_1.BATTLE_STATE.COMPLETE)
            return;
        this._state = BattleConst_1.BATTLE_STATE.COMPLETE;
        this.callSceneFunc("fightComplete", this._atkCnt, this._defCnt, true);
    };
    BattleCtrl.prototype.checkExtraActSkill = function () {
        if (this.updateNextSkill())
            return true;
        if (this.updateExtraSpecial())
            return true;
        if (this.updateExtraNormal())
            return true;
        if (this.updateDeathExtraSkills())
            return true;
        if (this.updateRebornUnits())
            return true;
        if (this.updateExtraAct())
            return true;
        return false;
    };
    // 死亡释放技能
    BattleCtrl.prototype.updateDeathExtraSkills = function () {
        var funcName = BattleConst_1.BattleConst.BUFF_FUNC.DEATH_EXTRA_SPECIAL;
        var info = this._actSkills[funcName];
        if (!info)
            return false;
        var seqs = Object.keys(info);
        if (seqs.length <= 0)
            return false;
        while (seqs.length > 0) {
            var seq = Number(seqs[0]);
            var atker = info[seq].atker;
            var sType = info[seq].sType;
            seqs.splice(0, 1);
            delete info[seq];
            if (atker.isAlive()) {
                var skill = (sType === BattleConst_1.BattleConst.SKILL_TYPE.NORMAL) ? atker.useNormalSkill() : atker.useSpecialSkill();
                if (skill) {
                    this.callSceneFunc("recordKeyWord", "DeathExtraSkills");
                    return true;
                }
            }
        }
        return false;
    };
    // 复活
    BattleCtrl.prototype.updateRebornUnits = function () {
        var funcName = BattleConst_1.BattleConst.BUFF_FUNC.REBORN;
        var info = this._actSkills[funcName];
        if (!info)
            return false;
        var seqs = Object.keys(info);
        if (seqs.length <= 0)
            return false;
        while (seqs.length > 0) {
            var seq = Number(seqs[0]);
            var atker = info[seq].atker;
            var ratio = Number(info[seq].sType);
            if (atker.getState() !== BattleConst_1.UNIT_STATE.WAIT_REBORN)
                break;
            seqs.splice(0, 1);
            delete info[seq];
            atker.stateReborn(ratio);
        }
        return true;
    };
    // 额外触发技能
    BattleCtrl.prototype.updateNextSkill = function () {
        var funcName = BattleConst_1.BattleConst.SKILL_FUNC.SKILL_NEXT;
        var info = this._actSkills[funcName];
        if (!info)
            return false;
        var seqs = Object.keys(info);
        if (seqs.length <= 0)
            return false;
        while (seqs.length > 0) {
            var seq = Number(seqs[0]);
            var atker = info[seq].atker;
            var sType = info[seq].sType;
            var nextCnt = info[seq].nextCnt;
            seqs.splice(0, 1);
            delete info[seq];
            if (atker.isAlive()) {
                var skill = (sType === BattleConst_1.BattleConst.SKILL_TYPE.NORMAL) ? atker.useNormalSkill(nextCnt) : atker.useSpecialSkill(nextCnt);
                if (skill) {
                    this.callSceneFunc("recordKeyWord", "NextSkill");
                    return true;
                }
            }
        }
        return false;
    };
    // 行动后放技能
    BattleCtrl.prototype.updateExtraSpecial = function () {
        var funcName = BattleConst_1.BattleConst.BUFF_FUNC.ACT_EXTRA_SPECIAL;
        var info = this._actSkills[funcName];
        if (!info)
            return false;
        var seqs = Object.keys(info);
        if (seqs.length <= 0)
            return false;
        while (seqs.length > 0) {
            var seq = Number(seqs[0]);
            var atker = info[seq].atker;
            seqs.splice(0, 1);
            delete info[seq];
            if (atker.isAlive()) {
                var skill = atker.useSpecialSkill();
                if (skill) {
                    this.callSceneFunc("recordKeyWord", "ExtraSpecial");
                    return true;
                }
            }
        }
        return false;
    };
    // 受击进行普攻
    BattleCtrl.prototype.updateExtraNormal = function () {
        var funcName = BattleConst_1.BattleConst.BUFF_FUNC.HIT_EXTRA_NORMAL;
        var info = this._actSkills[funcName];
        if (!info)
            return false;
        var seqs = Object.keys(info);
        if (seqs.length <= 0)
            return false;
        while (seqs.length > 0) {
            var seq = Number(seqs[0]);
            var atker = info[seq].atker;
            seqs.splice(0, 1);
            delete info[seq];
            if (atker.isAlive()) {
                var skill = atker.useNormalSkill();
                if (skill) {
                    this.callSceneFunc("recordKeyWord", "ExtraNormal");
                    return true;
                }
            }
        }
        return false;
    };
    // 额外行动
    BattleCtrl.prototype.updateExtraAct = function () {
        var funcName = BattleConst_1.BattleConst.BUFF_FUNC.ACT_EXTRA_ACT;
        var info = this._actSkills[funcName];
        if (!info)
            return false;
        var seqs = Object.keys(info);
        if (seqs.length <= 0)
            return false;
        while (seqs.length > 0) {
            var seq = Number(seqs[0]);
            var atker = info[seq].atker;
            seqs.splice(0, 1);
            delete info[seq];
            if (atker.isAlive()) {
                var skill = atker.actStart();
                if (skill) {
                    this.callSceneFunc("recordKeyWord", "ExtraAct");
                    return true;
                }
            }
        }
        return false;
    };
    BattleCtrl.prototype.stateStart = function () {
        this._readyIdxs = [];
        for (var i = 0; i < this._units.length; i++) {
            this._readyIdxs.push(i);
        }
        this.initStartBuffs(); //不同战斗附加属性buff
        this._state = BattleConst_1.BATTLE_STATE.START;
    };
    BattleCtrl.prototype.unitHpRefresh = function (unit, hp) {
        if (hp <= 0 && unit.immuneDeath > 0 && unit.immuneDeathCnt === 0) {
            unit.immuneDeathCnt++;
            return unit.hp;
        }
        return hp;
    };
    BattleCtrl.prototype.unitDie = function (unit) {
        this.buffCtrl.buff_deathDmg(unit);
        this.buffCtrl.buff_deathExtraBuff(unit);
        this.buffCtrl.buff_deathEnergy(unit);
        this.buffCtrl.deathDropBuff(unit);
        this.buffCtrl.buff_deathExtraSpecial(unit);
        switch (unit.group) {
            case BattleConst_1.BattleConst.UNIT_GROUP.ATTACKER:
                this._atkCnt--;
                break;
            case BattleConst_1.BattleConst.UNIT_GROUP.DEFENDER:
                this._defCnt--;
                break;
        }
    };
    BattleCtrl.prototype.addUnit = function (unit) {
        this._units.push(unit);
        if (unit.group == BattleConst_1.BattleConst.UNIT_GROUP.ATTACKER) {
            this._atkCnt++;
        }
        else {
            this._defCnt++;
        }
    };
    BattleCtrl.prototype.removeUnit = function (unit) {
        for (var i = this._units.length - 1; i >= 0; i--) {
            var u = this._units[i];
            if (u.seq === unit.seq) {
                this._units.splice(i, 1);
                break;
            }
        }
    };
    BattleCtrl.prototype.getUnits = function () {
        return this._units;
    };
    BattleCtrl.prototype.getMaxRound = function () {
        return this._maxRound;
    };
    BattleCtrl.prototype.pushExtraSkill = function (funcName, atker, sType, nextCnt) {
        var _a;
        var funcNameState = (_a = {},
            _a[BattleConst_1.BattleConst.BUFF_FUNC.ACT_EXTRA_SPECIAL] = BattleConst_1.BATTLE_STATE.FIGHT,
            _a[BattleConst_1.BattleConst.BUFF_FUNC.HIT_EXTRA_NORMAL] = BattleConst_1.BATTLE_STATE.FIGHT,
            _a[BattleConst_1.BattleConst.BUFF_FUNC.ACT_EXTRA_ACT] = BattleConst_1.BATTLE_STATE.FIGHT,
            _a[BattleConst_1.BattleConst.SKILL_FUNC.SKILL_NEXT] = BattleConst_1.BATTLE_STATE.FIGHT,
            _a);
        if (funcNameState[funcName] && funcNameState[funcName] != this._state) {
            return;
        }
        if (!this._actSkills[funcName])
            this._actSkills[funcName] = {};
        this._actSkills[funcName][atker.seq] = { atker: atker, sType: sType, nextCnt: nextCnt || 0 };
    };
    BattleCtrl.prototype.random = function (start, end) {
        return this._random.random(start, end);
    };
    //调用View层方法
    BattleCtrl.prototype.callSceneFunc = function (funcName) {
        var _a, _b;
        var param = [];
        for (var _i = 1; _i < arguments.length; _i++) {
            param[_i - 1] = arguments[_i];
        }
        if (this._scene && this._scene[funcName] && typeof this._scene[funcName] == 'function') {
            (_a = this._scene[funcName]).call.apply(_a, __spreadArray([this._scene], param));
        }
        if (this.recordCtrl && this.recordCtrl[funcName] && typeof this.recordCtrl[funcName] == 'function') {
            (_b = this.recordCtrl[funcName]).call.apply(_b, __spreadArray([this.recordCtrl], param));
        }
    };
    /** 随机数 */
    BattleCtrl.prototype.fRandomBy = function (under, over) {
        return Math.floor((Math.random() * (over - under + 1) + under));
    };
    return BattleCtrl;
}());
exports.BattleCtrl = BattleCtrl;
