package encoding_test

import (
	enc "github.com/ycjonlin/cdb/pkg/encoding"
	sch "github.com/ycjonlin/cdb/pkg/schema"
)

var TestMarshaller = func() *enc.Marshaller {
	m := enc.NewMarshaller()
	{
		r := m.MustPut(1, "Bool", Bool(false))
		r.SetType(sch.BoolType)
	}
	{
		r := m.MustPut(2, "Int", Int(0))
		r.SetType(sch.IntType)
	}
	{
		r := m.MustPut(3, "Uint", Uint(0))
		r.SetType(sch.UintType)
	}
	{
		r := m.MustPut(4, "Float", Float(0))
		r.SetType(sch.FloatType)
	}
	{
		r := m.MustPut(5, "String", String(""))
		r.SetType(sch.StringType)
	}
	{
		r := m.MustPut(6, "Bytes", make(Bytes, 0))
		r.SetType(sch.BytesType)
	}
	{
		r := m.MustPut(7, "Type", Type(0))
		r.SetType(m.GetByName("Int"))
	}
	{
		r := m.MustPut(11, "Enum", Enum(0))
		c := sch.NewCompositeType(r, sch.EnumKind, nil)
		c.MustPut(1, "Option1", nil)
		c.MustPut(2, "Option2", nil)
	}
	{
		r := m.MustPut(12, "Union", new(Union))
		c := sch.NewCompositeType(r, sch.UnionKind, new(Union))
		{
			r := c.MustPut(1, "Option1", new(Union_As_Option1))
			r.SetType(sch.BytesType)
		}
		{
			r := c.MustPut(2, "Option2", new(Union_As_Option2))
			r.SetType(m.GetByName("Type"))
		}
	}
	{
		r := m.MustPut(13, "Bitfield", new(Bitfield))
		c := sch.NewCompositeType(r, sch.BitfieldKind, nil)
		c.MustPut(1, "Field1", nil)
		c.MustPut(2, "Field2", nil)
	}
	{
		r := m.MustPut(14, "Struct", new(Struct))
		c := sch.NewCompositeType(r, sch.StructKind, nil)
		{
			r := c.MustPut(1, "Field1", nil)
			r.SetType(sch.BytesType)
		}
		{
			r := c.MustPut(2, "Field2", nil)
			r.SetType(m.GetByName("Type"))
		}
	}
	{
		r := m.MustPut(15, "Array", make(Array, 0))
		c := sch.NewContainerType(r, sch.ArrayKind)
		{
			r := c.Elem
			r.SetType(m.GetByName("Type"))
		}
	}
	{
		r := m.MustPut(16, "Map", make(Map))
		c := sch.NewContainerType(r, sch.MapKind)
		{
			r := c.Key
			r.SetType(m.GetByName("Type"))
		}
		{
			r := c.Elem
			r.SetType(m.GetByName("Type"))
		}
	}
	{
		r := m.MustPut(17, "Set", make(Set))
		c := sch.NewContainerType(r, sch.SetKind)
		{
			r := c.Key
			r.SetType(m.GetByName("Type"))
		}
	}
	{
		r := m.MustPut(21, "NestedUnion", new(NestedUnion))
		c := sch.NewCompositeType(r, sch.UnionKind, new(NestedUnion))
		{
			r := c.MustPut(1, "Enum", new(NestedUnion_As_Enum))
			r.SetType(m.GetByName("Enum"))
		}
		{
			r := c.MustPut(2, "Union", new(NestedUnion_As_Union))
			r.SetType(m.GetByName("Union"))
		}
		{
			r := c.MustPut(3, "Bitfield", new(NestedUnion_As_Bitfield))
			r.SetType(m.GetByName("Bitfield"))
		}
		{
			r := c.MustPut(4, "Struct", new(NestedUnion_As_Struct))
			r.SetType(m.GetByName("Struct"))
		}
		{
			r := c.MustPut(5, "Array", new(NestedUnion_As_Array))
			r.SetType(m.GetByName("Array"))
		}
		{
			r := c.MustPut(6, "Map", new(NestedUnion_As_Map))
			r.SetType(m.GetByName("Map"))
		}
		{
			r := c.MustPut(7, "Set", new(NestedUnion_As_Set))
			r.SetType(m.GetByName("Set"))
		}
	}
	{
		r := m.MustPut(22, "NestedStruct", new(NestedStruct))
		c := sch.NewCompositeType(r, sch.StructKind, nil)
		{
			r := c.MustPut(1, "Enum", nil)
			r.SetType(m.GetByName("Enum"))
		}
		{
			r := c.MustPut(2, "Union", nil)
			r.SetType(m.GetByName("Union"))
		}
		{
			r := c.MustPut(3, "Bitfield", nil)
			r.SetType(m.GetByName("Bitfield"))
		}
		{
			r := c.MustPut(4, "Struct", nil)
			r.SetType(m.GetByName("Struct"))
		}
		{
			r := c.MustPut(5, "Array", nil)
			r.SetType(m.GetByName("Array"))
		}
		{
			r := c.MustPut(6, "Map", nil)
			r.SetType(m.GetByName("Map"))
		}
		{
			r := c.MustPut(7, "Set", nil)
			r.SetType(m.GetByName("Set"))
		}
	}
	{
		r := m.MustPut(23, "NestedArray", new(NestedArray))
		c := sch.NewCompositeType(r, sch.UnionKind, new(NestedArray))
		{
			r := c.MustPut(1, "Enum", new(NestedArray_As_Enum))
			c := sch.NewContainerType(r, sch.ArrayKind)
			{
				r := c.Elem
				r.SetType(m.GetByName("Enum"))
			}
		}
		{
			r := c.MustPut(2, "Union", new(NestedArray_As_Union))
			c := sch.NewContainerType(r, sch.ArrayKind)
			{
				r := c.Elem
				r.SetType(m.GetByName("Union"))
			}
		}
		{
			r := c.MustPut(3, "Bitfield", new(NestedArray_As_Bitfield))
			c := sch.NewContainerType(r, sch.ArrayKind)
			{
				r := c.Elem
				r.SetType(m.GetByName("Bitfield"))
			}
		}
		{
			r := c.MustPut(4, "Struct", new(NestedArray_As_Struct))
			c := sch.NewContainerType(r, sch.ArrayKind)
			{
				r := c.Elem
				r.SetType(m.GetByName("Struct"))
			}
		}
		{
			r := c.MustPut(5, "Array", new(NestedArray_As_Array))
			c := sch.NewContainerType(r, sch.ArrayKind)
			{
				r := c.Elem
				r.SetType(m.GetByName("Array"))
			}
		}
		{
			r := c.MustPut(6, "Map", new(NestedArray_As_Map))
			c := sch.NewContainerType(r, sch.ArrayKind)
			{
				r := c.Elem
				r.SetType(m.GetByName("Map"))
			}
		}
		{
			r := c.MustPut(7, "Set", new(NestedArray_As_Set))
			c := sch.NewContainerType(r, sch.ArrayKind)
			{
				r := c.Elem
				r.SetType(m.GetByName("Set"))
			}
		}
	}
	{
		r := m.MustPut(24, "NestedMap", new(NestedMap))
		c := sch.NewCompositeType(r, sch.UnionKind, new(NestedMap))
		{
			r := c.MustPut(1, "Enum", new(NestedMap_As_Enum))
			c := sch.NewContainerType(r, sch.MapKind)
			{
				r := c.Key
				r.SetType(m.GetByName("Type"))
			}
			{
				r := c.Elem
				r.SetType(m.GetByName("Enum"))
			}
		}
		{
			r := c.MustPut(2, "Union", new(NestedMap_As_Union))
			c := sch.NewContainerType(r, sch.MapKind)
			{
				r := c.Key
				r.SetType(m.GetByName("Type"))
			}
			{
				r := c.Elem
				r.SetType(m.GetByName("Union"))
			}
		}
		{
			r := c.MustPut(3, "Bitfield", new(NestedMap_As_Bitfield))
			c := sch.NewContainerType(r, sch.MapKind)
			{
				r := c.Key
				r.SetType(m.GetByName("Type"))
			}
			{
				r := c.Elem
				r.SetType(m.GetByName("Bitfield"))
			}
		}
		{
			r := c.MustPut(4, "Struct", new(NestedMap_As_Struct))
			c := sch.NewContainerType(r, sch.MapKind)
			{
				r := c.Key
				r.SetType(m.GetByName("Type"))
			}
			{
				r := c.Elem
				r.SetType(m.GetByName("Struct"))
			}
		}
		{
			r := c.MustPut(5, "Array", new(NestedMap_As_Array))
			c := sch.NewContainerType(r, sch.MapKind)
			{
				r := c.Key
				r.SetType(m.GetByName("Type"))
			}
			{
				r := c.Elem
				r.SetType(m.GetByName("Array"))
			}
		}
		{
			r := c.MustPut(6, "Map", new(NestedMap_As_Map))
			c := sch.NewContainerType(r, sch.MapKind)
			{
				r := c.Key
				r.SetType(m.GetByName("Type"))
			}
			{
				r := c.Elem
				r.SetType(m.GetByName("Map"))
			}
		}
		{
			r := c.MustPut(7, "Set", new(NestedMap_As_Set))
			c := sch.NewContainerType(r, sch.MapKind)
			{
				r := c.Key
				r.SetType(m.GetByName("Type"))
			}
			{
				r := c.Elem
				r.SetType(m.GetByName("Set"))
			}
		}
	}
	err := m.Check()
	if err != nil {
		panic(err)
	}
	return m
}()
