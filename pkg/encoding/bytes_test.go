package encoding

import (
	"bytes"
	"fmt"
	"math"
	"runtime"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BytesList []Bytes

func (bs BytesList) Len() int {
	return len(bs)
}
func (bs BytesList) Less(i, j int) bool {
	return bytes.Compare(bs[i], bs[j]) < 0
}
func (bs BytesList) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

type Stats struct {
	mem0 runtime.MemStats
	mem1 runtime.MemStats

	Mallocs uint64
}

func NewStats() *Stats {
	var mem runtime.MemStats
	for i := 0; i < 8; i++ {
		runtime.ReadMemStats(&mem)
	}
	return &Stats{}
}

func (s *Stats) Record(callback func()) {
	runtime.ReadMemStats(&s.mem0)
	callback()
	runtime.ReadMemStats(&s.mem1)
	s.Mallocs += s.mem1.Mallocs - s.mem0.Mallocs
}

func TestUint8(t *testing.T) {
	var vs []uint8
	for i := 0; i <= 8; i++ {
		v := uint8(1 << uint(i))
		vs = append(vs, v, ^v)
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return vs[i] < vs[j]
	})
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, 1)
		vp, err := v, error(nil)
		stats.Record(func() {
			// encode
			b.EncodeUint8(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeUint8()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}

func TestUint16(t *testing.T) {
	var vs []uint16
	for i := 0; i <= 16; i++ {
		v := uint16(1 << uint(i))
		vs = append(vs, v, ^v)
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return vs[i] < vs[j]
	})
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, 2)
		vp, err := v, error(nil)
		stats.Record(func() {
			// encode
			b.EncodeUint16(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeUint16()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}

func TestUint32(t *testing.T) {
	var vs []uint32
	for i := 0; i <= 32; i++ {
		v := uint32(1 << uint(i))
		vs = append(vs, v, ^v)
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return vs[i] < vs[j]
	})
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, 4)
		vp, err := v, error(nil)
		stats.Record(func() {
			// encode
			b.EncodeUint32(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeUint32()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}

func TestUint64(t *testing.T) {
	var vs []uint64
	for i := 0; i <= 64; i++ {
		v := uint64(1 << uint(i))
		vs = append(vs, v, ^v)
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return vs[i] < vs[j]
	})
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, 8)
		vp, err := v, error(nil)
		stats.Record(func() {
			// encode
			b.EncodeUint64(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeUint64()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}

func TestUvarint(t *testing.T) {
	var vs []uint64
	for i := 0; i <= 64; i++ {
		v := uint64(1 << uint(i))
		vs = append(vs, v, ^v)
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return vs[i] < vs[j]
	})
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, 10)
		vp, err := v, error(nil)
		stats.Record(func() {
			// encode
			b.EncodeUvarint(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeUvarint()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}

func TestVarint(t *testing.T) {
	var vs []int64
	for i := 0; i <= 64; i++ {
		v := int64(1 << uint(i))
		vs = append(vs, v, ^v)
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return vs[i] < vs[j]
	})
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, 10)
		vp, err := v, error(nil)
		stats.Record(func() {
			// encode
			b.EncodeVarint(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeVarint()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}

func TestVarfloat(t *testing.T) {
	var vs []uint64
	for i := 0; i <= 64; i++ {
		v := uint64(1 << uint(i))
		vs = append(vs, v, ^v)
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return vs[i] < vs[j]
	})
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, 10)
		vp, err := v, error(nil)
		stats.Record(func() {
			// encode
			b.encodeVarfloat(v)
			bs[i] = b
			// decode
			vp, err = b.decodeVarfloat()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	for i := 0; i < len(bs)-1; i++ {
		if bytes.Compare(bs[i], bs[i+1]) > 0 {
			fmt.Println(bs[i], bs[i+1])
		}
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}

func TestFloat32(t *testing.T) {
	var vs []float32
	for i := 0; i < 32; i++ {
		v := uint32(1 << uint(i))
		vs = append(vs,
			math.Float32frombits(v),
			math.Float32frombits(v^(1<<31)))
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return vs[i] < vs[j]
	})
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, 10)
		vp, err := v, error(nil)
		stats.Record(func() {
			// encode
			b.EncodeFloat32(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeFloat32()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}

func TestFloat64(t *testing.T) {
	var vs []float64
	for i := 0; i < 64; i++ {
		v := uint64(1 << uint(i))
		vs = append(vs,
			math.Float64frombits(v),
			math.Float64frombits(v^(1<<63)))
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return vs[i] < vs[j]
	})
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, 10)
		vp, err := v, error(nil)
		stats.Record(func() {
			// encode
			b.EncodeFloat64(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeFloat64()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}

func TestSortingBytes(t *testing.T) {
	var vs [][]byte
	vs = append(vs, []byte{})
	j0, j1 := 0, 1
	for i := 0; i < 4; i++ {
		for j := j0; j < j1; j++ {
			v := vs[j]
			vs = append(vs,
				append([]byte{0x00}, v...),
				append([]byte{0x01}, v...))
		}
		j0, j1 = j1, len(vs)
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return bytes.Compare(vs[i], vs[j]) < 0
	})
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, len(v)*2+2)
		vp, err := make([]byte, 0, len(v)), error(nil)
		stats.Record(func() {
			// encode
			b.EncodeSortingBytes(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeSortingBytes(vp)
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}

func TestNonsortingBytes(t *testing.T) {
	var vs [][]byte
	vs = append(vs, []byte{})
	j0, j1 := 0, 1
	for i := 0; i < 4; i++ {
		for j := j0; j < j1; j++ {
			v := vs[j]
			vs = append(vs,
				append([]byte{0x00}, v...),
				append([]byte{0x01}, v...))
		}
		j0, j1 = j1, len(vs)
	}
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, len(v)+2)
		vp, err := make([]byte, 0, len(v)), error(nil)
		stats.Record(func() {
			// encode
			b.EncodeNonsortingBytes(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeNonsortingBytes(vp)
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
}

func TestBool(t *testing.T) {
	vs := []bool{false, true}
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, 1)
		vp, err := v, error(nil)
		stats.Record(func() {
			// encode
			b.EncodeBool(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeBool()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}

func TestSortingString(t *testing.T) {
	var vs []string
	vs = append(vs, "")
	j0, j1 := 0, 1
	for i := 0; i < 4; i++ {
		for j := j0; j < j1; j++ {
			v := vs[j]
			vs = append(vs, "\x00"+v, "\x01"+v)
		}
		j0, j1 = j1, len(vs)
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return vs[i] < vs[j]
	})
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, len(v)*2+2)
		vp, err := "", error(nil)
		stats.Record(func() {
			// encode
			b.EncodeSortingString(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeSortingString()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	//assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}

func TestNonsortingString(t *testing.T) {
	var vs []string
	vs = append(vs, "")
	j0, j1 := 0, 1
	for i := 0; i < 4; i++ {
		for j := j0; j < j1; j++ {
			v := vs[j]
			vs = append(vs, "\x00"+v, "\x01"+v)
		}
		j0, j1 = j1, len(vs)
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return vs[i] < vs[j]
	})
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, len(v)+2)
		vp, err := "", error(nil)
		stats.Record(func() {
			// encode
			b.EncodeNonsortingString(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeNonsortingString()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	//assert.Zero(t, stats.Mallocs)
}

func TestSize(t *testing.T) {
	vs := []int{}
	for i := 0; i < 31; i++ {
		v := int(1 << uint(i))
		vs = append(vs, v)
	}
	stats := NewStats()
	bs := make(BytesList, len(vs))
	for i, v := range vs {
		b := make(Bytes, 0, 9)
		vp, err := v, error(nil)
		stats.Record(func() {
			// encode
			b.EncodeSize(v)
			bs[i] = b
			// decode
			vp, err = b.DecodeSize()
		})
		assert.Nil(t, err)
		assert.Equal(t, v, vp)
		assert.Empty(t, b)
	}
	// mem stats
	assert.Zero(t, stats.Mallocs)
	// order
	assert.True(t, sort.IsSorted(bs))
}
