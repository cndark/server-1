package gamedata

var ConfLimitM = &limitTableM{}

type limit struct {
	MaxPlrLv    int32
	MaxWLevelLv int32
	MaxSvrGrpId int32
}

type limitTableM struct {
	limit limit
}

func (self *limitTableM) Load() {
	// plr lv
	for _, v := range ConfPlayerUp.items {
		if self.limit.MaxPlrLv == 0 {
			self.limit.MaxPlrLv = v.Level
		} else if self.limit.MaxPlrLv < v.Level {
			self.limit.MaxPlrLv = v.Level
		}
	}

	// world max lv
	for _, v := range ConfWorldLevel1000.items {
		if self.limit.MaxWLevelLv == 0 {
			self.limit.MaxWLevelLv = v.Id
		} else if self.limit.MaxWLevelLv < v.Id {
			self.limit.MaxWLevelLv = v.Id
		}
	}

	// max svrgrp id
	for _, v := range ConfCrossconf.items {
		if v.Id > self.limit.MaxSvrGrpId {
			self.limit.MaxSvrGrpId = v.Id
		}
	}
}

func (self *limitTableM) Query() *limit {
	return &self.limit
}
