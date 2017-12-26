package generator

import (
	sch "github.com/ycjonlin/cdb/pkg/schema"
)

type marshaller interface {
	putSchema(s *sch.Schema)
}

type marshallerImpl struct {
	*writer
	typer typer
}

func (c *marshallerImpl) putSchema(s *sch.Schema) {
	// Package
	c.putString("package ")
	c.putString(s.Name)
	c.putLine("")
	c.putLine("import enc \"github.com/ycjonlin/cdb/pkg/encoding\"")
	c.putLine("")
	c.putLine("type Marshaller struct{}")
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

	c.putUnmarshalFuncBegin(r)
	c.putUnmarshalType(r)
	c.putUnmarshalFuncEnd(r)

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
		c.putMarshalReferenceType(t, "d", "v")
	case *sch.PrimitiveType:
		c.putMarshalPrimitiveType(t, "d", "v")
	case *sch.CompositeType:
		c.putMarshalCompositeType(t)
	case *sch.ContainerType:
		c.putMarshalContainerType(t)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putUnmarshalType(t *sch.ReferenceType) {
	switch t := t.Type.(type) {
	case *sch.ReferenceType:
		c.putUnmarshalReferenceType(t, "d", "v")
	case *sch.PrimitiveType:
		c.putUnmarshalPrimitiveType(t, "d", "v")
	case *sch.CompositeType:
		c.putUnmarshalCompositeType(t)
	case *sch.ContainerType:
		c.putUnmarshalContainerType(t)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putMarshalSubType(r *sch.ReferenceType, e, v string) {
	switch t := r.Type.(type) {
	case *sch.ReferenceType:
		c.putMarshalReferenceType(t, e, v)
	case *sch.PrimitiveType:
		c.putMarshalPrimitiveType(t, e, v)
	case *sch.CompositeType:
		c.putMarshalReferenceType(r, e, v)
	case *sch.ContainerType:
		c.putMarshalReferenceType(r, e, v)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putUnmarshalSubType(r *sch.ReferenceType, e, v string) {
	switch t := r.Type.(type) {
	case *sch.ReferenceType:
		c.putUnmarshalReferenceType(t, e, v)
	case *sch.PrimitiveType:
		c.putUnmarshalPrimitiveType(t, e, v)
	case *sch.CompositeType:
		c.putUnmarshalReferenceType(r, e, v)
	case *sch.ContainerType:
		c.putUnmarshalReferenceType(r, e, v)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putMarshalKeyType(r *sch.ReferenceType, v string) {
	c.putLine("dk := d.Data(d.SetKey())")
	c.putMarshalSubType(r, "dk", v)
}

func (c *marshallerImpl) putUnmarshalKeyType(r *sch.ReferenceType, v string) {
	c.putLine("dk := d.Data(d.GetKey())")
	c.putUnmarshalSubType(r, "dk", v)
}

func (c *marshallerImpl) putMarshalReferenceType(t *sch.ReferenceType, e, v string) {
	c.putLine("err = m.Marshal")
	c.putCompoundName(t)
	c.putString("(")
	c.putString(e)
	c.putString(", ")
	c.typer.putConvertType(t, v)
	c.putString(")")
	c.putMarshalFuncReturnErr()
}

func (c *marshallerImpl) putUnmarshalReferenceType(t *sch.ReferenceType, e, v string) {
	c.putLine(v)
	c.putString(", err := m.Unmarshal")
	c.putCompoundName(t)
	c.putString("(")
	c.putString(e)
	c.putString(")")
}

func (c *marshallerImpl) putMarshalPrimitiveType(t *sch.PrimitiveType, e, v string) {
	c.putLine(e)
	c.putString(".Value().Encode")
	switch t {
	case sch.BoolType:
		c.putString("Bool(bool(")
	case sch.IntType:
		c.putString("Varint(int64(")
	case sch.UintType:
		c.putString("Uvarint(uint64(")
	case sch.FloatType:
		c.putString("Float64(float64(")
	case sch.StringType:
		c.putString("NonsortingString(string(")
	case sch.BytesType:
		c.putString("NonsortingBytes([]byte(")
	default:
		panic("")
	}
	c.putString(v)
	c.putString("))")
}

func (c *marshallerImpl) putUnmarshalPrimitiveType(t *sch.PrimitiveType, e, v string) {
	c.putLine(v)
	c.putString(", err := ")
	c.putString(e)
	c.putString(".Value().Decode")
	switch t {
	case sch.BoolType:
		c.putString("Bool()")
	case sch.IntType:
		c.putString("Varint()")
	case sch.UintType:
		c.putString("Uvarint()")
	case sch.FloatType:
		c.putString("Float64()")
	case sch.StringType:
		c.putString("NonsortingString()")
	case sch.BytesType:
		c.putString("NonsortingBytes(nil)")
	default:
		panic("")
	}
}

func (c *marshallerImpl) putMarshalCompositeType(t *sch.CompositeType) {
	switch t.Kind {
	case sch.EnumKind:
		c.putMarshalEnumType(t)
	case sch.BitfieldKind:
		c.putMarshalBitfieldType(t)
	case sch.UnionKind:
		c.putMarshalUnionType(t)
	case sch.StructKind:
		c.putMarshalStructType(t)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putUnmarshalCompositeType(t *sch.CompositeType) {
	switch t.Kind {
	case sch.EnumKind:
		c.putUnmarshalEnumType(t)
	case sch.BitfieldKind:
		c.putUnmarshalBitfieldType(t)
	case sch.UnionKind:
		c.putUnmarshalUnionType(t)
	case sch.StructKind:
		c.putUnmarshalStructType(t)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putMarshalEnumType(t *sch.CompositeType) {
	c.putLine("switch v {")
	for _, f := range t.Fields {
		c.putLine("case ")
		c.putCompoundName(f)
		c.putString(":")
		{
			c.pushIndent()
			c.putMarshalTag(f, true)
			c.putMarshalEmpty()
			c.popIndent()
		}
	}
	c.putMarshalFuncSwitchDefault(t)
	c.putLine("}")
	c.putMarshalEmpty()
	c.putLine("return nil")
}

func (c *marshallerImpl) putUnmarshalEnumType(t *sch.CompositeType) {
	c.putLine("var v ")
	c.putCompoundName(t.Ref)
	c.putUnmarshalTag(t, "t", true)
	c.putLine("switch t {")
	for _, f := range t.Fields {
		c.putUnmarshalFuncSwitchCase(f)
		{
			c.pushIndent()
			c.putUnmarshalEmpty()
			c.putLine("v = ")
			c.putCompoundName(f)
			c.popIndent()
		}
	}
	c.putUnmarshalFuncSwitchDefault(t)
	c.putLine("}")
	c.putUnmarshalEmpty()
	c.putLine("return v, nil")
}

func (c *marshallerImpl) putMarshalBitfieldType(t *sch.CompositeType) {
	c.putLine("if v == nil {")
	{
		c.pushIndent()
		c.putLine("goto end")
		c.popIndent()
	}
	c.putLine("}")
	for _, f := range t.Fields {
		c.putLine("if set, del := v.Set&")
		c.typer.putStructFieldFlag(f)
		c.putString(" != 0,")
		{
			c.pushIndent()
			c.putLine("v.Del&")
			c.typer.putStructFieldFlag(f)
			c.putString(" != 0; set || del {")
			c.putMarshalTag(f, false)
			c.putMarshalEmpty()
			c.popIndent()
		}
		c.putLine("}")
	}
	{
		c.popIndent()
		c.putLine("end:")
		c.pushIndent()
	}
	c.putMarshalEnd()
	c.putMarshalEmpty()
	c.putLine("return nil")
}

func (c *marshallerImpl) putUnmarshalBitfieldType(t *sch.CompositeType) {
	c.putLine("v := &")
	c.putCompoundName(t.Ref)
	c.putString("{}")
	c.putLine("for d.Next() {")
	{
		c.pushIndent()
		c.putUnmarshalTag(t, "t", false)
		c.putLine("switch t {")
		for _, f := range t.Fields {
			c.putUnmarshalFuncSwitchCase(f)
			{
				c.pushIndent()
				c.putLine("if set {")
				{
					c.pushIndent()
					c.putLine("v.Set |= ")
					c.typer.putStructFieldFlag(f)
					c.popIndent()
				}
				c.putLine("}")
				c.putLine("if del {")
				{
					c.pushIndent()
					c.putLine("v.Del |= ")
					c.typer.putStructFieldFlag(f)
					c.popIndent()
				}
				c.putLine("}")
				c.putUnmarshalEmpty()
				c.popIndent()
			}
		}
		c.putUnmarshalFuncSwitchDefault(t)
		c.putLine("}")
		c.putLine("if t == 0 {")
		{
			c.pushIndent()
			c.putLine("break")
			c.popIndent()
		}
		c.putLine("}")
		c.popIndent()
	}
	c.putLine("}")
	c.putUnmarshalEmpty()
	c.putLine("return v, nil")
}

func (c *marshallerImpl) putMarshalUnionType(t *sch.CompositeType) {
	c.putLine("switch v := v.(type) {")
	for _, f := range t.Fields {
		c.putLine("case *")
		c.typer.putUnionOptionTypeRef(f)
		c.putString(":")
		{
			c.pushIndent()
			c.putMarshalTag(f, true)
			c.putLine("f := v.")
			c.putName(f.Name)
			c.putMarshalSubType(f, "d", "f")
			c.popIndent()
		}
	}
	c.putMarshalFuncSwitchDefault(t)
	c.putLine("}")
	c.putMarshalEmpty()
	c.putLine("return nil")
}

func (c *marshallerImpl) putUnmarshalUnionType(t *sch.CompositeType) {
	c.putLine("var v ")
	c.putCompoundName(t.Ref)
	c.putUnmarshalTag(t, "t", true)
	c.putLine("switch t {")
	for _, f := range t.Fields {
		c.putUnmarshalFuncSwitchCase(f)
		{
			c.pushIndent()
			c.putUnmarshalSubType(f, "d", "f")
			c.putUnmarshalFuncReturnErr(t)
			c.putLine("v = &")
			c.typer.putUnionOptionTypeRef(f)
			c.putString("{f}")
			c.popIndent()
		}
	}
	c.putUnmarshalFuncSwitchDefault(t)
	c.putLine("}")
	c.putUnmarshalEmpty()
	c.putLine("return v, nil")
}

func (c *marshallerImpl) putMarshalStructType(t *sch.CompositeType) {
	c.putLine("if v == nil {")
	{
		c.pushIndent()
		c.putLine("goto end")
		c.popIndent()
	}
	c.putLine("}")
	for _, f := range t.Fields {
		c.putLine("if set, del := v.Set&")
		c.typer.putStructFieldFlag(f)
		c.putString(" != 0,")
		{
			c.pushIndent()
			c.putLine("v.Del&")
			c.typer.putStructFieldFlag(f)
			c.putString(" != 0; set || del {")
			c.putMarshalTag(f, false)
			c.putLine("if set {")
			{
				c.pushIndent()
				c.putLine("f := v.")
				c.putName(f.Name)
				c.putMarshalSubType(f, "d", "f")
				c.popIndent()
			}
			c.putLine("} else {")
			{
				c.pushIndent()
				c.putMarshalEmpty()
				c.popIndent()
			}
			c.putLine("}")
			c.popIndent()
		}
		c.putLine("}")
	}
	{
		c.popIndent()
		c.putLine("end:")
		c.pushIndent()
	}
	c.putMarshalEnd()
	c.putMarshalEmpty()
	c.putLine("return nil")
}

func (c *marshallerImpl) putUnmarshalStructType(t *sch.CompositeType) {
	c.putLine("v := &")
	c.putCompoundName(t.Ref)
	c.putString("{}")
	c.putLine("for d.Next() {")
	{
		c.pushIndent()
		c.putUnmarshalTag(t, "t", false)
		c.putLine("switch t {")
		for _, f := range t.Fields {
			c.putUnmarshalFuncSwitchCase(f)
			{
				c.pushIndent()
				c.putLine("if set {")
				{
					c.pushIndent()
					c.putLine("v.Set |= ")
					c.typer.putStructFieldFlag(f)
					c.putUnmarshalSubType(f, "d", "f")
					c.putUnmarshalFuncReturnErr(t)
					c.putLine("v.")
					c.putName(f.Name)
					c.putString(" = f")
					c.popIndent()
				}
				c.putLine("} else {")
				{
					c.pushIndent()
					c.putUnmarshalEmpty()
					c.popIndent()
				}
				c.putLine("}")
				c.putLine("if del {")
				{
					c.pushIndent()
					c.putLine("v.Del |= ")
					c.typer.putStructFieldFlag(f)
					c.popIndent()
				}
				c.putLine("}")
				c.popIndent()
			}
		}
		c.putUnmarshalFuncSwitchDefault(t)
		c.putLine("}")
		c.putLine("if t == 0 {")
		{
			c.pushIndent()
			c.putLine("break")
			c.popIndent()
		}
		c.putLine("}")
		c.popIndent()
	}
	c.putLine("}")
	c.putUnmarshalEmpty()
	c.putLine("return v, nil")
}

func (c *marshallerImpl) putMarshalContainerType(t *sch.ContainerType) {
	switch t.Kind {
	case sch.ArrayKind:
		c.putMarshalArrayType(t)
	case sch.MapKind:
		c.putMarshalMapType(t)
	case sch.SetKind:
		c.putMarshalSetType(t)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putUnmarshalContainerType(t *sch.ContainerType) {
	switch t.Kind {
	case sch.ArrayKind:
		c.putUnmarshalArrayType(t)
	case sch.MapKind:
		c.putUnmarshalMapType(t)
	case sch.SetKind:
		c.putUnmarshalSetType(t)
	default:
		panic("")
	}
}

func (c *marshallerImpl) putMarshalArrayType(t *sch.ContainerType) {
	c.putMarshalSize("v")
	c.putLine("for k, e := range v {")
	{
		c.pushIndent()
		c.putMarshalIndex("k")
		c.putMarshalSubType(t.Elem, "d", "e")
		c.popIndent()
	}
	c.putLine("}")
	c.putMarshalEnd()
	c.putMarshalEmpty()
}

func (c *marshallerImpl) putUnmarshalArrayType(t *sch.ContainerType) {
	c.putUnmarshalSize(t, "n")
	c.putLine("var v []")
	c.typer.putSubTypeRef(t.Elem)
	c.putLine("for i := 0; i < n && d.Next(); i++ {")
	{
		c.pushIndent()
		c.putUnmarshalIndex(t, "k")
		c.putUnmarshalSubType(t.Elem, "d", "e")
		c.putUnmarshalFuncReturnErr(t)
		c.putLine("for cap(v) <= k {")
		{
			c.pushIndent()
			c.putLine("v = append(v, ")
			c.typer.putZero(t.Elem)
			c.putString(")")
			c.popIndent()
		}
		c.putLine("}")
		c.putLine("v = v[:k+1]")
		c.putLine("v[k] = e")
		c.popIndent()
	}
	c.putLine("}")
	c.putUnmarshalEmpty()
	c.putLine("return v, nil")
}

func (c *marshallerImpl) putMarshalMapType(t *sch.ContainerType) {
	c.putMarshalSize("v")
	c.putLine("for k, e := range v {")
	{
		c.pushIndent()
		c.putMarshalKeyType(t.Key, "k")
		c.putMarshalSubType(t.Elem, "d", "e")
		c.popIndent()
	}
	c.putLine("}")
	c.putMarshalEmpty()
}

func (c *marshallerImpl) putUnmarshalMapType(t *sch.ContainerType) {
	c.putUnmarshalSize(t, "n")
	c.putLine("v := ")
	c.typer.putTypeRef(t.Ref)
	c.putString("{}")
	c.putLine("for i := 0; i < n && d.Next(); i++ {")
	{
		c.pushIndent()
		c.putUnmarshalKeyType(t.Key, "k")
		c.putUnmarshalFuncReturnErr(t)
		c.putUnmarshalSubType(t.Elem, "d", "e")
		c.putUnmarshalFuncReturnErr(t)
		c.putLine("v[k] = e")
		c.popIndent()
	}
	c.putLine("}")
	c.putUnmarshalEmpty()
	c.putLine("return v, nil")
}

func (c *marshallerImpl) putMarshalSetType(t *sch.ContainerType) {
	c.putMarshalSize("v")
	c.putLine("for k, _ := range v {")
	{
		c.pushIndent()
		c.putMarshalKeyType(t.Key, "k")
		c.putMarshalEmpty()
		c.popIndent()
	}
	c.putLine("}")
	c.putMarshalEmpty()
}

func (c *marshallerImpl) putUnmarshalSetType(t *sch.ContainerType) {
	c.putUnmarshalSize(t, "n")
	c.putLine("v := ")
	c.typer.putTypeRef(t.Ref)
	c.putString("{}")
	c.putLine("for i := 0; i < n && d.Next(); i++ {")
	{
		c.pushIndent()
		c.putUnmarshalKeyType(t.Key, "k")
		c.putUnmarshalFuncReturnErr(t)
		c.putUnmarshalEmpty()
		c.putLine("v[k] = struct{}{}")
		c.popIndent()
	}
	c.putLine("}")
	c.putUnmarshalEmpty()
	c.putLine("return v, nil")
}

// Snippets

func (c *marshallerImpl) putMarshalFuncBegin(t *sch.ReferenceType) {
	c.putLine("func (m *Marshaller) Marshal")
	c.putCompoundName(t)
	c.putString("(d enc.Data, v ")
	c.typer.putTypeRef(t)
	c.putString(") (err error) {")
	c.pushIndent()
}

func (c *marshallerImpl) putUnmarshalFuncBegin(t *sch.ReferenceType) {
	c.putLine("func (m *Marshaller) Unmarshal")
	c.putCompoundName(t)
	c.putString("(d enc.Data) (")
	c.typer.putTypeRef(t)
	c.putString(", error) {")
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

func (c *marshallerImpl) putMarshalFuncReturnErr() {
	c.putLine("if err != nil {")
	{
		c.pushIndent()
		c.putLine("return err")
		c.popIndent()
	}
	c.putLine("}")
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

// func (c *marshallerImpl) putMarshalFuncSwitchCase(f *sch.ReferenceType) {}

func (c *marshallerImpl) putUnmarshalFuncSwitchCase(f *sch.ReferenceType) {
	c.putLine("case ")
	c.putTag(f.Tag)
	c.putString(":")
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
		c.putMarshalTag(nil, true)
		c.putMarshalEmpty()
		c.popIndent()
	}
	c.putLine("default:")
	{
		c.pushIndent()
		c.putLine("return enc.ErrUnknownTag")
		c.popIndent()
	}
}
func (c *marshallerImpl) putUnmarshalFuncSwitchDefault(t *sch.CompositeType) {
	c.putLine("case 0:")
	{
		c.pushIndent()
		c.putUnmarshalEmpty()
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

// Utilities

func (c *marshallerImpl) putMarshalTag(f *sch.ReferenceType, constant bool) {
	c.putLine("d.SetKey().EncodeTag(")
	if f != nil {
		c.putTag(f.Tag)
	} else {
		c.putString("0")
	}
	if !constant {
		c.putString(", set, del)")
	} else {
		c.putString(", true, true)")
	}
}

func (c *marshallerImpl) putUnmarshalTag(t *sch.CompositeType, v string, constant bool) {
	c.putLine(v)
	if constant {
		c.putString(", _, _")
	} else {
		c.putString(", set, del")
	}
	c.putString(", err := d.GetKey().DecodeTag()")
	c.putUnmarshalFuncReturnErr(t)
}

func (c *marshallerImpl) putMarshalEmpty() {
	c.putLine("d.Value()")
}

func (c *marshallerImpl) putUnmarshalEmpty() {
	c.putLine("d.Value()")
}

func (c *marshallerImpl) putMarshalEnd() {
	c.putLine("if b := d.End(); b != nil {")
	{
		c.pushIndent()
		c.putLine("b.EncodeUvarint(0)")
		c.popIndent()
	}
	c.putLine("}")
}

// func (c *marshallerImpl) putUnmarshalEnd() {}

func (c *marshallerImpl) putMarshalIndex(v string) {
	c.putLine("d.SetKey().EncodeSize(")
	c.putString(v)
	c.putString(")")
}

func (c *marshallerImpl) putUnmarshalIndex(t *sch.ContainerType, v string) {
	c.putLine(v)
	c.putString(", err := d.GetKey().DecodeSize()")
	c.putUnmarshalFuncReturnErr(t)
}

func (c *marshallerImpl) putMarshalSize(v string) {
	c.putLine("if b := d.Len(); b != nil {")
	{
		c.pushIndent()
		c.putLine("b.EncodeSize(len(")
		c.putString(v)
		c.putString("))")
		c.popIndent()
	}
	c.putLine("}")
}

func (c *marshallerImpl) putUnmarshalSize(t *sch.ContainerType, v string) {
	c.putLine(v)
	c.putString(" := enc.MaxSize")
	c.putLine("if b := d.Len(); b != nil {")
	{
		c.pushIndent()
		c.putLine("var err error")
		c.putLine(v)
		c.putString(", err = b.DecodeSize()")
		c.putUnmarshalFuncReturnErr(t)
		c.popIndent()
	}
	c.putLine("}")
}
