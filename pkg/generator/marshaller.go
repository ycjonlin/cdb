package generator

import (
	"strings"

	sch "github.com/ycjonlin/cdb/pkg/schema"
)

type marshaller interface {
	putSchema(s *sch.Schema)
}

type marshallerImpl struct {
	*writer
	typer
}

func (c *marshallerImpl) putSchema(s *sch.Schema) {
	// Package
	c.putString("package ")
	c.putString(s.Name)
	c.putLine("")
	c.pushIndent("import (")
	c.putLine("enc \"github.com/ycjonlin/cdb/pkg/encoding\"")
	c.putLine("sch \"github.com/ycjonlin/cdb/pkg/schema\"")
	c.popIndent(")")
	c.putLine("")
	c.pushIndent("var ")
	c.putString(strings.Title(s.Name))
	c.putString("Marshaller = func() *enc.Marshaller {")
	c.putLine("m := enc.NewMarshaller()")
	// Types
	for _, f := range s.Types.Fields {
		c.pushIndent("{")
		c.putLine("r := m.MustPut(")
		c.putTag(f.Tag)
		c.putString(", \"")
		c.putName(f.Name)
		c.putString("\", ")
		c.putZero(f)
		c.putString(")")
		c.putType(f)
		c.popIndent("}")
	}
	c.putLine("err := m.Check()")
	c.pushIndent("if err != nil {")
	{
		c.putLine("panic(err)")
	}
	c.popIndent("}")
	c.putLine("return m")
	c.popIndent("}()")
	c.putLine("")
}

func (c *marshallerImpl) putZero(r *sch.ReferenceType) {
	switch t := r.Base.(type) {
	case *sch.PrimitiveType:
		c.putPrimitiveTypeZero(r, t)
	case *sch.CompositeType:
		c.putCompositeTypeZero(r, t)
	case *sch.ContainerType:
		c.putContainerTypeZero(r, t)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putPrimitiveTypeZero(r *sch.ReferenceType, t *sch.PrimitiveType) {
	switch t {
	case sch.BoolType:
		c.putCompoundName(r)
		c.putString("(false)")
	case sch.IntType:
		c.putCompoundName(r)
		c.putString("(0)")
	case sch.UintType:
		c.putCompoundName(r)
		c.putString("(0)")
	case sch.FloatType:
		c.putCompoundName(r)
		c.putString("(0)")
	case sch.StringType:
		c.putCompoundName(r)
		c.putString("(\"\")")
	case sch.BytesType:
		c.putString("make(")
		c.putCompoundName(r)
		c.putString(", 0)")
	default:
		panic("")
	}
}

func (c *marshallerImpl) putCompositeTypeZero(r *sch.ReferenceType, t *sch.CompositeType) {
	switch t.Kind {
	case sch.EnumKind:
		c.putCompoundName(r)
		c.putString("(0)")
	case sch.UnionKind:
		c.putString("new(")
		c.putCompoundName(r)
		c.putString(")")
	case sch.BitfieldKind:
		c.putString("new(")
		c.putCompoundName(r)
		c.putString(")")
	case sch.StructKind:
		c.putString("new(")
		c.putCompoundName(r)
		c.putString(")")
	default:
		panic("")
	}
}

func (c *marshallerImpl) putContainerTypeZero(r *sch.ReferenceType, t *sch.ContainerType) {
	switch t.Kind {
	case sch.ArrayKind:
		c.putString("make(")
		c.putCompoundName(r)
		c.putString(", 0)")
	case sch.MapKind:
		c.putString("make(")
		c.putCompoundName(r)
		c.putString(")")
	case sch.SetKind:
		c.putString("make(")
		c.putCompoundName(r)
		c.putString(")")
	default:
		panic("")
	}
}

func (c *marshallerImpl) putType(r *sch.ReferenceType) {
	switch t := r.Type.(type) {
	case *sch.ReferenceType:
		c.putReferenceType(r, t)
	case *sch.PrimitiveType:
		c.putPrimitiveType(r, t)
	case *sch.CompositeType:
		c.putCompositeType(r, t)
	case *sch.ContainerType:
		c.putContainerType(r, t)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putReferenceType(r *sch.ReferenceType, t *sch.ReferenceType) {
	c.putLine("r.SetType(m.GetByName(\"")
	c.putName(t.Name)
	c.putString("\"))")
}

func (c *marshallerImpl) putPrimitiveType(r *sch.ReferenceType, t *sch.PrimitiveType) {
	c.putLine("r.SetType(sch.")
	switch t {
	case sch.BoolType:
		c.putString("BoolType)")
	case sch.IntType:
		c.putString("IntType)")
	case sch.UintType:
		c.putString("UintType)")
	case sch.FloatType:
		c.putString("FloatType)")
	case sch.StringType:
		c.putString("StringType)")
	case sch.BytesType:
		c.putString("BytesType)")
	default:
		panic("")
	}
}

func (c *marshallerImpl) putCompositeType(r *sch.ReferenceType, t *sch.CompositeType) {
	switch t.Kind {
	case sch.EnumKind:
		c.putLine("c := sch.NewCompositeType(r, sch.EnumKind, nil)")
		for _, f := range t.Fields {
			c.putLine("c.MustPut(")
			c.putTag(f.Tag)
			c.putString(", \"")
			c.putName(f.Name)
			c.putString("\", nil)")
		}
	case sch.UnionKind:
		c.putLine("c := sch.NewCompositeType(r, sch.UnionKind, new(")
		c.putCompoundName(r)
		c.putString("))")
		for _, f := range t.Fields {
			c.pushIndent("{")
			c.putLine("r := c.MustPut(")
			c.putTag(f.Tag)
			c.putString(", \"")
			c.putName(f.Name)
			c.putString("\", new(")
			c.putUnionOptionTypeRef(f)
			c.putString("))")
			c.putType(f)
			c.popIndent("}")
		}
	case sch.BitfieldKind:
		c.putLine("c := sch.NewCompositeType(r, sch.BitfieldKind, nil)")
		for _, f := range t.Fields {
			c.putLine("c.MustPut(")
			c.putTag(f.Tag)
			c.putString(", \"")
			c.putName(f.Name)
			c.putString("\", nil)")
		}
	case sch.StructKind:
		c.putLine("c := sch.NewCompositeType(r, sch.StructKind, nil)")
		for _, f := range t.Fields {
			c.pushIndent("{")
			c.putLine("r := c.MustPut(")
			c.putTag(f.Tag)
			c.putString(", \"")
			c.putName(f.Name)
			c.putString("\", nil)")
			c.putType(f)
			c.popIndent("}")
		}
	default:
		panic("")
	}
}

func (c *marshallerImpl) putContainerType(r *sch.ReferenceType, t *sch.ContainerType) {
	switch t.Kind {
	case sch.ArrayKind:
		c.putLine("c := sch.NewContainerType(r, sch.ArrayKind)")
		c.pushIndent("{")
		c.putLine("r := c.Elem")
		c.putType(t.Elem)
		c.popIndent("}")
	case sch.MapKind:
		c.putLine("c := sch.NewContainerType(r, sch.MapKind)")
		c.pushIndent("{")
		c.putLine("r := c.Key")
		c.putType(t.Key)
		c.popIndent("}")
		c.pushIndent("{")
		c.putLine("r := c.Elem")
		c.putType(t.Elem)
		c.popIndent("}")
	case sch.SetKind:
		c.putLine("c := sch.NewContainerType(r, sch.SetKind)")
		c.pushIndent("{")
		c.putLine("r := c.Key")
		c.putType(t.Key)
		c.popIndent("}")
	default:
		panic("")
	}
}
