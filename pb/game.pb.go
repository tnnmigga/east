// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pb/tmp/game.proto

package pb

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type C2SPackage struct {
	UserID uint64 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
	Body   []byte `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (m *C2SPackage) Reset()         { *m = C2SPackage{} }
func (m *C2SPackage) String() string { return proto.CompactTextString(m) }
func (*C2SPackage) ProtoMessage()    {}
func (*C2SPackage) Descriptor() ([]byte, []int) {
	return fileDescriptor_70141d9832fb77cd, []int{0}
}
func (m *C2SPackage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *C2SPackage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_C2SPackage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *C2SPackage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_C2SPackage.Merge(m, src)
}
func (m *C2SPackage) XXX_Size() int {
	return m.Size()
}
func (m *C2SPackage) XXX_DiscardUnknown() {
	xxx_messageInfo_C2SPackage.DiscardUnknown(m)
}

var xxx_messageInfo_C2SPackage proto.InternalMessageInfo

func (m *C2SPackage) GetUserID() uint64 {
	if m != nil {
		return m.UserID
	}
	return 0
}

func (m *C2SPackage) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type S2CPackage struct {
	UserID uint64 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
	Body   []byte `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (m *S2CPackage) Reset()         { *m = S2CPackage{} }
func (m *S2CPackage) String() string { return proto.CompactTextString(m) }
func (*S2CPackage) ProtoMessage()    {}
func (*S2CPackage) Descriptor() ([]byte, []int) {
	return fileDescriptor_70141d9832fb77cd, []int{1}
}
func (m *S2CPackage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *S2CPackage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_S2CPackage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *S2CPackage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_S2CPackage.Merge(m, src)
}
func (m *S2CPackage) XXX_Size() int {
	return m.Size()
}
func (m *S2CPackage) XXX_DiscardUnknown() {
	xxx_messageInfo_S2CPackage.DiscardUnknown(m)
}

var xxx_messageInfo_S2CPackage proto.InternalMessageInfo

func (m *S2CPackage) GetUserID() uint64 {
	if m != nil {
		return m.UserID
	}
	return 0
}

func (m *S2CPackage) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type UserLoginReq struct {
	UserID uint64 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
}

func (m *UserLoginReq) Reset()         { *m = UserLoginReq{} }
func (m *UserLoginReq) String() string { return proto.CompactTextString(m) }
func (*UserLoginReq) ProtoMessage()    {}
func (*UserLoginReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_70141d9832fb77cd, []int{2}
}
func (m *UserLoginReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UserLoginReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UserLoginReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UserLoginReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserLoginReq.Merge(m, src)
}
func (m *UserLoginReq) XXX_Size() int {
	return m.Size()
}
func (m *UserLoginReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UserLoginReq.DiscardUnknown(m)
}

var xxx_messageInfo_UserLoginReq proto.InternalMessageInfo

func (m *UserLoginReq) GetUserID() uint64 {
	if m != nil {
		return m.UserID
	}
	return 0
}

type TestRPC struct {
}

func (m *TestRPC) Reset()         { *m = TestRPC{} }
func (m *TestRPC) String() string { return proto.CompactTextString(m) }
func (*TestRPC) ProtoMessage()    {}
func (*TestRPC) Descriptor() ([]byte, []int) {
	return fileDescriptor_70141d9832fb77cd, []int{3}
}
func (m *TestRPC) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TestRPC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TestRPC.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TestRPC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestRPC.Merge(m, src)
}
func (m *TestRPC) XXX_Size() int {
	return m.Size()
}
func (m *TestRPC) XXX_DiscardUnknown() {
	xxx_messageInfo_TestRPC.DiscardUnknown(m)
}

var xxx_messageInfo_TestRPC proto.InternalMessageInfo

type TestRPCRes struct {
	V int32 `protobuf:"varint,1,opt,name=v,proto3" json:"v,omitempty"`
}

func (m *TestRPCRes) Reset()         { *m = TestRPCRes{} }
func (m *TestRPCRes) String() string { return proto.CompactTextString(m) }
func (*TestRPCRes) ProtoMessage()    {}
func (*TestRPCRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_70141d9832fb77cd, []int{4}
}
func (m *TestRPCRes) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TestRPCRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TestRPCRes.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TestRPCRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestRPCRes.Merge(m, src)
}
func (m *TestRPCRes) XXX_Size() int {
	return m.Size()
}
func (m *TestRPCRes) XXX_DiscardUnknown() {
	xxx_messageInfo_TestRPCRes.DiscardUnknown(m)
}

var xxx_messageInfo_TestRPCRes proto.InternalMessageInfo

func (m *TestRPCRes) GetV() int32 {
	if m != nil {
		return m.V
	}
	return 0
}

func init() {
	proto.RegisterType((*C2SPackage)(nil), "pb.C2SPackage")
	proto.RegisterType((*S2CPackage)(nil), "pb.S2CPackage")
	proto.RegisterType((*UserLoginReq)(nil), "pb.UserLoginReq")
	proto.RegisterType((*TestRPC)(nil), "pb.TestRPC")
	proto.RegisterType((*TestRPCRes)(nil), "pb.TestRPCRes")
}

func init() { proto.RegisterFile("pb/tmp/game.proto", fileDescriptor_70141d9832fb77cd) }

