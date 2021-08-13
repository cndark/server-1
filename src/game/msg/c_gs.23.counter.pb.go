// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: c_gs.23.counter.proto

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

type CounterData struct {
	Cnt    map[int32]int64 `protobuf:"bytes,1,rep,name=Cnt" json:"Cnt,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	MaxCnt map[int32]int64 `protobuf:"bytes,2,rep,name=MaxCnt" json:"MaxCnt,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Ts     map[int32]int64 `protobuf:"bytes,4,rep,name=Ts" json:"Ts,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (m *CounterData) Reset()                    { *m = CounterData{} }
func (m *CounterData) String() string            { return proto.CompactTextString(m) }
func (*CounterData) ProtoMessage()               {}
func (*CounterData) Descriptor() ([]byte, []int) { return fileDescriptorCGs_23Counter, []int{0} }

// 推送计数器操作
type GS_CounterOpUpdate struct {
	Cnt    map[int32]int64 `protobuf:"bytes,1,rep,name=Cnt" json:"Cnt,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	MaxCnt map[int32]int64 `protobuf:"bytes,2,rep,name=MaxCnt" json:"MaxCnt,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Ts     map[int32]int64 `protobuf:"bytes,3,rep,name=Ts" json:"Ts,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (m *GS_CounterOpUpdate) Reset()                    { *m = GS_CounterOpUpdate{} }
func (m *GS_CounterOpUpdate) String() string            { return proto.CompactTextString(m) }
func (*GS_CounterOpUpdate) ProtoMessage()               {}
func (*GS_CounterOpUpdate) Descriptor() ([]byte, []int) { return fileDescriptorCGs_23Counter, []int{1} }

// 要求计算计数器恢复
type C_CounterRecover struct {
	Id int32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (m *C_CounterRecover) Reset()                    { *m = C_CounterRecover{} }
func (m *C_CounterRecover) String() string            { return proto.CompactTextString(m) }
func (*C_CounterRecover) ProtoMessage()               {}
func (*C_CounterRecover) Descriptor() ([]byte, []int) { return fileDescriptorCGs_23Counter, []int{2} }

type GS_CounterRecover_R struct {
	ErrorCode int32 `protobuf:"varint,1,opt,name=ErrorCode,proto3" json:"ErrorCode,omitempty"`
	Id        int32 `protobuf:"varint,2,opt,name=Id,proto3" json:"Id,omitempty"`
	Cnt       int64 `protobuf:"varint,3,opt,name=Cnt,proto3" json:"Cnt,omitempty"`
	Ts        int64 `protobuf:"varint,4,opt,name=Ts,proto3" json:"Ts,omitempty"`
}

func (m *GS_CounterRecover_R) Reset()                    { *m = GS_CounterRecover_R{} }
func (m *GS_CounterRecover_R) String() string            { return proto.CompactTextString(m) }
func (*GS_CounterRecover_R) ProtoMessage()               {}
func (*GS_CounterRecover_R) Descriptor() ([]byte, []int) { return fileDescriptorCGs_23Counter, []int{3} }

// 购买计数
type C_CounterBuy struct {
	Id int32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (m *C_CounterBuy) Reset()                    { *m = C_CounterBuy{} }
func (m *C_CounterBuy) String() string            { return proto.CompactTextString(m) }
func (*C_CounterBuy) ProtoMessage()               {}
func (*C_CounterBuy) Descriptor() ([]byte, []int) { return fileDescriptorCGs_23Counter, []int{4} }

type GS_CounterBuy_R struct {
	ErrorCode int32 `protobuf:"varint,1,opt,name=ErrorCode,proto3" json:"ErrorCode,omitempty"`
}

func (m *GS_CounterBuy_R) Reset()                    { *m = GS_CounterBuy_R{} }
func (m *GS_CounterBuy_R) String() string            { return proto.CompactTextString(m) }
func (*GS_CounterBuy_R) ProtoMessage()               {}
func (*GS_CounterBuy_R) Descriptor() ([]byte, []int) { return fileDescriptorCGs_23Counter, []int{5} }

func init() {
	proto.RegisterType((*CounterData)(nil), "msg.CounterData")
	proto.RegisterType((*GS_CounterOpUpdate)(nil), "msg.GS_CounterOpUpdate")
	proto.RegisterType((*C_CounterRecover)(nil), "msg.C_CounterRecover")
	proto.RegisterType((*GS_CounterRecover_R)(nil), "msg.GS_CounterRecover_R")
	proto.RegisterType((*C_CounterBuy)(nil), "msg.C_CounterBuy")
	proto.RegisterType((*GS_CounterBuy_R)(nil), "msg.GS_CounterBuy_R")
}
func (m *CounterData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CounterData) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Cnt) > 0 {
		for k, _ := range m.Cnt {
			dAtA[i] = 0xa
			i++
			v := m.Cnt[k]
			mapSize := 1 + sovCGs_23Counter(uint64(k)) + 1 + sovCGs_23Counter(uint64(v))
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(mapSize))
			dAtA[i] = 0x8
			i++
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(k))
			dAtA[i] = 0x10
			i++
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(v))
		}
	}
	if len(m.MaxCnt) > 0 {
		for k, _ := range m.MaxCnt {
			dAtA[i] = 0x12
			i++
			v := m.MaxCnt[k]
			mapSize := 1 + sovCGs_23Counter(uint64(k)) + 1 + sovCGs_23Counter(uint64(v))
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(mapSize))
			dAtA[i] = 0x8
			i++
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(k))
			dAtA[i] = 0x10
			i++
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(v))
		}
	}
	if len(m.Ts) > 0 {
		for k, _ := range m.Ts {
			dAtA[i] = 0x22
			i++
			v := m.Ts[k]
			mapSize := 1 + sovCGs_23Counter(uint64(k)) + 1 + sovCGs_23Counter(uint64(v))
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(mapSize))
			dAtA[i] = 0x8
			i++
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(k))
			dAtA[i] = 0x10
			i++
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(v))
		}
	}
	return i, nil
}

