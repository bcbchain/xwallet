package data

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

type jsonMapper struct {
	kindToType	map[string]reflect.Type
	typeToKind	map[reflect.Type]string
}

func newJsonMapper(base interface{}) *jsonMapper {
	return &jsonMapper{
		kindToType:	map[string]reflect.Type{},
		typeToKind:	map[reflect.Type]string{},
	}
}

func ToJSON(o interface{}) ([]byte, error) {
	d, err := json.MarshalIndent(o, "", "  ")
	return d, errors.WithStack(err)
}

func FromJSON(d []byte, o interface{}) error {
	return errors.WithStack(
		json.Unmarshal(d, o))
}

func (m *jsonMapper) registerImplementation(data interface{}, kind string, b byte) {
	typ := reflect.TypeOf(data)
	m.kindToType[kind] = typ
	m.typeToKind[typ] = kind
}

func (m *jsonMapper) getTarget(kind string) (interface{}, error) {
	typ, ok := m.kindToType[kind]
	if !ok {
		return nil, errors.Errorf("Unmarshaling into unknown type: %s", kind)
	}
	target := reflect.New(typ).Interface()
	return target, nil
}

func (m *jsonMapper) getKind(obj interface{}) (string, error) {
	typ := reflect.TypeOf(obj)
	kind, ok := m.typeToKind[typ]
	if !ok {
		return "", errors.Errorf("Marshalling from unknown type: %#v", obj)
	}
	return kind, nil
}

func (m *jsonMapper) FromJSON(data []byte) (interface{}, error) {

	if bytes.Equal(data, []byte("null")) {
		return nil, nil
	}
	e := envelope{
		Data: &json.RawMessage{},
	}
	err := json.Unmarshal(data, &e)
	if err != nil {
		return nil, err
	}

	bytes := *e.Data.(*json.RawMessage)
	res, err := m.getTarget(e.Kind)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &res)

	res = reflect.ValueOf(res).Elem().Interface()
	return res, err
}

func (m *jsonMapper) ToJSON(data interface{}) ([]byte, error) {

	if data == nil {
		return []byte("null"), nil
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	kind, err := m.getKind(data)
	if err != nil {
		return nil, err
	}
	msg := json.RawMessage(raw)
	e := envelope{
		Kind:	kind,
		Data:	&msg,
	}
	return json.Marshal(e)
}

type envelope struct {
	Kind	string		`json:"type"`
	Data	interface{}	`json:"data"`
}
