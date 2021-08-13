"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-05-10 16:37:28
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-05-19 09:50:21
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.RecordCtrl = void 0;
var BattleConst_1 = require("../../bundle0/battleCommon/BattleConst");
var RecordCtrl = /** @class */ (function () {
    function RecordCtrl(bCtrl) {
        this.bCtrl = null;
        this.content = [];
        this.bCtrl = bCtrl;
        this.content = [];
    }
    RecordCtrl.prototype.recordKeyWord = function (state) {
        var s = "[" + state + "]";
        this.content.push(s);
    };
    RecordCtrl.prototype.unitReborn = function (unit) {
        var s = "[Reborn]" + unit.order + "_" + unit.group + "_" + unit.id +
            " :hp=(" + unit.hp + ","
            + unit.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE] + ")";
        this.content.push(s);
    };
    RecordCtrl.prototype.unitSkillStart = function (skill, seqs) {
        var owner = skill.getOwner();
        var s = "[SkillStart]" + skill.id +
            " owner:" + owner.order + "_" + owner.group + "_" + owner.id;
        var extra = " targets:[";
        var units = this.bCtrl.getUnits();
        for (var i = 0; i < seqs.length; i++) {
            var seq = seqs[i];
            var exist = false;
            for (var j = 0; j < units.length; j++) {
                var unit = units[j];
                if (Number(seq) == unit.seq) {
                    extra += unit.order + "_" + unit.group + "_" + unit.id + ", ";
                    exist = true;
                    break;
                }
            }
            if (!exist) {
                extra += "cant find seq: " + seq + ", ";
            }
        }
        extra += "]";
        s += extra;
        this.content.push(s);
    };
    RecordCtrl.prototype.unitSkillComplete = function (skill) {
        var owner = skill.getOwner();
        var s = "[SkillComplete]" + skill.id +
            " owner:" + owner.order + "_" + owner.group + "_" + owner.id;
        this.content.push(s);
    };
    RecordCtrl.prototype.switchState = function (unit, preState, newState) {
        var s = "[SwitchState]" + unit.order + "_" + unit.group + "_" + unit.id +
            " , " + preState + "=>" + newState;
        this.content.push(s);
    };
    RecordCtrl.prototype.updateHp = function (unit, totalDmg, skill, buff, crit) {
        var reason = "";
        if (skill) {
            var owner = skill.getOwner();
            reason = "skill:" + skill.id + " atk:" + owner.order + "_" + owner.group + "_" + owner.id;
        }
        else if (buff) {
            var owner = buff.getOwner();
            reason = "buff:" + buff.id + " atk:" + owner.order + "_" + owner.group + "_" + owner.id;
        }
        else {
            reason = "forceUpdateHp";
        }
        var s = "[Hp]" + unit.order + "_" + unit.group
            + ":hp=(" + unit.hp + ","
            + unit.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE] + ")"
            + " reason:" + reason;
        this.content.push(s);
    };
    RecordCtrl.prototype.addBuff = function (buff) {
        var atker = buff.getAtker();
        var owner = buff.getOwner();
        var s = "[AddBuff]" + buff.id +
            " owner:" + owner.order + "_" + owner.group + "_" + owner.id +
            " atker:" + atker.order + "_" + atker.group + "_" + atker.id;
        this.content.push(s);
    };
    RecordCtrl.prototype.removeBuff = function (buff) {
        var owner = buff.getOwner();
        var s = "[RemoveBuff]" + buff.id +
            " owner:" + owner.order + "_" + owner.group + "_" + owner.id;
        this.content.push(s);
    };
    return RecordCtrl;
}());
exports.RecordCtrl = RecordCtrl;
