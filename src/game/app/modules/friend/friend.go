package friend

import (
	"fw/src/core/evtmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/chat"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"time"
)

// ============================================================================
// 好友
type Friend struct {
	Friends   friend_m        // 好友列表
	ApplyList friend_m        // 申请列表
	BlackList friend_m        // 黑名单
	GiveList  map[string]bool // 赠送列表
	RecvList  map[string]bool // 获赠列表

	plr IPlayer
}

type friend_m map[string]*friend_t
type friend_t struct {
	Ts time.Time // 添加时间

	Chat []*chat.ChatOne
}

// ============================================================================

func init() {
	// reset daily
	evtmgr.On(gconst.Evt_PlrResetDaily, func(args ...interface{}) {
		plr := args[0].(IPlayer)
		plr.GetFriend().reset_daily()
	})
}

func NewFriend() *Friend {
	return &Friend{
		Friends:   make(friend_m),
		ApplyList: make(friend_m),
		BlackList: make(friend_m),
		GiveList:  make(map[string]bool),
		RecvList:  make(map[string]bool),
	}
}

// ============================================================================

func (self *Friend) Init(plr IPlayer) {
	self.plr = plr
}

func (self *Friend) reset_daily() {
	self.GiveList = make(map[string]bool)
	self.RecvList = make(map[string]bool)
}

// 是否被搜索
func (self *Friend) IsSearch(plrid string) bool {
	if self.plr.GetId() == plrid {
		return false
	}

	_, ok := self.Friends[plrid]
	if ok {
		return false
	}

	_, ok = self.ApplyList[plrid]
	if ok {
		return false
	}

	_, ok = self.BlackList[plrid]
	if ok {
		return false
	}

	return true
}

func (self *Friend) IsFriend(plrid string) bool {
	_, ok := self.Friends[plrid]
	if ok {
		return true
	}

	return false
}

// ============================================================================
// 申请好友
func (self *Friend) Apply(plrid string) int32 {
	if plrid == self.plr.GetId() {
		return Err.Plr_SelfOpLimit
	}

	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g == nil {
		return Err.Failed
	}

	// check black list
	_, ok := self.BlackList[plrid]
	if ok {
		return Err.Friend_BlackList
	}

	_, ok = self.Friends[plrid]
	if ok {
		return Err.Friend_IsFriend
	}

	fplr := find_player(plrid)
	if fplr == nil {
		return Err.Plr_NotLoad
	}

	if len(fplr.GetFriend().ApplyList) >= int(conf_g.FriendApplyListMax) {
		return Err.Friend_ApplyListFull
	}

	_, ok = fplr.GetFriend().BlackList[self.plr.GetId()]
	if ok {
		return Err.Friend_BlackList
	}

	_, ok = fplr.GetFriend().ApplyList[self.plr.GetId()]
	if ok {
		return Err.Friend_ApplyDup
	}

	fplr.GetFriend().ApplyList[self.plr.GetId()] = &friend_t{Ts: time.Now()}

	// notify
	fplr.SendMsg(&msg.GS_FriendNewApplied{
		PlrId: self.plr.GetId(),
	})

	return Err.OK
}

// 接受好友
func (self *Friend) Accept(plrid string, accept bool) int32 {
	if plrid == self.plr.GetId() {
		return Err.Plr_SelfOpLimit
	}

	conf_p := gamedata.ConfPlayerUp.Query(self.plr.GetLevel())
	if conf_p == nil {
		return Err.Failed
	}

	_, ok := self.ApplyList[plrid]
	if !ok {
		return Err.Friend_NotApply
	}

	if !accept {
		delete(self.ApplyList, plrid)
		return Err.OK
	}

	_, ok = self.Friends[plrid]
	if ok {
		return Err.OK
	}

	if len(self.Friends) >= int(conf_p.FriendLimit) {
		return Err.Friend_Full
	}

	delete(self.ApplyList, plrid)

	// check black list
	_, ok = self.BlackList[plrid]
	if ok {
		return Err.Friend_BlackList
	}

	// add friend
	fplr := find_player(plrid)
	if fplr == nil {
		return Err.Plr_NotLoad
	}

	_, ok = fplr.GetFriend().BlackList[self.plr.GetId()]
	if ok {
		return Err.Friend_BlackList
	}

	// 机器人好友上限
	if fplr.GetSdk() == "soda.ai" {
		if len(fplr.GetFriend().Friends) >= 2000 {
			return Err.Friend_Full
		}
	} else {
		if len(fplr.GetFriend().Friends) >= int(conf_p.FriendLimit) {
			return Err.Friend_Full
		}
	}

	now := time.Now()
	self.Friends[plrid] = &friend_t{
		Ts: now,
	}

	fplr.GetFriend().Friends[self.plr.GetId()] = &friend_t{
		Ts: now,
	}

	// notify
	self.plr.SendMsg(&msg.GS_FriendNewFriend{
		Frd: self.plr.GetFriend().FriendInfo_ToMsg(fplr.GetId()),
	})

	fplr.SendMsg(&msg.GS_FriendNewFriend{
		Frd: fplr.GetFriend().FriendInfo_ToMsg(self.plr.GetId()),
	})

	evtmgr.Fire(gconst.Evt_FriendAdd, self.plr)
	evtmgr.Fire(gconst.Evt_FriendAdd, fplr)

	return Err.OK
}

