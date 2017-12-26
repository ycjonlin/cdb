package generator

import (
	sch "github.com/ycjonlin/cdb/pkg/schema"
)

type typer interface {
	putZero(t sch.Type)
	putTypeDef(r *sch.ReferenceType)
	putTypeRef(r *sch.ReferenceType)
	putSubTypeRef(r *sch.ReferenceType)
	putKeyTypeRef(r *sch.ReferenceType)
	putConvertType(r *sch.ReferenceType, v string)
	putUnionOptionTypeDef(f *sch.ReferenceType)
	putUnionOptionTypeRef(f *sch.ReferenceType)
	putStructFieldFlagType(f *sch.ReferenceType)
	putStructFieldFlag(f *sch.ReferenceType)
}

type typerImpl struct {
	*writer
}

func (c *typerImpl) putZero(t sch.Type) {
	if r, ok := t.(*sch.ReferenceType); ok {
		t = r.Base
	}
	switch t := t.(type) {
	case *sch.PrimitiveType:
		c.putPrimitiveTypeZero(t)
	case *sch.CompositeType:
		c.putCompositeTypeZero(t)
	case *sch.ContainerType:
		c.putString("nil")
	default:
		panic("")
	}
}

func (c *typerImpl) putPrimitiveTypeZero(t *sch.PrimitiveType) {
	switch t {
	case sch.BoolType:
		c.putString("false")
	case sch.IntType:
		c.putString("0")
	case sch.UintType:
		c.putString("0")
	case sch.FloatType:
		c.putString("0")
	case sch.BytesType:
		c.putString("nil")
	case sch.StringType:
		c.putString("\"\"")
	default:
		panic("")
	}
}

func (c *typerImpl) putCompositeTypeZero(t *sch.CompositeType) {
	switch t.Kind {
	case sch.EnumKind:
		c.putString("0")
	case sch.BitfieldKind:
		c.putString("nil")
	case sch.UnionKind:
		c.putString("nil")
	case sch.StructKind:
		c.putString("nil")
	default:
		panic("")
	}
}

func (c *typerImpl) putTypeDef(r *sch.ReferenceType) {
	switch t := r.Type.(type) {
	case *sch.ReferenceType:
		c.putReferenceType(t)
	case *sch.PrimitiveType:
		c.putPrimitiveType(t)
	case *sch.CompositeType:
		c.putCompositeType(t)
	case *sch.ContainerType:
		c.putContainerType(t)
	default:
		panic("")
	}
}

func (c *typerImpl) putTypeRef(r *sch.ReferenceType) {
	ct, ok := r.Base.(*sch.CompositeType)
	if ok && (ct.Kind == sch.StructKind || ct.Kind == sch.BitfieldKind) {
		c.putString("*")
	}
	c.putCompoundName(r)
	/*switch r.Type.(type) {
	case *sch.ReferenceType:
		c.putCompoundName(r)
	case *sch.PrimitiveType:
		c.putCompoundName(r)
	case *sch.CompositeType:
		c.putCompoundName(r)
	case *sch.ContainerType:
		//c.putContainerType(t)
		c.putCompoundName(r)
	default:
		panic("")
	}*/
}

func (c *typerImpl) putConvertType(r *sch.ReferenceType, v string) {
	ct, ok := r.Base.(*sch.CompositeType)
	if ok && (ct.Kind == sch.StructKind || ct.Kind == sch.BitfieldKind) {
		c.putString("(")
		c.putTypeRef(r)
		c.putString(")")
	} else {
		c.putTypeRef(r)
	}
	c.putString("(")
	c.putString(v)
	c.putString(")")
}

func (c *typerImpl) putPrimitiveType(t *sch.PrimitiveType) {
	n, ok := primitiveTypeName[t.Name]
	if !ok {
		panic("")
	}
	c.putString(n)
}

func (c *typerImpl) putReferenceType(t *sch.ReferenceType) {
	c.putCompoundName(t)
}

var primitiveTypeName = map[sch.Name]string{
	"void":      "struct{}",
	"bool":      "bool",
	"int":       "int64",
	"int8":      "int8",
	"int16":     "int16",
	"int32":     "int32",
	"int64":     "int64",
	"uint":      "uint64",
	"uint8":     "uint8",
	"uint16":    "uint16",
	"uint32":    "uint32",
	"uint64":    "uint64",
	"float":     "float64",
	"float32":   "float32",
	"float64":   "float64",
	"bytes":     "[]byte",
	"string":    "string",
	"timestamp": "int64",
	"duration":  "int64",
	"geopoint":  "uint64",
}

