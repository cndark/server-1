package gconst

// ============================================================================
// 异步操作队列类型

const (
	AQ_Default = 0
	AQ_Mail    = 1
	AQ_Bill    = 2
	AQ_Stats   = 3
	AQ_GLog    = 4
)

// ============================================================================
// 活跃玩家判定

const (
	PLAYER_ActiveDays = 10
)

// ============================================================================
// 对象类型

const (
	ObjType_Currency = 1 // 普通货币
	ObjType_Item     = 2 // 物品
	ObjType_Hero     = 3 // 英雄
	ObjType_Relic    = 4 // 神器
)

// ============================================================================

func ObjectType(id int32) int32 {
	return id / 10000
}

func IsCurrency(id int32) bool {
	return ObjectType(id) == ObjType_Currency
}

func IsItem(id int32) bool {
	return ObjectType(id) == ObjType_Item
}

func IsHero(id int32) bool {
	return ObjectType(id) == ObjType_Hero
}

func IsRelic(id int32) bool {
	return ObjectType(id) == ObjType_Relic
}

// ============================================================================
// 特殊货币

const (
	Gold      = 10001 // 金币
	Diamond   = 10002 // 钻石
	PlayerExp = 10003 // 玩家经验
	ArenaCoin = 10011 // 竞技币
	VipExp    = 10015 // Vip经验
	DrawHigh  = 20002 // 高级召唤卷

	ArenaCard = 20011 // 入场券
)

// ============================================================================
// 道具类型

const (
	ItemType_Armor = 40 // 装备
)

func IsArmor(tp int32) bool {
	return tp/10*10 == ItemType_Armor
}

func ArmorSlot(tp int32) int32 {
	return tp % 10
}

// ============================================================================
// 可穿戴

const (
	EquipGroup_Armor           = 1
	EquipGroup_Armor_SlotStart = 0
	EquipGroup_Armor_SlotEnd   = 3

	EquipGroup_Relic           = 2
	EquipGroup_Relic_SlotStart = 4
	EquipGroup_Relic_SlotEnd   = 4
)

// ============================================================================
// 充值发货类型

const (
	Bill_Normal      = 10 // 普通
	Bill_PrivCard    = 20 // 特权卡
	Bill_First       = 30 // 首冲
	Bill_Fund        = 40 // 基金
	Bill_Gift        = 50 // 礼包
	Bill_RushCross   = 60 // 冲榜
	Bill_MonthTicket = 70 // 月票
)

// cs充值透传csext=type_confid (下划线分割, 1_101)
const (
	Bill_CsExt_Type_GiftShop = "1" // 礼包充值商店
	Bill_CsExt_Type_PushGift = "2" // 推送礼包
	Bill_CsExt_Type_ActGift  = "3" // 活动礼包
)

const (
	Bill_PayId_SaleWeek    = 2001 // 折扣周卡
	Bill_PayId_SaleMonth   = 2003 // 折扣月卡
	Bill_PayId_GrowFund    = 4001 // 成长基金
	Bill_PayId_WLevelFund  = 4002 // 推图基金
	Bill_PayId_MonthTicket = 7128 // 月票

)

// ============================================================================
// battle

const (
	Team_MaxSlots = 6
)

const (
	TeamType_Dfd    = 1 // 普通防守阵容
	TeamType_GldFt  = 2 // 家族战防守阵容
	TeamType_WarCup = 3 // 杯赛防守阵容
	TeamType_Max    = 4 // 最大值
)

// ============================================================================
// attainid

const (
	AttainId_MaxAtkPwr = 1000
)

// ============================================================================
// 聊天

const (
	C_ChatType_World  = 1 // 世界
	C_ChatType_Cross  = 2 // 跨服
	C_ChatType_Guild  = 3 // 家族
	C_ChatType_Friend = 4 // 好友
	C_ChatType_GldZm  = 5 // 公会招募

	C_ChatExpireTs = 72 // 聊天消息有效期小时
)

// ============================================================================
// 特权卡

const (
	C_PrivCard_Month    = 2001 // 月卡
	C_PrivCard_LongLife = 3001 // 终身卡
)

// ============================================================================
// 活动礼包类型

const (
	C_ActGift_Plr = 1 // 个人礼包
	C_ActGift_Gld = 2 // 家族礼包
)
