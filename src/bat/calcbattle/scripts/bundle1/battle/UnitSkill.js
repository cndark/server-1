"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-03-29 20:04:45
 * @LastEditors: zyb
 * @LastEditTime: 2021-07-01 17:49:02
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.UnitSkill = void 0;
var Configs_1 = require("../../bundle0/framework/config/Configs");
var BattleConst_1 = require("../../bundle0/battleCommon/BattleConst");
var SkillBullet_1 = require("./SkillBullet");
var SkillFunc_1 = require("./SkillFunc");
var ActorTrigger_1 = require("./spine/ActorTrigger");
var UnitSkill = /** @class */ (function () {
    function UnitSkill(owner, id) {
        this._owner = null;
        this._bullet = null;
        this._funcs = null;
        this._duration = 0;
        this._extraDuration = 0;
        this._dt = 0;
        this._events = [];
        this.killed = false; // 击杀目标
        this.nextCnt = -1; // 触发上线
        this.isAct = false; // 英雄行动
        this.isComplete = false; // 技能完成
        this.id = 0;
        this.sType = "";
        this.anim = "";
        this.distance = 0;
        this.moveTime = 0.1;
        this._owner = owner;
        this.id = id;
        var sConf = Configs_1.Configs.skillConf[this.id];
        this.sType = sConf.skillType;
        this.anim = sConf.skillAnim;
        this.distance = sConf.distance;
    }
    UnitSkill.prototype.start = function () {
        this._dt = 0;
        this._duration = 0;
        this._events = [];
        this._bullet = null;
        this.killed = false;
        this.isComplete = false;
        this._funcs = [];
        var sConf = Configs_1.Configs.skillConf[this.id];
        var inheritTargets = [];
        for (var i = 0; i < sConf.skillFunc.length; i++) {
            var func_1 = new SkillFunc_1.SkillFunc(this, sConf.skillFunc[i], inheritTargets);
            this._funcs.push(func_1);
            inheritTargets = func_1.inheritTargets;
        }
        var func = this._funcs[0];
        if (!func || !func.isValid || func.targets.length <= 0) {
            return;
        }
        this._events.push({
            name: "animStart",
            time: 0,
        });
        var tConf = ActorTrigger_1.ActorTrigger[this._owner.spine];
        var data = tConf[this.anim];
        if (data) {
            this._duration = data.duration;
            for (var i = 0; i < data.events.length; i++) {
                var event_1 = data.events[i];
                this._events.push({
                    name: event_1.name,
                    time: event_1.time,
                });
            }
        }
        else {
            this._events.push({
                name: "e_gj",
                time: 0,
            });
        }
        if (this.distance === 1) {
            this._duration += this.moveTime * 2;
            for (var i = 0; i < this._events.length; i++) {
                var event_2 = this._events[i];
                event_2.time += this.moveTime;
            }
            this._events.splice(0, 0, {
                name: "moveForward",
                time: 0,
            });
            this._events.push({
                name: "moveBack",
                time: this._duration - this.moveTime,
            });
        }
        if (this.sType === BattleConst_1.BattleConst.SKILL_TYPE.SPECIAL) {
            this._extraDuration = 0.2;
        }
        else {
            this._extraDuration = 0;
        }
        if (this.isAct) {
            if (this._owner.bCtrl.buffCtrl.buff_actExtraSpecial(this._owner, this)) {
                this._owner.bCtrl.pushExtraSkill(BattleConst_1.BattleConst.BUFF_FUNC.ACT_EXTRA_SPECIAL, this._owner, BattleConst_1.BattleConst.SKILL_TYPE.SPECIAL);
            }
        }
        // 界面显示相关
        var seqMap = {};
        for (var i = 0; i < this._funcs.length; i++) {
            var func_2 = this._funcs[i];
            for (var j = 0; j < func_2.targets.length; j++) {
                var target = func_2.targets[j];
                seqMap[target.seq] = true;
            }
        }
        this._owner.bCtrl.callSceneFunc("unitSkillStart", this, Object.keys(seqMap));
        if (this.sType === BattleConst_1.BattleConst.SKILL_TYPE.SPECIAL) {
            this._owner.bCtrl.callSceneFunc("showShade");
        }
    };
    UnitSkill.prototype.nextSkill = function () {
        var sConf = Configs_1.Configs.skillConf[this.id];
        if (sConf.nextSkill === "" || this.nextCnt === 0)
            return;
        var info = sConf.nextSkill.split("~");
        if (info[0] === "kill") {
            var ratio = Number(info[1]);
            if (this.killed && this._owner.bCtrl.random() < ratio) {
                var id = Number(info[2]);
                var num = Number(info[3]);
                var nextCnt = this.nextCnt == -1 ? num : Math.min(num, this.nextCnt);
                return this._owner.useSkill(id, nextCnt);
            }
        }
    };
    UnitSkill.prototype.executeFuncs = function () {
        var dmg = 0;
        for (var i = 0; i < this._funcs.length; i++) {
            var func = this._funcs[i];
            // 检验condition
            if (func.isValid) {
                var conf = Configs_1.Configs.skillFuncConf[func.id];
                var info = conf.condition.split("_");
                if (info[0] === "killed") {
                    var killed = (i > 0) ? this._funcs[i - 1].killed : false;
                    if ((info[1] === "0" && killed === true) || (info[1] === "1" && killed === false)) {
                        func.isValid = false;
                    }
                }
                else if (info[0] === "buffType") {
                    var valid = false;
                    var units = [];
                    if (info[1] === "self") {
                        units = [this._owner];
                    }
                    else if (info[1] === "target") {
                        units = func.targets;
                    }
                    for (var j = 0; j < units.length; j++) {
                        var target = units[j];
                        if (target.hasBuffType([Number(info[2])])) {
                            valid = true;
                            break;
                        }
                    }
                    if (!valid) {
                        func.isValid = false;
                    }
                }
                else if (info[0] === "heroLive") {
                    var heroLive = false;
                    var heroId = Number(info[1]);
                    var units = this._owner.bCtrl.getUnits();
                    for (var j = 0; j < units.length; j++) {
                        var unit = units[j];
                        if (unit.id === heroId && unit.group === this._owner.group && unit.isAlive()) {
                            heroLive = true;
                            break;
                        }
                    }
                    if (!heroLive)
                        func.isValid = false;
                }
                else if (info[0] === "pos") {
                    if ((info[1] === "front" && this._owner.order > 1) ||
                        (info[1] === "back" && this._owner.order <= 1)) {
                        func.isValid = false;
                    }
                }
                else if (info[0] === "skillKill") {
                    if (!this.killed) {
                        func.isValid = false;
                    }
                }
            }
            dmg = func.execute(dmg);
            if (func.killed)
                this.killed = true;
        }
    };
    UnitSkill.prototype.update = function (dt) {
        if (this.isComplete)
            return;
        this._dt += dt;
        var bulletComplete = true;
        if (this._bullet) {
            bulletComplete = this._bullet.update(dt);
        }
        while (this._events.length > 0) {
            var event_3 = this._events[0];
            if (this._dt >= event_3.time) {
                if (event_3.name == "e_gj") {
                    if (this.isAct) {
                        var conf = Configs_1.Configs.globalBattleConf[1];
                        if (this.sType === BattleConst_1.BattleConst.SKILL_TYPE.NORMAL) {
                            this._owner.updateMp(conf.normalToEnergy);
                        }
                        else if (this.sType === BattleConst_1.BattleConst.SKILL_TYPE.SPECIAL) {
                            this._owner.updateMp(this._owner.mp * -1);
                        }
                    }
                    this._bullet = new SkillBullet_1.SkillBullet(this);
                    bulletComplete = false;
                    this._owner.bCtrl.callSceneFunc("addSkillBullet", this._bullet);
                }
                else if (event_3.name == "moveForward") {
                    this._owner.bCtrl.callSceneFunc("moveForward", this);
                }
                else if (event_3.name == "moveBack") {
                    this._owner.bCtrl.callSceneFunc("moveBack", this);
                }
                else if (event_3.name == "animStart") {
                    this._owner.bCtrl.callSceneFunc("skillAnimStart", this);
                }
                this._events.shift();
            }
            else {
                break;
            }
        }
        var isComplete = (this._duration < this._dt && bulletComplete);
        if (isComplete) {
            this._extraDuration -= dt;
        }
        this.isComplete = isComplete && this._extraDuration <= 0;
    };
    UnitSkill.prototype.valid = function () {
        if (this.sType === BattleConst_1.BattleConst.SKILL_TYPE.SPECIAL) {
            if (this._owner.noSkill > 0)
                return false;
            return this._owner.mp >= this._owner.maxMp;
        }
        else if (this.sType === BattleConst_1.BattleConst.SKILL_TYPE.NORMAL) {
            if (this._owner.noAttack > 0)
                return false;
            return true;
        }
        else {
            return true;
        }
    };
    UnitSkill.prototype.getOwner = function () {
        return this._owner;
    };
    UnitSkill.prototype.getTargets = function () {
        var func = this._funcs[0];
        if (func) {
            return func.targets;
        }
        else {
            return [];
        }
    };
    return UnitSkill;
}());
exports.UnitSkill = UnitSkill;
