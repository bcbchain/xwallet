package amino

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func (cdc *Codec) encodeReflectBinary(w io.Writer, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if rv.Kind() == reflect.Ptr {
		panic("should not happen")
	}
	if !rv.IsValid() {
		panic("should not happen")
	}
	if printLog {
		spew.Printf("(E) encodeReflectBinary(info: %v, rv: %#v (%v), opts: %v)\n",
			info, rv.Interface(), rv.Type(), opts)
		defer func() {
			fmt.Printf("(E) -> err: %v\n", err)
		}()
	}

	if info.Registered {
		var typ = typeToTyp4(info.Type, opts).Typ3()
		_, err = w.Write(info.Prefix.WithTyp3(typ).Bytes())
		if err != nil {
			return
		}
	}

	err = cdc._encodeReflectBinary(w, info, rv, opts)
	return
}

func (cdc *Codec) _encodeReflectBinary(w io.Writer, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if rv.Kind() == reflect.Ptr {
		panic("should not happen")
	}
	if !rv.IsValid() {
		panic("should not happen")
	}
	if printLog {
		spew.Printf("(_) _encodeReflectBinary(info: %v, rv: %#v (%v), opts: %v)\n",
			info, rv.Interface(), rv.Type(), opts)
		defer func() {
			fmt.Printf("(_) -> err: %v\n", err)
		}()
	}

	if info.IsAminoMarshaler {

		var rrv, rinfo = reflect.Value{}, (*TypeInfo)(nil)
		rrv, err = toReprObject(rv)
		if err != nil {
			return
		}
		rinfo, err = cdc.getTypeInfo_wlock(info.AminoMarshalReprType)
		if err != nil {
			return
		}

		err = cdc._encodeReflectBinary(w, rinfo, rrv, opts)
		return
	}

	switch info.Type.Kind() {

	case reflect.Interface:
		err = cdc.encodeReflectBinaryInterface(w, info, rv, opts)

	case reflect.Array:
		if info.Type.Elem().Kind() == reflect.Uint8 {
			err = cdc.encodeReflectBinaryByteArray(w, info, rv, opts)
		} else {
			err = cdc.encodeReflectBinaryList(w, info, rv, opts)
		}

	case reflect.Slice:
		if info.Type.Elem().Kind() == reflect.Uint8 {
			err = cdc.encodeReflectBinaryByteSlice(w, info, rv, opts)
		} else {
			err = cdc.encodeReflectBinaryList(w, info, rv, opts)
		}

	case reflect.Struct:
		err = cdc.encodeReflectBinaryStruct(w, info, rv, opts)

	case reflect.Int64:
		if opts.BinVarint {
			err = EncodeVarint(w, rv.Int())
		} else {
			err = EncodeInt64(w, rv.Int())
		}

	case reflect.Int32:
		err = EncodeInt32(w, int32(rv.Int()))

	case reflect.Int16:
		err = EncodeInt16(w, int16(rv.Int()))

	case reflect.Int8:
		err = EncodeInt8(w, int8(rv.Int()))

	case reflect.Int:
		err = EncodeVarint(w, rv.Int())

	case reflect.Uint64:
		if opts.BinVarint {
			err = EncodeUvarint(w, rv.Uint())
		} else {
			err = EncodeUint64(w, rv.Uint())
		}

	case reflect.Uint32:
		err = EncodeUint32(w, uint32(rv.Uint()))

	case reflect.Uint16:
		err = EncodeUint16(w, uint16(rv.Uint()))

	case reflect.Uint8:
		err = EncodeUint8(w, uint8(rv.Uint()))

	case reflect.Uint:
		err = EncodeUvarint(w, rv.Uint())

	case reflect.Bool:
		err = EncodeBool(w, rv.Bool())

	case reflect.Float64:
		if !opts.Unsafe {
			err = errors.New("Amino float* support requires `amino:\"unsafe\"`.")
			return
		}
		err = EncodeFloat64(w, rv.Float())

	case reflect.Float32:
		if !opts.Unsafe {
			err = errors.New("Amino float* support requires `amino:\"unsafe\"`.")
			return
		}
		err = EncodeFloat32(w, float32(rv.Float()))

	case reflect.String:
		err = EncodeString(w, rv.String())

	default:
		panic(fmt.Sprintf("unsupported type %v", info.Type.Kind()))
	}

	return
}

