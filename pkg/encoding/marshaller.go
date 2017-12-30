package encoding

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	sch "github.com/ycjonlin/cdb/pkg/schema"
)

// Err ...
var (
	ErrUnknownTag   = errors.New("unknown tag")
	ErrUnknownType  = errors.New("unknown type")
	ErrUnknownRType = errors.New("unknown rtype")
)

// Marshaller ...
type Marshaller struct {
	*sch.CompositeType
}

// NewMarshaller ...
func NewMarshaller() *Marshaller {
	return &Marshaller{
		sch.NewCompositeType(nil, sch.UnionKind, nil),
	}
}

// Check ...
func (m *Marshaller) Check() error {
	r := sch.NewReferenceType(0, 0, "", m.RType, nil)
	err := r.SetType(m.CompositeType)
	if err != nil {
		return nil
	}
	return check(r, m.RType)
}

// Marshal ...
func (m *Marshaller) Marshal(d Data, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr &&
		rv.Elem().Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	t := m.GetByRType(rv.Type())
	return marshal(d, rv, t)
}

// Unmarshal ...
func (m *Marshaller) Unmarshal(d Data, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr &&
		rv.Elem().Kind() != reflect.Array &&
		rv.Elem().Kind() != reflect.Struct {
		rv = rv.Elem()
	}
	t := m.GetByRType(rv.Type())
	return unmarshal(d, rv, t)
}

func check(r *sch.ReferenceType, rt reflect.Type) error {
	switch t := r.Base.(type) {
	case *sch.PrimitiveType:
		return checkPrimitiveType(t, rt)
	case *sch.CompositeType:
		return checkCompositeType(t, rt)
	case *sch.ContainerType:
		return checkContainerType(t, rt)
	default:
		return fmt.Errorf("unknown type %v", reflect.TypeOf(t))
	}
}

func checkPrimitiveType(t *sch.PrimitiveType, rt reflect.Type) error {
	switch t {
	case sch.BoolType:
		if rt.Kind() != reflect.Bool {
			return fmt.Errorf("Bool.Kind must be reflect.Bool, but got %v (%v)", rt.Kind(), rt)
		}
		return nil
	case sch.IntType:
		if rt.Kind() != reflect.Int64 {
			return fmt.Errorf("Int.Kind must be reflect.Int64, but got %v (%v)", rt.Kind(), rt)
		}
		return nil
	case sch.UintType:
		if rt.Kind() != reflect.Uint64 {
			return fmt.Errorf("Uint.Kind must be reflect.Uint64, but got %v (%v)", rt.Kind(), rt)
		}
		return nil
	case sch.FloatType:
		if rt.Kind() != reflect.Float64 {
			return fmt.Errorf("Float.Kind must be reflect.Float64, but got %v (%v)", rt.Kind(), rt)
		}
		return nil
	case sch.StringType:
		if rt.Kind() != reflect.String {
			return fmt.Errorf("String.Kind must be reflect.String, but got %v (%v)", rt.Kind(), rt)
		}
		return nil
	case sch.BytesType:
		if rt.Kind() != reflect.Slice {
			return fmt.Errorf("Bytes.Kind must be reflect.Slice, but got %v (%v)", rt.Kind(), rt)
		}
		if rt.Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("Bytes.Enum.Kind must be reflect.Uint8, but got %v (%v)", rt.Elem().Kind(), rt.Elem())
		}
		return nil
	default:
		return fmt.Errorf("unknown primitive type")
	}
}

func checkCompositeType(t *sch.CompositeType, rt reflect.Type) error {
	switch t.Kind {
	case sch.EnumKind:
		return checkEnumType(t, rt)
	case sch.UnionKind:
		return checkUnionType(t, rt)
	case sch.BitfieldKind:
		return checkBitfieldType(t, rt)
	case sch.StructKind:
		return checkStructType(t, rt)
	default:
		panic("unknown composite type")
	}
}

