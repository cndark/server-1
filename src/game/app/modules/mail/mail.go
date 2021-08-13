package mail

import (
	"fmt"
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/core/sched/loop"
	"fw/src/game/app/dbmgr"
	"fw/src/game/msg"
	Err "fw/src/proto/errorcode"
	"math/rand"
	"sort"
	"sync/atomic"
	"time"
)

// ============================================================================

const (
	C_mail_max             = 100
	C_mail_days_attachment = 15
	C_mail_days_unread     = 7
	C_mail_days_read       = 2
)

// ============================================================================

/*
MailBox struct
	seqid:   number
	gmailid: number
	mails:   {
		id: {
			contents
		}
	}
*/
type MailBox struct {
	seqid   int32           // mail id sequence
	gmailid int32           // last global mail id. !!! 如果转服、合服，gmailid 必须重新和新服同步
	mails   map[int32]*Mail // mails

	tid *core.Timer // timer for expire check

	plr IPlayer
}

/*
Mail struct
    Title, Text support both id and string formats
    id:     "n:123456"
    string: "s:helloworld"

    if no 'prefix' is found, it'll be treated as a plane-string

Deliver condition format:
	1. empty: 				no condition
	2. cond|cond|...		satisfy all 'cond'
	3. any unrecognized 'cond' will block the delivery

	cond: var op val

	var:
		lv:		player level
		cdate:  player creation date
		ldate:  player last login date

	op:   <, <=, >, >=, ==, !=

	val:
		number
		string
		date: 	format -> 2017-1-1 3:4:5
*/
type Mail struct {
	Id         int32             `bson:"id"`              // mail id
	Key        int32             `bson:"key"`             // mail key
	Sender     string            `bson:"sender"`          // sender
	Title      string            `bson:"title,omitempty"` // mail title (id | string)
	Text       string            `bson:"text,omitempty"`  // mail text  (id | string)
	Dict       map[string]string `bson:"dict,omitempty"`  // dict
	Attachment []*res_t          `bson:"a,omitempty"`     // attachment
	Read       bool              `bson:"read"`            // read flag
	Taken      bool              `bson:"taken"`           // taken flag
	Ts         time.Time         `bson:"ts"`              // create timestamp
	ExpireTs   time.Time         `bson:"ets"`             // expire timestamp
	Cond       string            `bson:",omitempty"`      // deliver condition for gmail
}

type res_t struct {
	Id int32
	N  float64
}

// ============================================================================

func NewMailBox() *MailBox {
	return &MailBox{}
}

func (self *MailBox) Init(plr IPlayer) {
	self.plr = plr

	// load mails
	self.load_mails()

	// start expire-check timer
	self.start_expire_check_timer()
}

func (self *MailBox) load_mails() {
	self.mails = make(map[int32]*Mail)

	var obj struct {
		MailBox *struct {
			SeqId   int32
			GmailId int32
			Mails   map[string]*Mail
		}
	}

	err := self.plr.DB().GetProjection(
		dbmgr.C_tabname_user,
		self.plr.GetId(),
		db.M{"mailbox": 1},
		&obj,
	)
	if err != nil {
		core.Panic("loading player mail-data failed:", self.plr.GetId())
	}

	// set
	if obj.MailBox == nil {
		self.seqid = 0
		self.gmailid = GMailBox.latest_id()
	} else {
		self.seqid = obj.MailBox.SeqId
		self.gmailid = obj.MailBox.GmailId

		for _, m := range obj.MailBox.Mails {
			self.mails[m.Id] = m
		}
	}
}

func (self *MailBox) start_expire_check_timer() {
	self.check_expire()
	self.tid = loop.SetTimeout(time.Now().Add(time.Duration(30+rand.Intn(30))*time.Minute), func() {
		self.start_expire_check_timer()
	})
}

func (self *MailBox) gen_mailid() int32 {
	id := atomic.AddInt32(&self.seqid, 1)

	flush(func() {
		err := self.plr.DB().Update(
			dbmgr.C_tabname_user,
			self.plr.GetId(),
			db.M{
				"$set": db.M{"mailbox.seqid": id},
			},
		)
		if err != nil {
			log.Error("MailBox.gen_mailid() failed:", err)
		}
	})

	return id
}

func (self *MailBox) check_expire() {
	var ids []int32

	// check expire
	now := time.Now()
	for id, m := range self.mails {
		if m.ExpireTs.Before(now) {
			ids = append(ids, id)
		}
	}

	for _, id := range ids {
		delete(self.mails, id)
	}

	// check max mails
	L := len(self.mails)
	if L > C_mail_max {
		arr := make([]*Mail, 0, L)
		for _, m := range self.mails {
			arr = append(arr, m)
		}

		sort.Slice(arr, func(i, j int) bool {
			return arr[i].ExpireTs.Before(arr[j].ExpireTs)
		})

		for i := 0; i < L-C_mail_max; i++ {
			id := arr[i].Id
			delete(self.mails, id)
			ids = append(ids, id)
		}
	}

	// check empty
	if ids == nil {
		return
	}

	// flush
	flush(func() {
		doc := db.M{}
		for _, id := range ids {
			doc[fmt.Sprintf("mailbox.mails.%d", id)] = 0
		}

		err := self.plr.DB().Update(
			dbmgr.C_tabname_user,
			self.plr.GetId(),
			db.M{"$unset": doc},
		)
		if err != nil {
			log.Error("MailBox.check_expire() failed:", err)
		}
	})

	// push: mail del
	self.plr.SendMsg(&msg.GS_MailDel{
		Ids: ids,
	})
}

