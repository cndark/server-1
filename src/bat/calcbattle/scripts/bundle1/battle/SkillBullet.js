"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-03-29 20:12:45
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-05-19 10:31:58
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.SkillBullet = void 0;
var Configs_1 = require("../../bundle0/framework/config/Configs");
var BattleConst_1 = require("../../bundle0/battleCommon/BattleConst");
var EffectTrigger_1 = require("./spine/EffectTrigger");
var SkillBullet = /** @class */ (function () {
    function SkillBullet(skill) {
        this._skill = null;
        this._owner = null;
        this._dt = 0;
        this._duration = 0;
        this._events = [];
        this.bulletJson = "";
        this.bulletType = BattleConst_1.BattleConst.BULLET_TYPE.NONE;
        this.anim = "";
        this._skill = skill;
        this._owner = skill.getOwner();
        this._dt = 0;
        this._duration = 0;
        this._events = [];
        var sConf = Configs_1.Configs.skillConf[skill.id];
        this.bulletJson = sConf.bulletJsonName;
        this.bulletType = sConf.bulletType;
        var eConf = EffectTrigger_1.EffectTrigger[this.bulletJson];
        if (eConf) {
            this.anim = Object.keys(eConf)[0];
        }
        else {
            this.anim = "";
        }
        if (this.bulletType === BattleConst_1.BattleConst.BULLET_TYPE.NORMAL ||
            this.bulletType === BattleConst_1.BattleConst.BULLET_TYPE.NONE) {
            this._events.push({
                name: "e_zp",
                time: 0,
            });
        }
        else if (this.bulletType === BattleConst_1.BattleConst.BULLET_TYPE.FLY_DIRECT ||
            this.bulletType === BattleConst_1.BattleConst.BULLET_TYPE.FLY_CURVE) {
            this._duration = 0.2;
            this._events.push({
                name: "e_zp",
                time: this._duration,
            });
        }
        else if (this.bulletType === BattleConst_1.BattleConst.BULLET_TYPE.RELEASE ||
            this.bulletType === BattleConst_1.BattleConst.BULLET_TYPE.RELEASE_ONLY) {
            this._duration = eConf[this.anim].duration;
            for (var i = 0; i < eConf[this.anim].events.length; i++) {
                var event_1 = eConf[this.anim].events[i];
                this._events.push(event_1);
            }
        }
        else if (this.bulletType === BattleConst_1.BattleConst.BULLET_TYPE.LASER) {
            this._events.push({
                name: "e_zp",
                time: 0,
            });
        }
        else {
            console.warn("unsupport bulletType: " + this.bulletType, this.bulletJson, this._skill.id);
        }
    }
    SkillBullet.prototype.update = function (dt) {
        this._dt += dt;
        while (this._events.length > 0) {
            var event_2 = this._events[0];
            if (this._dt >= event_2.time) {
                if (event_2.name == "e_zp") {
                    this._skill.executeFuncs();
                }
                this._events.shift();
            }
            else {
                break;
            }
        }
        return (this._duration < this._dt);
    };
    SkillBullet.prototype.isComplete = function () {
        return (this._duration < this._dt);
    };
    SkillBullet.prototype.getOwner = function () {
        return this._owner;
    };
    SkillBullet.prototype.getDuration = function () {
        return this._duration;
    };
    SkillBullet.prototype.getTargets = function () {
        return this._skill.getTargets();
    };
    return SkillBullet;
}());
exports.SkillBullet = SkillBullet;
