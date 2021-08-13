"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-15 16:51:49
 * @LastEditors: chenjie
 * @LastEditTime: 2021-04-27 17:58:11
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.BFFactory = void 0;
var BattleConst_1 = require("../../../bundle0/battleCommon/BattleConst");
var BFActEnergy_1 = require("./BFActEnergy");
var BFActExtraAct_1 = require("./BFActExtraAct");
var BFActExtraBuff_1 = require("./BFActExtraBuff");
var BFActExtraSpecial_1 = require("./BFActExtraSpecial");
var BFActNewSkill_1 = require("./BFActNewSkill");
var BFActStartDmgBuff_1 = require("./BFActStartDmgBuff");
var BFActStealEnergy_1 = require("./BFActStealEnergy");
var BFAtkDeath_1 = require("./BFAtkDeath");
var BFAtkEnergy_1 = require("./BFAtkEnergy");
var BFAtkExtraBuff_1 = require("./BFAtkExtraBuff");
var BFAtkExtraDmg_1 = require("./BFAtkExtraDmg");
var BFAtkTargetBack_1 = require("./BFAtkTargetBack");
var BFAttrCond_1 = require("./BFAttrCond");
var BFAttrToShield_1 = require("./BFAttrToShield");
var BFBuffGrpBuff_1 = require("./BFBuffGrpBuff");
var BFBuffTypeExtraBuff_1 = require("./BFBuffTypeExtraBuff");
var BFConfusion_1 = require("./BFConfusion");
var BFCrit_1 = require("./BFCrit");
var BFCritDamFixed_1 = require("./BFCritDamFixed");
var BFCritNot_1 = require("./BFCritNot");
var BFCureEnergy_1 = require("./BFCureEnergy");
var BFCureExtraBuff_1 = require("./BFCureExtraBuff");
var BFDeathDmg_1 = require("./BFDeathDmg");
var BFDeathEnergy_1 = require("./BFDeathEnergy");
var BFDeathExtraBuff_1 = require("./BFDeathExtraBuff");
var BFDeathExtraSpecial_1 = require("./BFDeathExtraSpecial");
var BFDmgRevise_1 = require("./BFDmgRevise");
var BFDot_1 = require("./BFDot");
var BFDotExtraBuff_1 = require("./BFDotExtraBuff");
var BFDotGamble_1 = require("./BFDotGamble");
var BFFirstAct_1 = require("./BFFirstAct");
var BFHitDmg_1 = require("./BFHitDmg");
var BFHitExtraNormal_1 = require("./BFHitExtraNormal");
var BFHpPointBuff_1 = require("./BFHpPointBuff");
var BFHpReduceBuff_1 = require("./BFHpReduceBuff");
var BFImmuneBuffType_1 = require("./BFImmuneBuffType");
var BFImmuneDeath_1 = require("./BFImmuneDeath");
var BFInvincible_1 = require("./BFInvincible");
var BFReborn_1 = require("./BFReborn");
var BFReboundBuff_1 = require("./BFReboundBuff");
var BFReboundDmg_1 = require("./BFReboundDmg");
var BFRoundChangeHp_1 = require("./BFRoundChangeHp");
var BFRoundEndBuff_1 = require("./BFRoundEndBuff");
var BFRoundEndDmgBuff_1 = require("./BFRoundEndDmgBuff");
var BFRoundStartBuff_1 = require("./BFRoundStartBuff");
var BFRoundStartDmgBuff_1 = require("./BFRoundStartDmgBuff");
var BFSteal_1 = require("./BFSteal");
var BFTaunt_1 = require("./BFTaunt");
var BFVampire_1 = require("./BFVampire");
var BuffFunc_1 = require("./BuffFunc");
var BFFactory = /** @class */ (function () {
    function BFFactory() {
    }
    BFFactory.createBuffFunc = function (buff, func) {
        var info = func.split("~");
        var funcName = info[0];
        var params = info.splice(1);
        var bfunc = null;
        if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.STEAL) {
            bfunc = new BFSteal_1.BFSteal(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ATK_EXTRA_DMG) {
            bfunc = new BFAtkExtraDmg_1.BFAtkExtraDmg(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.VAMPIRE) {
            bfunc = new BFVampire_1.BFVampire(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ATK_ENERGY) {
            bfunc = new BFAtkEnergy_1.BFAtkEnergy(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.CRIT_DAM_FIXED) {
            bfunc = new BFCritDamFixed_1.BFCritDamFixed(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.CRIT) {
            bfunc = new BFCrit_1.BFCrit(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.CRIT_NOT) {
            bfunc = new BFCritNot_1.BFCritNot(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.DMG_REVISE) {
            bfunc = new BFDmgRevise_1.BFDmgRevise(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ATK_EXTRA_BUFF) {
            bfunc = new BFAtkExtraBuff_1.BFAtkExtraBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.DEATH_EXTRA_BUFF) {
            bfunc = new BFDeathExtraBuff_1.BFDeathExtraBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ROUND_START_BUFF) {
            bfunc = new BFRoundStartBuff_1.BFRoundStartBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ROUND_END_BUFF) {
            bfunc = new BFRoundEndBuff_1.BFRoundEndBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ACT_EXTRA_BUFF) {
            bfunc = new BFActExtraBuff_1.BFActExtraBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.DOT) {
            bfunc = new BFDot_1.BFDot(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.DOT_EXTRA_BUFF) {
            bfunc = new BFDotExtraBuff_1.BFDotExtraBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.CURE_EXTRA_BUFF) {
            bfunc = new BFCureExtraBuff_1.BFCureExtraBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.BUFF_TYPE_EXTRA_BUFF) {
            bfunc = new BFBuffTypeExtraBuff_1.BFBuffTypeExtraBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.IMMUNE_BUFF_TYPE) {
            bfunc = new BFImmuneBuffType_1.BFImmuneBuffType(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.DEATH_ENERGY) {
            bfunc = new BFDeathEnergy_1.BFDeathEnergy(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.CURE_ENERGY) {
            bfunc = new BFCureEnergy_1.BFCureEnergy(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ACT_ENERGY) {
            bfunc = new BFActEnergy_1.BFActEnergy(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ACT_STEAL_ENERGY) {
            bfunc = new BFActStealEnergy_1.BFActStealEnergy(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ATTR_TO_SHIELD) {
            bfunc = new BFAttrToShield_1.BFAttrToShield(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.INVINCIBLE) {
            bfunc = new BFInvincible_1.BFInvincible(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.REBOUND_BUFF) {
            bfunc = new BFReboundBuff_1.BFReboundBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.REBOUND_DMG) {
            bfunc = new BFReboundDmg_1.BFReboundDmg(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ATK_DEATH) {
            bfunc = new BFAtkDeath_1.BFAtkDeath(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.HIT_DMG) {
            bfunc = new BFHitDmg_1.BFHitDmg(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ACT_START_DMG_BUFF) {
            bfunc = new BFActStartDmgBuff_1.BFActStartDmgBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ROUND_START_DMG_BUFF) {
            bfunc = new BFRoundStartDmgBuff_1.BFRoundStartDmgBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ROUND_END_DMG_BUFF) {
            bfunc = new BFRoundEndDmgBuff_1.BFRoundEndDmgBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.DEATH_DMG) {
            bfunc = new BFDeathDmg_1.BFDeathDmg(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.HP_REDUCE_BUFF) {
            bfunc = new BFHpReduceBuff_1.BFHpReduceBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.HP_POINT_BUFF) {
            bfunc = new BFHpPointBuff_1.BFHpPointBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ATTR_COND) {
            bfunc = new BFAttrCond_1.BFAttrCond(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.DOT_GAMBLE) {
            bfunc = new BFDotGamble_1.BFDotGamble(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ROUND_CHANGE_HP) {
            bfunc = new BFRoundChangeHp_1.BFRoundChangeHp(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.FIRST_ACT) {
            bfunc = new BFFirstAct_1.BFFirstAct(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.CONFUSION) {
            bfunc = new BFConfusion_1.BFConfusion(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ATK_TARGET_BACK) {
            bfunc = new BFAtkTargetBack_1.BFAtkTargetBack(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.TAUNT) {
            bfunc = new BFTaunt_1.BFTaunt(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.IMMUNE_DEATH) {
            bfunc = new BFImmuneDeath_1.BFImmuneDeath(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.BUFF_GRP_BUFF) {
            bfunc = new BFBuffGrpBuff_1.BFBuffGrpBuff(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ACT_NEW_SKILL) {
            bfunc = new BFActNewSkill_1.BFActNewSkill(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ACT_EXTRA_SPECIAL) {
            bfunc = new BFActExtraSpecial_1.BFActExtraSpecial(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.HIT_EXTRA_NORMAL) {
            bfunc = new BFHitExtraNormal_1.BFHitExtraNormal(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.DEATH_EXTRA_SPECIAL) {
            bfunc = new BFDeathExtraSpecial_1.BFDeathExtraSpecial(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.ACT_EXTRA_ACT) {
            bfunc = new BFActExtraAct_1.BFActExtraAct(buff, funcName, params);
        }
        else if (funcName === BattleConst_1.BattleConst.BUFF_FUNC.REBORN) {
            bfunc = new BFReborn_1.BFReborn(buff, funcName, params);
        }
        else {
            bfunc = new BuffFunc_1.BuffFunc(buff, funcName, params);
        }
        return bfunc;
    };
    return BFFactory;
}());
exports.BFFactory = BFFactory;
