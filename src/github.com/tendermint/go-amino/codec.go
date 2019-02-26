package amino

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

const (
	PrefixBytesLen	= 4
	DisambBytesLen	= 3
	DisfixBytesLen	= PrefixBytesLen + DisambBytesLen
)

type (
	PrefixBytes	[PrefixBytesLen]byte
	DisambBytes	[DisambBytesLen]byte
	DisfixBytes	[DisfixBytesLen]byte
)

func NewPrefixBytes(prefixBytes []byte) PrefixBytes {
	pb := PrefixBytes{}
	copy(pb[:], prefixBytes)
	return pb
}

func (pb PrefixBytes) Bytes() []byte			{ return pb[:] }
func (pb PrefixBytes) EqualBytes(bz []byte) bool	{ return bytes.Equal(pb[:], bz) }
func (pb PrefixBytes) WithTyp3(typ Typ3) PrefixBytes	{ pb[3] |= byte(typ); return pb }
func (pb PrefixBytes) SplitTyp3() (PrefixBytes, Typ3) {
	typ := Typ3(pb[3] & 0x07)
	pb[3] &= 0xF8
	return pb, typ
}
func (db DisambBytes) Bytes() []byte			{ return db[:] }
func (db DisambBytes) EqualBytes(bz []byte) bool	{ return bytes.Equal(db[:], bz) }
func (df DisfixBytes) Bytes() []byte			{ return df[:] }
func (df DisfixBytes) EqualBytes(bz []byte) bool	{ return bytes.Equal(df[:], bz) }

func NameToDisfix(name string) (db DisambBytes, pb PrefixBytes) {
	return nameToDisfix(name)
}

type TypeInfo struct {
	Type		reflect.Type
	PtrToType	reflect.Type
	ZeroValue	reflect.Value
	ZeroProto	interface{}
	InterfaceInfo
	ConcreteInfo
	StructInfo
}

type InterfaceInfo struct {
	Priority	[]DisfixBytes
	Implementers	map[PrefixBytes][]*TypeInfo
	InterfaceOptions
}

type InterfaceOptions struct {
	Priority		[]string
	AlwaysDisambiguate	bool
}

type ConcreteInfo struct {
	Registered		bool
	PointerPreferred	bool
	Name			string
	Disamb			DisambBytes
	Prefix			PrefixBytes
	ConcreteOptions

	IsAminoMarshaler	bool
	AminoMarshalReprType	reflect.Type
	IsAminoUnmarshaler	bool
	AminoUnmarshalReprType	reflect.Type
}

type StructInfo struct {
	Fields []FieldInfo
}

func (cinfo ConcreteInfo) GetDisfix() DisfixBytes {
	return toDisfix(cinfo.Disamb, cinfo.Prefix)
}

type ConcreteOptions struct {
}

type FieldInfo struct {
	Name		string
	Type		reflect.Type
	Index		int
	ZeroValue	reflect.Value
	FieldOptions
	BinTyp3	Typ3
}

type FieldOptions struct {
	JSONName	string
	JSONOmitEmpty	bool
	BinVarint	bool
	BinFieldNum	uint32
	Unsafe		bool
}

type Codec struct {
	mtx			sync.RWMutex
	typeInfos		map[reflect.Type]*TypeInfo
	interfaceInfos		[]*TypeInfo
	concreteInfos		[]*TypeInfo
	disfixToTypeInfo	map[DisfixBytes]*TypeInfo
}

func NewCodec() *Codec {
	cdc := &Codec{
		typeInfos:		make(map[reflect.Type]*TypeInfo),
		disfixToTypeInfo:	make(map[DisfixBytes]*TypeInfo),
	}
	return cdc
}

func (cdc *Codec) RegisterInterface(ptr interface{}, opts *InterfaceOptions) {

	rt := getTypeFromPointer(ptr)
	if rt.Kind() != reflect.Interface {
		panic(fmt.Sprintf("RegisterInterface expects an interface, got %v", rt))
	}

	var info = cdc.newTypeInfoFromInterfaceType(rt, opts)

	func() {
		cdc.mtx.Lock()
		defer cdc.mtx.Unlock()

		cdc.collectImplementers_nolock(info)
		err := cdc.checkConflictsInPrio_nolock(info)
		if err != nil {
			panic(err)
		}
		cdc.setTypeInfo_nolock(info)
	}()

}

