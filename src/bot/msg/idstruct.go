package msg

var MsgCreators = map[uint32]func() Message{
    1000: func() Message {
        return &C_Auth{}
    },
    1001: func() Message {
        return &GW_Auth_R{}
    },
    1002: func() Message {
        return &C_Login{}
    },
    1003: func() Message {
        return &GW_Login_R{}
    },
    1004: func() Message {
        return &C_TokenGet{}
    },
    1005: func() Message {
        return &GW_TokenGet_R{}
    },
    1006: func() Message {
        return &C_TokenAuth{}
    },
    1007: func() Message {
        return &GW_TokenAuth_R{}
    },
    5000: func() Message {
        return &GS_LoginError{}
    },
    5001: func() Message {
        return &GS_UserInfo{}
    },
    5002: func() Message {
        return &C_TimeSync{}
    },
    5003: func() Message {
        return &GS_TimeSync_R{}
    },
    5004: func() Message {
        return &GS_GameDataReloaded{}
    },
    5006: func() Message {
        return &C_UserExtInfo{}
    },
    5007: func() Message {
        return &GS_UserExtInfo_R{}
    },
    5100: func() Message {
        return &GS_PlayerUpdateLv{}
    },
    5101: func() Message {
        return &GS_PlayerUpdateAtkPwr{}
    },
    5102: func() Message {
        return &GS_PlayerUpdateHFrame{}
    },
    5132: func() Message {
        return &C_PlayerChangeName{}
    },
    5133: func() Message {
        return &GS_PlayerChangeName_R{}
    },
    5134: func() Message {
        return &C_PlayerInfo{}
    },
    5135: func() Message {
        return &GS_PlayerInfo_R{}
    },
    5138: func() Message {
        return &C_PlayerHFrameAdd{}
    },
    5139: func() Message {
        return &GS_PlayerHFrameAdd_R{}
    },
    5140: func() Message {
        return &C_PlayerHFrameSet{}
    },
    5141: func() Message {
        return &GS_PlayerHFrameSet_R{}
    },
    5142: func() Message {
        return &C_PlayerHeadSet{}
    },
    5143: func() Message {
        return &GS_PlayerHeadSet_R{}
    },
    5200: func() Message {
        return &GS_BagUpdate{}
    },
    5300: func() Message {
        return &GS_ArmorUpdate_HeroSeq{}
    },
    5310: func() Message {
        return &C_ItemExchange{}
    },
    5311: func() Message {
        return &GS_ItemExchange_R{}
    },
    5320: func() Message {
        return &C_ItemUse{}
    },
    5321: func() Message {
        return &GS_ItemUse_R{}
    },
    5322: func() Message {
        return &C_ArmorEquip{}
    },
    5323: func() Message {
        return &GS_ArmorEquip_R{}
    },
    5324: func() Message {
        return &C_ArmorUnequip{}
    },
    5325: func() Message {
        return &GS_ArmorUnequip_R{}
    },
    5326: func() Message {
        return &C_ArmorEquipOnekey{}
    },
    5327: func() Message {
        return &GS_ArmorEquipOnekey_R{}
    },
    5328: func() Message {
        return &C_ArmorUnequipOnekey{}
    },
    5329: func() Message {
        return &GS_ArmorUnequipOnekey_R{}
    },
    5330: func() Message {
        return &C_ArmorCompose{}
    },
    5331: func() Message {
        return &GS_ArmorCompose_R{}
    },
    5332: func() Message {
        return &C_ArmorComposeOnekey{}
    },
    5333: func() Message {
        return &GS_ArmorComposeOnekey_R{}
    },
    5334: func() Message {
        return &C_ItemChoose{}
    },
    5335: func() Message {
        return &GS_ItemChoose_R{}
    },
    5400: func() Message {
        return &GS_HeroUpdate{}
    },
    5410: func() Message {
        return &C_HeroLevelUp{}
    },
    5411: func() Message {
        return &GS_HeroLevelUp_R{}
    },
    5414: func() Message {
        return &C_HeroStarUp{}
    },
    5415: func() Message {
        return &GS_HeroStarUp_R{}
    },
    5416: func() Message {
        return &C_HeroLock{}
    },
    5417: func() Message {
        return &GS_HeroLock_R{}
    },
    5418: func() Message {
        return &C_HeroReset{}
    },
    5419: func() Message {
        return &GS_HeroReset_R{}
    },
    5420: func() Message {
        return &C_HeroDecompose{}
    },
    5421: func() Message {
        return &GS_HeroDecompose_R{}
    },
    5422: func() Message {
        return &C_HeroChange{}
    },
    5423: func() Message {
        return &GS_HeroChange_R{}
    },
    5424: func() Message {
        return &C_HeroChangeCancel{}
    },
    5425: func() Message {
        return &GS_HeroChangeCancel_R{}
    },
    5426: func() Message {
        return &C_HeroChangeApply{}
    },
    5427: func() Message {
        return &GS_HeroChangeApply_R{}
    },
    5428: func() Message {
        return &C_HeroTrinketUnlock{}
    },
    5429: func() Message {
        return &GS_HeroTrinketUnlock_R{}
    },
    5430: func() Message {
        return &C_HeroTrinketUp{}
    },
    5431: func() Message {
        return &GS_HeroTrinketUp_R{}
    },
    5432: func() Message {
        return &C_HeroTrinketTransformGen{}
    },
    5433: func() Message {
        return &GS_HeroTrinketTransformGen_R{}
    },
    5434: func() Message {
        return &C_HeroTrinketTransformCommit{}
    },
    5435: func() Message {
        return &GS_HeroTrinketTransformCommit_R{}
    },
    5436: func() Message {
        return &C_HeroBagBuy{}
    },
    5437: func() Message {
        return &GS_HeroBagBuy_R{}
    },
    5438: func() Message {
        return &C_HeroInherit{}
    },
    5439: func() Message {
        return &GS_HeroInherit_R{}
    },
    5440: func() Message {
        return &C_HeroSetSkin{}
    },
    5441: func() Message {
        return &GS_HeroSetSkin_R{}
    },
    5500: func() Message {
        return &GS_RelicUpdate_HeroSeq{}
    },
    5501: func() Message {
        return &GS_RelicUpdate_Star{}
    },
    5510: func() Message {
        return &C_RelicEquip{}
    },
    5511: func() Message {
        return &GS_RelicEquip_R{}
    },
    5512: func() Message {
        return &C_RelicUnequip{}
    },
    5513: func() Message {
        return &GS_RelicUnequip_R{}
    },
    5514: func() Message {
        return &C_RelicEat{}
    },
    5515: func() Message {
        return &GS_RelicEat_R{}
    },
    7000: func() Message {
        return &GS_MailNew{}
    },
    7001: func() Message {
        return &GS_MailDel{}
    },
    7010: func() Message {
        return &C_MailRead{}
    },
    7011: func() Message {
        return &GS_MailRead_R{}
    },
    7012: func() Message {
        return &C_MailDel{}
    },
    7013: func() Message {
        return &GS_MailDel_R{}
    },
    7014: func() Message {
        return &C_MailTakeAttachment{}
    },
    7015: func() Message {
        return &GS_MailTakeAttachment_R{}
    },
    7016: func() Message {
        return &C_MailTakeAttachmentAll{}
    },
    7017: func() Message {
        return &GS_MailTakeAttachmentAll_R{}
    },
    7018: func() Message {
        return &C_MailDelOnekey{}
    },
    7019: func() Message {
        return &GS_MailDelOnekey_R{}
    },
    7100: func() Message {
        return &C_CloudGet{}
    },
    7101: func() Message {
        return &GS_CloudGet_R{}
    },
    7102: func() Message {
        return &C_CloudSet{}
    },
    7103: func() Message {
        return &GS_CloudSet_R{}
    },
    7200: func() Message {
        return &C_TutorialSet{}
    },
    7201: func() Message {
        return &GS_TutorialSet_R{}
    },
    7300: func() Message {
        return &GS_ActStateChange{}
    },
    7320: func() Message {
        return &C_ActStateGet{}
    },
    7321: func() Message {
        return &GS_ActStateGet_R{}
    },
    7400: func() Message {
        return &GS_AttainTabObjValueChanged{}
    },
    7500: func() Message {
        return &GS_BillDone{}
    },
    7501: func() Message {
        return &GS_BillOrder{}
    },
    7510: func() Message {
        return &C_BillInfo{}
    },
    7511: func() Message {
        return &GS_BillInfo_R{}
    },
    7512: func() Message {
        return &C_BillRefundCodeGet{}
    },
    7513: func() Message {
        return &GS_BillRefundCodeGet_R{}
    },
    7514: func() Message {
        return &C_BillRefundTake{}
    },
    7515: func() Message {
        return &GS_BillRefundTake_R{}
    },
    7600: func() Message {
        return &GS_MiscGldActGift{}
    },
    7620: func() Message {
        return &C_MiscBillLocal{}
    },
    7621: func() Message {
        return &GS_MiscBillLocal_R{}
    },
    7622: func() Message {
        return &C_GiftExchange{}
    },
    7623: func() Message {
        return &GS_GiftExchange_R{}
    },
    7624: func() Message {
        return &C_MiscSkipTutorial{}
    },
    7625: func() Message {
        return &GS_MiscSkipTutorial_R{}
    },
    7626: func() Message {
        return &C_MiscGoldenHand{}
    },
    7627: func() Message {
        return &GS_MiscGoldenHand_R{}
    },
    7628: func() Message {
        return &C_MiscOnlineBoxTake{}
    },
    7629: func() Message {
        return &GS_MiscOnlineBoxTake_R{}
    },
    7630: func() Message {
        return &C_MiscSharedGame{}
    },
    7631: func() Message {
        return &GS_MiscSharedGame_R{}
    },
    7632: func() Message {
        return &C_MiscGldActGiftTake{}
    },
    7633: func() Message {
        return &GS_MiscGldActGiftTake_R{}
    },
    7700: func() Message {
        return &GS_GuildPlrLeaveTs{}
    },
    7701: func() Message {
        return &GS_Guild_Join{}
    },
    7702: func() Message {
        return &GS_Guild_Leave{}
    },
    7703: func() Message {
        return &GS_Guild_MbRank{}
    },
    7704: func() Message {
        return &GS_Guild_Lv{}
    },
    7705: func() Message {
        return &GS_Guild_Notice{}
    },
    7706: func() Message {
        return &GS_Guild_Icon{}
    },
    7707: func() Message {
        return &GS_Guild_NewApply{}
    },
    7730: func() Message {
        return &C_GuildCreate{}
    },
    7731: func() Message {
        return &GS_GuildCreate_R{}
    },
    7732: func() Message {
        return &C_GuildDestroy{}
    },
    7733: func() Message {
        return &GS_GuildDestroy_R{}
    },
    7734: func() Message {
        return &C_GuildChangeSetting{}
    },
    7735: func() Message {
        return &GS_GuildChangeSetting_R{}
    },
    7736: func() Message {
        return &C_GuildList{}
    },
    7737: func() Message {
        return &GS_GuildList_R{}
    },
    7738: func() Message {
        return &C_GuildPlrApplyList{}
    },
    7739: func() Message {
        return &GS_GuildPlrApplyList_R{}
    },
    7740: func() Message {
        return &C_GuildSearch{}
    },
    7741: func() Message {
        return &GS_GuildSearch_R{}
    },
    7742: func() Message {
        return &C_GuildApplyList{}
    },
    7743: func() Message {
        return &GS_GuildApplyList_R{}
    },
    7744: func() Message {
        return &C_GuildInfoFull{}
    },
    7745: func() Message {
        return &GS_GuildInfoFull_R{}
    },
    7746: func() Message {
        return &C_GuildApply{}
    },
    7747: func() Message {
        return &GS_GuildApply_R{}
    },
    7748: func() Message {
        return &C_GuildApplyCancel{}
    },
    7749: func() Message {
        return &GS_GuildApplyCancel_R{}
    },
    7750: func() Message {
        return &C_GuildApplyAccept{}
    },
    7751: func() Message {
        return &GS_GuildApplyAccept_R{}
    },
    7752: func() Message {
        return &C_GuildApplyDeny{}
    },
    7753: func() Message {
        return &GS_GuildApplyDeny_R{}
    },
    7754: func() Message {
        return &C_GuildLeave{}
    },
    7755: func() Message {
        return &GS_GuildLeave_R{}
    },
    7756: func() Message {
        return &C_GuildKick{}
    },
    7757: func() Message {
        return &GS_GuildKick_R{}
    },
    7758: func() Message {
        return &C_GuildSetRank{}
    },
    7759: func() Message {
        return &GS_GuildSetRank_R{}
    },
    7760: func() Message {
        return &C_GuildChangeName{}
    },
    7761: func() Message {
        return &GS_GuildChangeName_R{}
    },
    7762: func() Message {
        return &C_GuildApplyAcceptOneKey{}
    },
    7763: func() Message {
        return &GS_GuildApplyAcceptOneKey_R{}
    },
    7764: func() Message {
        return &C_GuildApplyDenyOneKey{}
    },
    7765: func() Message {
        return &GS_GuildApplyDenyOneKey_R{}
    },
    7766: func() Message {
        return &C_GuildKickOwner{}
    },
    7767: func() Message {
        return &GS_GuildKickOwner_R{}
    },
    7768: func() Message {
        return &C_GuildGetLog{}
    },
    7769: func() Message {
        return &GS_GuildGetLog_R{}
    },
    7770: func() Message {
        return &C_GuildSign{}
    },
    7771: func() Message {
        return &GS_GuildSign_R{}
    },
    7772: func() Message {
        return &C_GuildSetNotice{}
    },
    7773: func() Message {
        return &GS_GuildSetNotice_R{}
    },
    7774: func() Message {
        return &C_GuildSetIcon{}
    },
    7775: func() Message {
        return &GS_GuildSetIcon_R{}
    },
    7776: func() Message {
        return &C_GuildPublishZm{}
    },
    7777: func() Message {
        return &GS_GuildPublishZm_R{}
    },
    7800: func() Message {
        return &GS_GuildWishNew{}
    },
    7801: func() Message {
        return &GS_GuildWishFullHelp{}
    },
    7805: func() Message {
        return &C_GuildWishItem{}
    },
    7806: func() Message {
        return &GS_GuildWishItem_R{}
    },
    7807: func() Message {
        return &C_GuildWishHelp{}
    },
    7808: func() Message {
        return &GS_GuildWishHelp_R{}
    },
    7809: func() Message {
        return &C_GuildWishClose{}
    },
    7810: func() Message {
        return &GS_GuildWishClose_R{}
    },
    7811: func() Message {
        return &C_GuildWishList{}
    },
    7812: func() Message {
        return &GS_GuildWishList_R{}
    },
    7820: func() Message {
        return &GS_GuildHarborXpChange{}
    },
    7824: func() Message {
        return &C_GuildHarborDonate{}
    },
    7825: func() Message {
        return &GS_GuildHarborDonate_R{}
    },
    7826: func() Message {
        return &C_GuildHarborDonateList{}
    },
    7827: func() Message {
        return &GS_GuildHarborDonateList_R{}
    },
    7840: func() Message {
        return &C_GuildOrderGet{}
    },
    7841: func() Message {
        return &GS_GuildOrderGet_R{}
    },
    7842: func() Message {
        return &C_GuildOrderStarup{}
    },
    7843: func() Message {
        return &GS_GuildOrderStarup_R{}
    },
    7844: func() Message {
        return &C_GuildOrderStart{}
    },
    7845: func() Message {
        return &GS_GuildOrderStart_R{}
    },
    7846: func() Message {
        return &C_GuildOrderClose{}
    },
    7847: func() Message {
        return &GS_GuildOrderClose_R{}
    },
    7848: func() Message {
        return &C_GuildOrderList{}
    },
    7849: func() Message {
        return &GS_GuildOrderList_R{}
    },
    7860: func() Message {
        return &C_GuildTechLevelup{}
    },
    7861: func() Message {
        return &GS_GuildTechLevelup_R{}
    },
    7862: func() Message {
        return &C_GuildTechReset{}
    },
    7863: func() Message {
        return &GS_GuildTechReset_R{}
    },
    7864: func() Message {
        return &C_GuildTechGetInfo{}
    },
    7865: func() Message {
        return &GS_GuildTechGetInfo_R{}
    },
    7880: func() Message {
        return &C_GuildBossFight{}
    },
    7881: func() Message {
        return &GS_GuildBossFight_R{}
    },
    7882: func() Message {
        return &C_GuildBossGetCurrent{}
    },
    7883: func() Message {
        return &GS_GuildBossGetCurrent_R{}
    },
    7884: func() Message {
        return &C_GuildBossGetHistory{}
    },
    7885: func() Message {
        return &GS_GuildBossGetHistory_R{}
    },
    8000: func() Message {
        return &GS_MOpenModuleNew{}
    },
    8400: func() Message {
        return &GS_CounterOpUpdate{}
    },
    8410: func() Message {
        return &C_CounterRecover{}
    },
    8411: func() Message {
        return &GS_CounterRecover_R{}
    },
    8412: func() Message {
        return &C_CounterBuy{}
    },
    8413: func() Message {
        return &GS_CounterBuy_R{}
    },
    8450: func() Message {
        return &C_RankGet{}
    },
    8451: func() Message {
        return &GS_RankGet_R{}
    },
    8452: func() Message {
        return &C_RankLike{}
    },
    8453: func() Message {
        return &GS_RankLike_R{}
    },
    8600: func() Message {
        return &GS_TaskDailyValueChanged{}
    },
    8601: func() Message {
        return &GS_TaskDailyItemCompleted{}
    },
    8650: func() Message {
        return &C_TaskDailyInfo{}
    },
    8651: func() Message {
        return &GS_TaskDailyInfo_R{}
    },
    8654: func() Message {
        return &C_TaskDailyTakeBoxReward{}
    },
    8655: func() Message {
        return &GS_TaskDailyTakeBoxReward_R{}
    },
    8700: func() Message {
        return &C_TaskAchvTake{}
    },
    8701: func() Message {
        return &GS_TaskAchvTake_R{}
    },
    8800: func() Message {
        return &C_DrawGetInfo{}
    },
    8801: func() Message {
        return &GS_DrawGetInfo_R{}
    },
    8802: func() Message {
        return &C_DrawTp{}
    },
    8803: func() Message {
        return &GS_DrawTp_R{}
    },
    8804: func() Message {
        return &C_DrawScoreBoxTake{}
    },
    8805: func() Message {
        return &GS_DrawScoreBoxTake_R{}
    },
    8900: func() Message {
        return &C_WLevelFight{}
    },
    8901: func() Message {
        return &GS_WLevelFight_R{}
    },
    8902: func() Message {
        return &C_WLevelGJInfo{}
    },
    8903: func() Message {
        return &GS_WLevelGJInfo_R{}
    },
    8904: func() Message {
        return &C_WLevelGJTake{}
    },
    8905: func() Message {
        return &GS_WLevelGJTake_R{}
    },
    8906: func() Message {
        return &C_WLevelOneKeyGJTake{}
    },
    8907: func() Message {
        return &GS_WLevelOneKeyGJTake_R{}
    },
    8908: func() Message {
        return &C_WLevelFightOneKey{}
    },
    8909: func() Message {
        return &GS_WLevelFightOneKey_R{}
    },
    9000: func() Message {
        return &GS_AppointAddTask{}
    },
    9010: func() Message {
        return &C_AppointCheckAdd{}
    },
    9011: func() Message {
        return &GS_AppointCheckAdd_R{}
    },
    9012: func() Message {
        return &C_AppointRefresh{}
    },
    9013: func() Message {
        return &GS_AppointRefresh_R{}
    },
    9014: func() Message {
        return &C_AppointLock{}
    },
    9015: func() Message {
        return &GS_AppointLock_R{}
    },
    9016: func() Message {
        return &C_AppointSend{}
    },
    9017: func() Message {
        return &GS_AppointSend_R{}
    },
    9018: func() Message {
        return &C_AppointAcc{}
    },
    9019: func() Message {
        return &GS_AppointAcc_R{}
    },
    9020: func() Message {
        return &C_AppointTake{}
    },
    9021: func() Message {
        return &GS_AppointTake_R{}
    },
    9022: func() Message {
        return &C_AppointCancel{}
    },
    9023: func() Message {
        return &GS_AppointCancel_R{}
    },
    9100: func() Message {
        return &C_TowerFight{}
    },
    9101: func() Message {
        return &GS_TowerFight_R{}
    },
    9102: func() Message {
        return &C_TowerRaid{}
    },
    9103: func() Message {
        return &GS_TowerRaid_R{}
    },
    9104: func() Message {
        return &C_TowerRecord{}
    },
    9105: func() Message {
        return &GS_TowerRecord_R{}
    },
    9200: func() Message {
        return &C_SetTeam{}
    },
    9201: func() Message {
        return &GS_SetTeam_R{}
    },
    9300: func() Message {
        return &GS_ArenaStageUpdate{}
    },
    9301: func() Message {
        return &GS_ArenaFighted{}
    },
    9310: func() Message {
        return &C_ArenaUpdateEnemy{}
    },
    9311: func() Message {
        return &GS_ArenaUpdateEnemy_R{}
    },
    9312: func() Message {
        return &C_ArenaFight{}
    },
    9313: func() Message {
        return &GS_ArenaFight_R{}
    },
    9314: func() Message {
        return &C_ArenaRecordInfo{}
    },
    9315: func() Message {
        return &GS_ArenaRecordInfo_R{}
    },
    9316: func() Message {
        return &C_ArenaRank{}
    },
    9317: func() Message {
        return &GS_ArenaRank_R{}
    },
    9318: func() Message {
        return &C_ArenaReplayGet{}
    },
    9319: func() Message {
        return &GS_ArenaReplayGet_R{}
    },
    9400: func() Message {
        return &C_ShopGetInfo{}
    },
    9401: func() Message {
        return &GS_ShopGetInfo_R{}
    },
    9402: func() Message {
        return &C_ShopBuy{}
    },
    9403: func() Message {
        return &GS_ShopBuy_R{}
    },
    9404: func() Message {
        return &C_ShopRefresh{}
    },
    9405: func() Message {
        return &GS_ShopRefresh_R{}
    },
    9500: func() Message {
        return &C_MarvelRollInfo{}
    },
    9501: func() Message {
        return &GS_MarvelRollInfo_R{}
    },
    9502: func() Message {
        return &C_MarvelRollRefresh{}
    },
    9503: func() Message {
        return &GS_MarvelRollRefresh_R{}
    },
    9504: func() Message {
        return &C_MarvelRollTake{}
    },
    9505: func() Message {
        return &GS_MarvelRollTake_R{}
    },
    9600: func() Message {
        return &GS_FriendNewApplied{}
    },
    9601: func() Message {
        return &GS_FriendNewFriend{}
    },
    9602: func() Message {
        return &GS_FriendRemoveFriend{}
    },
    9603: func() Message {
        return &GS_FriendRecv{}
    },
    9610: func() Message {
        return &C_FriendGetFrds{}
    },
    9611: func() Message {
        return &GS_FriendGetFrds_R{}
    },
    9612: func() Message {
        return &C_FriendSearchFrds{}
    },
    9613: func() Message {
        return &GS_FriendSearchFrds_R{}
    },
    9614: func() Message {
        return &C_FriendRemoveFrds{}
    },
    9615: func() Message {
        return &GS_FriendRemoveFrds_R{}
    },
    9616: func() Message {
        return &C_FriendGetApplyList{}
    },
    9617: func() Message {
        return &GS_FriendGetApplyList_R{}
    },
    9618: func() Message {
        return &C_FriendApply{}
    },
    9619: func() Message {
        return &GS_FriendApply_R{}
    },
    9620: func() Message {
        return &C_FriendAccept{}
    },
    9621: func() Message {
        return &GS_FriendAccept_R{}
    },
    9622: func() Message {
        return &C_FriendGetBlackList{}
    },
    9623: func() Message {
        return &GS_FriendGetBlackList_R{}
    },
    9624: func() Message {
        return &C_FriendAddBlackList{}
    },
    9625: func() Message {
        return &GS_FriendAddBlackList_R{}
    },
    9626: func() Message {
        return &C_FriendGiveAndRecv{}
    },
    9627: func() Message {
        return &GS_FriendGiveAndRecv_R{}
    },
    9700: func() Message {
        return &GS_CrusadeStageUpdate{}
    },
    9710: func() Message {
        return &C_CrusadeGetInfo{}
    },
    9711: func() Message {
        return &GS_CrusadeGetInfo_R{}
    },
    9712: func() Message {
        return &C_CrusadeBoxTake{}
    },
    9713: func() Message {
        return &GS_CrusadeBoxTake_R{}
    },
    9714: func() Message {
        return &C_CrusadeFight{}
    },
    9715: func() Message {
        return &GS_CrusadeFight_R{}
    },
    9800: func() Message {
        return &GS_VipUpdate{}
    },
    10000: func() Message {
        return &C_ActRushLocalGetInfo{}
    },
    10001: func() Message {
        return &GS_ActRushLocalGetInfo_R{}
    },
    10002: func() Message {
        return &C_ActRushLocalTake{}
    },
    10003: func() Message {
        return &GS_ActRushLocalTake_R{}
    },
    10010: func() Message {
        return &GS_ActBillLtTotal{}
    },
    10011: func() Message {
        return &C_ActBillLtTotalInfo{}
    },
    10012: func() Message {
        return &GS_ActBillLtTotalInfo_R{}
    },
    10013: func() Message {
        return &C_ActBillLtTotalTake{}
    },
    10014: func() Message {
        return &GS_ActBillLtTotalTake_R{}
    },
    10020: func() Message {
        return &C_ActBillLtDayInfo{}
    },
    10021: func() Message {
        return &GS_ActBillLtDayInfo_R{}
    },
    10022: func() Message {
        return &C_ActBillLtDayTake{}
    },
    10023: func() Message {
        return &GS_ActBillLtDayTake_R{}
    },
    10031: func() Message {
        return &GS_ActGiftNew{}
    },
    10040: func() Message {
        return &C_ActSummonInfo{}
    },
    10041: func() Message {
        return &GS_ActSummonInfo_R{}
    },
    10042: func() Message {
        return &C_ActSummonPick{}
    },
    10043: func() Message {
        return &GS_ActSummonPick_R{}
    },
    10044: func() Message {
        return &C_ActSummonDraw{}
    },
    10045: func() Message {
        return &GS_ActSummonDraw_R{}
    },
    10050: func() Message {
        return &GS_ActTargetTaskObjValueChanged{}
    },
    10052: func() Message {
        return &C_ActTargetTaskInfo{}
    },
    10053: func() Message {
        return &GS_ActTargetTaskInfo_R{}
    },
    10054: func() Message {
        return &C_ActTargetTaskTake{}
    },
    10055: func() Message {
        return &GS_ActTargetTaskTake_R{}
    },
    10060: func() Message {
        return &C_ActHeroSkinInfo{}
    },
    10061: func() Message {
        return &GS_ActHeroSkinInfo_R{}
    },
    10070: func() Message {
        return &C_ActMSummonInfo{}
    },
    10071: func() Message {
        return &GS_ActMSummonInfo_R{}
    },
    10072: func() Message {
        return &C_ActMSummonDraw{}
    },
    10073: func() Message {
        return &GS_ActMSummonDraw_R{}
    },
    10080: func() Message {
        return &GS_ActMonopolyObjValueChanged{}
    },
    10081: func() Message {
        return &GS_ActMonoPolyNextLv{}
    },
    10082: func() Message {
        return &GS_ActMonopolyBoxReward{}
    },
    10090: func() Message {
        return &C_ActMonopolyInfo{}
    },
    10091: func() Message {
        return &GS_ActMonoPolyInfo_R{}
    },
    10092: func() Message {
        return &C_ActMonopolyNPos{}
    },
    10093: func() Message {
        return &GS_ActMonopolyNPos_R{}
    },
    10094: func() Message {
        return &C_ActMonopolyAnswer{}
    },
    10095: func() Message {
        return &GS_ActMonopolyAnswer_R{}
    },
    10096: func() Message {
        return &C_ActMonopolyBuy{}
    },
    10097: func() Message {
        return &GS_ActMonopolyBuy_R{}
    },
    10098: func() Message {
        return &C_ActMonopolyBattle{}
    },
    10099: func() Message {
        return &GS_ActMonopolyBattle_R{}
    },
    10100: func() Message {
        return &C_ActMonopolyTaskTake{}
    },
    10101: func() Message {
        return &GS_ActMonopolyTaskTake_R{}
    },
    10110: func() Message {
        return &GS_ActMazeObjValueChanged{}
    },
    10111: func() Message {
        return &GS_ActMazeTaskTaken{}
    },
    10115: func() Message {
        return &C_ActMazeInfo{}
    },
    10116: func() Message {
        return &GS_ActMazeInfo_R{}
    },
    10117: func() Message {
        return &C_ActMazeClick{}
    },
    10118: func() Message {
        return &GS_ActMazeClick_R{}
    },
    10119: func() Message {
        return &C_ActMazeClickNext{}
    },
    10120: func() Message {
        return &GS_ActMazeClickNext_R{}
    },
    10121: func() Message {
        return &C_ActMazeClickTrade{}
    },
    10122: func() Message {
        return &GS_ActMazeClickTrade_R{}
    },
    10123: func() Message {
        return &C_ActMazeReset{}
    },
    10124: func() Message {
        return &GS_ActMazeReset_R{}
    },
    10125: func() Message {
        return &C_ActMazeClickThing{}
    },
    10126: func() Message {
        return &GS_ActMazeClickThing_R{}
    },
    10127: func() Message {
        return &C_ActMazeClickBattle{}
    },
    10128: func() Message {
        return &GS_ActMazeClickBattle_R{}
    },
    10129: func() Message {
        return &C_ActMazeTakeTask{}
    },
    10130: func() Message {
        return &GS_ActMazeTakeTask_R{}
    },
    10131: func() Message {
        return &C_ActMazeBuff{}
    },
    10132: func() Message {
        return &GS_ActMazeBuff_R{}
    },
    13000: func() Message {
        return &GS_ChatMsg{}
    },
    13010: func() Message {
        return &C_ChatSend{}
    },
    13011: func() Message {
        return &GS_ChatSend_R{}
    },
    13200: func() Message {
        return &GS_MonthTicketValueChanged{}
    },
    13201: func() Message {
        return &GS_MonthTicketItemCompleted{}
    },
    13202: func() Message {
        return &GS_MonthTicketIsBuy{}
    },
    13221: func() Message {
        return &C_MonthTicketInfo{}
    },
    13222: func() Message {
        return &GS_MonthTicketInfo_R{}
    },
    13223: func() Message {
        return &C_MonthTicketTakeOneKey{}
    },
    13224: func() Message {
        return &GS_MonthTicketTakeOneKey_R{}
    },
    13225: func() Message {
        return &C_MonthTicketBuyUp{}
    },
    13226: func() Message {
        return &GS_MonthTicketBuyUp_R{}
    },
    13227: func() Message {
        return &C_MonthTicketTake{}
    },
    13228: func() Message {
        return &GS_MonthTicketTake_R{}
    },
    13229: func() Message {
        return &C_MonthTicketTaskTake{}
    },
    13230: func() Message {
        return &GS_MonthTicketTaskTake_R{}
    },
    13300: func() Message {
        return &GS_PushGiftNew{}
    },
    13301: func() Message {
        return &GS_PushGiftRewards{}
    },
    13310: func() Message {
        return &C_PushGiftSetCreateTs{}
    },
    13311: func() Message {
        return &GS_PushGiftSetCreateTs_R{}
    },
    13400: func() Message {
        return &GS_GiftShopNew{}
    },
    13410: func() Message {
        return &C_GiftShopTake{}
    },
    13411: func() Message {
        return &GS_GiftShopTake_R{}
    },
    13500: func() Message {
        return &GS_PrivCardNew{}
    },
    13510: func() Message {
        return &C_PrivCardTake{}
    },
    13511: func() Message {
        return &GS_PrivCardTake_R{}
    },
    13601: func() Message {
        return &C_SignDailySign{}
    },
    13602: func() Message {
        return &GS_SignDailySign_R{}
    },
    13700: func() Message {
        return &GS_TaskMonthValueChanged{}
    },
    13701: func() Message {
        return &GS_TaskMonthItemCompleted{}
    },
    13702: func() Message {
        return &C_TaskMonthInfo{}
    },
    13703: func() Message {
        return &GS_TaskMonthInfo_R{}
    },
    13704: func() Message {
        return &C_TaskMonthTake{}
    },
    13705: func() Message {
        return &GS_TaskMonthTask_R{}
    },
    13820: func() Message {
        return &C_DaySignTake{}
    },
    13821: func() Message {
        return &GS_DaySignTake_R{}
    },
    13900: func() Message {
        return &C_TargetDaysTake{}
    },
    13901: func() Message {
        return &GS_TargetDaysTake_R{}
    },
    13902: func() Message {
        return &C_TargetDaysBuy{}
    },
    13903: func() Message {
        return &GS_TargetDaysBuy_R{}
    },
    14030: func() Message {
        return &C_TaskGrowTake{}
    },
    14031: func() Message {
        return &GS_TaskGrowTake_R{}
    },
    14000: func() Message {
        return &GS_WLevelFundNew{}
    },
    14010: func() Message {
        return &C_WLevelFundTake{}
    },
    14011: func() Message {
        return &GS_WLevelFundTake_R{}
    },
    14100: func() Message {
        return &GS_GrowFundNew{}
    },
    14110: func() Message {
        return &C_GrowFundInfo{}
    },
    14111: func() Message {
        return &GS_GrowFundInfo_R{}
    },
    14112: func() Message {
        return &C_GrowFundTakeLv{}
    },
    14113: func() Message {
        return &GS_GrowFundTakeLv_R{}
    },
    14114: func() Message {
        return &C_GrowFundTakeSvr{}
    },
    14115: func() Message {
        return &GS_GrowFundTakeSvr_R{}
    },
    14201: func() Message {
        return &GS_BillFirstNew{}
    },
    14222: func() Message {
        return &C_BillFirstTake{}
    },
    14223: func() Message {
        return &GS_BillFirstTake_R{}
    },
    14300: func() Message {
        return &GS_LampMsg{}
    },
    14400: func() Message {
        return &GS_GWarStageChange{}
    },
    14401: func() Message {
        return &GS_GWarNewG2{}
    },
    14410: func() Message {
        return &C_GWarGetSummary{}
    },
    14411: func() Message {
        return &GS_GWarGetSummary_R{}
    },
    14412: func() Message {
        return &C_GWarGetG2Members{}
    },
    14413: func() Message {
        return &GS_GWarGetG2Members_R{}
    },
    14414: func() Message {
        return &C_GWarGetGuildRank{}
    },
    14415: func() Message {
        return &GS_GWarGetGuildRank_R{}
    },
    14416: func() Message {
        return &C_GWarGetPlrRank{}
    },
    14417: func() Message {
        return &GS_GWarGetPlrRank_R{}
    },
    14418: func() Message {
        return &C_GWarFight{}
    },
    14419: func() Message {
        return &GS_GWarFight_R{}
    },
    14500: func() Message {
        return &GS_RiftMonsterNew{}
    },
    14501: func() Message {
        return &GS_RiftMineNew{}
    },
    14502: func() Message {
        return &GS_RiftMineOccupied{}
    },
    14503: func() Message {
        return &GS_RiftBoxNew{}
    },
    14504: func() Message {
        return &GS_RiftBoxRewards{}
    },
    14505: func() Message {
        return &GS_RiftBoxOccupied{}
    },
    14560: func() Message {
        return &C_RiftExplore{}
    },
    14561: func() Message {
        return &GS_RiftExplore_R{}
    },
    14562: func() Message {
        return &C_RiftMonsterFight{}
    },
    14563: func() Message {
        return &GS_RiftMonsterFight_R{}
    },
    14564: func() Message {
        return &C_RiftMineInfo{}
    },
    14565: func() Message {
        return &GS_RiftMineInfo_R{}
    },
    14566: func() Message {
        return &C_RiftMineOccupy{}
    },
    14567: func() Message {
        return &GS_RiftMineOccupy_R{}
    },
    14568: func() Message {
        return &C_RiftMineCancel{}
    },
    14569: func() Message {
        return &GS_RiftMineCancel_R{}
    },
    14570: func() Message {
        return &C_RiftMineTakeRewards{}
    },
    14571: func() Message {
        return &GS_RiftMineTakeRewards_R{}
    },
    14572: func() Message {
        return &C_RiftBoxOccupy{}
    },
    14573: func() Message {
        return &GS_RiftBoxOccupy_R{}
    },
    14574: func() Message {
        return &C_RiftBoxInfo{}
    },
    14575: func() Message {
        return &GS_RiftBoxInfo_R{}
    },
    14600: func() Message {
        return &GS_LadderStageChange{}
    },
    14610: func() Message {
        return &C_LadderGetSummary{}
    },
    14611: func() Message {
        return &GS_LadderGetSummary_R{}
    },
    14612: func() Message {
        return &C_LadderMatch{}
    },
    14613: func() Message {
        return &GS_LadderMatch_R{}
    },
    14614: func() Message {
        return &C_LadderFight{}
    },
    14615: func() Message {
        return &GS_LadderFight_R{}
    },
    14616: func() Message {
        return &C_LadderGetRank{}
    },
    14617: func() Message {
        return &GS_LadderGetRank_R{}
    },
    14618: func() Message {
        return &C_LadderGetReplayList{}
    },
    14619: func() Message {
        return &GS_LadderGetReplayList_R{}
    },
    14620: func() Message {
        return &C_LadderGetReplay{}
    },
    14621: func() Message {
        return &GS_LadderGetReplay_R{}
    },
    14712: func() Message {
        return &C_HeroSkinAdd{}
    },
    14713: func() Message {
        return &GS_HeroSkinAdd_R{}
    },
    14714: func() Message {
        return &C_HeroSkinLvUp{}
    },
    14715: func() Message {
        return &GS_HeroSkinLvUp_R{}
    },
    14800: func() Message {
        return &C_WLevelDrawDraw{}
    },
    14801: func() Message {
        return &GS_WLevelDrawDraw_R{}
    },
    14802: func() Message {
        return &C_WLevelDrawTake{}
    },
    14803: func() Message {
        return &GS_WLevelDrawTake_R{}
    },
    14900: func() Message {
        return &GS_WarCupStageUpdate{}
    },
    14901: func() Message {
        return &GS_WarCupChat{}
    },
    14902: func() Message {
        return &GS_WarCupGuessRatio{}
    },
    14903: func() Message {
        return &GS_WarCupAttainObjValueChanged{}
    },
    14920: func() Message {
        return &C_WarCupGuessInfo{}
    },
    14921: func() Message {
        return &GS_WarCupGuessInfo_R{}
    },
    14922: func() Message {
        return &C_WarCupSelfInfo{}
    },
    14923: func() Message {
        return &GS_WarCupSelfInfo_R{}
    },
    14924: func() Message {
        return &C_WarCupTop64Info{}
    },
    14925: func() Message {
        return &GS_WarCupTop64Info_R{}
    },
    14926: func() Message {
        return &C_WarCupTop8Info{}
    },
    14927: func() Message {
        return &GS_WarCupTop8Info_R{}
    },
    14928: func() Message {
        return &C_WarCupGuess{}
    },
    14929: func() Message {
        return &GS_WarCupGuess_R{}
    },
    14930: func() Message {
        return &C_WarCupAuditionRank{}
    },
    14931: func() Message {
        return &GS_WarCupAuditionRank_R{}
    },
    14932: func() Message {
        return &C_WarCupChatSend{}
    },
    14933: func() Message {
        return &GS_WarCupChatSend_R{}
    },
    14934: func() Message {
        return &C_WarCupGetReplay{}
    },
    14935: func() Message {
        return &GS_WarCupGetReplay_R{}
    },
    14936: func() Message {
        return &C_WarCupTop1Info{}
    },
    14937: func() Message {
        return &GS_WarCupTop1Info_R{}
    },
    14938: func() Message {
        return &C_WarCupTaskTake{}
    },
    14939: func() Message {
        return &GS_WarCupTaskTake_R{}
    },
    14940: func() Message {
        return &C_WarCupWatch{}
    },
    14942: func() Message {
        return &C_WarCupGuessRecords{}
    },
    14943: func() Message {
        return &GS_WarCupGuessRecords_R{}
    },
    15000: func() Message {
        return &GS_WBossStageChange{}
    },
    15010: func() Message {
        return &C_WBossGetSummary{}
    },
    15011: func() Message {
        return &GS_WBossGetSummary_R{}
    },
    15012: func() Message {
        return &C_WBossFight{}
    },
    15013: func() Message {
        return &GS_WBossFight_R{}
    },
    15014: func() Message {
        return &C_WBossTakeMaxDmgRwd{}
    },
    15015: func() Message {
        return &GS_WBossTakeMaxDmgRwd_R{}
    },
    15016: func() Message {
        return &C_WBossGetRank{}
    },
    15017: func() Message {
        return &GS_WBossGetRank_R{}
    },
    15018: func() Message {
        return &C_WBossGetMaxDmgInfo{}
    },
    15019: func() Message {
        return &GS_WBossGetMaxDmgInfo_R{}
    },
    15020: func() Message {
        return &C_WBossGetSelfRank{}
    },
    15021: func() Message {
        return &GS_WBossGetSelfRank_R{}
    },
    15130: func() Message {
        return &C_InviteInfo{}
    },
    15131: func() Message {
        return &GS_InviteInfo_R{}
    },
    15132: func() Message {
        return &C_InviteTake{}
    },
    15133: func() Message {
        return &GS_InviteTake_R{}
    },
}

