package schema

import (
	"fmt"
	"reflect"
)

// TypeDef ...
type TypeDef interface {
	isTypeDef()
}

func (NameRef) isTypeDef()     {}
func (TagRef) isTypeDef()      {}
func (*BuiltIn) isTypeDef()    {}
func (EnumDef) isTypeDef()     {}
func (UnionDef) isTypeDef()    {}
func (BitfieldDef) isTypeDef() {}
func (StructDef) isTypeDef()   {}
func (*ArrayDef) isTypeDef()   {}
func (*MapDef) isTypeDef()     {}
func (*SetDef) isTypeDef()     {}

// ConstDef ...
type ConstDef struct {
	Tag  int
	Name string
}

// FieldDef ...
type FieldDef struct {
	Tag   int
	Name  string
	Value interface{}
	Type  TypeDef
}

// NameRef ...
type NameRef string

// TagRef ...
type TagRef int

// BuiltIn ...
type BuiltIn struct{ Type Type }

// EnumDef ...
type EnumDef []ConstDef

// UnionDef ...
type UnionDef []FieldDef

// BitfieldDef ...
type BitfieldDef []ConstDef

// StructDef ...
type StructDef []FieldDef

// ArrayDef ...
type ArrayDef struct{ Elem TypeDef }

// MapDef ...
type MapDef struct{ Key, Elem TypeDef }

// SetDef ...
type SetDef struct{ Key TypeDef }

type assign struct {
	ref  *ReferenceType
	name string
}

type compiler struct {
	assigns []assign
}

type checker struct {
	checked map[*ReferenceType]bool
}

