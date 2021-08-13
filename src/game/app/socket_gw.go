package app

import (
	"fw/src/core/log"
	"fw/src/core/net/tcp"
	"fw/src/core/packet"
	"fw/src/core/sched/loop"
	"fw/src/game/msg"
	"sync/atomic"
)

// ============================================================================

var (
	C_op_register_gate = (&msg.GW_RegisterGate{}).MsgId()
)

// ============================================================================

type SocketGW struct {
	id      int32
	sock    *tcp.Socket
	preader *packet.Reader
	pwriter *packet.Writer
}

// ============================================================================

func new_socket_gw(sock *tcp.Socket) *SocketGW {
	return &SocketGW{
		id:      atomic.AddInt32(&seq_gateid, 1),
		sock:    sock,
		preader: packet.NewReader(),
		pwriter: packet.NewWriter(),
	}
}

func (self *SocketGW) Close() {
	self.sock.Close()
}

func (self *SocketGW) SendPacket(p packet.Packet) {
	buf := self.pwriter.Write(p)
	self.sock.Send(buf)
}

func (self *SocketGW) SendMsg(message msg.Message) {
	body, err := msg.Marshal(message)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return
	}

	p := packet.Assemble_B8(message.MsgId(), body, 0)

	self.SendPacket(p)
}

func (self *SocketGW) Dispatch(p packet.Packet) {
	// !Note: in net-thread

	op := p.Op()

	// check register
	if !self.is_registered() && op != C_op_register_gate {
		return
	}

	// get msg creator
	f := msg.MsgCreators[op]
	if f == nil {
		return
	}

	// remove sid
	sid := p.Remove_B8()

	// unmarshal
	message := f()
	err := msg.Unmarshal(p.Body(), message)
	if err != nil {
		return
	}

	// find handler
	h := msg.MsgHandlers[op]
	if h != nil {
		loop.Push(func() {
			// set ctx
			var ctx interface{}
			if sid == 0 {
				// directly from gate
				ctx = self
			} else {
				// from player
				plr := PlayerMgr.FindPlayerBySid(sid) // this function MUST be run in loop thread
				if plr == nil {
					return
				} else {
					ctx = plr
				}
			}

			h(message, ctx)
		})
	}
}

// ============================================================================

func (self *SocketGW) is_registered() bool {
	return self.id < C_max_gateid
}