func (self *C_Auth) MsgId() uint32 {
    return 1000
}

func (self *GW_Auth_R) MsgId() uint32 {
    return 1001
}

func (self *C_Login) MsgId() uint32 {
    return 1002
}

func (self *GW_Login_R) MsgId() uint32 {
    return 1003
}

func (self *C_TokenGet) MsgId() uint32 {
    return 1004
}

func (self *GW_TokenGet_R) MsgId() uint32 {
    return 1005
}

func (self *C_TokenAuth) MsgId() uint32 {
    return 1006
}

func (self *GW_TokenAuth_R) MsgId() uint32 {
    return 1007
}

func (self *GS_LoginError) MsgId() uint32 {
    return 5000
}

func (self *GS_UserInfo) MsgId() uint32 {
    return 5001
}

func (self *C_TimeSync) MsgId() uint32 {
    return 5002
}

func (self *GS_TimeSync_R) MsgId() uint32 {
    return 5003
}

func (self *GS_GameDataReloaded) MsgId() uint32 {
    return 5004
}

func (self *C_UserExtInfo) MsgId() uint32 {
    return 5006
}

func (self *GS_UserExtInfo_R) MsgId() uint32 {
    return 5007
}

func (self *GS_PlayerUpdateLv) MsgId() uint32 {
    return 5100
}