// Compile ...
func Compile(d UnionDef) (*CompositeType, error) {
	t := NewCompositeType(nil, UnionKind, nil)
	var c compiler
	for _, f := range d {
		r, err := t.Put(f.Tag, f.Name, f.Value)
		if err != nil {
			return nil, err
		}
		err = c.compile(f.Type, r)
		if err != nil {
			return nil, err
		}
	}
	for _, a := range c.assigns {
		r := t.GetByName(a.name)
		if r == nil {
			return nil, ErrIllegalName
		}
		err := a.ref.SetType(r)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

// Check ...
func Check(t *CompositeType) error {
	var c checker
	c.checked = map[*ReferenceType]bool{}
	for _, f := range t.Fields {
		rf := f.RType
		err := c.check(f, rf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *compiler) compile(d TypeDef, r *ReferenceType) error {
	switch d := d.(type) {
	case NameRef:
		c.assigns = append(c.assigns, assign{r, string(d)})
		return nil
	case *BuiltIn:
		r.SetType(d.Type)
		return nil
	case EnumDef:
		t := NewCompositeType(r, EnumKind, nil)
		return c.compileConsts(d, t)
	case UnionDef:
		t := NewCompositeType(r, UnionKind, nil)
		return c.compileFields(d, t)
	case BitfieldDef:
		t := NewCompositeType(r, BitfieldKind, nil)
		return c.compileConsts(d, t)
	case StructDef:
		t := NewCompositeType(r, StructKind, nil)
		return c.compileFields(d, t)
	case *ArrayDef:
		t := NewContainerType(r, ArrayKind)
		err := c.compile(d.Elem, t.Elem)
		if err != nil {
			return err
		}
		return nil
	case *MapDef:
		t := NewContainerType(r, MapKind)
		err := c.compile(d.Key, t.Key)
		if err != nil {
			return err
		}
		err = c.compile(d.Elem, t.Elem)
		if err != nil {
			return err
		}
		return nil
	case *SetDef:
		t := NewContainerType(r, SetKind)
		err := c.compile(d.Key, t.Key)
		if err != nil {
			return err
		}
		return nil
	default:
		panic("")
	}
}

func (c *compiler) compileConsts(d []ConstDef, t *CompositeType) error {
	for _, c := range d {
		_, err := t.Put(c.Tag, c.Name, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *compiler) compileFields(d []FieldDef, t *CompositeType) error {
	for _, f := range d {
		r, err := t.Put(f.Tag, f.Name, f.Value)
		if err != nil {
			return err
		}
		err = c.compile(f.Type, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *checker) check(r *ReferenceType, rt reflect.Type) error {
	if c.checked[r] {
		return nil
	}
	c.checked[r] = true
	if rt == nil {
		return fmt.Errorf("Type must not be nil")
	}
	switch t := r.Base.(type) {
	case *PrimitiveType:
		return c.checkPrimitiveType(t, rt)
	case *CompositeType:
		return c.checkCompositeType(t, rt)
	case *ContainerType:
		return c.checkContainerType(t, rt)
	default:
		return fmt.Errorf("unknown type %v", reflect.TypeOf(t))
	}
}

func (c *checker) checkPrimitiveType(t *PrimitiveType, rt reflect.Type) error {
	switch t {
	case BoolType:
		if rt.Kind() != reflect.Bool {
			return fmt.Errorf("Bool.Kind must be reflect.Bool, but got %v (%v)", rt.Kind(), rt)
		}
		return nil
	case IntType:
		if rt.Kind() != reflect.Int64 {
			return fmt.Errorf("Int.Kind must be reflect.Int64, but got %v (%v)", rt.Kind(), rt)
		}
		return nil
	case UintType:
		if rt.Kind() != reflect.Uint64 {
			return fmt.Errorf("Uint.Kind must be reflect.Uint64, but got %v (%v)", rt.Kind(), rt)
		}
		return nil
	case FloatType:
		if rt.Kind() != reflect.Float64 {
			return fmt.Errorf("Float.Kind must be reflect.Float64, but got %v (%v)", rt.Kind(), rt)
		}
		return nil
	case StringType:
		if rt.Kind() != reflect.String {
			return fmt.Errorf("String.Kind must be reflect.String, but got %v (%v)", rt.Kind(), rt)
		}
		return nil
	case BytesType:
		if rt.Kind() != reflect.Slice {
			return fmt.Errorf("Bytes.Kind must be reflect.Slice, but got %v (%v)", rt.Kind(), rt)
		}
		if rt.Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("Bytes.Enum.Kind must be reflect.Uint8, but got %v (%v)", rt.Elem().Kind(), rt.Elem())
		}
		return nil
	default:
		return fmt.Errorf("unknown primitive type")
	}
}

func (c *checker) checkCompositeType(t *CompositeType, rt reflect.Type) error {
	switch t.Kind {
	case EnumKind:
		return c.checkEnumType(t, rt)
	case UnionKind:
		return c.checkUnionType(t, rt)
	case BitfieldKind:
		return c.checkBitfieldType(t, rt)
	case StructKind:
		return c.checkStructType(t, rt)
	default:
		panic("unknown composite type")
	}
}

func (c *checker) checkEnumType(t *CompositeType, rt reflect.Type) error {
	if rt.Kind() != reflect.Int {
		return fmt.Errorf("Enum.Kind must be reflect.Int, but got %v (%v)", rt.Kind(), rt)
	}
	return nil
}

func (c *checker) checkUnionType(t *CompositeType, rt reflect.Type) error {
	if rt.Kind() != reflect.Interface {
		return fmt.Errorf("Union.Kind must be reflect.Interface, but got %v (%v)", rt.Kind(), rt)
	}
	for _, f := range t.Fields {
		rf := f.RType
		if !rf.Implements(rt) {
			return fmt.Errorf("Union_Option must be implement Union, but %v does not implement %v", rf, rt)
		}
		if rf.Kind() != reflect.Ptr {
			return fmt.Errorf("Struct.Kind must be reflect.Ptr, but got %v (%v)", rf.Kind(), rf)
		}
		if rf.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("Union_Option.Kind must be reflect.Struct, but got %v (%v)", rf.Elem().Kind(), rf.Elem())
		}
		if rf.Elem().NumField() != 1 {
			return fmt.Errorf("Struct.Elem.NumField must be 1, but got %d (%v)", rf.Elem().NumField(), rf.Elem())
		}
		rf = rf.Elem().Field(0).Type
		err := c.check(f, rf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *checker) checkBitfieldType(t *CompositeType, rt reflect.Type) error {
	if rt.Kind() != reflect.Ptr {
		return fmt.Errorf("Bitfield.Kind must be reflect.Ptr, but got %v (%v)", rt.Kind(), rt)
	}
	if rt.Elem().Kind() != reflect.Array {
		return fmt.Errorf("Bitfield.Elem.Kind must be reflect.Array, but got %v (%v)", rt.Elem().Kind(), rt.Elem())
	}
	if rt.Elem().Elem().Kind() != reflect.Uint8 {
		return fmt.Errorf("Bitfield.Elem.Kind.Elem must be reflect.Uint8, but got %v (%v)", rt.Elem().Elem().Kind(), rt.Elem().Elem())
	}
	if rt.Elem().Len() != (len(t.Fields)+7)/8 {
		return fmt.Errorf("Bitfield.Elem.Len must be (len(Fields)+7)/8 = %d, but got %d (%v)", (len(t.Fields)+7)/8, rt.Elem().Len(), rt.Elem())
	}
	return nil
}

func (c *checker) checkStructType(t *CompositeType, rt reflect.Type) error {
	if rt.Kind() != reflect.Ptr {
		return fmt.Errorf("Struct.Kind must be reflect.Ptr, but got %v (%v)", rt.Kind(), rt)
	}
	if rt.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Struct.Elem.Kind must be reflect.Struct, but got %v (%v)", rt.Elem().Kind(), rt.Elem())
	}
	if rt.Elem().NumField() != len(t.Fields) {
		return fmt.Errorf("Struct.Elem.NumField must be len(Fields) = %d, but got %d (%v)", len(t.Fields), rt.Elem().NumField(), rt.Elem())
	}
	for _, f := range t.Fields {
		rf := rt.Elem().Field(f.Index).Type
		err := c.check(f, rf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *checker) checkContainerType(t *ContainerType, rt reflect.Type) error {
	switch t.Kind {
	case ArrayKind:
		return c.checkArrayType(t, rt)
	case MapKind:
		return c.checkMapType(t, rt)
	case SetKind:
		return c.checkSetType(t, rt)
	default:
		panic("unknown composite type")
	}
}

func (c *checker) checkArrayType(t *ContainerType, rt reflect.Type) error {
	if rt.Kind() != reflect.Slice {
		return fmt.Errorf("Array.Kind must be reflect.Slice, but got %v (%v)", rt.Kind(), rt)
	}
	err := c.check(t.Elem, rt.Elem())
	if err != nil {
		return err
	}
	return nil
}

func (c *checker) checkMapType(t *ContainerType, rt reflect.Type) error {
	if rt.Kind() != reflect.Map {
		return fmt.Errorf("Map.Kind must be reflect.Map, but got %v (%v)", rt.Kind(), rt)
	}
	err := c.check(t.Key, rt.Key())
	if err != nil {
		return err
	}
	err = c.check(t.Elem, rt.Elem())
	if err != nil {
		return err
	}
	return nil
}

func (c *checker) checkSetType(t *ContainerType, rt reflect.Type) error {
	if rt.Kind() != reflect.Map {
		return fmt.Errorf("Set.Kind must be reflect.Map, but got %v (%v)", rt.Kind(), rt)
	}
	if rt.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Set.Elem.Kind must be reflect.Struct, but got %v (%v)", rt.Elem().Kind(), rt.Elem())
	}
	if rt.Elem().NumField() != 0 {
		return fmt.Errorf("Set.Elem.NumField must be 0, but got %d (%v)", rt.Elem().NumField(), rt.Elem())
	}
	err := c.check(t.Key, rt.Key())
	if err != nil {
		return err
	}
	return nil
}
