package generator

import (
	sch "github.com/ycjonlin/cdb/pkg/schema"
)

type marshallerDef interface {
	putMarshallerDef()
	putMarshallerRef()

	putMarshalPrimitiveType(t *sch.PrimitiveType, v string)
	putMarshalTagBegin(f *sch.ReferenceType)
	putMarshalTagEnd(f *sch.ReferenceType)
	putMarshalZeroTag()
	putMarshalZeroValue()
	putMarshalLen()

	putUnmarshalPrimitiveType(t *sch.PrimitiveType, v string)
	putUnmarshalTag(t *sch.CompositeType)
	putUnmarshalLen(t *sch.ContainerType)
}

type marshaller interface {
	putSchema(s *sch.Schema)
}

type marshallerImpl struct {
	*writer
	typer typer
	def   marshallerDef
}

func (c *marshallerImpl) putSchema(s *sch.Schema) {
	// Package
	c.putString("package ")
	c.putString(s.Name)
	c.putLine("")
	c.putLine("import enc \"../pkg/encoding\"")
	c.putLine("")
	c.putLine("type ")
	c.def.putMarshallerRef()
	c.putString(" ")
	c.def.putMarshallerDef()
	c.putLine("")
	// Types
	for _, r := range s.Types.Fields {
		c.putType(r)
	}
}

