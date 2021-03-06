// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: c_gs.54.lamp.proto

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

// 跑马灯数据
type LampData struct {
	Data []*LampOne `protobuf:"bytes,1,rep,name=Data" json:"Data,omitempty"`
}

func (m *LampData) Reset()                    { *m = LampData{} }
func (m *LampData) String() string            { return proto.CompactTextString(m) }
func (*LampData) ProtoMessage()               {}
func (*LampData) Descriptor() ([]byte, []int) { return fileDescriptorCGs_54Lamp, []int{0} }

type LampOne struct {
	Id    int32             `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Param map[string]string `protobuf:"bytes,2,rep,name=Param" json:"Param,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Ts    int64             `protobuf:"varint,3,opt,name=Ts,proto3" json:"Ts,omitempty"`
}

func (m *LampOne) Reset()                    { *m = LampOne{} }
func (m *LampOne) String() string            { return proto.CompactTextString(m) }
func (*LampOne) ProtoMessage()               {}
func (*LampOne) Descriptor() ([]byte, []int) { return fileDescriptorCGs_54Lamp, []int{1} }

// ============================================================================
// 通知
type GS_LampMsg struct {
	One *LampOne `protobuf:"bytes,1,opt,name=One" json:"One,omitempty"`
}

func (m *GS_LampMsg) Reset()                    { *m = GS_LampMsg{} }
func (m *GS_LampMsg) String() string            { return proto.CompactTextString(m) }
func (*GS_LampMsg) ProtoMessage()               {}
func (*GS_LampMsg) Descriptor() ([]byte, []int) { return fileDescriptorCGs_54Lamp, []int{2} }

