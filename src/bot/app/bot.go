package app

import (
	"fmt"
	"fw/src/bot/botconf"
	"fw/src/bot/msg"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/packet"
	Err "fw/src/proto/errorcode"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// ============================================================================

type Bot struct {
	Id          int32
	sock        socket_t
	preader     *packet.Reader
	pwriter     *packet.Writer
	locker_send sync.Mutex

	// msgq
	msgq chan func()
	quit chan bool

	// opt
	opt *bot_opt

	// data storage
	dstore map[string]interface{}

	// jobs
	wt_sum int32
	jobs   []*job_t
}

type bot_opt struct {
	svr        string
	sdk        string
	model      string
	job_prefix []string
	job_itv    []int32
	args       map[string]string
	grp        int32
}

type job_t struct {
	name string
	f    func(*Bot)
	wt   int32
}

// ============================================================================

func new_bot(sock socket_t, id int32, opt *bot_opt) *Bot {
	bot := &Bot{
		Id:      id,
		sock:    sock,
		preader: packet.NewReader(),
		pwriter: packet.NewWriter(),

		msgq: make(chan func(), 1000),
		quit: make(chan bool),

		opt: opt,

		dstore: make(map[string]interface{}),
	}

	bot.preader.SetDecryptor(NewRc4())
	bot.pwriter.SetEncryptor(NewRc4())

	return bot
}

func (self *Bot) SendPacket(p packet.Packet) {
	self.locker_send.Lock()
	defer self.locker_send.Unlock()

	buf := self.pwriter.Write(p)
	self.sock.Send(buf)
}

func (self *Bot) SendMsg(message msg.Message) {
	body, err := msg.Marshal(message)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return
	}

	p := packet.Assemble(message.MsgId(), body)

	self.SendPacket(p)
}

func (self *Bot) Close() {
	self.sock.Close()
}

func (self *Bot) Dispatch(p packet.Packet) {
	// !Note: in net-thread

	op := p.Op()
	f := msg.MsgCreators[op]
	if f != nil {
		message := f()
		err := msg.Unmarshal(p.Body(), message)
		if err != nil {
			log.Warning("unmarshal msg failed:", err)
			return
		}

		h := msg.MsgHandlers[op]
		if h != nil {
			self.Push(func() {
				h(message, self)
			})
		}

		// fire
		self.Push(func() {
			evtmgr.Fire(fmt.Sprintf("msg.%d", op), self, message)
		})
	}
}

// ============================================================================

func (self *Bot) Push(f func()) {
	defer func() { recover() }()

	select {
	case self.msgq <- f:
	default:
		log.Error("msgq FULL. push discarded")
	}
}

func (self *Bot) run() {
	go func() {
		defer close(self.msgq)

		for {
			select {
			case <-self.quit:
				return

			case f := <-self.msgq:
				core.PCall(f)
			}
		}
	}()
}

func (self *Bot) stop() {
	close(self.quit)
}

// ============================================================================

func (self *Bot) GetArg(key string) string {
	return self.opt.args[key]
}

func (self *Bot) GetData(key string) interface{} {
	return self.dstore[key]
}

func (self *Bot) SetData(key string, d interface{}) {
	self.dstore[key] = d
}

// ============================================================================

func (self *Bot) JobAdd(name string, f func(*Bot)) {
	if !self.job_check_prefix(name) {
		return
	}

	if !self.job_check_grp_binding(name) {
		return
	}

	self.jobs = append(self.jobs, &job_t{
		name: name,
		f:    f,
		wt:   100, // default weight
	})
}

func (self *Bot) job_check_prefix(name string) bool {
	for _, prefix := range self.opt.job_prefix {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

func (self *Bot) job_check_grp_binding(name string) bool {
	arr := botconf.JobCtl.GrpBinding[name]
	if arr == nil { // no binding
		return true
	}

	// has binding: run only if grp matched
	for _, grp := range arr {
		if self.opt.grp == grp {
			return true
		}
	}

	return false
}

func (self *Bot) job_start() {
	// do nothing if no job
	L := len(self.jobs)
	if L == 0 {
		return
	}

	// prepare job weights
	self.job_prepare_weight()

	// heartbeat thread
	core.Go(func() {
		for {
			select {
			case <-self.quit:
				return

			default:
				time.Sleep(time.Second * 15)
				self.SendMsg(&msg.C_TimeSync{})
			}
		}
	})

	// job thread
	core.Go(func() {
		a := self.opt.job_itv[0]
		b := self.opt.job_itv[1]

		for {
			select {
			case <-self.quit:
				return

			default:
				// wait
				time.Sleep(time.Second * time.Duration(core.RandInt32(a, b)))

				// randomly pick a job based on job weight
				job := self.job_select()

				// check trigger prob
				if !self.job_check_tri_prob(job.name) {
					continue
				}

				// run
				job.f(self)
			}
		}
	})
}

func (self *Bot) job_prepare_weight() {
	self.wt_sum = 0

	for _, v := range self.jobs {
		wt := botconf.JobCtl.SelectWeight[v.name]
		if wt > 0 {
			v.wt = wt
		} else {
			wt = v.wt
		}

		self.wt_sum += wt
	}
}

func (self *Bot) job_select() *job_t {
	p := rand.Int31n(self.wt_sum)
	for _, v := range self.jobs {
		p -= v.wt
		if p < 0 {
			return v
		}
	}

	return nil
}

func (self *Bot) job_check_tri_prob(name string) bool {
	p := botconf.JobCtl.TriggerProb[name]
	if p == 0 { // not set
		return true
	}

	return rand.Float32() < p
}

// ============================================================================

func (self *Bot) OnConnected() {
	log.Info("bot connected:", self.Id)

	self.SendMsg(&msg.C_Auth{
		AuthId:    core.I32toa(self.Id),
		AuthToken: botconf.Password,
		Sdk:       self.opt.sdk,
		Model:     self.opt.model,
		DevId:     fmt.Sprintf("aibot-%d", self.Id),
		Os:        "linux",
		OsVer:     "7",
		VerMajor:  botconf.Ver_Major,
		VerMinor:  botconf.Ver_Minor,
		VerBuild:  botconf.Ver_Build,
	})
}

func (self *Bot) OnDisconnected() {
	log.Info("bot disconnected:", self.Id)
}

func (self *Bot) OnAuth(req *msg.GW_Auth_R) {
	if req.ErrorCode != Err.OK {
		log.Error("bot auth failed:", self.Id, "ErrorCode:", req.ErrorCode)
		self.Close()
		return
	}

	self.SendMsg(&msg.C_Login{
		Svr0:     self.opt.svr,
		ChgSvr:   false,
		Language: "cn",
	})
}

func (self *Bot) OnLogin(req *msg.GW_Login_R) {
	if req.ErrorCode != Err.OK {
		log.Error("bot login failed:", self.Id, "ErrorCode:", req.ErrorCode)
		self.Close()
		return
	}
}

func (self *Bot) OnUserInfo(user *msg.GS_UserInfo) {
	log.Info("user info:", user.UserId, user.Name)

	// fire
	evtmgr.Fire("userinfo", self, user)

	// start job chain
	self.job_start()
}
