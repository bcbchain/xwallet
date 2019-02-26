package amino

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

func (cdc *Codec) decodeReflectJSON(bz []byte, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if info.Type.Kind() == reflect.Interface && rv.Kind() == reflect.Ptr {
		panic("should not happen")
	}
	if printLog {
		spew.Printf("(D) decodeReflectJSON(bz: %s, info: %v, rv: %#v (%v), opts: %v)\n",
			bz, info, rv.Interface(), rv.Type(), opts)
		defer func() {
			fmt.Printf("(D) -> err: %v\n", err)
		}()
	}

	if info.Registered {

		var disfix DisfixBytes
		disfix, bz, err = decodeDisfixJSON(bz)
		if err != nil {
			return
		}
		if !info.GetDisfix().EqualBytes(disfix[:]) {
			err = fmt.Errorf("Expected disfix bytes %X but got %X", info.GetDisfix(), disfix)
			return
		}
	}

	err = cdc._decodeReflectJSON(bz, info, rv, opts)
	return
}

func (cdc *Codec) _decodeReflectJSON(bz []byte, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if info.Type.Kind() == reflect.Interface && rv.Kind() == reflect.Ptr {
		panic("should not happen")
	}
	if printLog {
		spew.Printf("(_) _decodeReflectJSON(bz: %s, info: %v, rv: %#v (%v), opts: %v)\n",
			bz, info, rv.Interface(), rv.Type(), opts)
		defer func() {
			fmt.Printf("(_) -> err: %v\n", err)
		}()
	}

	if nullBytes(bz) {
		switch rv.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Array:
			rv.Set(reflect.Zero(rv.Type()))
			return
		}
	}

	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			newPtr := reflect.New(rv.Type().Elem())
			rv.Set(newPtr)
		}
		rv = rv.Elem()
	}

	if rv.Addr().Type().Implements(jsonUnmarshalerType) {
		err = rv.Addr().Interface().(json.Unmarshaler).UnmarshalJSON(bz)
		return
	}

	if info.IsAminoUnmarshaler {

		rrv, rinfo := reflect.New(info.AminoUnmarshalReprType).Elem(), (*TypeInfo)(nil)
		rinfo, err = cdc.getTypeInfo_wlock(info.AminoUnmarshalReprType)
		if err != nil {
			return
		}
		err = cdc._decodeReflectJSON(bz, rinfo, rrv, opts)
		if err != nil {
			return
		}

		uwrm := rv.Addr().MethodByName("UnmarshalAmino")
		uwouts := uwrm.Call([]reflect.Value{rrv})
		err = uwouts[0].Interface().(error)
		return
	}

	switch ikind := info.Type.Kind(); ikind {

	case reflect.Interface:
		err = cdc.decodeReflectJSONInterface(bz, info, rv, opts)

	case reflect.Array:
		err = cdc.decodeReflectJSONArray(bz, info, rv, opts)

	case reflect.Slice:
		err = cdc.decodeReflectJSONSlice(bz, info, rv, opts)

	case reflect.Struct:
		err = cdc.decodeReflectJSONStruct(bz, info, rv, opts)

	case reflect.Map:
		err = cdc.decodeReflectJSONMap(bz, info, rv, opts)

	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int,
		reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
		err = invokeStdlibJSONUnmarshal(bz, rv, opts)

	case reflect.Float32, reflect.Float64:
		if !opts.Unsafe {
			return errors.New("Amino.JSON float* support requires `amino:\"unsafe\"`.")
		}
		fallthrough
	case reflect.Bool, reflect.String:
		err = invokeStdlibJSONUnmarshal(bz, rv, opts)

	default:
		panic(fmt.Sprintf("unsupported type %v", info.Type.Kind()))
	}

	return
}

func invokeStdlibJSONUnmarshal(bz []byte, rv reflect.Value, opts FieldOptions) error {
	if !rv.CanAddr() && rv.Kind() != reflect.Ptr {
		panic("rv not addressable nor pointer")
	}

	var rrv reflect.Value = rv
	if rv.Kind() != reflect.Ptr {
		rrv = reflect.New(rv.Type())
	}

	if err := json.Unmarshal(bz, rrv.Interface()); err != nil {
		return err
	}
	rv.Set(rrv.Elem())
	return nil
}

func (cdc *Codec) decodeReflectJSONInterface(bz []byte, iinfo *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if printLog {
		fmt.Println("(d) decodeReflectJSONInterface")
		defer func() {
			fmt.Printf("(d) -> err: %v\n", err)
		}()
	}

	if !rv.IsNil() {

		rv.Set(iinfo.ZeroValue)
	}

	disfix, bz, err := decodeDisfixJSON(bz)
	if err != nil {
		return
	}

	var cinfo *TypeInfo
	cinfo, err = cdc.getTypeInfoFromDisfix_rlock(disfix)
	if err != nil {
		return
	}

	var crv, irvSet = constructConcreteType(cinfo)

	err = cdc._decodeReflectJSON(bz, cinfo, crv, opts)
	if err != nil {
		rv.Set(irvSet)
		return
	}

	rv.Set(irvSet)
	return
}