func (c *marshallerImpl) putType(r *sch.ReferenceType) {
	c.putMarshalFuncBegin(r)
	c.putMarshalType(r)
	c.putMarshalFuncEnd(r)

	/*c.putUnmarshalFuncBegin(r)
	c.putUnmarshalType(r)
	c.putUnmarshalFuncEnd(r)*/

	switch t := r.Type.(type) {
	case *sch.ReferenceType:
		// do nothing
	case *sch.PrimitiveType:
		// do nothing
	case *sch.CompositeType:
		c.putCompositeType(t)
	case *sch.ContainerType:
		c.putContainerType(t)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putCompositeType(t *sch.CompositeType) {
	switch t.Kind {
	case sch.EnumKind:
		// do nothing
	case sch.BitfieldKind:
		// do nothing
	case sch.UnionKind:
		for _, f := range t.Fields {
			c.putSubType(f)
		}
	case sch.StructKind:
		for _, f := range t.Fields {
			c.putSubType(f)
		}
	default:
		panic("")
	}
}

func (c *marshallerImpl) putContainerType(t *sch.ContainerType) {
	switch t.Kind {
	case sch.ArrayKind:
		c.putSubType(t.Elem)
	case sch.MapKind:
		c.putSubType(t.Key)
		c.putSubType(t.Elem)
	case sch.SetKind:
		c.putSubType(t.Key)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putSubType(r *sch.ReferenceType) {
	switch r.Type.(type) {
	case *sch.ReferenceType:
		// do nothing
	case *sch.PrimitiveType:
		// do nothing
	case *sch.CompositeType:
		c.putType(r)
	case *sch.ContainerType:
		c.putType(r)
	default:
		panic("")
	}
}

// Marshalling

func (c *marshallerImpl) putMarshalType(t *sch.ReferenceType) {
	switch t := t.Type.(type) {
	case *sch.ReferenceType:
		c.putMarshalReferenceType(t, "v")
	case *sch.PrimitiveType:
		c.def.putMarshalPrimitiveType(t, "v")
	case *sch.CompositeType:
		c.putMarshalCompositeType(t)
	case *sch.ContainerType:
		c.putMarshalContainerType(t)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putMarshalSubType(r *sch.ReferenceType, v string) {
	switch t := r.Type.(type) {
	case *sch.ReferenceType:
		c.putMarshalReferenceType(t, v)
	case *sch.PrimitiveType:
		c.def.putMarshalPrimitiveType(t, v)
	case *sch.CompositeType:
		c.putMarshalReferenceType(r, v)
	case *sch.ContainerType:
		c.putMarshalReferenceType(r, v)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putMarshalKeyType(r *sch.ReferenceType, v string) {
	switch t := r.Base.(type) {
	case *sch.PrimitiveType:
		c.def.putMarshalPrimitiveType(t, v)
	case *sch.CompositeType:
		switch t.Kind {
		case sch.EnumKind, sch.BitfieldKind:
			c.putMarshalSubType(r, v)
		default:
			panic("")
		}
	default:
		panic("")
	}
}

func (c *marshallerImpl) putMarshalReferenceType(t *sch.ReferenceType, v string) {
	c.putLine("err = m.Marshal")
	c.putCompoundName(t)
	c.putString("(")
	c.typer.putConvertType(t, v)
	c.putString(")")
	c.putMarshalFuncReturnErr()
}

func (c *marshallerImpl) putMarshalCompositeType(t *sch.CompositeType) {
	switch t.Kind {
	case sch.EnumKind:
		c.putLine("switch v {")
		for _, f := range t.Fields {
			c.putLine("case ")
			c.putCompoundName(f)
			c.putString(":")
			{
				c.pushIndent()
				c.def.putMarshalTagBegin(f)
				c.def.putMarshalZeroValue()
				c.def.putMarshalTagEnd(f)
				c.putLine("return nil")
				c.popIndent()
			}
		}
		c.putMarshalFuncSwitchDefault(t)
		c.putLine("}")
	case sch.BitfieldKind:
		for _, f := range t.Fields {
			c.putLine("if v&")
			c.putCompoundName(f)
			c.putString(" != 0 {")
			{
				c.pushIndent()
				c.def.putMarshalTagBegin(f)
				c.def.putMarshalZeroValue()
				c.def.putMarshalTagEnd(f)
				c.popIndent()
			}
			c.putLine("}")
		}
		c.def.putMarshalZeroTag()
		c.putLine("return nil")
	case sch.UnionKind:
		c.putLine("switch u := v.(type) {")
		for _, f := range t.Fields {
			c.putLine("case *")
			c.typer.putUnionOptionTypeRef(f)
			c.putString(":")
			{
				c.pushIndent()
				c.def.putMarshalTagBegin(f)
				c.putLine("f := u.")
				c.putName(f.Name)
				c.putMarshalSubType(f, "f")
				c.def.putMarshalTagEnd(f)
				c.putLine("return nil")
				c.popIndent()
			}
		}
		c.putMarshalFuncSwitchDefault(t)
		c.putLine("}")
	case sch.StructKind:
		for _, f := range t.Fields {
			c.putLine("if f := v.")
			c.putName(f.Name)
			c.putString("; f != ")
			c.typer.putZero(f)
			c.putString(" {")
			{
				c.pushIndent()
				c.def.putMarshalTagBegin(f)
				c.putMarshalSubType(f, "f")
				c.def.putMarshalTagEnd(f)
				c.popIndent()
			}
			c.putLine("}")
		}
		c.def.putMarshalZeroTag()
		c.putLine("return nil")
	default:
		panic("")
	}
}

func (c *marshallerImpl) putMarshalContainerType(t *sch.ContainerType) {
	switch t.Kind {
	case sch.ArrayKind:
		c.def.putMarshalLen()
		c.putLine("for i, n := 0, len(v); i < n; i++ {")
		{
			c.pushIndent()
			c.putLine("e := v[i]")
			c.putMarshalSubType(t.Elem, "e")
			c.popIndent()
		}
		c.putLine("}")
	case sch.MapKind:
		c.def.putMarshalLen()
		c.putLine("for k, e := range v {")
		{
			c.pushIndent()
			c.putMarshalKeyType(t.Key, "k")
			c.putMarshalSubType(t.Elem, "e")
			c.popIndent()
		}
		c.putLine("}")
	case sch.SetKind:
		c.def.putMarshalLen()
		c.putLine("for k, _ := range v {")
		{
			c.pushIndent()
			c.putMarshalKeyType(t.Key, "k")
			c.popIndent()
		}
		c.putLine("}")
	default:
		panic("")
	}
}

func (c *marshallerImpl) putMarshalFuncBegin(t *sch.ReferenceType) {
	c.putLine("func (m *")
	c.def.putMarshallerRef()
	c.putString(") Marshal")
	c.putCompoundName(t)
	c.putString("(v ")
	c.typer.putTypeRef(t)
	c.putString(") (err error) {")
	c.pushIndent()
}

func (c *marshallerImpl) putMarshalFuncEnd(t *sch.ReferenceType) {
	if _, ok := t.Type.(*sch.CompositeType); !ok {
		c.putLine("return nil")
	}
	c.popIndent()
	c.putLine("}")
	c.putLine("")
}

func (c *marshallerImpl) putMarshalFuncReturnErr() {
	c.putLine("if err != nil {")
	{
		c.pushIndent()
		c.putLine("return err")
		c.popIndent()
	}
	c.putLine("}")
}

func (c *marshallerImpl) putMarshalFuncSwitchDefault(t *sch.CompositeType) {
	switch t.Kind {
	case sch.EnumKind:
		c.putLine("case 0:")
	case sch.UnionKind:
		c.putLine("case nil:")
	default:
		panic("")
	}
	{
		c.pushIndent()
		c.def.putMarshalZeroTag()
		c.putLine("return nil")
		c.popIndent()
	}
	c.putLine("default:")
	{
		c.pushIndent()
		c.putLine("return enc.ErrUnknownTag")
		c.popIndent()
	}
}

// Unmarshalling

func (c *marshallerImpl) putUnmarshalType(t *sch.ReferenceType) {
	switch t := t.Type.(type) {
	case *sch.ReferenceType:
		c.putUnmarshalReferenceType(t, "v")
	case *sch.PrimitiveType:
		c.def.putUnmarshalPrimitiveType(t, "v")
	case *sch.CompositeType:
		c.putUnmarshalCompositeType(t)
	case *sch.ContainerType:
		c.putUnmarshalContainerType(t)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putUnmarshalSubType(r *sch.ReferenceType, v string) {
	switch t := r.Type.(type) {
	case *sch.ReferenceType:
		c.putUnmarshalReferenceType(t, v)
	case *sch.PrimitiveType:
		c.def.putUnmarshalPrimitiveType(t, v)
	case *sch.CompositeType:
		c.putUnmarshalReferenceType(r, v)
	case *sch.ContainerType:
		c.putUnmarshalReferenceType(r, v)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putUnmarshalReferenceType(t *sch.ReferenceType, v string) {
	c.putLine(v)
	c.putString(", err := m.Unmarshal")
	c.putCompoundName(t)
	c.putString("()")
}

func (c *marshallerImpl) putUnmarshalCompositeType(t *sch.CompositeType) {
	switch t.Kind {
	case sch.EnumKind:
		c.def.putUnmarshalTag(t)
		c.putUnmarshalFuncReturnErr(t)
		c.putLine("switch t {")
		for _, f := range t.Fields {
			c.putUnmarshalFuncSwitchCase(f)
			{
				c.pushIndent()
				c.putLine("return ")
				c.putCompoundName(f)
				c.putString(", nil")
				c.popIndent()
			}
		}
		c.putUnmarshalFuncSwitchDefault(t)
		c.putLine("}")
	case sch.BitfieldKind:
		c.putLine("var v ")
		c.putCompoundName(t.Ref)
		c.putLine("for {")
		{
			c.pushIndent()
			c.def.putUnmarshalTag(t)
			c.putUnmarshalFuncReturnErr(t)
			c.putLine("switch t {")
			for _, f := range t.Fields {
				c.putUnmarshalFuncSwitchCase(f)
				{
					c.pushIndent()
					c.putLine("v |= ")
					c.putCompoundName(f)
					c.popIndent()
				}
			}
			c.putUnmarshalFuncSwitchDefault(t)
			c.putLine("}")
			c.popIndent()
		}
		c.putLine("}")
	case sch.UnionKind:
		c.def.putUnmarshalTag(t)
		c.putUnmarshalFuncReturnErr(t)
		c.putLine("switch t {")
		for _, f := range t.Fields {
			c.putUnmarshalFuncSwitchCase(f)
			{
				c.pushIndent()
				c.putUnmarshalSubType(f, "u")
				c.putUnmarshalFuncReturnErr(t)
				c.putLine("return &")
				c.typer.putUnionOptionTypeRef(f)
				c.putString("{u}, nil")
				c.popIndent()
			}
		}
		c.putUnmarshalFuncSwitchDefault(t)
		c.putLine("}")
	case sch.StructKind:
		c.putLine("v := &")
		c.putCompoundName(t.Ref)
		c.putString("{}")
		c.putLine("for {")
		{
			c.pushIndent()
			c.def.putUnmarshalTag(t)
			c.putUnmarshalFuncReturnErr(t)
			c.putLine("switch t {")
			for _, f := range t.Fields {
				c.putUnmarshalFuncSwitchCase(f)
				{
					c.pushIndent()
					c.putUnmarshalSubType(f, "f")
					c.putUnmarshalFuncReturnErr(t)
					c.putLine("v.")
					c.putName(f.Name)
					c.putString(" = f")
					c.popIndent()
				}
			}
			c.putUnmarshalFuncSwitchDefault(t)
			c.putLine("}")
			c.popIndent()
		}
		c.putLine("}")
	default:
		panic("")
	}
}

func (c *marshallerImpl) putUnmarshalContainerType(t *sch.ContainerType) {
	c.def.putUnmarshalLen(t)
	c.putUnmarshalFuncReturnErr(t)
	switch t.Kind {
	case sch.ArrayKind:
		c.putLine("v := make([]")
		c.typer.putSubTypeRef(t.Elem)
		c.putString(", n)")
		c.putLine("for i := 0; i < n; i++ {")
		{
			c.pushIndent()
			c.putUnmarshalSubType(t.Elem, "e")
			c.putUnmarshalFuncReturnErr(t)
			c.putLine("v[i] = e")
			c.popIndent()
		}
		c.putLine("}")
		c.putLine("return v, nil")
	case sch.MapKind:
		c.putLine("v := ")
		c.typer.putTypeRef(t.Ref)
		c.putString("{}")
		c.putLine("for i := 0; i < n; i++ {")
		{
			c.pushIndent()
			c.putUnmarshalSubType(t.Key, "k")
			c.putUnmarshalFuncReturnErr(t)
			c.putUnmarshalSubType(t.Elem, "e")
			c.putUnmarshalFuncReturnErr(t)
			c.putLine("v[k] = e")
			c.popIndent()
		}
		c.putLine("}")
		c.putLine("return v, nil")
	case sch.SetKind:
		c.putLine("v := ")
		c.typer.putTypeRef(t.Ref)
		c.putString("{}")
		c.putLine("for i := 0; i < n; i++ {")
		{
			c.pushIndent()
			c.putUnmarshalSubType(t.Key, "k")
			c.putUnmarshalFuncReturnErr(t)
			c.putLine("v[k] = struct{}{}")
			c.popIndent()
		}
		c.putLine("}")
		c.putLine("return v, nil")
	default:
		panic("")
	}
}

func (c *marshallerImpl) putUnmarshalFuncBegin(t *sch.ReferenceType) {
	c.putLine("func (m *")
	c.def.putMarshallerRef()
	c.putString(") Unmarshal")
	c.putCompoundName(t)
	c.putString("() (")
	c.typer.putTypeRef(t)
	c.putString(", error) {")
	c.pushIndent()
}

func (c *marshallerImpl) putUnmarshalFuncEnd(t *sch.ReferenceType) {
	switch t.Type.(type) {
	case *sch.ReferenceType, *sch.PrimitiveType:
		c.putLine("return (")
		c.typer.putTypeRef(t)
		c.putString(")(v), err")
	}
	c.popIndent()
	c.putLine("}")
	c.putLine("")
}

func (c *marshallerImpl) putUnmarshalFuncReturnErr(t sch.Type) {
	c.putLine("if err != nil {")
	{
		c.pushIndent()
		c.putLine("return ")
		c.typer.putZero(t)
		c.putString(", err")
		c.popIndent()
	}
	c.putLine("}")
}

func (c *marshallerImpl) putUnmarshalFuncSwitchCase(f *sch.ReferenceType) {
	c.putLine("case ")
	c.putTag(f.Tag)
	c.putString(":")
}

func (c *marshallerImpl) putUnmarshalFuncSwitchDefault(t *sch.CompositeType) {
	c.putLine("case 0:")
	{
		c.pushIndent()
		c.putLine("return ")
		if t.Kind == sch.EnumKind || t.Kind == sch.UnionKind {
			c.typer.putZero(t)
		} else {
			c.putString("v")
		}
		c.putString(", nil")
		c.popIndent()
	}
	c.putLine("default:")
	{
		c.pushIndent()
		c.putLine("return ")
		c.typer.putZero(t)
		c.putString(", enc.ErrUnknownTag")
		c.popIndent()
	}
}