func (cdc *Codec) RegisterConcrete(o interface{}, name string, opts *ConcreteOptions) {

	var pointerPreferred bool

	rt := reflect.TypeOf(o)
	if rt.Kind() == reflect.Interface {
		panic(fmt.Sprintf("expected a non-interface: %v", rt))
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		if rt.Kind() == reflect.Ptr {

			panic(fmt.Sprintf("registering pointer-pointers not yet supported: *%v", rt))
		}
		if rt.Kind() == reflect.Interface {

			panic(fmt.Sprintf("registering interface-pointers not yet supported: *%v", rt))
		}
		pointerPreferred = true
	}

	var info = cdc.newTypeInfoFromRegisteredConcreteType(rt, pointerPreferred, name, opts)

	func() {
		cdc.mtx.Lock()
		defer cdc.mtx.Unlock()

		cdc.addCheckConflictsWithConcrete_nolock(info)
		cdc.setTypeInfo_nolock(info)
	}()
}

func (cdc *Codec) setTypeInfo_nolock(info *TypeInfo) {

	if info.Type.Kind() == reflect.Ptr {
		panic(fmt.Sprintf("unexpected pointer type"))
	}
	if _, ok := cdc.typeInfos[info.Type]; ok {
		panic(fmt.Sprintf("TypeInfo already exists for %v", info.Type))
	}

	cdc.typeInfos[info.Type] = info
	if info.Type.Kind() == reflect.Interface {
		cdc.interfaceInfos = append(cdc.interfaceInfos, info)
	} else if info.Registered {
		cdc.concreteInfos = append(cdc.concreteInfos, info)
		disfix := info.GetDisfix()
		if existing, ok := cdc.disfixToTypeInfo[disfix]; ok {
			panic(fmt.Sprintf("disfix <%X> already registered for %v", disfix, existing.Type))
		}
		cdc.disfixToTypeInfo[disfix] = info

	}
}

func (cdc *Codec) getTypeInfo_wlock(rt reflect.Type) (info *TypeInfo, err error) {
	cdc.mtx.Lock()
	defer cdc.mtx.Unlock()

	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	info, ok := cdc.typeInfos[rt]
	if !ok {
		if rt.Kind() == reflect.Interface {
			err = fmt.Errorf("Unregistered interface %v", rt)
			return
		}

		info = cdc.newTypeInfoUnregistered(rt)
		cdc.setTypeInfo_nolock(info)
	}
	return info, nil
}

func (cdc *Codec) getTypeInfoFromPrefix_rlock(iinfo *TypeInfo, pb PrefixBytes) (info *TypeInfo, err error) {
	cdc.mtx.RLock()
	defer cdc.mtx.RUnlock()

	infos, ok := iinfo.Implementers[pb]
	if !ok {
		err = fmt.Errorf("unrecognized prefix bytes %X", pb)
		return
	}
	if len(infos) > 1 {
		err = fmt.Errorf("Conflicting concrete types registered for %X: e.g. %v and %v.", pb, infos[0].Type, infos[1].Type)
		return
	}
	info = infos[0]
	return
}

func (cdc *Codec) getTypeInfoFromDisfix_rlock(df DisfixBytes) (info *TypeInfo, err error) {
	cdc.mtx.RLock()
	defer cdc.mtx.RUnlock()

	info, ok := cdc.disfixToTypeInfo[df]
	if !ok {
		err = fmt.Errorf("unrecognized disambiguation+prefix bytes %X", df)
		return
	}
	return
}

func (cdc *Codec) parseStructInfo(rt reflect.Type) (sinfo StructInfo) {
	if rt.Kind() != reflect.Struct {
		panic("should not happen")
	}

	var infos = make([]FieldInfo, 0, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		var field = rt.Field(i)
		var ftype = field.Type
		if !isExported(field) {
			continue
		}
		skip, opts := cdc.parseFieldOptions(field)
		if skip {
			continue
		}

		opts.BinFieldNum = uint32(len(infos) + 1)
		fieldInfo := FieldInfo{
			Name:		field.Name,
			Index:		i,
			Type:		ftype,
			ZeroValue:	reflect.Zero(ftype),
			FieldOptions:	opts,
			BinTyp3:	typeToTyp4(ftype, opts).Typ3(),
		}
		checkUnsafe(fieldInfo)
		infos = append(infos, fieldInfo)
	}
	sinfo = StructInfo{infos}
	return
}

func (cdc *Codec) parseFieldOptions(field reflect.StructField) (skip bool, opts FieldOptions) {
	binTag := field.Tag.Get("binary")
	aminoTag := field.Tag.Get("amino")
	jsonTag := field.Tag.Get("json")

	if jsonTag == "-" {
		skip = true
		return
	}

	jsonTagParts := strings.Split(jsonTag, ",")
	if jsonTagParts[0] == "" {
		opts.JSONName = field.Name
	} else {
		opts.JSONName = jsonTagParts[0]
	}

	if len(jsonTagParts) > 1 {
		if jsonTagParts[1] == "omitempty" {
			opts.JSONOmitEmpty = true
		}
	}

	if binTag == "varint" {
		opts.BinVarint = true
	}

	if aminoTag == "unsafe" {
		opts.Unsafe = true
	}

	return
}

