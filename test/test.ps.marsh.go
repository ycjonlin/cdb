package test

import (
  enc "github.com/ycjonlin/cdb/pkg/encoding"
  sch "github.com/ycjonlin/cdb/pkg/schema"
)

var TestMarshaller = enc.NewMarshaller(sch.UnionDef{
  {1, "Bool", Bool(false), &sch.BuiltIn{sch.BoolType}},
  {2, "Int", Int(0), &sch.BuiltIn{sch.IntType}},
  {3, "Uint", Uint(0), &sch.BuiltIn{sch.UintType}},
  {4, "Float", Float(0), &sch.BuiltIn{sch.FloatType}},
  {5, "String", String(""), &sch.BuiltIn{sch.StringType}},
  {6, "Bytes", make(Bytes, 0), &sch.BuiltIn{sch.BytesType}},
  {7, "Type", Type(0), sch.NameRef("Int")},
  {11, "Enum", Enum(0), sch.EnumDef{
    {1, "Option1"},
    {2, "Option2"},
  }},
  {12, "Union", new(Union), sch.UnionDef{
    {1, "Option1", new(Union_As_Option1), &sch.BuiltIn{sch.BytesType}},
    {2, "Option2", new(Union_As_Option2), sch.NameRef("Type")},
  }},
  {13, "Bitfield", new(Bitfield), sch.BitfieldDef{
    {1, "Field1"},
    {2, "Field2"},
  }},
  {14, "Struct", new(Struct), sch.StructDef{
    {1, "Field1", nil, &sch.BuiltIn{sch.BytesType}},
    {2, "Field2", nil, sch.NameRef("Type")},
  }},
  {15, "Array", make(Array, 0), &sch.ArrayDef{sch.NameRef("Type")}},
  {16, "Map", make(Map), &sch.MapDef{sch.NameRef("Type"), sch.NameRef("Type")}},
  {17, "Set", make(Set), &sch.SetDef{sch.NameRef("Type")}},
  {21, "SelfRef", new(SelfRef), sch.StructDef{
    {1, "Field", nil, sch.NameRef("SelfRef")},
  }},
  {22, "MutRef1", new(MutRef1), sch.StructDef{
    {1, "Field", nil, sch.NameRef("MutRef2")},
  }},
  {23, "MutRef2", new(MutRef2), sch.StructDef{
    {1, "Field", nil, sch.NameRef("MutRef1")},
  }},
  {31, "NestedUnion", new(NestedUnion), sch.UnionDef{
    {1, "Bool", new(NestedUnion_As_Bool), &sch.BuiltIn{sch.BoolType}},
    {2, "Int", new(NestedUnion_As_Int), &sch.BuiltIn{sch.IntType}},
    {3, "Uint", new(NestedUnion_As_Uint), &sch.BuiltIn{sch.UintType}},
    {4, "Float", new(NestedUnion_As_Float), &sch.BuiltIn{sch.FloatType}},
    {5, "String", new(NestedUnion_As_String), &sch.BuiltIn{sch.StringType}},
    {6, "Bytes", new(NestedUnion_As_Bytes), &sch.BuiltIn{sch.BytesType}},
    {11, "Enum", new(NestedUnion_As_Enum), sch.NameRef("Enum")},
    {12, "Union", new(NestedUnion_As_Union), sch.NameRef("Union")},
    {13, "Bitfield", new(NestedUnion_As_Bitfield), sch.NameRef("Bitfield")},
    {14, "Struct", new(NestedUnion_As_Struct), sch.NameRef("Struct")},
    {15, "Array", new(NestedUnion_As_Array), sch.NameRef("Array")},
    {16, "Map", new(NestedUnion_As_Map), sch.NameRef("Map")},
    {17, "Set", new(NestedUnion_As_Set), sch.NameRef("Set")},
  }},
  {32, "NestedStruct", new(NestedStruct), sch.StructDef{
    {1, "Bool", nil, &sch.BuiltIn{sch.BoolType}},
    {2, "Int", nil, &sch.BuiltIn{sch.IntType}},
    {3, "Uint", nil, &sch.BuiltIn{sch.UintType}},
    {4, "Float", nil, &sch.BuiltIn{sch.FloatType}},
    {5, "String", nil, &sch.BuiltIn{sch.StringType}},
    {6, "Bytes", nil, &sch.BuiltIn{sch.BytesType}},
    {11, "Enum", nil, sch.NameRef("Enum")},
    {12, "Union", nil, sch.NameRef("Union")},
    {13, "Bitfield", nil, sch.NameRef("Bitfield")},
    {14, "Struct", nil, sch.NameRef("Struct")},
    {15, "Array", nil, sch.NameRef("Array")},
    {16, "Map", nil, sch.NameRef("Map")},
    {17, "Set", nil, sch.NameRef("Set")},
  }},
  {33, "NestedArray", new(NestedArray), sch.UnionDef{
    {1, "Bool", new(NestedArray_As_Bool), &sch.ArrayDef{&sch.BuiltIn{sch.BoolType}}},
    {2, "Int", new(NestedArray_As_Int), &sch.ArrayDef{&sch.BuiltIn{sch.IntType}}},
    {3, "Uint", new(NestedArray_As_Uint), &sch.ArrayDef{&sch.BuiltIn{sch.UintType}}},
    {4, "Float", new(NestedArray_As_Float), &sch.ArrayDef{&sch.BuiltIn{sch.FloatType}}},
    {5, "String", new(NestedArray_As_String), &sch.ArrayDef{&sch.BuiltIn{sch.StringType}}},
    {6, "Bytes", new(NestedArray_As_Bytes), &sch.ArrayDef{&sch.BuiltIn{sch.BytesType}}},
    {11, "Enum", new(NestedArray_As_Enum), &sch.ArrayDef{sch.NameRef("Enum")}},
    {12, "Union", new(NestedArray_As_Union), &sch.ArrayDef{sch.NameRef("Union")}},
    {13, "Bitfield", new(NestedArray_As_Bitfield), &sch.ArrayDef{sch.NameRef("Bitfield")}},
    {14, "Struct", new(NestedArray_As_Struct), &sch.ArrayDef{sch.NameRef("Struct")}},
    {15, "Array", new(NestedArray_As_Array), &sch.ArrayDef{sch.NameRef("Array")}},
    {16, "Map", new(NestedArray_As_Map), &sch.ArrayDef{sch.NameRef("Map")}},
    {17, "Set", new(NestedArray_As_Set), &sch.ArrayDef{sch.NameRef("Set")}},
  }},
  {34, "NestedMap", new(NestedMap), sch.UnionDef{
    {1, "Bool", new(NestedMap_As_Bool), &sch.MapDef{sch.NameRef("Type"), &sch.BuiltIn{sch.BoolType}}},
    {2, "Int", new(NestedMap_As_Int), &sch.MapDef{sch.NameRef("Type"), &sch.BuiltIn{sch.IntType}}},
    {3, "Uint", new(NestedMap_As_Uint), &sch.MapDef{sch.NameRef("Type"), &sch.BuiltIn{sch.UintType}}},
    {4, "Float", new(NestedMap_As_Float), &sch.MapDef{sch.NameRef("Type"), &sch.BuiltIn{sch.FloatType}}},
    {5, "String", new(NestedMap_As_String), &sch.MapDef{sch.NameRef("Type"), &sch.BuiltIn{sch.StringType}}},
    {6, "Bytes", new(NestedMap_As_Bytes), &sch.MapDef{sch.NameRef("Type"), &sch.BuiltIn{sch.BytesType}}},
    {11, "Enum", new(NestedMap_As_Enum), &sch.MapDef{sch.NameRef("Type"), sch.NameRef("Enum")}},
    {12, "Union", new(NestedMap_As_Union), &sch.MapDef{sch.NameRef("Type"), sch.NameRef("Union")}},
    {13, "Bitfield", new(NestedMap_As_Bitfield), &sch.MapDef{sch.NameRef("Type"), sch.NameRef("Bitfield")}},
    {14, "Struct", new(NestedMap_As_Struct), &sch.MapDef{sch.NameRef("Type"), sch.NameRef("Struct")}},
    {15, "Array", new(NestedMap_As_Array), &sch.MapDef{sch.NameRef("Type"), sch.NameRef("Array")}},
    {16, "Map", new(NestedMap_As_Map), &sch.MapDef{sch.NameRef("Type"), sch.NameRef("Map")}},
    {17, "Set", new(NestedMap_As_Set), &sch.MapDef{sch.NameRef("Type"), sch.NameRef("Set")}},
  }},
})