func (self *GS_PlayerUpdateAtkPwr) MsgId() uint32 {
    return 5101
}

func (self *GS_PlayerUpdateHFrame) MsgId() uint32 {
    return 5102
}

func (self *C_PlayerChangeName) MsgId() uint32 {
    return 5132
}

func (self *GS_PlayerChangeName_R) MsgId() uint32 {
    return 5133
}

func (self *C_PlayerInfo) MsgId() uint32 {
    return 5134
}

func (self *GS_PlayerInfo_R) MsgId() uint32 {
    return 5135
}

func (self *C_PlayerHFrameAdd) MsgId() uint32 {
    return 5138
}

func (self *GS_PlayerHFrameAdd_R) MsgId() uint32 {
    return 5139
}

func (self *C_PlayerHFrameSet) MsgId() uint32 {
    return 5140
}

func (self *GS_PlayerHFrameSet_R) MsgId() uint32 {
    return 5141
}

func (self *C_PlayerHeadSet) MsgId() uint32 {
    return 5142
}

func (self *GS_PlayerHeadSet_R) MsgId() uint32 {
    return 5143
}

func (self *GS_BagUpdate) MsgId() uint32 {
    return 5200
}

func (self *GS_ArmorUpdate_HeroSeq) MsgId() uint32 {
    return 5300
}

