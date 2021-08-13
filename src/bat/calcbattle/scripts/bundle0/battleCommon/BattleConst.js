"use strict";
/*
 * @Description: 战斗常量
 * @Autor: qiaomingwu
 * @Date: 2021-03-25 17:24:09
 * @LastEditors: YueHongTing
 * @LastEditTime: 2021-07-21 18:36:03
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.BattleConst = exports.UNIT_STATE = exports.DMG_TYPE = exports.BATTLE_STATE = void 0;
/** 战斗状态 */
var BATTLE_STATE;
(function (BATTLE_STATE) {
    BATTLE_STATE[BATTLE_STATE["NONE"] = 0] = "NONE";
    BATTLE_STATE[BATTLE_STATE["START"] = 1] = "START";
    BATTLE_STATE[BATTLE_STATE["ROUND_START"] = 2] = "ROUND_START";
    BATTLE_STATE[BATTLE_STATE["SORT_ACT"] = 3] = "SORT_ACT";
    BATTLE_STATE[BATTLE_STATE["FIGHT"] = 4] = "FIGHT";
    BATTLE_STATE[BATTLE_STATE["ROUND_END"] = 5] = "ROUND_END";
    BATTLE_STATE[BATTLE_STATE["COMPLETE"] = 6] = "COMPLETE";
})(BATTLE_STATE = exports.BATTLE_STATE || (exports.BATTLE_STATE = {}));
;
/** 伤害类型 */
var DMG_TYPE;
(function (DMG_TYPE) {
    DMG_TYPE[DMG_TYPE["NONE"] = 0] = "NONE";
    DMG_TYPE[DMG_TYPE["PHYIC"] = 1] = "PHYIC";
    DMG_TYPE[DMG_TYPE["CURE"] = 2] = "CURE";
    DMG_TYPE[DMG_TYPE["INHERIT"] = 3] = "INHERIT";
})(DMG_TYPE = exports.DMG_TYPE || (exports.DMG_TYPE = {}));
/** 战斗单位状态 */
var UNIT_STATE;
(function (UNIT_STATE) {
    UNIT_STATE[UNIT_STATE["IDLE"] = 0] = "IDLE";
    UNIT_STATE[UNIT_STATE["SKILL"] = 1] = "SKILL";
    UNIT_STATE[UNIT_STATE["DIE"] = 2] = "DIE";
    UNIT_STATE[UNIT_STATE["REBORN_DIE"] = 3] = "REBORN_DIE";
    UNIT_STATE[UNIT_STATE["WAIT_REBORN"] = 4] = "WAIT_REBORN";
    UNIT_STATE[UNIT_STATE["REBORN"] = 5] = "REBORN";
})(UNIT_STATE = exports.UNIT_STATE || (exports.UNIT_STATE = {}));
var BattleConst = /** @class */ (function () {
    function BattleConst() {
    }
    /** 属性id */
    BattleConst.ATTR_TO_ID = {
        ATK: 10,
        DEF: 12,
        MAXLIFE: 14,
        SPEED: 15,
        ATK_COEF: 1000,
        DEF_COEF: 1200,
        MAXLIFE_COEF: 1400,
        SPEED_COEF: 1500,
        CRIT: 20,
        CRIT_DAM: 21,
        CALM: 22,
        CRIT_DOWN: 23,
        HIT: 24,
        BLOCK: 25,
        ENERTY_RATIO: 26,
        CURE_RATIO: 30,
        HEALING_RATIO: 31,
        DAM_UP: 40,
        DAM_DOWN: 41,
        SKILL_DAM_UP: 42,
        SKILL_DAM_DOWN: 43,
        DOT_DAM_UP: 44,
        DOT_DAM_DOWN: 45,
        TAN_DAM_UP: 46,
        ZHAN_DAM_UP: 47,
        FA_DAM_UP: 48,
        XIA_DAM_UP: 49,
        FU_DAM_UP: 50,
        TAN_DAM_DOWN: 51,
        ZHAN_DAM_DOWN: 52,
        FA_DAM_DOWN: 53,
        XIA_DAM_DOWN: 54,
        FU_DAM_DOWN: 55,
        KONG_RES: 70,
        YUN_RES: 71,
        SHI_RES: 72,
        BING_RES: 73,
        MA_RES: 74,
        NU_RES: 75,
        MO_RES: 76,
        KONG_HIT: 79,
        DOT_RES: 80,
        FIRE_DOT_RES: 81,
        BLOOD_DOT_RES: 82,
        DU_DOT_RES: 83,
        ICE_DOT_RES: 84,
        DOT_HIT: 89, //	症状命中
    };
    /** 英雄元素类型 */
    BattleConst.HERO_ELEM = {
        ALL: "0",
        WATER: "1",
        SHOUND: "2",
        GOLD: "3",
        FIRE: "4",
        LIGHT: "5",
        DARK: "6", //暗
    };
    /** 单位阵营 */
    BattleConst.UNIT_GROUP = {
        ATTACKER: 1,
        DEFENDER: 2, //防守方
    };
    /** 战斗类型 */
    BattleConst.BATTLE_TYPE = {
        AUTO_FIGHT: "AUTO_FIGHT",
        RANDOM_FIGHT: "RANDOM_FIGHT",
        DEBUG: "DEBUG",
        WLEVEL: "WLEVEL",
        AUTO_WLEVEL: "AUTO_WLEVEL",
        TOWER: "TOWER",
        GUILD_BOSS: "GUILD_BOSS",
        AP_ARENA: "AP_ARENA",
        CRUSADE: "CRUSADE",
        GUILD_WAR_DEFEN: "GUILD_WAR_DEFEN",
        GUILD_WAR: "GUILD_WAR",
        RIFT_MINE: "RIFT_MINE",
        RIFT_MONSTER: "RIFT_MONSTER",
        RIFT_BOX: "RIFT_BOX",
        LADDER: "LADDER",
        MONOPOLY_MONSTER: "MONOPOLY_MONSTER",
        MONOPOLY_HERO: "MONOPOLY_HERO",
        WAR_CUP: "WAR_CUP",
        WORLD_BOSS: "WORLD_BOSS",
        MAZE_BATTLE: "MAZE_BATTLE", //迷宫战斗
    };
    /** 技能类型 */
    BattleConst.SKILL_TYPE = {
        NORMAL: "normal",
        SPECIAL: "special",
        PASSIVE: "passive", //被动技能
    };
    /** 子弹类型 */
    BattleConst.BULLET_TYPE = {
        NONE: 0,
        NORMAL: 1,
        FLY_DIRECT: 2,
        LASER: 3,
        RELEASE: 4,
        RELEASE_ONLY: 5,
        FLY_CURVE: 6, // 抛物线飞行
    };
    /** 所有职业 */
    BattleConst.HERO_JOB = {
        TAN: 1,
        ZHAN: 2,
        FA: 3,
        XIA: 4,
        FU: 5, //辅助
    };
    BattleConst.BUFF_FUNC = {
        STEAL: "buff_steal",
        ATK_EXTRA_DMG: "buff_atkExtraDmg",
        VAMPIRE: "buff_vampire",
        ATK_ENERGY: "buff_atkEnergy",
        CRIT_DAM_FIXED: "buff_critDamFixed",
        CRIT: "buff_crit",
        CRIT_NOT: "buff_critNot",
        DMG_REVISE: "buff_dmgRevise",
        ATK_EXTRA_BUFF: "buff_atkExtraBuff",
        DEATH_EXTRA_BUFF: "buff_deathExtraBuff",
        ROUND_START_BUFF: "buff_roundStartBuff",
        ROUND_END_BUFF: "buff_roundEndBuff",
        ACT_EXTRA_BUFF: "buff_actExtraBuff",
        DOT: "buff_dot",
        DOT_EXTRA_BUFF: "buff_dotExtraBuff",
        CURE_EXTRA_BUFF: "buff_cureExtraBuff",
        BUFF_TYPE_EXTRA_BUFF: "buff_buffTypeExtraBuff",
        IMMUNE_BUFF_TYPE: "buff_immuneBuffType",
        DEATH_ENERGY: "buff_deathEnergy",
        CURE_ENERGY: "buff_cureEnergy",
        ACT_ENERGY: "buff_actEnergy",
        ACT_STEAL_ENERGY: "buff_actStealEnergy",
        ATTR_TO_SHIELD: "buff_attrToShield",
        INVINCIBLE: "buff_invincible",
        HP_REDUCE_BUFF: "buff_hpReduceBuff",
        HP_POINT_BUFF: "buff_hpPointBuff",
        REBOUND_BUFF: "buff_reboundBuff",
        REBOUND_DMG: "buff_reboundDmg",
        ATK_DEATH: "buff_atkDeath",
        HIT_DMG: "buff_hitDmg",
        ACT_START_DMG_BUFF: "buff_actStartDmgBuff",
        ROUND_START_DMG_BUFF: "buff_roundStartDmgBuff",
        ROUND_END_DMG_BUFF: "buff_roundEndDmgBuff",
        DEATH_DMG: "buff_deathDmg",
        ATTR_COND: "buff_attrCond",
        DOT_GAMBLE: "buff_dotGamble",
        ROUND_CHANGE_HP: "buff_roundChangeHp",
        FIRST_ACT: "buff_firstAct",
        CONFUSION: "buff_confusion",
        ATK_TARGET_BACK: "buff_atkTargetBack",
        TAUNT: "buff_taunt",
        IMMUNE_DEATH: "buff_immuneDeath",
        BUFF_GRP_BUFF: "buff_buffGrpBuff",
        ACT_NEW_SKILL: "buff_actNewSkill",
        ACT_EXTRA_SPECIAL: "buff_actExtraSpecial",
        HIT_EXTRA_NORMAL: "buff_hitExtraNormal",
        DEATH_EXTRA_SPECIAL: "buff_deathExtraSpecial",
        ACT_EXTRA_ACT: "buff_actExtraAct",
        REBORN: "buff_reborn",
        REBORN_NOT: "buff_rebornNot",
    };
    BattleConst.SKILL_FUNC = {
        SKILL_NEXT: "skill_next",
    };
    BattleConst.BATTLE_BASE_ROUND = 15;
    return BattleConst;
}());
exports.BattleConst = BattleConst;
;
