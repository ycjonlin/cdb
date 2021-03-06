package generator

import (
	sch "github.com/ycjonlin/cdb/pkg/schema"
)

type declarator interface {
	putSchema(s *sch.Schema)
}

type declaratorImpl struct {
	*writer
	typer
}

func (c *declaratorImpl) putSchema(s *sch.Schema) {
	// Package
	c.putString("package ")
	c.putString(s.Name)
	c.putLine("")
	// Types
	for _, r := range s.Types.Fields {
		if r.Tag == 0 || r.Name == "" {
			continue
		}
		c.putType(r)
	}
}

func (c *declaratorImpl) putType(r *sch.ReferenceType) {
	c.putLine("type ")
	c.putCompoundName(r)
	c.putString(" ")
	c.putTypeDef(r)
	c.putLine("")

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

func (c *declaratorImpl) putCompositeType(t *sch.CompositeType) {
	switch t.Kind {
	case sch.EnumKind:
		c.putEnumType(t)
	case sch.UnionKind:
		c.putUnionType(t)
		for _, f := range t.Fields {
			c.putSubType(f)
		}
	case sch.BitfieldKind:
		c.putBitfieldType(t)
	case sch.StructKind:
		c.putStructType(t)
		for _, f := range t.Fields {
			c.putSubType(f)
		}
	default:
		panic("")
	}
}

func (c *declaratorImpl) putEnumType(t *sch.CompositeType) {
	// Options
	c.pushIndent("const (")
	for i, f := range t.Fields {
		c.putLine("")
		c.putCompoundName(f)
		c.putString(" ")
		c.putCompoundName(f.Super)
		c.putString(" = ")
		c.putUint(uint64(i + 1))
	}
	c.popIndent(")")
	c.putLine("")
}

func (c *declaratorImpl) putUnionType(t *sch.CompositeType) {
	// Options
	for _, f := range t.Fields {
		c.putLine("type ")
		c.putUnionOptionTypeRef(f)
		c.putString(" ")
		c.putUnionOptionTypeDef(f)
	}
	c.putLine("")
	for _, f := range t.Fields {
		c.putLine("func (*")
		c.putUnionOptionTypeRef(f)
		c.putString(") is")
		c.putCompoundName(f.Super)
		c.putString("() {}")
	}
	c.putLine("")
}

func (c *declaratorImpl) putBitfieldType(t *sch.CompositeType) {
	c.putGetMethod(t)
	c.putSetMethod(t)
	c.putDelMethod(t)
}

func (c *declaratorImpl) putStructType(t *sch.CompositeType) {
}

func (c *declaratorImpl) putContainerType(t *sch.ContainerType) {
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

func (c *declaratorImpl) putSubType(r *sch.ReferenceType) {
	switch t := r.Type.(type) {
	case *sch.ReferenceType:
		// do nothing
	case *sch.PrimitiveType:
		// do nothing
	case *sch.CompositeType:
		c.putType(r)
	case *sch.ContainerType:
		//c.putType(r)
		c.putContainerType(t) // skip
	default:
		panic("")
	}
}

func (c *declaratorImpl) putGetMethod(t *sch.CompositeType) {
	for i, f := range t.Fields {
		c.putLine("func (c ")
		c.putTypeRef(t.Ref)
		c.putString(") Get")
		c.putName(f.Name)
		c.putString("() bool { return ")
		c.putString("c[")
		c.putUint(uint64(i >> 3))
		c.putString("]&0x")
		c.putHexUint(1 << uint(i&0x7))
		c.putString(" != 0 }")
	}
	c.putLine("")
}

func (c *declaratorImpl) putSetMethod(t *sch.CompositeType) {
	for i, f := range t.Fields {
		c.putLine("func (c ")
		c.putTypeRef(t.Ref)
		c.putString(") Set")
		c.putName(f.Name)
		c.putString("() { c[")
		c.putUint(uint64(i >> 3))
		c.putString("] |= 0x")
		c.putHexUint(1 << uint(i&0x7))
		c.putString(" }")
	}
	c.putLine("")
}

func (c *declaratorImpl) putDelMethod(t *sch.CompositeType) {
	for i, f := range t.Fields {
		c.putLine("func (c ")
		c.putTypeRef(t.Ref)
		c.putString(") Del")
		c.putName(f.Name)
		c.putString("() { c[")
		c.putUint(uint64(i >> 3))
		c.putString("] &^= 0x")
		c.putHexUint(1 << uint(i&0x7))
		c.putString(" }")
	}
	c.putLine("")
}
