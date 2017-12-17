package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReferenceType_SetType(t *testing.T) {
	p := NewPrimitiveType("")
	r1 := NewReferenceType(0, "", nil)
	r2 := NewReferenceType(0, "", nil)
	err := r2.SetType(r1)
	assert.Equal(t, err, ErrCircularTypeDefinition)
	err = r1.SetType(p)
	assert.Nil(t, err)
	assert.Equal(t, r1.Type, p)
	assert.Equal(t, r1.Base, p)
	err = r2.SetType(r1)
	assert.Nil(t, err)
	assert.Equal(t, r2.Type, r1)
	assert.Equal(t, r2.Base, p)
}

func Test_ReferenceType_Match(t *testing.T) {
	r1 := NewReferenceType(0, "", nil)
	r2 := NewReferenceType(0, "", nil)
	assert.True(t, r1.Match(r1))
	assert.False(t, r1.Match(r2))
}

func Test_PrimitiveType_Match(t *testing.T) {
	p1 := NewPrimitiveType("")
	p2 := NewPrimitiveType("")
	assert.True(t, p1.Match(p1))
	assert.False(t, p1.Match(p2))
}

func Test_NewCompositeType(t *testing.T) {
	r := NewReferenceType(0, "", nil)
	c := NewCompositeType(r, EnumKind)
	assert.Equal(t, c.Ref, r)
	assert.Equal(t, c.Kind, EnumKind)
	assert.Equal(t, r.Type, c)
	assert.Equal(t, r.Base, c)
}

func Test_CompositeType_Match(t *testing.T) {
	p := NewPrimitiveType("")

	e := NewCompositeType(nil, EnumKind)
	e.Put(1, "a")

	s := NewCompositeType(nil, StructKind)
	f, _ := s.Put(1, "a")
	f.SetType(p)

	assert.False(t, e.Match(s))

	{
		ep := NewCompositeType(nil, EnumKind)
		ep.Put(1, "a")
		assert.True(t, e.Match(ep))
	}
	{
		ep := NewCompositeType(nil, EnumKind)
		assert.False(t, e.Match(ep))
	}
	{
		ep := NewCompositeType(nil, EnumKind)
		ep.Put(2, "a")
		assert.False(t, e.Match(ep))
	}
	{
		ep := NewCompositeType(nil, EnumKind)
		ep.Put(1, "b")
		assert.False(t, e.Match(ep))
	}
	{
		ep := NewCompositeType(nil, EnumKind)
		ep.Put(1, "a")
		ep.Put(2, "b")
		assert.False(t, e.Match(ep))
	}
	{
		sp := NewCompositeType(nil, StructKind)
		f, _ := sp.Put(1, "a")
		f.SetType(p)
		assert.True(t, s.Match(sp))
	}
	{
		sp := NewCompositeType(nil, StructKind)
		f, _ := sp.Put(1, "a")
		f.SetType(e)
		assert.False(t, s.Match(sp))
	}
}

func Test_CompositeType_Put(t *testing.T) {
	c := NewCompositeType(nil, EnumKind)
	f, err := c.Put(1, "a")
	assert.Equal(t, f.Tag, Tag(1))
	assert.Equal(t, f.Name, Name("a"))
	assert.Nil(t, err)
	_, err = c.Put(0, "a")
	assert.Equal(t, err, ErrZeroTypeTag)
	_, err = c.Put(1, "")
	assert.Equal(t, err, ErrZeroTypeName)
	_, err = c.Put(1, "b")
	assert.Equal(t, err, ErrDuplicatedTypeTag)
	_, err = c.Put(2, "a")
	assert.Equal(t, err, ErrDuplicatedTypeName)
}

func Test_CompositeType_GetByTag(t *testing.T) {
	c := NewCompositeType(nil, EnumKind)
	f, err := c.Put(1, "a")
	assert.Nil(t, err)
	fp := c.GetByTag(1)
	assert.Equal(t, fp, f)
	fp = c.GetByTag(2)
	assert.Nil(t, fp)
}

func Test_CompositeType_GetByName(t *testing.T) {
	c := NewCompositeType(nil, EnumKind)
	f, err := c.Put(1, "a")
	assert.Nil(t, err)
	fp := c.GetByName("a")
	assert.Equal(t, fp, f)
	fp = c.GetByName("b")
	assert.Nil(t, fp)
}

func Test_NewContainerType(t *testing.T) {
	r := NewReferenceType(0, "", nil)
	c := NewContainerType(r, MapKind)
	assert.Equal(t, c.Ref, r)
	assert.Equal(t, c.Kind, MapKind)
	assert.Equal(t, c.Key.Tag, Tag(0))
	assert.Equal(t, c.Key.Name, Name("Key"))
	assert.Equal(t, c.Elem.Tag, Tag(0))
	assert.Equal(t, c.Elem.Name, Name("Elem"))
	assert.Equal(t, r.Type, c)
	assert.Equal(t, r.Base, c)
}

func Test_ContainerType_Match(t *testing.T) {
	p1 := NewPrimitiveType("")
	p2 := NewPrimitiveType("")

	m := NewContainerType(nil, MapKind)
	m.Key.SetType(p1)
	m.Elem.SetType(p1)

	a := NewContainerType(nil, ArrayKind)
	a.Elem.SetType(p1)
	assert.False(t, m.Match(a))

	{
		mp := NewContainerType(nil, MapKind)
		mp.Key.SetType(p1)
		mp.Elem.SetType(p1)
		assert.True(t, m.Match(mp))
	}
	{
		mp := NewContainerType(nil, MapKind)
		mp.Key.SetType(p2)
		mp.Elem.SetType(p1)
		assert.False(t, m.Match(mp))
	}
	{
		mp := NewContainerType(nil, MapKind)
		mp.Key.SetType(p1)
		mp.Elem.SetType(p2)
		assert.False(t, m.Match(mp))
	}
}
