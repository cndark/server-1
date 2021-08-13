package mail

import (
	"fw/src/core/sched/async"
	"fw/src/game/app/gconst"
)

// ============================================================================

type sendmail_t struct {
	m   *Mail
	plr IPlayer
}

// ============================================================================

func flush(f func()) {
	async.PushQ(gconst.AQ_Mail, f)
}

// ============================================================================

func New(toplr interface{}) *sendmail_t {
	if toplr == nil {
		// gmail
		return &sendmail_t{
			m: GMailBox.create_mail(),
		}
	} else {
		// pmail
		plr := toplr.(IPlayer)

		return &sendmail_t{
			m:   plr.GetMailBox().create_mail(),
			plr: plr,
		}
	}
}

// ============================================================================

func (self *sendmail_t) SetKey(k int32) *sendmail_t {
	self.m.set_key(k)
	return self
}

func (self *sendmail_t) SetSender(v string) *sendmail_t {
	self.m.set_sender(v)
	return self
}

func (self *sendmail_t) SetTitle(v string) *sendmail_t {
	self.m.set_title(v)
	return self
}

func (self *sendmail_t) SetText(v string) *sendmail_t {
	self.m.set_text(v)
	return self
}

func (self *sendmail_t) AddDict(k, v string) *sendmail_t {
	self.m.add_dict(k, v)
	return self
}

func (self *sendmail_t) AddDictInt32(k string, v int32) *sendmail_t {
	self.m.add_dict_int32(k, v)
	return self
}

func (self *sendmail_t) AddAttachment(id int32, n float64) *sendmail_t {
	self.m.add_attachment(id, n)
	return self
}

func (self *sendmail_t) SetCond(v string) *sendmail_t {
	// only for gmail
	if self.plr == nil {
		self.m.set_cond(v)
	}
	return self
}

func (self *sendmail_t) Send() {
	if self.plr == nil {
		GMailBox.add(self.m)
	} else {
		self.plr.GetMailBox().add(self.m)
	}
}