func (cdc *Codec) newTypeInfoUnregistered(rt reflect.Type) *TypeInfo {
	if rt.Kind() == reflect.Ptr {
		panic("unexpected pointer type")
	}
	if rt.Kind() == reflect.Interface {
		panic("unexpected interface type")
	}

	var info = new(TypeInfo)
	info.Type = rt
	info.PtrToType = reflect.PtrTo(rt)
	info.ZeroValue = reflect.Zero(rt)
	info.ZeroProto = reflect.Zero(rt).Interface()
	if rt.Kind() == reflect.Struct {
		info.StructInfo = cdc.parseStructInfo(rt)
	}
	if rm, ok := rt.MethodByName("MarshalAmino"); ok {
		info.ConcreteInfo.IsAminoMarshaler = true
		info.ConcreteInfo.AminoMarshalReprType = marshalAminoReprType(rm)
	}
	if rm, ok := rt.MethodByName("UnmarshalAmino"); ok {
		info.ConcreteInfo.IsAminoUnmarshaler = true
		info.ConcreteInfo.AminoUnmarshalReprType = unmarshalAminoReprType(rm)
	}
	return info
}

func (cdc *Codec) newTypeInfoFromInterfaceType(rt reflect.Type, opts *InterfaceOptions) *TypeInfo {
	if rt.Kind() != reflect.Interface {
		panic(fmt.Sprintf("expected interface type, got %v", rt))
	}

	var info = new(TypeInfo)
	info.Type = rt
	info.PtrToType = reflect.PtrTo(rt)
	info.ZeroValue = reflect.Zero(rt)
	info.ZeroProto = reflect.Zero(rt).Interface()
	info.InterfaceInfo.Implementers = make(map[PrefixBytes][]*TypeInfo)
	if opts != nil {
		info.InterfaceInfo.InterfaceOptions = *opts
		info.InterfaceInfo.Priority = make([]DisfixBytes, len(opts.Priority))

		for i, name := range opts.Priority {
			disamb, prefix := nameToDisfix(name)
			disfix := toDisfix(disamb, prefix)
			info.InterfaceInfo.Priority[i] = disfix
		}
	}
	return info
}

func (cdc *Codec) newTypeInfoFromRegisteredConcreteType(rt reflect.Type, pointerPreferred bool, name string, opts *ConcreteOptions) *TypeInfo {
	if rt.Kind() == reflect.Interface ||
		rt.Kind() == reflect.Ptr {
		panic(fmt.Sprintf("expected non-interface non-pointer concrete type, got %v", rt))
	}

	var info = cdc.newTypeInfoUnregistered(rt)
	info.ConcreteInfo.Registered = true
	info.ConcreteInfo.PointerPreferred = pointerPreferred
	info.ConcreteInfo.Name = name
	info.ConcreteInfo.Disamb = nameToDisamb(name)
	info.ConcreteInfo.Prefix = nameToPrefix(name)
	if opts != nil {
		info.ConcreteOptions = *opts
	}
	return info
}

func (cdc *Codec) collectImplementers_nolock(info *TypeInfo) {
	for _, cinfo := range cdc.concreteInfos {
		if cinfo.PtrToType.Implements(info.Type) {
			info.Implementers[cinfo.Prefix] = append(
				info.Implementers[cinfo.Prefix], cinfo)
		}
	}
}

func (cdc *Codec) checkConflictsInPrio_nolock(iinfo *TypeInfo) error {

	for _, cinfos := range iinfo.Implementers {
		if len(cinfos) < 2 {
			continue
		}
		for _, cinfo := range cinfos {
			var inPrio = false
			for _, disfix := range iinfo.InterfaceInfo.Priority {
				if cinfo.GetDisfix() == disfix {
					inPrio = true
				}
			}
			if !inPrio {
				return fmt.Errorf("%v conflicts with %v other(s). Add it to the priority list for %v.",
					cinfo.Type, len(cinfos), iinfo.Type)
			}
		}
	}
	return nil
}

func (cdc *Codec) addCheckConflictsWithConcrete_nolock(cinfo *TypeInfo) {

	for _, iinfo := range cdc.interfaceInfos {
		if !cinfo.PtrToType.Implements(iinfo.Type) {
			continue
		}

		var origImpls = iinfo.Implementers[cinfo.Prefix]
		iinfo.Implementers[cinfo.Prefix] = append(origImpls, cinfo)

		err := cdc.checkConflictsInPrio_nolock(iinfo)
		if err != nil {

			iinfo.Implementers[cinfo.Prefix] = origImpls
			panic(err)
		}
	}
}

