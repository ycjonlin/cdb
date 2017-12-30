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
	doType()
	Match(t Type) bool
}

type isType struct{}

func (isType) doType() {}

// Tag ...
type Tag int

// Name ...
type Name string

// ReferenceType ...
type ReferenceType struct {
	isType

	Index int
	Tag   Tag
	Name  Name
	Type  Type
	RType RType

	Base  Type
	Super *ReferenceType
}

// PrimitiveType ...
type PrimitiveType struct {
	isType

	Name Name
}

// CompositeType ...
type CompositeType struct {
	isType

	Ref   *ReferenceType
	Kind  Kind
	RType RType

	Fields []*ReferenceType

	tagMap   map[Tag]*ReferenceType
	nameMap  map[Name]*ReferenceType
	rtypeMap map[RType]*ReferenceType
}

// ContainerType ...
type ContainerType struct {
	isType

	Ref  *ReferenceType
	Kind Kind
	Key  *ReferenceType
	Elem *ReferenceType
}

// NewReferenceType ...
func NewReferenceType(index int, tag Tag, name Name, rtype RType, super *ReferenceType) *ReferenceType {
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

// Match ...
func (r *ReferenceType) Match(t Type) bool {
	rp, ok := t.(*ReferenceType)
	if !ok || r != rp {
		return false
	}
	return true
}

// NewPrimitiveType ...
func NewPrimitiveType(name Name) *PrimitiveType {
	return &PrimitiveType{
		Name: name,
	}
}

// Match ...
func (p *PrimitiveType) Match(t Type) bool {
	pp, ok := t.(*PrimitiveType)
	if !ok || p != pp {
		return false
	}
	return true
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

		tagMap:   map[Tag]*ReferenceType{},
		nameMap:  map[Name]*ReferenceType{},
		rtypeMap: map[RType]*ReferenceType{},
	}
	if r != nil {
		r.SetType(c)
	}
	return c
}

// Match ...
func (c *CompositeType) Match(t Type) bool {
	cp, ok := t.(*CompositeType)
	if !ok || c.Kind != cp.Kind || len(c.Fields) != len(cp.Fields) {
		return false
	}
	for _, f := range c.Fields {
		fp := cp.GetByTag(f.Tag)
		if fp == nil || f.Name != fp.Name {
			return false
		}
		if c.Kind == UnionKind || c.Kind == StructKind {
			if !f.Type.Match(fp.Type) {
				return false
			}
		}
	}
	return true
}

// Put ...
func (c *CompositeType) Put(tag Tag, name Name, v interface{}) (*ReferenceType, error) {
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

// MustPut ...
func (c *CompositeType) MustPut(tag Tag, name Name, v interface{}) *ReferenceType {
	r, err := c.Put(tag, name, v)
	if err != nil {
		panic(err)
	}
	return r
}

// GetByIndex ...
func (c *CompositeType) GetByIndex(index int) *ReferenceType {
	if index < 0 || index >= len(c.Fields) {
		return nil
	}
	return c.Fields[index]
}

// GetByTag ...
func (c *CompositeType) GetByTag(tag Tag) *ReferenceType {
	return c.tagMap[tag]
}

// GetByName ...
func (c *CompositeType) GetByName(name Name) *ReferenceType {
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

// Match ...
func (c *ContainerType) Match(t Type) bool {
	cp, ok := t.(*ContainerType)
	if !ok || c.Kind != cp.Kind {
		return false
	}
	if c.Kind == MapKind || c.Kind == SetKind {
		if !c.Key.Type.Match(cp.Key.Type) {
			return false
		}
	}
	if c.Kind == ArrayKind || c.Kind == MapKind {
		if !c.Elem.Type.Match(cp.Elem.Type) {
			return false
		}
	}
	return true
}