func (m *GS_CounterOpUpdate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GS_CounterOpUpdate) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Cnt) > 0 {
		for k, _ := range m.Cnt {
			dAtA[i] = 0xa
			i++
			v := m.Cnt[k]
			mapSize := 1 + sovCGs_23Counter(uint64(k)) + 1 + sovCGs_23Counter(uint64(v))
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(mapSize))
			dAtA[i] = 0x8
			i++
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(k))
			dAtA[i] = 0x10
			i++
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(v))
		}
	}
	if len(m.MaxCnt) > 0 {
		for k, _ := range m.MaxCnt {
			dAtA[i] = 0x12
			i++
			v := m.MaxCnt[k]
			mapSize := 1 + sovCGs_23Counter(uint64(k)) + 1 + sovCGs_23Counter(uint64(v))
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(mapSize))
			dAtA[i] = 0x8
			i++
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(k))
			dAtA[i] = 0x10
			i++
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(v))
		}
	}
	if len(m.Ts) > 0 {
		for k, _ := range m.Ts {
			dAtA[i] = 0x1a
			i++
			v := m.Ts[k]
			mapSize := 1 + sovCGs_23Counter(uint64(k)) + 1 + sovCGs_23Counter(uint64(v))
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(mapSize))
			dAtA[i] = 0x8
			i++
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(k))
			dAtA[i] = 0x10
			i++
			i = encodeVarintCGs_23Counter(dAtA, i, uint64(v))
		}
	}
	return i, nil
}

func (m *C_CounterRecover) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *C_CounterRecover) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCGs_23Counter(dAtA, i, uint64(m.Id))
	}
	return i, nil
}

func (m *GS_CounterRecover_R) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GS_CounterRecover_R) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ErrorCode != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCGs_23Counter(dAtA, i, uint64(m.ErrorCode))
	}
	if m.Id != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintCGs_23Counter(dAtA, i, uint64(m.Id))
	}
	if m.Cnt != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintCGs_23Counter(dAtA, i, uint64(m.Cnt))
	}
	if m.Ts != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintCGs_23Counter(dAtA, i, uint64(m.Ts))
	}
	return i, nil
}

func (m *C_CounterBuy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *C_CounterBuy) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCGs_23Counter(dAtA, i, uint64(m.Id))
	}
	return i, nil
}

func (m *GS_CounterBuy_R) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GS_CounterBuy_R) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ErrorCode != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCGs_23Counter(dAtA, i, uint64(m.ErrorCode))
	}
	return i, nil
}