func (ti TypeInfo) String() string {
	buf := new(bytes.Buffer)
	buf.Write([]byte("TypeInfo{"))
	buf.Write([]byte(fmt.Sprintf("Type:%v,", ti.Type)))
	if ti.Type.Kind() == reflect.Interface {
		buf.Write([]byte(fmt.Sprintf("Priority:%v,", ti.Priority)))
		buf.Write([]byte("Implementers:{"))
		for pb, cinfos := range ti.Implementers {
			buf.Write([]byte(fmt.Sprintf("\"%X\":", pb)))
			buf.Write([]byte(fmt.Sprintf("%v,", cinfos)))
		}
		buf.Write([]byte("}"))
		buf.Write([]byte(fmt.Sprintf("Priority:%v,", ti.InterfaceOptions.Priority)))
		buf.Write([]byte(fmt.Sprintf("AlwaysDisambiguate:%v,", ti.InterfaceOptions.AlwaysDisambiguate)))
	}
	if ti.Type.Kind() != reflect.Interface {
		if ti.ConcreteInfo.Registered {
			buf.Write([]byte("Registered:true,"))
			buf.Write([]byte(fmt.Sprintf("PointerPreferred:%v,", ti.PointerPreferred)))
			buf.Write([]byte(fmt.Sprintf("Name:\"%v\",", ti.Name)))
			buf.Write([]byte(fmt.Sprintf("Disamb:\"%X\",", ti.Disamb)))
			buf.Write([]byte(fmt.Sprintf("Prefix:\"%X\",", ti.Prefix)))
		} else {
			buf.Write([]byte("Registered:false,"))
		}
		buf.Write([]byte(fmt.Sprintf("AminoMarshalReprType:\"%X\",", ti.AminoMarshalReprType)))
		buf.Write([]byte(fmt.Sprintf("AminoUnmarshalReprType:\"%X\",", ti.AminoUnmarshalReprType)))
		if ti.Type.Kind() == reflect.Struct {
			buf.Write([]byte(fmt.Sprintf("Fields:%v,", ti.Fields)))
		}
	}
	buf.Write([]byte("}"))
	return buf.String()
}

func isExported(field reflect.StructField) bool {

	if field.PkgPath != "" {
		return false
	}

	var first rune
	for _, c := range field.Name {
		first = c
		break
	}

	if !unicode.IsUpper(first) {
		return false
	}

	return true
}

func nameToDisamb(name string) (db DisambBytes) {
	db, _ = nameToDisfix(name)
	return
}

func nameToPrefix(name string) (pb PrefixBytes) {
	_, pb = nameToDisfix(name)
	return
}

func nameToDisfix(name string) (db DisambBytes, pb PrefixBytes) {
	hasher := sha256.New()
	hasher.Write([]byte(name))
	bz := hasher.Sum(nil)
	for bz[0] == 0x00 {
		bz = bz[1:]
	}
	copy(db[:], bz[0:3])
	bz = bz[3:]
	for bz[0] == 0x00 {
		bz = bz[1:]
	}
	copy(pb[:], bz[0:4])

	pb[3] &= 0xF8
	return
}

func toDisfix(db DisambBytes, pb PrefixBytes) (df DisfixBytes) {
	copy(df[0:3], db[0:3])
	copy(df[3:7], pb[0:4])
	return
}

func marshalAminoReprType(rm reflect.Method) (rrt reflect.Type) {

	if rm.Type.NumIn() != 1 {
		panic(fmt.Sprintf("MarshalAmino should have 1 input parameters (including receiver); got %v", rm.Type))
	}
	if rm.Type.NumOut() != 2 {
		panic(fmt.Sprintf("MarshalAmino should have 2 output parameters; got %v", rm.Type))
	}
	if out := rm.Type.Out(1); out != errorType {
		panic(fmt.Sprintf("MarshalAmino should have second output parameter of error type, got %v", out))
	}
	rrt = rm.Type.Out(0)
	if rrt.Kind() == reflect.Ptr {
		panic(fmt.Sprintf("Representative objects cannot be pointers; got %v", rrt))
	}
	return
}

func unmarshalAminoReprType(rm reflect.Method) (rrt reflect.Type) {

	if rm.Type.NumIn() != 2 {
		panic(fmt.Sprintf("UnmarshalAmino should have 2 input parameters (including receiver); got %v", rm.Type))
	}
	if in1 := rm.Type.In(0); in1.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("UnmarshalAmino first input parameter should be pointer type but got %v", in1))
	}
	if rm.Type.NumOut() != 1 {
		panic(fmt.Sprintf("UnmarshalAmino should have 1 output parameters; got %v", rm.Type))
	}
	if out := rm.Type.Out(0); out != errorType {
		panic(fmt.Sprintf("UnmarshalAmino should have first output parameter of error type, got %v", out))
	}
	rrt = rm.Type.In(0)
	if rrt.Kind() == reflect.Ptr {
		panic(fmt.Sprintf("Representative objects cannot be pointers; got %v", rrt))
	}
	return
}
