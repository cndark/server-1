package Err

// ============================================================================

const (
	OK     = 0 // 没问题
	Failed = 1 // 通用错误码(最好不弹窗)

	// --------------------------------
	// 通用 [2, 100)
	// --------------------------------

	Common_NoStatusHandler = 3  // 条件未发现
	Common_Timeout         = 5  // 超时
	Common_TimeNotUp       = 8  // 时间未到
	Common_BattleResError  = 9  // 战斗结异常
	Common_NotSetTeam      = 10 // 防守阵容为空
	Common_NeedMore        = 11 // 需要更多玩家
	Common_BattleNotFound  = 12 // 战斗没发现

	// --------------------------------
	// 登录 [100, 200)
	// --------------------------------

	Auth_InvalidVersion  = 100 // 版本不匹配
	Auth_Failed          = 101 // 认证失败
	Login_Failed         = 102 // 登录失败
	Login_UserSvr        = 103 // 玩家 服Id 出错
	Login_CloseReg       = 104 // 已关闭注册
	Login_CreateUserInfo = 105 // 创建玩家信息失败
	Login_UserBanned     = 106 // 玩家已被封号
	Login_GetUserInfo    = 107 // 获取玩家信息失败

	// --------------------------------
	// 玩家 [200, 300)
	// --------------------------------

	Plr_NameLen            = 200 // 昵称长度不符，请重新输入
	Plr_InvalidName        = 201 // 该昵称违规，请重新输入
	Plr_SameName           = 202 // 该昵称不能使用，请重新输入
	Plr_Offline            = 204 // 玩家已离线
	Plr_NotLoad            = 205 // 玩家未登陆
	Plr_DupName            = 206 // 名字被占用
	Plr_ModuleLocked       = 208 // 模块未解锁
	Plr_LowLevel           = 209 // 玩家等级不足
	Plr_WLevelLow          = 210 // 主线关卡不足
	Plr_NotFound           = 213 // 玩家没发现
	Plr_NotNew             = 215 // 不是新号
	Plr_FullLevel          = 216 // 满级了
	Plr_AttainTabNotEnough = 222 // 条件不足
	Plr_LowVipLevel        = 223 // VIP等级不足
	Plr_TeamInvalid        = 224 // 上阵阵容有误
	Plr_TakenBefore        = 225 // 奖励已经领取
	Plr_IsSelf             = 226 // 是自己
	Plr_HFrameError        = 227 // 头像框不可用
	Plr_SelfOpLimit        = 228 // 不能操作自己
	Plr_TeamAllDeath       = 229 // 队伍全死
	Plr_PrivCardInValid    = 230 // 特权卡没激活
	Plr_PrivCardNotFound   = 231 // 特权不存在
	Plr_NotPay             = 232 // 没有付费
	Plr_CondLimited        = 233 // 条件没达到
	Plr_InTeam             = 234 // 在防守阵容里面

	// --------------------------------
	// 背包 [300, 400)
	// --------------------------------

	Bag_Full           = 300 // 背包已满
	Bag_NotEnoughHero  = 301 // 您的英雄不足
	Bag_NotEnoughRelic = 302 // 您的神器不足

	// --------------------------------
	// 英雄养成 [500, 600)
	// --------------------------------

	Hero_NotFound               = 500 // 英雄没有找到
	Hero_FullLevel              = 501 // 英雄等级已满
	Hero_MaxLevelReached        = 502 // 英雄当前最大等级已到
	Hero_LowLevel               = 503 // 英雄等级不够
	Hero_FullStar               = 504 // 英雄星级已满
	Hero_InvalidCostHero        = 505 // 消耗的英雄不正确
	Hero_Locked                 = 506 // 英雄已经锁定
	Hero_TrinketNotUnlocked     = 507 // 饰品没有解锁
	Hero_TrinketAlreadyUnlocked = 508 // 饰品已经解锁
	Hero_TrinketFullLevel       = 509 // 饰品满级
	Hero_LowStar                = 510 // 英雄星级不够
	Hero_TrinketNoPropToCommit  = 511 // 饰品没有可保存的属性
	Hero_SkinExist              = 512 // 英雄皮肤已经存在
	Hero_SkinNotExist           = 513 // 英雄皮肤不存在
	Hero_SkinLvFull             = 514 // 英雄皮肤等级满

	// --------------------------------
	// 可装备容器 [600, 650)
	// --------------------------------

	Equip_NotFound        = 600 // 没有找到要装备的物件
	Equip_AlreadyEquipped = 601 // 已被装备
	Equip_NotEquipped     = 602 // 没有被装备
	Equip_InvalidSlot     = 603 // 槽位信息错误
	Equip_SlotIsEmpty     = 604 // 装备位为空

	// --------------------------------
	// 装备 [650, 700)
	// --------------------------------

	Armor_Uncomposable = 650 // 该装备不能合成

	// --------------------------------
	// 神器 [700, 750)
	// --------------------------------

	Relic_FullStar = 700 // 神器星级已满

	// --------------------------------
	// 计数器 [1000, 1100)
	// --------------------------------

	Counter_NotFound  = 1000 // 计数器没有找到
	Counter_CanNotBuy = 1001 // 计数器不能购买
	Counter_Enough    = 1002 // 计数器还有次数

	// --------------------------------
	// 邮件 [1200, 1300)
	// --------------------------------

	Mail_NoMail       = 1200 // 没有邮件
	Mail_NotFound     = 1201 // 没找到邮件
	Mail_AlreadyRead  = 1202 // 已读
	Mail_AlreadyTaken = 1203 // 已领取
	Mail_NoAttachment = 1204 // 没有附件

	// --------------------------------
	// 公会 [1600, 1900)
	// --------------------------------

	Guild_NameTooLong         = 1600 // 家族名称过长，请重新输入
	Guild_InvalidName         = 1601 // 无法命名该名称，请重新输入
	Guild_SameName            = 1602 // 此命名重复，请重新输入
	Guild_DupName             = 1603 // 名字被占用
	Guild_InvalidHead         = 1604 // 该家族头像暂不能使用
	Guild_HeadUnavailable     = 1605 // 该家族头像还不可以使用
	Guild_NoticeTooLong       = 1606 // 输入家族宣言过长，请修改家族宣言
	Guild_InvalidNotice       = 1607 // 宣言中有敏感词汇，请重新命名
	Guild_PlrInGuild          = 1608 // 玩家已经在家族内
	Guild_PlrLowLevel         = 1609 // 玩家等级过低
	Guild_AlreadyApplied      = 1610 // 家族重复申请
	Guild_FullPlrApply        = 1611 // 家族申请列表中已满
	Guild_FullMember          = 1612 // 家族内玩家已经满了
	Guild_FullVice            = 1613 // 家族内副会长人数已满
	Guild_FullElite           = 1614 // 家族内副会长人数已满
	Guild_JoinDenied          = 1615 // 申请加入已经被拒绝
	Guild_InvalidApplyMode    = 1616 // 应用模式无效
	Guild_OperateOnSelf       = 1617 // 本人无法操作
	Guild_InvalidRank         = 1618 // 排名无效
	Guild_NotApplied          = 1619 // 尚未申请
	Guild_NotAMember          = 1620 // 不是家族成员
	Guild_OwnerCantLeave      = 1621 // 家族会长不能脱离家族
	Guild_NotMeetCreateCond   = 1622 // 忍阶不足3阶，不能创建家族
	Guild_NotFound            = 1623 // 家族没有找到
	Guild_LowPriv             = 1624 // 家族成员权限不足
	Guild_JoinTimeLimited     = 1647 // 家族加入时间限制
	Guild_Locked              = 1648 // 家族没开
	Guild_KickOwnerTimeError  = 1661 // 族长弹劾时间没到
	Guild_LeaveTimeLimited    = 1662 // 这个时间段内不能离开家族
	Guild_DonateSameDay       = 1663 // 今天已经捐献过了
	Guild_LowLevel            = 1679 // 家族等级不足
	Guild_TodayLeave          = 1685 // 今天离开过家族
	Guild_SignAlready         = 1686 // 今天已签到
	Guild_WishCD              = 1687 // 许愿cd
	Guild_WishItemCnt         = 1688 // 许愿物品次数上限
	Guild_WishLimit           = 1689 // 许愿条目上限
	Guild_WishNotFound        = 1690 // 许愿条目不存在
	Guild_WishFullHelp        = 1691 // 许愿条目助力已满
	Guild_WishNotYours        = 1692 // 许愿条目不是你的
	Guild_WishNotFullHelp     = 1693 // 许愿条目助力不够
	Guild_OrderCd             = 1694 // 订单cd
	Guild_OrderLimit          = 1695 // 订单条目上限
	Guild_OrderNotFound       = 1696 // 订单没有找到
	Guild_OrderAlreadyStarted = 1697 // 订单已经启动
	Guild_OrderFullStar       = 1698 // 订单满星
	Guild_OrderNotEnd         = 1699 // 订单为结束
	Guild_TechPreCond         = 1700 // 科技没有满足前置条件
	Guild_TechFullLevel       = 1701 // 科技满级
	Guild_BossHistNotFound    = 1702 // 副本 boss 历史未找到
	Guild_ZmCd                = 1703 // 招募cd中

	// --------------------------------
	// 新手 [1900, 2000)
	// --------------------------------

	Tut_KeyAlreadyUsed = 1900 // 激活码已经被使用了哟，请再换一个吧
	Tut_KeyNotFound    = 1901 // 激活码没有找到，请再确认一下呢？

	// --------------------------------
	// 云数据 [2000, 2100)
	// --------------------------------

	Cld_InvalidKey = 2000 // 该激活码已经存在

	// --------------------------------
	// 活动管理器 [2200, 2300)
	// --------------------------------

	Act_ActNotFound        = 2200 // 活动没有找到
	Act_ActPlrDataNotFound = 2201 // 玩家活动数据未找到
	Act_ActSvrDataNotFound = 2202 // 活动的数据未找到
	Act_ActClosed          = 2203 // 活动已关闭
	Act_ConfGrp            = 2204 // 活动confgrp不对
	Act_StageError         = 2205 // 活动阶段有误

	// --------------------------------
	// 礼包码 [2700, 2800)
	// --------------------------------

	Gift_NoCode      = 2700 // 兑换码不存在
	Gift_AreaLimit   = 2701 // 区域限定
	Gift_CodeExpired = 2702 // 兑换码已过期
	Gift_CodeUsed    = 2703 // 兑换码已使用

	// --------------------------------
	// 充值 [2900, 3000)
	// --------------------------------

	Bill_NoCard                 = 2900 // 没有该月卡
	Bill_CardRewardAlreadyTaken = 2901 // 月卡奖励今日已领
	Bill_RefundCodeError        = 2910 // 返利码不存在
	Bill_RefundCodeTaken        = 2911 // 返利码已经领取
	Bill_RefundUpdateFailed     = 2912 // 返利码更新领取失败
	Bill_NotEnough              = 2913 // 充值不足
	Bill_NotPay                 = 2914 // 没有充值

	// --------------------------------
	// 开服庆典任务(7日目标) [3800, 3900)
	// --------------------------------

	TargetDays_Closed       = 3800 // 开服庆典已经结束
	TargetDays_NotCompleted = 3801 // 开服庆典任务还未完成
	TargetDays_BuyCntMax    = 3802 // 开服庆典购买最大限制

	// --------------------------------
	// 日常任务 [4100, 4200)
	// --------------------------------

	TaskDaily_NotFound        = 4100 // 任务没有找到
	TaskDaily_NotCompleted    = 4101 // 任务还没有完成
	TaskDaily_AlreadRewarded  = 4102 // 任务奖励已经领取
	TaskDaily_ActiveNotEnough = 4103 // 任务积分不足

	// --------------------------------
	// 成就任务 [4200, 4300)
	// --------------------------------

	TaskAchv_NotFound       = 4200 // 成就没有找到
	TaskAchv_NotCompleted   = 4201 // 成就还没有完成
	TaskAchv_AlreadRewarded = 4202 // 成就奖励已经领取

	// --------------------------------
	// 抽卡 [4300, 4400)
	// --------------------------------

	Draw_TpNotFound             = 4300 // tp没发现
	Draw_ScoreBoxTaken          = 4301 // 宝箱奖励已领取
	Draw_ScoreBoxScoreNotEnough = 4302 // 宝箱奖励积分不足

	// --------------------------------
	// 推图 [4400, 4500)
	// --------------------------------

	WLevel_GJEmpty    = 4400 // 挂机背包为空
	WLevel_LowLvNum   = 4401 // 通过关卡不足
	WLevel_DrawBefore = 4402 // 推图十连-已经抽过
	WLevel_DrawNull   = 4403 // 推图十连-无掉落

	// --------------------------------
	// 酒馆派遣 [4500, 4600)
	// --------------------------------

	Appoint_AllLock         = 4500 // 酒馆--任务全部锁定
	Appoint_TaskNotFound    = 4501 // 酒馆--任务没发现
	Appoint_TaskSending     = 4502 // 酒馆--任务正在进行
	Appoint_HeroSending     = 4503 // 酒馆--英雄正在进行
	Appoint_HeroCondLimit   = 4504 // 酒馆--英雄条件不足
	Appoint_TooManyHero     = 4505 // 酒馆--英雄过多
	Appoint_TaskNotSend     = 4506 // 酒馆--任务没有开始
	Appoint_TaskFinished    = 4507 // 酒馆--任务已经完成
	Appoint_TaskNotFinished = 4508 // 酒馆--任务未完成
	Appoint_TaskCntLimit    = 4509 // 酒馆--任务数量限制

	// --------------------------------
	// 竞技场 [4600, 4700)
	// --------------------------------

	Arena_EnemyNotFound  = 4600 // 竞技场-没有发现对手
	Arena_PlayerNotFound = 4601 // 竞技场-对手有误
	Arena_RevengeError   = 4602 // 竞技场-冤冤相报何时了

	// --------------------------------
	// 商店 [4700, 4800)
	// --------------------------------

	Shop_NotFound        = 4700 // 商城没有找到
	Shop_ItemIdError     = 4701 // 商品没有找到
	Shop_BuyCntTop       = 4702 // 购买次数达到上限
	Shop_PriceError      = 4703 // 价格错误
	Shop_RefreshCntError = 4704 // 刷新次数错误
	Shop_RefreshCntFull  = 4705 // 刷新次数到达上限
	Shop_BuyBlankLock    = 4706 // 购买的格子被锁定
	Shop_NotOpen         = 4707 // 商城没有开启

	// --------------------------------
	// 奇迹之盘 [4800, 4900)
	// --------------------------------

	MarvelRoll_GroupNotFound = 4800 // 分组没找到
	MarvelRoll_RollNothing   = 4801 // 没随机到东西

	// --------------------------------
	// 好友 [4900, 5000)
	// --------------------------------

	Friend_BlackList     = 4950 // 黑名单
	Friend_ApplyDup      = 4951 // 重复申请
	Friend_IsFriend      = 4952 // 已经是好友
	Friend_IsNotFriend   = 4953 // 非好友
	Friend_NotApply      = 4954 // 不在申请列表
	Friend_Full          = 4955 // 好友满了
	Friend_BlackListFull = 4956 // 黑名单满了
	Friend_ApplyListFull = 4957 // 对方申请好友已满

	// --------------------------------
	// 远征 [5000, 5100)
	// --------------------------------

	Crusade_PassAll       = 5000 // 已经通关
	Crusade_EnemyNotFound = 5001 // 对手没找到
	Crusade_NotPass       = 5002 // 没有通过
	Crusade_End           = 5003 // 已经结束

	// --------------------------------
	// 活动 [5200, 7000)
	// --------------------------------

	Activity_TimeLimit         = 5200 // 活动时间限制
	Activity_TakeBefore        = 5201 // 活动奖励已经领取
	Activity_NotJoin           = 5202 // 活动未参与
	Activity_RankNotFound      = 5203 // 榜单未发现
	Activity_BillNotEnough     = 5204 // 充值不足
	Activity_CondLimited       = 5205 // 条件不足
	Activity_BuyBefore         = 5206 // 已经购买
	Activity_NotBuy            = 5207 // 没有购买
	Activity_BillDayNotEnough  = 5208 // 充值天数不足
	Activity_SummonHeroNotPick = 5209 // 主题召唤没设置心愿英雄
	Activity_SummonDiamLimit   = 5210 // 主题召唤钻石消耗上限
	// Activity_SummonNoPick      = 5211 // 主题召唤未选择英雄
	// Activity_SummonDiamLimit   = 5212 // 主题召唤钻石消耗上限

	Activity_MonopolyLvMax      = 5213 // 关卡达到上限或配置为空
	Activity_MonopolyPNotFound  = 5214 // 没有对应的答题奇遇
	Activity_MonopolySNotFound  = 5215 // 没有对应的商店折扣奇遇
	Activity_MonopolyExpired    = 5216 // 奇遇过期
	Activity_MonopolyTpNotFound = 5217 // 没有该奇遇类型

	Activity_MazeAlreadyClicked = 5218 // 迷宫点已经点击过了
	Activity_MazeBuyCntLimit    = 5219 // 体力购买次数限制
	Activity_MazeTradeLimit     = 5220 // 商人兑换上限
	Activity_MazeAlreadyTaken   = 5221 // 已领取
	Activity_MazeNotClicked     = 5222 // 格子没有点击
	Activity_MazePosLimit       = 5223 // 位置为(0~25]
	Activity_MazeNoBattle       = 5224 // 该点没有Battle

	// --------------------------------
	// 聊天 [7000, 7100)
	// --------------------------------

	Chat_TypeNotFound      = 7000 // 聊天类型有误
	Chat_ContentLenLimited = 7001 // 聊天内容长度有误

	// --------------------------------
	// 月票 [7100, 7200)
	// --------------------------------

	MonthTicket_NotFound        = 7100 // 任务没有找到
	MonthTicket_NotCompleted    = 7101 // 任务还没有完成
	MonthTicket_AlreadRewarded  = 7102 // 任务奖励已经领取
	MonthTicket_ActiveNotEnough = 7103 // 任务积分不足

	// --------------------------------
	// 签到 [7200,7300)
	// --------------------------------

	SignDaily_AlreadRewarded = 7200 // 已经签到过

	// -------------------------------
	// 每月任务 [7300,7400)
	// -------------------------------

	TaskMonth_AlreadRewarded = 7300 // 任务奖励已经领取
	TaskMonth_NotCompleted   = 7301 // 任务还没有完成
	TaskMonth_NotFound       = 7302 // 任务没有找到

	// -------------------------------
	// 在线宝箱 [7400,7500)
	// -------------------------------

	OnlineBox_NotCompleted   = 7400 // 在线时长不足
	OnlineBox_AlreadRewarded = 7401 // 奖励已经领取

	// -------------------------------
	// 七日之约 [7500,7600)
	// -------------------------------

	DaySign_AlreadSigned   = 7500 // 已经签过了
	DaySign_NotSigned      = 7501 // 没有签到
	DaySign_AlreadRewarded = 7502 // 奖励已经领取
	DaySign_NotFound       = 7503 // 签到没有找到
	DaySign_NotCond        = 7504 // 签到条件不足

	// -------------------------------
	// 榜单玩法 [7600,7700)
	// -------------------------------

	RankPlay_RankNotReady    = 7600 // 榜单未准备好
	RankPlay_PlayerNotOnList = 7601 // 玩家没有上榜

	// -------------------------------
	// 进阶之路 [7700,7800)
	// -------------------------------

	TaskGrow_NotFound       = 7700 // 成就没有找到
	TaskGrow_NotCompleted   = 7701 // 成就还没有完成
	TaskGrow_AlreadRewarded = 7702 // 成就奖励已经领取

	// -------------------------------
	// 超值首充 [7800,7900)
	// -------------------------------

	BillFirst_NotBill        = 7801 // 没有充值
	BillFirst_NotCond        = 7802 // 领取未到时间
	BillFirst_AlreadRewarded = 7803 // 奖励已经领取
	BillFirst_NotFound       = 7804 // 没有找到

	// -------------------------------
	// 工会战 [7900,8000)
	// -------------------------------

	GWar_NotEnrolled    = 7900 // 没有入围
	GWar_NoMatch        = 7901 // 没有匹配
	GWar_TargetNotFound = 7902 // 目标没有找到
	GWar_TargetDead     = 7903 // 目标已经死亡
	GWar_TargetLocked   = 7904 // 目标锁定中

	// -------------------------------
	// 裂隙 [8000,8100)
	// -------------------------------

	Rift_MonsterNotFound   = 8000 // 怪物未发现
	Rift_MineExploreBefore = 8010 // 请先探索矿
	Rift_MineOccupyNoMore  = 8011 // 该类型的矿不能占领
	Rift_MineNotFound      = 8012 // 矿没找到
	Rift_MineFin           = 8013 // 矿已经完成了
	Rift_MineNotFin        = 8014 // 矿时间未到
	Rift_MineIsSelf        = 8015 // 自己的矿不能抢
	Rift_MineNotSelf       = 8015 // 不是自己的矿
	Rift_BoxNotFound       = 8016 // 宝箱没找到
	Rift_BoxIsSelf         = 8017 // 自己的宝箱不能抢
	Rift_BoxIsBattle       = 8018 // 宝箱正在抢夺
	Rift_BoxIsEnd          = 8019 // 宝箱已经结束
	Rift_MineIsBattle      = 8020 // 矿正在抢夺

	// -------------------------------
	// 天梯 [8100,8200)
	// -------------------------------

	Ladder_TargetNotFound = 8100 // 目标未找到
	Ladder_InFight        = 8101 // 自己或目标正在战斗中 (被打)
	Ladder_TargetChanged  = 8102 // 目标已改变

	// -------------------------------
	// 本服杯赛 [8200,8300)
	// -------------------------------

	WarCup_NotInGuess  = 8200 // 当前不是竞猜期
	WarCup_GuessBefore = 8201 // 已经竞猜过

	// -------------------------------
	// wboss [8300,8400)
	// -------------------------------

	WBoss_MaxDmgRwdTaken      = 8300 // 最大伤害奖已经领取
	WBoss_MaxDmgRwdNotAllowed = 8301 // 最大伤害奖没有满足要求
)

// ============================================================================

// 返回错误码：货币，道具，装备，counter等不足
func NotEnoughObject(objid int32) int32 {
	return 100_000_000 + int32(objid)

}
