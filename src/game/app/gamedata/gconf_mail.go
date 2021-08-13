package gamedata

var ConfMail = &mailTable{}

type mail struct {
	Key      int32  `json:"key"`  // 邮件key
	Skey     string `json:"skey"` // 邮件skey
	MailItem []*struct {
		Id int32 `json:"Id"`
		N  int64 `json:"n"`
	} `json:"mailItem"` // 邮件附件
}

type mailTable struct {
	items map[int32]*mail
}

func (self *mailTable) Load() {
	var arr []*mail
	if !load_json("mail.json", &arr) {
		return
	}

	items := make(map[int32]*mail)

	for _, v := range arr {
		items[v.Key] = v
	}

	self.items = items
}

func (self *mailTable) Query(key int32) *mail {
	return self.items[key]
}

func (self *mailTable) Items() map[int32]*mail {
	return self.items
}
