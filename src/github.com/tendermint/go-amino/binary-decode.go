package amino

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func (cdc *Codec) decodeReflectBinary(bz []byte, info *TypeInfo, rv reflect.Value, opts FieldOptions) (n int, err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if info.Type.Kind() == reflect.Interface && rv.Kind() == reflect.Ptr {
		panic("should not happen")
	}
	if printLog {
		spew.Printf("(D) decodeReflectBinary(bz: %X, info: %v, rv: %#v (%v), opts: %v)\n",
			bz, info, rv.Interface(), rv.Type(), opts)
		defer func() {
			fmt.Printf("(D) -> n: %v, err: %v\n", n, err)
		}()
	}

	if info.Registered {
		if len(bz) < PrefixBytesLen {
			err = errors.New("EOF skipping prefix bytes.")
			return
		}

		prefix3 := NewPrefixBytes(bz[:PrefixBytesLen])
		var prefix, typ = prefix3.SplitTyp3()
		if info.Prefix != prefix {
			panic("should not happen")
		}

		err = checkTyp3(info.Type, typ, opts)
		if err != nil {
			return
		}

		bz = bz[PrefixBytesLen:]
		n += PrefixBytesLen
	}

	_n := 0
	_n, err = cdc._decodeReflectBinary(bz, info, rv, opts)
	slide(&bz, &n, _n)
	return
}

