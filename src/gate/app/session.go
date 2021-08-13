package app

import (
	"fw/src/core/log"
	"fw/src/core/packet"
	"fw/src/gate/msg"
	"fw/src/shared/config"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================

var seq_sid uint64

func InitSession() {
	seq_sid = uint64(config.CurGate.Id) << 41
}

// ============================================================================

const (
	AuthState_None = iota
	AuthState_InProgress
	AuthState_OK
)

const (
	LoginState_None = iota
	LoginState_InProgress
)

type Session struct {
	id      uint64 // session id
	gsid    int32  // game id
	sock    socket_t
	preader *packet.Reader
	pwriter *packet.Writer

	// auth
	locker_auth sync.Mutex
	state_auth  int
	AuthReq     *msg.C_Auth
	AuthRet     map[string]string

	// login
	locker_login sync.Mutex
	state_login  int

	// token
	locker_tk sync.Mutex
	tk_get_ts time.Time
}

// ============================================================================

func new_session(sock socket_t) *Session {
	s := &Session{
		id:          atomic.AddUint64(&seq_sid, 1),
		gsid:        0,
		sock:        sock,
		preader:     packet.NewReader(),
		pwriter:     packet.NewWriter(),
		state_auth:  AuthState_None,
		state_login: LoginState_None,
		AuthRet:     make(map[string]string),
	}

	s.preader.SetMaxPacketLen(1024 * 200) // 200 KB

	s.preader.SetDecryptor(NewRc4())
	s.pwriter.SetEncryptor(NewRc4())

	return s
}

func (self *Session) Close() {
	self.sock.Close()
}

func (self *Session) SendPacket(p packet.Packet) {
	buf := self.pwriter.Write(p)
	self.sock.Send(buf)
}

func (self *Session) SendMsg(message msg.Message) {
	body, err := msg.Marshal(message)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return
	}

	p := packet.Assemble(message.MsgId(), body)

	self.SendPacket(p)
}

func (self *Session) Dispatch(p packet.Packet) {
	// !Note: in net-thread

	op := p.Op()
	f := msg.MsgCreators[op]
	if f == nil {
		// op NOT found. forward to gs. MUST be authenticated
		if !self.IsAuthenticated() {
			return
		}

		NetMgr.Forward2GS(self.gsid, self.id, p)

	} else {
		// gate local handler
		message := f()
		err := msg.Unmarshal(p.Body(), message)
		if err != nil {
			self.Close()
			return
		}

		h := msg.MsgHandlers[op]
		if h != nil {
			h(message, self)
		}
	}
}

// ============================================================================

func (self *Session) GetId() uint64 {
	return self.id
}

func (self *Session) GetIP() string {
	return self.sock.RemoteIP()
}

func (self *Session) SetGsId(id int32) {
	if self.gsid != 0 {
		self.LogoutPlayer()
	}
	self.gsid = id
}

func (self *Session) AuthBegin() bool {
	self.locker_auth.Lock()
	defer self.locker_auth.Unlock()

	if self.state_auth != AuthState_None {
		return false
	}

	self.state_auth = AuthState_InProgress
	return true
}

func (self *Session) AuthEnd() {
	self.locker_auth.Lock()
	defer self.locker_auth.Unlock()

	if self.AuthReq == nil {
		self.state_auth = AuthState_None
	} else {
		self.state_auth = AuthState_OK
	}
}

func (self *Session) IsAuthenticated() bool {
	return self.state_auth == AuthState_OK // no need to lock
}

func (self *Session) LoginBegin() bool {
	self.locker_login.Lock()
	defer self.locker_login.Unlock()

	if self.state_login != LoginState_None {
		return false
	}

	self.state_login = LoginState_InProgress
	return true
}

func (self *Session) LoginEnd() {
	self.locker_login.Lock()
	defer self.locker_login.Unlock()

	self.state_login = LoginState_None
}

func (self *Session) LoginPlayer(gsid int32, m *msg.GW_UserOnline) bool {
	self.SetGsId(gsid)
	return NetMgr.Send2GS(self.gsid, m)
}

func (self *Session) LogoutPlayer() {
	NetMgr.Send2GS(self.gsid, &msg.GW_LogoutPlayer{
		Sid: self.id,
	})
}

func (self *Session) CheckTokenGet() bool {
	self.locker_tk.Lock()
	defer self.locker_tk.Unlock()

	now := time.Now()

	if now.Sub(self.tk_get_ts).Minutes() < 30 {
		return false
	}

	self.tk_get_ts = now
	return true
}
