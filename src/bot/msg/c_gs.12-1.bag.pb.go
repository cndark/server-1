// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: c_gs.12-1.bag.proto

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

// 背包数据
type BagData struct {
	Currency []*Ccy   `protobuf:"bytes,1,rep,name=Currency" json:"Currency,omitempty"`
	Items    []*Item  `protobuf:"bytes,2,rep,name=Items" json:"Items,omitempty"`
	Heroes   []*Hero  `protobuf:"bytes,3,rep,name=Heroes" json:"Heroes,omitempty"`
	Armors   []*Armor `protobuf:"bytes,4,rep,name=Armors" json:"Armors,omitempty"`
	Relics   []*Relic `protobuf:"bytes,5,rep,name=Relics" json:"Relics,omitempty"`
}

func (m *BagData) Reset()                    { *m = BagData{} }
func (m *BagData) String() string            { return proto.CompactTextString(m) }
func (*BagData) ProtoMessage()               {}
func (*BagData) Descriptor() ([]byte, []int) { return fileDescriptorCGs_12_1Bag, []int{0} }

// 背包变化推送
type GS_BagUpdate struct {
	Currency  []*Ccy   `protobuf:"bytes,1,rep,name=Currency" json:"Currency,omitempty"`
	Items     []*Item  `protobuf:"bytes,2,rep,name=Items" json:"Items,omitempty"`
	Heroes    []*Hero  `protobuf:"bytes,3,rep,name=Heroes" json:"Heroes,omitempty"`
	HeroesDel []int64  `protobuf:"varint,4,rep,packed,name=HeroesDel" json:"HeroesDel,omitempty"`
	Relics    []*Relic `protobuf:"bytes,5,rep,name=Relics" json:"Relics,omitempty"`
	RelicsDel []int64  `protobuf:"varint,6,rep,packed,name=RelicsDel" json:"RelicsDel,omitempty"`
}

func (m *GS_BagUpdate) Reset()                    { *m = GS_BagUpdate{} }
func (m *GS_BagUpdate) String() string            { return proto.CompactTextString(m) }
func (*GS_BagUpdate) ProtoMessage()               {}
func (*GS_BagUpdate) Descriptor() ([]byte, []int) { return fileDescriptorCGs_12_1Bag, []int{1} }

