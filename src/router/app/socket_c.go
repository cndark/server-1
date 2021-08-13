package app

import (
	"fw/src/core/log"
	"fw/src/core/net/tcp"
	"fw/src/core/packet"
	"fw/src/router/msg"
	"sync/atomic"
)

// ============================================================================

type SocketC struct {
	id      int32
	sock    *tcp.Socket
	preader *packet.Reader
	pwriter *packet.Writer
}

// ============================================================================

func new_socket_c(sock *tcp.Socket) *SocketC {
	return &SocketC{
		id:      atomic.AddInt32(&seq_cid, 1),
		sock:    sock,
		preader: packet.NewReader(),
		pwriter: packet.NewWriter(),
	}
}

func (self *SocketC) Close() {
	self.sock.Close()
}

func (self *SocketC) SendPacket(p packet.Packet) {
	buf := self.pwriter.Write(p)
	self.sock.Send(buf)
}

func (self *SocketC) SendMsg(message msg.Message) {
	// marshal
	body, err := msg.Marshal(message)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return
	}

	// assemble
	p := packet.Assemble_B8_Str(message.MsgId(), body, 0, "")

	// send
	self.SendPacket(p)
}

func (self *SocketC) Dispatch(p packet.Packet) {
	// !Note: in net-thread

	// check route
	if p = NetMgr.check_route(p); p == nil {
		return
	}

	// local msg. get msg creator
	op := p.Op()
	f := msg.MsgCreators[op]
	if f == nil {
		return
	}

	// unmarshal
	message := f()
	err := msg.Unmarshal(p.Body(), message)
	if err != nil {
		return
	}

	// find handler
	h := msg.MsgHandlers[op]
	if h != nil {
		h(message, self)
	}
}

// ============================================================================

func (self *SocketC) is_registered() bool {
	return self.id < C_max_cid
}