func (self *MailBox) create_mail() *Mail {

	return &Mail{
		Id:    self.gen_mailid(),
		Read:  false,
		Taken: false,
		Ts:    time.Now(),
	}
}

func (self *MailBox) add(m *Mail) {
	m.update_expire_ts()
	self.mails[m.Id] = m

	// flush
	flush(func() {
		err := self.plr.DB().Update(
			dbmgr.C_tabname_user,
			self.plr.GetId(),
			db.M{
				"$set": db.M{fmt.Sprintf("mailbox.mails.%d", m.Id): m},
			},
		)
		if err != nil {
			log.Error("MailBox.Add() failed:", err)
		}
	})

	// push: mail new
	self.plr.SendMsg(&msg.GS_MailNew{
		M: m.ToMsg(),
	})
}

func (self *MailBox) Remove(id int32) int32 {
	if id == 0 {
		if len(self.mails) == 0 {
			return Err.Mail_NoMail
		}

		self.mails = make(map[int32]*Mail)

		// flush
		flush(func() {
			err := self.plr.DB().Update(
				dbmgr.C_tabname_user,
				self.plr.GetId(),
				db.M{"$set": db.M{"mailbox.mails": db.M{}}},
			)
			if err != nil {
				log.Error("MailBox.Remove() failed:", err)
			}
		})

	} else {
		m := self.mails[id]
		if m == nil {
			return Err.Mail_NotFound
		}

		delete(self.mails, id)

		// flush
		flush(func() {
			err := self.plr.DB().Update(
				dbmgr.C_tabname_user,
				self.plr.GetId(),
				db.M{"$unset": db.M{fmt.Sprintf("mailbox.mails.%d", id): 0}},
			)
			if err != nil {
				log.Error("MailBox.Remove() failed:", err)
			}
		})
	}

	return Err.OK
}

func (self *MailBox) Read(id int32) (ec int32, affected *Mail) {
	m := self.mails[id]
	if m == nil {
		ec = Err.Mail_NotFound
		return
	}

	if m.Read {
		ec = Err.Mail_AlreadyRead
		return
	}

	m.Read = true
	m.update_expire_ts()

	// flush
	doc := db.M{
		fmt.Sprintf("mailbox.mails.%d.read", id): true,
		fmt.Sprintf("mailbox.mails.%d.ets", id):  m.ExpireTs,
	}

	flush(func() {
		err := self.plr.DB().Update(
			dbmgr.C_tabname_user,
			self.plr.GetId(),
			db.M{"$set": doc},
		)
		if err != nil {
			log.Error("MailBox.Read() failed:", err)
		}
	})

	// ok
	affected = m
	return
}

func (self *MailBox) TakeAttachment(id int32) (ec int32, a []*res_t, affected *Mail) {
	m := self.mails[id]
	if m == nil {
		ec = Err.Mail_NotFound
		return
	}

	if m.Taken {
		ec = Err.Mail_AlreadyTaken
		return
	}

	if m.Attachment == nil {
		ec = Err.Mail_NoAttachment
		return
	}

	a = m.Attachment
	m.Read = true
	m.Taken = true
	m.update_expire_ts()

	// flush
	doc := db.M{
		fmt.Sprintf("mailbox.mails.%d.read", id):  true,
		fmt.Sprintf("mailbox.mails.%d.taken", id): true,
		fmt.Sprintf("mailbox.mails.%d.ets", id):   m.ExpireTs,
	}

	flush(func() {
		err := self.plr.DB().Update(
			dbmgr.C_tabname_user,
			self.plr.GetId(),
			db.M{
				"$set": doc,
			},
		)
		if err != nil {
			log.Error("MailBox.TakeAttachment() failed:", err)
		}
	})

	// ok
	affected = m
	return
}

func (self *MailBox) TakeAttachmentAll() (ec int32, a []*res_t, affected []*Mail) {
	if len(self.mails) == 0 {
		ec = Err.Mail_NoMail
		return
	}

	doc := db.M{}

	for _, m := range self.mails {
		if m.Taken || m.Attachment == nil {
			continue
		}

		a = append(a, m.Attachment...)
		m.Read = true
		m.Taken = true
		m.update_expire_ts()

		affected = append(affected, m)

		doc[fmt.Sprintf("mailbox.mails.%d.read", m.Id)] = true
		doc[fmt.Sprintf("mailbox.mails.%d.taken", m.Id)] = true
		doc[fmt.Sprintf("mailbox.mails.%d.ets", m.Id)] = m.ExpireTs
	}

	if a == nil {
		ec = Err.Mail_NoAttachment
		return
	}

	// flush
	flush(func() {
		err := self.plr.DB().Update(
			dbmgr.C_tabname_user,
			self.plr.GetId(),
			db.M{
				"$set": doc,
			},
		)
		if err != nil {
			log.Error("MailBox.TakeAttachmentAll() failed:", err)
		}
	})

	//ok
	return
}

