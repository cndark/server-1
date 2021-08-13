"use strict";
/*
 * @Description: 战斗服teamCtrl
 * @Autor: chenjie
 * @Date: 2021-04-08 21:34:38
 * @LastEditors: zyb
 * @LastEditTime: 2021-07-09 17:57:54
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.SvrTeamCtrl = exports.ROUND_TYPE = void 0;
var Configs_1 = require("../framework/config/Configs");
var ROUND_TYPE;
(function (ROUND_TYPE) {
    ROUND_TYPE["MONSTER"] = "0";
    ROUND_TYPE["BOSS"] = "1";
    ROUND_TYPE["ARENA"] = "2";
})(ROUND_TYPE = exports.ROUND_TYPE || (exports.ROUND_TYPE = {}));
;
var SvrTeamCtrl = /** @class */ (function () {
    function SvrTeamCtrl() {
    }
    // 获取hero战斗属性
    SvrTeamCtrl.getHeroPropsByLvStar = function (id, lv, star) {
        var props = {};
        var mConf = Configs_1.Configs.monsterConf[id];
        var starConf = Configs_1.Configs.heroStarUpConf[star || mConf.star];
        var baseProps = mConf.heroBaseProps;
        var lvAddProps = mConf.heroPropGrowth;
        var starRatio = starConf.propsRatio;
        //基础属性
        baseProps.forEach(function (info) {
            props[info.id] = (props[info.id] || 0) + info.val;
        });
        //升级属性
        lvAddProps.forEach(function (info) {
            var add = info.val * lv;
            props[info.id] = (props[info.id] || 0) + add;
        });
        //升星属性
        starRatio.forEach(function (info) {
            if (props[info.id]) {
                props[info.id] *= (1 + info.val);
            }
        });
        return props;
    };
    // 获取monster战斗属性
    SvrTeamCtrl.getMonsterAttr = function (id, lv) {
        var props = {};
        var monsterConf = Configs_1.Configs.monsterConf[id];
        var monsterPowerConf = Configs_1.Configs.monsterPowerConf[lv];
        var gConf = Configs_1.Configs.globalBattleConf[1];
        for (var i = 0; i < monsterConf.monsterPropsRatio.length; i++) {
            var ratio = monsterConf.monsterPropsRatio[i];
            var prop = { id: 0, val: 0 };
            if (monsterPowerConf.baseProps[i]) {
                prop.id = monsterPowerConf.baseProps[i].id;
                prop.val = monsterPowerConf.baseProps[i].val;
            }
            if (gConf.monsterBaseProps[i]) {
                if (!prop.id) {
                    prop.id = gConf.monsterBaseProps[i].id;
                    prop.val = gConf.monsterBaseProps[i].val;
                }
            }
            else {
                prop.val += gConf.monsterBaseProps[i].val;
            }
            if (!prop.id)
                break;
            props[prop.id] = prop.val * ratio;
        }
        return props;
    };
    // 从monster表中获取英雄信息
    SvrTeamCtrl.getFighterByConf = function (monsterInfo) {
        var fighters = [];
        for (var i = 0; i < monsterInfo.length; i++) {
            var info = monsterInfo[i];
            if (info.id <= 0)
                continue;
            var mpConf = Configs_1.Configs.monsterPowerConf[info.lv];
            var star = info.star || mpConf.star;
            fighters.push({
                Id: info.id,
                Lv: info.lv,
                Star: star,
                Props: info.id >= 31000 ? this.getMonsterAttr(info.id, info.lv) : this.getHeroPropsByLvStar(info.id, info.lv, star),
                Pos: i, // start from 0
            });
        }
        return fighters;
    };
    return SvrTeamCtrl;
}());
exports.SvrTeamCtrl = SvrTeamCtrl;