func (self *C_ItemExchange) MsgId() uint32 {
    return 5310
}

func (self *GS_ItemExchange_R) MsgId() uint32 {
    return 5311
}

func (self *C_ItemUse) MsgId() uint32 {
    return 5320
}

func (self *GS_ItemUse_R) MsgId() uint32 {
    return 5321
}

func (self *C_ArmorEquip) MsgId() uint32 {
    return 5322
}

func (self *GS_ArmorEquip_R) MsgId() uint32 {
    return 5323
}

func (self *C_ArmorUnequip) MsgId() uint32 {
    return 5324
}

func (self *GS_ArmorUnequip_R) MsgId() uint32 {
    return 5325
}

func (self *C_ArmorEquipOnekey) MsgId() uint32 {
    return 5326
}

func (self *GS_ArmorEquipOnekey_R) MsgId() uint32 {
    return 5327
}

func (self *C_ArmorUnequipOnekey) MsgId() uint32 {
    return 5328
}

func (self *GS_ArmorUnequipOnekey_R) MsgId() uint32 {
    return 5329
}

func (self *C_ArmorCompose) MsgId() uint32 {
    return 5330
}

func (self *GS_ArmorCompose_R) MsgId() uint32 {
    return 5331
}

func (self *C_ArmorComposeOnekey) MsgId() uint32 {
    return 5332
}