func (self *MailBox) RemoveOnekey() (ids []int32) {
	// remove those read and no attachment left
	for _, m := range self.mails {
		if m.Read && (m.Attachment == nil || m.Taken) {
			ids = append(ids, m.Id)
		}
	}

	if len(ids) == 0 {
		return
	}

	// delete from mails
	for _, id := range ids {
		delete(self.mails, id)
	}

	// unset doc
	doc := db.M{}

	for _, id := range ids {
		doc[fmt.Sprintf("mailbox.mails.%d", id)] = 0
	}

	// flush
	flush(func() {
		err := self.plr.DB().Update(
			dbmgr.C_tabname_user,
			self.plr.GetId(),
			db.M{"$unset": doc},
		)
		if err != nil {
			log.Error("MailBox.RemoveOnekey() failed:", err)
		}
	})

	return
}

func (self *MailBox) SyncGMails() {
	mails := GMailBox.query_latest_mails(self.gmailid)
	if mails == nil {
		return
	}

	// add new gmails
	for _, m := range mails {
		// check cond
		if !self.plr.CheckGmailDeliverCond(m.Cond) {
			continue
		}

		// convert
		m := m.gmail_to_mail(self.gen_mailid())

		// add to mailbox
		self.add(m)
	}

	// update gmailid
	new_gmailid := mails[0].Id
	self.gmailid = new_gmailid

	// flush
	flush(func() {
		err := self.plr.DB().Update(
			dbmgr.C_tabname_user,
			self.plr.GetId(),
			db.M{
				"$set": db.M{"mailbox.gmailid": new_gmailid},
			},
		)
		if err != nil {
			log.Error("update gmailid failed:", err)
		}
	})
}

func (self *MailBox) ToMsg() (ret []*msg.Mail) {
	for _, m := range self.mails {
		ret = append(ret, m.ToMsg())
	}

	return
}

// ============================================================================

func (self *Mail) set_key(k int32) *Mail {
	self.Key = k
	return self
}

func (self *Mail) set_sender(v string) *Mail {
	self.Sender = v
	return self
}

func (self *Mail) set_title(v string) *Mail {
	self.Title = v
	return self
}

func (self *Mail) set_text(v string) *Mail {
	self.Text = v
	return self
}

func (self *Mail) add_dict(k, v string) *Mail {
	if self.Dict == nil {
		self.Dict = make(map[string]string)
	}

	self.Dict[k] = v
	return self
}

func (self *Mail) add_dict_int32(k string, v int32) *Mail {
	return self.add_dict(k, core.I32toa(v))
}

func (self *Mail) add_attachment(id int32, n float64) *Mail {
	self.Attachment = append(self.Attachment, &res_t{id, n})
	return self
}

func (self *Mail) set_cond(v string) *Mail {
	self.Cond = v
	return self
}

func (self *Mail) update_expire_ts() {
	if !self.Taken && len(self.Attachment) > 0 {
		self.ExpireTs = self.Ts.Add(C_mail_days_attachment * 24 * time.Hour)
	} else if self.Read {
		self.ExpireTs = time.Now().Add(C_mail_days_read * 24 * time.Hour)
	} else {
		self.ExpireTs = self.Ts.Add(C_mail_days_unread * 24 * time.Hour)
	}
}

func (self *Mail) gmail_to_mail(newid int32) *Mail {

	m := &Mail{
		Id:         newid,
		Key:        self.Key,
		Sender:     self.Sender,
		Title:      self.Title,
		Text:       self.Text,
		Dict:       nil,
		Attachment: nil,
		Read:       self.Read,
		Taken:      self.Taken,
		Ts:         self.Ts,
		ExpireTs:   self.ExpireTs,
	}

	// dict
	for k, v := range self.Dict {
		m.Dict[k] = v
	}

	// attachment
	for _, res := range self.Attachment {
		m.Attachment = append(m.Attachment, &res_t{res.Id, res.N})
	}

	return m
}

func (self *Mail) ToMsg() *msg.Mail {
	ret := &msg.Mail{
		Id:     self.Id,
		Key:    self.Key,
		Sender: self.Sender,
		Title:  self.Title,
		Text:   self.Text,
		Dict:   self.Dict,
		Read:   self.Read,
		Taken:  self.Taken,
		Ts:     self.Ts.Unix(),
		ETs:    self.ExpireTs.Unix(),
	}

	for _, res := range self.Attachment {
		ret.A = append(ret.A, &msg.MailRes{res.Id, res.N})
	}

	return ret
}
