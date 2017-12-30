package schema

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Type", func() {
	Describe("ReferenceType", func() {
		Describe("NewReferenceType", func() {
			It("should success", func() {

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
		})
	})
	Describe("PrimitiveType", func() {
		Describe("NewPrimitiveType", func() {
			It("should sucess", func() {

				p := NewPrimitiveType("name")
				Expect(p.Name).To(Equal("name"))
			})
		})
	})
	Describe("CompositeType", func() {
		Describe("NewCompositeType", func() {
			It("should sucess", func() {

				r := NewReferenceType(-1, 0, "", nil, nil)
				c := NewCompositeType(r, EnumKind, nil)
				Expect(c.Ref).To(Equal(r))
				Expect(c.Kind).To(Equal(EnumKind))
				Expect(r.Type).To(Equal(c))
				Expect(r.Base).To(Equal(c))
			})
		})
		Describe(".Put", func() {
			It("should sucess", func() {

				c := NewCompositeType(nil, EnumKind, nil)
				f, err := c.Put(1, "a", nil)
				Expect(f.Tag).To(Equal(1))
				Expect(f.Name).To(Equal("a"))
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
		})
		Describe(".GetByIndex", func() {
			It("should sucess", func() {

				c := NewCompositeType(nil, EnumKind, nil)
				f, err := c.Put(1, "a", nil)
				Expect(err).To(BeNil())
				fp := c.GetByIndex(0)
				Expect(fp).To(Equal(f))
				fp = c.GetByIndex(1)
				Expect(fp).To(BeNil())
			})
		})
		Describe(".GetByTag", func() {
			It("should sucess", func() {

				c := NewCompositeType(nil, EnumKind, nil)
				f, err := c.Put(1, "a", nil)
				Expect(err).To(BeNil())
				fp := c.GetByTag(1)
				Expect(fp).To(Equal(f))
				fp = c.GetByTag(2)
				Expect(fp).To(BeNil())
			})
		})
		Describe(".GetByName", func() {
			It("should sucess", func() {

				c := NewCompositeType(nil, EnumKind, nil)
				f, err := c.Put(1, "a", nil)
				Expect(err).To(BeNil())
				fp := c.GetByName("a")
				Expect(fp).To(Equal(f))
				fp = c.GetByName("b")
				Expect(fp).To(BeNil())
			})
		})
		Describe(".GetByRType", func() {
			It("should sucess", func() {
				var v interface{}
				t1 := reflect.TypeOf(&v).Elem()
				t2 := reflect.TypeOf(&v)

				c := NewCompositeType(nil, EnumKind, nil)
				f, err := c.Put(1, "a", &v)
				Expect(err).To(BeNil())
				fp := c.GetByRType(t1)
				Expect(fp).To(Equal(f))
				fp = c.GetByRType(t2)
				Expect(fp).To(BeNil())
			})
		})
	})
	Describe("ContainerType", func() {
		Describe("NewContainerType", func() {
			It("should sucess", func() {

				r := NewReferenceType(-1, 0, "", nil, nil)
				c := NewContainerType(r, MapKind)
				Expect(c.Ref).To(Equal(r))
				Expect(c.Kind).To(Equal(MapKind))
				Expect(c.Key.Tag).To(Equal(0))
				Expect(c.Key.Name).To(Equal("Key"))
				Expect(c.Elem.Tag).To(Equal(0))
				Expect(c.Elem.Name).To(Equal("Elem"))
				Expect(r.Type).To(Equal(c))
				Expect(r.Base).To(Equal(c))
			})
		})
	})
})
