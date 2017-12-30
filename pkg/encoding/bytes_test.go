package encoding_test

import (
	"bytes"
	"math"
	"runtime"
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	enc "github.com/ycjonlin/cdb/pkg/encoding"
)

var _ = Describe("Bytes", func() {
	It("Uint8", func() {

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
			b := make(enc.Bytes, 0, 1)
			vp, err := v, error(nil)
			stats.Record(func() {
				b.EncodeUint8(v)
				bs[i] = b
				vp, err = b.DecodeUint8()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
	It("Uint16", func() {

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
			b := make(enc.Bytes, 0, 2)
			vp, err := v, error(nil)
			stats.Record(func() {
				b.EncodeUint16(v)
				bs[i] = b
				vp, err = b.DecodeUint16()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
	It("Uint32", func() {

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
			b := make(enc.Bytes, 0, 4)
			vp, err := v, error(nil)
			stats.Record(func() {
				b.EncodeUint32(v)
				bs[i] = b
				vp, err = b.DecodeUint32()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
	It("Uint64", func() {

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
			b := make(enc.Bytes, 0, 8)
			vp, err := v, error(nil)
			stats.Record(func() {
				b.EncodeUint64(v)
				bs[i] = b
				vp, err = b.DecodeUint64()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
	It("Uvarint", func() {

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
			b := make(enc.Bytes, 0, 10)
			vp, err := v, error(nil)
			stats.Record(func() {
				b.EncodeUvarint(v)
				bs[i] = b
				vp, err = b.DecodeUvarint()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
	It("Varint", func() {

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
			b := make(enc.Bytes, 0, 10)
			vp, err := v, error(nil)
			stats.Record(func() {
				b.EncodeVarint(v)
				bs[i] = b
				vp, err = b.DecodeVarint()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
	It("Float32", func() {

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
			b := make(enc.Bytes, 0, 10)
			vp, err := v, error(nil)
			stats.Record(func() {
				b.EncodeFloat32(v)
				bs[i] = b
				vp, err = b.DecodeFloat32()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
	It("Float64", func() {

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
			b := make(enc.Bytes, 0, 10)
			vp, err := v, error(nil)
			stats.Record(func() {
				b.EncodeFloat64(v)
				bs[i] = b
				vp, err = b.DecodeFloat64()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
	It("Varfloat", func() {

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
			b := make(enc.Bytes, 0, 10)
			vp, err := v, error(nil)
			stats.Record(func() {
				b.EncodeVarfloat(v)
				bs[i] = b
				vp, err = b.DecodeVarfloat()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
	It("Sorting Bytes", func() {

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
			b := make(enc.Bytes, 0, len(v)*2+2)
			vp, err := make([]byte, 0, len(v)), error(nil)
			stats.Record(func() {
				b.EncodeSortingBytes(v)
				bs[i] = b
				vp, err = b.DecodeSortingBytes(vp)
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
	It("Nonsorting Bytes", func() {

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
			b := make(enc.Bytes, 0, len(v)+2)
			vp, err := make([]byte, 0, len(v)), error(nil)
			stats.Record(func() {
				b.EncodeNonsortingBytes(v)
				bs[i] = b
				vp, err = b.DecodeNonsortingBytes(vp)
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
	})
	It("Bool", func() {

		vs := []bool{false, true}
		stats := NewStats()
		bs := make(BytesList, len(vs))
		for i, v := range vs {
			b := make(enc.Bytes, 0, 1)
			vp, err := v, error(nil)
			stats.Record(func() {
				b.EncodeBool(v)
				bs[i] = b
				vp, err = b.DecodeBool()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
	It("Sorting String", func() {

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
			b := make(enc.Bytes, 0, len(v)*2+2)
			vp, err := "", error(nil)
			stats.Record(func() {
				b.EncodeSortingString(v)
				bs[i] = b
				vp, err = b.DecodeSortingString()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
	It("Nonsorting String", func() {

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
			b := make(enc.Bytes, 0, len(v)+2)
			vp, err := "", error(nil)
			stats.Record(func() {
				b.EncodeNonsortingString(v)
				bs[i] = b
				vp, err = b.DecodeNonsortingString()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
	})
	It("Size", func() {

		vs := []int{}
		for i := 0; i < 31; i++ {
			v := int(1 << uint(i))
			vs = append(vs, v)
		}
		stats := NewStats()
		bs := make(BytesList, len(vs))
		for i, v := range vs {
			b := make(enc.Bytes, 0, 9)
			vp, err := v, error(nil)
			stats.Record(func() {
				b.EncodeSize(v)
				bs[i] = b
				vp, err = b.DecodeSize()
			})
			Expect(err).To(BeNil())
			Expect(vp).To(Equal(v))
			Expect(b).To(BeEmpty())
		}
		Expect(stats.Mallocs).To(BeZero())
		Expect(sort.IsSorted(bs)).To(BeTrue())
	})
})

type BytesList []enc.Bytes

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
	mem0    runtime.MemStats
	mem1    runtime.MemStats
	Mallocs uint64
}

func NewStats() *Stats {
	var mem runtime.MemStats
	for i := 0; i < 64; i++ {
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
