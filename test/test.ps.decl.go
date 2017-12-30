package test

type Bool bool

type Int int64

type Uint uint64

type Float float64

type String string

type Bytes []byte

type Type Int

type Enum int

const (
  Enum_Option1 Enum = 1
  Enum_Option2 Enum = 2
)

type Union interface {
  isUnion()
}

type Union_As_Option1 struct{ Option1 []byte }
type Union_As_Option2 struct{ Option2 Type }

func (*Union_As_Option1) isUnion() {}
func (*Union_As_Option2) isUnion() {}

type Bitfield [1]byte

func (c *Bitfield) GetField1() bool { return c[0]&0x1 != 0 }
func (c *Bitfield) GetField2() bool { return c[0]&0x2 != 0 }

func (c *Bitfield) SetField1() { c[0] |= 0x1 }
func (c *Bitfield) SetField2() { c[0] |= 0x2 }

func (c *Bitfield) DelField1() { c[0] &^= 0x1 }
func (c *Bitfield) DelField2() { c[0] &^= 0x2 }

type Struct struct {
  Field1 []byte
  Field2 Type
}

type Array []Type

type Map map[Type]Type

type Set map[Type]struct{}

type SelfRef struct {
  Field *SelfRef
}

type MutRef1 struct {
  Field *MutRef2
}

type MutRef2 struct {
  Field *MutRef1
}

type NestedUnion interface {
  isNestedUnion()
}

type NestedUnion_As_Bool struct{ Bool bool }
type NestedUnion_As_Int struct{ Int int64 }
type NestedUnion_As_Uint struct{ Uint uint64 }
type NestedUnion_As_Float struct{ Float float64 }
type NestedUnion_As_String struct{ String string }
type NestedUnion_As_Bytes struct{ Bytes []byte }
type NestedUnion_As_Enum struct{ Enum Enum }
type NestedUnion_As_Union struct{ Union Union }
type NestedUnion_As_Bitfield struct{ Bitfield *Bitfield }
type NestedUnion_As_Struct struct{ Struct *Struct }
type NestedUnion_As_Array struct{ Array Array }
type NestedUnion_As_Map struct{ Map Map }
type NestedUnion_As_Set struct{ Set Set }

func (*NestedUnion_As_Bool) isNestedUnion() {}
func (*NestedUnion_As_Int) isNestedUnion() {}
func (*NestedUnion_As_Uint) isNestedUnion() {}
func (*NestedUnion_As_Float) isNestedUnion() {}
func (*NestedUnion_As_String) isNestedUnion() {}
func (*NestedUnion_As_Bytes) isNestedUnion() {}
func (*NestedUnion_As_Enum) isNestedUnion() {}
func (*NestedUnion_As_Union) isNestedUnion() {}
func (*NestedUnion_As_Bitfield) isNestedUnion() {}
func (*NestedUnion_As_Struct) isNestedUnion() {}
func (*NestedUnion_As_Array) isNestedUnion() {}
func (*NestedUnion_As_Map) isNestedUnion() {}
func (*NestedUnion_As_Set) isNestedUnion() {}

type NestedStruct struct {
  Bool bool
  Int int64
  Uint uint64
  Float float64
  String string
  Bytes []byte
  Enum Enum
  Union Union
  Bitfield *Bitfield
  Struct *Struct
  Array Array
  Map Map
  Set Set
}

type NestedArray interface {
  isNestedArray()
}

type NestedArray_As_Bool struct{ Bool []bool }
type NestedArray_As_Int struct{ Int []int64 }
type NestedArray_As_Uint struct{ Uint []uint64 }
type NestedArray_As_Float struct{ Float []float64 }
type NestedArray_As_String struct{ String []string }
type NestedArray_As_Bytes struct{ Bytes [][]byte }
type NestedArray_As_Enum struct{ Enum []Enum }
type NestedArray_As_Union struct{ Union []Union }
type NestedArray_As_Bitfield struct{ Bitfield []*Bitfield }
type NestedArray_As_Struct struct{ Struct []*Struct }
type NestedArray_As_Array struct{ Array []Array }
type NestedArray_As_Map struct{ Map []Map }
type NestedArray_As_Set struct{ Set []Set }

func (*NestedArray_As_Bool) isNestedArray() {}
func (*NestedArray_As_Int) isNestedArray() {}
func (*NestedArray_As_Uint) isNestedArray() {}
func (*NestedArray_As_Float) isNestedArray() {}
func (*NestedArray_As_String) isNestedArray() {}
func (*NestedArray_As_Bytes) isNestedArray() {}
func (*NestedArray_As_Enum) isNestedArray() {}
func (*NestedArray_As_Union) isNestedArray() {}
func (*NestedArray_As_Bitfield) isNestedArray() {}
func (*NestedArray_As_Struct) isNestedArray() {}
func (*NestedArray_As_Array) isNestedArray() {}
func (*NestedArray_As_Map) isNestedArray() {}
func (*NestedArray_As_Set) isNestedArray() {}

type NestedMap interface {
  isNestedMap()
}

type NestedMap_As_Bool struct{ Bool map[Type]bool }
type NestedMap_As_Int struct{ Int map[Type]int64 }
type NestedMap_As_Uint struct{ Uint map[Type]uint64 }
type NestedMap_As_Float struct{ Float map[Type]float64 }
type NestedMap_As_String struct{ String map[Type]string }
type NestedMap_As_Bytes struct{ Bytes map[Type][]byte }
type NestedMap_As_Enum struct{ Enum map[Type]Enum }
type NestedMap_As_Union struct{ Union map[Type]Union }
type NestedMap_As_Bitfield struct{ Bitfield map[Type]*Bitfield }
type NestedMap_As_Struct struct{ Struct map[Type]*Struct }
type NestedMap_As_Array struct{ Array map[Type]Array }
type NestedMap_As_Map struct{ Map map[Type]Map }
type NestedMap_As_Set struct{ Set map[Type]Set }

func (*NestedMap_As_Bool) isNestedMap() {}
func (*NestedMap_As_Int) isNestedMap() {}
func (*NestedMap_As_Uint) isNestedMap() {}
func (*NestedMap_As_Float) isNestedMap() {}
func (*NestedMap_As_String) isNestedMap() {}
func (*NestedMap_As_Bytes) isNestedMap() {}
func (*NestedMap_As_Enum) isNestedMap() {}
func (*NestedMap_As_Union) isNestedMap() {}
func (*NestedMap_As_Bitfield) isNestedMap() {}
func (*NestedMap_As_Struct) isNestedMap() {}
func (*NestedMap_As_Array) isNestedMap() {}
func (*NestedMap_As_Map) isNestedMap() {}
func (*NestedMap_As_Set) isNestedMap() {}
