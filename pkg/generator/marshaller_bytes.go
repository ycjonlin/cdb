package generator

import (
	sch "github.com/ycjonlin/cdb/pkg/schema"
)

type bytesMarshaller struct {
	*writer
	typer typer
}

func (c *bytesMarshaller) putMarshallerDef() {
	c.putString("struct{ b enc.Bytes }")
}

func (c *bytesMarshaller) putMarshallerRef() {
	c.putString("BytesMarshaller")
}

func (c *bytesMarshaller) putMarshalPrimitiveType(t *sch.PrimitiveType, v string) {
	switch t {
	case sch.BoolType:
		c.putLine("m.b.EncodeBool(bool(")
	case sch.IntType:
		c.putLine("m.b.EncodeVarint(int64(")
	case sch.UintType:
		c.putLine("m.b.EncodeUvarint(uint64(")
	case sch.FloatType:
		c.putLine("m.b.EncodeFloat64(float64(")
	case sch.StringType:
		c.putLine("m.b.EncodeNonsortingString(string(")
	case sch.BytesType:
		c.putLine("m.b.EncodeNonsortingBytes([]byte(")
	default:
		panic("")
	}
	c.putString(v)
	c.putString("))")
}

func (c *bytesMarshaller) putMarshalTagBegin(f *sch.ReferenceType) {
	c.putLine("m.b.EncodeUvarint(")
	c.putTag(f.Tag)
	c.putString(")")
}

func (c *bytesMarshaller) putMarshalTagEnd(f *sch.ReferenceType) {
}

func (c *bytesMarshaller) putMarshalZeroValue() {}

func (c *bytesMarshaller) putMarshalZeroTag() {
	c.putLine("m.b.EncodeUvarint(0)")
}

func (c *bytesMarshaller) putMarshalLen() {
	c.putLine("m.b.EncodeUvarint(uint64(len(v)))")
}

func (c *bytesMarshaller) putUnmarshalPrimitiveType(t *sch.PrimitiveType, v string) {
	c.putLine(v)
	switch t {
	case sch.BoolType:
		c.putString(", err := m.b.DecodeBool()")
	case sch.IntType:
		c.putString(", err := m.b.DecodeVarint()")
	case sch.UintType:
		c.putString(", err := m.b.DecodeUvarint()")
	case sch.FloatType:
		c.putString(", err := m.b.DecodeFloat64()")
	case sch.StringType:
		c.putString(", err := m.b.DecodeNonsortingString()")
	case sch.BytesType:
		c.putString(", err := m.b.DecodeNonsortingBytes(nil)")
	default:
		panic("")
	}
}

func (c *bytesMarshaller) putUnmarshalTag(t *sch.CompositeType) {
	c.putLine("t, err := m.b.DecodeUvarint()")
}

func (c *bytesMarshaller) putUnmarshalLen(t *sch.ContainerType) {
	c.putLine("n, err := m.b.DecodeSize()")
}
