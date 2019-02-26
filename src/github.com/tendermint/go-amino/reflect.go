package amino

import (
	"fmt"
	"reflect"
	"time"
)

const printLog = false

const RFC3339Millis = "2006-01-02T15:04:05.000Z"

var timeType = reflect.TypeOf(time.Time{})

func getTypeFromPointer(ptr interface{}) reflect.Type {
	rt := reflect.TypeOf(ptr)
	if rt.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("expected pointer, got %v", rt))
	}
	return rt.Elem()
}

func checkUnsafe(field FieldInfo) {
	if field.Unsafe {
		return
	}
	switch field.Type.Kind() {
	case reflect.Float32, reflect.Float64:
		panic("floating point types are unsafe for go-amino")
	}
}

func slide(bz *[]byte, n *int, _n int) bool {
	if _n < 0 || _n > len(*bz) {
		panic(fmt.Sprintf("impossible slide: len:%v _n:%v", len(*bz), _n))
	}
	*bz = (*bz)[_n:]
	*n += _n
	return true
}

func derefPointers(rv reflect.Value) (drv reflect.Value, isPtr bool, isNilPtr bool) {
	for rv.Kind() == reflect.Ptr {
		isPtr = true
		if rv.IsNil() {
			isNilPtr = true
			return
		}
		rv = rv.Elem()
	}
	drv = rv
	return
}

func isVoid(rv reflect.Value) (erv reflect.Value, isVoid bool) {
	rv, _, isNilPtr := derefPointers(rv)
	if isNilPtr {
		return rv, true
	} else {
		switch rv.Kind() {
		case reflect.String:
			return rv, rv.Len() == 0
		case reflect.Chan, reflect.Map, reflect.Slice:
			return rv, rv.IsNil() || rv.Len() == 0
		case reflect.Func, reflect.Interface:
			return rv, rv.IsNil()
		default:
			return rv, false
		}
	}
}

func constructConcreteType(cinfo *TypeInfo) (crv, irvSet reflect.Value) {

	if cinfo.PointerPreferred {
		cPtrRv := reflect.New(cinfo.Type)
		crv = cPtrRv.Elem()
		irvSet = cPtrRv
	} else {
		crv = reflect.New(cinfo.Type).Elem()
		irvSet = crv
	}
	return
}

func typeToTyp4(rt reflect.Type, opts FieldOptions) (typ Typ4) {

	var pointer = false
	for rt.Kind() == reflect.Ptr {
		pointer = true
		rt = rt.Elem()
	}

	typ = Typ4(typeToTyp3(rt, opts))

	if pointer {
		typ |= Typ4_Pointer
	}
	return
}

func typeToTyp3(rt reflect.Type, opts FieldOptions) Typ3 {
	switch rt.Kind() {
	case reflect.Interface:
		return Typ3_Interface
	case reflect.Array, reflect.Slice:
		ert := rt.Elem()
		switch ert.Kind() {
		case reflect.Uint8:
			return Typ3_ByteLength
		default:
			return Typ3_List
		}
	case reflect.String:
		return Typ3_ByteLength
	case reflect.Struct, reflect.Map:
		return Typ3_Struct
	case reflect.Int64, reflect.Uint64:
		if opts.BinVarint {
			return Typ3_Varint
		}
		return Typ3_8Byte
	case reflect.Float64:
		return Typ3_8Byte
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return Typ3_4Byte
	case reflect.Int16, reflect.Int8, reflect.Int,
		reflect.Uint16, reflect.Uint8, reflect.Uint, reflect.Bool:
		return Typ3_Varint
	default:
		panic(fmt.Sprintf("unsupported field type %v", rt))
	}
}

func toReprObject(rv reflect.Value) (rrv reflect.Value, err error) {
	var mwrm reflect.Value
	if rv.CanAddr() {
		mwrm = rv.Addr().MethodByName("MarshalAmino")
	} else {
		mwrm = rv.MethodByName("MarshalAmino")
	}
	mwouts := mwrm.Call(nil)
	if !mwouts[1].IsNil() {
		err = mwouts[1].Interface().(error)
		if err != nil {
			return
		}
	}
	rrv = mwouts[0]
	return
}
