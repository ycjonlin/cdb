package encoding_test

import (
	"fmt"
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	enc "github.com/ycjonlin/cdb/pkg/encoding"
	. "github.com/ycjonlin/cdb/test"
)

var _ = Describe("Marshalling", func() {
	m := TestMarshaller
	for _, c := range []struct {
		n string
		v func() enc.Data
	}{
		{"BytesData", enc.NewBytesData},
		{"MapData", enc.NewMapData},
		{"TreeData", enc.NewTreeData},
	} {
		v := c.v
		Describe("between "+c.n+" and", func() {
			Run := func(name string, marshal, unmarshal func(enc.Data) (interface{}, error)) {
				Context("value "+name, func() {
					It("should be successful", func() {
						d := v()
						d.SetKey()
						v, err := marshal(d)
						fmt.Fprintf(GinkgoWriter, "\nencoded: %v", d)
						Expect(err).To(BeNil())
						d.GetKey()
						vp, err := unmarshal(d)
						fmt.Fprintf(GinkgoWriter, "\ndecoded: %v", d)
						Expect(err).To(BeNil())
						if v != nil {
							Expect(vp).To(Equal(v))
						} else {
							Expect(vp).To(BeNil())
						}
					})
				})
			}
			Describe("Bool", func() {
				for _, c := range []struct {
					n string
					v Bool
				}{
					{"false", false},
					{"true", true},
				} {
					v := c.v
					var vp Bool
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("Int", func() {
				for _, c := range []struct {
					n string
					v Int
				}{
					{"0", 0},
					{"1", 1},
					{"-1", -1},
					{"1000000", 1000000},
					{"-1000000", -1000000},
					{"MaxInt64", math.MaxInt64},
					{"MinInt64", math.MinInt64},
				} {
					v := c.v
					var vp Int
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("Uint", func() {
				for _, c := range []struct {
					n string
					v Uint
				}{
					{"0", 0},
					{"1", 1},
					{"1000000", 1000000},
					{"MaxUint64", math.MaxUint64},
				} {
					v := c.v
					var vp Uint
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("Float", func() {
				for _, c := range []struct {
					n string
					v Float
				}{
					{"0", 0},
					{"1", 1},
					{"-1", -1},
					{"1000000", 1000000},
					{"-1000000", -1000000},
					{"1/1000000", 1.0 / 1000000},
					{"-1/1000000", -1.0 / 1000000},
					{"Pi", math.Pi},
					{"-Pi", -math.Pi},
					{"MaxFloat64", math.MaxFloat64},
					{"-MaxFloat64", -math.MaxFloat64},
					{"SmallestNonzeroFloat64", math.SmallestNonzeroFloat64},
					{"-SmallestNonzeroFloat64", -math.SmallestNonzeroFloat64},
				} {
					v := c.v
					var vp Float
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("String", func() {
				for _, c := range []struct {
					n string
					v String
				}{
					{"\"\"", ""},
					{"\"a\"", "a"},
					{"\"ab\"", "ab"},
				} {
					v := c.v
					var vp String
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("Bytes", func() {
				for _, c := range []struct {
					n string
					v Bytes
				}{
					{"{}", nil},
					{"{1}", Bytes{1}},
					{"{1, 2}", Bytes{1, 2}},
				} {
					v := c.v
					var vp Bytes
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("Enum", func() {
				for _, c := range []struct {
					n string
					v Enum
				}{
					{"0", 0},
					{"Option1", Enum_Option1},
					{"Option2", Enum_Option2},
				} {
					v := c.v
					var vp Enum
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("Union", func() {
				for _, c := range []struct {
					n string
					v Union
				}{
					{"nil", nil},
					{"{Option1: {1}}", &Union_As_Option1{
						Option1: []byte{1},
					}},
					{"{Option2: 2}", &Union_As_Option2{
						Option2: 2,
					}},
				} {
					v := c.v
					var vp Union
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, &v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("Bitfield", func() {
				for _, c := range []struct {
					n string
					v Bitfield
				}{
					{"{}", Bitfield{}},
					{"{Field1}", Bitfield{0x1}},
					{"{Field2}", Bitfield{0x2}},
					{"{Field1, Field2}", Bitfield{0x3}},
				} {
					v := c.v
					var vp Bitfield
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, &v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("Struct", func() {
				for _, c := range []struct {
					n string
					v *Struct
				}{
					{"{}", &Struct{}},
					{"{Field1: {1}}", &Struct{
						Field1: []byte{1},
					}},
					{"{Field2: 2}", &Struct{
						Field2: 2,
					}},
					{"{Field1: {1}, Field2: 2}", &Struct{
						Field1: []byte{1}, Field2: 2,
					}},
				} {
					v := c.v
					vp := &Struct{}
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, vp)
					})
				}
			})
			Describe("Array", func() {
				for _, c := range []struct {
					n string
					v Array
				}{
					{"{}", Array{}},
					{"{1}", Array{1}},
					{"{2}", Array{2}},
					{"{1, 2}", Array{1, 2}},
				} {
					v := c.v
					var vp Array
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("Map", func() {
				for _, c := range []struct {
					n string
					v Map
				}{
					{"{}", Map{}},
					{"{1: 1}", Map{1: 1}},
					{"{2: 2}", Map{2: 2}},
					{"{1: 1, 2: 2}", Map{1: 1, 2: 2}},
				} {
					v := c.v
					var vp Map
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("Set", func() {
				for _, c := range []struct {
					n string
					v Set
				}{
					{"{}", Set{}},
					{"{1}", Set{1: struct{}{}}},
					{"{2}", Set{2: struct{}{}}},
					{"{1, 2}", Set{1: struct{}{}, 2: struct{}{}}},
				} {
					v := c.v
					var vp Set
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("SelfRef", func() {
				for _, c := range []struct {
					n string
					v *SelfRef
				}{
					{"{}", &SelfRef{}},
					{"{{}}", &SelfRef{&SelfRef{}}},
					{"{{{}}}", &SelfRef{&SelfRef{&SelfRef{}}}},
				} {
					v := c.v
					vp := &SelfRef{}
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, vp)
					})
				}
			})
			Describe("MulRef", func() {
				for _, c := range []struct {
					n string
					v *MutRef1
				}{
					{"{}", &MutRef1{}},
					{"{{}}", &MutRef1{&MutRef2{}}},
					{"{{{}}}", &MutRef1{&MutRef2{&MutRef1{}}}},
				} {
					v := c.v
					vp := &MutRef1{}
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, vp)
					})
				}
			})
			Describe("NestedUnion", func() {
				for _, c := range []struct {
					n string
					v NestedUnion
				}{
					{"nil", nil},
					{"{Bool: true}", &NestedUnion_As_Bool{true}},
					{"{Int: 1}", &NestedUnion_As_Int{1}},
					{"{Uint: 1}", &NestedUnion_As_Uint{1}},
					{"{Float: 1}", &NestedUnion_As_Float{1}},
					{"{String: \"1\"}", &NestedUnion_As_String{"1"}},
					{"{Bytes: {1}}", &NestedUnion_As_Bytes{[]byte{1}}},
					{"{Enum: Option1}", &NestedUnion_As_Enum{
						Enum_Option1,
					}},
					{"{Union: {Option1: {1}}}", &NestedUnion_As_Union{
						&Union_As_Option1{Option1: []byte{1}},
					}},
					{"{Bitfield: {Field1, Field2}}", &NestedUnion_As_Bitfield{
						&Bitfield{0x3},
					}},
					{"{Struct: {Field1: {1}, Field2: 2}}", &NestedUnion_As_Struct{
						&Struct{Field1: []byte{1}, Field2: 2},
					}},
					{"{Array: {1, 2}", &NestedUnion_As_Array{
						Array{1, 2},
					}},
					{"{Map: {1:1, 2:2}", &NestedUnion_As_Map{
						Map{1: 1, 2: 2},
					}},
					{"{Set: {1, 2}", &NestedUnion_As_Set{
						Set{1: struct{}{}, 2: struct{}{}},
					}},
				} {
					v := c.v
					var vp NestedUnion
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, &v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("NestedStruct", func() {
				for _, c := range []struct {
					n string
					v *NestedStruct
				}{
					{"{}", &NestedStruct{}},
					{"{Bool: true}", &NestedStruct{Bool: true}},
					{"{Int: 1}", &NestedStruct{Int: 1}},
					{"{Uint: 1}", &NestedStruct{Uint: 1}},
					{"{Float: 1}", &NestedStruct{Float: 1}},
					{"{String: \"1\"}", &NestedStruct{String: "1"}},
					{"{Bytes: {1}}", &NestedStruct{Bytes: []byte{1}}},
					{"{Enum: Option1}", &NestedStruct{
						Enum: Enum_Option1,
					}},
					{"{Union: {Option1: {1}}}", &NestedStruct{
						Union: &Union_As_Option1{Option1: []byte{1}},
					}},
					{"{Bitfield: {Field1, Field2}}", &NestedStruct{
						Bitfield: &Bitfield{0x3},
					}},
					{"{Struct: {Field1: {1}, Field2: 2}}", &NestedStruct{
						Struct: &Struct{Field1: []byte{1}, Field2: 2},
					}},
					{"{Array: {1, 2}}}", &NestedStruct{
						Array: Array{1, 2},
					}},
					{"{Map: {1:1, 2:2}}", &NestedStruct{
						Map: Map{1: 1, 2: 2},
					}},
					{"{Set: {1, 2}}", &NestedStruct{
						Set: Set{1: struct{}{}, 2: struct{}{}},
					}},
					{"{Bool: true, " +
						"Int: 1, " +
						"Uint: 1, " +
						"Float: 1, " +
						"String: \"1\", " +
						"Bytes: {1}, " +
						"Enum: Option1, " +
						"Bitfield: {Field1, Field2}, " +
						"Union: {Option1: {1}}, " +
						"Struct: {Field1: {1}, Field2: 2}, " +
						"Array: {1, 2}}, " +
						"Map: {1:1, 2:2}, " +
						"Set: {1, 2}}", &NestedStruct{
						Bool:     true,
						Int:      1,
						Uint:     1,
						Float:    1,
						String:   "1",
						Bytes:    []byte{1},
						Enum:     Enum_Option1,
						Union:    &Union_As_Option1{Option1: []byte{1}},
						Bitfield: &Bitfield{0x3},
						Struct:   &Struct{Field1: []byte{1}, Field2: 2},
						Array:    Array{1, 2},
						Map:      Map{1: 1, 2: 2},
						Set:      Set{1: struct{}{}, 2: struct{}{}},
					}},
				} {
					v := c.v
					vp := &NestedStruct{}
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, vp)
					})
				}
			})
			Describe("NestedArray", func() {
				for _, c := range []struct {
					n string
					v NestedArray
				}{
					{"nil", nil},
					{"{Bool: {true}}", &NestedArray_As_Bool{[]bool{true}}},
					{"{Int: {1}}", &NestedArray_As_Int{[]int64{1}}},
					{"{Uint: {1}}", &NestedArray_As_Uint{[]uint64{1}}},
					{"{Float: {1}}", &NestedArray_As_Float{[]float64{1}}},
					{"{String: {\"1\"}}", &NestedArray_As_String{[]string{"1"}}},
					{"{Bytes: {{1}}}", &NestedArray_As_Bytes{[][]byte{{1}}}},
					{"{Enum: Option1}", &NestedArray_As_Enum{
						[]Enum{Enum_Option1},
					}},
					{"{Union: {Option1: {1}}}", &NestedArray_As_Union{
						[]Union{&Union_As_Option1{Option1: []byte{1}}},
					}},
					{"{Bitfield: {Field1, Field2}}", &NestedArray_As_Bitfield{
						[]*Bitfield{{0x3}},
					}},
					{"{Struct: {Field1: {1}, Field2: 2}}", &NestedArray_As_Struct{
						[]*Struct{{Field1: []byte{1}, Field2: 2}},
					}},
					{"{Array: {1, 2}", &NestedArray_As_Array{
						[]Array{{1, 2}},
					}},
					{"{Map: {1:1, 2:2}", &NestedArray_As_Map{
						[]Map{{1: 1, 2: 2}},
					}},
					{"{Set: {1, 2}", &NestedArray_As_Set{
						[]Set{{1: struct{}{}, 2: struct{}{}}},
					}},
				} {
					v := c.v
					var vp NestedArray
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, &v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
			Describe("NestedMap", func() {
				for _, c := range []struct {
					n string
					v NestedMap
				}{
					{"nil", nil},
					{"{Bool: {1: true}}", &NestedMap_As_Bool{map[Type]bool{1: true}}},
					{"{Int: {1: 1}}", &NestedMap_As_Int{map[Type]int64{1: 1}}},
					{"{Uint: {1: 1}}", &NestedMap_As_Uint{map[Type]uint64{1: 1}}},
					{"{Float: {1: 1}}", &NestedMap_As_Float{map[Type]float64{1: 1}}},
					{"{String: {1: \"1\"}}", &NestedMap_As_String{map[Type]string{1: "1"}}},
					{"{Bytes: {1: {1}}}", &NestedMap_As_Bytes{map[Type][]byte{1: {1}}}},
					{"{Enum: Option1}", &NestedMap_As_Enum{
						map[Type]Enum{0: Enum_Option1},
					}},
					{"{Union: {Option1: {1}}}", &NestedMap_As_Union{
						map[Type]Union{0: &Union_As_Option1{Option1: []byte{1}}},
					}},
					{"{Bitfield: {Field1, Field2}}", &NestedMap_As_Bitfield{
						map[Type]*Bitfield{0: {0x3}},
					}},
					{"{Struct: {Field1: {1}, Field2: 2}}", &NestedMap_As_Struct{
						map[Type]*Struct{0: {Field1: []byte{1}, Field2: 2}},
					}},
					{"{Array: {1, 2}", &NestedMap_As_Array{
						map[Type]Array{0: {1, 2}},
					}},
					{"{Map: {1:1, 2:2}", &NestedMap_As_Map{
						map[Type]Map{0: {1: 1, 2: 2}},
					}},
					{"{Set: {1, 2}", &NestedMap_As_Set{
						map[Type]Set{0: {1: struct{}{}, 2: struct{}{}}},
					}},
				} {
					v := c.v
					var vp NestedMap
					Run(c.n, func(d enc.Data) (interface{}, error) {
						return &v, m.Marshal(d, &v)
					}, func(d enc.Data) (interface{}, error) {
						return &vp, m.Unmarshal(d, &vp)
					})
				}
			})
		})
	}
})