func encodeVarintCGs_23Counter(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *CounterData) Size() (n int) {
	var l int
	_ = l
	if len(m.Cnt) > 0 {
		for k, v := range m.Cnt {
			_ = k
			_ = v
			mapEntrySize := 1 + sovCGs_23Counter(uint64(k)) + 1 + sovCGs_23Counter(uint64(v))
			n += mapEntrySize + 1 + sovCGs_23Counter(uint64(mapEntrySize))
		}
	}
	if len(m.MaxCnt) > 0 {
		for k, v := range m.MaxCnt {
			_ = k
			_ = v
			mapEntrySize := 1 + sovCGs_23Counter(uint64(k)) + 1 + sovCGs_23Counter(uint64(v))
			n += mapEntrySize + 1 + sovCGs_23Counter(uint64(mapEntrySize))
		}
	}
	if len(m.Ts) > 0 {
		for k, v := range m.Ts {
			_ = k
			_ = v
			mapEntrySize := 1 + sovCGs_23Counter(uint64(k)) + 1 + sovCGs_23Counter(uint64(v))
			n += mapEntrySize + 1 + sovCGs_23Counter(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *GS_CounterOpUpdate) Size() (n int) {
	var l int
	_ = l
	if len(m.Cnt) > 0 {
		for k, v := range m.Cnt {
			_ = k
			_ = v
			mapEntrySize := 1 + sovCGs_23Counter(uint64(k)) + 1 + sovCGs_23Counter(uint64(v))
			n += mapEntrySize + 1 + sovCGs_23Counter(uint64(mapEntrySize))
		}
	}
	if len(m.MaxCnt) > 0 {
		for k, v := range m.MaxCnt {
			_ = k
			_ = v
			mapEntrySize := 1 + sovCGs_23Counter(uint64(k)) + 1 + sovCGs_23Counter(uint64(v))
			n += mapEntrySize + 1 + sovCGs_23Counter(uint64(mapEntrySize))
		}
	}
	if len(m.Ts) > 0 {
		for k, v := range m.Ts {
			_ = k
			_ = v
			mapEntrySize := 1 + sovCGs_23Counter(uint64(k)) + 1 + sovCGs_23Counter(uint64(v))
			n += mapEntrySize + 1 + sovCGs_23Counter(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *C_CounterRecover) Size() (n int) {
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovCGs_23Counter(uint64(m.Id))
	}
	return n
}

func (m *GS_CounterRecover_R) Size() (n int) {
	var l int
	_ = l
	if m.ErrorCode != 0 {
		n += 1 + sovCGs_23Counter(uint64(m.ErrorCode))
	}
	if m.Id != 0 {
		n += 1 + sovCGs_23Counter(uint64(m.Id))
	}
	if m.Cnt != 0 {
		n += 1 + sovCGs_23Counter(uint64(m.Cnt))
	}
	if m.Ts != 0 {
		n += 1 + sovCGs_23Counter(uint64(m.Ts))
	}
	return n
}

func (m *C_CounterBuy) Size() (n int) {
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovCGs_23Counter(uint64(m.Id))
	}
	return n
}

func (m *GS_CounterBuy_R) Size() (n int) {
	var l int
	_ = l
	if m.ErrorCode != 0 {
		n += 1 + sovCGs_23Counter(uint64(m.ErrorCode))
	}
	return n
}

func sovCGs_23Counter(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCGs_23Counter(x uint64) (n int) {
	return sovCGs_23Counter(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *CounterData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_23Counter
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
			return fmt.Errorf("proto: CounterData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CounterData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cnt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_23Counter
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
				return ErrInvalidLengthCGs_23Counter
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Cnt == nil {
				m.Cnt = make(map[int32]int64)
			}
			var mapkey int32
			var mapvalue int64
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_23Counter
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
							return ErrIntOverflowCGs_23Counter
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
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_23Counter
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= (int64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipCGs_23Counter(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthCGs_23Counter
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Cnt[mapkey] = mapvalue
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxCnt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_23Counter
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
				return ErrInvalidLengthCGs_23Counter
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MaxCnt == nil {
				m.MaxCnt = make(map[int32]int64)
			}
			var mapkey int32
			var mapvalue int64
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_23Counter
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
							return ErrIntOverflowCGs_23Counter
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
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_23Counter
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= (int64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipCGs_23Counter(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthCGs_23Counter
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.MaxCnt[mapkey] = mapvalue
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_23Counter
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
				return ErrInvalidLengthCGs_23Counter
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Ts == nil {
				m.Ts = make(map[int32]int64)
			}
			var mapkey int32
			var mapvalue int64
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_23Counter
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
							return ErrIntOverflowCGs_23Counter
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
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_23Counter
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= (int64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipCGs_23Counter(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthCGs_23Counter
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Ts[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_23Counter(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_23Counter
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
func (m *GS_CounterOpUpdate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_23Counter
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
			return fmt.Errorf("proto: GS_CounterOpUpdate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GS_CounterOpUpdate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cnt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_23Counter
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
				return ErrInvalidLengthCGs_23Counter
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Cnt == nil {
				m.Cnt = make(map[int32]int64)
			}
			var mapkey int32
			var mapvalue int64
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_23Counter
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
							return ErrIntOverflowCGs_23Counter
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
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_23Counter
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= (int64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipCGs_23Counter(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthCGs_23Counter
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Cnt[mapkey] = mapvalue
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxCnt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_23Counter
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
				return ErrInvalidLengthCGs_23Counter
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MaxCnt == nil {
				m.MaxCnt = make(map[int32]int64)
			}
			var mapkey int32
			var mapvalue int64
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_23Counter
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
							return ErrIntOverflowCGs_23Counter
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
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_23Counter
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= (int64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipCGs_23Counter(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthCGs_23Counter
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.MaxCnt[mapkey] = mapvalue
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_23Counter
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
				return ErrInvalidLengthCGs_23Counter
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Ts == nil {
				m.Ts = make(map[int32]int64)
			}
			var mapkey int32
			var mapvalue int64
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCGs_23Counter
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
							return ErrIntOverflowCGs_23Counter
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
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCGs_23Counter
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= (int64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipCGs_23Counter(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthCGs_23Counter
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Ts[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_23Counter(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_23Counter
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
func (m *C_CounterRecover) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_23Counter
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
			return fmt.Errorf("proto: C_CounterRecover: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: C_CounterRecover: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_23Counter
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
			skippy, err := skipCGs_23Counter(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_23Counter
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
func (m *GS_CounterRecover_R) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_23Counter
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
			return fmt.Errorf("proto: GS_CounterRecover_R: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GS_CounterRecover_R: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrorCode", wireType)
			}
			m.ErrorCode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_23Counter
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
					return ErrIntOverflowCGs_23Counter
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cnt", wireType)
			}
			m.Cnt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_23Counter
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Cnt |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ts", wireType)
			}
			m.Ts = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_23Counter
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
			skippy, err := skipCGs_23Counter(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_23Counter
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
func (m *C_CounterBuy) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_23Counter
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
			return fmt.Errorf("proto: C_CounterBuy: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: C_CounterBuy: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_23Counter
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
			skippy, err := skipCGs_23Counter(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_23Counter
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
func (m *GS_CounterBuy_R) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCGs_23Counter
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
			return fmt.Errorf("proto: GS_CounterBuy_R: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GS_CounterBuy_R: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrorCode", wireType)
			}
			m.ErrorCode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCGs_23Counter
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
		default:
			iNdEx = preIndex
			skippy, err := skipCGs_23Counter(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCGs_23Counter
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
func skipCGs_23Counter(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCGs_23Counter
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
					return 0, ErrIntOverflowCGs_23Counter
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
					return 0, ErrIntOverflowCGs_23Counter
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
				return 0, ErrInvalidLengthCGs_23Counter
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCGs_23Counter
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
				next, err := skipCGs_23Counter(dAtA[start:])
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
	ErrInvalidLengthCGs_23Counter = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCGs_23Counter   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("c_gs.23.counter.proto", fileDescriptorCGs_23Counter) }

var fileDescriptorCGs_23Counter = []byte{
	// 381 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x93, 0xcb, 0x4e, 0xf2, 0x40,
	0x18, 0x86, 0x99, 0xe9, 0x0f, 0xbf, 0x7e, 0x18, 0x24, 0x23, 0x26, 0x95, 0x90, 0xda, 0xd4, 0x4d,
	0x13, 0x93, 0x92, 0x80, 0x1a, 0x0f, 0x3b, 0x2a, 0x31, 0x2c, 0x8c, 0x49, 0xad, 0x6b, 0x52, 0xe9,
	0xa4, 0x0b, 0xa5, 0x43, 0x7a, 0x20, 0xf6, 0x2e, 0xbc, 0x28, 0x17, 0x2c, 0xbd, 0x04, 0xc4, 0x1b,
	0x31, 0x9d, 0x0e, 0x42, 0xa8, 0xc4, 0xb0, 0x74, 0xd7, 0xce, 0x3c, 0xcf, 0xf4, 0xfd, 0xde, 0xb6,
	0xb0, 0x3f, 0xe8, 0x7b, 0xa1, 0xd1, 0x6a, 0x1b, 0x03, 0x16, 0xfb, 0x11, 0x0d, 0x8c, 0x51, 0xc0,
	0x22, 0x46, 0xa4, 0x61, 0xe8, 0xd5, 0xc1, 0x63, 0x1e, 0xcb, 0x16, 0xb4, 0x37, 0x0c, 0x65, 0x33,
	0x43, 0xae, 0x9d, 0xc8, 0x21, 0xc7, 0x20, 0x99, 0x7e, 0x24, 0x23, 0x55, 0xd2, 0xcb, 0xad, 0x03,
	0x63, 0x18, 0x7a, 0xc6, 0xd2, 0xb6, 0x61, 0xfa, 0x51, 0xd7, 0x8f, 0x82, 0xc4, 0x4a, 0x29, 0x72,
	0x02, 0xa5, 0x5b, 0xe7, 0x25, 0xe5, 0x31, 0xe7, 0x1b, 0x39, 0x3e, 0xdb, 0xce, 0x14, 0xc1, 0x12,
	0x1d, 0xb0, 0x1d, 0xca, 0xff, 0xb8, 0x21, 0xe7, 0x0c, 0x3b, 0xcc, 0x68, 0x6c, 0x87, 0xf5, 0x33,
	0xd8, 0x9a, 0xdb, 0xa4, 0x0a, 0xd2, 0x13, 0x4d, 0x64, 0xa4, 0x22, 0xbd, 0x68, 0xa5, 0x97, 0xa4,
	0x06, 0xc5, 0xb1, 0xf3, 0x1c, 0x53, 0x19, 0xab, 0x48, 0x97, 0xac, 0xec, 0xe6, 0x12, 0x9f, 0xa3,
	0xfa, 0x05, 0x94, 0x97, 0x1e, 0xbc, 0x91, 0x7a, 0x0a, 0xff, 0x45, 0x82, 0x4d, 0x34, 0x6d, 0x8a,
	0x81, 0xdc, 0xdc, 0xf7, 0xc5, 0x20, 0x77, 0xa3, 0x87, 0x91, 0xeb, 0x44, 0x94, 0xb4, 0x96, 0xdb,
	0x54, 0xf9, 0xac, 0x79, 0x6a, 0xa5, 0xd4, 0xab, 0x95, 0x52, 0x8f, 0xd6, 0x69, 0x3f, 0x75, 0xdb,
	0xe4, 0xdd, 0x4a, 0x5c, 0x3c, 0x5c, 0x27, 0xfe, 0xcd, 0x8a, 0x35, 0xa8, 0x9a, 0xf3, 0x51, 0x2c,
	0x3a, 0x60, 0x63, 0x1a, 0x90, 0x0a, 0xe0, 0x9e, 0x2b, 0x74, 0xdc, 0x73, 0x35, 0x0a, 0x7b, 0x8b,
	0x79, 0x05, 0xd4, 0xb7, 0x48, 0x03, 0xb6, 0xbb, 0x41, 0xc0, 0x02, 0x93, 0xb9, 0x54, 0xd0, 0x8b,
	0x05, 0x71, 0x08, 0x9e, 0x1f, 0x92, 0x86, 0x4a, 0xdb, 0x97, 0x78, 0x00, 0xfe, 0x4a, 0x2a, 0xe2,
	0x8b, 0x4d, 0x17, 0xb0, 0x1d, 0x6a, 0x0a, 0xec, 0x7c, 0x47, 0xe9, 0xc4, 0x49, 0x2e, 0x46, 0x13,
	0x76, 0x17, 0x31, 0x3a, 0x71, 0xf2, 0x5b, 0x84, 0x4e, 0x6d, 0xf2, 0xa1, 0x14, 0x26, 0x33, 0x05,
	0xbd, 0xcf, 0x14, 0x34, 0x9d, 0x29, 0xe8, 0xf5, 0x53, 0x29, 0x3c, 0x96, 0xf8, 0x2f, 0xda, 0xfe,
	0x0a, 0x00, 0x00, 0xff, 0xff, 0x4f, 0x6f, 0x7a, 0xdd, 0xcc, 0x03, 0x00, 0x00,
}
