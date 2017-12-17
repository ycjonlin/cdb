package schema

import (
	"errors"
)

var (
	ErrZeroTypeTag            = errors.New("zero type tag")
	ErrZeroTypeName           = errors.New("zero type name")
	ErrDuplicatedTypeTag      = errors.New("duplicated type tag")
	ErrDuplicatedTypeName     = errors.New("duplicated type name")
	ErrCircularTypeDefinition = errors.New("circular type definition")
	ErrNonPrimitiveKeyType    = errors.New("non-primitive key type")
)

// Kind ...
type Kind int

const (
	EnumKind     Kind = 1
	BitfieldKind Kind = 2
	UnionKind    Kind = 3
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
type Tag uint32

// Name ...
type Name string

// ReferenceType ...
type ReferenceType struct {
	isType

	Tag  Tag
	Name Name
	Type Type

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

	Ref      *ReferenceType
	IsAtomic bool
	Kind     Kind
	Fields   []*ReferenceType

	tagMap  map[Tag]*ReferenceType
	nameMap map[Name]*ReferenceType
}

// ContainerType ...
type ContainerType struct {
	isType

	Ref      *ReferenceType
	IsAtomic bool
	Kind     Kind
	Key      *ReferenceType
	Elem     *ReferenceType
}

// NewReferenceType ...
func NewReferenceType(tag Tag, name Name, super *ReferenceType) *ReferenceType {
	return &ReferenceType{
		Tag:   tag,
		Name:  name,
		Super: super,
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
	if !ok || r.Tag != rp.Tag {
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
func NewCompositeType(r *ReferenceType, kind Kind) *CompositeType {
	c := &CompositeType{
		Ref:  r,
		Kind: kind,

		tagMap:  map[Tag]*ReferenceType{},
		nameMap: map[Name]*ReferenceType{},
	}
	if r != nil {
		r.SetType(c)
	}
	return c
}

// Match ...
func (c *CompositeType) Match(t Type) bool {
	cp, ok := t.(*CompositeType)
	if !ok || c.IsAtomic != cp.IsAtomic || c.Kind != cp.Kind {
		return false
	}
	for _, f := range c.Fields {
		fp := cp.GetByTag(f.Tag)
		if fp == nil || !f.Type.Match(fp.Type) {
			return false
		}
	}
	return true
}

// Put ...
func (c *CompositeType) Put(tag Tag, name Name) (*ReferenceType, error) {
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
	f := NewReferenceType(tag, name, c.Ref)
	c.Fields = append(c.Fields, f)
	c.tagMap[tag] = f
	c.nameMap[name] = f
	return f, nil
}

// GetByTag ...
func (c *CompositeType) GetByTag(tag Tag) *ReferenceType {
	return c.tagMap[tag]
}

// GetByName ...
func (c *CompositeType) GetByName(name Name) *ReferenceType {
	return c.nameMap[name]
}

// NewContainerType ...
func NewContainerType(r *ReferenceType, kind Kind) *ContainerType {
	c := &ContainerType{
		Ref:  r,
		Kind: kind,
		Key:  NewReferenceType(0, "Key", r),
		Elem: NewReferenceType(0, "Elem", r),
	}
	if r != nil {
		r.SetType(c)
	}
	return c
}

// Match ...
func (c *ContainerType) Match(t Type) bool {
	cp, ok := t.(*ContainerType)
	if !ok || c.IsAtomic != cp.IsAtomic || c.Kind != cp.Kind {
		return false
	}
	if !c.Key.Match(cp.Key) {
		return false
	}
	if !c.Elem.Match(cp.Elem) {
		return false
	}
	return true
}
