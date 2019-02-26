package amino

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

func (cdc *Codec) encodeReflectJSON(w io.Writer, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if !rv.IsValid() {
		panic("should not happen")
	}
	if printLog {
		spew.Printf("(E) encodeReflectJSON(info: %v, rv: %#v (%v), opts: %v)\n",
			info, rv.Interface(), rv.Type(), opts)
		defer func() {
			fmt.Printf("(E) -> err: %v\n", err)
		}()
	}

	if info.Registered {

		disfix := toDisfix(info.Disamb, info.Prefix)
		err = writeStr(w, _fmt(`{"type":"%X","value":`, disfix))
		if err != nil {
			return
		}

		defer func() {
			if err != nil {
				return
			}
			err = writeStr(w, `}`)
		}()
	}

	err = cdc._encodeReflectJSON(w, info, rv, opts)
	return
}

func (cdc *Codec) _encodeReflectJSON(w io.Writer, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if !rv.IsValid() {
		panic("should not happen")
	}
	if printLog {
		spew.Printf("(_) _encodeReflectJSON(info: %v, rv: %#v (%v), opts: %v)\n",
			info, rv.Interface(), rv.Type(), opts)
		defer func() {
			fmt.Printf("(_) -> err: %v\n", err)
		}()
	}

	var isNilPtr bool
	rv, _, isNilPtr = derefPointers(rv)

	if isNilPtr {
		err = writeStr(w, `null`)
		return
	}

	if rv.CanAddr() {
		if rv.Addr().Type().Implements(jsonMarshalerType) {
			err = invokeMarshalJSON(w, rv.Addr())
			return
		}
	} else if rv.Type().Implements(jsonMarshalerType) {
		err = invokeMarshalJSON(w, rv)
		return
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

		err = cdc._encodeReflectJSON(w, rinfo, rrv, opts)
		return
	}

	switch info.Type.Kind() {

	case reflect.Interface:
		return cdc.encodeReflectJSONInterface(w, info, rv, opts)

	case reflect.Array, reflect.Slice:
		return cdc.encodeReflectJSONList(w, info, rv, opts)

	case reflect.Struct:
		return cdc.encodeReflectJSONStruct(w, info, rv, opts)

	case reflect.Map:
		return cdc.encodeReflectJSONMap(w, info, rv, opts)

	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int,
		reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
		return invokeStdlibJSONMarshal(w, rv.Interface())

	case reflect.Float64, reflect.Float32:
		if !opts.Unsafe {
			return errors.New("Amino.JSON float* support requires `amino:\"unsafe\"`.")
		}
		fallthrough
	case reflect.Bool, reflect.String:
		return invokeStdlibJSONMarshal(w, rv.Interface())

	default:
		panic(fmt.Sprintf("unsupported type %v", info.Type.Kind()))
	}
}

func (cdc *Codec) encodeReflectJSONInterface(w io.Writer, iinfo *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if printLog {
		fmt.Println("(e) encodeReflectJSONInterface")
		defer func() {
			fmt.Printf("(e) -> err: %v\n", err)
		}()
	}

	if rv.IsNil() {
		err = writeStr(w, `null`)
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

	disfix := toDisfix(cinfo.Disamb, cinfo.Prefix)
	err = writeStr(w, _fmt(`{"type":"%X","value":`, disfix))
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			return
		}
		err = writeStr(w, `}`)
	}()

	err = cdc._encodeReflectJSON(w, cinfo, crv, opts)
	return
}

func (cdc *Codec) encodeReflectJSONList(w io.Writer, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if printLog {
		fmt.Println("(e) encodeReflectJSONList")
		defer func() {
			fmt.Printf("(e) -> err: %v\n", err)
		}()
	}

	if rv.Kind() == reflect.Slice && rv.IsNil() {
		err = writeStr(w, `null`)
		return
	}

	ert := info.Type.Elem()
	length := rv.Len()

	switch ert.Kind() {

	case reflect.Uint8:

		bz := []byte(nil)
		if rv.CanAddr() {
			bz = rv.Slice(0, length).Bytes()
		} else {
			bz = make([]byte, length)
			reflect.Copy(reflect.ValueOf(bz), rv)
		}
		jsonBytes := []byte(nil)
		jsonBytes, err = json.Marshal(bz)
		if err != nil {
			return
		}
		_, err = w.Write(jsonBytes)
		return

	default:

		err = writeStr(w, `[`)
		if err != nil {
			return
		}

		var einfo *TypeInfo
		einfo, err = cdc.getTypeInfo_wlock(ert)
		if err != nil {
			return
		}
		for i := 0; i < length; i++ {

			var erv, _, isNil = derefPointers(rv.Index(i))
			if isNil {
				err = writeStr(w, `null`)
			} else {
				err = cdc.encodeReflectJSON(w, einfo, erv, opts)
			}
			if err != nil {
				return
			}

			if i != length-1 {
				err = writeStr(w, `,`)
				if err != nil {
					return
				}
			}
		}

		defer func() {
			err = writeStr(w, `]`)
		}()
		return
	}
}

