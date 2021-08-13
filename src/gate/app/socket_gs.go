package app

import (
	"fw/src/core/log"
	"fw/src/core/net/tcp"
	"fw/src/core/packet"
	"fw/src/gate/msg"
	"fw/src/shared/config"
)

// ============================================================================

type SocketGS struct {
	Id      int32
	sock    *tcp.Socket
	preader *packet.Reader
	pwriter *packet.Writer
}

// ============================================================================

func new_socket_gs(sock *tcp.Socket, id int32) *SocketGS {
	return &SocketGS{
		Id:      id,
		sock:    sock,
		preader: packet.NewReader(),
		pwriter: packet.NewWriter(),
	}
}

func (self *SocketGS) Open() {
	// register
	self.SendMsg(&msg.GW_RegisterGate{
		Id: config.CurGate.Id,
	})
}

func (self *SocketGS) Close() {
	self.sock.Close()
}

func (self *SocketGS) SendPacket(p packet.Packet) {
	buf := self.pwriter.Write(p)
	self.sock.Send(buf)
}

func (self *SocketGS) SendMsg(message msg.Message) {
	// marshal
	body, err := msg.Marshal(message)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return
	}

	// assemble
	p := packet.Assemble_B8(message.MsgId(), body, 0)

	// send
	self.SendPacket(p)
}

func (self *SocketGS) Dispatch(p packet.Packet) {
	// !Note: in net-thread

	sid := p.Remove_B8()
	if sid == 0 {
		// gate local msg
		op := p.Op()
		f := msg.MsgCreators[op]
		if f == nil {
			return
		}

		// unmarshal
		message := f()
		err := msg.Unmarshal(p.Body(), message)
		if err != nil {
			log.Error("unmarshal msg failed:", err)
			self.Close()
			return
		}

		h := msg.MsgHandlers[op]
		if h != nil {
			h(message, self)
		}
	} else {
		// forward to session
		NetMgr.Forward2Session(sid, p)
	}
}