func checkEnumType(t *sch.CompositeType, rt reflect.Type) error {
	if rt.Kind() != reflect.Int {
		return fmt.Errorf("Enum.Kind must be reflect.Int, but got %v (%v)", rt.Kind(), rt)
	}
	return nil
}

func checkUnionType(t *sch.CompositeType, rt reflect.Type) error {
	if rt != nil && rt.Kind() != reflect.Interface {
		return fmt.Errorf("Union.Kind must be reflect.Interface, but got %v (%v)", rt.Kind(), rt)
	}
	for _, f := range t.Fields {
		rf := f.RType
		if rt != nil {
			if !rf.Implements(rt) {
				return fmt.Errorf("Union_Option must be implement Union, but %v does not implement %v", rf, rt)
			}
			if rf.Kind() != reflect.Ptr {
				return fmt.Errorf("Struct.Kind must be reflect.Ptr, but got %v (%v)", rf.Kind(), rf)
			}
			if rf.Elem().Kind() != reflect.Struct {
				return fmt.Errorf("Union_Option.Kind must be reflect.Struct, but got %v (%v)", rf.Elem().Kind(), rf.Elem())
			}
			if rf.Elem().NumField() != 1 {
				return fmt.Errorf("Struct.Elem.NumField must be 1, but got %d (%v)", rf.Elem().NumField(), rf.Elem())
			}
			rf = rf.Elem().Field(0).Type
		}
		err := check(f, rf)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkBitfieldType(t *sch.CompositeType, rt reflect.Type) error {
	if rt.Kind() != reflect.Ptr {
		return fmt.Errorf("Bitfield.Kind must be reflect.Ptr, but got %v (%v)", rt.Kind(), rt)
	}
	if rt.Elem().Kind() != reflect.Array {
		return fmt.Errorf("Bitfield.Elem.Kind must be reflect.Array, but got %v (%v)", rt.Elem().Kind(), rt.Elem())
	}
	if rt.Elem().Elem().Kind() != reflect.Uint8 {
		return fmt.Errorf("Bitfield.Elem.Kind.Elem must be reflect.Uint8, but got %v (%v)", rt.Elem().Elem().Kind(), rt.Elem().Elem())
	}
	if rt.Elem().Len() != (len(t.Fields)+7)/8 {
		return fmt.Errorf("Bitfield.Elem.Len must be (len(Fields)+7)/8 = %d, but got %d (%v)", (len(t.Fields)+7)/8, rt.Elem().Len(), rt.Elem())
	}
	return nil
}

func checkStructType(t *sch.CompositeType, rt reflect.Type) error {
	if rt.Kind() != reflect.Ptr {
		return fmt.Errorf("Struct.Kind must be reflect.Ptr, but got %v (%v)", rt.Kind(), rt)
	}
	if rt.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Struct.Elem.Kind must be reflect.Struct, but got %v (%v)", rt.Elem().Kind(), rt.Elem())
	}
	if rt.Elem().NumField() != len(t.Fields) {
		return fmt.Errorf("Struct.Elem.NumField must be len(Fields) = %d, but got %d (%v)", len(t.Fields), rt.Elem().NumField(), rt.Elem())
	}
	for _, f := range t.Fields {
		rf := rt.Elem().Field(f.Index).Type
		err := check(f, rf)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkContainerType(t *sch.ContainerType, rt reflect.Type) error {
	switch t.Kind {
	case sch.ArrayKind:
		return checkArrayType(t, rt)
	case sch.MapKind:
		return checkMapType(t, rt)
	case sch.SetKind:
		return checkSetType(t, rt)
	default:
		panic("unknown composite type")
	}
}

func checkArrayType(t *sch.ContainerType, rt reflect.Type) error {
	if rt.Kind() != reflect.Slice {
		return fmt.Errorf("Array.Kind must be reflect.Slice, but got %v (%v)", rt.Kind(), rt)
	}
	err := check(t.Elem, rt.Elem())
	if err != nil {
		return err
	}
	return nil
}

func checkMapType(t *sch.ContainerType, rt reflect.Type) error {
	if rt.Kind() != reflect.Map {
		return fmt.Errorf("Map.Kind must be reflect.Map, but got %v (%v)", rt.Kind(), rt)
	}
	err := check(t.Key, rt.Key())
	if err != nil {
		return err
	}
	err = check(t.Elem, rt.Elem())
	if err != nil {
		return err
	}
	return nil
}

func checkSetType(t *sch.ContainerType, rt reflect.Type) error {
	if rt.Kind() != reflect.Map {
		return fmt.Errorf("Set.Kind must be reflect.Map, but got %v (%v)", rt.Kind(), rt)
	}
	if rt.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Set.Elem.Kind must be reflect.Struct, but got %v (%v)", rt.Elem().Kind(), rt.Elem())
	}
	if rt.Elem().NumField() != 0 {
		return fmt.Errorf("Set.Elem.NumField must be 0, but got %d (%v)", rt.Elem().NumField(), rt.Elem())
	}
	err := check(t.Key, rt.Key())
	if err != nil {
		return err
	}
	return nil
}

func isZero(rv reflect.Value, r *sch.ReferenceType) bool {
	switch t := r.Base.(type) {
	case *sch.PrimitiveType:
		return isPrimitiveTypeZero(rv, t)
	case *sch.CompositeType:
		return isCompositeTypeZero(rv, t)
	case *sch.ContainerType:
		return isContainerTypeZero(rv, t)
	case nil:
		panic("nil type")
	default:
		panic("unknown type")
	}
}

func isPrimitiveTypeZero(rv reflect.Value, t *sch.PrimitiveType) bool {
	switch t {
	case sch.BoolType:
		return !rv.Bool()
	case sch.IntType:
		return rv.Int() == 0
	case sch.UintType:
		return rv.Uint() == 0
	case sch.FloatType:
		return rv.Float() == 0
	case sch.StringType:
		return rv.String() == ""
	case sch.BytesType:
		return rv.IsNil() || rv.Len() == 0
	default:
		panic("unknown primitive type")
	}
}

func isCompositeTypeZero(rv reflect.Value, t *sch.CompositeType) bool {
	switch t.Kind {
	case sch.EnumKind:
		return rv.Int() == 0
	case sch.UnionKind:
		return rv.IsNil()
	case sch.BitfieldKind:
		return rv.IsNil()
	case sch.StructKind:
		return rv.IsNil()
	default:
		panic("unknown composite type")
	}
}

func isContainerTypeZero(rv reflect.Value, t *sch.ContainerType) bool {
	switch t.Kind {
	case sch.ArrayKind:
		return rv.IsNil() || rv.Len() == 0
	case sch.MapKind:
		return rv.IsNil() || rv.Len() == 0
	case sch.SetKind:
		return rv.IsNil() || rv.Len() == 0
	default:
		panic("unknown composite type")
	}
}

func marshal(d Data, rv reflect.Value, t *sch.ReferenceType) error {
	if t == nil {
		return ErrUnknownType
	}
	switch t := t.Base.(type) {
	case *sch.PrimitiveType:
		return marshalPrimitiveType(d, rv, t)
	case *sch.CompositeType:
		return marshalCompositeType(d, rv, t)
	case *sch.ContainerType:
		return marshalContainerType(d, rv, t)
	default:
		return ErrUnknownType
	}
}

func unmarshal(d Data, rv reflect.Value, t *sch.ReferenceType) error {
	if t == nil {
		return ErrUnknownType
	}
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		rv.Set(reflect.New(rv.Type().Elem()))
	}
	switch t := t.Base.(type) {
	case *sch.PrimitiveType:
		return unmarshalPrimitiveType(d, rv, t)
	case *sch.CompositeType:
		return unmarshalCompositeType(d, rv, t)
	case *sch.ContainerType:
		return unmarshalContainerType(d, rv, t)
	default:
		return ErrUnknownType
	}
}

func marshalPrimitiveType(d Data, rv reflect.Value, t *sch.PrimitiveType) error {
	switch t {
	case sch.BoolType:
		d.Value().EncodeBool(rv.Bool())
		return nil
	case sch.IntType:
		d.Value().EncodeVarint(rv.Int())
		return nil
	case sch.UintType:
		d.Value().EncodeUvarint(rv.Uint())
		return nil
	case sch.FloatType:
		d.Value().EncodeVarfloat(rv.Float())
		return nil
	case sch.StringType:
		d.Value().EncodeNonsortingString(rv.String())
		return nil
	case sch.BytesType:
		d.Value().EncodeNonsortingBytes(rv.Bytes())
		return nil
	default:
		return ErrUnknownType
	}
}

func unmarshalPrimitiveType(d Data, rv reflect.Value, t *sch.PrimitiveType) error {
	switch t {
	case sch.BoolType:
		v, err := d.Value().DecodeBool()
		rv.SetBool(v)
		return err
	case sch.IntType:
		v, err := d.Value().DecodeVarint()
		rv.SetInt(v)
		return err
	case sch.UintType:
		v, err := d.Value().DecodeUvarint()
		rv.SetUint(v)
		return err
	case sch.FloatType:
		v, err := d.Value().DecodeVarfloat()
		rv.SetFloat(v)
		return err
	case sch.StringType:
		v, err := d.Value().DecodeNonsortingString()
		rv.SetString(v)
		return err
	case sch.BytesType:
		v, err := d.Value().DecodeNonsortingBytes(nil)
		rv.SetBytes(v)
		return err
	default:
		return ErrUnknownType
	}
}

func marshalCompositeType(d Data, rv reflect.Value, t *sch.CompositeType) error {
	switch t.Kind {
	case sch.EnumKind:
		return marshalEnumType(d, rv, t)
	case sch.UnionKind:
		return marshalUnionType(d, rv, t)
	case sch.BitfieldKind:
		return marshalBitfieldType(d, rv, t)
	case sch.StructKind:
		return marshalStructType(d, rv, t)
	default:
		return ErrUnknownType
	}
}

func unmarshalCompositeType(d Data, rv reflect.Value, t *sch.CompositeType) error {
	switch t.Kind {
	case sch.EnumKind:
		return unmarshalEnumType(d, rv, t)
	case sch.UnionKind:
		return unmarshalUnionType(d, rv, t)
	case sch.BitfieldKind:
		return unmarshalBitfieldType(d, rv, t)
	case sch.StructKind:
		return unmarshalStructType(d, rv, t)
	default:
		return ErrUnknownType
	}
}

func marshalContainerType(d Data, rv reflect.Value, t *sch.ContainerType) error {
	switch t.Kind {
	case sch.ArrayKind:
		return marshalArrayType(d, rv, t)
	case sch.MapKind:
		return marshalMapType(d, rv, t)
	case sch.SetKind:
		return marshalSetType(d, rv, t)
	default:
		return ErrUnknownType
	}
}

func unmarshalContainerType(d Data, rv reflect.Value, t *sch.ContainerType) error {
	switch t.Kind {
	case sch.ArrayKind:
		return unmarshalArrayType(d, rv, t)
	case sch.MapKind:
		return unmarshalMapType(d, rv, t)
	case sch.SetKind:
		return unmarshalSetType(d, rv, t)
	default:
		return ErrUnknownType
	}
}

func marshalEnumType(d Data, rv reflect.Value, t *sch.CompositeType) error {
	v := rv.Int()
	if v < 0 || v > math.MaxInt32 {
		return ErrExceededSizeLimit
	}
	if v == 0 {
		d.SetKey().EncodeSize(0)
		d.Value()
		d.Value()
		return nil
	}
	f := t.GetByIndex(int(v) - 1)
	if f == nil {
		return ErrUnknownTag
	}
	tag := int(f.Tag)
	d.SetKey().EncodeSize(tag)
	d.Value()
	d.Value()
	return nil
}

func unmarshalEnumType(d Data, rv reflect.Value, t *sch.CompositeType) error {
	tag, err := d.GetKey().DecodeSize()
	if err != nil {
		return err
	}
	if tag == 0 {
		d.Value()
		d.Value()
		rv.SetInt(0)
		return nil
	}
	f := t.GetByTag(sch.Tag(tag))
	if f == nil {
		return ErrUnknownTag
	}
	d.Value()
	d.Value()
	v := int64(f.Index + 1)
	rv.SetInt(v)
	return nil
}

func marshalUnionType(d Data, rv reflect.Value, t *sch.CompositeType) error {
	if rv.IsNil() {
		d.SetKey().EncodeSize(0)
		d.Value()
		d.Value()
		return nil
	}
	rv = rv.Elem()
	f := t.GetByRType(rv.Type())
	if f == nil {
		return ErrUnknownRType
	}
	tag := int(f.Tag)
	d.SetKey().EncodeSize(tag)
	rf := rv.Elem().Field(0)
	err := marshal(d, rf, f)
	if err != nil {
		return err
	}
	d.Value()
	return nil
}

func unmarshalUnionType(d Data, rv reflect.Value, t *sch.CompositeType) error {
	tag, err := d.GetKey().DecodeSize()
	if err != nil {
		return err
	}
	if tag == 0 {
		d.Value()
		d.Value()
		rv.Set(reflect.Zero(rv.Type()))
		return nil
	}
	f := t.GetByTag(sch.Tag(tag))
	if f == nil {
		return ErrUnknownTag
	}
	rvo := reflect.New(f.RType.Elem())
	rf := rvo.Elem().Field(0)
	err = unmarshal(d, rf, f)
	if err != nil {
		return err
	}
	d.Value()
	rv.Set(rvo)
	return nil
}

func marshalBitfieldType(d Data, rv reflect.Value, t *sch.CompositeType) error {
	rv = rv.Elem()
	b := rv.Slice(0, rv.Len()).Bytes()
	for _, f := range t.Fields {
		i := f.Index
		if b[i>>3]&(1<<uint(i&7)) != 0 {
			tag := int(f.Tag)
			d.SetKey().EncodeSize(tag)
			d.Value()
		}
	}
	if b := d.End(); b != nil {
		b.EncodeSize(0)
	}
	d.Value()
	return nil
}

func unmarshalBitfieldType(d Data, rv reflect.Value, t *sch.CompositeType) error {
	rv = rv.Elem()
	b := rv.Slice(0, rv.Len()).Bytes()
	for d.Next() {
		tag, err := d.GetKey().DecodeSize()
		if err != nil {
			return err
		}
		if tag == 0 {
			d.Value()
			break
		}
		f := t.GetByTag(sch.Tag(tag))
		if f == nil {
			return ErrUnknownTag
		}
		d.Value()
		i := f.Index
		b[i>>3] |= (1 << uint(i&7))
	}
	d.Value()
	return nil
}

func marshalStructType(d Data, rv reflect.Value, t *sch.CompositeType) error {
	rv = rv.Elem()
	for _, f := range t.Fields {
		rf := rv.Field(f.Index)
		if isZero(rf, f) {
			continue
		}
		tag := int(f.Tag)
		d.SetKey().EncodeSize(tag)
		err := marshal(d, rf, f)
		if err != nil {
			return err
		}
	}
	if b := d.End(); b != nil {
		b.EncodeSize(0)
	}
	d.Value()
	return nil
}

func unmarshalStructType(d Data, rv reflect.Value, t *sch.CompositeType) error {
	rv = rv.Elem()
	for d.Next() {
		tag, err := d.GetKey().DecodeSize()
		if err != nil {
			return err
		}
		if tag == 0 {
			d.Value()
			break
		}
		f := t.GetByTag(sch.Tag(tag))
		if f == nil {
			return ErrUnknownTag
		}
		rf := rv.Field(f.Index)
		err = unmarshal(d, rf, f)
		if err != nil {
			return err
		}
	}
	d.Value()
	return nil
}

func marshalArrayType(d Data, rv reflect.Value, t *sch.ContainerType) error {
	if b := d.Len(); b != nil {
		b.EncodeSize(rv.Len())
	}
	for k, n := 0, rv.Len(); k < n; k++ {
		re := rv.Index(k)
		d.SetKey().EncodeSize(k)
		err := marshal(d, re, t.Elem)
		if err != nil {
			return err
		}
	}
	d.Value()
	return nil
}

func unmarshalArrayType(d Data, rv reflect.Value, t *sch.ContainerType) error {
	n := MaxSize
	if b := d.Len(); b != nil {
		var err error
		n, err = b.DecodeSize()
		if err != nil {
			return err
		}
	}
	rt := rv.Type()
	rv.Set(reflect.MakeSlice(rt, 0, 0))
	for i := 0; i < n && d.Next(); i++ {
		k, err := d.GetKey().DecodeSize()
		if err != nil {
			return err
		}
		if rv.Cap() <= k {
			c := rv.Cap() * 2
			if c < k+1 {
				c = k + 1
			}
			rvp := reflect.MakeSlice(rt, rv.Len(), c)
			reflect.Copy(rvp, rv)
			rv.Set(rvp)
		}
		rv.Set(rv.Slice(0, k+1))
		re := rv.Index(k)
		err = unmarshal(d, re, t.Elem)
		if err != nil {
			return err
		}
	}
	d.Value()
	return nil
}

func marshalMapType(d Data, rv reflect.Value, t *sch.ContainerType) error {
	if b := d.Len(); b != nil {
		b.EncodeSize(rv.Len())
	}
	for _, rk := range rv.MapKeys() {
		re := rv.MapIndex(rk)
		dk := d.Data(d.SetKey())
		err := marshal(dk, rk, t.Key)
		if err != nil {
			return err
		}
		err = marshal(d, re, t.Elem)
		if err != nil {
			return err
		}
	}
	d.Value()
	return nil
}

func unmarshalMapType(d Data, rv reflect.Value, t *sch.ContainerType) error {
	n := MaxSize
	if b := d.Len(); b != nil {
		var err error
		n, err = b.DecodeSize()
		if err != nil {
			return err
		}
	}
	rt := rv.Type()
	rv.Set(reflect.MakeMap(rt))
	rk := reflect.New(rt.Key()).Elem()
	re := reflect.New(rt.Elem()).Elem()
	for i := 0; i < n && d.Next(); i++ {
		dk := d.Data(d.GetKey())
		err := unmarshal(dk, rk, t.Key)
		if err != nil {
			return err
		}
		err = unmarshal(d, re, t.Elem)
		if err != nil {
			return err
		}
		rv.SetMapIndex(rk, re)
	}
	d.Value()
	return nil
}

func marshalSetType(d Data, rv reflect.Value, t *sch.ContainerType) error {
	if b := d.Len(); b != nil {
		b.EncodeSize(rv.Len())
	}
	for _, rk := range rv.MapKeys() {
		dk := d.Data(d.SetKey())
		err := marshal(dk, rk, t.Key)
		if err != nil {
			return err
		}
		d.Value()
	}
	d.Value()
	return nil
}

func unmarshalSetType(d Data, rv reflect.Value, t *sch.ContainerType) error {
	n := MaxSize
	if b := d.Len(); b != nil {
		var err error
		n, err = b.DecodeSize()
		if err != nil {
			return err
		}
	}
	rt := rv.Type()
	rtk := rt.Key()
	rv.Set(reflect.MakeMap(rt))
	rk := reflect.New(rtk).Elem()
	re := reflect.ValueOf(struct{}{})
	for i := 0; i < n && d.Next(); i++ {
		dk := d.Data(d.GetKey())
		err := unmarshal(dk, rk, t.Key)
		if err != nil {
			return err
		}
		d.Value()
		rv.SetMapIndex(rk, re)
	}
	d.Value()
	return nil
}
