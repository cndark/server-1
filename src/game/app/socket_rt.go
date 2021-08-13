package app

import (
	"fw/src/core/log"
	"fw/src/core/net/tcp"
	"fw/src/core/packet"
	"fw/src/core/sched/loop"
	"fw/src/game/msg"
	"fw/src/shared/config"
)

// ============================================================================

type SocketRt struct {
	Id      int32
	sock    *tcp.Socket
	preader *packet.Reader
	pwriter *packet.Writer
}

// ============================================================================

func new_socket_rt(sock *tcp.Socket, id int32) *SocketRt {
	return &SocketRt{
		Id:      id,
		sock:    sock,
		preader: packet.NewReader(),
		pwriter: packet.NewWriter(),
	}
}

func (self *SocketRt) Open() {
	// register
	self.SendMsg(&msg.GS_RegisterGame{
		Id: config.CurGame.Id,
	})
}

func (self *SocketRt) Close() {
	self.sock.Close()
}

func (self *SocketRt) SendPacket(p packet.Packet) {
	buf := self.pwriter.Write(p)
	self.sock.Send(buf)
}

func (self *SocketRt) SendMsg(message msg.Message) {
	p := self.to_packet(message, 0, "")
	if p == nil {
		return
	}

	// send
	self.SendPacket(p)
}

func (self *SocketRt) to_packet(message msg.Message, b8 uint64, str string) packet.Packet {
	// marshal
	body, err := msg.Marshal(message)
	if err != nil {
		log.Error("marshal msg failed:", message.MsgId(), err)
		return nil
	}

	// assemble
	p := packet.Assemble_B8_Str(message.MsgId(), body, b8, str)

	return p
}

func (self *SocketRt) Dispatch(p packet.Packet) {
	// !Note: in net-thread

	// check route
	if p = NetMgr.check_route(p); p == nil {
		return
	}

	// get msg creator
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
		loop.Push(func() {
			h(message, self)
		})
	}
}