func (cdc *Codec) _decodeReflectBinary(bz []byte, info *TypeInfo, rv reflect.Value, opts FieldOptions) (n int, err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if info.Type.Kind() == reflect.Interface && rv.Kind() == reflect.Ptr {
		panic("should not happen")
	}
	if printLog {
		spew.Printf("(_) _decodeReflectBinary(bz: %X, info: %v, rv: %#v (%v), opts: %v)\n",
			bz, info, rv.Interface(), rv.Type(), opts)
		defer func() {
			fmt.Printf("(_) -> n: %v, err: %v\n", n, err)
		}()
	}
	var _n int

	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			newPtr := reflect.New(rv.Type().Elem())
			rv.Set(newPtr)
		}
		rv = rv.Elem()
	}

	if info.IsAminoUnmarshaler {

		rrv, rinfo := reflect.New(info.AminoUnmarshalReprType).Elem(), (*TypeInfo)(nil)
		rinfo, err = cdc.getTypeInfo_wlock(info.AminoUnmarshalReprType)
		if err != nil {
			return
		}
		_n, err = cdc._decodeReflectBinary(bz, rinfo, rrv, opts)
		if slide(&bz, &n, _n) && err != nil {
			return
		}

		uwrm := rv.Addr().MethodByName("UnmarshalAmino")
		uwouts := uwrm.Call([]reflect.Value{rrv})
		err = uwouts[0].Interface().(error)
		return
	}

	switch info.Type.Kind() {

	case reflect.Interface:
		_n, err = cdc.decodeReflectBinaryInterface(bz, info, rv, opts)
		n += _n
		return

	case reflect.Array:
		ert := info.Type.Elem()
		if ert.Kind() == reflect.Uint8 {
			_n, err = cdc.decodeReflectBinaryByteArray(bz, info, rv, opts)
			n += _n
		} else {
			_n, err = cdc.decodeReflectBinaryArray(bz, info, rv, opts)
			n += _n
		}
		return

	case reflect.Slice:
		ert := info.Type.Elem()
		if ert.Kind() == reflect.Uint8 {
			_n, err = cdc.decodeReflectBinaryByteSlice(bz, info, rv, opts)
			n += _n
		} else {
			_n, err = cdc.decodeReflectBinarySlice(bz, info, rv, opts)
			n += _n
		}
		return

	case reflect.Struct:
		_n, err = cdc.decodeReflectBinaryStruct(bz, info, rv, opts)
		n += _n
		return

	case reflect.Int64:
		var num int64
		if opts.BinVarint {
			num, _n, err = DecodeVarint(bz)
			if slide(&bz, &n, _n) && err != nil {
				return
			}
			rv.SetInt(num)
		} else {
			num, _n, err = DecodeInt64(bz)
			if slide(&bz, &n, _n) && err != nil {
				return
			}
			rv.SetInt(num)
		}
		return

	case reflect.Int32:
		var num int32
		num, _n, err = DecodeInt32(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.SetInt(int64(num))
		return

	case reflect.Int16:
		var num int16
		num, _n, err = DecodeInt16(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.SetInt(int64(num))
		return

	case reflect.Int8:
		var num int8
		num, _n, err = DecodeInt8(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.SetInt(int64(num))
		return

	case reflect.Int:
		var num int64
		num, _n, err = DecodeVarint(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.SetInt(num)
		return

	case reflect.Uint64:
		var num uint64
		if opts.BinVarint {
			num, _n, err = DecodeUvarint(bz)
			if slide(&bz, &n, _n) && err != nil {
				return
			}
			rv.SetUint(num)
		} else {
			num, _n, err = DecodeUint64(bz)
			if slide(&bz, &n, _n) && err != nil {
				return
			}
			rv.SetUint(num)
		}
		return

	case reflect.Uint32:
		var num uint32
		num, _n, err = DecodeUint32(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.SetUint(uint64(num))
		return

	case reflect.Uint16:
		var num uint16
		num, _n, err = DecodeUint16(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.SetUint(uint64(num))
		return

	case reflect.Uint8:
		var num uint8
		num, _n, err = DecodeUint8(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.SetUint(uint64(num))
		return

	case reflect.Uint:
		var num uint64
		num, _n, err = DecodeUvarint(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.SetUint(num)
		return

	case reflect.Bool:
		var b bool
		b, _n, err = DecodeBool(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.SetBool(b)
		return

	case reflect.Float64:
		var f float64
		if !opts.Unsafe {
			err = errors.New("Float support requires `amino:\"unsafe\"`.")
			return
		}
		f, _n, err = DecodeFloat64(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.SetFloat(f)
		return

	case reflect.Float32:
		var f float32
		if !opts.Unsafe {
			err = errors.New("Float support requires `amino:\"unsafe\"`.")
			return
		}
		f, _n, err = DecodeFloat32(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.SetFloat(float64(f))
		return

	case reflect.String:
		var str string
		str, _n, err = DecodeString(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.SetString(str)
		return

	default:
		panic(fmt.Sprintf("unknown field type %v", info.Type.Kind()))
	}

}

func (cdc *Codec) decodeReflectBinaryInterface(bz []byte, iinfo *TypeInfo, rv reflect.Value, opts FieldOptions) (n int, err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if printLog {
		fmt.Println("(d) decodeReflectBinaryInterface")
		defer func() {
			fmt.Printf("(d) -> err: %v\n", err)
		}()
	}
	if !rv.IsNil() {

		err = errors.New("Decoding to a non-nil interface is not supported yet")
		return
	}

	disamb, hasDisamb, prefix, typ, hasPrefix, isNil, _n, err := DecodeDisambPrefixBytes(bz)
	if slide(&bz, &n, _n) && err != nil {
		return
	}

	if isNil {
		rv.Set(iinfo.ZeroValue)
		return
	}

	var cinfo *TypeInfo
	if hasDisamb {
		cinfo, err = cdc.getTypeInfoFromDisfix_rlock(toDisfix(disamb, prefix))
	} else if hasPrefix {
		cinfo, err = cdc.getTypeInfoFromPrefix_rlock(iinfo, prefix)
	} else {
		err = errors.New("Expected disambiguation or prefix bytes.")
	}
	if err != nil {
		return
	}

	err = checkTyp3(cinfo.Type, typ, opts)
	if err != nil {
		return
	}

	var crv, irvSet = constructConcreteType(cinfo)

	_n, err = cdc._decodeReflectBinary(bz, cinfo, crv, opts)
	if slide(&bz, &n, _n) && err != nil {
		rv.Set(irvSet)
		return
	}

	rv.Set(irvSet)
	return
}

func (cdc *Codec) decodeReflectBinaryByteArray(bz []byte, info *TypeInfo, rv reflect.Value, opts FieldOptions) (n int, err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if printLog {
		fmt.Println("(d) decodeReflectBinaryByteArray")
		defer func() {
			fmt.Printf("(d) -> err: %v\n", err)
		}()
	}
	ert := info.Type.Elem()
	if ert.Kind() != reflect.Uint8 {
		panic("should not happen")
	}
	length := info.Type.Len()
	if len(bz) < length {
		return 0, fmt.Errorf("Insufficient bytes to decode [%v]byte.", length)
	}

	var byteslice, _n = []byte(nil), int(0)
	byteslice, _n, err = DecodeByteSlice(bz)
	if slide(&bz, &n, _n) && err != nil {
		return
	}
	if len(byteslice) != length {
		err = fmt.Errorf("Mismatched byte array length: Expected %v, got %v",
			length, len(byteslice))
		return
	}

	reflect.Copy(rv, reflect.ValueOf(byteslice))
	return
}

func (cdc *Codec) decodeReflectBinaryArray(bz []byte, info *TypeInfo, rv reflect.Value, opts FieldOptions) (n int, err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if printLog {
		fmt.Println("(d) decodeReflectBinaryArray")
		defer func() {
			fmt.Printf("(d) -> err: %v\n", err)
		}()
	}
	ert := info.Type.Elem()
	if ert.Kind() == reflect.Uint8 {
		panic("should not happen")
	}
	length := info.Type.Len()
	einfo := (*TypeInfo)(nil)
	einfo, err = cdc.getTypeInfo_wlock(ert)
	if err != nil {
		return
	}

	var ptr, _n = false, int(0)
	ptr, _n, err = decodeTyp4AndCheck(ert, bz, opts)
	if slide(&bz, &n, _n) && err != nil {
		return
	}

	var count = uint64(0)
	count, _n, err = DecodeUvarint(bz)
	if slide(&bz, &n, _n) && err != nil {
		return
	}
	if int(count) != length {
		err = fmt.Errorf("Expected num items of %v, decoded %v", length, count)
		return
	}

	for i := 0; i < length; i++ {
		var erv, _n = rv.Index(i), int(0)

		if ptr {
			numNil := int64(0)
			numNil, _n, err = decodeNumNilBytes(bz)
			if slide(&bz, &n, _n) && err != nil {
				return
			}
			if numNil == 0 {

			} else if numNil == 1 {

				erv.Set(reflect.Zero(erv.Type()))
				continue
			} else {
				panic("should not happen")
			}
		}

		_n, err = cdc.decodeReflectBinary(bz, einfo, erv, opts)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
	}
	return
}

func (cdc *Codec) decodeReflectBinaryByteSlice(bz []byte, info *TypeInfo, rv reflect.Value, opts FieldOptions) (n int, err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if printLog {
		fmt.Println("(d) decodeReflectByteSlice")
		defer func() {
			fmt.Printf("(d) -> err: %v\n", err)
		}()
	}
	ert := info.Type.Elem()
	if ert.Kind() != reflect.Uint8 {
		panic("should not happen")
	}

	var byteslice, _n = []byte(nil), int(0)
	byteslice, _n, err = DecodeByteSlice(bz)
	if slide(&bz, &n, _n) && err != nil {
		return
	}
	if len(byteslice) == 0 {

		rv.Set(info.ZeroValue)
	} else {
		rv.Set(reflect.ValueOf(byteslice))
	}
	return
}

func (cdc *Codec) decodeReflectBinarySlice(bz []byte, info *TypeInfo, rv reflect.Value, opts FieldOptions) (n int, err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if printLog {
		fmt.Println("(d) decodeReflectBinarySlice")
		defer func() {
			fmt.Printf("(d) -> err: %v\n", err)
		}()
	}
	ert := info.Type.Elem()
	if ert.Kind() == reflect.Uint8 {
		panic("should not happen")
	}
	einfo := (*TypeInfo)(nil)
	einfo, err = cdc.getTypeInfo_wlock(ert)
	if err != nil {
		return
	}

	var ptr, _n = false, int(0)
	ptr, _n, err = decodeTyp4AndCheck(ert, bz, opts)
	if slide(&bz, &n, _n) && err != nil {
		return
	}

	var count = uint64(0)
	count, _n, err = DecodeUvarint(bz)
	if slide(&bz, &n, _n) && err != nil {
		return
	}
	if int(count) < 0 {
		err = fmt.Errorf("Impossible number of elements (%v)", count)
		return
	}
	if int(count) > len(bz) {
		err = fmt.Errorf("Impossible number of elements (%v) compared to buffer length (%v)",
			count, len(bz))
		return
	}

	if count == 0 {
		rv.Set(info.ZeroValue)
		return
	}

	var esrt = reflect.SliceOf(ert)
	var srv = reflect.MakeSlice(esrt, int(count), int(count))
	for i := 0; i < int(count); i++ {
		var erv, _n = srv.Index(i), int(0)

		if ptr {
			var numNil = int64(0)
			numNil, _n, err = decodeNumNilBytes(bz)
			if slide(&bz, &n, _n) && err != nil {
				return
			}
			if numNil == 0 {

			} else if numNil == 1 {

				erv.Set(reflect.Zero(erv.Type()))
				continue
			} else {
				panic("should not happen")
			}
		}

		_n, err = cdc.decodeReflectBinary(bz, einfo, erv, opts)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
	}
	rv.Set(srv)
	return
}

func (cdc *Codec) decodeReflectBinaryStruct(bz []byte, info *TypeInfo, rv reflect.Value, _ FieldOptions) (n int, err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if printLog {
		fmt.Println("(d) decodeReflectBinaryStruct")
		defer func() {
			fmt.Printf("(d) -> err: %v\n", err)
		}()
	}
	_n := 0

	switch info.Type {

	case timeType:

		var t time.Time
		t, _n, err = DecodeTime(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		rv.Set(reflect.ValueOf(t))
		return

	default:

		for _, field := range info.Fields {

			var frv = rv.Field(field.Index)
			var finfo *TypeInfo
			finfo, err = cdc.getTypeInfo_wlock(field.Type)
			if err != nil {
				return
			}

			var fieldNum, typ = uint32(0), Typ3(0x00)
			fieldNum, typ, _n, err = decodeFieldNumberAndTyp3(bz)
			if field.BinFieldNum < fieldNum {

				frv.Set(reflect.Zero(frv.Type()))
				continue

			}
			if fieldNum == 0 {

				break
			}
			if slide(&bz, &n, _n) && err != nil {
				return
			}

			if field.BinFieldNum != fieldNum {
				err = errors.New(fmt.Sprintf("Expected field number %v, got %v", field.BinFieldNum, fieldNum))
				return
			}
			typWanted := typeToTyp4(field.Type, field.FieldOptions).Typ3()
			if typ != typWanted {
				err = errors.New(fmt.Sprintf("Expected field type %X, got %X", typWanted, typ))
				return
			}

			_n, err = cdc.decodeReflectBinary(bz, finfo, frv, field.FieldOptions)
			if slide(&bz, &n, _n) && err != nil {
				return
			}
		}

		var typ = Typ3(0x00)
		typ, _n, err = decodeTyp3(bz)
		if slide(&bz, &n, _n) && err != nil {
			return
		}
		if typ != Typ3_StructTerm {
			err = errors.New(fmt.Sprintf("Expected StructTerm typ3 byte, got %X", typ))
			return
		}
		return
	}
}

func DecodeDisambPrefixBytes(bz []byte) (db DisambBytes, hasDb bool, pb PrefixBytes, typ Typ3, hasPb bool, isNil bool, n int, err error) {

	if len(bz) >= 2 && bz[0] == 0x00 && bz[1] == 0x00 {
		isNil = true
		n = 2
		return
	}

	if len(bz) < 4 {
		err = errors.New("EOF reading prefix bytes.")
		return
	}
	if bz[0] == 0x00 {

		if len(bz) < 8 {
			err = errors.New("EOF reading disamb bytes.")
			return
		}
		copy(db[0:3], bz[1:4])
		copy(pb[0:4], bz[4:8])
		pb, typ = pb.SplitTyp3()
		hasDb = true
		hasPb = true
		n = 8
		return
	} else {

		copy(pb[0:4], bz[0:4])
		pb, typ = pb.SplitTyp3()
		hasDb = false
		hasPb = true
		n = 4
		return
	}
}

func decodeFieldNumberAndTyp3(bz []byte) (num uint32, typ Typ3, n int, err error) {

	var value64 = uint64(0)
	value64, n, err = DecodeUvarint(bz)
	if err != nil {
		return
	}

	typ = Typ3(value64 & 0x07)

	var num64 uint64
	num64 = value64 >> 3
	if num64 > (1<<29 - 1) {
		err = errors.New(fmt.Sprintf("invalid field num %v", num64))
		return
	}
	num = uint32(num64)
	return
}

func decodeTyp4AndCheck(rt reflect.Type, bz []byte, opts FieldOptions) (ptr bool, n int, err error) {
	var typ = Typ4(0x00)
	typ, n, err = decodeTyp4(bz)
	if err != nil {
		return
	}
	var typWanted = typeToTyp4(rt, opts)
	if typWanted != typ {
		err = errors.New(fmt.Sprintf("Typ4 mismatch.  Expected %X, got %X", typWanted, typ))
		return
	}
	ptr = (typ & 0x08) != 0
	return
}

func decodeTyp4(bz []byte) (typ Typ4, n int, err error) {
	if len(bz) == 0 {
		err = errors.New(fmt.Sprintf("EOF reading typ4 byte"))
		return
	}
	if bz[0]&0xF0 != 0 {
		err = errors.New(fmt.Sprintf("Invalid non-zero nibble reading typ4 byte"))
		return
	}
	typ = Typ4(bz[0])
	n = 1
	return
}

func checkTyp3(rt reflect.Type, typ Typ3, opts FieldOptions) (err error) {
	typWanted := typeToTyp3(rt, opts)
	if typ != typWanted {
		err = fmt.Errorf("Typ3 mismatch.  Expected %X, got %X", typWanted, typ)
	}
	return
}

func decodeTyp3(bz []byte) (typ Typ3, n int, err error) {
	if len(bz) == 0 {
		err = fmt.Errorf("EOF reading typ3 byte")
		return
	}
	if bz[0]&0xF8 != 0 {
		err = fmt.Errorf("Invalid typ3 byte")
		return
	}
	typ = Typ3(bz[0])
	n = 1
	return
}

func decodeNumNilBytes(bz []byte) (numNil int64, n int, err error) {
	if len(bz) == 0 {
		err = errors.New("EOF reading nil byte(s)")
		return
	}
	if bz[0] == 0x00 {
		numNil, n = 0, 1
		return
	}
	if bz[0] == 0x01 {
		numNil, n = 1, 1
		return
	}
	n, err = 0, fmt.Errorf("Unexpected nil byte %X (sparse lists not supported)", bz[0])
	return
}
