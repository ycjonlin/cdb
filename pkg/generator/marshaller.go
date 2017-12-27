package generator

import (
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

func (c *marshallerImpl) putUnmarshalSubType(t sch.Type, r *sch.ReferenceType, e, v string) {
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
	c.putUnmarshalFuncReturnErr(t)
}

func (c *marshallerImpl) putMarshalKeyType(r *sch.ReferenceType) {
	c.putLine("dk := d.Data(d.SetKey())")
	c.putMarshalSubType(r, "dk", "k")
}

func (c *marshallerImpl) putUnmarshalKeyType(t sch.Type, r *sch.ReferenceType) {
	c.putLine("dk := d.Data(d.GetKey())")
	c.putUnmarshalSubType(t, r, "dk", "k")
}

func (c *marshallerImpl) putMarshalReferenceType(t *sch.ReferenceType, e, v string) {
	c.putLine("err = m.Marshal")
	c.putCompoundName(t)
	c.putString("(")
	c.putString(e)
	c.putString(", ")
	c.putConvertType(t, v)
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
	c.putMarshalFieldOneSwitchBegin(t)
	for _, f := range t.Fields {
		c.putMarshalFieldOneSwitchCase(t, f)
		c.putMarshalEmpty()
	}
	c.putMarshalFieldOneSwitchEnd(t)
}

func (c *marshallerImpl) putUnmarshalEnumType(t *sch.CompositeType) {
	c.putUnmarshalFieldOneSwitchBegin(t)
	for _, f := range t.Fields {
		c.putUnmarshalFieldOneSwitchCase(f)
		c.putUnmarshalEmpty()
		c.putLine("v = ")
		c.putCompoundName(f)
	}
	c.putUnmarshalFieldOneSwitchEnd(t)
}

func (c *marshallerImpl) putMarshalUnionType(t *sch.CompositeType) {
	c.putMarshalFieldOneSwitchBegin(t)
	for _, f := range t.Fields {
		c.putMarshalFieldOneSwitchCase(t, f)
		c.putLine("f := v.")
		c.putName(f.Name)
		c.putMarshalSubType(f, "d", "f")
	}
	c.putMarshalFieldOneSwitchEnd(t)
}

func (c *marshallerImpl) putUnmarshalUnionType(t *sch.CompositeType) {
	c.putUnmarshalFieldOneSwitchBegin(t)
	for _, f := range t.Fields {
		c.putUnmarshalFieldOneSwitchCase(f)
		c.putUnmarshalSubType(t, f, "d", "f")
		c.putLine("v = &")
		c.putUnionOptionTypeRef(f)
		c.putString("{f}")
	}
	c.putUnmarshalFieldOneSwitchEnd(t)
}

func (c *marshallerImpl) putMarshalBitfieldType(t *sch.CompositeType) {
	c.putMarshalFieldManySwitchBegin(t)
	for i, f := range t.Fields {
		c.putMarshalFieldManySwitchCase(f, i)
		c.putMarshalEmpty()
	}
	c.putMarshalFieldManySwitchEnd(t)
}

func (c *marshallerImpl) putUnmarshalBitfieldType(t *sch.CompositeType) {
	c.putUnmarshalFieldManySwitchBegin(t)
	for i, f := range t.Fields {
		c.putUnmarshalFieldManySwitchCase(f, i)
		c.putUnmarshalEmpty()
	}
	c.putUnmarshalFieldManySwitchEnd(t)
}

func (c *marshallerImpl) putMarshalStructType(t *sch.CompositeType) {
	c.putMarshalFieldManySwitchBegin(t)
	for i, f := range t.Fields {
		c.putMarshalFieldManySwitchCase(f, i)
		c.pushIndent("if set {")
		{
			c.putLine("f := v.")
			c.putName(f.Name)
			c.putMarshalSubType(f, "d", "f")
		}
		c.nextIndent("} else {")
		{
			c.putMarshalEmpty()
		}
		c.popIndent("}")
	}
	c.putMarshalFieldManySwitchEnd(t)
}

func (c *marshallerImpl) putUnmarshalStructType(t *sch.CompositeType) {
	c.putUnmarshalFieldManySwitchBegin(t)
	for i, f := range t.Fields {
		c.putUnmarshalFieldManySwitchCase(f, i)
		c.pushIndent("if set {")
		{
			c.putUnmarshalSubType(t, f, "d", "f")
			c.putLine("v.")
			c.putName(f.Name)
			c.putString(" = f")
		}
		c.nextIndent("} else {")
		{
			c.putUnmarshalEmpty()
		}
		c.popIndent("}")
	}
	c.putUnmarshalFieldManySwitchEnd(t)
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
	c.putMarshalContainerBegin(t)
	{
		c.putLine("d.SetKey().EncodeSize(k)")
		c.putMarshalSubType(t.Elem, "d", "e")
	}
	c.putMarshalContainerEnd(t)
}

func (c *marshallerImpl) putUnmarshalArrayType(t *sch.ContainerType) {
	c.putUnmarshalContainerBegin(t)
	{
		c.putLine("k, err := d.GetKey().DecodeSize()")
		c.putUnmarshalFuncReturnErr(t)
		c.putUnmarshalSubType(t, t.Elem, "d", "e")
		c.pushIndent("for cap(v) <= k {")
		{
			c.putLine("v = append(v, ")
			c.putZero(t.Elem)
			c.putString(")")
		}
		c.popIndent("}")
		c.putLine("v = v[:k+1]")
		c.putLine("v[k] = e")
	}
	c.putUnmarshalContainerEnd(t)
}

func (c *marshallerImpl) putMarshalMapType(t *sch.ContainerType) {
	c.putMarshalContainerBegin(t)
	{
		c.putMarshalKeyType(t.Key)
		c.putMarshalSubType(t.Elem, "d", "e")
	}
	c.putMarshalContainerEnd(t)
}

func (c *marshallerImpl) putUnmarshalMapType(t *sch.ContainerType) {
	c.putUnmarshalContainerBegin(t)
	{
		c.putUnmarshalKeyType(t, t.Key)
		c.putUnmarshalSubType(t, t.Elem, "d", "e")
		c.putLine("v[k] = e")
	}
	c.putUnmarshalContainerEnd(t)
}

func (c *marshallerImpl) putMarshalSetType(t *sch.ContainerType) {
	c.putMarshalContainerBegin(t)
	{
		c.putMarshalKeyType(t.Key)
		c.putMarshalEmpty()
	}
	c.putMarshalContainerEnd(t)
}

func (c *marshallerImpl) putUnmarshalSetType(t *sch.ContainerType) {
	c.putUnmarshalContainerBegin(t)
	{
		c.putUnmarshalKeyType(t, t.Key)
		c.putUnmarshalEmpty()
		c.putLine("v[k] = struct{}{}")
	}
	c.putUnmarshalContainerEnd(t)
}

func (c *marshallerImpl) putMarshalContainerBegin(t *sch.ContainerType) {
	c.pushIndent("if b := d.Len(); b != nil {")
	{
		c.putLine("b.EncodeSize(len(v))")
	}
	c.popIndent("}")
	c.pushIndent("for k")
	if t.Kind != sch.SetKind {
		c.putString(", e")
	}
	c.putString(" := range v {")
}

func (c *marshallerImpl) putMarshalContainerEnd(t *sch.ContainerType) {
	c.popIndent("}")
	c.putMarshalEmpty()
}

func (c *marshallerImpl) putUnmarshalContainerBegin(t *sch.ContainerType) {
	c.putLine("n := enc.MaxSize")
	c.pushIndent("if b := d.Len(); b != nil {")
	{
		c.putLine("var err error")
		c.putLine("n, err = b.DecodeSize()")
		c.putUnmarshalFuncReturnErr(t)
	}
	c.popIndent("}")
	c.putLine("v := ")
	c.putTypeRef(t.Ref)
	c.putString("{}")
	c.pushIndent("for i := 0; i < n && d.Next(); i++ {")
}

func (c *marshallerImpl) putUnmarshalContainerEnd(t *sch.ContainerType) {
	c.popIndent("}")
	c.putUnmarshalEmpty()
	c.putLine("return v, nil")
}

// Snippets

func (c *marshallerImpl) putMarshalFuncBegin(t *sch.ReferenceType) {
	c.pushIndent("func (m *Marshaller) Marshal")
	c.putCompoundName(t)
	c.putString("(d enc.Data, v ")
	c.putTypeRef(t)
	c.putString(") (err error) {")
}

func (c *marshallerImpl) putUnmarshalFuncBegin(t *sch.ReferenceType) {
	c.pushIndent("func (m *Marshaller) Unmarshal")
	c.putCompoundName(t)
	c.putString("(d enc.Data) (")
	c.putTypeRef(t)
	c.putString(", error) {")
}

func (c *marshallerImpl) putMarshalFuncEnd(t *sch.ReferenceType) {
	if _, ok := t.Type.(*sch.CompositeType); !ok {
		c.putLine("return nil")
	}
	c.popIndent("}")
	c.putLine("")
}

func (c *marshallerImpl) putUnmarshalFuncEnd(t *sch.ReferenceType) {
	switch t.Type.(type) {
	case *sch.ReferenceType, *sch.PrimitiveType:
		c.putLine("return (")
		c.putTypeRef(t)
		c.putString(")(v), err")
	}
	c.popIndent("}")
	c.putLine("")
}

func (c *marshallerImpl) putMarshalFuncReturnErr() {
	c.pushIndent("if err != nil {")
	{
		c.putLine("return err")
	}
	c.popIndent("}")
}

func (c *marshallerImpl) putUnmarshalFuncReturnErr(t sch.Type) {
	c.pushIndent("if err != nil {")
	{
		c.putLine("return ")
		c.putZero(t)
		c.putString(", err")
	}
	c.popIndent("}")
}

func (c *marshallerImpl) putMarshalFieldOneSwitchBegin(t *sch.CompositeType) {
	c.pushIndent("switch v ")
	if t.Kind == sch.UnionKind {
		c.putString(":= v.(type) ")
	}
	c.putString("{")
}

func (c *marshallerImpl) putMarshalFieldOneSwitchCase(t *sch.CompositeType, f *sch.ReferenceType) {
	c.nextIndent("case ")
	if t.Kind == sch.UnionKind {
		c.putString("*")
		c.putUnionOptionTypeRef(f)
	} else {
		c.putCompoundName(f)
	}
	c.putString(":")
	{
		c.putMarshalTag(f)
	}
}

func (c *marshallerImpl) putMarshalFieldOneSwitchEnd(t *sch.CompositeType) {
	c.nextIndent("case ")
	c.putZero(t)
	c.putString(":")
	{
		c.putMarshalTag(nil)
		c.putMarshalEmpty()
	}
	c.nextIndent("default:")
	{
		c.putLine("return enc.ErrUnknownTag")
	}
	c.popIndent("}")
	c.putMarshalEmpty()
	c.putLine("return nil")
}

func (c *marshallerImpl) putUnmarshalFieldOneSwitchBegin(t *sch.CompositeType) {
	c.putLine("var v ")
	c.putCompoundName(t.Ref)
	c.putUnmarshalTag(t)
	c.pushIndent("switch t {")
}

func (c *marshallerImpl) putUnmarshalFieldOneSwitchCase(f *sch.ReferenceType) {
	c.nextIndent("case ")
	c.putTag(f.Tag)
	c.putString(":")
}

func (c *marshallerImpl) putUnmarshalFieldOneSwitchEnd(t *sch.CompositeType) {
	c.nextIndent("case 0:")
	{
		c.putUnmarshalEmpty()
	}
	c.putUnmarshalFieldSwitchDefault(t)
	c.popIndent("}")
	c.putUnmarshalEmpty()
	c.putLine("return v, nil")
}

func (c *marshallerImpl) putMarshalFieldManySwitchBegin(t *sch.CompositeType) {
	c.pushIndent("if v == nil {")
	{
		c.putLine("goto end")
	}
}

func (c *marshallerImpl) putMarshalFieldManySwitchCase(f *sch.ReferenceType, i int) {
	c.popIndent("}")
	c.pushIndent("if set, del := v.Set[")
	c.putUint(uint64(i >> 3))
	c.putString("]&0x")
	c.putHexUint(1 << uint(i&0x7))
	c.putString(" != 0, v.Del[")
	c.putUint(uint64(i >> 3))
	c.putString("]&0x")
	c.putHexUint(1 << uint(i&0x7))
	c.putString(" != 0; set || del {")
	{
		c.putMarshalTag(f)
		c.putLine("d.SetFlags(set, del)")
	}
}

func (c *marshallerImpl) putMarshalFieldManySwitchEnd(t *sch.CompositeType) {
	c.popIndent("}")
	c.nextIndent("end:")
	c.pushIndent("if b := d.End(); b != nil {")
	{
		c.putLine("b.EncodeUvarint(0)")
	}
	c.popIndent("}")
	c.putMarshalEmpty()
	c.putLine("return nil")
}

func (c *marshallerImpl) putUnmarshalFieldManySwitchBegin(t *sch.CompositeType) {
	c.putLine("v := &")
	c.putCompoundName(t.Ref)
	c.putString("{}")
	c.pushIndent("for d.Next() {")
	{
		c.putUnmarshalTag(t)
		c.pushIndent("if t == 0 {")
		{
			c.putUnmarshalEmpty()
			c.putLine("break")
		}
		c.popIndent("}")
		c.putLine("set, del, err := d.GetFlags()")
		c.putUnmarshalFuncReturnErr(t)
		c.pushIndent("switch t {")
	}
}

func (c *marshallerImpl) putUnmarshalFieldManySwitchCase(f *sch.ReferenceType, i int) {
	{
		c.nextIndent("case ")
		c.putTag(f.Tag)
		c.putString(":")
		{
			c.pushIndent("if del {")
			{
				c.putLine("v.Del[")
				c.putUint(uint64(i >> 3))
				c.putString("] |= 0x")
				c.putHexUint(1 << uint(i&0x7))
			}
			c.popIndent("}")
			c.pushIndent("if set {")
			{
				c.putLine("v.Set[")
				c.putUint(uint64(i >> 3))
				c.putString("] |= 0x")
				c.putHexUint(1 << uint(i&0x7))
			}
			c.popIndent("}")
		}
	}
}

func (c *marshallerImpl) putUnmarshalFieldManySwitchEnd(t *sch.CompositeType) {
	{
		c.putUnmarshalFieldSwitchDefault(t)
		c.popIndent("}")
	}
	c.popIndent("}")
	c.putUnmarshalEmpty()
	c.putLine("return v, nil")
}

func (c *marshallerImpl) putUnmarshalFieldSwitchDefault(t *sch.CompositeType) {
	c.nextIndent("default:")
	{
		c.putLine("return ")
		c.putZero(t)
		c.putString(", enc.ErrUnknownTag")
	}
}

// Utilities

func (c *marshallerImpl) putMarshalTag(f *sch.ReferenceType) {
	c.putLine("d.SetKey().EncodeSize(")
	if f != nil {
		c.putTag(f.Tag)
	} else {
		c.putString("0")
	}
	c.putString(")")
}

func (c *marshallerImpl) putUnmarshalTag(t *sch.CompositeType) {
	c.putLine("t, err := d.GetKey().DecodeSize()")
	c.putUnmarshalFuncReturnErr(t)
}

func (c *marshallerImpl) putMarshalEmpty() {
	c.putLine("d.Value()")
}

func (c *marshallerImpl) putUnmarshalEmpty() {
	c.putLine("d.Value()")
}