func (cdc *Codec) encodeReflectBinaryInterface(w io.Writer, iinfo *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if printLog {
		fmt.Println("(e) encodeReflectBinaryInterface")
		defer func() {
			fmt.Printf("(e) -> err: %v\n", err)
		}()
	}

	if rv.IsNil() {
		_, err = w.Write([]byte{0x00, 0x00})
		return
	}

	var crv, isPtr, isNilPtr = derefPointers(rv.Elem())
	if isPtr && crv.Kind() == reflect.Interface {

		panic("should not happen")
	}
	if isNilPtr {
		panic(fmt.Sprintf("Illegal nil-pointer of type %v for registered interface %v. "+
			"For compatibility with other languages, nil-pointer interface values are forbidden.", crv.Type(), iinfo.Type))
	}
	var crt = crv.Type()

	var cinfo *TypeInfo
	cinfo, err = cdc.getTypeInfo_wlock(crt)
	if err != nil {
		return
	}
	if !cinfo.Registered {
		err = fmt.Errorf("Cannot encode unregistered concrete type %v.", crt)
		return
	}

	var needDisamb bool = false
	if iinfo.AlwaysDisambiguate {
		needDisamb = true
	} else if len(iinfo.Implementers[cinfo.Prefix]) > 1 {
		needDisamb = true
	}
	if needDisamb {
		_, err = w.Write(append([]byte{0x00}, cinfo.Disamb[:]...))
		if err != nil {
			return
		}
	}

	var typ = typeToTyp3(crt, opts)
	_, err = w.Write(cinfo.Prefix.WithTyp3(typ).Bytes())
	if err != nil {
		return
	}

	err = cdc._encodeReflectBinary(w, cinfo, crv, opts)
	return
}

func (cdc *Codec) encodeReflectBinaryByteArray(w io.Writer, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	ert := info.Type.Elem()
	if ert.Kind() != reflect.Uint8 {
		panic("should not happen")
	}
	length := info.Type.Len()

	var byteslice = []byte(nil)
	if rv.CanAddr() {
		byteslice = rv.Slice(0, length).Bytes()
	} else {
		byteslice = make([]byte, length)
		reflect.Copy(reflect.ValueOf(byteslice), rv)
	}

	err = EncodeByteSlice(w, byteslice)
	return
}

func (cdc *Codec) encodeReflectBinaryList(w io.Writer, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if printLog {
		fmt.Println("(e) encodeReflectBinaryList")
		defer func() {
			fmt.Printf("(e) -> err: %v\n", err)
		}()
	}
	ert := info.Type.Elem()
	if ert.Kind() == reflect.Uint8 {
		panic("should not happen")
	}

	var typ = typeToTyp4(ert, opts)
	err = EncodeByte(w, byte(typ))
	if err != nil {
		return
	}

	err = EncodeUvarint(w, uint64(rv.Len()))
	if err != nil {
		return
	}

	var einfo *TypeInfo
	einfo, err = cdc.getTypeInfo_wlock(ert)
	if err != nil {
		return
	}
	for i := 0; i < rv.Len(); i++ {

		var erv, void = isVoid(rv.Index(i))
		if typ.IsPointer() {

			if void {

				_, err = w.Write([]byte{0x01})
				continue
			} else {

				_, err = w.Write([]byte{0x00})
			}
		}

		err = cdc.encodeReflectBinary(w, einfo, erv, opts)
		if err != nil {
			return
		}
	}
	return
}

func (cdc *Codec) encodeReflectBinaryByteSlice(w io.Writer, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if printLog {
		fmt.Println("(e) encodeReflectBinaryByteSlice")
		defer func() {
			fmt.Printf("(e) -> err: %v\n", err)
		}()
	}
	ert := info.Type.Elem()
	if ert.Kind() != reflect.Uint8 {
		panic("should not happen")
	}

	var byteslice = rv.Bytes()
	err = EncodeByteSlice(w, byteslice)
	return
}

func (cdc *Codec) encodeReflectBinaryStruct(w io.Writer, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if printLog {
		fmt.Println("(e) encodeReflectBinaryBinaryStruct")
		defer func() {
			fmt.Printf("(e) -> err: %v\n", err)
		}()
	}

	switch info.Type {

	case timeType:

		err = EncodeTime(w, rv.Interface().(time.Time))
		return

	default:
		for _, field := range info.Fields {

			var frv, void = isVoid(rv.Field(field.Index))
			if void {

				continue
			}
			var finfo *TypeInfo
			finfo, err = cdc.getTypeInfo_wlock(field.Type)
			if err != nil {
				return
			}

			err = encodeFieldNumberAndTyp3(w, field.BinFieldNum, field.BinTyp3)
			if err != nil {
				return
			}

			err = cdc.encodeReflectBinary(w, finfo, frv, field.FieldOptions)
			if err != nil {
				return
			}
		}

		err = EncodeByte(w, byte(Typ3_StructTerm))
		if err != nil {
			return
		}
		return

	}

}

func encodeFieldNumberAndTyp3(w io.Writer, num uint32, typ Typ3) (err error) {
	if (typ & 0xF8) != 0 {
		panic(fmt.Sprintf("invalid Typ3 byte %X", typ))
	}
	if num < 0 || num > (1<<29-1) {
		panic(fmt.Sprintf("invalid field number %v", num))
	}

	var value64 = (uint64(num) << 3) | uint64(typ)

	var buf [10]byte
	n := binary.PutUvarint(buf[:], value64)
	_, err = w.Write(buf[0:n])
	return
}