func init() {
	proto.RegisterType((*BagData)(nil), "msg.BagData")
	proto.RegisterType((*GS_BagUpdate)(nil), "msg.GS_BagUpdate")
}
func (m *BagData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BagData) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Currency) > 0 {
		for _, msg := range m.Currency {
			dAtA[i] = 0xa
			i++
			i = encodeVarintCGs_12_1Bag(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Items) > 0 {
		for _, msg := range m.Items {
			dAtA[i] = 0x12
			i++
			i = encodeVarintCGs_12_1Bag(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Heroes) > 0 {
		for _, msg := range m.Heroes {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintCGs_12_1Bag(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Armors) > 0 {
		for _, msg := range m.Armors {
			dAtA[i] = 0x22
			i++
			i = encodeVarintCGs_12_1Bag(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Relics) > 0 {
		for _, msg := range m.Relics {
			dAtA[i] = 0x2a
			i++
			i = encodeVarintCGs_12_1Bag(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *GS_BagUpdate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GS_BagUpdate) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Currency) > 0 {
		for _, msg := range m.Currency {
			dAtA[i] = 0xa
			i++
			i = encodeVarintCGs_12_1Bag(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Items) > 0 {
		for _, msg := range m.Items {
			dAtA[i] = 0x12
			i++
			i = encodeVarintCGs_12_1Bag(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Heroes) > 0 {
		for _, msg := range m.Heroes {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintCGs_12_1Bag(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.HeroesDel) > 0 {
		dAtA2 := make([]byte, len(m.HeroesDel)*10)
		var j1 int
		for _, num1 := range m.HeroesDel {
			num := uint64(num1)
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		dAtA[i] = 0x22
		i++
		i = encodeVarintCGs_12_1Bag(dAtA, i, uint64(j1))
		i += copy(dAtA[i:], dAtA2[:j1])
	}
	if len(m.Relics) > 0 {
		for _, msg := range m.Relics {
			dAtA[i] = 0x2a
			i++
			i = encodeVarintCGs_12_1Bag(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.RelicsDel) > 0 {
		dAtA4 := make([]byte, len(m.RelicsDel)*10)
		var j3 int
		for _, num1 := range m.RelicsDel {
			num := uint64(num1)
			for num >= 1<<7 {
				dAtA4[j3] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j3++
			}
			dAtA4[j3] = uint8(num)
			j3++
		}
		dAtA[i] = 0x32
		i++
		i = encodeVarintCGs_12_1Bag(dAtA, i, uint64(j3))
		i += copy(dAtA[i:], dAtA4[:j3])
	}
	return i, nil
}

func encodeVarintCGs_12_1Bag(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *BagData) Size() (n int) {
	var l int
	_ = l
	if len(m.Currency) > 0 {
		for _, e := range m.Currency {
			l = e.Size()
			n += 1 + l + sovCGs_12_1Bag(uint64(l))
		}
	}
	if len(m.Items) > 0 {
		for _, e := range m.Items {
			l = e.Size()
			n += 1 + l + sovCGs_12_1Bag(uint64(l))
		}
	}
	if len(m.Heroes) > 0 {
		for _, e := range m.Heroes {
			l = e.Size()
			n += 1 + l + sovCGs_12_1Bag(uint64(l))
		}
	}
	if len(m.Armors) > 0 {
		for _, e := range m.Armors {
			l = e.Size()
			n += 1 + l + sovCGs_12_1Bag(uint64(l))
		}
	}
	if len(m.Relics) > 0 {
		for _, e := range m.Relics {
			l = e.Size()
			n += 1 + l + sovCGs_12_1Bag(uint64(l))
		}
	}
	return n
}

func (m *GS_BagUpdate) Size() (n int) {
	var l int
	_ = l
	if len(m.Currency) > 0 {
		for _, e := range m.Currency {
			l = e.Size()
			n += 1 + l + sovCGs_12_1Bag(uint64(l))
		}
	}
	if len(m.Items) > 0 {
		for _, e := range m.Items {
			l = e.Size()
			n += 1 + l + sovCGs_12_1Bag(uint64(l))
		}
	}
	if len(m.Heroes) > 0 {
		for _, e := range m.Heroes {
			l = e.Size()
			n += 1 + l + sovCGs_12_1Bag(uint64(l))
		}
	}
	if len(m.HeroesDel) > 0 {
		l = 0
		for _, e := range m.HeroesDel {
			l += sovCGs_12_1Bag(uint64(e))
		}
		n += 1 + sovCGs_12_1Bag(uint64(l)) + l
	}
	if len(m.Relics) > 0 {
		for _, e := range m.Relics {
			l = e.Size()
			n += 1 + l + sovCGs_12_1Bag(uint64(l))
		}
	}
	if len(m.RelicsDel) > 0 {
		l = 0
		for _, e := range m.RelicsDel {
			l += sovCGs_12_1Bag(uint64(e))
		}
		n += 1 + sovCGs_12_1Bag(uint64(l)) + l
	}
	return n
}

func sovCGs_12_1Bag(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCGs_12_1Bag(x uint64) (n int) {
	return sovCGs_12_1Bag(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BagData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_12_1Bag
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
			return fmt.Errorf("proto: BagData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BagData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Currency", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_12_1Bag
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
				return ErrInvalidLengthCGs_12_1Bag
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Currency = append(m.Currency, &Ccy{})
			if err := m.Currency[len(m.Currency)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Items", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_12_1Bag
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
				return ErrInvalidLengthCGs_12_1Bag
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Items = append(m.Items, &Item{})
			if err := m.Items[len(m.Items)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Heroes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_12_1Bag
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
				return ErrInvalidLengthCGs_12_1Bag
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Heroes = append(m.Heroes, &Hero{})
			if err := m.Heroes[len(m.Heroes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Armors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_12_1Bag
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
				return ErrInvalidLengthCGs_12_1Bag
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Armors = append(m.Armors, &Armor{})
			if err := m.Armors[len(m.Armors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Relics", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_12_1Bag
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
				return ErrInvalidLengthCGs_12_1Bag
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Relics = append(m.Relics, &Relic{})
			if err := m.Relics[len(m.Relics)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_12_1Bag(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_12_1Bag
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
func (m *GS_BagUpdate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_12_1Bag
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
			return fmt.Errorf("proto: GS_BagUpdate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GS_BagUpdate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Currency", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_12_1Bag
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
				return ErrInvalidLengthCGs_12_1Bag
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Currency = append(m.Currency, &Ccy{})
			if err := m.Currency[len(m.Currency)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Items", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_12_1Bag
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
				return ErrInvalidLengthCGs_12_1Bag
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Items = append(m.Items, &Item{})
			if err := m.Items[len(m.Items)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Heroes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_12_1Bag
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
				return ErrInvalidLengthCGs_12_1Bag
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Heroes = append(m.Heroes, &Hero{})
			if err := m.Heroes[len(m.Heroes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType == 0 {
				var v int64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_12_1Bag
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (int64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.HeroesDel = append(m.HeroesDel, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_12_1Bag
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthCGs_12_1Bag
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v int64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_12_1Bag
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (int64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.HeroesDel = append(m.HeroesDel, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field HeroesDel", wireType)
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Relics", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_12_1Bag
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
				return ErrInvalidLengthCGs_12_1Bag
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Relics = append(m.Relics, &Relic{})
			if err := m.Relics[len(m.Relics)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType == 0 {
				var v int64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_12_1Bag
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (int64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.RelicsDel = append(m.RelicsDel, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_12_1Bag
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthCGs_12_1Bag
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v int64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_12_1Bag
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (int64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.RelicsDel = append(m.RelicsDel, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field RelicsDel", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_12_1Bag(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_12_1Bag
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
func skipCGs_12_1Bag(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCGs_12_1Bag
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
					return 0, ErrIntOverflowCGs_12_1Bag
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
					return 0, ErrIntOverflowCGs_12_1Bag
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
				return 0, ErrInvalidLengthCGs_12_1Bag
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCGs_12_1Bag
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
				next, err := skipCGs_12_1Bag(dAtA[start:])
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
	ErrInvalidLengthCGs_12_1Bag = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCGs_12_1Bag   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("c_gs.12-1.bag.proto", fileDescriptorCGs_12_1Bag) }

var fileDescriptorCGs_12_1Bag = []byte{
	// 270 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x91, 0xb1, 0x4a, 0x03, 0x31,
	0x1c, 0xc6, 0x1b, 0xcf, 0x9e, 0xed, 0x5f, 0xa7, 0xa8, 0x10, 0x8a, 0x44, 0x3d, 0x1c, 0x5c, 0x0c,
	0x6d, 0x7d, 0x02, 0xaf, 0x05, 0x75, 0x8d, 0x38, 0x97, 0x34, 0x86, 0x20, 0xf4, 0x4c, 0x49, 0xd2,
	0xa1, 0x6f, 0xe1, 0xfb, 0xf8, 0x02, 0x1d, 0x1d, 0x1d, 0xf5, 0x7c, 0x11, 0x49, 0x72, 0xd4, 0x1b,
	0x9d, 0xdc, 0xfe, 0xfc, 0xbe, 0x1f, 0x1f, 0xf9, 0x08, 0x1c, 0xca, 0x99, 0x76, 0x6c, 0x34, 0xbe,
	0x1a, 0xb1, 0xb9, 0xd0, 0x6c, 0x69, 0x8d, 0x37, 0x38, 0xab, 0x9c, 0x1e, 0x80, 0x36, 0xda, 0x24,
	0x30, 0x38, 0x8e, 0xd6, 0x70, 0xc8, 0x9c, 0xb7, 0x2b, 0xe9, 0x5d, 0xc2, 0xc5, 0x1b, 0x82, 0xbd,
	0x52, 0xe8, 0xa9, 0xf0, 0x02, 0x5f, 0x40, 0x6f, 0xb2, 0xb2, 0x56, 0xbd, 0xc8, 0x35, 0x41, 0x67,
	0xd9, 0xe5, 0xfe, 0xb8, 0xc7, 0x2a, 0xa7, 0xd9, 0x44, 0xae, 0xf9, 0x36, 0xc1, 0xa7, 0xd0, 0xbd,
	0xf7, 0xaa, 0x72, 0x64, 0x27, 0x2a, 0xfd, 0xa8, 0x04, 0xc2, 0x13, 0xc7, 0xe7, 0x90, 0xdf, 0x29,
	0x6b, 0x94, 0x23, 0x59, 0xcb, 0x08, 0x88, 0x37, 0x01, 0x2e, 0x20, 0xbf, 0xb1, 0x95, 0xb1, 0x8e,
	0xec, 0x46, 0x05, 0xa2, 0x12, 0x11, 0x6f, 0x92, 0xe0, 0x70, 0xb5, 0x78, 0x96, 0x8e, 0x74, 0x5b,
	0x4e, 0x44, 0xbc, 0x49, 0x8a, 0x0f, 0x04, 0x07, 0xb7, 0x0f, 0xb3, 0x52, 0xe8, 0xc7, 0xe5, 0x93,
	0xf0, 0xea, 0x1f, 0x27, 0x9c, 0x40, 0x3f, 0x5d, 0x53, 0xb5, 0x88, 0x2b, 0x32, 0xfe, 0x0b, 0xfe,
	0xf2, 0xf8, 0xd0, 0x90, 0xae, 0xd0, 0x90, 0xa7, 0x86, 0x2d, 0x28, 0x8f, 0x36, 0x5f, 0xb4, 0xb3,
	0xa9, 0x29, 0x7a, 0xaf, 0x29, 0xfa, 0xac, 0x29, 0x7a, 0xfd, 0xa6, 0x9d, 0x79, 0x1e, 0x7f, 0xed,
	0xfa, 0x27, 0x00, 0x00, 0xff, 0xff, 0x5f, 0x9b, 0xe8, 0xbe, 0xf4, 0x01, 0x00, 0x00,
}
