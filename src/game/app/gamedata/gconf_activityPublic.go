package gamedata

var ConfActivityPublic = &activityPublicTable{}

type activityPublic struct {
	StructureNum       int32 `json:"structureNum"`    // 结构序列
	ActBillLtDayMin    int32 `json:"actBillLtDayMin"` // 累天充值最低额度
	RushLocalDrawScore []*struct {
		Tp string `json:"tp"`
		Sc int32  `json:"sc"`
	} `json:"rushLocalDrawScore"` // 召唤所得冲榜积分
	RushLocalMarvelRollScore []*struct {
		Tp string `json:"tp"`
		Sc int32  `json:"sc"`
	} `json:"rushLocalMarvelRollScore"` // 轮盘祈愿所得冲榜积分
	RushLocalArenaScore []int32 `json:"rushLocalArenaScore"` // 竞技场挑战所得冲榜积分
	MonopolyConsNormal  []*struct {
		Cg int32 `json:"cg"`
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"monopolyConsNormal"` // 大富翁活动普通消耗道具
	MonopolyConsRandom []*struct {
		Cg int32 `json:"cg"`
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"monopolyConsRandom"` // 大富翁活动随机骰子消耗道具
	MonopolyAdvSplitNum int32   `json:"monopolyAdvSplitNum"` // 大富翁奇遇点位间隔数
	MonopolyBattleRatio float32 `json:"monopolyBattleRatio"` // 大富翁怪物强度系数
	MazePowerItem       []*struct {
		Id    int32 `json:"id"`
		Power int32 `json:"power"`
	} `json:"mazePowerItem"` // 迷宫体力道具回复
	MazeBloodItem []*struct {
		Id    int32 `json:"id"`
		Blood int32 `json:"blood"`
	} `json:"mazeBloodItem"` // 迷宫生命道具回复
	MazeReviveCost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"mazeReviveCost"` // 迷宫复活消耗
	MazeLoseBlood int32 `json:"mazeLoseBlood"` // 战斗失败扣除生命
	MazeTrapBlood int32 `json:"mazeTrapBlood"` // 迷宫陷阱扣除生命
	MazePowerFree int32 `json:"mazePowerFree"` // 迷宫体力增加值
	MazeKeyReward []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"mazeKeyReward"` // 钥匙事件奖励
	MazeMail      int32   `json:"mazeMail"`      // 迷宫奖励邮件
	MazeItemClear []int32 `json:"mazeItemClear"` // 迷宫结束清除道具
}

type activityPublicTable struct {
	items map[int32]*activityPublic
}

func (self *activityPublicTable) Load() {
	var arr []*activityPublic
	if !load_json("activityPublic.json", &arr) {
		return
	}

	items := make(map[int32]*activityPublic)

	for _, v := range arr {
		items[v.StructureNum] = v
	}

	self.items = items
}

func (self *activityPublicTable) Query(structureNum int32) *activityPublic {
	return self.items[structureNum]
}

func (self *activityPublicTable) Items() map[int32]*activityPublic {
	return self.items
}
