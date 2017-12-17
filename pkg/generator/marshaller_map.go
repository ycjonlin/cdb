package generator

import (
	sch "github.com/ycjonlin/cdb/pkg/schema"
)

type mapMarshaller struct {
	*writer
	typer typer
}

func (c *mapMarshaller) putMarshallerDef() {
	c.putString("struct{ m enc.Map }")
}

func (c *mapMarshaller) putMarshallerRef() {
	c.putString("MapMarshaller")
}

func (c *mapMarshaller) putMarshalPrimitiveType(t *sch.PrimitiveType, v string) {
	switch t {
	case sch.BoolType:
		c.putLine("m.m.Data.EncodeBool(bool(")
	case sch.IntType:
		c.putLine("m.m.Data.EncodeVarint(int64(")
	case sch.UintType:
		c.putLine("m.m.Data.EncodeUvarint(uint64(")
	case sch.FloatType:
		c.putLine("m.m.Data.EncodeFloat64(float64(")
	case sch.StringType:
		c.putLine("m.m.Data.EncodeNonsortingString(string(")
	case sch.BytesType:
		c.putLine("m.m.Data.EncodeNonsortingBytes([]byte(")
	default:
		panic("")
	}
	c.putString(v)
	c.putString("))")
	c.putLine("m.m.Put()")
}

func (c *mapMarshaller) putMarshalTagBegin(f *sch.ReferenceType) {
	c.putLine("m.m.Path.Push()")
	c.putLine("m.m.Path.EncodeUvarint(")
	c.putTag(f.Tag)
	c.putString(")")
}

func (c *mapMarshaller) putMarshalTagEnd(f *sch.ReferenceType) {
	c.putLine("m.m.Path.Pop()")
}

func (c *mapMarshaller) putMarshalZeroValue() {
	c.putLine("m.m.Put()")
}

func (c *mapMarshaller) putMarshalZeroTag() {}

func (c *mapMarshaller) putMarshalLen() {}

func (c *mapMarshaller) putUnmarshalPrimitiveType(t *sch.PrimitiveType, v string) {
	c.putLine(v)
	switch t {
	case sch.BoolType:
		c.putString(", err := m.m.Data.DecodeBool()")
	case sch.IntType:
		c.putString(", err := m.m.Data.DecodeVarint()")
	case sch.UintType:
		c.putString(", err := m.m.Data.DecodeUvarint()")
	case sch.FloatType:
		c.putString(", err := m.m.Data.DecodeFloat64()")
	case sch.StringType:
		c.putString(", err := m.m.Data.DecodeNonsortingString()")
	case sch.BytesType:
		c.putString(", err := m.m.Data.DecodeNonsortingBytes(nil)")
	default:
		panic("")
	}
}

func (c *mapMarshaller) putUnmarshalTag(t *sch.CompositeType) {
	c.putLine("t, err := 0, error(nil)")
}

func (c *mapMarshaller) putUnmarshalLen(t *sch.ContainerType) {
	c.putLine("n, err := 0, error(nil)")
}
