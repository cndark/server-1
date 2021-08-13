package gamedata

var ConfAppoint = &appointTable{}

type appoint struct {
	Id   int32 `json:"id"`   // id
	Star int32 `json:"star"` // 星级
	Cost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"cost"` // 探索消耗
	Weight  int32 `json:"weight"` // 任务权重
	AccCost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"accCost"` // 加速消耗消耗
	HeroNum int32 `json:"heroNum"` // 最大英雄数
	Time    int32 `json:"time"`    // 探索时间(秒)
	Reward  []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"reward"` // 奖励
	Require []*struct {
		Type int32 `json:"type"`
		N    int32 `json:"n"`
	} `json:"require"` // 探索要求
	HeroElemRange []int32 `json:"heroElemRange"` // 英雄元素阵营范围
}

type appointTable struct {
	items map[int32]*appoint
}

func (self *appointTable) Load() {
	var arr []*appoint
	if !load_json("appoint.json", &arr) {
		return
	}

	items := make(map[int32]*appoint)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *appointTable) Query(id int32) *appoint {
	return self.items[id]
}

func (self *appointTable) Items() map[int32]*appoint {
	return self.items
}
