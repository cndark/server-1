package chat

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/guild"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"fw/src/shared/config"
	"time"

	Err "fw/src/proto/errorcode"
)

// ============================================================================
// 聊天管理
var ChatMgr = &Chat{
	Guild: make(map[string][]*ChatOne),
}

// ============================================================================

const (
	C_Content_Max_Cnt = 30
	C_Content_Max_Len = 140
)

// ============================================================================
// 聊天
type Chat struct {
	World    []*ChatOne            // 世界
	Guild    map[string][]*ChatOne // 家族
	CrossAll []*msg.ChatOne        // 跨服通服
}

// 本服聊天结构
type ChatOne struct {
	Tp      int32
	FromId  string
	Content string
	Ts      time.Time
	ToId    string
	GId     string
	GLv     int32
}

// ============================================================================
func init() {
	//remote
	evtmgr.On(gconst.Evt_GsPush_ChatCross, func(args ...interface{}) {
		oarg := args[1].([]byte)

		// unmarshal object arg
		var res *msg.GS_ChatMsg
		err := utils.UnmarshalArg(oarg, &res)
		if err != nil {
			log.Error("unmarshal GS_ChatMsg failed:", err)
			return
		}

		utils.BroadcastPlayers(res)

		// add
		ChatMgr.CrossAll = append(ChatMgr.CrossAll, res.One)
		if l := len(ChatMgr.CrossAll); l > C_Content_Max_Cnt {
			ChatMgr.CrossAll = ChatMgr.CrossAll[l-C_Content_Max_Cnt:]
		}
	})

	//guild destroy
	evtmgr.On(gconst.Evt_GuildDestroy, func(args ...interface{}) {
		gld := args[0].(*guild.Guild)

		delete(ChatMgr.Guild, gld.GetId())
	})

}

// ============================================================================
func Open() {
	load_data()
}

func Close() {
	save()
}

func load_data() {
	err := dbmgr.DBGame.GetObject(
		dbmgr.C_tabname_chat,
		1,
		&ChatMgr,
	)

	if ChatMgr.Guild == nil {
		ChatMgr.Guild = make(map[string][]*ChatOne)
	}

	if err != nil && !db.IsNotFound(err) {
		log.Warning("load chat data failed:", err)
	}
}

func save() {
	doc := save_gen_doc()

	err := dbmgr.DBGame.Upsert(
		dbmgr.C_tabname_chat,
		1,
		doc,
	)

	if err != nil {
		log.Warning("save chat data failed:", err)
	}

	log.Info("save chat")
}

func save_gen_doc() db.M {
	return db.M{
		"$set": db.M{
			"world":    core.CloneBsonArray(ChatMgr.World),
			"crossall": core.CloneBsonArray(ChatMgr.CrossAll),
			"guild":    core.CloneBsonArray(ChatMgr.Guild),
		},
	}
}

// ============================================================================

func Add(plr IPlayer, ci *ChatOne) (ec int32) {
	ec = Err.OK

	switch ci.Tp {
	case gconst.C_ChatType_World, gconst.C_ChatType_GldZm:
		ec = ChatMgr.to_world(plr, ci)
	case gconst.C_ChatType_Cross:
		ec = ChatMgr.to_cross(plr, ci)
	case gconst.C_ChatType_Guild:
		ec = ChatMgr.to_guild(plr, ci)
	default:
		ec = Err.Chat_TypeNotFound
	}

	return
}

func (self *Chat) to_world(plr IPlayer, ci *ChatOne) int32 {
	// broadcast
	utils.BroadcastPlayers(&msg.GS_ChatMsg{
		One: &msg.ChatOne{
			Tp:      ci.Tp,
			From:    plr.ToMsg_SimpleInfo(),
			Content: ci.Content,
			Ts:      ci.Ts.Unix(),
			GId:     ci.GId,
			GLv:     ci.GLv,
		},
	})

	// add
	self.World = append(self.World, ci)
	if l := len(self.World); l > C_Content_Max_Cnt {
		self.World = self.World[l-C_Content_Max_Cnt:]
	}

	return Err.OK
}

func (self *Chat) to_cross(plr IPlayer, ci *ChatOne) int32 {

	utils.GsPushAll(gconst.Evt_GsPush_ChatCross, nil, &msg.GS_ChatMsg{One: &msg.ChatOne{
		Tp:      ci.Tp,
		From:    plr.ToMsg_SimpleInfo(),
		Content: ci.Content,
		Ts:      time.Now().Unix(),
	}})

	return Err.OK
}

func (self *Chat) to_guild(plr IPlayer, ci *ChatOne) int32 {
	gldid := plr.GetGuildId()
	gld := guild.GuildMgr.FindGuild(gldid)
	if gld == nil {
		return Err.Guild_NotFound
	}

	// broadcast
	gld.Broadcast(&msg.GS_ChatMsg{
		One: &msg.ChatOne{
			Tp:      ci.Tp,
			From:    plr.ToMsg_SimpleInfo(),
			Content: ci.Content,
			Ts:      ci.Ts.Unix(),
		},
	})

	// add
	self.Guild[gldid] = append(self.Guild[gldid], ci)
	if l := len(self.Guild[gldid]); l > C_Content_Max_Cnt {
		self.Guild[gldid] = self.Guild[gldid][l-C_Content_Max_Cnt:]
	}

	return Err.OK
}

// ============================================================================

func ChatHistory_ToMsg(plr IPlayer) *msg.ChatData {
	ret := &msg.ChatData{}
	now := time.Now()

	for i := len(ChatMgr.World) - 1; i >= 0; i-- {
		c := ChatMgr.World[i]
		if c.Ts.Add(time.Duration(gconst.C_ChatExpireTs) * time.Hour).Before(now) {
			ChatMgr.World = ChatMgr.World[i+1:]
			break
		}

		cplr := find_player(c.FromId)
		if cplr != nil && !cplr.IsBan() {
			ret.Data = append(ret.Data, &msg.ChatOne{
				Tp:      c.Tp,
				From:    cplr.ToMsg_SimpleInfo(),
				Content: c.Content,
				Ts:      c.Ts.Unix(),
				GId:     c.GId,
				GLv:     c.GLv,
			})
		}
	}

	for i := len(ChatMgr.CrossAll) - 1; i >= 0; i-- {
		c := ChatMgr.CrossAll[i]
		if c.Ts+gconst.C_ChatExpireTs*3600 < now.Unix() {
			ChatMgr.CrossAll = ChatMgr.CrossAll[i+1:]
			break
		}

		if c.From.SvrId == config.CurGame.Id {
			cplr := find_player(c.From.Id)
			if cplr != nil && !cplr.IsBan() {
				c.From = cplr.ToMsg_SimpleInfo()
			}
		}

		ret.Data = append(ret.Data, c)
	}

	for gid, v := range ChatMgr.Guild {
		if gid == plr.GetGuildId() {
			for i := len(v) - 1; i >= 0; i-- {
				c := v[i]
				if c.Ts.Add(time.Duration(gconst.C_ChatExpireTs) * time.Hour).Before(now) {
					ChatMgr.Guild[gid] = ChatMgr.Guild[gid][i+1:]
					break
				}

				cplr := find_player(c.FromId)
				if cplr != nil && !cplr.IsBan() {
					ret.Data = append(ret.Data, &msg.ChatOne{
						Tp:      c.Tp,
						From:    cplr.ToMsg_SimpleInfo(),
						Content: c.Content,
						Ts:      c.Ts.Unix(),
					})
				}
			}
		}
	}

	return ret
}
