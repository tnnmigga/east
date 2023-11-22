// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: core/pb/package.proto

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

type Package struct {
	ServerID int32  `protobuf:"varint,1,opt,name=ServerID,proto3" json:"ServerID,omitempty"`
	Body     []byte `protobuf:"bytes,2,opt,name=Body,proto3" json:"Body,omitempty"`
}

func (m *Package) Reset()         { *m = Package{} }
func (m *Package) String() string { return proto.CompactTextString(m) }
func (*Package) ProtoMessage()    {}
func (*Package) Descriptor() ([]byte, []int) {
	return fileDescriptor_c0d7f795a8914964, []int{0}
}
func (m *Package) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Package) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Package.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Package) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Package.Merge(m, src)
}
func (m *Package) XXX_Size() int {
	return m.Size()
}
func (m *Package) XXX_DiscardUnknown() {
	xxx_messageInfo_Package.DiscardUnknown(m)
}

var xxx_messageInfo_Package proto.InternalMessageInfo

func (m *Package) GetServerID() int32 {
	if m != nil {
		return m.ServerID
	}
	return 0
}

func (m *Package) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func init() {
	proto.RegisterType((*Package)(nil), "pb.Package")
}

func init() { proto.RegisterFile("core/pb/package.proto", fileDescriptor_c0d7f795a8914964) }

var fileDescriptor_c0d7f795a8914964 = []byte{
	// 178 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4d, 0xce, 0x2f, 0x4a,
	0xd5, 0x2f, 0x48, 0xd2, 0x2f, 0x48, 0x4c, 0xce, 0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0x62, 0x2a, 0x48, 0x92, 0x32, 0x29, 0x4b, 0xcd, 0x4b, 0xc9, 0x2f, 0xd2, 0x4f, 0xcf,
	0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xcf, 0x4f, 0xcf, 0xd7, 0x07, 0xab,
	0x48, 0x2a, 0x4d, 0x03, 0xf3, 0xc0, 0x1c, 0x30, 0x0b, 0xa2, 0x53, 0xc9, 0x92, 0x8b, 0x3d, 0x00,
	0x62, 0x94, 0x90, 0x14, 0x17, 0x47, 0x70, 0x6a, 0x51, 0x59, 0x6a, 0x91, 0xa7, 0x8b, 0x04, 0xa3,
	0x02, 0xa3, 0x06, 0x6b, 0x10, 0x9c, 0x2f, 0x24, 0xc4, 0xc5, 0xe2, 0x94, 0x9f, 0x52, 0x29, 0xc1,
	0xa4, 0xc0, 0xa8, 0xc1, 0x13, 0x04, 0x66, 0x3b, 0xa9, 0x5c, 0x78, 0x28, 0xc7, 0x70, 0xe2, 0x91,
	0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0x4e, 0x78, 0x2c, 0xc7, 0x30, 0xe3,
	0xb1, 0x1c, 0xc3, 0x85, 0xc7, 0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x44, 0x31, 0x15, 0x24, 0x25,
	0xb1, 0x81, 0xed, 0x31, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xc4, 0xbf, 0x88, 0xd4, 0xba, 0x00,
	0x00, 0x00,
}

func (m *Package) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Package) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Package) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Body) > 0 {
		i -= len(m.Body)
		copy(dAtA[i:], m.Body)
		i = encodeVarintPackage(dAtA, i, uint64(len(m.Body)))
		i--
		dAtA[i] = 0x12
	}
	if m.ServerID != 0 {
		i = encodeVarintPackage(dAtA, i, uint64(m.ServerID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintPackage(dAtA []byte, offset int, v uint64) int {
	offset -= sovPackage(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Package) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ServerID != 0 {
		n += 1 + sovPackage(uint64(m.ServerID))
	}
	l = len(m.Body)
	if l > 0 {
		n += 1 + l + sovPackage(uint64(l))
	}
	return n
}

func sovPackage(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPackage(x uint64) (n int) {
	return sovPackage(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Package) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPackage
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
			return fmt.Errorf("proto: Package: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Package: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ServerID", wireType)
			}
			m.ServerID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPackage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ServerID |= int32(b&0x7F) << shift
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
					return ErrIntOverflowPackage
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
				return ErrInvalidLengthPackage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthPackage
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
			skippy, err := skipPackage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPackage
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPackage
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
func skipPackage(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPackage
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
					return 0, ErrIntOverflowPackage
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
					return 0, ErrIntOverflowPackage
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
				return 0, ErrInvalidLengthPackage
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPackage
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPackage
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPackage        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPackage          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPackage = fmt.Errorf("proto: unexpected end of group")
)
