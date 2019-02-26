package amino

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type Typ3 uint8
type Typ4 uint8

const (
	Typ3_Varint	= Typ3(0)
	Typ3_8Byte	= Typ3(1)
	Typ3_ByteLength	= Typ3(2)
	Typ3_Struct	= Typ3(3)
	Typ3_StructTerm	= Typ3(4)
	Typ3_4Byte	= Typ3(5)
	Typ3_List	= Typ3(6)
	Typ3_Interface	= Typ3(7)

	Typ4_Pointer	= Typ4(0x08)
)

func (typ Typ3) String() string {
	switch typ {
	case Typ3_Varint:
		return "Varint"
	case Typ3_8Byte:
		return "8Byte"
	case Typ3_ByteLength:
		return "ByteLength"
	case Typ3_Struct:
		return "Struct"
	case Typ3_StructTerm:
		return "StructTerm"
	case Typ3_4Byte:
		return "4Byte"
	case Typ3_List:
		return "List"
	case Typ3_Interface:
		return "Interface"
	default:
		return fmt.Sprintf("<Invalid Typ3 %X>", byte(typ))
	}
}

func (typ Typ4) Typ3() Typ3		{ return Typ3(typ & 0x07) }
func (typ Typ4) IsPointer() bool	{ return (typ & 0x08) > 0 }
func (typ Typ4) String() string {
	if typ&0xF0 != 0 {
		return fmt.Sprintf("<Invalid Typ4 %X>", byte(typ))
	}
	if typ&0x08 != 0 {
		return "*" + Typ3(typ&0x07).String()
	} else {
		return Typ3(typ).String()
	}
}

