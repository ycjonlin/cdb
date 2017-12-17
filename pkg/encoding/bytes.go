package encoding

import (
	"bytes"
	"errors"
	"math"
	"reflect"
	"unsafe"
)

// Err ...
var (
	ErrUnexpectedEOF     = errors.New("unexpected EOF")
	ErrUnknownTag        = errors.New("unknown tag")
	ErrSizeLimitExceeded = errors.New("size limit exceeded")
)

// Bytes ...
type Bytes []byte

// EncodeUint8 ...
func (b *Bytes) EncodeUint8(v uint8) {
	*b = append(*b, byte(v))
}

// EncodeUint16 ...
func (b *Bytes) EncodeUint16(v uint16) {
	*b = append(*b, byte(v>>8), byte(v))
}

// EncodeUint32 ...
func (b *Bytes) EncodeUint32(v uint32) {
	*b = append(*b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// EncodeUint64 ...
func (b *Bytes) EncodeUint64(v uint64) {
	*b = append(*b, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
		byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
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
		*b = append(*b, 0xf0|byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8),
			byte(v))
		return
	case v>>42 == 0:
		*b = append(*b, 0xf8|byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16),
			byte(v>>8), byte(v))
		return
	case v>>49 == 0:
		*b = append(*b, 0xfc|byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24),
			byte(v>>16), byte(v>>8), byte(v))
		return
	case v>>56 == 0:
		*b = append(*b, 0xfe, byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24),
			byte(v>>16), byte(v>>8), byte(v))
		return
	case v>>63 == 0:
		*b = append(*b, 0xff, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
			byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	default:
		*b = append(*b, 0xff, 0x80, byte(v>>56), byte(v>>48), byte(v>>40),
			byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	}
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
		*b = append(*b, 0xf8^byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8),
			byte(v))
		return
	case u>>41 == 0:
		*b = append(*b, 0xfc^byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16),
			byte(v>>8), byte(v))
		return
	case u>>48 == 0:
		*b = append(*b, 0xfe^s, byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16),
			byte(v>>8), byte(v))
		return
	case u>>55 == 0:
		*b = append(*b, 0xff^s, byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24),
			byte(v>>16), byte(v>>8), byte(v))
		return
	case u>>62 == 0:
		*b = append(*b, 0xff^s, 0x80^byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
			byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	default:
		*b = append(*b, 0xff^s, 0xc0^s, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
			byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		return
	}
}

// EncodeFloat32 ...
func (b *Bytes) EncodeFloat32(v float32) {
	n := math.Float32bits(v)
	if n&(1<<31) == 0 {
		n ^= (1 << 31)
	} else {
		n = ^n
	}
	b.EncodeUint32(n)
}

// EncodeFloat64 ...
func (b *Bytes) EncodeFloat64(v float64) {
	n := math.Float64bits(v)
	if n&(1<<63) == 0 {
		n ^= (1 << 63)
	} else {
		n = ^n
	}
	b.EncodeUint64(n)
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

// EncodeNonsortingBytes ...
func (b *Bytes) EncodeNonsortingBytes(v []byte) {
	b.EncodeUvarint(uint64(len(v)))
	*b = append(*b, v...)
}

// EncodeBool ...
func (b *Bytes) EncodeBool(v bool) {
	if v {
		*b = append(*b, 0xff)
	} else {
		*b = append(*b, 0x00)
	}
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

// EncodeSize ...
func (b *Bytes) EncodeSize(v int) {
	b.EncodeUvarint(uint64(v))
}

////////////////////////
// Decode
////////////////////////

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
		if l&(1<<uint((^i)&7)) == 0 {
			v &= (1<<uint((i<<3)-i+7) - 1)
			*b = a[i+1:]
			return v, nil
		}
	}
	return 0, ErrUnexpectedEOF
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

// DecodeFloat32 ...
func (b *Bytes) DecodeFloat32() (float32, error) {
	n, err := b.DecodeUint32()
	if err != nil {
		return 0, ErrUnexpectedEOF
	}
	if n&(1<<31) != 0 {
		n ^= (1 << 31)
	} else {
		n = ^n
	}
	v := math.Float32frombits(n)
	return v, nil
}

// DecodeFloat64 ...
func (b *Bytes) DecodeFloat64() (float64, error) {
	n, err := b.DecodeUint64()
	if err != nil {
		return 0, ErrUnexpectedEOF
	}
	if n&(1<<63) != 0 {
		n ^= (1 << 63)
	} else {
		n = ^n
	}
	v := math.Float64frombits(n)
	return v, nil
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

// DecodeNonsortingBytes ...
func (b *Bytes) DecodeNonsortingBytes(v []byte) ([]byte, error) {
	u, err := b.DecodeUvarint()
	if err != nil || u>>31 != 0 {
		return v, ErrUnexpectedEOF
	}
	n := int(u)
	a := *b
	if len(a) < n {
		return v, ErrUnexpectedEOF
	}
	v = append(v, a[:n]...)
	*b = a[n:]
	return v, nil
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

// DecodeSortingString ...
func (b *Bytes) DecodeSortingString(v []byte) (string, []byte, error) {
	v, err := b.DecodeSortingBytes(v)
	return *(*string)(unsafe.Pointer(&v)), v, err
}

// DecodeNonsortingString ...
func (b *Bytes) DecodeNonsortingString(v []byte) (string, []byte, error) {
	v, err := b.DecodeNonsortingBytes(v)
	return *(*string)(unsafe.Pointer(&v)), v, err
}

// DecodeSize ...
func (b *Bytes) DecodeSize() (int, error) {
	v, err := b.DecodeUvarint()
	if err != nil {
		return 0, err
	}
	if v>>31 != 0 {
		return 0, ErrSizeLimitExceeded
	}
	return int(v), nil
}