// 移除好友
func (self *Friend) Remove(plrid string) int32 {
	if plrid == self.plr.GetId() {
		return Err.Plr_SelfOpLimit
	}

	_, ok := self.Friends[plrid]
	if !ok {
		return Err.Friend_IsNotFriend
	}

	fplr := load_player(plrid)
	if fplr != nil {
		delete(fplr.GetFriend().Friends, self.plr.GetId())
	}

	delete(self.Friends, plrid)

	// notify
	fplr.SendMsg(&msg.GS_FriendRemoveFriend{
		PlrId: self.plr.GetId(),
	})

	return Err.OK
}

// 添加黑名单
func (self *Friend) AddBlackList(plrid string) int32 {
	if plrid == self.plr.GetId() {
		return Err.Plr_SelfOpLimit
	}

	// check black list
	_, ok := self.BlackList[plrid]
	if ok {
		return Err.Friend_BlackList
	}

	conf := gamedata.ConfGlobalPublic.Query(1)
	if conf == nil {
		return Err.Failed
	}

	if len(self.BlackList) >= int(conf.FriendShieldLimit) {
		return Err.Friend_BlackListFull
	}

	_, ok = self.Friends[plrid]
	if ok {
		self.Remove(plrid)
	}

	_, ok = self.ApplyList[plrid]
	if ok {
		delete(self.ApplyList, plrid)
	}

	self.BlackList[plrid] = &friend_t{Ts: time.Now()}

	return Err.OK
}

// 移除黑名单
func (self *Friend) RemoveBlackList(plrid string) int32 {
	_, ok := self.BlackList[plrid]
	if ok {
		delete(self.BlackList, plrid)
	}

	return Err.OK
}

// 送礼和收礼
func (self *Friend) GiveAndRecv(plrids []string) (cnt int32) {
	cop := self.plr.GetCounter().NewOp(gconst.ObjFrom_FriendGive)

	n := int32(0)
	for _, plrid := range plrids {
		fplr := load_player(plrid)
		if fplr == nil {
			continue
		}

		_, ok := self.Friends[plrid]
		if !ok {
			continue
		}

		if self.RecvList[plrid] {
			cop.IncCounter(gconst.Cnt_PlayerStrength, 1)
			delete(self.RecvList, plrid)
			cnt++
		}

		if self.GiveList[plrid] {
			continue
		}

		self.GiveList[plrid] = true
		fplr.GetFriend().RecvList[self.plr.GetId()] = true

		n++

		// notify
		fplr.SendMsg(&msg.GS_FriendRecv{
			PlrId: self.plr.GetId(),
		})
	}

	cop.Apply()

	if n > 0 {
		evtmgr.Fire(gconst.Evt_FriendGive, self.plr, n)
	}

	return
}

func (self *Friend) AddChat(c *chat.ChatOne) {
	fid := c.FromId
	if self.plr.GetId() == c.FromId {
		fid = c.ToId
	}

	frd, ok := self.Friends[fid]
	if !ok {
		return
	}

	frd.Chat = append(frd.Chat, c)

	if l := len(frd.Chat); l > 10 {
		frd.Chat = frd.Chat[l-10:]
	}

	// notify
	fplr := load_player(c.FromId)
	if fplr != nil {
		self.plr.SendMsg(&msg.GS_ChatMsg{
			One: &msg.ChatOne{
				Tp:      c.Tp,
				From:    fplr.ToMsg_SimpleInfo(),
				ToId:    c.ToId,
				Content: c.Content,
				Ts:      time.Now().Unix(),
			}})
	}
}

// ============================================================================

func (self *Friend) ToMsg() *msg.FriendData {
	ret := &msg.FriendData{}

	for plrid := range self.Friends {
		m := self.FriendInfo_ToMsg(plrid)
		if m != nil {
			ret.Friends = append(ret.Friends, m)
		}
	}

	for plrid := range self.ApplyList {
		plr := find_player(plrid)
		if plr != nil {
			ret.ApplyList = append(ret.ApplyList, plr.ToMsg_SimpleInfo())
		}
	}

	for plrid := range self.BlackList {
		plr := find_player(plrid)
		if plr != nil {
			ret.BlackList = append(ret.BlackList, plr.ToMsg_SimpleInfo())
		}
	}

	for v := range self.GiveList {
		ret.GiveList = append(ret.GiveList, v)
	}

	for v := range self.RecvList {
		ret.RecvList = append(ret.RecvList, v)
	}

	return ret
}

func (self *Friend) FriendInfo_ToMsg(plrid string) *msg.FriendInfo {
	frd := self.Friends[plrid]
	if frd == nil {
		return nil
	}

	plr := find_player(plrid)
	if plr == nil {
		self.Remove(plrid)
		return nil
	}

	ret := &msg.FriendInfo{
		Info:     plr.ToMsg_SimpleInfo(),
		IsOnline: plr.IsOnline(),
		OffTs:    plr.GetOfflineTs().Unix(),
		Ts:       frd.Ts.Unix(),
	}

	ret.ChatData = self.Chat_ToMsg(frd)

	return ret
}

func (self *Friend) Chat_ToMsg(frd *friend_t) []*msg.ChatOne {
	ret := []*msg.ChatOne{}
	for i := len(frd.Chat) - 1; i >= 0; i-- {
		c := frd.Chat[i]
		if c.Ts.Add(time.Duration(gconst.C_ChatExpireTs) * time.Hour).Before(time.Now()) {
			frd.Chat = frd.Chat[i+1:]
			break
		}

		fplr := load_player(c.FromId)
		if fplr != nil {
			if fplr.IsBan() {
				continue
			}

			ret = append(ret, &msg.ChatOne{
				Tp:      gconst.C_ChatType_Friend,
				From:    fplr.ToMsg_SimpleInfo(),
				ToId:    c.ToId,
				Content: c.Content,
				Ts:      c.Ts.Unix(),
			})
		}
	}

	return ret
}
