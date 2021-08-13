// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: c_gs.16.act.proto

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

type ActState struct {
	Name    string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Stage   string `protobuf:"bytes,2,opt,name=Stage,proto3" json:"Stage,omitempty"`
	T1      int64  `protobuf:"varint,3,opt,name=T1,proto3" json:"T1,omitempty"`
	T2      int64  `protobuf:"varint,4,opt,name=T2,proto3" json:"T2,omitempty"`
	ConfGrp int32  `protobuf:"varint,5,opt,name=ConfGrp,proto3" json:"ConfGrp,omitempty"`
}

func (m *ActState) Reset()                    { *m = ActState{} }
func (m *ActState) String() string            { return proto.CompactTextString(m) }
func (*ActState) ProtoMessage()               {}
func (*ActState) Descriptor() ([]byte, []int) { return fileDescriptorCGs_16Act, []int{0} }

type GS_ActStateChange struct {
	Act *ActState `protobuf:"bytes,1,opt,name=Act" json:"Act,omitempty"`
}

func (m *GS_ActStateChange) Reset()                    { *m = GS_ActStateChange{} }
func (m *GS_ActStateChange) String() string            { return proto.CompactTextString(m) }
func (*GS_ActStateChange) ProtoMessage()               {}
func (*GS_ActStateChange) Descriptor() ([]byte, []int) { return fileDescriptorCGs_16Act, []int{1} }

type C_ActStateGet struct {
}

func (m *C_ActStateGet) Reset()                    { *m = C_ActStateGet{} }
func (m *C_ActStateGet) String() string            { return proto.CompactTextString(m) }
func (*C_ActStateGet) ProtoMessage()               {}
func (*C_ActStateGet) Descriptor() ([]byte, []int) { return fileDescriptorCGs_16Act, []int{2} }

type GS_ActStateGet_R struct {
	ErrorCode int32       `protobuf:"varint,1,opt,name=ErrorCode,proto3" json:"ErrorCode,omitempty"`
	Acts      []*ActState `protobuf:"bytes,2,rep,name=Acts" json:"Acts,omitempty"`
}

func (m *GS_ActStateGet_R) Reset()                    { *m = GS_ActStateGet_R{} }
func (m *GS_ActStateGet_R) String() string            { return proto.CompactTextString(m) }
func (*GS_ActStateGet_R) ProtoMessage()               {}
func (*GS_ActStateGet_R) Descriptor() ([]byte, []int) { return fileDescriptorCGs_16Act, []int{3} }

func init() {
	proto.RegisterType((*ActState)(nil), "msg.ActState")
	proto.RegisterType((*GS_ActStateChange)(nil), "msg.GS_ActStateChange")
	proto.RegisterType((*C_ActStateGet)(nil), "msg.C_ActStateGet")
	proto.RegisterType((*GS_ActStateGet_R)(nil), "msg.GS_ActStateGet_R")
}
func (m *ActState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ActState) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCGs_16Act(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.Stage) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintCGs_16Act(dAtA, i, uint64(len(m.Stage)))
		i += copy(dAtA[i:], m.Stage)
	}
	if m.T1 != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintCGs_16Act(dAtA, i, uint64(m.T1))
	}
	if m.T2 != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintCGs_16Act(dAtA, i, uint64(m.T2))
	}
	if m.ConfGrp != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintCGs_16Act(dAtA, i, uint64(m.ConfGrp))
	}
	return i, nil
}

