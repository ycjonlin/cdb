package schema

import (
	"errors"
	"reflect"
)

// Err ...
var (
	ErrZeroTypeTag            = errors.New("zero type tag")
	ErrZeroTypeName           = errors.New("zero type name")
	ErrDuplicatedTypeTag      = errors.New("duplicated type tag")
	ErrDuplicatedTypeName     = errors.New("duplicated type name")
	ErrCircularTypeDefinition = errors.New("circular type definition")
	ErrNonPrimitiveKeyType    = errors.New("non-primitive key type")
	ErrUnionNotImplemented    = errors.New("union not implemented")
)

// RType ...
type RType reflect.Type

// Kind ...
type Kind int

// Kind ...
const (
	EnumKind     Kind = 1
	UnionKind    Kind = 2
	BitfieldKind Kind = 3
	StructKind   Kind = 4

	ArrayKind Kind = 5
	MapKind   Kind = 6
	SetKind   Kind = 7
)

// Type ...
type Type interface {
	isType()
}

func (*ReferenceType) isType() {}
func (*PrimitiveType) isType() {}
func (*CompositeType) isType() {}
func (*ContainerType) isType() {}

// ReferenceType ...
type ReferenceType struct {
	Index int
	Tag   int
	Name  string
	Type  Type
	RType RType

	Base  Type
	Super *ReferenceType
}

// PrimitiveType ...
type PrimitiveType struct {
	Name string
}

// CompositeType ...
type CompositeType struct {
	Kind   Kind
	RType  RType
	Fields []*ReferenceType

	Ref *ReferenceType

	tagMap   map[int]*ReferenceType
	nameMap  map[string]*ReferenceType
	rtypeMap map[RType]*ReferenceType
}

// ContainerType ...
type ContainerType struct {
	Kind Kind
	Key  *ReferenceType
	Elem *ReferenceType

	Ref *ReferenceType
}

// NewReferenceType ...
func NewReferenceType(index int, tag int, name string, rtype RType, super *ReferenceType) *ReferenceType {
	return &ReferenceType{
		Index: index,
		Tag:   tag,
		Name:  name,
		Super: super,
		RType: rtype,
	}
}

// SetType ...
func (r *ReferenceType) SetType(typ Type) error {
	base := typ
	if r, ok := typ.(*ReferenceType); ok {
		if r.Base == nil {
			return ErrCircularTypeDefinition
		}
		base = r.Base
	}
	r.Type = typ
	r.Base = base
	return nil
}

// NewPrimitiveType ...
func NewPrimitiveType(name string) *PrimitiveType {
	return &PrimitiveType{
		Name: name,
	}
}

// NewCompositeType ...
func NewCompositeType(r *ReferenceType, kind Kind, v interface{}) *CompositeType {
	var rtype RType
	if kind == UnionKind && v != nil {
		rtype = reflect.TypeOf(v).Elem()
	}
	c := &CompositeType{
		Ref:   r,
		Kind:  kind,
		RType: rtype,

		tagMap:   map[int]*ReferenceType{},
		nameMap:  map[string]*ReferenceType{},
		rtypeMap: map[RType]*ReferenceType{},
	}
	if r != nil {
		r.SetType(c)
	}
	return c
}

// Put ...
func (c *CompositeType) Put(tag int, name string, v interface{}) (*ReferenceType, error) {
	if tag == 0 {
		return nil, ErrZeroTypeTag
	}
	if name == "" {
		return nil, ErrZeroTypeName
	}
	if _, ok := c.tagMap[tag]; ok {
		return nil, ErrDuplicatedTypeTag
	}
	if _, ok := c.nameMap[name]; ok {
		return nil, ErrDuplicatedTypeName
	}
	index := len(c.Fields)
	rtype := reflect.TypeOf(v)
	if rtype != nil {
		// interfaces could only be passed in as their pointers.
		if rtype.Kind() == reflect.Ptr &&
			rtype.Elem().Kind() == reflect.Interface {
			rtype = rtype.Elem()
		}
		if c.Kind == UnionKind && c.RType != nil {
			if !rtype.Implements(c.RType) {
				return nil, ErrUnionNotImplemented
			}
		}
	}
	f := NewReferenceType(index, tag, name, rtype, c.Ref)
	c.Fields = append(c.Fields, f)
	c.tagMap[tag] = f
	c.nameMap[name] = f
	if rtype != nil {
		c.rtypeMap[rtype] = f
	}
	return f, nil
}

// GetByIndex ...
func (c *CompositeType) GetByIndex(index int) *ReferenceType {
	if index < 0 || index >= len(c.Fields) {
		return nil
	}
	return c.Fields[index]
}

// GetByTag ...
func (c *CompositeType) GetByTag(tag int) *ReferenceType {
	return c.tagMap[tag]
}

// GetByName ...
func (c *CompositeType) GetByName(name string) *ReferenceType {
	return c.nameMap[name]
}

// GetByRType ...
func (c *CompositeType) GetByRType(rtype RType) *ReferenceType {
	return c.rtypeMap[rtype]
}

// NewContainerType ...
func NewContainerType(r *ReferenceType, kind Kind) *ContainerType {
	c := &ContainerType{
		Ref:  r,
		Kind: kind,
		Key:  NewReferenceType(0, 0, "Key", nil, r),
		Elem: NewReferenceType(0, 0, "Elem", nil, r),
	}
	if r != nil {
		r.SetType(c)
	}
	return c
}