func (c *typerImpl) putCompositeType(t *sch.CompositeType) {
	switch t.Kind {
	case sch.EnumKind:
		c.putString("uint32")
	case sch.BitfieldKind:
		c.putBitfieldType(t)
	case sch.UnionKind:
		c.putUnionType(t)
	case sch.StructKind:
		c.putStructType(t)
	default:
		panic("")
	}
}

func (c *typerImpl) putBitfieldType(t *sch.CompositeType) {
	c.putString("struct{ Del, Set ")
	c.putStructFieldFlagType(t.Ref)
	c.putString(" }")
}

func (c *typerImpl) putUnionType(t *sch.CompositeType) {
	c.putString("interface {")
	c.pushIndent()
	{
		c.putLine("is")
		c.putCompoundName(t.Ref)
		c.putString("()")
	}
	c.popIndent()
	c.putLine("}")
}

func (c *typerImpl) putUnionOptionTypeDef(f *sch.ReferenceType) {
	c.putString("struct{ ")
	c.putName(f.Name)
	c.putString(" ")
	c.putSubTypeRef(f)
	c.putString(" }")
}

func (c *typerImpl) putUnionOptionTypeRef(r *sch.ReferenceType) {
	c.putCompoundName(r.Super)
	c.putString("_As_")
	c.putName(r.Name)
}

func (c *typerImpl) putStructFieldFlagType(r *sch.ReferenceType) {
	c.putCompoundName(r)
	c.putString("_Flags")
}

func (c *typerImpl) putStructFieldFlag(f *sch.ReferenceType) {
	c.putCompoundName(f)
	c.putString("_Flag")
}

func (c *typerImpl) putStructType(t *sch.CompositeType) {
	c.putString("struct {")
	c.pushIndent()
	c.putLine("Set, Del ")
	c.putCompoundName(t.Ref)
	c.putString("_Flags")
	for _, f := range t.Fields {
		c.putLine("")
		c.putName(f.Name)
		c.putString(" ")
		c.putSubTypeRef(f)
	}
	c.popIndent()
	c.putLine("}")
}

func (c *typerImpl) putStructFieldTypeDef(f *sch.ReferenceType) {
	c.putString("uint64")
}

func (c *typerImpl) putStructFieldTypeRef(r *sch.ReferenceType) {
	c.putCompoundName(r)
	c.putString("_Flag")
}

func (c *typerImpl) putStructFieldTypeIns(r *sch.ReferenceType) {
	c.putCompoundName(r.Super)
	c.putString("_Flag_")
	c.putName(r.Name)
}

func (c *typerImpl) putContainerType(t *sch.ContainerType) {
	switch t.Kind {
	case sch.ArrayKind:
		c.putString("[]")
		c.putSubTypeRef(t.Elem)
	case sch.MapKind:
		c.putString("map[")
		c.putKeyTypeRef(t.Key)
		c.putString("]")
		c.putSubTypeRef(t.Elem)
	case sch.SetKind:
		c.putString("map[")
		c.putKeyTypeRef(t.Key)
		c.putString("]struct{}")
	default:
		panic("")
	}
}

func (c *typerImpl) putKeyTypeRef(r *sch.ReferenceType) {
	switch t := r.Base.(type) {
	case *sch.PrimitiveType:
		if t == sch.BytesType {
			c.putString("string")
		} else {
			c.putSubTypeRef(r)
		}
	case *sch.CompositeType:
		switch t.Kind {
		case sch.EnumKind, sch.BitfieldKind:
			c.putSubTypeRef(r)
		case sch.UnionKind, sch.StructKind:
			c.putString("string")
		default:
			panic("")
		}
	case *sch.ContainerType:
		c.putString("string")
	default:
		panic("")
	}
}

func (c *typerImpl) putSubTypeRef(r *sch.ReferenceType) {
	ct, ok := r.Base.(*sch.CompositeType)
	if ok && (ct.Kind == sch.StructKind || ct.Kind == sch.BitfieldKind) {
		c.putString("*")
	}
	switch t := r.Type.(type) {
	case *sch.ReferenceType:
		c.putReferenceType(t)
	case *sch.PrimitiveType:
		c.putPrimitiveType(t)
	case *sch.CompositeType:
		c.putCompoundName(r)
	case *sch.ContainerType:
		c.putContainerType(t)
	default:
		panic("")
	}
}
