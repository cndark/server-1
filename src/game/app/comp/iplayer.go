package comp

import (
	"fw/src/core/db"
	"fw/src/game/msg"
	"time"
)

type IPlayer interface {
	GetAuthId() string
	GetId() string
	GetName() string
	GetSvr0() string
	GetSdk() string
	GetModel() string
	GetDevId() string
	GetOs() string
	GetOsVer() string
	GetLoginIP() string
	GetCreateTs() time.Time
	GetHead() string
	GetHFrame() int32
	GetLevel() int32
	GetAtkPwr() int32
	IsOnline() bool
	GetLoginTs() time.Time
	GetOnlineDur() int32
	GetOfflineTs() time.Time
	GetLoginSumDays() int32
	IsActive() bool
	GetAuthRet() map[string]string
	GetBanTs() time.Time
	IsBan() bool
	AccOnlineDur() int32
	AddExp(v int32)
	SetHFrame(id int32)

	GetVipLevel() int32
	AddVipExp(v int32)

	DB() *db.Database
	SendMsg(message msg.Message)

	GetBag() *Bag
	GetCounter() *Counter

	GetBillTotalBaseCcy() int64

	GetGuildId() string
	GetGuildName() string
	GetGuildRank() int32

	GetAttainObjVal(oid int32) float64
	IsModuleOpen(mid int32) bool
	GetActRawData(name string) interface{}

	IsSetTeam(tp int32) bool
	GetTeam(tp int32) *msg.TeamFormation
	GetTeamAtkPwr(tf *msg.TeamFormation) int32
	IsTeamFormationValid(tf *msg.TeamFormation) bool

	ToMsg_SimpleInfo(atkpwr ...int32) *msg.PlayerSimpleInfo
	ToMsg_BattleTeam(tf *msg.TeamFormation, teamAtkPwr ...bool) *msg.BattleTeam

	IsFriend(plrid string) bool

	IsPrivCardValid(id int32) bool
	AddPrivCard(id int32)
}