func init() {
	proto.RegisterType((*LampData)(nil), "msg.LampData")
	proto.RegisterType((*LampOne)(nil), "msg.LampOne")
	proto.RegisterType((*GS_LampMsg)(nil), "msg.GS_LampMsg")
}
func (m *LampData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LampData) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		for _, msg := range m.Data {
			dAtA[i] = 0xa
			i++
			i = encodeVarintCGs_54Lamp(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *LampOne) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LampOne) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCGs_54Lamp(dAtA, i, uint64(m.Id))
	}
	if len(m.Param) > 0 {
		for k, _ := range m.Param {
			dAtA[i] = 0x12
			i++
			v := m.Param[k]
			mapSize := 1 + len(k) + sovCGs_54Lamp(uint64(len(k))) + 1 + len(v) + sovCGs_54Lamp(uint64(len(v)))
			i = encodeVarintCGs_54Lamp(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintCGs_54Lamp(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintCGs_54Lamp(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if m.Ts != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintCGs_54Lamp(dAtA, i, uint64(m.Ts))
	}
	return i, nil
}

func (m *GS_LampMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GS_LampMsg) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.One != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCGs_54Lamp(dAtA, i, uint64(m.One.Size()))
		n1, err := m.One.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func encodeVarintCGs_54Lamp(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *LampData) Size() (n int) {
	var l int
	_ = l
	if len(m.Data) > 0 {
		for _, e := range m.Data {
			l = e.Size()
			n += 1 + l + sovCGs_54Lamp(uint64(l))
		}
	}
	return n
}

func (m *LampOne) Size() (n int) {
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovCGs_54Lamp(uint64(m.Id))
	}
	if len(m.Param) > 0 {
		for k, v := range m.Param {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovCGs_54Lamp(uint64(len(k))) + 1 + len(v) + sovCGs_54Lamp(uint64(len(v)))
			n += mapEntrySize + 1 + sovCGs_54Lamp(uint64(mapEntrySize))
		}
	}
	if m.Ts != 0 {
		n += 1 + sovCGs_54Lamp(uint64(m.Ts))
	}
	return n
}

func (m *GS_LampMsg) Size() (n int) {
	var l int
	_ = l
	if m.One != nil {
		l = m.One.Size()
		n += 1 + l + sovCGs_54Lamp(uint64(l))
	}
	return n
}

func sovCGs_54Lamp(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCGs_54Lamp(x uint64) (n int) {
	return sovCGs_54Lamp(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LampData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_54Lamp
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
			return fmt.Errorf("proto: LampData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LampData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_54Lamp
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
				return ErrInvalidLengthCGs_54Lamp
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data, &LampOne{})
			if err := m.Data[len(m.Data)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_54Lamp(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_54Lamp
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
func (m *LampOne) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_54Lamp
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
			return fmt.Errorf("proto: LampOne: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LampOne: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_54Lamp
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Param", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_54Lamp
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
				return ErrInvalidLengthCGs_54Lamp
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Param == nil {
				m.Param = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_54Lamp
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
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_54Lamp
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthCGs_54Lamp
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_54Lamp
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthCGs_54Lamp
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipCGs_54Lamp(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthCGs_54Lamp
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Param[mapkey] = mapvalue
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ts", wireType)
			}
			m.Ts = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_54Lamp
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Ts |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_54Lamp(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_54Lamp
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
func (m *GS_LampMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_54Lamp
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
			return fmt.Errorf("proto: GS_LampMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GS_LampMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field One", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_54Lamp
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
				return ErrInvalidLengthCGs_54Lamp
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.One == nil {
				m.One = &LampOne{}
			}
			if err := m.One.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_54Lamp(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_54Lamp
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
func skipCGs_54Lamp(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCGs_54Lamp
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
					return 0, ErrIntOverflowCGs_54Lamp
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
					return 0, ErrIntOverflowCGs_54Lamp
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
				return 0, ErrInvalidLengthCGs_54Lamp
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCGs_54Lamp
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
				next, err := skipCGs_54Lamp(dAtA[start:])
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
	ErrInvalidLengthCGs_54Lamp = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCGs_54Lamp   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("c_gs.54.lamp.proto", fileDescriptorCGs_54Lamp) }

var fileDescriptorCGs_54Lamp = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4a, 0x8e, 0x4f, 0x2f,
	0xd6, 0x33, 0x35, 0xd1, 0xcb, 0x49, 0xcc, 0x2d, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0xce, 0x2d, 0x4e, 0x97, 0xe2, 0x4a, 0xcf, 0x4f, 0xcf, 0x87, 0x08, 0x28, 0xe9, 0x70, 0x71, 0xf8,
	0x24, 0xe6, 0x16, 0xb8, 0x24, 0x96, 0x24, 0x0a, 0x29, 0x70, 0xb1, 0x80, 0x68, 0x09, 0x46, 0x05,
	0x66, 0x0d, 0x6e, 0x23, 0x1e, 0xbd, 0xdc, 0xe2, 0x74, 0x3d, 0x90, 0xa4, 0x7f, 0x5e, 0x6a, 0x10,
	0x58, 0x46, 0x69, 0x12, 0x23, 0x17, 0x3b, 0x54, 0x44, 0x88, 0x8f, 0x8b, 0xc9, 0x33, 0x45, 0x82,
	0x51, 0x81, 0x51, 0x83, 0x35, 0x88, 0xc9, 0x33, 0x45, 0x48, 0x97, 0x8b, 0x35, 0x20, 0xb1, 0x28,
	0x31, 0x57, 0x82, 0x09, 0xac, 0x5d, 0x1c, 0x59, 0xbb, 0x1e, 0x58, 0xc6, 0x35, 0xaf, 0xa4, 0xa8,
	0x32, 0x08, 0xa2, 0x0a, 0xa4, 0x3d, 0xa4, 0x58, 0x82, 0x59, 0x81, 0x51, 0x83, 0x39, 0x88, 0x29,
	0xa4, 0x58, 0xca, 0x82, 0x8b, 0x0b, 0xa1, 0x48, 0x48, 0x80, 0x8b, 0x39, 0x3b, 0xb5, 0x12, 0x6c,
	0x3a, 0x67, 0x10, 0x88, 0x29, 0x24, 0xc2, 0xc5, 0x5a, 0x96, 0x98, 0x53, 0x9a, 0x2a, 0xc1, 0x04,
	0x16, 0x83, 0x70, 0xac, 0x98, 0x2c, 0x18, 0x95, 0x74, 0xb8, 0xb8, 0xdc, 0x83, 0xe3, 0x41, 0x36,
	0xf9, 0x16, 0xa7, 0x0b, 0xc9, 0x71, 0x31, 0xfb, 0xe7, 0xa5, 0x82, 0x75, 0xa2, 0xfb, 0x01, 0x24,
	0xe1, 0x24, 0x72, 0xe2, 0xa1, 0x1c, 0xc3, 0x89, 0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e,
	0x78, 0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0x43, 0x12, 0x1b, 0x38, 0x34, 0x8c, 0x01, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xe2, 0xc5, 0x28, 0xd1, 0x34, 0x01, 0x00, 0x00,
}
