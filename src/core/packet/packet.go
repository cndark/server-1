package packet

import (
	"encoding/binary"
	"errors"
)

// ============================================================================

/*
	packet format:
		length 4 bytes	(packet total length)
		op     4 bytes
		body   []byte
		b8     8 bytes  (optional)
		str    []byte   (optional; if exists, b8 MUST exist)
		lstr   1 byte   (optional; 'str' length; exists if 'str' exists)
*/

// ============================================================================

var (
	Err_PacketLength = errors.New("invalid packet length")
)

// ============================================================================

type IEncryptor interface {
	Encrypt(dst, src []byte)
}

type IDecryptor interface {
	Decrypt(dst, src []byte)
}

// ============================================================================

type Reader struct {
	data []byte
	ptr  []byte
	l    uint32

	max_pkt_len uint32 // 0: unlimited

	decryptor IDecryptor
}

func NewReader() *Reader {
	r := &Reader{
		max_pkt_len: 10 * 1024 * 1024, // 10 MB by default
	}

	r.reset()

	return r
}

func (self *Reader) reset() {
	self.data = make([]byte, 4)
	self.ptr = self.data
	self.l = 0
}

func (self *Reader) SetMaxPacketLen(n uint32) {
	self.max_pkt_len = n
}

func (self *Reader) SetDecryptor(v IDecryptor) {
	self.decryptor = v
}

func (self *Reader) Read(buf []byte) (p Packet, buf_ptr []byte, err error) {

	buf_ptr = buf

	// read length
	if self.l == 0 {
		n := copy(self.ptr, buf_ptr)
		self.ptr, buf_ptr = self.ptr[n:], buf_ptr[n:]

		// partial
		if len(self.ptr) > 0 {
			return
		}

		// decrypt length
		if self.decryptor != nil {
			self.decryptor.Decrypt(self.data, self.data)
		}

		// get length
		self.l = binary.BigEndian.Uint32(self.data)

		// check length
		if self.l < 8 || (self.max_pkt_len > 0 && self.l > self.max_pkt_len) {
			self.l = 0 // for safety
			return nil, nil, Err_PacketLength
		}

		// alloc packet buffer
		//	+8: provide more room for appending 'b8' without buffer re-allocation
		new_data := make([]byte, self.l, self.l+8)
		copy(new_data, self.data)
		self.data = new_data
		self.ptr = self.data[4:]
	}

	// read rest
	n := copy(self.ptr, buf_ptr)
	self.ptr, buf_ptr = self.ptr[n:], buf_ptr[n:]

	// partial
	if len(self.ptr) > 0 {
		return
	}

	// full packet
	p = self.data

	// decrypt rest
	if self.decryptor != nil {
		self.decryptor.Decrypt(p[4:], p[4:])
	}

	// reset
	self.reset()

	return
}

// ============================================================================

type Writer struct {
	encryptor IEncryptor
}

func NewWriter() *Writer {
	return &Writer{}
}

func (self *Writer) SetEncryptor(v IEncryptor) {
	self.encryptor = v
}

func (self *Writer) Write(p Packet) (buf []byte) {
	if self.encryptor != nil {
		self.encryptor.Encrypt(p, p)
	}

	return p
}

// ============================================================================

type Packet []byte

func (self Packet) Op() uint32 {
	return binary.BigEndian.Uint32(self[4:8])
}

func (self Packet) Body() []byte {
	return self[8:]
}

func (self *Packet) Add_B8(b8 uint64) {
	l := len(*self)

	// check
	if l+8 > cap(*self) {
		return
	}

	// add
	*self = (*self)[:l+8]
	binary.BigEndian.PutUint64((*self)[l:l+8], b8)
	binary.BigEndian.PutUint32((*self)[:4], uint32(l+8))
}

func (self *Packet) Remove_B8() uint64 {
	l := len(*self)

	// check
	if l-8 < 8 {
		return 0
	}

	// remove
	b8 := binary.BigEndian.Uint64((*self)[l-8:])

	*self = (*self)[:l-8]
	binary.BigEndian.PutUint32((*self)[:4], uint32(l-8))

	return b8
}

func (self *Packet) Remove_B8_Str() (b8 uint64, str string) {
	l := len(*self)

	// check
	if l-9 < 8 {
		return
	}

	// remove
	lstr := int((*self)[l-1])

	b8 = binary.BigEndian.Uint64((*self)[l-1-lstr-8:])
	str = string((*self)[l-1-lstr : l-1])

	*self = (*self)[:l-1-lstr-8]
	binary.BigEndian.PutUint32((*self)[:4], uint32(l-1-lstr-8))

	return
}

func (self *Packet) Peek_B8() uint64 {
	l := len(*self)

	// check
	if l-8 < 8 {
		return 0
	}

	// peek
	return binary.BigEndian.Uint64((*self)[l-8:])
}

func (self *Packet) Peek_B8_Str() (b8 uint64, str string) {
	l := len(*self)

	// check
	if l-9 < 8 {
		return
	}

	// peek
	lstr := int((*self)[l-1])

	b8 = binary.BigEndian.Uint64((*self)[l-1-lstr-8:])
	str = string((*self)[l-1-lstr : l-1])

	return
}

// ============================================================================

func Assemble(op uint32, body []byte) Packet {
	l := 4 + 4 + len(body)

	buf := make([]byte, l)

	binary.BigEndian.PutUint32(buf[:4], uint32(l))
	binary.BigEndian.PutUint32(buf[4:8], op)
	copy(buf[8:], body)

	return Packet(buf)
}

func Assemble_B8(op uint32, body []byte, b8 uint64) Packet {
	l := 4 + 4 + len(body) + 8

	buf := make([]byte, l)

	binary.BigEndian.PutUint32(buf[:4], uint32(l))
	binary.BigEndian.PutUint32(buf[4:8], op)
	copy(buf[8:], body)
	binary.BigEndian.PutUint64(buf[l-8:], b8)

	return Packet(buf)
}

func Assemble_B8_Str(op uint32, body []byte, b8 uint64, str string) Packet {
	lstr := len(str)
	l := 4 + 4 + len(body) + 8 + lstr + 1

	buf := make([]byte, l)

	binary.BigEndian.PutUint32(buf[:4], uint32(l))
	binary.BigEndian.PutUint32(buf[4:8], op)
	copy(buf[8:], body)
	binary.BigEndian.PutUint64(buf[l-1-lstr-8:], b8)
	copy(buf[l-1-lstr:], []byte(str))
	buf[l-1] = byte(lstr)

	return Packet(buf)
}