func (self *GS_ArmorComposeOnekey_R) MsgId() uint32 {
    return 5333
}

func (self *C_ItemChoose) MsgId() uint32 {
    return 5334
}

func (self *GS_ItemChoose_R) MsgId() uint32 {
    return 5335
}

func (self *GS_HeroUpdate) MsgId() uint32 {
    return 5400
}

func (self *C_HeroLevelUp) MsgId() uint32 {
    return 5410
}

func (self *GS_HeroLevelUp_R) MsgId() uint32 {
    return 5411
}

func (self *C_HeroStarUp) MsgId() uint32 {
    return 5414
}

func (self *GS_HeroStarUp_R) MsgId() uint32 {
    return 5415
}

func (self *C_HeroLock) MsgId() uint32 {
    return 5416
}

func (self *GS_HeroLock_R) MsgId() uint32 {
    return 5417
}

func (self *C_HeroReset) MsgId() uint32 {
    return 5418
}

func (self *GS_HeroReset_R) MsgId() uint32 {
    return 5419
}

func (self *C_HeroDecompose) MsgId() uint32 {
    return 5420
}

func (self *GS_HeroDecompose_R) MsgId() uint32 {
    return 5421
}

func (self *C_HeroChange) MsgId() uint32 {
    return 5422
}

func (self *GS_HeroChange_R) MsgId() uint32 {
    return 5423
}

func (self *C_HeroChangeCancel) MsgId() uint32 {
    return 5424
}

func (self *GS_HeroChangeCancel_R) MsgId() uint32 {
    return 5425
}

func (self *C_HeroChangeApply) MsgId() uint32 {
    return 5426
}

func (self *GS_HeroChangeApply_R) MsgId() uint32 {
    return 5427
}

func (self *C_HeroTrinketUnlock) MsgId() uint32 {
    return 5428
}

func (self *GS_HeroTrinketUnlock_R) MsgId() uint32 {
    return 5429
}

func (self *C_HeroTrinketUp) MsgId() uint32 {
    return 5430
}

func (self *GS_HeroTrinketUp_R) MsgId() uint32 {
    return 5431
}

func (self *C_HeroTrinketTransformGen) MsgId() uint32 {
    return 5432
}

func (self *GS_HeroTrinketTransformGen_R) MsgId() uint32 {
    return 5433
}

func (self *C_HeroTrinketTransformCommit) MsgId() uint32 {
    return 5434
}

func (self *GS_HeroTrinketTransformCommit_R) MsgId() uint32 {
    return 5435
}

func (self *C_HeroBagBuy) MsgId() uint32 {
    return 5436
}

func (self *GS_HeroBagBuy_R) MsgId() uint32 {
    return 5437
}

func (self *C_HeroInherit) MsgId() uint32 {
    return 5438
}

func (self *GS_HeroInherit_R) MsgId() uint32 {
    return 5439
}

func (self *C_HeroSetSkin) MsgId() uint32 {
    return 5440
}

func (self *GS_HeroSetSkin_R) MsgId() uint32 {
    return 5441
}

func (self *GS_RelicUpdate_HeroSeq) MsgId() uint32 {
    return 5500
}

func (self *GS_RelicUpdate_Star) MsgId() uint32 {
    return 5501
}

func (self *C_RelicEquip) MsgId() uint32 {
    return 5510
}

func (self *GS_RelicEquip_R) MsgId() uint32 {
    return 5511
}

func (self *C_RelicUnequip) MsgId() uint32 {
    return 5512
}

func (self *GS_RelicUnequip_R) MsgId() uint32 {
    return 5513
}

func (self *C_RelicEat) MsgId() uint32 {
    return 5514
}

func (self *GS_RelicEat_R) MsgId() uint32 {
    return 5515
}

func (self *GS_MailNew) MsgId() uint32 {
    return 7000
}

func (self *GS_MailDel) MsgId() uint32 {
    return 7001
}

func (self *C_MailRead) MsgId() uint32 {
    return 7010
}

func (self *GS_MailRead_R) MsgId() uint32 {
    return 7011
}

func (self *C_MailDel) MsgId() uint32 {
    return 7012
}

func (self *GS_MailDel_R) MsgId() uint32 {
    return 7013
}

func (self *C_MailTakeAttachment) MsgId() uint32 {
    return 7014
}

func (self *GS_MailTakeAttachment_R) MsgId() uint32 {
    return 7015
}

func (self *C_MailTakeAttachmentAll) MsgId() uint32 {
    return 7016
}

func (self *GS_MailTakeAttachmentAll_R) MsgId() uint32 {
    return 7017
}

func (self *C_MailDelOnekey) MsgId() uint32 {
    return 7018
}

func (self *GS_MailDelOnekey_R) MsgId() uint32 {
    return 7019
}

func (self *C_CloudGet) MsgId() uint32 {
    return 7100
}

func (self *GS_CloudGet_R) MsgId() uint32 {
    return 7101
}

func (self *C_CloudSet) MsgId() uint32 {
    return 7102
}

func (self *GS_CloudSet_R) MsgId() uint32 {
    return 7103
}

func (self *C_TutorialSet) MsgId() uint32 {
    return 7200
}

func (self *GS_TutorialSet_R) MsgId() uint32 {
    return 7201
}

func (self *GS_ActStateChange) MsgId() uint32 {
    return 7300
}

func (self *C_ActStateGet) MsgId() uint32 {
    return 7320
}

func (self *GS_ActStateGet_R) MsgId() uint32 {
    return 7321
}

func (self *GS_AttainTabObjValueChanged) MsgId() uint32 {
    return 7400
}

func (self *GS_BillDone) MsgId() uint32 {
    return 7500
}

func (self *GS_BillOrder) MsgId() uint32 {
    return 7501
}

func (self *C_BillInfo) MsgId() uint32 {
    return 7510
}

func (self *GS_BillInfo_R) MsgId() uint32 {
    return 7511
}

func (self *C_BillRefundCodeGet) MsgId() uint32 {
    return 7512
}

func (self *GS_BillRefundCodeGet_R) MsgId() uint32 {
    return 7513
}

func (self *C_BillRefundTake) MsgId() uint32 {
    return 7514
}

func (self *GS_BillRefundTake_R) MsgId() uint32 {
    return 7515
}

func (self *GS_MiscGldActGift) MsgId() uint32 {
    return 7600
}

func (self *C_MiscBillLocal) MsgId() uint32 {
    return 7620
}

func (self *GS_MiscBillLocal_R) MsgId() uint32 {
    return 7621
}

func (self *C_GiftExchange) MsgId() uint32 {
    return 7622
}

func (self *GS_GiftExchange_R) MsgId() uint32 {
    return 7623
}

func (self *C_MiscSkipTutorial) MsgId() uint32 {
    return 7624
}

func (self *GS_MiscSkipTutorial_R) MsgId() uint32 {
    return 7625
}

func (self *C_MiscGoldenHand) MsgId() uint32 {
    return 7626
}

func (self *GS_MiscGoldenHand_R) MsgId() uint32 {
    return 7627
}

func (self *C_MiscOnlineBoxTake) MsgId() uint32 {
    return 7628
}

func (self *GS_MiscOnlineBoxTake_R) MsgId() uint32 {
    return 7629
}

func (self *C_MiscSharedGame) MsgId() uint32 {
    return 7630
}

func (self *GS_MiscSharedGame_R) MsgId() uint32 {
    return 7631
}

func (self *C_MiscGldActGiftTake) MsgId() uint32 {
    return 7632
}

func (self *GS_MiscGldActGiftTake_R) MsgId() uint32 {
    return 7633
}

func (self *GS_GuildPlrLeaveTs) MsgId() uint32 {
    return 7700
}

func (self *GS_Guild_Join) MsgId() uint32 {
    return 7701
}

func (self *GS_Guild_Leave) MsgId() uint32 {
    return 7702
}

func (self *GS_Guild_MbRank) MsgId() uint32 {
    return 7703
}

func (self *GS_Guild_Lv) MsgId() uint32 {
    return 7704
}

func (self *GS_Guild_Notice) MsgId() uint32 {
    return 7705
}

func (self *GS_Guild_Icon) MsgId() uint32 {
    return 7706
}

func (self *GS_Guild_NewApply) MsgId() uint32 {
    return 7707
}

func (self *C_GuildCreate) MsgId() uint32 {
    return 7730
}

func (self *GS_GuildCreate_R) MsgId() uint32 {
    return 7731
}

func (self *C_GuildDestroy) MsgId() uint32 {
    return 7732
}

func (self *GS_GuildDestroy_R) MsgId() uint32 {
    return 7733
}

func (self *C_GuildChangeSetting) MsgId() uint32 {
    return 7734
}

func (self *GS_GuildChangeSetting_R) MsgId() uint32 {
    return 7735
}

func (self *C_GuildList) MsgId() uint32 {
    return 7736
}

func (self *GS_GuildList_R) MsgId() uint32 {
    return 7737
}

func (self *C_GuildPlrApplyList) MsgId() uint32 {
    return 7738
}

func (self *GS_GuildPlrApplyList_R) MsgId() uint32 {
    return 7739
}

func (self *C_GuildSearch) MsgId() uint32 {
    return 7740
}

func (self *GS_GuildSearch_R) MsgId() uint32 {
    return 7741
}

func (self *C_GuildApplyList) MsgId() uint32 {
    return 7742
}

func (self *GS_GuildApplyList_R) MsgId() uint32 {
    return 7743
}

func (self *C_GuildInfoFull) MsgId() uint32 {
    return 7744
}

func (self *GS_GuildInfoFull_R) MsgId() uint32 {
    return 7745
}

func (self *C_GuildApply) MsgId() uint32 {
    return 7746
}

func (self *GS_GuildApply_R) MsgId() uint32 {
    return 7747
}

func (self *C_GuildApplyCancel) MsgId() uint32 {
    return 7748
}

func (self *GS_GuildApplyCancel_R) MsgId() uint32 {
    return 7749
}

func (self *C_GuildApplyAccept) MsgId() uint32 {
    return 7750
}

func (self *GS_GuildApplyAccept_R) MsgId() uint32 {
    return 7751
}

func (self *C_GuildApplyDeny) MsgId() uint32 {
    return 7752
}

func (self *GS_GuildApplyDeny_R) MsgId() uint32 {
    return 7753
}

func (self *C_GuildLeave) MsgId() uint32 {
    return 7754
}

func (self *GS_GuildLeave_R) MsgId() uint32 {
    return 7755
}

func (self *C_GuildKick) MsgId() uint32 {
    return 7756
}

func (self *GS_GuildKick_R) MsgId() uint32 {
    return 7757
}

func (self *C_GuildSetRank) MsgId() uint32 {
    return 7758
}

func (self *GS_GuildSetRank_R) MsgId() uint32 {
    return 7759
}

func (self *C_GuildChangeName) MsgId() uint32 {
    return 7760
}

func (self *GS_GuildChangeName_R) MsgId() uint32 {
    return 7761
}

