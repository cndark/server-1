// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: c_gs.45.privcard.proto

package msg

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type PrivCardData struct {
	Cards []*PrivCard `protobuf:"bytes,1,rep,name=Cards" json:"Cards,omitempty"`
}

func (m *PrivCardData) Reset()                    { *m = PrivCardData{} }
func (m *PrivCardData) String() string            { return proto.CompactTextString(m) }
func (*PrivCardData) ProtoMessage()               {}
func (*PrivCardData) Descriptor() ([]byte, []int) { return fileDescriptorCGs_45Privcard, []int{0} }

type PrivCard struct {
	Id       int32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	ExpireTs int64 `protobuf:"varint,2,opt,name=ExpireTs,proto3" json:"ExpireTs,omitempty"`
	IsAward  bool  `protobuf:"varint,3,opt,name=IsAward,proto3" json:"IsAward,omitempty"`
	AddCnt   int32 `protobuf:"varint,4,opt,name=AddCnt,proto3" json:"AddCnt,omitempty"`
}

func (m *PrivCard) Reset()                    { *m = PrivCard{} }
func (m *PrivCard) String() string            { return proto.CompactTextString(m) }
func (*PrivCard) ProtoMessage()               {}
func (*PrivCard) Descriptor() ([]byte, []int) { return fileDescriptorCGs_45Privcard, []int{1} }

type GS_PrivCardNew struct {
	Card    *PrivCard `protobuf:"bytes,1,opt,name=Card" json:"Card,omitempty"`
	Rewards *Rewards  `protobuf:"bytes,2,opt,name=Rewards" json:"Rewards,omitempty"`
}

func (m *GS_PrivCardNew) Reset()                    { *m = GS_PrivCardNew{} }
func (m *GS_PrivCardNew) String() string            { return proto.CompactTextString(m) }
func (*GS_PrivCardNew) ProtoMessage()               {}
func (*GS_PrivCardNew) Descriptor() ([]byte, []int) { return fileDescriptorCGs_45Privcard, []int{2} }

// ============================================================================
type C_PrivCardTake struct {
	Id int32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (m *C_PrivCardTake) Reset()                    { *m = C_PrivCardTake{} }
func (m *C_PrivCardTake) String() string            { return proto.CompactTextString(m) }
func (*C_PrivCardTake) ProtoMessage()               {}
func (*C_PrivCardTake) Descriptor() ([]byte, []int) { return fileDescriptorCGs_45Privcard, []int{3} }

type GS_PrivCardTake_R struct {
	ErrorCode int32    `protobuf:"varint,1,opt,name=ErrorCode,proto3" json:"ErrorCode,omitempty"`
	Id        int32    `protobuf:"varint,2,opt,name=Id,proto3" json:"Id,omitempty"`
	Rewards   *Rewards `protobuf:"bytes,3,opt,name=Rewards" json:"Rewards,omitempty"`
}

func (m *GS_PrivCardTake_R) Reset()                    { *m = GS_PrivCardTake_R{} }
func (m *GS_PrivCardTake_R) String() string            { return proto.CompactTextString(m) }
func (*GS_PrivCardTake_R) ProtoMessage()               {}
func (*GS_PrivCardTake_R) Descriptor() ([]byte, []int) { return fileDescriptorCGs_45Privcard, []int{4} }

func init() {
	proto.RegisterType((*PrivCardData)(nil), "msg.PrivCardData")
	proto.RegisterType((*PrivCard)(nil), "msg.PrivCard")
	proto.RegisterType((*GS_PrivCardNew)(nil), "msg.GS_PrivCardNew")
	proto.RegisterType((*C_PrivCardTake)(nil), "msg.C_PrivCardTake")
	proto.RegisterType((*GS_PrivCardTake_R)(nil), "msg.GS_PrivCardTake_R")
}
func (m *PrivCardData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PrivCardData) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Cards) > 0 {
		for _, msg := range m.Cards {
			dAtA[i] = 0xa
			i++
			i = encodeVarintCGs_45Privcard(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *PrivCard) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PrivCard) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCGs_45Privcard(dAtA, i, uint64(m.Id))
	}
	if m.ExpireTs != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintCGs_45Privcard(dAtA, i, uint64(m.ExpireTs))
	}
	if m.IsAward {
		dAtA[i] = 0x18
		i++
		if m.IsAward {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.AddCnt != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintCGs_45Privcard(dAtA, i, uint64(m.AddCnt))
	}
	return i, nil
}

