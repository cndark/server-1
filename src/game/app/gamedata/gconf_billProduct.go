package gamedata

var ConfBillProduct = &billProductTable{}

type billProduct struct {
	Id       int32  `json:"id"`       // 充值ID
	Sdk      string `json:"sdk"`      // 平台ID
	TypeId   int32  `json:"typeId"`   // 类型ID
	PayId    int32  `json:"payId"`    // 支付ID
	Txt_Name string `json:"txt_Name"` // 名字
	BaseCcy  int32  `json:"baseCcy"`  // 基准货币值
	Price    int32  `json:"price"`    // 售价
	Ccy      string `json:"ccy"`      // 货币类型
	Goods    []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"goods"` // 获得物品
	VipExp    int32 `json:"vipExp"` // 获得vip经验
	ExtGoods1 []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"extGoods1"` // 首次额外获得物品
	ExtGoods2 []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"extGoods2"` // 多次购买奖励物品
	ResetDaily int32 `json:"resetDaily"` // 是否每日重置
	MailId     int32 `json:"mailId"`     // 邮件ID
}

type billProductTable struct {
	items map[int32]*billProduct
}

func (self *billProductTable) Load() {
	var arr []*billProduct
	if !load_json("billProduct.json", &arr) {
		return
	}

	items := make(map[int32]*billProduct)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *billProductTable) Query(id int32) *billProduct {
	return self.items[id]
}

func (self *billProductTable) Items() map[int32]*billProduct {
	return self.items
}