func (self *C_GuildApplyAcceptOneKey) MsgId() uint32 {
    return 7762
}

func (self *GS_GuildApplyAcceptOneKey_R) MsgId() uint32 {
    return 7763
}

func (self *C_GuildApplyDenyOneKey) MsgId() uint32 {
    return 7764
}

func (self *GS_GuildApplyDenyOneKey_R) MsgId() uint32 {
    return 7765
}

func (self *C_GuildKickOwner) MsgId() uint32 {
    return 7766
}

func (self *GS_GuildKickOwner_R) MsgId() uint32 {
    return 7767
}

func (self *C_GuildGetLog) MsgId() uint32 {
    return 7768
}

func (self *GS_GuildGetLog_R) MsgId() uint32 {
    return 7769
}

func (self *C_GuildSign) MsgId() uint32 {
    return 7770
}

func (self *GS_GuildSign_R) MsgId() uint32 {
    return 7771
}

func (self *C_GuildSetNotice) MsgId() uint32 {
    return 7772
}

func (self *GS_GuildSetNotice_R) MsgId() uint32 {
    return 7773
}

func (self *C_GuildSetIcon) MsgId() uint32 {
    return 7774
}

func (self *GS_GuildSetIcon_R) MsgId() uint32 {
    return 7775
}

func (self *C_GuildPublishZm) MsgId() uint32 {
    return 7776
}

func (self *GS_GuildPublishZm_R) MsgId() uint32 {
    return 7777
}

func (self *GS_GuildWishNew) MsgId() uint32 {
    return 7800
}

func (self *GS_GuildWishFullHelp) MsgId() uint32 {
    return 7801
}

func (self *C_GuildWishItem) MsgId() uint32 {
    return 7805
}

func (self *GS_GuildWishItem_R) MsgId() uint32 {
    return 7806
}

func (self *C_GuildWishHelp) MsgId() uint32 {
    return 7807
}

func (self *GS_GuildWishHelp_R) MsgId() uint32 {
    return 7808
}

func (self *C_GuildWishClose) MsgId() uint32 {
    return 7809
}

func (self *GS_GuildWishClose_R) MsgId() uint32 {
    return 7810
}

func (self *C_GuildWishList) MsgId() uint32 {
    return 7811
}

func (self *GS_GuildWishList_R) MsgId() uint32 {
    return 7812
}

func (self *GS_GuildHarborXpChange) MsgId() uint32 {
    return 7820
}

func (self *C_GuildHarborDonate) MsgId() uint32 {
    return 7824
}

func (self *GS_GuildHarborDonate_R) MsgId() uint32 {
    return 7825
}

func (self *C_GuildHarborDonateList) MsgId() uint32 {
    return 7826
}

func (self *GS_GuildHarborDonateList_R) MsgId() uint32 {
    return 7827
}

func (self *C_GuildOrderGet) MsgId() uint32 {
    return 7840
}

func (self *GS_GuildOrderGet_R) MsgId() uint32 {
    return 7841
}

func (self *C_GuildOrderStarup) MsgId() uint32 {
    return 7842
}

func (self *GS_GuildOrderStarup_R) MsgId() uint32 {
    return 7843
}

func (self *C_GuildOrderStart) MsgId() uint32 {
    return 7844
}

func (self *GS_GuildOrderStart_R) MsgId() uint32 {
    return 7845
}

func (self *C_GuildOrderClose) MsgId() uint32 {
    return 7846
}

func (self *GS_GuildOrderClose_R) MsgId() uint32 {
    return 7847
}

func (self *C_GuildOrderList) MsgId() uint32 {
    return 7848
}

func (self *GS_GuildOrderList_R) MsgId() uint32 {
    return 7849
}

func (self *C_GuildTechLevelup) MsgId() uint32 {
    return 7860
}

func (self *GS_GuildTechLevelup_R) MsgId() uint32 {
    return 7861
}

func (self *C_GuildTechReset) MsgId() uint32 {
    return 7862
}

func (self *GS_GuildTechReset_R) MsgId() uint32 {
    return 7863
}

func (self *C_GuildTechGetInfo) MsgId() uint32 {
    return 7864
}

func (self *GS_GuildTechGetInfo_R) MsgId() uint32 {
    return 7865
}

func (self *C_GuildBossFight) MsgId() uint32 {
    return 7880
}

func (self *GS_GuildBossFight_R) MsgId() uint32 {
    return 7881
}

func (self *C_GuildBossGetCurrent) MsgId() uint32 {
    return 7882
}

func (self *GS_GuildBossGetCurrent_R) MsgId() uint32 {
    return 7883
}

func (self *C_GuildBossGetHistory) MsgId() uint32 {
    return 7884
}

func (self *GS_GuildBossGetHistory_R) MsgId() uint32 {
    return 7885
}

func (self *GS_MOpenModuleNew) MsgId() uint32 {
    return 8000
}

func (self *GS_CounterOpUpdate) MsgId() uint32 {
    return 8400
}

func (self *C_CounterRecover) MsgId() uint32 {
    return 8410
}

func (self *GS_CounterRecover_R) MsgId() uint32 {
    return 8411
}

func (self *C_CounterBuy) MsgId() uint32 {
    return 8412
}

func (self *GS_CounterBuy_R) MsgId() uint32 {
    return 8413
}

func (self *C_RankGet) MsgId() uint32 {
    return 8450
}

func (self *GS_RankGet_R) MsgId() uint32 {
    return 8451
}

func (self *C_RankLike) MsgId() uint32 {
    return 8452
}

func (self *GS_RankLike_R) MsgId() uint32 {
    return 8453
}

func (self *GS_TaskDailyValueChanged) MsgId() uint32 {
    return 8600
}

func (self *GS_TaskDailyItemCompleted) MsgId() uint32 {
    return 8601
}

func (self *C_TaskDailyInfo) MsgId() uint32 {
    return 8650
}

func (self *GS_TaskDailyInfo_R) MsgId() uint32 {
    return 8651
}

func (self *C_TaskDailyTakeBoxReward) MsgId() uint32 {
    return 8654
}

func (self *GS_TaskDailyTakeBoxReward_R) MsgId() uint32 {
    return 8655
}

func (self *C_TaskAchvTake) MsgId() uint32 {
    return 8700
}

func (self *GS_TaskAchvTake_R) MsgId() uint32 {
    return 8701
}

func (self *C_DrawGetInfo) MsgId() uint32 {
    return 8800
}

func (self *GS_DrawGetInfo_R) MsgId() uint32 {
    return 8801
}

func (self *C_DrawTp) MsgId() uint32 {
    return 8802
}

func (self *GS_DrawTp_R) MsgId() uint32 {
    return 8803
}

func (self *C_DrawScoreBoxTake) MsgId() uint32 {
    return 8804
}

func (self *GS_DrawScoreBoxTake_R) MsgId() uint32 {
    return 8805
}

func (self *C_WLevelFight) MsgId() uint32 {
    return 8900
}

func (self *GS_WLevelFight_R) MsgId() uint32 {
    return 8901
}

func (self *C_WLevelGJInfo) MsgId() uint32 {
    return 8902
}

func (self *GS_WLevelGJInfo_R) MsgId() uint32 {
    return 8903
}

func (self *C_WLevelGJTake) MsgId() uint32 {
    return 8904
}

func (self *GS_WLevelGJTake_R) MsgId() uint32 {
    return 8905
}

func (self *C_WLevelOneKeyGJTake) MsgId() uint32 {
    return 8906
}

func (self *GS_WLevelOneKeyGJTake_R) MsgId() uint32 {
    return 8907
}

func (self *C_WLevelFightOneKey) MsgId() uint32 {
    return 8908
}

func (self *GS_WLevelFightOneKey_R) MsgId() uint32 {
    return 8909
}

func (self *GS_AppointAddTask) MsgId() uint32 {
    return 9000
}

func (self *C_AppointCheckAdd) MsgId() uint32 {
    return 9010
}

func (self *GS_AppointCheckAdd_R) MsgId() uint32 {
    return 9011
}

func (self *C_AppointRefresh) MsgId() uint32 {
    return 9012
}

func (self *GS_AppointRefresh_R) MsgId() uint32 {
    return 9013
}

func (self *C_AppointLock) MsgId() uint32 {
    return 9014
}

func (self *GS_AppointLock_R) MsgId() uint32 {
    return 9015
}

func (self *C_AppointSend) MsgId() uint32 {
    return 9016
}

func (self *GS_AppointSend_R) MsgId() uint32 {
    return 9017
}

func (self *C_AppointAcc) MsgId() uint32 {
    return 9018
}

func (self *GS_AppointAcc_R) MsgId() uint32 {
    return 9019
}

func (self *C_AppointTake) MsgId() uint32 {
    return 9020
}

func (self *GS_AppointTake_R) MsgId() uint32 {
    return 9021
}

func (self *C_AppointCancel) MsgId() uint32 {
    return 9022
}

func (self *GS_AppointCancel_R) MsgId() uint32 {
    return 9023
}

func (self *C_TowerFight) MsgId() uint32 {
    return 9100
}

func (self *GS_TowerFight_R) MsgId() uint32 {
    return 9101
}

func (self *C_TowerRaid) MsgId() uint32 {
    return 9102
}

func (self *GS_TowerRaid_R) MsgId() uint32 {
    return 9103
}

func (self *C_TowerRecord) MsgId() uint32 {
    return 9104
}

func (self *GS_TowerRecord_R) MsgId() uint32 {
    return 9105
}

func (self *C_SetTeam) MsgId() uint32 {
    return 9200
}

func (self *GS_SetTeam_R) MsgId() uint32 {
    return 9201
}

func (self *GS_ArenaStageUpdate) MsgId() uint32 {
    return 9300
}

func (self *GS_ArenaFighted) MsgId() uint32 {
    return 9301
}

func (self *C_ArenaUpdateEnemy) MsgId() uint32 {
    return 9310
}

func (self *GS_ArenaUpdateEnemy_R) MsgId() uint32 {
    return 9311
}

func (self *C_ArenaFight) MsgId() uint32 {
    return 9312
}

func (self *GS_ArenaFight_R) MsgId() uint32 {
    return 9313
}

func (self *C_ArenaRecordInfo) MsgId() uint32 {
    return 9314
}

func (self *GS_ArenaRecordInfo_R) MsgId() uint32 {
    return 9315
}

func (self *C_ArenaRank) MsgId() uint32 {
    return 9316
}

func (self *GS_ArenaRank_R) MsgId() uint32 {
    return 9317
}

func (self *C_ArenaReplayGet) MsgId() uint32 {
    return 9318
}

func (self *GS_ArenaReplayGet_R) MsgId() uint32 {
    return 9319
}

func (self *C_ShopGetInfo) MsgId() uint32 {
    return 9400
}

func (self *GS_ShopGetInfo_R) MsgId() uint32 {
    return 9401
}

func (self *C_ShopBuy) MsgId() uint32 {
    return 9402
}

func (self *GS_ShopBuy_R) MsgId() uint32 {
    return 9403
}

func (self *C_ShopRefresh) MsgId() uint32 {
    return 9404
}

func (self *GS_ShopRefresh_R) MsgId() uint32 {
    return 9405
}

func (self *C_MarvelRollInfo) MsgId() uint32 {
    return 9500
}

func (self *GS_MarvelRollInfo_R) MsgId() uint32 {
    return 9501
}

