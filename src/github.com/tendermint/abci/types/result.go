package types

import (
	"bytes"
	"encoding/json"
	"github.com/gogo/protobuf/jsonpb"
	"bcbchain.io/bcerrors"
)

const (
	CodeTypeOK uint32 = bcerrors.ErrCodeOK
)

func (r ResponseCheckTx) IsOK() bool {
	return r.Code == CodeTypeOK
}

func (r ResponseCheckTx) IsErr() bool {
	return r.Code != CodeTypeOK
}

func (r ResponseDeliverTx) IsOK() bool {
	return r.Code == CodeTypeOK
}

func (r ResponseDeliverTx) IsErr() bool {
	return r.Code != CodeTypeOK
}

func (r ResponseQuery) IsOK() bool {
	return r.Code == CodeTypeOK
}

func (r ResponseQuery) IsErr() bool {
	return r.Code != CodeTypeOK
}

var (
	jsonpbMarshaller	= jsonpb.Marshaler{
		EnumsAsInts:	true,
		EmitDefaults:	false,
	}
	jsonpbUnmarshaller	= jsonpb.Unmarshaler{}
)

func (r *ResponseSetOption) UnmarshalJSON(b []byte) error {
	reader := bytes.NewBuffer(b)
	return jsonpbUnmarshaller.Unmarshal(reader, r)
}

func (r *ResponseCheckTx) UnmarshalJSON(b []byte) error {
	reader := bytes.NewBuffer(b)
	return jsonpbUnmarshaller.Unmarshal(reader, r)
}

func (r *ResponseDeliverTx) UnmarshalJSON(b []byte) error {
	reader := bytes.NewBuffer(b)
	return jsonpbUnmarshaller.Unmarshal(reader, r)
}

func (r *ResponseQuery) UnmarshalJSON(b []byte) error {
	reader := bytes.NewBuffer(b)
	return jsonpbUnmarshaller.Unmarshal(reader, r)
}

func (r *ResponseCommit) UnmarshalJSON(b []byte) error {
	reader := bytes.NewBuffer(b)
	return jsonpbUnmarshaller.Unmarshal(reader, r)
}

type jsonRoundTripper interface {
	json.Marshaler
	json.Unmarshaler
}
