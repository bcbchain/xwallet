package data

import (
	"github.com/pkg/errors"
	wire "github.com/tendermint/go-wire"
)

type binaryMapper struct {
	base	interface{}
	impls	[]wire.ConcreteType
}

func newBinaryMapper(base interface{}) *binaryMapper {
	return &binaryMapper{
		base: base,
	}
}

func (m *binaryMapper) registerImplementation(data interface{}, kind string, b byte) {
	m.impls = append(m.impls, wire.ConcreteType{O: data, Byte: b})
	wire.RegisterInterface(m.base, m.impls...)
}

func ToWire(o interface{}) ([]byte, error) {
	return wire.BinaryBytes(o), nil
}

func FromWire(d []byte, o interface{}) error {
	return errors.WithStack(
		wire.ReadBinaryBytes(d, o))
}
