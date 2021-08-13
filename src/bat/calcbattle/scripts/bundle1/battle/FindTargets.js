"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-07 12:11:10
 * @LastEditors: chenjie
 * @LastEditTime: 2021-05-19 14:59:33
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.FindTargets = void 0;
var BattleConst_1 = require("../../bundle0/battleCommon/BattleConst");
var FindTargets = /** @class */ (function () {
    function FindTargets() {
    }
    FindTargets.getTargets = function (owner, isFriendly, targetType, allTargets, errStr) {
        if (errStr === void 0) { errStr = ""; }
        var info = targetType.split("|");
        if (info[0] === "mixAnd") {
            return this.getMixAndTargets(info, owner, isFriendly, allTargets, errStr);
        }
        else if (info[0] === "mixOr") {
            return this.getMixOrTargets(info, owner, isFriendly, allTargets, errStr);
        }
        else if (info[0] === "mixFirst") {
            return this.getMixFirstTargets(info, owner, isFriendly, allTargets, errStr);
        }
        else {
            return this.getSingleTargets(targetType, owner, isFriendly, allTargets, errStr);
        }
    };
    FindTargets.getMixAndTargets = function (info, owner, isFriendly, allTargets, errStr) {
        if (errStr === void 0) { errStr = ""; }
        var ret = [];
        for (var i = 1; i < info.length; i++) {
            var targetType = info[i];
            var targets = FindTargets.getSingleTargets(targetType, owner, isFriendly, i == 1 ? allTargets : ret, errStr);
            ret = targets;
        }
        return ret;
    };
    FindTargets.getMixOrTargets = function (info, owner, isFriendly, allTargets, errStr) {
        if (errStr === void 0) { errStr = ""; }
        var seqMap = {};
        var ret = [];
        for (var i = 1; i < info.length; i++) {
            var targetType = info[i];
            var targets = FindTargets.getSingleTargets(targetType, owner, isFriendly, allTargets, errStr);
            if (i == 1) {
                for (var j = 0; j < targets.length; j++) {
                    var target = targets[j];
                    ret.push(target);
                    seqMap[target.seq] = true;
                }
            }
            else {
                for (var j = 0; j < targets.length; j++) {
                    var target = targets[j];
                    if (!seqMap[target.seq]) {
                        seqMap[target.seq] = true;
                        ret.push(target);
                    }
                }
            }
        }
        return ret;
    };
    FindTargets.getMixFirstTargets = function (info, owner, isFriendly, allTargets, errStr) {
        if (errStr === void 0) { errStr = ""; }
        for (var i = 1; i < info.length; i++) {
            var targetType = info[i];
            var targets = FindTargets.getSingleTargets(targetType, owner, isFriendly, allTargets, errStr);
            if (targets.length > 0)
                return targets;
        }
        return [];
    };
    FindTargets.getSingleTargets = function (targetType, owner, isFriendly, allTargets, errStr) {
        if (errStr === void 0) { errStr = ""; }
        var ret = [];
        var info = targetType.split("_");
        var funcName = "find_" + info[0];
        var obj = FindTargets;
        if (obj[funcName] && typeof obj[funcName] == 'function') {
            ret = obj[funcName](owner, isFriendly, info, allTargets);
        }
        else {
            console.error("Find Targets invalid TargetType: " + targetType, errStr);
        }
        return ret;
    };
    FindTargets.getCommon = function (owner, isFriendly, allTargets) {
        var ret = [];
        var group = owner.group;
        if (isFriendly === 0) {
            group++;
            if (group > BattleConst_1.BattleConst.UNIT_GROUP.DEFENDER) {
                group -= BattleConst_1.BattleConst.UNIT_GROUP.DEFENDER;
            }
        }
        for (var i = 0; i < allTargets.length; i++) {
            var unit = allTargets[i];
            if (unit.isAlive() && unit.group == group) {
                ret.push(unit);
            }
        }
        return ret;
    };
    FindTargets.getRandom = function (ret, cnt, owner) {
        if (cnt < ret.length) {
            var tmp = [];
            while (tmp.length < cnt) {
                var idx = Math.floor(owner.bCtrl.random() * ret.length);
                tmp.push(ret[idx]);
                ret.splice(idx, 1);
            }
            return tmp;
        }
        else {
            return ret;
        }
    };
    // 最小最大站位
    FindTargets.find_order = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        if (info[2] === "min") {
            ret.sort(function (a, b) { return a.order - b.order; });
        }
        else {
            ret.sort(function (a, b) { return b.order - a.order; });
        }
        while (ret.length > cnt) {
            ret.pop();
        }
        return ret;
    };
    // 随机个数
    FindTargets.find_random = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        ret = FindTargets.getRandom(ret, cnt, owner);
        return ret;
    };
    // 自己
    FindTargets.find_self = function (owner, isFriendly, info, allTargets) {
        var ret = [owner];
        return ret;
    };
    //队友
    FindTargets.find_teammate = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var tmp = [];
        for (var j = 0; j < ret.length; j++) {
            var unit = ret[j];
            if (unit.seq !== owner.seq && unit.group === owner.group) {
                tmp.push(unit);
            }
        }
        return FindTargets.getRandom(tmp, cnt, owner);
    };
    // 前排
    FindTargets.find_front = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var findBack = info[2] === "1";
        var tmp = [];
        for (var j = 0; j < ret.length; j++) {
            var unit = ret[j];
            if (unit.order <= 1) {
                tmp.push(unit);
            }
        }
        if (tmp.length <= 0) {
            if (findBack)
                return FindTargets.getRandom(ret, cnt, owner);
            else
                return [];
        }
        else {
            return FindTargets.getRandom(tmp, cnt, owner);
        }
    };
    // 后排
    FindTargets.find_back = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var findFront = info[2] === "1";
        var tmp = [];
        for (var j = 0; j < ret.length; j++) {
            var unit = ret[j];
            if (unit.order > 1) {
                tmp.push(unit);
            }
        }
        if (tmp.length <= 0) {
            if (findFront)
                return FindTargets.getRandom(ret, cnt, owner);
            else
                return [];
        }
        else {
            return FindTargets.getRandom(tmp, cnt, owner);
        }
    };
    // 阵营
    FindTargets.find_elem = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var elemMap = {};
        for (var i = 2; i < info.length; i++) {
            elemMap[Number(info[i])] = true;
        }
        var tmp = [];
        for (var j = 0; j < ret.length; j++) {
            var unit = ret[j];
            if (elemMap[unit.elem]) {
                tmp.push(unit);
            }
        }
        return FindTargets.getRandom(tmp, cnt, owner);
    };
    // 职业
    FindTargets.find_job = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var jobMap = {};
        for (var i = 2; i < info.length; i++) {
            jobMap[Number(info[i])] = true;
        }
        var tmp = [];
        for (var j = 0; j < ret.length; j++) {
            var unit = ret[j];
            if (jobMap[unit.job]) {
                tmp.push(unit);
            }
        }
        return FindTargets.getRandom(tmp, cnt, owner);
    };
    // 血量百分比最低最高
    FindTargets.find_hpRatio = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        if (info[2] === "min") {
            ret.sort(function (a, b) { return a.hp / a.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE] - b.hp / b.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE]; });
        }
        else {
            ret.sort(function (a, b) { return b.hp / b.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE] - a.hp / a.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE]; });
        }
        while (ret.length > cnt) {
            ret.pop();
        }
        return ret;
    };
    // 血量绝对值最低最高
    FindTargets.find_hp = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        if (info[2] === "min") {
            ret.sort(function (a, b) { return a.hp - b.hp; });
        }
        else {
            ret.sort(function (a, b) { return b.hp - a.hp; });
        }
        while (ret.length > cnt) {
            ret.pop();
        }
        return ret;
    };
    // 血量百分比高于低于指定百分比
    FindTargets.find_hpPoint = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var lt = (info[2] === "lt");
        var ratio = Number(info[3]);
        var tmp = [];
        for (var j = 0; j < ret.length; j++) {
            var unit = ret[j];
            var hpRatio = unit.hp / unit.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
            if ((lt && hpRatio < ratio) || (!lt && hpRatio > ratio)) {
                tmp.push(unit);
            }
        }
        return FindTargets.getRandom(tmp, cnt, owner);
    };
    // 血量百分比高于低于自身
    FindTargets.find_hpSelf = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var lt = (info[2] === "lt");
        var ratio = owner.hp / owner.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
        var tmp = [];
        for (var j = 0; j < ret.length; j++) {
            var unit = ret[j];
            var hpRatio = unit.hp / unit.attrs[BattleConst_1.BattleConst.ATTR_TO_ID.MAXLIFE];
            if ((lt && hpRatio < ratio) || (!lt && hpRatio > ratio)) {
                tmp.push(unit);
            }
        }
        return FindTargets.getRandom(tmp, cnt, owner);
    };
    // 指定属性最低最高
    FindTargets.find_attr = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var attrId = Number(info[3]);
        if (info[2] === "min") {
            ret.sort(function (a, b) { return a.attrs[attrId] - b.attrs[attrId]; });
        }
        else {
            ret.sort(function (a, b) { return b.attrs[attrId] - a.attrs[attrId]; });
        }
        while (ret.length > cnt) {
            ret.pop();
        }
        return ret;
    };
    // 指定属性高于低于自身指定百分比
    FindTargets.find_attrSelf = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var lt = (info[2] === "lt");
        var attrId = Number(info[3]);
        var ratio = Number(info[4]);
        var value = owner.attrs[attrId];
        var tmp = [];
        for (var j = 0; j < ret.length; j++) {
            var unit = ret[j];
            var attrVal = unit.attrs[attrId];
            if (lt && attrVal < value && (value - attrVal) / value > ratio) {
                tmp.push(unit);
            }
            else if (!lt && attrVal > value && (value <= 0 || (attrVal - value) / value > ratio)) {
                tmp.push(unit);
            }
        }
        return FindTargets.getRandom(tmp, cnt, owner);
    };
    // 指定近战远程
    FindTargets.find_distance = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var tmp = [];
        var dis = Number(info[2]);
        for (var j = 0; j < ret.length; j++) {
            var unit = ret[j];
            if (unit.distance === dis) {
                tmp.push(unit);
            }
        }
        return FindTargets.getRandom(tmp, cnt, owner);
    };
    // 能量最高最低
    FindTargets.find_mp = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        if (info[2] === "min") {
            ret.sort(function (a, b) { return a.mp - b.mp; });
        }
        else {
            ret.sort(function (a, b) { return b.mp - a.mp; });
        }
        while (ret.length > cnt) {
            ret.pop();
        }
        return ret;
    };
    // 满能量
    FindTargets.find_mpFull = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var tmp = [];
        for (var j = 0; j < ret.length; j++) {
            var unit = ret[j];
            if (unit.mp >= unit.maxMp) {
                tmp.push(unit);
            }
        }
        return FindTargets.getRandom(tmp, cnt, owner);
    };
    // 指定英雄
    FindTargets.find_hero = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var heroId = Number(info[2]);
        var tmp = [];
        for (var j = 0; j < ret.length; j++) {
            var unit = ret[j];
            if (unit.id === heroId) {
                tmp.push(unit);
            }
        }
        return FindTargets.getRandom(tmp, cnt, owner);
    };
    // 已死亡
    FindTargets.find_died = function (owner, isFriendly, info, allTargets) {
        var ret = [];
        var cnt = Number(info[1]);
        var group = owner.group;
        if (isFriendly === 0) {
            group++;
            if (group > BattleConst_1.BattleConst.UNIT_GROUP.DEFENDER) {
                group -= BattleConst_1.BattleConst.UNIT_GROUP.DEFENDER;
            }
        }
        for (var i = 0; i < allTargets.length; i++) {
            var unit = allTargets[i];
            if (!unit.isAlive() && unit.group == group) {
                ret.push(unit);
            }
        }
        return FindTargets.getRandom(ret, cnt, owner);
    };
    // 身上有指定buffType的目标
    FindTargets.find_buffType = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var bTypes = [Number(info[2])];
        var tmp = [];
        for (var i = 0; i < ret.length; i++) {
            var unit = ret[i];
            if (unit.hasBuffType(bTypes))
                tmp.push(unit);
        }
        return FindTargets.getRandom(tmp, cnt, owner);
    };
    //身上有指定buffGrp的目标
    FindTargets.find_buffGrp = function (owner, isFriendly, info, allTargets) {
        var ret = FindTargets.getCommon(owner, isFriendly, allTargets);
        var cnt = Number(info[1]);
        var grpId = Number(info[2]);
        var tmp = [];
        for (var i = 0; i < ret.length; i++) {
            var unit = ret[i];
            var num = unit.buffGrpCnt[grpId];
            if (num && num > 0)
                tmp.push(unit);
        }
        return FindTargets.getRandom(tmp, cnt, owner);
    };
    return FindTargets;
}());
exports.FindTargets = FindTargets;
