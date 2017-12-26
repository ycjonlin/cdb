package encoding

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math"
	"reflect"
	"unsafe"
)

// Err ...
var (
	ErrUnexpectedEOF     = errors.New("unexpected EOF")
	ErrUnknownTag        = errors.New("unknown tag")
	ErrExceededSizeLimit = errors.New("exceeded size limit")
	ErrDataOverflow      = errors.New("data overflow")
)

// Bytes ...
type Bytes []byte

// String ...
func (b Bytes) String() string {
	return hex.EncodeToString(b)
}

// EncodeUint8 ...
func (b *Bytes) EncodeUint8(v uint8) {
	*b = append(*b, byte(v))
}

// DecodeUint8 ...
func (b *Bytes) DecodeUint8() (uint8, error) {
	a := *b
	if len(a) < 1 {
		return 0, ErrUnexpectedEOF
	}
	*b = a[1:]
	v := uint8(a[0])
	return v, nil
}

// EncodeUint16 ...
func (b *Bytes) EncodeUint16(v uint16) {
	*b = append(*b, byte(v>>8), byte(v))
}

// DecodeUint16 ...
func (b *Bytes) DecodeUint16() (uint16, error) {
	a := *b
	if len(a) < 2 {
		return 0, ErrUnexpectedEOF
	}
	*b = a[2:]
	v := uint16(a[0])<<8 | uint16(a[1])
	return v, nil
}

