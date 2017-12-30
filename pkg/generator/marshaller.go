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
	c.putString("Marshaller = enc.NewMarshaller(sch.UnionDef{")
	// Types
	for _, f := range s.Types.Fields {
		c.putLine("{")
		c.putTag(f.Tag)
		c.putString(", \"")
		c.putName(f.Name)
		c.putString("\", ")
		c.putZero(f)
		c.putString(", ")
		c.putType(f)
		c.putString("},")
	}
	c.popIndent("})")
	c.putLine("")
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
	c.putString("sch.NameRef(\"")
	c.putName(t.Name)
	c.putString("\")")
}

func (c *marshallerImpl) putPrimitiveType(r *sch.ReferenceType, t *sch.PrimitiveType) {
	c.putString("&sch.BuiltIn{sch.")
	switch t {
	case sch.BoolType:
		c.putString("BoolType}")
	case sch.IntType:
		c.putString("IntType}")
	case sch.UintType:
		c.putString("UintType}")
	case sch.FloatType:
		c.putString("FloatType}")
	case sch.StringType:
		c.putString("StringType}")
	case sch.BytesType:
		c.putString("BytesType}")
	default:
		panic("")
	}
}

func (c *marshallerImpl) putCompositeType(r *sch.ReferenceType, t *sch.CompositeType) {
	switch t.Kind {
	case sch.EnumKind:
		c.pushIndentCont("sch.EnumDef{")
		for _, f := range t.Fields {
			c.putLine("{")
			c.putTag(f.Tag)
			c.putString(", \"")
			c.putName(f.Name)
			c.putString("\"},")
		}
		c.popIndent("}")
	case sch.UnionKind:
		c.pushIndentCont("sch.UnionDef{")
		for _, f := range t.Fields {
			c.putLine("{")
			c.putTag(f.Tag)
			c.putString(", \"")
			c.putName(f.Name)
			c.putString("\", new(")
			c.putUnionOptionTypeRef(f)
			c.putString("), ")
			c.putType(f)
			c.putString("},")
		}
		c.popIndent("}")
	case sch.BitfieldKind:
		c.pushIndentCont("sch.BitfieldDef{")
		for _, f := range t.Fields {
			c.putLine("{")
			c.putTag(f.Tag)
			c.putString(", \"")
			c.putName(f.Name)
			c.putString("\"},")
		}
		c.popIndent("}")
	case sch.StructKind:
		c.pushIndentCont("sch.StructDef{")
		for _, f := range t.Fields {
			c.putLine("{")
			c.putTag(f.Tag)
			c.putString(", \"")
			c.putName(f.Name)
			c.putString("\", nil, ")
			c.putType(f)
			c.putString("},")
		}
		c.popIndent("}")
	default:
		panic("")
	}
}

func (c *marshallerImpl) putContainerType(r *sch.ReferenceType, t *sch.ContainerType) {
	switch t.Kind {
	case sch.ArrayKind:
		c.putString("&sch.ArrayDef{")
		c.putType(t.Elem)
		c.putString("}")
	case sch.MapKind:
		c.putString("&sch.MapDef{")
		c.putType(t.Key)
		c.putString(", ")
		c.putType(t.Elem)
		c.putString("}")
	case sch.SetKind:
		c.putString("&sch.SetDef{")
		c.putType(t.Key)
		c.putString("}")
	default:
		panic("")
	}
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
