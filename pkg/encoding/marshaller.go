package encoding

import (
	"errors"
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
	types *sch.CompositeType
}

// NewMarshaller ...
func NewMarshaller(d sch.UnionDef) *Marshaller {
	c, err := sch.Compile(d)
	if err != nil {
		panic(err)
	}
	err = sch.Check(c)
	if err != nil {
		panic(err)
	}
	return &Marshaller{c}
}

// Marshal ...
func (m *Marshaller) Marshal(d Data, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr &&
		rv.Elem().Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	t := m.types.GetByRType(rv.Type())
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
	t := m.types.GetByRType(rv.Type())
	return unmarshal(d, rv, t)
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
	d.SetKey().EncodeSize(f.Tag)
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
	f := t.GetByTag(tag)
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
	d.SetKey().EncodeSize(f.Tag)
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
	f := t.GetByTag(tag)
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
			d.SetKey().EncodeSize(f.Tag)
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
		f := t.GetByTag(tag)
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
		f := t.GetByTag(tag)
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
