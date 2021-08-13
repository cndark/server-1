package gamedata

var ConfGlobalPublic = &globalPublicTable{}

type globalPublic struct {
	StructureNum     int32  `json:"structureNum"`  // 结构序列
	HeroBagLimit     int32  `json:"heroBagLimit"`  // 卡牌背包上限
	SysMailSender    string `json:"sysMailSender"` // 系统邮件发送者
	RefundSwitch     int32  `json:"refundSwitch"`  // 返利码开关
	PlayerNameChange []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"playerNameChange"` // 改名消耗
	PlayerHeadCost []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"playerHeadCost"` // 玩家头像更换消耗
	DrawScore []*struct {
		Vip int32 `json:"vip"`
		N   int32 `json:"n"`
		Id  int32 `json:"id"`
	} `json:"drawScore"` // 召唤积分进度及奖励
	ExploreTimeLimit  int32 `json:"exploreTimeLimit"`  // 挂机最大时长，小时
	ExploreOnekeyTime int32 `json:"exploreOnekeyTime"` // 快速挂机收益时间，小时
	HeroResetCost     []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"heroResetCost"` // 英雄重置等级消耗
	AppointRefreshCost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"appointRefreshCost"` // 酒馆刷新消耗单价
	AppointAddTaskTime int32   `json:"appointAddTaskTime"` // 酒馆新增1个任务时间,分钟
	HeroJobRange       []int32 `json:"heroJobRange"`       // 英雄职业范围
	GuildCreateCost    []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"guildCreateCost"` // 公会创建消耗
	GuildNameChange []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"guildNameChange"` // 公会改名消耗
	GuildJoinCD []int32 `json:"guildJoinCD"` // 公会加入冷却（分钟）
	GuildWish   []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"guildWish"` // 公会祈愿列表
	GuildDisMail  int32 `json:"guildDisMail"`  // 公会解散后给该公会全体玩家发的邮件
	GuildKickHour int32 `json:"guildKickHour"` // 公会弹劾需要满足的小时数
	GuildKickCost []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"guildKickCost"` // 公会弹劾需要满足的消耗
	GuildExitMail     int32    `json:"guildExitMail"`     // 逐出公会消息
	NotExitGuild      []string `json:"notExitGuild"`      // 不允许离开公会的时间
	GuildWishLimit    int32    `json:"guildWishLimit"`    // 公会祈愿上限
	GuildWishCd       int32    `json:"guildWishCd"`       // 公会许愿冷却(分钟)
	GuildOrderLimit   int32    `json:"guildOrderLimit"`   // 公会订单列表上限
	HarborOrderCd     int32    `json:"harborOrderCd"`     // 港口订单获取冷却(小时)
	HarborDonateValue int32    `json:"harborDonateValue"` // 公会港口单次捐赠量
	GuildHelpReward   []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"guildHelpReward"` // 助力奖励
	CommonFragment     []int32 `json:"commonFragment"`     // 升星返还通用碎片(禁止修改)
	FriendShieldLimit  int32   `json:"friendShieldLimit"`  // 好友屏蔽数量
	FriendApplyListMax int32   `json:"friendApplyListMax"` // 好友申请列表上限
	FriendPointGiveMax int32   `json:"friendPointGiveMax"` // 友情点单日赠送上限
	FriendPointAddMax  int32   `json:"friendPointAddMax"`  // 友情点单日获得上限
	CrusadeLevelMin    int32   `json:"crusadeLevelMin"`    // 英灵试炼最低等级要求
	BuyGoldCritAdd     []*struct {
		Min float32 `json:"min"`
		Max float32 `json:"max"`
	} `json:"buyGoldCritAdd"` // 点金暴击率增加随机区间
	BuyGoldBasic        int32 `json:"buyGoldBasic"`        // 点金初始金币
	BuyGoldLevelRatio   int32 `json:"buyGoldLevelRatio"`   // 点金等级系数
	BuyGoldLimit        int32 `json:"buyGoldLimit"`        // 点金上限
	CrusadeRewardMailId int32 `json:"crusadeRewardMailId"` // 英灵试炼奖励补发邮件id
	VipUpRewardMailId   int32 `json:"VipUpRewardMailId"`   // VIP奖励邮件
	ChatCrsSvrCost      []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"chatCrsSvrCost"` // 跨服聊天消耗
	MonthTicketTaskNum int32 `json:"monthTicketTaskNum"` // 月票任务随机个数
	RankLikeReward     []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"rankLikeReward"` // 排行榜点赞奖励
	TargetDaysDur int32 `json:"targetDaysDur"` // 开服庆典持续时间（天）
	ArenaMatching []*struct {
		A float32 `json:"a"`
		B float32 `json:"b"`
	} `json:"arenaMatching"` // 竞技场积分匹配范围
	HeroChangeCost []*struct {
		Star int32 `json:"star"`
		Id   int32 `json:"id"`
		N    int32 `json:"n"`
	} `json:"heroChangeCost"` // 英雄置换消耗
	HeroStarInheritCost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"heroStarInheritCost"` // 英雄星级继承消耗
	TutorialDraw       []int32 `json:"tutorialDraw"`       // 新手引导10连奖励
	TutorialDrawNormal int32   `json:"tutorialDrawNormal"` // 新手引导1次基础召唤
	TutorialDrawSenior int32   `json:"tutorialDrawSenior"` // 新手引导1次高级召唤
	TutorialGj         []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"tutorialGj"` // 新手引导挂机奖励
	GuildWarKillScore  int32   `json:"guildWarKillScore"`  // 公会战强制击败值
	GuildWarKillShield int32   `json:"guildWarKillShield"` // 公会战强制击败护盾
	RiftEvent          []int32 `json:"riftEvent"`          // 裂隙事件权重
	RiftCost           []int32 `json:"riftCost"`           // 裂隙事件消耗
	RiftMineLv         []*struct {
		Min float32 `json:"min"`
		Max float32 `json:"max"`
	} `json:"riftMineLv"` // 裂隙矿产等级范围
	RiftMineFastCost int32 `json:"riftMineFastCost"` // 裂隙矿产快速领取分钟钻石
	RiftMineRobMail  int32 `json:"riftMineRobMail"`  // 裂隙矿产被抢邮件
	RiftMonsterLv    []*struct {
		Min float32 `json:"min"`
		Max float32 `json:"max"`
	} `json:"riftMonsterLv"` // 裂隙怪物等级范围
	RiftBoxRefresh    int32 `json:"riftBoxRefresh"`    // 裂隙宝箱刷新时间
	RiftBoxNum        int32 `json:"riftBoxNum"`        // 裂隙宝箱刷新个数
	RiftBoxLifeMinute int32 `json:"riftBoxLifeMinute"` // 裂隙宝箱持续时间（分钟）
	RiftBoxOpenMinute int32 `json:"riftBoxOpenMinute"` // 裂隙宝箱开启时间（分钟）
	RiftBoxReward     []*struct {
		Id   int32   `json:"id"`
		N    int32   `json:"n"`
		Odds float32 `json:"odds"`
	} `json:"riftBoxReward"` // 裂隙宝箱开启基础奖励
	RiftBoxMail      int32     `json:"riftBoxMail"`      // 裂隙宝箱奖励邮件
	RiftMineNorepeat int32     `json:"riftMineNorepeat"` // 裂隙矿去重次数
	SkinLvLimit      int32     `json:"skinLvLimit"`      // 皮肤等级上限
	LadderMatch      []float64 `json:"ladderMatch"`      // 天梯向上匹配参数
	LadderRobot      []*struct {
		Lv int32 `json:"lv"`
		N  int32 `json:"n"`
	} `json:"ladderRobot"` // 天梯机器人
	LadderMatchNum       int32   `json:"ladderMatchNum"`       // 天梯匹配对手数量
	LadderRewardMailId   int32   `json:"ladderRewardMailId"`   // 天梯奖励邮件id
	RiftFirstMonster     int32   `json:"riftFirstMonster"`     // 空间裂隙第一次探索怪物
	HeroResetFreeLv      int32   `json:"heroResetFreeLv"`      // 免费英雄重生等级上限
	WarCupLossRatio      float32 `json:"warCupLossRatio"`      // 杯赛竞猜赔付率
	WarCupGuessBaseOdds  float32 `json:"warCupGuessBaseOdds"`  // 杯赛竞猜基准赔率
	WarCupGuessOddsLimit float32 `json:"warCupGuessOddsLimit"` // 杯赛竞猜赔率上限
	WarCupGuessFloors    float32 `json:"warCupGuessFloors"`    // 杯赛竞猜赔率保底
	WarCupAFreePoints    []int32 `json:"warCupAFreePoints"`    // 杯赛海选赛免费投放竞猜积分
	WarCupKFreePoints    []int32 `json:"warCupKFreePoints"`    // 杯赛淘汰赛免费投放竞猜积分
	WarCupFFreePoints    []int32 `json:"warCupFFreePoints"`    // 杯赛决赛赛免费投放竞猜积分
	WarCupGuessLimit     int32   `json:"warCupGuessLimit"`     // 单场竞猜参与竞猜积分上限
	WarCupCoinRatio      float32 `json:"warCupCoinRatio"`      // 杯赛竞猜积分兑换系数
	WarCupRankMailId     int32   `json:"warCupRankMailId"`     // 杯赛排名奖励id
	WarCupGuessMailId    int32   `json:"warCupGuessMailId"`    // 杯赛竞猜积分兑换奖励id
	WarCupWinBaseScore   int32   `json:"warCupWinBaseScore"`   // 杯赛胜利基础积分
	WarCupWinScoreP1     int32   `json:"warCupWinScoreP1"`     // 杯赛积分计算参数1
	WarCupWinScoreP2     float32 `json:"warCupWinScoreP2"`     // 杯赛积分计算参数2
	WarCupWinFloorsRatio float32 `json:"warCupWinFloorsRatio"` // 杯赛胜利积分保底系数
	WarCupForcastMail    []*struct {
		MailId int32 `json:"mailId"`
		Sec    int64 `json:"sec"`
	} `json:"warCupForcastMail"` // 杯赛预告邮件
	WorldBossArrange     []int32 `json:"worldBossArrange"`     // 世界boss期数boss安排
	WorldBossRankMailId  int32   `json:"worldBossRankMailId"`  // 世界boss排名奖励id
	WorldBossPointFloors int32   `json:"worldBossPointFloors"` // boss伤害保底积分
	WorldBossPointRatio  float32 `json:"worldBossPointRatio"`  // boss伤害换算系数
}

type globalPublicTable struct {
	items map[int32]*globalPublic
}

func (self *globalPublicTable) Load() {
	var arr []*globalPublic
	if !load_json("globalPublic.json", &arr) {
		return
	}

	items := make(map[int32]*globalPublic)

	for _, v := range arr {
		items[v.StructureNum] = v
	}

	self.items = items
}

func (self *globalPublicTable) Query(structureNum int32) *globalPublic {
	return self.items[structureNum]
}

func (self *globalPublicTable) Items() map[int32]*globalPublic {
	return self.items
}
