package schema

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Type", func() {
	It("reference type set type", func() {

		p := NewPrimitiveType("")
		r1 := NewReferenceType(-1, 0, "", nil, nil)
		r2 := NewReferenceType(-1, 0, "", nil, nil)

		err := r2.SetType(r1)
		Expect(err).To(Equal(ErrCircularTypeDefinition))

		err = r1.SetType(p)
		Expect(err).To(BeNil())
		Expect(r1.Type).To(Equal(p))
		Expect(r1.Base).To(Equal(p))

		err = r2.SetType(r1)
		Expect(err).To(BeNil())
		Expect(r2.Type).To(Equal(r1))
		Expect(r2.Base).To(Equal(p))
	})
	It("reference type match", func() {

		r1 := NewReferenceType(-1, 0, "", nil, nil)
		r2 := NewReferenceType(-1, 0, "", nil, nil)
		Expect(r1.Match(r1)).To(BeTrue())
		Expect(r1.Match(r2)).To(BeFalse())
	})
	It("primitive type match", func() {

		p1 := NewPrimitiveType("")
		p2 := NewPrimitiveType("")
		Expect(p1.Match(p1)).To(BeTrue())
		Expect(p1.Match(p2)).To(BeFalse())
	})
	It("new composite type", func() {

		r := NewReferenceType(-1, 0, "", nil, nil)
		c := NewCompositeType(r, EnumKind, nil)
		Expect(c.Ref).To(Equal(r))
		Expect(c.Kind).To(Equal(EnumKind))
		Expect(r.Type).To(Equal(c))
		Expect(r.Base).To(Equal(c))
	})
	It("composite type match", func() {

		p := NewPrimitiveType("")

		e := NewCompositeType(nil, EnumKind, nil)
		e.Put(1, "a", nil)

		s := NewCompositeType(nil, StructKind, nil)
		f, _ := s.Put(1, "a", nil)
		f.SetType(p)

		Expect(e.Match(s)).To(BeFalse())

		{
			ep := NewCompositeType(nil, EnumKind, nil)
			ep.Put(1, "a", nil)
			Expect(e.Match(ep)).To(BeTrue())
		}
		{
			ep := NewCompositeType(nil, EnumKind, nil)
			Expect(e.Match(ep)).To(BeFalse())
		}
		{
			ep := NewCompositeType(nil, EnumKind, nil)
			ep.Put(2, "a", nil)
			Expect(e.Match(ep)).To(BeFalse())
		}
		{
			ep := NewCompositeType(nil, EnumKind, nil)
			ep.Put(1, "b", nil)
			Expect(e.Match(ep)).To(BeFalse())
		}
		{
			ep := NewCompositeType(nil, EnumKind, nil)
			ep.Put(1, "a", nil)
			ep.Put(2, "b", nil)
			Expect(e.Match(ep)).To(BeFalse())
		}
		{
			sp := NewCompositeType(nil, StructKind, nil)
			f, _ := sp.Put(1, "a", nil)
			f.SetType(p)
			Expect(s.Match(sp)).To(BeTrue())
		}
		{
			sp := NewCompositeType(nil, StructKind, nil)
			f, _ := sp.Put(1, "a", nil)
			f.SetType(e)
			Expect(s.Match(sp)).To(BeFalse())
		}
	})
	It("composite type put", func() {

		c := NewCompositeType(nil, EnumKind, nil)
		f, err := c.Put(1, "a", nil)
		Expect(f.Tag).To(Equal(Tag(1)))
		Expect(f.Name).To(Equal(Name("a")))
		Expect(err).To(BeNil())
		_, err = c.Put(0, "a", nil)
		Expect(err).To(Equal(ErrZeroTypeTag))
		_, err = c.Put(1, "", nil)
		Expect(err).To(Equal(ErrZeroTypeName))
		_, err = c.Put(1, "b", nil)
		Expect(err).To(Equal(ErrDuplicatedTypeTag))
		_, err = c.Put(2, "a", nil)
		Expect(err).To(Equal(ErrDuplicatedTypeName))
	})
	It("composite type get by tag", func() {

		c := NewCompositeType(nil, EnumKind, nil)
		f, err := c.Put(1, "a", nil)
		Expect(err).To(BeNil())
		fp := c.GetByTag(1)
		Expect(fp).To(Equal(f))
		fp = c.GetByTag(2)
		Expect(fp).To(BeNil())
	})
	It("composite type get by name", func() {

		c := NewCompositeType(nil, EnumKind, nil)
		f, err := c.Put(1, "a", nil)
		Expect(err).To(BeNil())
		fp := c.GetByName("a")
		Expect(fp).To(Equal(f))
		fp = c.GetByName("b")
		Expect(fp).To(BeNil())
	})
	It("new container type", func() {

		r := NewReferenceType(-1, 0, "", nil, nil)
		c := NewContainerType(r, MapKind)
		Expect(c.Ref).To(Equal(r))
		Expect(c.Kind).To(Equal(MapKind))
		Expect(c.Key.Tag).To(Equal(Tag(0)))
		Expect(c.Key.Name).To(Equal(Name("Key")))
		Expect(c.Elem.Tag).To(Equal(Tag(0)))
		Expect(c.Elem.Name).To(Equal(Name("Elem")))
		Expect(r.Type).To(Equal(c))
		Expect(r.Base).To(Equal(c))
	})
	It("container type match", func() {

		p1 := NewPrimitiveType("")
		p2 := NewPrimitiveType("")

		m := NewContainerType(nil, MapKind)
		m.Key.SetType(p1)
		m.Elem.SetType(p1)

		a := NewContainerType(nil, ArrayKind)
		a.Elem.SetType(p1)
		Expect(m.Match(a)).To(BeFalse())

		{
			mp := NewContainerType(nil, MapKind)
			mp.Key.SetType(p1)
			mp.Elem.SetType(p1)
			Expect(m.Match(mp)).To(BeTrue())
		}
		{
			mp := NewContainerType(nil, MapKind)
			mp.Key.SetType(p2)
			mp.Elem.SetType(p1)
			Expect(m.Match(mp)).To(BeFalse())
		}
		{
			mp := NewContainerType(nil, MapKind)
			mp.Key.SetType(p1)
			mp.Elem.SetType(p2)
			Expect(m.Match(mp)).To(BeFalse())
		}
	})
})
