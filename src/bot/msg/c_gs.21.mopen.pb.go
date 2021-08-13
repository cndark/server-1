// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: c_gs.21.mopen.proto

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

type MOpenData struct {
	M map[int32]bool `protobuf:"bytes,1,rep,name=M" json:"M,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (m *MOpenData) Reset()                    { *m = MOpenData{} }
func (m *MOpenData) String() string            { return proto.CompactTextString(m) }
func (*MOpenData) ProtoMessage()               {}
func (*MOpenData) Descriptor() ([]byte, []int) { return fileDescriptorCGs_21Mopen, []int{0} }

// 推送
type GS_MOpenModuleNew struct {
	MId     int32    `protobuf:"varint,1,opt,name=MId,proto3" json:"MId,omitempty"`
	Rewards *Rewards `protobuf:"bytes,2,opt,name=Rewards" json:"Rewards,omitempty"`
}

func (m *GS_MOpenModuleNew) Reset()                    { *m = GS_MOpenModuleNew{} }
func (m *GS_MOpenModuleNew) String() string            { return proto.CompactTextString(m) }
func (*GS_MOpenModuleNew) ProtoMessage()               {}
func (*GS_MOpenModuleNew) Descriptor() ([]byte, []int) { return fileDescriptorCGs_21Mopen, []int{1} }

func init() {
	proto.RegisterType((*MOpenData)(nil), "msg.MOpenData")
	proto.RegisterType((*GS_MOpenModuleNew)(nil), "msg.GS_MOpenModuleNew")
}
func (m *MOpenData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MOpenData) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.M) > 0 {
		for k, _ := range m.M {
			dAtA[i] = 0xa
			i++
			v := m.M[k]
			mapSize := 1 + sovCGs_21Mopen(uint64(k)) + 1 + 1
			i = encodeVarintCGs_21Mopen(dAtA, i, uint64(mapSize))
			dAtA[i] = 0x8
			i++
			i = encodeVarintCGs_21Mopen(dAtA, i, uint64(k))
			dAtA[i] = 0x10
			i++
			if v {
				dAtA[i] = 1
			} else {
				dAtA[i] = 0
			}
			i++
		}
	}
	return i, nil
}

func (m *GS_MOpenModuleNew) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GS_MOpenModuleNew) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.MId != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCGs_21Mopen(dAtA, i, uint64(m.MId))
	}
	if m.Rewards != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintCGs_21Mopen(dAtA, i, uint64(m.Rewards.Size()))
		n1, err := m.Rewards.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func encodeVarintCGs_21Mopen(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *MOpenData) Size() (n int) {
	var l int
	_ = l
	if len(m.M) > 0 {
		for k, v := range m.M {
			_ = k
			_ = v
			mapEntrySize := 1 + sovCGs_21Mopen(uint64(k)) + 1 + 1
			n += mapEntrySize + 1 + sovCGs_21Mopen(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *GS_MOpenModuleNew) Size() (n int) {
	var l int
	_ = l
	if m.MId != 0 {
		n += 1 + sovCGs_21Mopen(uint64(m.MId))
	}
	if m.Rewards != nil {
		l = m.Rewards.Size()
		n += 1 + l + sovCGs_21Mopen(uint64(l))
	}
	return n
}

func sovCGs_21Mopen(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCGs_21Mopen(x uint64) (n int) {
	return sovCGs_21Mopen(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MOpenData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_21Mopen
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
			return fmt.Errorf("proto: MOpenData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MOpenData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field M", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_21Mopen
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
				return ErrInvalidLengthCGs_21Mopen
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.M == nil {
				m.M = make(map[int32]bool)
			}
			var mapkey int32
			var mapvalue bool
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_21Mopen
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
				if fieldNum == 1 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_21Mopen
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= (int32(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					var mapvaluetemp int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_21Mopen
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvaluetemp |= (int(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					mapvalue = bool(mapvaluetemp != 0)
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipCGs_21Mopen(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthCGs_21Mopen
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.M[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_21Mopen(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_21Mopen
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
func (m *GS_MOpenModuleNew) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_21Mopen
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
			return fmt.Errorf("proto: GS_MOpenModuleNew: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GS_MOpenModuleNew: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MId", wireType)
			}
			m.MId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_21Mopen
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MId |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rewards", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_21Mopen
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
				return ErrInvalidLengthCGs_21Mopen
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
			skippy, err := skipCGs_21Mopen(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_21Mopen
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
func skipCGs_21Mopen(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCGs_21Mopen
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
					return 0, ErrIntOverflowCGs_21Mopen
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
					return 0, ErrIntOverflowCGs_21Mopen
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
				return 0, ErrInvalidLengthCGs_21Mopen
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCGs_21Mopen
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
				next, err := skipCGs_21Mopen(dAtA[start:])
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
	ErrInvalidLengthCGs_21Mopen = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCGs_21Mopen   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("c_gs.21.mopen.proto", fileDescriptorCGs_21Mopen) }

var fileDescriptorCGs_21Mopen = []byte{
	// 240 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4e, 0x8e, 0x4f, 0x2f,
	0xd6, 0x33, 0x32, 0xd4, 0xcb, 0xcd, 0x2f, 0x48, 0xcd, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0x62, 0xce, 0x2d, 0x4e, 0x97, 0xe2, 0x4a, 0xcf, 0x4f, 0xcf, 0x87, 0x08, 0x48, 0x89, 0x82, 0x55,
	0x19, 0x18, 0xe8, 0x15, 0x97, 0x14, 0x95, 0x26, 0x97, 0x14, 0x43, 0x84, 0x95, 0xd2, 0xb8, 0x38,
	0x7d, 0xfd, 0x0b, 0x52, 0xf3, 0x5c, 0x12, 0x4b, 0x12, 0x85, 0x94, 0xb9, 0x18, 0x7d, 0x25, 0x18,
	0x15, 0x98, 0x35, 0xb8, 0x8d, 0x44, 0xf5, 0x72, 0x8b, 0xd3, 0xf5, 0xe0, 0x52, 0x7a, 0xbe, 0xae,
	0x79, 0x25, 0x45, 0x95, 0x41, 0x8c, 0xbe, 0x52, 0x26, 0x5c, 0x6c, 0x10, 0x8e, 0x90, 0x00, 0x17,
	0x73, 0x76, 0x6a, 0xa5, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x6b, 0x10, 0x88, 0x29, 0x24, 0xc2, 0xc5,
	0x5a, 0x96, 0x98, 0x53, 0x9a, 0x2a, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x11, 0x04, 0xe1, 0x58, 0x31,
	0x59, 0x30, 0x2a, 0xf9, 0x72, 0x09, 0xba, 0x07, 0xc7, 0x83, 0xcd, 0xf3, 0xcd, 0x4f, 0x29, 0xcd,
	0x49, 0xf5, 0x4b, 0x2d, 0x07, 0x19, 0xe0, 0xeb, 0x99, 0x02, 0x33, 0xc0, 0xd7, 0x33, 0x45, 0x48,
	0x8d, 0x8b, 0x3d, 0x28, 0xb5, 0x3c, 0xb1, 0x28, 0xa5, 0x18, 0x6c, 0x04, 0xb7, 0x11, 0x0f, 0xd8,
	0x1d, 0x50, 0xb1, 0x20, 0x98, 0xa4, 0x93, 0xc8, 0x89, 0x87, 0x72, 0x0c, 0x27, 0x1e, 0xc9, 0x31,
	0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x84, 0xc7, 0x72, 0x0c, 0x49, 0x6c, 0x60,
	0x3f, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xbb, 0xe1, 0xc0, 0xa5, 0x12, 0x01, 0x00, 0x00,
}