func (m *GS_PrivCardNew) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GS_PrivCardNew) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Card != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCGs_45Privcard(dAtA, i, uint64(m.Card.Size()))
		n1, err := m.Card.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.Rewards != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintCGs_45Privcard(dAtA, i, uint64(m.Rewards.Size()))
		n2, err := m.Rewards.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func (m *C_PrivCardTake) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *C_PrivCardTake) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCGs_45Privcard(dAtA, i, uint64(m.Id))
	}
	return i, nil
}

func (m *GS_PrivCardTake_R) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GS_PrivCardTake_R) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ErrorCode != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCGs_45Privcard(dAtA, i, uint64(m.ErrorCode))
	}
	if m.Id != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintCGs_45Privcard(dAtA, i, uint64(m.Id))
	}
	if m.Rewards != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintCGs_45Privcard(dAtA, i, uint64(m.Rewards.Size()))
		n3, err := m.Rewards.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}

func encodeVarintCGs_45Privcard(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *PrivCardData) Size() (n int) {
	var l int
	_ = l
	if len(m.Cards) > 0 {
		for _, e := range m.Cards {
			l = e.Size()
			n += 1 + l + sovCGs_45Privcard(uint64(l))
		}
	}
	return n
}

func (m *PrivCard) Size() (n int) {
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovCGs_45Privcard(uint64(m.Id))
	}
	if m.ExpireTs != 0 {
		n += 1 + sovCGs_45Privcard(uint64(m.ExpireTs))
	}
	if m.IsAward {
		n += 2
	}
	if m.AddCnt != 0 {
		n += 1 + sovCGs_45Privcard(uint64(m.AddCnt))
	}
	return n
}

func (m *GS_PrivCardNew) Size() (n int) {
	var l int
	_ = l
	if m.Card != nil {
		l = m.Card.Size()
		n += 1 + l + sovCGs_45Privcard(uint64(l))
	}
	if m.Rewards != nil {
		l = m.Rewards.Size()
		n += 1 + l + sovCGs_45Privcard(uint64(l))
	}
	return n
}

func (m *C_PrivCardTake) Size() (n int) {
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovCGs_45Privcard(uint64(m.Id))
	}
	return n
}

func (m *GS_PrivCardTake_R) Size() (n int) {
	var l int
	_ = l
	if m.ErrorCode != 0 {
		n += 1 + sovCGs_45Privcard(uint64(m.ErrorCode))
	}
	if m.Id != 0 {
		n += 1 + sovCGs_45Privcard(uint64(m.Id))
	}
	if m.Rewards != nil {
		l = m.Rewards.Size()
		n += 1 + l + sovCGs_45Privcard(uint64(l))
	}
	return n
}