func (cdc *Codec) MarshalBinary(o interface{}) ([]byte, error) {

	var buf = new(bytes.Buffer)

	bz, err := cdc.MarshalBinaryBare(o)
	if err != nil {
		return nil, err
	}

	err = EncodeUvarint(buf, uint64(len(bz)))
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(bz)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (cdc *Codec) MarshalBinaryWriter(w io.Writer, o interface{}) (n int64, err error) {
	var bz, _n = []byte(nil), int(0)
	bz, err = cdc.MarshalBinary(o)
	if err != nil {
		return 0, err
	}
	_n, err = w.Write(bz)
	n = int64(_n)
	return
}

func (cdc *Codec) MustMarshalBinary(o interface{}) []byte {
	bz, err := cdc.MarshalBinary(o)
	if err != nil {
		panic(err)
	}
	return bz
}

func (cdc *Codec) MarshalBinaryBare(o interface{}) ([]byte, error) {

	var rv, _, isNilPtr = derefPointers(reflect.ValueOf(o))
	if isNilPtr {

		panic("MarshalBinary cannot marshal a nil pointer directly. Try wrapping in a struct?")
	}

	var bz []byte
	buf := new(bytes.Buffer)
	rt := rv.Type()
	info, err := cdc.getTypeInfo_wlock(rt)
	if err != nil {
		return nil, err
	}
	err = cdc.encodeReflectBinary(buf, info, rv, FieldOptions{})
	if err != nil {
		return nil, err
	}
	bz = buf.Bytes()

	return bz, nil
}

func (cdc *Codec) MustMarshalBinaryBare(o interface{}) []byte {
	bz, err := cdc.MarshalBinaryBare(o)
	if err != nil {
		panic(err)
	}
	return bz
}

func (cdc *Codec) UnmarshalBinary(bz []byte, ptr interface{}) error {
	if len(bz) == 0 {
		return errors.New("UnmarshalBinary cannot decode empty bytes")
	}

	u64, n := binary.Uvarint(bz)
	if n < 0 {
		return fmt.Errorf("Error reading msg byte-length prefix: got code %v", n)
	}
	if u64 > uint64(len(bz)-n) {
		return fmt.Errorf("Not enough bytes to read in UnmarshalBinary, want %v more bytes but only have %v",
			u64, len(bz)-n)
	} else if u64 < uint64(len(bz)-n) {
		return fmt.Errorf("Bytes left over in UnmarshalBinary, should read %v more bytes but have %v",
			u64, len(bz)-n)
	}
	bz = bz[n:]

	return cdc.UnmarshalBinaryBare(bz, ptr)
}

func (cdc *Codec) UnmarshalBinaryReader(r io.Reader, ptr interface{}, maxSize int64) (n int64, err error) {
	if maxSize < 0 {
		panic("maxSize cannot be negative.")
	}

	var l int64
	var buf [binary.MaxVarintLen64]byte
	for i := 0; i < len(buf); i++ {
		_, err = r.Read(buf[i : i+1])
		if err != nil {
			return
		}
		n += 1
		if buf[i]&0x80 == 0 {
			break
		}
		if n >= maxSize {
			err = fmt.Errorf("Read overflow, maxSize is %v but uvarint(length-prefix) is itself greater than maxSize.", maxSize)
		}
	}
	u64, _ := binary.Uvarint(buf[:])
	if err != nil {
		return
	}
	if maxSize > 0 {
		if uint64(maxSize) < u64 {
			err = fmt.Errorf("Read overflow, maxSize is %v but this amino binary object is %v bytes.", maxSize, u64)
			return
		}
		if (maxSize - n) < int64(u64) {
			err = fmt.Errorf("Read overflow, maxSize is %v but this length-prefixed amino binary object is %v+%v bytes.", maxSize, n, u64)
			return
		}
	}
	l = int64(u64)
	if l < 0 {
		err = fmt.Errorf("Read overflow, this implementation can't read this because, why would anyone have this much data? Hello from 2018.")
	}

	var bz = make([]byte, l, l)
	_, err = io.ReadFull(r, bz)
	if err != nil {
		return
	}
	n += l

	err = cdc.UnmarshalBinaryBare(bz, ptr)
	return
}

func (cdc *Codec) MustUnmarshalBinary(bz []byte, ptr interface{}) {
	err := cdc.UnmarshalBinary(bz, ptr)
	if err != nil {
		panic(err)
	}
}

func (cdc *Codec) UnmarshalBinaryBare(bz []byte, ptr interface{}) error {
	if len(bz) == 0 {
		return errors.New("UnmarshalBinaryBare cannot decode empty bytes")
	}

	rv, rt := reflect.ValueOf(ptr), reflect.TypeOf(ptr)
	if rv.Kind() != reflect.Ptr {
		panic("Unmarshal expects a pointer")
	}
	rv, rt = rv.Elem(), rt.Elem()
	info, err := cdc.getTypeInfo_wlock(rt)
	if err != nil {
		return err
	}
	n, err := cdc.decodeReflectBinary(bz, info, rv, FieldOptions{})
	if err != nil {
		return err
	}
	if n != len(bz) {
		return fmt.Errorf("Unmarshal didn't read all bytes. Expected to read %v, only read %v", len(bz), n)
	}
	return nil
}

func (cdc *Codec) MustUnmarshalBinaryBare(bz []byte, ptr interface{}) {
	err := cdc.UnmarshalBinaryBare(bz, ptr)
	if err != nil {
		panic(err)
	}
}

func (cdc *Codec) MarshalJSON(o interface{}) ([]byte, error) {
	rv := reflect.ValueOf(o)
	if rv.Kind() == reflect.Invalid {
		return []byte("null"), nil
	}
	rt := rv.Type()

	w := new(bytes.Buffer)
	info, err := cdc.getTypeInfo_wlock(rt)
	if err != nil {
		return nil, err
	}
	if err := cdc.encodeReflectJSON(w, info, rv, FieldOptions{}); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (cdc *Codec) UnmarshalJSON(bz []byte, ptr interface{}) error {
	if len(bz) == 0 {
		return errors.New("UnmarshalJSON cannot decode empty bytes")
	}

	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return errors.New("UnmarshalJSON expects a pointer")
	}

	rv = rv.Elem()
	rt := rv.Type()
	info, err := cdc.getTypeInfo_wlock(rt)
	if err != nil {
		return err
	}
	return cdc.decodeReflectJSON(bz, info, rv, FieldOptions{})
}

func (cdc *Codec) MarshalJSONIndent(o interface{}, prefix, indent string) ([]byte, error) {
	bz, err := cdc.MarshalJSON(o)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	err = json.Indent(&out, bz, prefix, indent)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

var (
	jsonMarshalerType	= reflect.TypeOf(new(json.Marshaler)).Elem()
	jsonUnmarshalerType	= reflect.TypeOf(new(json.Unmarshaler)).Elem()
	errorType		= reflect.TypeOf(new(error)).Elem()
)
