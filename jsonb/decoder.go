package jsonb

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/si3nloong/sqlike/core"
	"github.com/si3nloong/sqlike/reflext"
	"golang.org/x/xerrors"
)

// Decoder :
type Decoder struct {
	registry *Registry
}

// SetDecoders :
func (dec Decoder) SetDecoders(rg *Registry) {
	rg.SetTypeDecoder(reflect.TypeOf([]byte{}), dec.DecodeByte)
	rg.SetTypeDecoder(reflect.TypeOf(time.Time{}), dec.DecodeTime)
	rg.SetTypeDecoder(reflect.TypeOf(json.RawMessage{}), dec.DecodeJSONRaw)
	rg.SetKindDecoder(reflect.String, dec.DecodeString)
	rg.SetKindDecoder(reflect.Bool, dec.DecodeBool)
	rg.SetKindDecoder(reflect.Int, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Int8, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Int16, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Int32, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Int64, dec.DecodeInt)
	rg.SetKindDecoder(reflect.Uint, dec.DecodeUint)
	rg.SetKindDecoder(reflect.Uint8, dec.DecodeUint)
	rg.SetKindDecoder(reflect.Uint16, dec.DecodeUint)
	rg.SetKindDecoder(reflect.Uint32, dec.DecodeUint)
	rg.SetKindDecoder(reflect.Uint64, dec.DecodeUint)
	rg.SetKindDecoder(reflect.Float32, dec.DecodeFloat)
	rg.SetKindDecoder(reflect.Float64, dec.DecodeFloat)
	rg.SetKindDecoder(reflect.Struct, dec.DecodeStruct)
	rg.SetKindDecoder(reflect.Array, dec.DecodeArray)
	rg.SetKindDecoder(reflect.Slice, dec.DecodeArray)
	rg.SetKindDecoder(reflect.Interface, dec.DecodeInterface)
	dec.registry = rg
}

// DecodeByte :
func (dec Decoder) DecodeByte(r *Reader, v reflect.Value) error {
	x, err := r.ReadString()
	if err != nil {
		return err
	}
	var b []byte
	if x != "" {
		b, err = base64.StdEncoding.DecodeString(x)
		if err != nil {
			return err
		}
	}
	v.SetBytes(b)
	return nil
}

// DecodeTime :
func (dec Decoder) DecodeTime(r *Reader, v reflect.Value) error {
	b, err := r.ReadBytes()
	if err != nil {
		return err
	}
	if string(b) == null {
		v.Set(reflect.ValueOf(time.Time{}))
		return nil
	}
	x, err := time.Parse(`"`+time.RFC3339Nano+`"`, string(b))
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(x))
	return nil
}

// DecodeJSONRaw :
func (dec Decoder) DecodeJSONRaw(r *Reader, v reflect.Value) error {
	v.SetBytes(r.Bytes())
	return nil
}

// DecodeString :
func (dec Decoder) DecodeString(r *Reader, v reflect.Value) error {
	x, err := r.ReadEscapeString()
	if err != nil {
		return err
	}
	v.SetString(x)
	return nil
}

// DecodeBool :
func (dec Decoder) DecodeBool(r *Reader, v reflect.Value) error {
	x, err := r.ReadBoolean()
	if err != nil {
		return err
	}
	v.SetBool(x)
	return nil
}

// DecodeInt :
func (dec Decoder) DecodeInt(r *Reader, v reflect.Value) error {
	x, err := r.ReadNumber()
	if err != nil {
		return err
	}
	if v.OverflowInt(x) {
		return xerrors.New("integer overflow")
	}
	v.SetInt(x)
	return nil
}

// DecodeUint :
func (dec Decoder) DecodeUint(r *Reader, v reflect.Value) error {
	x, err := r.ReadNumber()
	if err != nil {
		return err
	}
	if x < 0 {
		return xerrors.New("number is not unsigned")
	}
	num := uint64(x)
	if v.OverflowUint(num) {
		return xerrors.New("unsigned integer overflow")
	}
	v.SetUint(num)
	return nil
}

// DecodeFloat :
func (dec Decoder) DecodeFloat(r *Reader, v reflect.Value) error {
	x, err := strconv.ParseFloat(string(r.Bytes()), 64)
	if err != nil {
		return err
	}
	if v.OverflowFloat(x) {
		return xerrors.New("float overflow")
	}
	v.SetFloat(x)
	return nil
}

// DecodeStruct :
func (dec *Decoder) DecodeStruct(r *Reader, v reflect.Value) error {
	mapper := core.DefaultMapper
	if r.IsNull() {
		v.Set(reflect.Zero(v.Type()))
		return r.skipNull()
	}

	return r.ReadFlattenObject(func(it *Reader, k string) error {
		vv, exists := mapper.LookUpFieldByName(v, k)
		if !exists {
			return nil
		}
		decoder, err := dec.registry.LookupDecoder(vv.Type())
		if err != nil {
			return err
		}
		return decoder(it, vv)
	})
}

// DecodeArray :
func (dec *Decoder) DecodeArray(r *Reader, v reflect.Value) error {
	t := v.Type()
	if r.IsNull() {
		v.Set(reflect.Zero(t))
		return r.skipNull()
	}

	v.Set(reflect.MakeSlice(t, 0, 0))
	t = t.Elem()
	return r.ReadArray(func(it *Reader) error {
		v.Set(reflect.Append(v, reflext.Zero(t)))
		vv := v.Index(v.Len() - 1)
		decoder, err := dec.registry.LookupDecoder(t)
		if err != nil {
			return err
		}
		return decoder(it, vv)
	})
}

// DecodeInterface :
func (dec Decoder) DecodeInterface(r *Reader, v reflect.Value) error {
	x, err := r.ReadValue()
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(x))
	return nil
}