func (cdc *Codec) encodeReflectJSONStruct(w io.Writer, info *TypeInfo, rv reflect.Value, _ FieldOptions) (err error) {
	if printLog {
		fmt.Println("(e) encodeReflectJSONStruct")
		defer func() {
			fmt.Printf("(e) -> err: %v\n", err)
		}()
	}

	err = writeStr(w, `{`)
	if err != nil {
		return
	}

	defer func() {
		if err == nil {
			err = writeStr(w, `}`)
		}
	}()

	var writeComma = false
	for _, field := range info.Fields {

		var frv, _, isNil = derefPointers(rv.Field(field.Index))
		var finfo *TypeInfo
		finfo, err = cdc.getTypeInfo_wlock(field.Type)
		if err != nil {
			return
		}

		if field.JSONOmitEmpty && isEmpty(frv, field.ZeroValue) {
			continue
		}

		if writeComma {
			err = writeStr(w, `,`)
			if err != nil {
				return
			}
			writeComma = false
		}

		err = invokeStdlibJSONMarshal(w, field.JSONName)
		if err != nil {
			return
		}

		err = writeStr(w, `:`)
		if err != nil {
			return
		}

		if isNil {
			err = writeStr(w, `null`)
		} else {
			err = cdc.encodeReflectJSON(w, finfo, frv, field.FieldOptions)
		}
		if err != nil {
			return
		}
		writeComma = true
	}
	return
}

func (cdc *Codec) encodeReflectJSONMap(w io.Writer, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if printLog {
		fmt.Println("(e) encodeReflectJSONMap")
		defer func() {
			fmt.Printf("(e) -> err: %v\n", err)
		}()
	}

	err = writeStr(w, `{`)
	if err != nil {
		return
	}

	defer func() {
		if err == nil {
			err = writeStr(w, `}`)
		}
	}()

	if rv.Type().Key().Kind() != reflect.String {
		err = errors.New("encodeReflectJSONMap: map key type must be a string")
		return
	}

	var writeComma = false
	for _, krv := range rv.MapKeys() {

		var vrv, _, isNil = derefPointers(rv.MapIndex(krv))

		if writeComma {
			err = writeStr(w, `,`)
			if err != nil {
				return
			}
			writeComma = false
		}

		err = invokeStdlibJSONMarshal(w, krv.Interface())
		if err != nil {
			return
		}

		err = writeStr(w, `:`)
		if err != nil {
			return
		}

		if isNil {
			err = writeStr(w, `null`)
		} else {
			var vinfo *TypeInfo
			vinfo, err = cdc.getTypeInfo_wlock(vrv.Type())
			if err != nil {
				return
			}
			err = cdc.encodeReflectJSON(w, vinfo, vrv, opts)
		}
		if err != nil {
			return
		}
		writeComma = true
	}
	return

}

func invokeMarshalJSON(w io.Writer, rv reflect.Value) error {
	blob, err := rv.Interface().(json.Marshaler).MarshalJSON()
	if err != nil {
		return err
	}
	_, err = w.Write(blob)
	return err
}

func invokeStdlibJSONMarshal(w io.Writer, v interface{}) error {

	blob, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(blob)
	return err
}

func writeStr(w io.Writer, s string) (err error) {
	_, err = w.Write([]byte(s))
	return
}

func _fmt(s string, args ...interface{}) string {
	return fmt.Sprintf(s, args...)
}

func isEmpty(rv reflect.Value, zrv reflect.Value) bool {
	if !rv.IsValid() {
		return true
	}
	if reflect.DeepEqual(rv.Interface(), zrv.Interface()) {
		return true
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array, reflect.String:
		if rv.Len() == 0 {
			return true
		}
	}
	return false
}