var fileDescriptor_70141d9832fb77cd = []byte{
	// 205 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0x48, 0xd2, 0x2f,
	0xc9, 0x2d, 0xd0, 0x4f, 0x4f, 0xcc, 0x4d, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a,
	0x48, 0x92, 0xe2, 0x4a, 0xcf, 0x4f, 0xcf, 0x87, 0xf0, 0x95, 0x2c, 0xb8, 0xb8, 0x9c, 0x8d, 0x82,
	0x03, 0x12, 0x93, 0xb3, 0x13, 0xd3, 0x53, 0x85, 0xc4, 0xb8, 0xd8, 0x4a, 0x8b, 0x53, 0x8b, 0x3c,
	0x5d, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x82, 0xa0, 0x3c, 0x21, 0x21, 0x2e, 0x96, 0xa4, 0xfc,
	0x94, 0x4a, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x30, 0x1b, 0xa4, 0x33, 0xd8, 0xc8, 0x99,
	0x1c, 0x9d, 0x6a, 0x5c, 0x3c, 0xa1, 0xc5, 0xa9, 0x45, 0x3e, 0xf9, 0xe9, 0x99, 0x79, 0x41, 0xa9,
	0x85, 0xb8, 0xf4, 0x2a, 0x71, 0x72, 0xb1, 0x87, 0xa4, 0x16, 0x97, 0x04, 0x05, 0x38, 0x2b, 0x49,
	0x71, 0x71, 0x41, 0x99, 0x41, 0xa9, 0xc5, 0x42, 0x3c, 0x5c, 0x8c, 0x65, 0x60, 0xb5, 0xac, 0x41,
	0x8c, 0x65, 0x4e, 0x2a, 0x17, 0x1e, 0xca, 0x31, 0x9c, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c,
	0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb, 0x31, 0xcc, 0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1,
	0x1c, 0xc3, 0x8d, 0xc7, 0x72, 0x0c, 0x51, 0x4c, 0x05, 0x49, 0x49, 0x6c, 0x60, 0xff, 0x1a, 0x03,
	0x02, 0x00, 0x00, 0xff, 0xff, 0xad, 0xf8, 0xfb, 0x63, 0x14, 0x01, 0x00, 0x00,
}

func (m *C2SPackage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *C2SPackage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *C2SPackage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Body) > 0 {
		i -= len(m.Body)
		copy(dAtA[i:], m.Body)
		i = encodeVarintGame(dAtA, i, uint64(len(m.Body)))
		i--
		dAtA[i] = 0x12
	}
	if m.UserID != 0 {
		i = encodeVarintGame(dAtA, i, uint64(m.UserID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *S2CPackage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2CPackage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2CPackage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Body) > 0 {
		i -= len(m.Body)
		copy(dAtA[i:], m.Body)
		i = encodeVarintGame(dAtA, i, uint64(len(m.Body)))
		i--
		dAtA[i] = 0x12
	}
	if m.UserID != 0 {
		i = encodeVarintGame(dAtA, i, uint64(m.UserID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *UserLoginReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UserLoginReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UserLoginReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.UserID != 0 {
		i = encodeVarintGame(dAtA, i, uint64(m.UserID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TestRPC) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TestRPC) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TestRPC) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *TestRPCRes) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TestRPCRes) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TestRPCRes) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.V != 0 {
		i = encodeVarintGame(dAtA, i, uint64(m.V))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintGame(dAtA []byte, offset int, v uint64) int {
	offset -= sovGame(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *C2SPackage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.UserID != 0 {
		n += 1 + sovGame(uint64(m.UserID))
	}
	l = len(m.Body)
	if l > 0 {
		n += 1 + l + sovGame(uint64(l))
	}
	return n
}

func (m *S2CPackage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.UserID != 0 {
		n += 1 + sovGame(uint64(m.UserID))
	}
	l = len(m.Body)
	if l > 0 {
		n += 1 + l + sovGame(uint64(l))
	}
	return n
}

func (m *UserLoginReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.UserID != 0 {
		n += 1 + sovGame(uint64(m.UserID))
	}
	return n
}

func (m *TestRPC) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *TestRPCRes) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.V != 0 {
		n += 1 + sovGame(uint64(m.V))
	}
	return n
}

func sovGame(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGame(x uint64) (n int) {
	return sovGame(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *C2SPackage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGame
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: C2SPackage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: C2SPackage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserID", wireType)
			}
			m.UserID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGame
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UserID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Body", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGame
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGame
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGame
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Body = append(m.Body[:0], dAtA[iNdEx:postIndex]...)
			if m.Body == nil {
				m.Body = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGame(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGame
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
func (m *S2CPackage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGame
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: S2CPackage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: S2CPackage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserID", wireType)
			}
			m.UserID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGame
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UserID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Body", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGame
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGame
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGame
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Body = append(m.Body[:0], dAtA[iNdEx:postIndex]...)
			if m.Body == nil {
				m.Body = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGame(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGame
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
func (m *UserLoginReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGame
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: UserLoginReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UserLoginReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserID", wireType)
			}
			m.UserID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGame
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UserID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGame(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGame
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
func (m *TestRPC) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGame
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TestRPC: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TestRPC: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipGame(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGame
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
func (m *TestRPCRes) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGame
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TestRPCRes: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TestRPCRes: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field V", wireType)
			}
			m.V = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGame
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.V |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGame(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGame
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
func skipGame(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGame
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
					return 0, ErrIntOverflowGame
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGame
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
			if length < 0 {
				return 0, ErrInvalidLengthGame
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGame
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGame
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGame        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGame          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGame = fmt.Errorf("proto: unexpected end of group")
)