func sovCGs_45Privcard(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCGs_45Privcard(x uint64) (n int) {
	return sovCGs_45Privcard(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PrivCardData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_45Privcard
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PrivCardData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PrivCardData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cards", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCGs_45Privcard
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Cards = append(m.Cards, &PrivCard{})
			if err := m.Cards[len(m.Cards)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_45Privcard(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_45Privcard
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PrivCard) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_45Privcard
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PrivCard: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PrivCard: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpireTs", wireType)
			}
			m.ExpireTs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExpireTs |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsAward", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsAward = bool(v != 0)
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddCnt", wireType)
			}
			m.AddCnt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AddCnt |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_45Privcard(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_45Privcard
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GS_PrivCardNew) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_45Privcard
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GS_PrivCardNew: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GS_PrivCardNew: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Card", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCGs_45Privcard
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Card == nil {
				m.Card = &PrivCard{}
			}
			if err := m.Card.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rewards", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCGs_45Privcard
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Rewards == nil {
				m.Rewards = &Rewards{}
			}
			if err := m.Rewards.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_45Privcard(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_45Privcard
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *C_PrivCardTake) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_45Privcard
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: C_PrivCardTake: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: C_PrivCardTake: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_45Privcard(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_45Privcard
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GS_PrivCardTake_R) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_45Privcard
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GS_PrivCardTake_R: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GS_PrivCardTake_R: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrorCode", wireType)
			}
			m.ErrorCode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ErrorCode |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rewards", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCGs_45Privcard
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Rewards == nil {
				m.Rewards = &Rewards{}
			}
			if err := m.Rewards.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_45Privcard(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_45Privcard
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipCGs_45Privcard(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCGs_45Privcard
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCGs_45Privcard
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthCGs_45Privcard
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCGs_45Privcard
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipCGs_45Privcard(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthCGs_45Privcard = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCGs_45Privcard   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("c_gs.45.privcard.proto", fileDescriptorCGs_45Privcard) }

var fileDescriptorCGs_45Privcard = []byte{
	// 317 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xdf, 0x4a, 0xc3, 0x30,
	0x14, 0xc6, 0x97, 0x76, 0xff, 0x3c, 0xce, 0x81, 0x41, 0x47, 0x18, 0x52, 0x6a, 0x05, 0xe9, 0x55,
	0x18, 0x4e, 0x1f, 0x60, 0xd6, 0x21, 0xbb, 0x11, 0x89, 0xbb, 0xf3, 0x62, 0xc4, 0xa5, 0xd4, 0x22,
	0x33, 0x23, 0xa9, 0x9b, 0x8f, 0xe1, 0x63, 0xed, 0xd2, 0x47, 0xd0, 0xfa, 0x22, 0xd2, 0xf4, 0x8f,
	0x43, 0xc1, 0xbb, 0x73, 0xbe, 0xf3, 0xe3, 0x7c, 0xdf, 0x49, 0xa0, 0x37, 0x9f, 0x45, 0x9a, 0x9e,
	0x5f, 0xd0, 0xa5, 0x8a, 0x57, 0x73, 0xae, 0x04, 0x5d, 0x2a, 0x99, 0x48, 0x6c, 0x2f, 0x74, 0xd4,
	0x87, 0x48, 0x46, 0x32, 0x17, 0xfa, 0x87, 0x06, 0x1c, 0x0c, 0xa8, 0x4e, 0xd4, 0xcb, 0x3c, 0xd1,
	0xb9, 0xec, 0x0d, 0xa1, 0x73, 0xab, 0xe2, 0x55, 0xc0, 0x95, 0xb8, 0xe2, 0x09, 0xc7, 0x27, 0xd0,
	0xc8, 0x6a, 0x4d, 0x90, 0x6b, 0xfb, 0xbb, 0x67, 0x7b, 0x74, 0xa1, 0x23, 0x5a, 0x12, 0x2c, 0x9f,
	0x79, 0x8f, 0xd0, 0x2e, 0x25, 0xdc, 0x05, 0x6b, 0x22, 0x08, 0x72, 0x91, 0xdf, 0x60, 0xd6, 0x44,
	0xe0, 0x3e, 0xb4, 0xc7, 0xaf, 0xcb, 0x58, 0x85, 0x53, 0x4d, 0x2c, 0x17, 0xf9, 0x36, 0xab, 0x7a,
	0x4c, 0xa0, 0x35, 0xd1, 0xa3, 0x35, 0x57, 0x82, 0xd8, 0x2e, 0xf2, 0xdb, 0xac, 0x6c, 0x71, 0x0f,
	0x9a, 0x23, 0x21, 0x82, 0xe7, 0x84, 0xd4, 0xcd, 0xa6, 0xa2, 0xf3, 0xee, 0xa1, 0x7b, 0x7d, 0x37,
	0x2b, 0xcd, 0x6e, 0xc2, 0x35, 0x3e, 0x86, 0x7a, 0x56, 0x1a, 0xc7, 0x3f, 0xf9, 0xcc, 0x08, 0x9f,
	0x42, 0x8b, 0x85, 0x6b, 0x73, 0x85, 0x65, 0xa8, 0x8e, 0xa1, 0x0a, 0x8d, 0x95, 0x43, 0xcf, 0x85,
	0x6e, 0x50, 0xed, 0x9e, 0xf2, 0xa7, 0xf0, 0xf7, 0x31, 0x5e, 0x0c, 0xfb, 0x5b, 0xf6, 0x19, 0x32,
	0x63, 0xf8, 0x08, 0x76, 0xc6, 0x4a, 0x49, 0x15, 0x48, 0x11, 0x16, 0xec, 0x8f, 0x50, 0xac, 0xb0,
	0xaa, 0xf7, 0xd8, 0x0a, 0x63, 0xff, 0x13, 0xe6, 0xf2, 0x60, 0xf3, 0xe9, 0xd4, 0x36, 0xa9, 0x83,
	0xde, 0x53, 0x07, 0x7d, 0xa4, 0x0e, 0x7a, 0xfb, 0x72, 0x6a, 0x0f, 0x4d, 0xf3, 0x4b, 0xc3, 0xef,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xfb, 0x4b, 0x3c, 0x21, 0xe7, 0x01, 0x00, 0x00,
}
