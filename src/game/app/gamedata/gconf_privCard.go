package gamedata

var ConfPrivCard = &privCardTable{}

type privCard struct {
	Id               int32 `json:"id"` // 特权卡
	CounterTimeLimit []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"counterTimeLimit"` // 有效期内次数加成
	CounterOnce []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"counterOnce"` // 只生效1次次数永久加成
	Reward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"reward"` // 奖励
	DailyReward []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"dailyReward"` // 每日奖励
	ExpireDay     int32 `json:"expireDay"`     // 有效期（天）
	ExtNeedAddCnt int32 `json:"extNeedAddCnt"` // 获得几次该卡之后额外奖励生效
	ExtDay        int32 `json:"extDay"`        // 额外天数
	ExtReward     []*struct {
		Id int32 `json:"id"`
		N  int64 `json:"n"`
	} `json:"extReward"` // 获得额外奖励
	ExtCounter []*struct {
		Id int32 `json:"id"`
		N  int32 `json:"n"`
	} `json:"extCounter"` // 获得额外特权
	DisPayId int32 `json:"disPayId"` // 折扣充值id
	PayId    int32 `json:"payId"`    // 充值id
}

type privCardTable struct {
	items map[int32]*privCard
}

func (self *privCardTable) Load() {
	var arr []*privCard
	if !load_json("privCard.json", &arr) {
		return
	}

	items := make(map[int32]*privCard)

	for _, v := range arr {
		items[v.Id] = v
	}

	self.items = items
}

func (self *privCardTable) Query(id int32) *privCard {
	return self.items[id]
}

func (self *privCardTable) Items() map[int32]*privCard {
	return self.items
}
