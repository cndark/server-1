package pt

import (
	"fw/src/bot/msg"
	"fw/src/core/evtmgr"
	"math/rand"
	"sync"
)

// ============================================================================

var (
	v_plrlist      []string
	v_plrlist_lock sync.RWMutex
)

// ============================================================================

func init() {
	evtmgr.On("userinfo", func(args ...interface{}) {
		info := args[1].(*msg.GS_UserInfo)

		v_plrlist_lock.Lock()
		v_plrlist = append(v_plrlist, info.UserId)
		v_plrlist_lock.Unlock()
	})
}

func global_rand_plrid() string {
	v_plrlist_lock.RLock()
	defer v_plrlist_lock.RUnlock()

	L := len(v_plrlist)
	if L == 0 {
		return ""
	}

	return v_plrlist[rand.Intn(L)]
}