func (m *GS_ActStateChange) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GS_ActStateChange) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Act != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCGs_16Act(dAtA, i, uint64(m.Act.Size()))
		n1, err := m.Act.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func (m *C_ActStateGet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *C_ActStateGet) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *GS_ActStateGet_R) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GS_ActStateGet_R) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ErrorCode != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCGs_16Act(dAtA, i, uint64(m.ErrorCode))
	}
	if len(m.Acts) > 0 {
		for _, msg := range m.Acts {
			dAtA[i] = 0x12
			i++
			i = encodeVarintCGs_16Act(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeVarintCGs_16Act(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ActState) Size() (n int) {
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovCGs_16Act(uint64(l))
	}
	l = len(m.Stage)
	if l > 0 {
		n += 1 + l + sovCGs_16Act(uint64(l))
	}
	if m.T1 != 0 {
		n += 1 + sovCGs_16Act(uint64(m.T1))
	}
	if m.T2 != 0 {
		n += 1 + sovCGs_16Act(uint64(m.T2))
	}
	if m.ConfGrp != 0 {
		n += 1 + sovCGs_16Act(uint64(m.ConfGrp))
	}
	return n
}

func (m *GS_ActStateChange) Size() (n int) {
	var l int
	_ = l
	if m.Act != nil {
		l = m.Act.Size()
		n += 1 + l + sovCGs_16Act(uint64(l))
	}
	return n
}

func (m *C_ActStateGet) Size() (n int) {
	var l int
	_ = l
	return n
}

func (m *GS_ActStateGet_R) Size() (n int) {
	var l int
	_ = l
	if m.ErrorCode != 0 {
		n += 1 + sovCGs_16Act(uint64(m.ErrorCode))
	}
	if len(m.Acts) > 0 {
		for _, e := range m.Acts {
			l = e.Size()
			n += 1 + l + sovCGs_16Act(uint64(l))
		}
	}
	return n
}

func sovCGs_16Act(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCGs_16Act(x uint64) (n int) {
	return sovCGs_16Act(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ActState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_16Act
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
			return fmt.Errorf("proto: ActState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ActState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_16Act
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCGs_16Act
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_16Act
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCGs_16Act
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Stage = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field T1", wireType)
			}
			m.T1 = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_16Act
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.T1 |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field T2", wireType)
			}
			m.T2 = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_16Act
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.T2 |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConfGrp", wireType)
			}
			m.ConfGrp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_16Act
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ConfGrp |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_16Act(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_16Act
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
func (m *GS_ActStateChange) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_16Act
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
			return fmt.Errorf("proto: GS_ActStateChange: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GS_ActStateChange: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Act", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_16Act
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
				return ErrInvalidLengthCGs_16Act
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Act == nil {
				m.Act = &ActState{}
			}
			if err := m.Act.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_16Act(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_16Act
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
func (m *C_ActStateGet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_16Act
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
			return fmt.Errorf("proto: C_ActStateGet: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: C_ActStateGet: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_16Act(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_16Act
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
func (m *GS_ActStateGet_R) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_16Act
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
			return fmt.Errorf("proto: GS_ActStateGet_R: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GS_ActStateGet_R: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrorCode", wireType)
			}
			m.ErrorCode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_16Act
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Acts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_16Act
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
				return ErrInvalidLengthCGs_16Act
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Acts = append(m.Acts, &ActState{})
			if err := m.Acts[len(m.Acts)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_16Act(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_16Act
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
func skipCGs_16Act(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCGs_16Act
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
					return 0, ErrIntOverflowCGs_16Act
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
					return 0, ErrIntOverflowCGs_16Act
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
				return 0, ErrInvalidLengthCGs_16Act
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCGs_16Act
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
				next, err := skipCGs_16Act(dAtA[start:])
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
	ErrInvalidLengthCGs_16Act = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCGs_16Act   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("c_gs.16.act.proto", fileDescriptorCGs_16Act) }

var fileDescriptorCGs_16Act = []byte{
	// 260 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0x3f, 0x4e, 0xf4, 0x30,
	0x10, 0xc5, 0xd7, 0xf9, 0xf3, 0x7d, 0xec, 0xa0, 0x05, 0xd6, 0xda, 0xc2, 0x42, 0x28, 0x84, 0x54,
	0xa9, 0x22, 0x25, 0x20, 0xfa, 0x10, 0xa1, 0x74, 0x14, 0xce, 0xf6, 0x91, 0x31, 0xc6, 0x34, 0x89,
	0x57, 0xce, 0x1c, 0x84, 0x63, 0x6d, 0xc9, 0x11, 0x20, 0x5c, 0x04, 0xe1, 0x10, 0x40, 0xa2, 0x9b,
	0xdf, 0x7b, 0x33, 0xf3, 0x34, 0x03, 0x6b, 0xd9, 0xea, 0x21, 0xcb, 0xaf, 0x33, 0x21, 0x31, 0xdb,
	0x59, 0x83, 0x86, 0xfa, 0xdd, 0xa0, 0x4f, 0x41, 0x1b, 0x6d, 0x26, 0x21, 0xe9, 0xe1, 0xa0, 0x94,
	0xd8, 0xa0, 0x40, 0x45, 0x29, 0x04, 0x77, 0xa2, 0x53, 0x8c, 0xc4, 0x24, 0x5d, 0x72, 0x57, 0xd3,
	0x0d, 0x84, 0x0d, 0x0a, 0xad, 0x98, 0xe7, 0xc4, 0x09, 0xe8, 0x11, 0x78, 0xdb, 0x9c, 0xf9, 0x31,
	0x49, 0x7d, 0xee, 0x6d, 0x73, 0xc7, 0x05, 0x0b, 0xbe, 0xb8, 0xa0, 0x0c, 0xfe, 0x57, 0xa6, 0x7f,
	0xac, 0xed, 0x8e, 0x85, 0x31, 0x49, 0x43, 0x3e, 0x63, 0x72, 0x05, 0xeb, 0xba, 0x69, 0xe7, 0xc8,
	0xea, 0x49, 0xf4, 0x5a, 0xd1, 0x73, 0xf0, 0x4b, 0x89, 0x2e, 0xf7, 0xb0, 0x58, 0x65, 0xdd, 0xa0,
	0xb3, 0xb9, 0x83, 0x7f, 0x3a, 0xc9, 0x31, 0xac, 0xaa, 0xef, 0xa1, 0x5a, 0x61, 0xd2, 0xc0, 0xc9,
	0xaf, 0x35, 0xb5, 0xc2, 0x96, 0xd3, 0x33, 0x58, 0xde, 0x5a, 0x6b, 0x6c, 0x65, 0x1e, 0xa6, 0x1b,
	0x42, 0xfe, 0x23, 0xd0, 0x0b, 0x08, 0x4a, 0x89, 0x03, 0xf3, 0x62, 0xff, 0x6f, 0x88, 0xb3, 0x6e,
	0x36, 0xfb, 0xb7, 0x68, 0xb1, 0x1f, 0x23, 0xf2, 0x32, 0x46, 0xe4, 0x75, 0x8c, 0xc8, 0xf3, 0x7b,
	0xb4, 0xb8, 0xff, 0xe7, 0x1e, 0x75, 0xf9, 0x11, 0x00, 0x00, 0xff, 0xff, 0xc3, 0x7c, 0x99, 0x4c,
	0x4e, 0x01, 0x00, 0x00,
}
