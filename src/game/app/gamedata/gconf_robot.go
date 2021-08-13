package gamedata

var ConfRobot = &robotTable{}

type robot struct {
	GroupId           int32   `json:"groupId"`           // bot组ID
	Lv                int32   `json:"lv"`                // bot强度
	MaxLv             int32   `json:"maxLv"`             // bot强度上限
	Head              string  `json:"head"`              // 头像
	HFrame            int32   `json:"hFrame"`            // 头像框
	HFigure           int32   `json:"hFigure"`           // 形象
	Power             int32   `json:"power"`             // 机器人战力
	ArenaInitialScore int32   `json:"arenaInitialScore"` // 机器人积分
	Num               int32   `json:"num"`               // 生成数量
	Monster1          []int32 `json:"monster1"`          // 第一位模型随机
	Monster2          []int32 `json:"monster2"`          // 第二位模型随机
	Monster3          []int32 `json:"monster3"`          // 第三位模型随机
	Monster4          []int32 `json:"monster4"`          // 第四位模型随机
	Monster5          []int32 `json:"monster5"`          // 第五位模型随机
	Monster6          []int32 `json:"monster6"`          // 第六位模型随机
}

type robotTable struct {
	items map[int32]*robot
}

func (self *robotTable) Load() {
	var arr []*robot
	if !load_json("robot.json", &arr) {
		return
	}

	items := make(map[int32]*robot)

	for _, v := range arr {
		items[v.GroupId] = v
	}

	self.items = items
}

func (self *robotTable) Query(groupId int32) *robot {
	return self.items[groupId]
}

func (self *robotTable) Items() map[int32]*robot {
	return self.items
}
