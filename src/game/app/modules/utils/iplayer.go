package utils

import (
	"fw/src/game/msg"
)

// ============================================================================

var iplrmgr IPlayerMgr

// ============================================================================

type IPlayer interface {
	SendMsg(message msg.Message)
}

type IPlayerMgr interface {
	LoadIPlayer(uid string) interface{}
	FindIPlayerById(uid string) interface{}
	FindIPlayerByName(name string) interface{}
	ForEachLoadedIPlayer(f func(plr interface{}))
	ForEachOnlineIPlayer(f func(plr interface{}))
	ForEachLoadedIPlayerBreakable(f func(plr interface{}) bool)
	ForEachOnlineIPlayerBreakable(f func(plr interface{}) bool)

	NumLoaded() int32
	NumOnline() int32
}

// ============================================================================

func LoadPlayer(uid string) interface{} {
	plr := iplrmgr.LoadIPlayer(uid)
	if plr == nil {
		return nil
	} else {
		return plr
	}
}

func FindPlayer(uid string) interface{} {
	plr := iplrmgr.FindIPlayerById(uid)
	if plr == nil {
		return nil
	} else {
		return plr
	}
}

func FindPlayerByName(name string) interface{} {
	plr := iplrmgr.FindIPlayerByName(name)
	if plr == nil {
		return nil
	} else {
		return plr
	}
}

func ForEachLoadedPlayer(f func(plr interface{})) {
	iplrmgr.ForEachLoadedIPlayer(f)
}

func ForEachOnlinePlayer(f func(plr interface{})) {
	iplrmgr.ForEachOnlineIPlayer(f)
}

func ForEachLoadedPlayerBreakable(f func(plr interface{}) bool) {
	iplrmgr.ForEachLoadedIPlayerBreakable(f)
}

func ForEachOnlinePlayerBreakable(f func(plr interface{}) bool) {
	iplrmgr.ForEachOnlineIPlayerBreakable(f)
}

func PlayerNumLoaded() int32 {
	return iplrmgr.NumLoaded()
}

func PlayerNumOnline() int32 {
	return iplrmgr.NumOnline()
}

func Send2Player(uid string, message msg.Message) {
	plr := iplrmgr.FindIPlayerById(uid)
	if plr != nil {
		plr.(IPlayer).SendMsg(message)
	}
}

func BroadcastPlayers(message msg.Message) {
	iplrmgr.ForEachOnlineIPlayer(func(plr interface{}) {
		plr.(IPlayer).SendMsg(message)
	})
}
