// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: c_gs.32.teammgr.proto

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

// 阵容管理
type TeamMgrData struct {
	Teams map[int32]*TeamFormation `protobuf:"bytes,1,rep,name=Teams" json:"Teams,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *TeamMgrData) Reset()                    { *m = TeamMgrData{} }
func (m *TeamMgrData) String() string            { return proto.CompactTextString(m) }
func (*TeamMgrData) ProtoMessage()               {}
func (*TeamMgrData) Descriptor() ([]byte, []int) { return fileDescriptorCGs_32Teammgr, []int{0} }

type C_SetTeam struct {
	Tp int32          `protobuf:"varint,1,opt,name=Tp,proto3" json:"Tp,omitempty"`
	T  *TeamFormation `protobuf:"bytes,2,opt,name=T" json:"T,omitempty"`
}

func (m *C_SetTeam) Reset()                    { *m = C_SetTeam{} }
func (m *C_SetTeam) String() string            { return proto.CompactTextString(m) }
func (*C_SetTeam) ProtoMessage()               {}
func (*C_SetTeam) Descriptor() ([]byte, []int) { return fileDescriptorCGs_32Teammgr, []int{1} }

type GS_SetTeam_R struct {
	ErrorCode int32          `protobuf:"varint,1,opt,name=ErrorCode,proto3" json:"ErrorCode,omitempty"`
	Tp        int32          `protobuf:"varint,2,opt,name=Tp,proto3" json:"Tp,omitempty"`
	T         *TeamFormation `protobuf:"bytes,3,opt,name=T" json:"T,omitempty"`
}

func (m *GS_SetTeam_R) Reset()                    { *m = GS_SetTeam_R{} }
func (m *GS_SetTeam_R) String() string            { return proto.CompactTextString(m) }
func (*GS_SetTeam_R) ProtoMessage()               {}
func (*GS_SetTeam_R) Descriptor() ([]byte, []int) { return fileDescriptorCGs_32Teammgr, []int{2} }

func init() {
	proto.RegisterType((*TeamMgrData)(nil), "msg.TeamMgrData")
	proto.RegisterType((*C_SetTeam)(nil), "msg.C_SetTeam")
	proto.RegisterType((*GS_SetTeam_R)(nil), "msg.GS_SetTeam_R")
}
func (m *TeamMgrData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TeamMgrData) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Teams) > 0 {
		for k, _ := range m.Teams {
			dAtA[i] = 0xa
			i++
			v := m.Teams[k]
			msgSize := 0
			if v != nil {
				msgSize = v.Size()
				msgSize += 1 + sovCGs_32Teammgr(uint64(msgSize))
			}
			mapSize := 1 + sovCGs_32Teammgr(uint64(k)) + msgSize
			i = encodeVarintCGs_32Teammgr(dAtA, i, uint64(mapSize))
			dAtA[i] = 0x8
			i++
			i = encodeVarintCGs_32Teammgr(dAtA, i, uint64(k))
			if v != nil {
				dAtA[i] = 0x12
				i++
				i = encodeVarintCGs_32Teammgr(dAtA, i, uint64(v.Size()))
				n1, err := v.MarshalTo(dAtA[i:])
				if err != nil {
					return 0, err
				}
				i += n1
			}
		}
	}
	return i, nil
}

func (m *C_SetTeam) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *C_SetTeam) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Tp != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCGs_32Teammgr(dAtA, i, uint64(m.Tp))
	}
	if m.T != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintCGs_32Teammgr(dAtA, i, uint64(m.T.Size()))
		n2, err := m.T.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func (m *GS_SetTeam_R) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GS_SetTeam_R) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ErrorCode != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCGs_32Teammgr(dAtA, i, uint64(m.ErrorCode))
	}
	if m.Tp != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintCGs_32Teammgr(dAtA, i, uint64(m.Tp))
	}
	if m.T != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintCGs_32Teammgr(dAtA, i, uint64(m.T.Size()))
		n3, err := m.T.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}

func encodeVarintCGs_32Teammgr(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *TeamMgrData) Size() (n int) {
	var l int
	_ = l
	if len(m.Teams) > 0 {
		for k, v := range m.Teams {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovCGs_32Teammgr(uint64(l))
			}
			mapEntrySize := 1 + sovCGs_32Teammgr(uint64(k)) + l
			n += mapEntrySize + 1 + sovCGs_32Teammgr(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *C_SetTeam) Size() (n int) {
	var l int
	_ = l
	if m.Tp != 0 {
		n += 1 + sovCGs_32Teammgr(uint64(m.Tp))
	}
	if m.T != nil {
		l = m.T.Size()
		n += 1 + l + sovCGs_32Teammgr(uint64(l))
	}
	return n
}

func (m *GS_SetTeam_R) Size() (n int) {
	var l int
	_ = l
	if m.ErrorCode != 0 {
		n += 1 + sovCGs_32Teammgr(uint64(m.ErrorCode))
	}
	if m.Tp != 0 {
		n += 1 + sovCGs_32Teammgr(uint64(m.Tp))
	}
	if m.T != nil {
		l = m.T.Size()
		n += 1 + l + sovCGs_32Teammgr(uint64(l))
	}
	return n
}

func sovCGs_32Teammgr(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCGs_32Teammgr(x uint64) (n int) {
	return sovCGs_32Teammgr(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TeamMgrData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_32Teammgr
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
			return fmt.Errorf("proto: TeamMgrData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TeamMgrData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Teams", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_32Teammgr
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
				return ErrInvalidLengthCGs_32Teammgr
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Teams == nil {
				m.Teams = make(map[int32]*TeamFormation)
			}
			var mapkey int32
			var mapvalue *TeamFormation
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_32Teammgr
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
							return ErrIntOverflowCGs_32Teammgr
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
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_32Teammgr
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= (int(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthCGs_32Teammgr
					}
					postmsgIndex := iNdEx + mapmsglen
					if mapmsglen < 0 {
						return ErrInvalidLengthCGs_32Teammgr
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &TeamFormation{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipCGs_32Teammgr(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthCGs_32Teammgr
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Teams[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_32Teammgr(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_32Teammgr
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
func (m *C_SetTeam) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_32Teammgr
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
			return fmt.Errorf("proto: C_SetTeam: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: C_SetTeam: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tp", wireType)
			}
			m.Tp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_32Teammgr
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Tp |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field T", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_32Teammgr
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
				return ErrInvalidLengthCGs_32Teammgr
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.T == nil {
				m.T = &TeamFormation{}
			}
			if err := m.T.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_32Teammgr(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_32Teammgr
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
func (m *GS_SetTeam_R) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_32Teammgr
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
			return fmt.Errorf("proto: GS_SetTeam_R: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GS_SetTeam_R: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrorCode", wireType)
			}
			m.ErrorCode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_32Teammgr
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
				return fmt.Errorf("proto: wrong wireType = %d for field Tp", wireType)
			}
			m.Tp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_32Teammgr
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Tp |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field T", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_32Teammgr
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
				return ErrInvalidLengthCGs_32Teammgr
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.T == nil {
				m.T = &TeamFormation{}
			}
			if err := m.T.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_32Teammgr(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_32Teammgr
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
func skipCGs_32Teammgr(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCGs_32Teammgr
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
					return 0, ErrIntOverflowCGs_32Teammgr
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
					return 0, ErrIntOverflowCGs_32Teammgr
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
				return 0, ErrInvalidLengthCGs_32Teammgr
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCGs_32Teammgr
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
				next, err := skipCGs_32Teammgr(dAtA[start:])
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
	ErrInvalidLengthCGs_32Teammgr = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCGs_32Teammgr   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("c_gs.32.teammgr.proto", fileDescriptorCGs_32Teammgr) }

var fileDescriptorCGs_32Teammgr = []byte{
	// 273 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xc1, 0x4a, 0xc3, 0x30,
	0x18, 0xc7, 0xf7, 0xb5, 0x4c, 0xd8, 0x57, 0x11, 0x09, 0x0a, 0xa5, 0x4a, 0x29, 0x3d, 0xf5, 0x14,
	0x66, 0x77, 0x11, 0xc1, 0x8b, 0x73, 0x7a, 0xd1, 0x4b, 0xd7, 0xb3, 0x25, 0xce, 0x10, 0x44, 0xb3,
	0x94, 0x24, 0x13, 0xf6, 0x12, 0xe2, 0x63, 0xed, 0xe8, 0x23, 0x68, 0x7d, 0x11, 0xe9, 0x52, 0x9d,
	0x5e, 0x76, 0xfb, 0xe7, 0xc7, 0x2f, 0xbf, 0xc3, 0x87, 0x87, 0xb3, 0x4a, 0x18, 0x3a, 0xca, 0xa9,
	0xe5, 0x4c, 0x4a, 0xa1, 0x69, 0xad, 0x95, 0x55, 0xc4, 0x97, 0x46, 0x44, 0x28, 0x94, 0x50, 0x0e,
	0x44, 0xce, 0x1b, 0x0e, 0xa9, 0xb1, 0x7a, 0x31, 0xb3, 0xc6, 0xe1, 0xf4, 0x15, 0x30, 0x28, 0x39,
	0x93, 0xb7, 0x42, 0x5f, 0x32, 0xcb, 0xc8, 0x09, 0xf6, 0xdb, 0xa7, 0x09, 0x21, 0xf1, 0xb3, 0x20,
	0x3f, 0xa2, 0xd2, 0x08, 0xfa, 0x47, 0x58, 0x6f, 0x33, 0x99, 0x5b, 0xbd, 0x2c, 0x9c, 0x19, 0xdd,
	0x20, 0x6e, 0x20, 0xd9, 0x47, 0xff, 0x89, 0x2f, 0x43, 0x48, 0x20, 0xeb, 0x17, 0xed, 0x24, 0x19,
	0xf6, 0x5f, 0xd8, 0xf3, 0x82, 0x87, 0x5e, 0x02, 0x59, 0x90, 0x93, 0xdf, 0xe4, 0x95, 0xd2, 0x92,
	0xd9, 0x47, 0x35, 0x2f, 0x9c, 0x70, 0xe6, 0x9d, 0x42, 0x7a, 0x8e, 0x83, 0x71, 0x35, 0xe5, 0xb6,
	0x15, 0xc8, 0x1e, 0x7a, 0x65, 0xdd, 0xb5, 0xbc, 0xb2, 0x26, 0x09, 0x42, 0xb9, 0x25, 0x03, 0x65,
	0x7a, 0x87, 0xbb, 0xd7, 0xd3, 0x9f, 0xff, 0x55, 0x41, 0x8e, 0x71, 0x30, 0xd1, 0x5a, 0xe9, 0xb1,
	0x7a, 0xe0, 0x5d, 0x68, 0x03, 0xba, 0xbe, 0xf7, 0xbf, 0xef, 0x6f, 0xe9, 0x5f, 0x1c, 0xac, 0x3e,
	0xe3, 0xde, 0xaa, 0x89, 0xe1, 0xbd, 0x89, 0xe1, 0xa3, 0x89, 0xe1, 0xed, 0x2b, 0xee, 0xdd, 0xef,
	0xac, 0x8f, 0x39, 0xfa, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x57, 0x46, 0xf1, 0x8d, 0x01, 0x00,
	0x00,
}