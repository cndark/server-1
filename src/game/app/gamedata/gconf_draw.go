package gamedata

var ConfDraw = &drawTable{}

type draw struct {
	DrawGroup string `json:"drawGroup"` // 抽奖组
	ModuleId  int32  `json:"moduleId"`  // 功能板块，openId
	Cost      []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"cost"` // 基础单价
	EquilCost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"equilCost"` // 等价物单价
	Discount         int32   `json:"discount"`         // 10连折扣
	CostFreeTime     int32   `json:"costFreeTime"`     // 免费时限(秒)
	InitialTimes     int32   `json:"initialTimes"`     // 初始保底次数
	InitialOdds      float32 `json:"initialOdds"`      // 初始保底几率
	InitialGroup     string  `json:"initialGroup"`     // 保底组
	InitialBaseGroup string  `json:"InitialBaseGroup"` // 非5星保底组
	FixedTimes       []int32 `json:"fixedTimes"`       // 次数必出
	FixedGroup       string  `json:"fixedGroup"`       // 保底组
	ExtraTimes       int32   `json:"extraTimes"`       // 保底组次数
	ExtraGroup       string  `json:"extraGroup"`       // 保底组
	Score            int32   `json:"score"`            // 单次获得积分
}

type drawTable struct {
	items map[string]*draw
}

func (self *drawTable) Load() {
	var arr []*draw
	if !load_json("draw.json", &arr) {
		return
	}

	items := make(map[string]*draw)

	for _, v := range arr {
		items[v.DrawGroup] = v
	}

	self.items = items
}

func (self *drawTable) Query(drawGroup string) *draw {
	return self.items[drawGroup]
}

func (self *drawTable) Items() map[string]*draw {
	return self.items
}
