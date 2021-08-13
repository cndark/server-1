package cloud

import (
	"strings"
)

// ============================================================================

type Cloud struct {
	Data map[string]string

	plr IPlayer
}

// ============================================================================

func NewCloud() *Cloud {
	return &Cloud{
		Data: make(map[string]string),
	}
}

func (self *Cloud) Init(plr IPlayer) {
	self.plr = plr

	if self.Data == nil {
		self.Data = make(map[string]string)
	}
}

func (self *Cloud) Get(key string) string {
	return self.Data[key]
}

func (self *Cloud) Set(key, val string) bool {
	// check key
	if key == "" || strings.Contains(key, ".") {
		return false
	}

	if val == "" {
		// remove
		delete(self.Data, key)
	} else {
		// update
		self.Data[key] = val
	}

	return true
}