func (cdc *Codec) decodeReflectJSONArray(bz []byte, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if printLog {
		fmt.Println("(d) decodeReflectJSONArray")
		defer func() {
			fmt.Printf("(d) -> err: %v\n", err)
		}()
	}
	ert := info.Type.Elem()
	length := info.Type.Len()

	switch ert.Kind() {

	case reflect.Uint8:
		var buf []byte
		err = json.Unmarshal(bz, &buf)
		if err != nil {
			return
		}
		if len(buf) != length {
			err = fmt.Errorf("decodeReflectJSONArray: byte-length mismatch, got %v want %v",
				len(buf), length)
		}
		reflect.Copy(rv, reflect.ValueOf(buf))
		return

	default:
		var einfo *TypeInfo
		einfo, err = cdc.getTypeInfo_wlock(ert)
		if err != nil {
			return
		}

		var rawSlice []json.RawMessage
		if err = json.Unmarshal(bz, &rawSlice); err != nil {
			return
		}
		if len(rawSlice) != length {
			err = fmt.Errorf("decodeReflectJSONArray: length mismatch, got %v want %v", len(rawSlice), length)
			return
		}

		for i := 0; i < length; i++ {
			erv := rv.Index(i)
			ebz := rawSlice[i]
			err = cdc.decodeReflectJSON(ebz, einfo, erv, opts)
			if err != nil {
				return
			}
		}
		return
	}
}

func (cdc *Codec) decodeReflectJSONSlice(bz []byte, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if printLog {
		fmt.Println("(d) decodeReflectJSONSlice")
		defer func() {
			fmt.Printf("(d) -> err: %v\n", err)
		}()
	}

	var ert = info.Type.Elem()

	switch ert.Kind() {

	case reflect.Uint8:
		err = json.Unmarshal(bz, rv.Addr().Interface())
		if err != nil {
			return
		}
		if rv.Len() == 0 {

			rv.Set(info.ZeroValue)
		} else {

		}
		return

	default:
		var einfo *TypeInfo
		einfo, err = cdc.getTypeInfo_wlock(ert)
		if err != nil {
			return
		}

		var rawSlice []json.RawMessage
		if err = json.Unmarshal(bz, &rawSlice); err != nil {
			return
		}

		var length = len(rawSlice)
		if length == 0 {
			rv.Set(info.ZeroValue)
			return
		}

		var esrt = reflect.SliceOf(ert)
		var srv = reflect.MakeSlice(esrt, length, length)
		for i := 0; i < length; i++ {
			erv := srv.Index(i)
			ebz := rawSlice[i]
			err = cdc.decodeReflectJSON(ebz, einfo, erv, opts)
			if err != nil {
				return
			}
		}

		rv.Set(srv)
		return
	}
}

func (cdc *Codec) decodeReflectJSONStruct(bz []byte, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if printLog {
		fmt.Println("(d) decodeReflectJSONStruct")
		defer func() {
			fmt.Printf("(d) -> err: %v\n", err)
		}()
	}

	var rawMap = make(map[string]json.RawMessage)
	err = json.Unmarshal(bz, &rawMap)
	if err != nil {
		return
	}

	for _, field := range info.Fields {

		var frv = rv.Field(field.Index)
		var finfo *TypeInfo
		finfo, err = cdc.getTypeInfo_wlock(field.Type)
		if err != nil {
			return
		}

		var valueBytes = rawMap[field.JSONName]
		if len(valueBytes) == 0 {

			frv.Set(reflect.Zero(frv.Type()))
			continue
		}

		err = cdc.decodeReflectJSON(valueBytes, finfo, frv, opts)
		if err != nil {
			return
		}
	}

	return nil
}

func (cdc *Codec) decodeReflectJSONMap(bz []byte, info *TypeInfo, rv reflect.Value, opts FieldOptions) (err error) {
	if !rv.CanAddr() {
		panic("rv not addressable")
	}
	if printLog {
		fmt.Println("(d) decodeReflectJSONMap")
		defer func() {
			fmt.Printf("(d) -> err: %v\n", err)
		}()
	}

	var rawMap = make(map[string]json.RawMessage)
	err = json.Unmarshal(bz, &rawMap)
	if err != nil {
		return
	}

	var krt = rv.Type().Key()
	if krt.Kind() != reflect.String {
		err = fmt.Errorf("decodeReflectJSONMap: key type must be string")
		return
	}
	var vinfo *TypeInfo
	vinfo, err = cdc.getTypeInfo_wlock(rv.Type().Elem())
	if err != nil {
		return
	}

	var mrv = reflect.MakeMapWithSize(rv.Type(), len(rawMap))
	for key, valueBytes := range rawMap {

		vrv := reflect.New(mrv.Type().Elem()).Elem()

		err = cdc.decodeReflectJSON(valueBytes, vinfo, vrv, opts)
		if err != nil {
			return
		}

		krv := reflect.New(reflect.TypeOf("")).Elem()
		krv.SetString(key)
		mrv.SetMapIndex(krv, vrv)
	}
	rv.Set(mrv)

	return nil
}

type disfixWrapper struct {
	Disfix	string		`json:"type"`
	Data	json.RawMessage	`json:"value"`
}

func decodeDisfixJSON(bz []byte) (df DisfixBytes, data []byte, err error) {
	if string(bz) == "null" {
		panic("yay")
	}
	dfw := new(disfixWrapper)
	err = json.Unmarshal(bz, dfw)
	if err != nil {
		err = fmt.Errorf("Cannot parse disfix JSON wrapper: %v", err)
		return
	}
	dfBytes, err := hex.DecodeString(dfw.Disfix)
	if err != nil {
		return
	}

	if g, w := len(dfBytes), DisfixBytesLen; g != w {
		err = fmt.Errorf("Disfix length got=%d want=%d data=%s", g, w, bz)
		return
	}
	copy(df[:], dfBytes)
	if (DisfixBytes{}).EqualBytes(df[:]) {
		err = errors.New("Unexpected zero disfix in JSON")
		return
	}

	if len(dfw.Data) == 0 {
		err = errors.New("Disfix JSON wrapper should have non-empty value field")
		return
	}
	data = dfw.Data
	return
}

func nullBytes(b []byte) bool {
	return bytes.Equal(b, []byte(`null`))
}