// EncodeUint32 ...
func (b *Bytes) EncodeUint32(v uint32) {
	*b = append(*b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// DecodeUint32 ...
func (b *Bytes) DecodeUint32() (uint32, error) {
	a := *b
	if len(a) < 4 {
		return 0, ErrUnexpectedEOF
	}
	*b = a[4:]
	v := uint32(a[0])<<24 | uint32(a[1])<<16 | uint32(a[2])<<8 | uint32(a[3])
	return v, nil
}

// EncodeUint64 ...
func (b *Bytes) EncodeUint64(v uint64) {
	*b = append(*b, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
		byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// DecodeUint64 ...
func (b *Bytes) DecodeUint64() (uint64, error) {
	a := *b
	if len(a) < 8 {
		return 0, ErrUnexpectedEOF
	}
	*b = a[8:]
	v := uint64(a[0])<<56 | uint64(a[1])<<48 | uint64(a[2])<<40 | uint64(a[3])<<32 |
		uint64(a[4])<<24 | uint64(a[5])<<16 | uint64(a[6])<<8 | uint64(a[7])
	return v, nil
}

// EncodeUvarint ...
func (b *Bytes) EncodeUvarint(v uint64) {
	switch {
	case v>>7 == 0:
		*b = append(*b, byte(v))
		return
	case v>>14 == 0:
		*b = append(*b, 0x80|byte(v>>8), byte(v))
		return
	case v>>21 == 0:
		*b = append(*b, 0xc0|byte(v>>16), byte(v>>8), byte(v))
		return
	case v>>28 == 0:
		*b = append(*b, 0xe0|byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	case v>>35 == 0:
		*b = append(*b, 0xf0|byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	case v>>42 == 0:
		*b = append(*b, 0xf8|byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	case v>>49 == 0:
		*b = append(*b, 0xfc|byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	case v>>56 == 0:
		*b = append(*b, 0xfe, byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	case v>>63 == 0:
		*b = append(*b, 0xff, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	default:
		*b = append(*b, 0xff, 0x80, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	}
}

// DecodeUvarint ...
func (b *Bytes) DecodeUvarint() (uint64, error) {
	a := *b
	v := uint64(0)
	l := byte(0)
	for i, n := 0, len(a); i < n; i++ {
		v = v<<8 | uint64(a[i])
		if i&7 == 0 {
			l = a[i>>3]
		}
		if i == 9 && l != 0x80 {
			return 0, ErrDataOverflow
		}
		if l&(1<<uint((^i)&7)) == 0 {
			v &= (1<<uint((i<<3)-i+7) - 1)
			*b = a[i+1:]
			return v, nil
		}
	}
	return 0, ErrUnexpectedEOF
}

// EncodeVarint ...
func (b *Bytes) EncodeVarint(v int64) {
	u := uint64(v)
	s := byte(0)
	if v < 0 {
		u = ^u
		s = 0xff
	}
	switch {
	case u>>6 == 0:
		*b = append(*b, 0x80^byte(v))
		return
	case u>>13 == 0:
		*b = append(*b, 0xc0^byte(v>>8), byte(v))
		return
	case u>>20 == 0:
		*b = append(*b, 0xe0^byte(v>>16), byte(v>>8), byte(v))
		return
	case u>>27 == 0:
		*b = append(*b, 0xf0^byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	case u>>34 == 0:
		*b = append(*b, 0xf8^byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	case u>>41 == 0:
		*b = append(*b, 0xfc^byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	case u>>48 == 0:
		*b = append(*b, 0xfe^s, byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	case u>>55 == 0:
		*b = append(*b, 0xff^s, byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	case u>>62 == 0:
		*b = append(*b, 0xff^s, 0x80^byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	default:
		*b = append(*b, 0xff^s, 0xc0^s, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	}
}

// DecodeVarint ...
func (b *Bytes) DecodeVarint() (int64, error) {
	a := *b
	v := uint64(0)
	l := byte(0)
	for i, n := 0, len(a); i < n; i++ {
		v = v<<8 | uint64(a[i])
		if i&7 == 0 {
			l = a[i>>3]
			l ^= l << 1
		} else if i&7 == 7 {
			l ^= a[i>>3+1] >> 7
		}
		if i == 9 && l^a[0]>>7 != 0x41 {
			return 0, ErrDataOverflow
		}
		if l&(1<<uint((^i)&7)) != 0 {
			if a[0]&0x80 != 0 {
				v &= (1<<uint((i<<3)-i+7) - 1)
			} else {
				v |= ^(1<<uint((i<<3)-i+7) - 1)
			}
			*b = a[i+1:]
			return int64(v), nil
		}
	}
	return 0, ErrUnexpectedEOF
}

var c = [128]byte{
	0x01, 0x02, 0x05, 0x06, 0x09, 0x0a, 0x0d, 0x0e, 0x11, 0x12, 0x15, 0x16, 0x19, 0x1a, 0x1d, 0x1e,
	0x21, 0x22, 0x25, 0x26, 0x29, 0x2a, 0x2d, 0x2e, 0x31, 0x32, 0x35, 0x36, 0x39, 0x3a, 0x3d, 0x3e,
	0x41, 0x42, 0x45, 0x46, 0x49, 0x4a, 0x4d, 0x4e, 0x51, 0x52, 0x55, 0x56, 0x59, 0x5a, 0x5d, 0x5e,
	0x61, 0x62, 0x65, 0x66, 0x69, 0x6a, 0x6d, 0x6e, 0x71, 0x72, 0x75, 0x76, 0x79, 0x7a, 0x7d, 0x7e,
	0x81, 0x82, 0x85, 0x86, 0x89, 0x8a, 0x8d, 0x8e, 0x91, 0x92, 0x95, 0x96, 0x99, 0x9a, 0x9d, 0x9e,
	0xa1, 0xa2, 0xa5, 0xa6, 0xa9, 0xaa, 0xad, 0xae, 0xb1, 0xb2, 0xb5, 0xb6, 0xb9, 0xba, 0xbd, 0xbe,
	0xc1, 0xc2, 0xc5, 0xc6, 0xc9, 0xca, 0xcd, 0xce, 0xd1, 0xd2, 0xd5, 0xd6, 0xd9, 0xda, 0xdd, 0xde,
	0xe1, 0xe2, 0xe5, 0xe6, 0xe9, 0xea, 0xed, 0xee, 0xf1, 0xf2, 0xf5, 0xf6, 0xf9, 0xfa, 0xfd, 0xfe,
}

func (b *Bytes) encodeVarfloat(v uint64) {
	u := v
	if v&1 != 0 {
		u = ^u
	}
	switch {
	case u<<6 == 0:
		*b = append(*b, byte(v>>56))
		return
	case u<<13 == 0:
		*b = append(*b, c[(v>>57)&0x7f], byte(v>>49))
		return
	case u<<20 == 0:
		*b = append(*b, c[(v>>57)&0x7f], c[(v>>50)&0x7f], byte(v>>42))
		return
	case u<<27 == 0:
		*b = append(*b, c[(v>>57)&0x7f], c[(v>>50)&0x7f], c[(v>>43)&0x7f], byte(v>>35))
		return
	case u<<34 == 0:
		*b = append(*b, c[(v>>57)&0x7f], c[(v>>50)&0x7f], c[(v>>43)&0x7f], c[(v>>36)&0x7f], byte(v>>28))
		return
	case u<<41 == 0:
		*b = append(*b, c[(v>>57)&0x7f], c[(v>>50)&0x7f], c[(v>>43)&0x7f], c[(v>>36)&0x7f], c[(v>>29)&0x7f], byte(v>>21))
		return
	case u<<48 == 0:
		*b = append(*b, c[(v>>57)&0x7f], c[(v>>50)&0x7f], c[(v>>43)&0x7f], c[(v>>36)&0x7f], c[(v>>29)&0x7f], c[(v>>22)&0x7f], byte(v>>14))
		return
	case u<<55 == 0:
		*b = append(*b, c[(v>>57)&0x7f], c[(v>>50)&0x7f], c[(v>>43)&0x7f], c[(v>>36)&0x7f], c[(v>>29)&0x7f], c[(v>>22)&0x7f], c[(v>>15)&0x7f], byte(v>>7))
		return
	case u<<62 == 0:
		*b = append(*b, c[(v>>57)&0x7f], c[(v>>50)&0x7f], c[(v>>43)&0x7f], c[(v>>36)&0x7f], c[(v>>29)&0x7f], c[(v>>22)&0x7f], c[(v>>15)&0x7f], c[(v>>8)&0x7f], byte(v))
		return
	default:
		*b = append(*b, c[(v>>57)&0x7f], c[(v>>50)&0x7f], c[(v>>43)&0x7f], c[(v>>36)&0x7f], c[(v>>29)&0x7f], c[(v>>22)&0x7f], c[(v>>15)&0x7f], c[(v>>8)&0x7f], c[(v>>1)&0x7f], byte(v<<7))
		return

	}
}

func (b *Bytes) decodeVarfloat() (uint64, error) {
	a := *b
	v := uint64(0)
	for i, n := 0, len(a); i < n; i++ {
		if i != 9 {
			v = v<<7 | uint64(a[i]>>1)
		} else {
			if a[i]&0x7f != 0 {
				return 0, ErrDataOverflow
			}
			v = v<<1 | uint64(a[i]>>7)
			*b = a[i+1:]
			return v, nil
		}
		if (a[i]^(a[i]>>1))&0x01 == 0 {
			v <<= uint(57 - (i<<3 - i))
			if a[i]&0x02 != 0 {
				v |= 1<<uint(57-(i<<3-i)) - 1
			}
			*b = a[i+1:]
			return v, nil
		}
	}
	return 0, ErrUnexpectedEOF
}

// EncodeFloat32 ...
func (b *Bytes) EncodeFloat32(v float32) {
	m := math.Float32bits(v)
	n := uint64(m) << 32
	if n&(1<<63) == 0 {
		n = n ^ 1<<63
	} else {
		n = ^n
	}
	b.encodeVarfloat(n)
}

// DecodeFloat32 ...
func (b *Bytes) DecodeFloat32() (float32, error) {
	n, err := b.decodeVarfloat()
	if err != nil {
		return 0, ErrUnexpectedEOF
	}
	if n&(1<<63) != 0 {
		n = n ^ 1<<63
	} else {
		n = ^n
	}
	if n<<32 != 0 {
		return 0, ErrDataOverflow
	}
	v := math.Float32frombits(uint32(n >> 32))
	return v, nil
}

// EncodeFloat64 ...
func (b *Bytes) EncodeFloat64(v float64) {
	n := math.Float64bits(v)
	if n&(1<<63) == 0 {
		n = n ^ 1<<63
	} else {
		n = ^n
	}
	b.encodeVarfloat(n)
}

// DecodeFloat64 ...
func (b *Bytes) DecodeFloat64() (float64, error) {
	n, err := b.decodeVarfloat()
	if err != nil {
		return 0, ErrUnexpectedEOF
	}
	if n&(1<<63) != 0 {
		n = n ^ 1<<63
	} else {
		n = ^n
	}
	v := math.Float64frombits(n)
	return v, nil
}

// EncodeSortingBytes ...
func (b *Bytes) EncodeSortingBytes(v []byte) {
	a := *b
	for {
		i := bytes.IndexByte(v, 0x00)
		if i == -1 {
			a = append(a, v...)
			a = append(a, 0x00, 0x00)
			*b = a
			return
		}
		a = append(a, v[:i]...)
		a = append(a, 0x00, 0x01)
		v = v[i+1:]
	}
}

// DecodeSortingBytes ...
func (b *Bytes) DecodeSortingBytes(v []byte) ([]byte, error) {
	a := *b
	for {
		i := bytes.IndexByte(a, 0x00)
		if i == -1 {
			return v, ErrUnexpectedEOF
		}
		if len(a) < i+2 {
			return v, ErrUnexpectedEOF
		}
		switch a[i+1] {
		case 0x00:
			v = append(v, a[:i]...)
			*b = a[i+2:]
			return v, nil
		case 0x01:
			v = append(v, a[:i+1]...)
			a = a[i+2:]
		default:
			return v, ErrUnexpectedEOF
		}
	}
}

// EncodeNonsortingBytes ...
func (b *Bytes) EncodeNonsortingBytes(v []byte) {
	b.EncodeSize(len(v))
	*b = append(*b, v...)
}

// DecodeNonsortingBytes ...
func (b *Bytes) DecodeNonsortingBytes(v []byte) ([]byte, error) {
	n, err := b.DecodeSize()
	if err != nil {
		return v, ErrUnexpectedEOF
	}
	a := *b
	if len(a) < n {
		return v, ErrUnexpectedEOF
	}
	v = append(v, a[:n]...)
	*b = a[n:]
	return v, nil
}

// EncodeBool ...
func (b *Bytes) EncodeBool(v bool) {
	if v {
		*b = append(*b, 0xff)
	} else {
		*b = append(*b, 0x00)
	}
}

// DecodeBool ...
func (b *Bytes) DecodeBool() (bool, error) {
	a := *b
	if len(a) < 1 {
		return false, ErrUnexpectedEOF
	}
	v := a[0] != 0
	*b = a[1:]
	return v, nil
}

// EncodeSortingString ...
func (b *Bytes) EncodeSortingString(v string) {
	if len(v) == 0 {
		b.EncodeSortingBytes(nil)
		return
	}
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&v))
	arg := (*[0x7fffffff]byte)(unsafe.Pointer(hdr.Data))[:len(v):len(v)]
	b.EncodeSortingBytes(arg)
}

// DecodeSortingString ...
func (b *Bytes) DecodeSortingString() (string, error) {
	v, err := b.DecodeSortingBytes(nil)
	return *(*string)(unsafe.Pointer(&v)), err
}

// EncodeNonsortingString ...
func (b *Bytes) EncodeNonsortingString(v string) {
	if len(v) == 0 {
		b.EncodeNonsortingBytes(nil)
		return
	}
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&v))
	arg := (*[0x7fffffff]byte)(unsafe.Pointer(hdr.Data))[:len(v):len(v)]
	b.EncodeNonsortingBytes(arg)
}

// DecodeNonsortingString ...
func (b *Bytes) DecodeNonsortingString() (string, error) {
	v, err := b.DecodeNonsortingBytes(nil)
	return *(*string)(unsafe.Pointer(&v)), err
}

// MaxSize ...
const MaxSize = 1<<31 - 1

// EncodeSize ...
func (b *Bytes) EncodeSize(v int) {
	b.EncodeUvarint(uint64(v))
}

// DecodeSize ...
func (b *Bytes) DecodeSize() (int, error) {
	v, err := b.DecodeUvarint()
	if err != nil {
		return 0, err
	}
	if v>>31 != 0 {
		return 0, ErrExceededSizeLimit
	}
	return int(v), nil
}

// EncodeTag ...
func (b *Bytes) EncodeTag(v int, set, del bool) {
	u := uint64(v) << 2
	if set {
		u |= 1
	}
	if del {
		u |= 2
	}
	b.EncodeUvarint(u)
}

// DecodeTag ...
func (b *Bytes) DecodeTag() (tag int, set, del bool, err error) {
	v, err := b.DecodeUvarint()
	if err != nil {
		return 0, false, false, err
	}
	if v>>33 != 0 {
		return 0, false, false, ErrExceededSizeLimit
	}
	return int(v >> 2), v&1 != 0, v&2 != 0, nil
}
