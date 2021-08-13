package mail

import (
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/modules/utils"
	"sync/atomic"
	"time"
)

// ============================================================================

const (
	C_gmail_max = 300
)

// ============================================================================

var (
	GMailBox = &gmail_box_t{}
)

// ============================================================================

type gmail_box_t struct {
	seqid int32
	mails []*Mail
}

// ============================================================================

func (self *gmail_box_t) Init() {
	// load gmails
	var obj struct {
		SeqId int32
		Mails []*Mail
	}

	err := dbmgr.DBGame.GetObject(
		dbmgr.C_tabname_gmail,
		1,
		&obj,
	)
	if err == nil {

		self.seqid = obj.SeqId
		self.mails = obj.Mails

	} else if db.IsNotFound(err) {

		self.seqid = 0
		self.mails = nil

	} else {
		log.Warning("loading gmail-data failed:", err)
	}
}

func (self *gmail_box_t) gen_mailid() int32 {
	id := atomic.AddInt32(&self.seqid, 1)

	flush(func() {
		err := dbmgr.DBGame.Upsert(
			dbmgr.C_tabname_gmail,
			1,
			db.M{
				"$set": db.M{"seqid": id},
			},
		)
		if err != nil {
			log.Error("gmail_box_t.gen_mailid() failed:", err)
		}
	})

	return id
}

func (self *gmail_box_t) create_mail() *Mail {

	return &Mail{
		Id:    self.gen_mailid(),
		Read:  false,
		Taken: false,
		Ts:    time.Now(),
	}
}

func (self *gmail_box_t) add(m *Mail) {
	self.mails = append(self.mails, m)

	// limit
	L := len(self.mails)
	if L > C_gmail_max {
		self.mails = self.mails[L-C_gmail_max:]
	}

	// flush
	err := dbmgr.DBGame.Upsert(
		dbmgr.C_tabname_gmail,
		1,
		db.M{
			"$push": db.M{
				"mails": db.M{
					"$each":  []*Mail{m},
					"$slice": -C_gmail_max,
				},
			},
		},
	)
	if err != nil {
		log.Error("gmail_box_t.Add() failed:", err)
	}

	// sync gmail for all online players
	utils.ForEachOnlinePlayer(func(plr interface{}) {
		plr.(IPlayer).GetMailBox().SyncGMails()
	})
}

func (self *gmail_box_t) latest_id() int32 {
	return self.seqid
}

func (self *gmail_box_t) query_latest_mails(fromid int32) (ret []*Mail) {
	L := len(self.mails)

	for i := L - 1; i >= 0; i-- {
		m := self.mails[i]

		if m.Id <= fromid {
			break
		}

		ret = append(ret, m)
	}

	return
}
