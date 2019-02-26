package data

import "fmt"

type Mapper struct {
	*jsonMapper
	*binaryMapper
}

func NewMapper(base interface{}) Mapper {
	return Mapper{
		jsonMapper:	newJsonMapper(base),
		binaryMapper:	newBinaryMapper(base),
	}
}

func (m Mapper) RegisterImplementation(data interface{}, kind string, b byte) Mapper {
	m.jsonMapper.registerImplementation(data, kind, b)
	m.binaryMapper.registerImplementation(data, kind, b)
	return m
}

func ToText(o interface{}) (string, error) {
	d, err := ToJSON(o)
	if err != nil {
		return "", err
	}

	var s string
	err = FromJSON(d, &s)
	if err == nil {
		return s, nil
	}

	text := textEnv{}
	err = FromJSON(d, &text)
	if err != nil {
		return "", err
	}
	res := fmt.Sprintf("%s:%s", text.Kind, text.Data)
	return res, nil
}

type textEnv struct {
	Kind	string	`json:"type"`
	Data	string	`json:"data"`
}