func (self *C_MarvelRollRefresh) MsgId() uint32 {
    return 9502
}

func (self *GS_MarvelRollRefresh_R) MsgId() uint32 {
    return 9503
}

func (self *C_MarvelRollTake) MsgId() uint32 {
    return 9504
}

func (self *GS_MarvelRollTake_R) MsgId() uint32 {
    return 9505
}

func (self *GS_FriendNewApplied) MsgId() uint32 {
    return 9600
}

func (self *GS_FriendNewFriend) MsgId() uint32 {
    return 9601
}

func (self *GS_FriendRemoveFriend) MsgId() uint32 {
    return 9602
}

func (self *GS_FriendRecv) MsgId() uint32 {
    return 9603
}

func (self *C_FriendGetFrds) MsgId() uint32 {
    return 9610
}

func (self *GS_FriendGetFrds_R) MsgId() uint32 {
    return 9611
}

func (self *C_FriendSearchFrds) MsgId() uint32 {
    return 9612
}

func (self *GS_FriendSearchFrds_R) MsgId() uint32 {
    return 9613
}

func (self *C_FriendRemoveFrds) MsgId() uint32 {
    return 9614
}

func (self *GS_FriendRemoveFrds_R) MsgId() uint32 {
    return 9615
}

func (self *C_FriendGetApplyList) MsgId() uint32 {
    return 9616
}

func (self *GS_FriendGetApplyList_R) MsgId() uint32 {
    return 9617
}

func (self *C_FriendApply) MsgId() uint32 {
    return 9618
}

func (self *GS_FriendApply_R) MsgId() uint32 {
    return 9619
}

func (self *C_FriendAccept) MsgId() uint32 {
    return 9620
}

func (self *GS_FriendAccept_R) MsgId() uint32 {
    return 9621
}

func (self *C_FriendGetBlackList) MsgId() uint32 {
    return 9622
}

func (self *GS_FriendGetBlackList_R) MsgId() uint32 {
    return 9623
}

func (self *C_FriendAddBlackList) MsgId() uint32 {
    return 9624
}

func (self *GS_FriendAddBlackList_R) MsgId() uint32 {
    return 9625
}

func (self *C_FriendGiveAndRecv) MsgId() uint32 {
    return 9626
}

func (self *GS_FriendGiveAndRecv_R) MsgId() uint32 {
    return 9627
}

func (self *GS_CrusadeStageUpdate) MsgId() uint32 {
    return 9700
}

func (self *C_CrusadeGetInfo) MsgId() uint32 {
    return 9710
}

func (self *GS_CrusadeGetInfo_R) MsgId() uint32 {
    return 9711
}

func (self *C_CrusadeBoxTake) MsgId() uint32 {
    return 9712
}

func (self *GS_CrusadeBoxTake_R) MsgId() uint32 {
    return 9713
}

func (self *C_CrusadeFight) MsgId() uint32 {
    return 9714
}

func (self *GS_CrusadeFight_R) MsgId() uint32 {
    return 9715
}

func (self *GS_VipUpdate) MsgId() uint32 {
    return 9800
}

func (self *C_ActRushLocalGetInfo) MsgId() uint32 {
    return 10000
}

func (self *GS_ActRushLocalGetInfo_R) MsgId() uint32 {
    return 10001
}

func (self *C_ActRushLocalTake) MsgId() uint32 {
    return 10002
}

func (self *GS_ActRushLocalTake_R) MsgId() uint32 {
    return 10003
}

func (self *GS_ActBillLtTotal) MsgId() uint32 {
    return 10010
}

func (self *C_ActBillLtTotalInfo) MsgId() uint32 {
    return 10011
}

func (self *GS_ActBillLtTotalInfo_R) MsgId() uint32 {
    return 10012
}

func (self *C_ActBillLtTotalTake) MsgId() uint32 {
    return 10013
}

func (self *GS_ActBillLtTotalTake_R) MsgId() uint32 {
    return 10014
}

func (self *C_ActBillLtDayInfo) MsgId() uint32 {
    return 10020
}

func (self *GS_ActBillLtDayInfo_R) MsgId() uint32 {
    return 10021
}

func (self *C_ActBillLtDayTake) MsgId() uint32 {
    return 10022
}

func (self *GS_ActBillLtDayTake_R) MsgId() uint32 {
    return 10023
}

func (self *GS_ActGiftNew) MsgId() uint32 {
    return 10031
}

func (self *C_ActSummonInfo) MsgId() uint32 {
    return 10040
}

func (self *GS_ActSummonInfo_R) MsgId() uint32 {
    return 10041
}

func (self *C_ActSummonPick) MsgId() uint32 {
    return 10042
}

func (self *GS_ActSummonPick_R) MsgId() uint32 {
    return 10043
}

func (self *C_ActSummonDraw) MsgId() uint32 {
    return 10044
}

func (self *GS_ActSummonDraw_R) MsgId() uint32 {
    return 10045
}

func (self *GS_ActTargetTaskObjValueChanged) MsgId() uint32 {
    return 10050
}

func (self *C_ActTargetTaskInfo) MsgId() uint32 {
    return 10052
}

func (self *GS_ActTargetTaskInfo_R) MsgId() uint32 {
    return 10053
}

func (self *C_ActTargetTaskTake) MsgId() uint32 {
    return 10054
}

func (self *GS_ActTargetTaskTake_R) MsgId() uint32 {
    return 10055
}

func (self *C_ActHeroSkinInfo) MsgId() uint32 {
    return 10060
}

func (self *GS_ActHeroSkinInfo_R) MsgId() uint32 {
    return 10061
}

func (self *C_ActMSummonInfo) MsgId() uint32 {
    return 10070
}

func (self *GS_ActMSummonInfo_R) MsgId() uint32 {
    return 10071
}

func (self *C_ActMSummonDraw) MsgId() uint32 {
    return 10072
}

func (self *GS_ActMSummonDraw_R) MsgId() uint32 {
    return 10073
}

func (self *GS_ActMonopolyObjValueChanged) MsgId() uint32 {
    return 10080
}

func (self *GS_ActMonoPolyNextLv) MsgId() uint32 {
    return 10081
}

func (self *GS_ActMonopolyBoxReward) MsgId() uint32 {
    return 10082
}

func (self *C_ActMonopolyInfo) MsgId() uint32 {
    return 10090
}

func (self *GS_ActMonoPolyInfo_R) MsgId() uint32 {
    return 10091
}

func (self *C_ActMonopolyNPos) MsgId() uint32 {
    return 10092
}

func (self *GS_ActMonopolyNPos_R) MsgId() uint32 {
    return 10093
}

func (self *C_ActMonopolyAnswer) MsgId() uint32 {
    return 10094
}

func (self *GS_ActMonopolyAnswer_R) MsgId() uint32 {
    return 10095
}

func (self *C_ActMonopolyBuy) MsgId() uint32 {
    return 10096
}

func (self *GS_ActMonopolyBuy_R) MsgId() uint32 {
    return 10097
}

func (self *C_ActMonopolyBattle) MsgId() uint32 {
    return 10098
}

func (self *GS_ActMonopolyBattle_R) MsgId() uint32 {
    return 10099
}

func (self *C_ActMonopolyTaskTake) MsgId() uint32 {
    return 10100
}

func (self *GS_ActMonopolyTaskTake_R) MsgId() uint32 {
    return 10101
}

func (self *GS_ActMazeObjValueChanged) MsgId() uint32 {
    return 10110
}

func (self *GS_ActMazeTaskTaken) MsgId() uint32 {
    return 10111
}

func (self *C_ActMazeInfo) MsgId() uint32 {
    return 10115
}

func (self *GS_ActMazeInfo_R) MsgId() uint32 {
    return 10116
}

func (self *C_ActMazeClick) MsgId() uint32 {
    return 10117
}

func (self *GS_ActMazeClick_R) MsgId() uint32 {
    return 10118
}

func (self *C_ActMazeClickNext) MsgId() uint32 {
    return 10119
}

func (self *GS_ActMazeClickNext_R) MsgId() uint32 {
    return 10120
}

func (self *C_ActMazeClickTrade) MsgId() uint32 {
    return 10121
}

func (self *GS_ActMazeClickTrade_R) MsgId() uint32 {
    return 10122
}

func (self *C_ActMazeReset) MsgId() uint32 {
    return 10123
}

func (self *GS_ActMazeReset_R) MsgId() uint32 {
    return 10124
}

func (self *C_ActMazeClickThing) MsgId() uint32 {
    return 10125
}

func (self *GS_ActMazeClickThing_R) MsgId() uint32 {
    return 10126
}

func (self *C_ActMazeClickBattle) MsgId() uint32 {
    return 10127
}

func (self *GS_ActMazeClickBattle_R) MsgId() uint32 {
    return 10128
}

func (self *C_ActMazeTakeTask) MsgId() uint32 {
    return 10129
}

func (self *GS_ActMazeTakeTask_R) MsgId() uint32 {
    return 10130
}

func (self *C_ActMazeBuff) MsgId() uint32 {
    return 10131
}

func (self *GS_ActMazeBuff_R) MsgId() uint32 {
    return 10132
}

func (self *GS_ChatMsg) MsgId() uint32 {
    return 13000
}

func (self *C_ChatSend) MsgId() uint32 {
    return 13010
}

func (self *GS_ChatSend_R) MsgId() uint32 {
    return 13011
}

func (self *GS_MonthTicketValueChanged) MsgId() uint32 {
    return 13200
}

func (self *GS_MonthTicketItemCompleted) MsgId() uint32 {
    return 13201
}

func (self *GS_MonthTicketIsBuy) MsgId() uint32 {
    return 13202
}

func (self *C_MonthTicketInfo) MsgId() uint32 {
    return 13221
}

func (self *GS_MonthTicketInfo_R) MsgId() uint32 {
    return 13222
}

func (self *C_MonthTicketTakeOneKey) MsgId() uint32 {
    return 13223
}

func (self *GS_MonthTicketTakeOneKey_R) MsgId() uint32 {
    return 13224
}

func (self *C_MonthTicketBuyUp) MsgId() uint32 {
    return 13225
}

func (self *GS_MonthTicketBuyUp_R) MsgId() uint32 {
    return 13226
}

func (self *C_MonthTicketTake) MsgId() uint32 {
    return 13227
}

func (self *GS_MonthTicketTake_R) MsgId() uint32 {
    return 13228
}

func (self *C_MonthTicketTaskTake) MsgId() uint32 {
    return 13229
}

func (self *GS_MonthTicketTaskTake_R) MsgId() uint32 {
    return 13230
}

func (self *GS_PushGiftNew) MsgId() uint32 {
    return 13300
}

func (self *GS_PushGiftRewards) MsgId() uint32 {
    return 13301
}

func (self *C_PushGiftSetCreateTs) MsgId() uint32 {
    return 13310
}

func (self *GS_PushGiftSetCreateTs_R) MsgId() uint32 {
    return 13311
}

func (self *GS_GiftShopNew) MsgId() uint32 {
    return 13400
}

func (self *C_GiftShopTake) MsgId() uint32 {
    return 13410
}

func (self *GS_GiftShopTake_R) MsgId() uint32 {
    return 13411
}

func (self *GS_PrivCardNew) MsgId() uint32 {
    return 13500
}

func (self *C_PrivCardTake) MsgId() uint32 {
    return 13510
}

func (self *GS_PrivCardTake_R) MsgId() uint32 {
    return 13511
}

