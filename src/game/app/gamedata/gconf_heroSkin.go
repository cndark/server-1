package gamedata

var ConfHeroSkin = &heroSkinTable{}

type heroSkin struct {
	Id   int32 `json:"id"`   // id
	Hero int32 `json:"hero"` // 关联英雄
	Cost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"cost"` // 激活消耗
	LvUpCost []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"lvUpCost"` // 升级消耗
	BaseProps []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"baseProps"` // 属性
	BasePropGrowth []*struct {
		Id  int32   `json:"id"`
		Val float32 `json:"val"`
	} `json:"basePropGrowth"` // 属性增长
}

type heroSkinTable struct {
	items map[int32]*heroSkin
}

func (self *heroSkinTable) Load() {
	var arr []*heroSkin
	if !load_json("heroSkin.json", &arr) {
		return
	}

	items := make(map[int32]*heroSkin)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *heroSkinTable) Query(id int32) *heroSkin {
	return self.items[id]
}

func (self *heroSkinTable) Items() map[int32]*heroSkin {
	return self.items
}