func (self *C_SignDailySign) MsgId() uint32 {
    return 13601
}

func (self *GS_SignDailySign_R) MsgId() uint32 {
    return 13602
}

func (self *GS_TaskMonthValueChanged) MsgId() uint32 {
    return 13700
}

func (self *GS_TaskMonthItemCompleted) MsgId() uint32 {
    return 13701
}

func (self *C_TaskMonthInfo) MsgId() uint32 {
    return 13702
}

func (self *GS_TaskMonthInfo_R) MsgId() uint32 {
    return 13703
}

func (self *C_TaskMonthTake) MsgId() uint32 {
    return 13704
}

func (self *GS_TaskMonthTask_R) MsgId() uint32 {
    return 13705
}

func (self *C_DaySignTake) MsgId() uint32 {
    return 13820
}

func (self *GS_DaySignTake_R) MsgId() uint32 {
    return 13821
}

func (self *C_TargetDaysTake) MsgId() uint32 {
    return 13900
}

func (self *GS_TargetDaysTake_R) MsgId() uint32 {
    return 13901
}

func (self *C_TargetDaysBuy) MsgId() uint32 {
    return 13902
}

func (self *GS_TargetDaysBuy_R) MsgId() uint32 {
    return 13903
}

func (self *C_TaskGrowTake) MsgId() uint32 {
    return 14030
}

func (self *GS_TaskGrowTake_R) MsgId() uint32 {
    return 14031
}

func (self *GS_WLevelFundNew) MsgId() uint32 {
    return 14000
}

func (self *C_WLevelFundTake) MsgId() uint32 {
    return 14010
}

func (self *GS_WLevelFundTake_R) MsgId() uint32 {
    return 14011
}

func (self *GS_GrowFundNew) MsgId() uint32 {
    return 14100
}

func (self *C_GrowFundInfo) MsgId() uint32 {
    return 14110
}

func (self *GS_GrowFundInfo_R) MsgId() uint32 {
    return 14111
}

func (self *C_GrowFundTakeLv) MsgId() uint32 {
    return 14112
}

func (self *GS_GrowFundTakeLv_R) MsgId() uint32 {
    return 14113
}

func (self *C_GrowFundTakeSvr) MsgId() uint32 {
    return 14114
}

func (self *GS_GrowFundTakeSvr_R) MsgId() uint32 {
    return 14115
}

func (self *GS_BillFirstNew) MsgId() uint32 {
    return 14201
}

func (self *C_BillFirstTake) MsgId() uint32 {
    return 14222
}

func (self *GS_BillFirstTake_R) MsgId() uint32 {
    return 14223
}

func (self *GS_LampMsg) MsgId() uint32 {
    return 14300
}

func (self *GS_GWarStageChange) MsgId() uint32 {
    return 14400
}

func (self *GS_GWarNewG2) MsgId() uint32 {
    return 14401
}

func (self *C_GWarGetSummary) MsgId() uint32 {
    return 14410
}

func (self *GS_GWarGetSummary_R) MsgId() uint32 {
    return 14411
}

func (self *C_GWarGetG2Members) MsgId() uint32 {
    return 14412
}

func (self *GS_GWarGetG2Members_R) MsgId() uint32 {
    return 14413
}

func (self *C_GWarGetGuildRank) MsgId() uint32 {
    return 14414
}

func (self *GS_GWarGetGuildRank_R) MsgId() uint32 {
    return 14415
}

func (self *C_GWarGetPlrRank) MsgId() uint32 {
    return 14416
}

func (self *GS_GWarGetPlrRank_R) MsgId() uint32 {
    return 14417
}

func (self *C_GWarFight) MsgId() uint32 {
    return 14418
}

func (self *GS_GWarFight_R) MsgId() uint32 {
    return 14419
}

func (self *GS_RiftMonsterNew) MsgId() uint32 {
    return 14500
}

func (self *GS_RiftMineNew) MsgId() uint32 {
    return 14501
}

func (self *GS_RiftMineOccupied) MsgId() uint32 {
    return 14502
}

func (self *GS_RiftBoxNew) MsgId() uint32 {
    return 14503
}

func (self *GS_RiftBoxRewards) MsgId() uint32 {
    return 14504
}

func (self *GS_RiftBoxOccupied) MsgId() uint32 {
    return 14505
}

func (self *C_RiftExplore) MsgId() uint32 {
    return 14560
}

func (self *GS_RiftExplore_R) MsgId() uint32 {
    return 14561
}

func (self *C_RiftMonsterFight) MsgId() uint32 {
    return 14562
}

func (self *GS_RiftMonsterFight_R) MsgId() uint32 {
    return 14563
}

func (self *C_RiftMineInfo) MsgId() uint32 {
    return 14564
}

func (self *GS_RiftMineInfo_R) MsgId() uint32 {
    return 14565
}

func (self *C_RiftMineOccupy) MsgId() uint32 {
    return 14566
}

func (self *GS_RiftMineOccupy_R) MsgId() uint32 {
    return 14567
}

func (self *C_RiftMineCancel) MsgId() uint32 {
    return 14568
}

func (self *GS_RiftMineCancel_R) MsgId() uint32 {
    return 14569
}

func (self *C_RiftMineTakeRewards) MsgId() uint32 {
    return 14570
}

func (self *GS_RiftMineTakeRewards_R) MsgId() uint32 {
    return 14571
}

func (self *C_RiftBoxOccupy) MsgId() uint32 {
    return 14572
}

func (self *GS_RiftBoxOccupy_R) MsgId() uint32 {
    return 14573
}

func (self *C_RiftBoxInfo) MsgId() uint32 {
    return 14574
}

func (self *GS_RiftBoxInfo_R) MsgId() uint32 {
    return 14575
}

func (self *GS_LadderStageChange) MsgId() uint32 {
    return 14600
}

func (self *C_LadderGetSummary) MsgId() uint32 {
    return 14610
}

func (self *GS_LadderGetSummary_R) MsgId() uint32 {
    return 14611
}

func (self *C_LadderMatch) MsgId() uint32 {
    return 14612
}

func (self *GS_LadderMatch_R) MsgId() uint32 {
    return 14613
}

func (self *C_LadderFight) MsgId() uint32 {
    return 14614
}

func (self *GS_LadderFight_R) MsgId() uint32 {
    return 14615
}

func (self *C_LadderGetRank) MsgId() uint32 {
    return 14616
}

func (self *GS_LadderGetRank_R) MsgId() uint32 {
    return 14617
}

func (self *C_LadderGetReplayList) MsgId() uint32 {
    return 14618
}

func (self *GS_LadderGetReplayList_R) MsgId() uint32 {
    return 14619
}

func (self *C_LadderGetReplay) MsgId() uint32 {
    return 14620
}

func (self *GS_LadderGetReplay_R) MsgId() uint32 {
    return 14621
}

func (self *C_HeroSkinAdd) MsgId() uint32 {
    return 14712
}

func (self *GS_HeroSkinAdd_R) MsgId() uint32 {
    return 14713
}

func (self *C_HeroSkinLvUp) MsgId() uint32 {
    return 14714
}

func (self *GS_HeroSkinLvUp_R) MsgId() uint32 {
    return 14715
}

func (self *C_WLevelDrawDraw) MsgId() uint32 {
    return 14800
}

func (self *GS_WLevelDrawDraw_R) MsgId() uint32 {
    return 14801
}

func (self *C_WLevelDrawTake) MsgId() uint32 {
    return 14802
}

func (self *GS_WLevelDrawTake_R) MsgId() uint32 {
    return 14803
}

func (self *GS_WarCupStageUpdate) MsgId() uint32 {
    return 14900
}

func (self *GS_WarCupChat) MsgId() uint32 {
    return 14901
}

func (self *GS_WarCupGuessRatio) MsgId() uint32 {
    return 14902
}

func (self *GS_WarCupAttainObjValueChanged) MsgId() uint32 {
    return 14903
}

func (self *C_WarCupGuessInfo) MsgId() uint32 {
    return 14920
}

func (self *GS_WarCupGuessInfo_R) MsgId() uint32 {
    return 14921
}

func (self *C_WarCupSelfInfo) MsgId() uint32 {
    return 14922
}

func (self *GS_WarCupSelfInfo_R) MsgId() uint32 {
    return 14923
}

func (self *C_WarCupTop64Info) MsgId() uint32 {
    return 14924
}

func (self *GS_WarCupTop64Info_R) MsgId() uint32 {
    return 14925
}

func (self *C_WarCupTop8Info) MsgId() uint32 {
    return 14926
}

func (self *GS_WarCupTop8Info_R) MsgId() uint32 {
    return 14927
}

func (self *C_WarCupGuess) MsgId() uint32 {
    return 14928
}

func (self *GS_WarCupGuess_R) MsgId() uint32 {
    return 14929
}

func (self *C_WarCupAuditionRank) MsgId() uint32 {
    return 14930
}

func (self *GS_WarCupAuditionRank_R) MsgId() uint32 {
    return 14931
}

func (self *C_WarCupChatSend) MsgId() uint32 {
    return 14932
}

func (self *GS_WarCupChatSend_R) MsgId() uint32 {
    return 14933
}

func (self *C_WarCupGetReplay) MsgId() uint32 {
    return 14934
}

func (self *GS_WarCupGetReplay_R) MsgId() uint32 {
    return 14935
}

func (self *C_WarCupTop1Info) MsgId() uint32 {
    return 14936
}

func (self *GS_WarCupTop1Info_R) MsgId() uint32 {
    return 14937
}

func (self *C_WarCupTaskTake) MsgId() uint32 {
    return 14938
}

func (self *GS_WarCupTaskTake_R) MsgId() uint32 {
    return 14939
}

func (self *C_WarCupWatch) MsgId() uint32 {
    return 14940
}

func (self *C_WarCupGuessRecords) MsgId() uint32 {
    return 14942
}

func (self *GS_WarCupGuessRecords_R) MsgId() uint32 {
    return 14943
}

func (self *GS_WBossStageChange) MsgId() uint32 {
    return 15000
}

func (self *C_WBossGetSummary) MsgId() uint32 {
    return 15010
}

func (self *GS_WBossGetSummary_R) MsgId() uint32 {
    return 15011
}

func (self *C_WBossFight) MsgId() uint32 {
    return 15012
}

func (self *GS_WBossFight_R) MsgId() uint32 {
    return 15013
}

func (self *C_WBossTakeMaxDmgRwd) MsgId() uint32 {
    return 15014
}

func (self *GS_WBossTakeMaxDmgRwd_R) MsgId() uint32 {
    return 15015
}

func (self *C_WBossGetRank) MsgId() uint32 {
    return 15016
}

func (self *GS_WBossGetRank_R) MsgId() uint32 {
    return 15017
}

func (self *C_WBossGetMaxDmgInfo) MsgId() uint32 {
    return 15018
}

func (self *GS_WBossGetMaxDmgInfo_R) MsgId() uint32 {
    return 15019
}

func (self *C_WBossGetSelfRank) MsgId() uint32 {
    return 15020
}

func (self *GS_WBossGetSelfRank_R) MsgId() uint32 {
    return 15021
}

func (self *C_InviteInfo) MsgId() uint32 {
    return 15130
}

func (self *GS_InviteInfo_R) MsgId() uint32 {
    return 15131
}

func (self *C_InviteTake) MsgId() uint32 {
    return 15132
}

func (self *GS_InviteTake_R) MsgId() uint32 {
    return 15133
}
