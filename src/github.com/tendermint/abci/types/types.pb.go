package types

import (
	"encoding/json"
	"fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tmlibs/common"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math"
)

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

const _ = proto.GoGoProtoPackageIsVersion2

type Request struct {
	Value isRequest_Value `protobuf_oneof:"value"`
}

type AppState struct {
	BlockHeight	int64			`json:"block_height,omitempty"`
	AppHash		crypto.Hash		`json:"app_hash,omitempty"`
	TxsHashList	[]crypto.Hash		`json:"txs_hash_list,omitempty"`
	Rewards		[]common.KVPair		`json:"rewards,omitempty"`
	Fee		uint64			`json:"fee,omitempty"`
	BeginBlock	RequestBeginBlock	`json:"beginBlock,omitempty"`
}

func ByteToAppState(appstate []byte) *AppState {
	var a AppState
	json.Unmarshal(appstate, &a)
	return &a
}

func AppStateToByte(appstate *AppState) []byte {
	b, _ := json.Marshal(appstate)
	return b
}

func (m *Request) Reset()			{ *m = Request{} }
func (m *Request) String() string		{ return proto.CompactTextString(m) }
func (*Request) ProtoMessage()			{}
func (*Request) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{0} }

type isRequest_Value interface {
	isRequest_Value()
}

type Request_Echo struct {
	Echo *RequestEcho `protobuf:"bytes,2,opt,name=echo,oneof"`
}
type Request_Flush struct {
	Flush *RequestFlush `protobuf:"bytes,3,opt,name=flush,oneof"`
}
type Request_Info struct {
	Info *RequestInfo `protobuf:"bytes,4,opt,name=info,oneof"`
}
type Request_SetOption struct {
	SetOption *RequestSetOption `protobuf:"bytes,5,opt,name=set_option,json=setOption,oneof"`
}
type Request_InitChain struct {
	InitChain *RequestInitChain `protobuf:"bytes,6,opt,name=init_chain,json=initChain,oneof"`
}
type Request_Query struct {
	Query *RequestQuery `protobuf:"bytes,7,opt,name=query,oneof"`
}
type Request_BeginBlock struct {
	BeginBlock *RequestBeginBlock `protobuf:"bytes,8,opt,name=begin_block,json=beginBlock,oneof"`
}
type Request_CheckTx struct {
	CheckTx *RequestCheckTx `protobuf:"bytes,9,opt,name=check_tx,json=checkTx,oneof"`
}
type Request_DeliverTx struct {
	DeliverTx *RequestDeliverTx `protobuf:"bytes,19,opt,name=deliver_tx,json=deliverTx,oneof"`
}
type Request_EndBlock struct {
	EndBlock *RequestEndBlock `protobuf:"bytes,11,opt,name=end_block,json=endBlock,oneof"`
}
type Request_Commit struct {
	Commit *RequestCommit `protobuf:"bytes,12,opt,name=commit,oneof"`
}

func (*Request_Echo) isRequest_Value()		{}
func (*Request_Flush) isRequest_Value()		{}
func (*Request_Info) isRequest_Value()		{}
func (*Request_SetOption) isRequest_Value()	{}
func (*Request_InitChain) isRequest_Value()	{}
func (*Request_Query) isRequest_Value()		{}
func (*Request_BeginBlock) isRequest_Value()	{}
func (*Request_CheckTx) isRequest_Value()	{}
func (*Request_DeliverTx) isRequest_Value()	{}
func (*Request_EndBlock) isRequest_Value()	{}
func (*Request_Commit) isRequest_Value()	{}

func (m *Request) GetValue() isRequest_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Request) GetEcho() *RequestEcho {
	if x, ok := m.GetValue().(*Request_Echo); ok {
		return x.Echo
	}
	return nil
}

func (m *Request) GetFlush() *RequestFlush {
	if x, ok := m.GetValue().(*Request_Flush); ok {
		return x.Flush
	}
	return nil
}

func (m *Request) GetInfo() *RequestInfo {
	if x, ok := m.GetValue().(*Request_Info); ok {
		return x.Info
	}
	return nil
}

func (m *Request) GetSetOption() *RequestSetOption {
	if x, ok := m.GetValue().(*Request_SetOption); ok {
		return x.SetOption
	}
	return nil
}

func (m *Request) GetInitChain() *RequestInitChain {
	if x, ok := m.GetValue().(*Request_InitChain); ok {
		return x.InitChain
	}
	return nil
}

func (m *Request) GetQuery() *RequestQuery {
	if x, ok := m.GetValue().(*Request_Query); ok {
		return x.Query
	}
	return nil
}

func (m *Request) GetBeginBlock() *RequestBeginBlock {
	if x, ok := m.GetValue().(*Request_BeginBlock); ok {
		return x.BeginBlock
	}
	return nil
}

func (m *Request) GetCheckTx() *RequestCheckTx {
	if x, ok := m.GetValue().(*Request_CheckTx); ok {
		return x.CheckTx
	}
	return nil
}

func (m *Request) GetDeliverTx() *RequestDeliverTx {
	if x, ok := m.GetValue().(*Request_DeliverTx); ok {
		return x.DeliverTx
	}
	return nil
}

func (m *Request) GetEndBlock() *RequestEndBlock {
	if x, ok := m.GetValue().(*Request_EndBlock); ok {
		return x.EndBlock
	}
	return nil
}

func (m *Request) GetCommit() *RequestCommit {
	if x, ok := m.GetValue().(*Request_Commit); ok {
		return x.Commit
	}
	return nil
}

func (*Request) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Request_OneofMarshaler, _Request_OneofUnmarshaler, _Request_OneofSizer, []interface{}{
		(*Request_Echo)(nil),
		(*Request_Flush)(nil),
		(*Request_Info)(nil),
		(*Request_SetOption)(nil),
		(*Request_InitChain)(nil),
		(*Request_Query)(nil),
		(*Request_BeginBlock)(nil),
		(*Request_CheckTx)(nil),
		(*Request_DeliverTx)(nil),
		(*Request_EndBlock)(nil),
		(*Request_Commit)(nil),
	}
}

func _Request_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Request)

	switch x := m.Value.(type) {
	case *Request_Echo:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Echo); err != nil {
			return err
		}
	case *Request_Flush:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Flush); err != nil {
			return err
		}
	case *Request_Info:
		_ = b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Info); err != nil {
			return err
		}
	case *Request_SetOption:
		_ = b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.SetOption); err != nil {
			return err
		}
	case *Request_InitChain:
		_ = b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.InitChain); err != nil {
			return err
		}
	case *Request_Query:
		_ = b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Query); err != nil {
			return err
		}
	case *Request_BeginBlock:
		_ = b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.BeginBlock); err != nil {
			return err
		}
	case *Request_CheckTx:
		_ = b.EncodeVarint(9<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CheckTx); err != nil {
			return err
		}
	case *Request_DeliverTx:
		_ = b.EncodeVarint(19<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.DeliverTx); err != nil {
			return err
		}
	case *Request_EndBlock:
		_ = b.EncodeVarint(11<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.EndBlock); err != nil {
			return err
		}
	case *Request_Commit:
		_ = b.EncodeVarint(12<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Commit); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Request.Value has unexpected type %T", x)
	}
	return nil
}

func _Request_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Request)
	switch tag {
	case 2:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RequestEcho)
		err := b.DecodeMessage(msg)
		m.Value = &Request_Echo{msg}
		return true, err
	case 3:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RequestFlush)
		err := b.DecodeMessage(msg)
		m.Value = &Request_Flush{msg}
		return true, err
	case 4:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RequestInfo)
		err := b.DecodeMessage(msg)
		m.Value = &Request_Info{msg}
		return true, err
	case 5:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RequestSetOption)
		err := b.DecodeMessage(msg)
		m.Value = &Request_SetOption{msg}
		return true, err
	case 6:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RequestInitChain)
		err := b.DecodeMessage(msg)
		m.Value = &Request_InitChain{msg}
		return true, err
	case 7:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RequestQuery)
		err := b.DecodeMessage(msg)
		m.Value = &Request_Query{msg}
		return true, err
	case 8:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RequestBeginBlock)
		err := b.DecodeMessage(msg)
		m.Value = &Request_BeginBlock{msg}
		return true, err
	case 9:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RequestCheckTx)
		err := b.DecodeMessage(msg)
		m.Value = &Request_CheckTx{msg}
		return true, err
	case 19:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RequestDeliverTx)
		err := b.DecodeMessage(msg)
		m.Value = &Request_DeliverTx{msg}
		return true, err
	case 11:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RequestEndBlock)
		err := b.DecodeMessage(msg)
		m.Value = &Request_EndBlock{msg}
		return true, err
	case 12:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RequestCommit)
		err := b.DecodeMessage(msg)
		m.Value = &Request_Commit{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Request_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Request)

	switch x := m.Value.(type) {
	case *Request_Echo:
		s := proto.Size(x.Echo)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Request_Flush:
		s := proto.Size(x.Flush)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Request_Info:
		s := proto.Size(x.Info)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Request_SetOption:
		s := proto.Size(x.SetOption)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Request_InitChain:
		s := proto.Size(x.InitChain)
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Request_Query:
		s := proto.Size(x.Query)
		n += proto.SizeVarint(7<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Request_BeginBlock:
		s := proto.Size(x.BeginBlock)
		n += proto.SizeVarint(8<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Request_CheckTx:
		s := proto.Size(x.CheckTx)
		n += proto.SizeVarint(9<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Request_DeliverTx:
		s := proto.Size(x.DeliverTx)
		n += proto.SizeVarint(19<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Request_EndBlock:
		s := proto.Size(x.EndBlock)
		n += proto.SizeVarint(11<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Request_Commit:
		s := proto.Size(x.Commit)
		n += proto.SizeVarint(12<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type RequestEcho struct {
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *RequestEcho) Reset()				{ *m = RequestEcho{} }
func (m *RequestEcho) String() string			{ return proto.CompactTextString(m) }
func (*RequestEcho) ProtoMessage()			{}
func (*RequestEcho) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{1} }

func (m *RequestEcho) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type RequestFlush struct {
}

func (m *RequestFlush) Reset()				{ *m = RequestFlush{} }
func (m *RequestFlush) String() string			{ return proto.CompactTextString(m) }
func (*RequestFlush) ProtoMessage()			{}
func (*RequestFlush) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{2} }

type RequestInfo struct {
	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (m *RequestInfo) Reset()				{ *m = RequestInfo{} }
func (m *RequestInfo) String() string			{ return proto.CompactTextString(m) }
func (*RequestInfo) ProtoMessage()			{}
func (*RequestInfo) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{3} }

func (m *RequestInfo) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type RequestSetOption struct {
	Key	string	`protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value	string	`protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *RequestSetOption) Reset()			{ *m = RequestSetOption{} }
func (m *RequestSetOption) String() string		{ return proto.CompactTextString(m) }
func (*RequestSetOption) ProtoMessage()			{}
func (*RequestSetOption) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{4} }

func (m *RequestSetOption) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *RequestSetOption) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type RequestInitChain struct {
	Validators	[]Validator	`protobuf:"bytes,1,rep,name=validators" json:"validators"`
	ChainId		string		`protobuf:"bytes,2,rep,name=chain_id" json:"chain_id"`
	AppStateBytes	[]byte		`protobuf:"bytes,3,opt,name=app_state_bytes,json=appStateBytes,proto3" json:"app_state_bytes,omitempty"`
}

func (m *RequestInitChain) Reset()			{ *m = RequestInitChain{} }
func (m *RequestInitChain) String() string		{ return proto.CompactTextString(m) }
func (*RequestInitChain) ProtoMessage()			{}
func (*RequestInitChain) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{5} }

func (m *RequestInitChain) GetValidators() []Validator {
	if m != nil {
		return m.Validators
	}
	return nil
}

func (m *RequestInitChain) GetAppStateBytes() []byte {
	if m != nil {
		return m.AppStateBytes
	}
	return nil
}

type RequestQuery struct {
	Data	[]byte	`protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Path	string	`protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	Height	int64	`protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	Prove	bool	`protobuf:"varint,4,opt,name=prove,proto3" json:"prove,omitempty"`
}

func (m *RequestQuery) Reset()				{ *m = RequestQuery{} }
func (m *RequestQuery) String() string			{ return proto.CompactTextString(m) }
func (*RequestQuery) ProtoMessage()			{}
func (*RequestQuery) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{6} }

func (m *RequestQuery) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *RequestQuery) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *RequestQuery) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *RequestQuery) GetProve() bool {
	if m != nil {
		return m.Prove
	}
	return false
}

type RequestBeginBlock struct {
	Hash			[]byte		`protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Header			Header		`protobuf:"bytes,2,opt,name=header" json:"header"`
	AbsentValidators	[]int32		`protobuf:"varint,3,rep,packed,name=absent_validators,json=absentValidators" json:"absent_validators,omitempty"`
	ByzantineValidators	[]Evidence	`protobuf:"bytes,4,rep,name=byzantine_validators,json=byzantineValidators" json:"byzantine_validators"`
}

func (m *RequestBeginBlock) Reset()			{ *m = RequestBeginBlock{} }
func (m *RequestBeginBlock) String() string		{ return proto.CompactTextString(m) }
func (*RequestBeginBlock) ProtoMessage()		{}
func (*RequestBeginBlock) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{7} }

func (m *RequestBeginBlock) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *RequestBeginBlock) GetHeader() Header {
	if m != nil {
		return m.Header
	}
	return Header{}
}

func (m *RequestBeginBlock) GetAbsentValidators() []int32 {
	if m != nil {
		return m.AbsentValidators
	}
	return nil
}

func (m *RequestBeginBlock) GetByzantineValidators() []Evidence {
	if m != nil {
		return m.ByzantineValidators
	}
	return nil
}

type RequestCheckTx struct {
	Tx []byte `protobuf:"bytes,1,opt,name=tx,proto3" json:"tx,omitempty"`
}

func (m *RequestCheckTx) Reset()			{ *m = RequestCheckTx{} }
func (m *RequestCheckTx) String() string		{ return proto.CompactTextString(m) }
func (*RequestCheckTx) ProtoMessage()			{}
func (*RequestCheckTx) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{8} }

func (m *RequestCheckTx) GetTx() []byte {
	if m != nil {
		return m.Tx
	}
	return nil
}

type RequestDeliverTx struct {
	Tx []byte `protobuf:"bytes,1,opt,name=tx,proto3" json:"tx,omitempty"`
}

func (m *RequestDeliverTx) Reset()			{ *m = RequestDeliverTx{} }
func (m *RequestDeliverTx) String() string		{ return proto.CompactTextString(m) }
func (*RequestDeliverTx) ProtoMessage()			{}
func (*RequestDeliverTx) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{9} }

func (m *RequestDeliverTx) GetTx() []byte {
	if m != nil {
		return m.Tx
	}
	return nil
}

type RequestEndBlock struct {
	Height int64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *RequestEndBlock) Reset()			{ *m = RequestEndBlock{} }
func (m *RequestEndBlock) String() string		{ return proto.CompactTextString(m) }
func (*RequestEndBlock) ProtoMessage()			{}
func (*RequestEndBlock) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{10} }

func (m *RequestEndBlock) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type RequestCommit struct {
}

func (m *RequestCommit) Reset()				{ *m = RequestCommit{} }
func (m *RequestCommit) String() string			{ return proto.CompactTextString(m) }
func (*RequestCommit) ProtoMessage()			{}
func (*RequestCommit) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{11} }

type Response struct {
	Value isResponse_Value `protobuf_oneof:"value"`
}

func (m *Response) Reset()			{ *m = Response{} }
func (m *Response) String() string		{ return proto.CompactTextString(m) }
func (*Response) ProtoMessage()			{}
func (*Response) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{12} }

type isResponse_Value interface {
	isResponse_Value()
}

type Response_Exception struct {
	Exception *ResponseException `protobuf:"bytes,1,opt,name=exception,oneof"`
}
type Response_Echo struct {
	Echo *ResponseEcho `protobuf:"bytes,2,opt,name=echo,oneof"`
}
type Response_Flush struct {
	Flush *ResponseFlush `protobuf:"bytes,3,opt,name=flush,oneof"`
}
type Response_Info struct {
	Info *ResponseInfo `protobuf:"bytes,4,opt,name=info,oneof"`
}
type Response_SetOption struct {
	SetOption *ResponseSetOption `protobuf:"bytes,5,opt,name=set_option,json=setOption,oneof"`
}
type Response_InitChain struct {
	InitChain *ResponseInitChain `protobuf:"bytes,6,opt,name=init_chain,json=initChain,oneof"`
}
type Response_Query struct {
	Query *ResponseQuery `protobuf:"bytes,7,opt,name=query,oneof"`
}
type Response_BeginBlock struct {
	BeginBlock *ResponseBeginBlock `protobuf:"bytes,8,opt,name=begin_block,json=beginBlock,oneof"`
}
type Response_CheckTx struct {
	CheckTx *ResponseCheckTx `protobuf:"bytes,9,opt,name=check_tx,json=checkTx,oneof"`
}
type Response_DeliverTx struct {
	DeliverTx *ResponseDeliverTx `protobuf:"bytes,10,opt,name=deliver_tx,json=deliverTx,oneof"`
}
type Response_EndBlock struct {
	EndBlock *ResponseEndBlock `protobuf:"bytes,11,opt,name=end_block,json=endBlock,oneof"`
}
type Response_Commit struct {
	Commit *ResponseCommit `protobuf:"bytes,12,opt,name=commit,oneof"`
}

func (*Response_Exception) isResponse_Value()	{}
func (*Response_Echo) isResponse_Value()	{}
func (*Response_Flush) isResponse_Value()	{}
func (*Response_Info) isResponse_Value()	{}
func (*Response_SetOption) isResponse_Value()	{}
func (*Response_InitChain) isResponse_Value()	{}
func (*Response_Query) isResponse_Value()	{}
func (*Response_BeginBlock) isResponse_Value()	{}
func (*Response_CheckTx) isResponse_Value()	{}
func (*Response_DeliverTx) isResponse_Value()	{}
func (*Response_EndBlock) isResponse_Value()	{}
func (*Response_Commit) isResponse_Value()	{}

func (m *Response) GetValue() isResponse_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Response) GetException() *ResponseException {
	if x, ok := m.GetValue().(*Response_Exception); ok {
		return x.Exception
	}
	return nil
}

func (m *Response) GetEcho() *ResponseEcho {
	if x, ok := m.GetValue().(*Response_Echo); ok {
		return x.Echo
	}
	return nil
}

func (m *Response) GetFlush() *ResponseFlush {
	if x, ok := m.GetValue().(*Response_Flush); ok {
		return x.Flush
	}
	return nil
}

func (m *Response) GetInfo() *ResponseInfo {
	if x, ok := m.GetValue().(*Response_Info); ok {
		return x.Info
	}
	return nil
}

func (m *Response) GetSetOption() *ResponseSetOption {
	if x, ok := m.GetValue().(*Response_SetOption); ok {
		return x.SetOption
	}
	return nil
}

func (m *Response) GetInitChain() *ResponseInitChain {
	if x, ok := m.GetValue().(*Response_InitChain); ok {
		return x.InitChain
	}
	return nil
}

func (m *Response) GetQuery() *ResponseQuery {
	if x, ok := m.GetValue().(*Response_Query); ok {
		return x.Query
	}
	return nil
}

func (m *Response) GetBeginBlock() *ResponseBeginBlock {
	if x, ok := m.GetValue().(*Response_BeginBlock); ok {
		return x.BeginBlock
	}
	return nil
}

func (m *Response) GetCheckTx() *ResponseCheckTx {
	if x, ok := m.GetValue().(*Response_CheckTx); ok {
		return x.CheckTx
	}
	return nil
}

func (m *Response) GetDeliverTx() *ResponseDeliverTx {
	if x, ok := m.GetValue().(*Response_DeliverTx); ok {
		return x.DeliverTx
	}
	return nil
}

func (m *Response) GetEndBlock() *ResponseEndBlock {
	if x, ok := m.GetValue().(*Response_EndBlock); ok {
		return x.EndBlock
	}
	return nil
}

func (m *Response) GetCommit() *ResponseCommit {
	if x, ok := m.GetValue().(*Response_Commit); ok {
		return x.Commit
	}
	return nil
}

func (*Response) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Response_OneofMarshaler, _Response_OneofUnmarshaler, _Response_OneofSizer, []interface{}{
		(*Response_Exception)(nil),
		(*Response_Echo)(nil),
		(*Response_Flush)(nil),
		(*Response_Info)(nil),
		(*Response_SetOption)(nil),
		(*Response_InitChain)(nil),
		(*Response_Query)(nil),
		(*Response_BeginBlock)(nil),
		(*Response_CheckTx)(nil),
		(*Response_DeliverTx)(nil),
		(*Response_EndBlock)(nil),
		(*Response_Commit)(nil),
	}
}

func _Response_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Response)

	switch x := m.Value.(type) {
	case *Response_Exception:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Exception); err != nil {
			return err
		}
	case *Response_Echo:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Echo); err != nil {
			return err
		}
	case *Response_Flush:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Flush); err != nil {
			return err
		}
	case *Response_Info:
		_ = b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Info); err != nil {
			return err
		}
	case *Response_SetOption:
		_ = b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.SetOption); err != nil {
			return err
		}
	case *Response_InitChain:
		_ = b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.InitChain); err != nil {
			return err
		}
	case *Response_Query:
		_ = b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Query); err != nil {
			return err
		}
	case *Response_BeginBlock:
		_ = b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.BeginBlock); err != nil {
			return err
		}
	case *Response_CheckTx:
		_ = b.EncodeVarint(9<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CheckTx); err != nil {
			return err
		}
	case *Response_DeliverTx:
		_ = b.EncodeVarint(10<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.DeliverTx); err != nil {
			return err
		}
	case *Response_EndBlock:
		_ = b.EncodeVarint(11<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.EndBlock); err != nil {
			return err
		}
	case *Response_Commit:
		_ = b.EncodeVarint(12<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Commit); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Response.Value has unexpected type %T", x)
	}
	return nil
}

func _Response_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Response)
	switch tag {
	case 1:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ResponseException)
		err := b.DecodeMessage(msg)
		m.Value = &Response_Exception{msg}
		return true, err
	case 2:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ResponseEcho)
		err := b.DecodeMessage(msg)
		m.Value = &Response_Echo{msg}
		return true, err
	case 3:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ResponseFlush)
		err := b.DecodeMessage(msg)
		m.Value = &Response_Flush{msg}
		return true, err
	case 4:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ResponseInfo)
		err := b.DecodeMessage(msg)
		m.Value = &Response_Info{msg}
		return true, err
	case 5:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ResponseSetOption)
		err := b.DecodeMessage(msg)
		m.Value = &Response_SetOption{msg}
		return true, err
	case 6:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ResponseInitChain)
		err := b.DecodeMessage(msg)
		m.Value = &Response_InitChain{msg}
		return true, err
	case 7:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ResponseQuery)
		err := b.DecodeMessage(msg)
		m.Value = &Response_Query{msg}
		return true, err
	case 8:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ResponseBeginBlock)
		err := b.DecodeMessage(msg)
		m.Value = &Response_BeginBlock{msg}
		return true, err
	case 9:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ResponseCheckTx)
		err := b.DecodeMessage(msg)
		m.Value = &Response_CheckTx{msg}
		return true, err
	case 10:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ResponseDeliverTx)
		err := b.DecodeMessage(msg)
		m.Value = &Response_DeliverTx{msg}
		return true, err
	case 11:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ResponseEndBlock)
		err := b.DecodeMessage(msg)
		m.Value = &Response_EndBlock{msg}
		return true, err
	case 12:
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ResponseCommit)
		err := b.DecodeMessage(msg)
		m.Value = &Response_Commit{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Response_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Response)

	switch x := m.Value.(type) {
	case *Response_Exception:
		s := proto.Size(x.Exception)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Response_Echo:
		s := proto.Size(x.Echo)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Response_Flush:
		s := proto.Size(x.Flush)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Response_Info:
		s := proto.Size(x.Info)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Response_SetOption:
		s := proto.Size(x.SetOption)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Response_InitChain:
		s := proto.Size(x.InitChain)
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Response_Query:
		s := proto.Size(x.Query)
		n += proto.SizeVarint(7<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Response_BeginBlock:
		s := proto.Size(x.BeginBlock)
		n += proto.SizeVarint(8<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Response_CheckTx:
		s := proto.Size(x.CheckTx)
		n += proto.SizeVarint(9<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Response_DeliverTx:
		s := proto.Size(x.DeliverTx)
		n += proto.SizeVarint(10<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Response_EndBlock:
		s := proto.Size(x.EndBlock)
		n += proto.SizeVarint(11<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Response_Commit:
		s := proto.Size(x.Commit)
		n += proto.SizeVarint(12<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ResponseException struct {
	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (m *ResponseException) Reset()			{ *m = ResponseException{} }
func (m *ResponseException) String() string		{ return proto.CompactTextString(m) }
func (*ResponseException) ProtoMessage()		{}
func (*ResponseException) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{13} }

func (m *ResponseException) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type ResponseEcho struct {
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *ResponseEcho) Reset()				{ *m = ResponseEcho{} }
func (m *ResponseEcho) String() string			{ return proto.CompactTextString(m) }
func (*ResponseEcho) ProtoMessage()			{}
func (*ResponseEcho) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{14} }

func (m *ResponseEcho) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type ResponseFlush struct {
}

func (m *ResponseFlush) Reset()				{ *m = ResponseFlush{} }
func (m *ResponseFlush) String() string			{ return proto.CompactTextString(m) }
func (*ResponseFlush) ProtoMessage()			{}
func (*ResponseFlush) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{15} }

type ResponseInfo struct {
	Data		string	`protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Version		string	`protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	LastBlockHeight	int64	`protobuf:"varint,3,opt,name=last_block_height,json=lastBlockHeight,proto3" json:"last_block_height,omitempty"`
	LastAppState	[]byte	`protobuf:"bytes,4,opt,name=last_app_state,json=lastAppState,proto3" json:"last_app_state,omitempty"`
}

func (m *ResponseInfo) Reset()				{ *m = ResponseInfo{} }
func (m *ResponseInfo) String() string			{ return proto.CompactTextString(m) }
func (*ResponseInfo) ProtoMessage()			{}
func (*ResponseInfo) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{16} }

func (m *ResponseInfo) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *ResponseInfo) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *ResponseInfo) GetLastBlockHeight() int64 {
	if m != nil {
		return m.LastBlockHeight
	}
	return 0
}

func (m *ResponseInfo) GetLastBlockAppHash() []byte {
	if m != nil {
		appState := ByteToAppState(m.LastAppState)
		return appState.AppHash
	}
	return nil
}

type ResponseSetOption struct {
	Code	uint32	`protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`

	Log	string	`protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	Info	string	`protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
}

func (m *ResponseSetOption) Reset()			{ *m = ResponseSetOption{} }
func (m *ResponseSetOption) String() string		{ return proto.CompactTextString(m) }
func (*ResponseSetOption) ProtoMessage()		{}
func (*ResponseSetOption) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{17} }

func (m *ResponseSetOption) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *ResponseSetOption) GetLog() string {
	if m != nil {
		return m.Log
	}
	return ""
}

func (m *ResponseSetOption) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

type ResponseInitChain struct {
	Code		uint32	`protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Log		string	`protobuf:"bytes,2,opt,name=log,proto3" json:"log,omitempty"`
	GenAppState	[]byte	`protobuf:"bytes,3,opt,name=gen_app_state" json:"gen_app_state,omitempty"`
}

func (m *ResponseInitChain) Reset()			{ *m = ResponseInitChain{} }
func (m *ResponseInitChain) String() string		{ return proto.CompactTextString(m) }
func (*ResponseInitChain) ProtoMessage()		{}
func (*ResponseInitChain) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{18} }

type ResponseQuery struct {
	Code	uint32	`protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`

	Log	string	`protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	Info	string	`protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	Index	int64	`protobuf:"varint,5,opt,name=index,proto3" json:"index,omitempty"`
	Key	[]byte	`protobuf:"bytes,6,opt,name=key,proto3" json:"key,omitempty"`
	Value	[]byte	`protobuf:"bytes,7,opt,name=value,proto3" json:"value,omitempty"`
	Proof	[]byte	`protobuf:"bytes,8,opt,name=proof,proto3" json:"proof,omitempty"`
	Height	int64	`protobuf:"varint,9,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *ResponseQuery) Reset()				{ *m = ResponseQuery{} }
func (m *ResponseQuery) String() string			{ return proto.CompactTextString(m) }
func (*ResponseQuery) ProtoMessage()			{}
func (*ResponseQuery) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{19} }

func (m *ResponseQuery) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *ResponseQuery) GetLog() string {
	if m != nil {
		return m.Log
	}
	return ""
}

func (m *ResponseQuery) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

func (m *ResponseQuery) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *ResponseQuery) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *ResponseQuery) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *ResponseQuery) GetProof() []byte {
	if m != nil {
		return m.Proof
	}
	return nil
}

func (m *ResponseQuery) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type ResponseBeginBlock struct {
	Code	uint32	`protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Log	string	`protobuf:"bytes,2,opt,name=log,proto3" json:"log,omitempty"`
}

func (m *ResponseBeginBlock) Reset()			{ *m = ResponseBeginBlock{} }
func (m *ResponseBeginBlock) String() string		{ return proto.CompactTextString(m) }
func (*ResponseBeginBlock) ProtoMessage()		{}
func (*ResponseBeginBlock) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{20} }

type ResponseCheckTx struct {
	Code		uint32		`protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Data		string		`protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Log		string		`protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	Info		string		`protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	GasLimit	uint64		`protobuf:"varint,5,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit,omitempty"`
	GasUsed		uint64		`protobuf:"varint,6,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	Fee		uint64		`protobuf:"varint,7,opt,name=fee,json=fee,proto3" json:"fee,omitempty"`
	Tags		[]common.KVPair	`protobuf:"bytes,8,rep,name=tags" json:"tags,omitempty"`
	TxHash		common.HexBytes	`protobuf:"bytes,9,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	Height		int64		`protobuf:"varint,10,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *ResponseCheckTx) Reset()			{ *m = ResponseCheckTx{} }
func (m *ResponseCheckTx) String() string		{ return proto.CompactTextString(m) }
func (*ResponseCheckTx) ProtoMessage()			{}
func (*ResponseCheckTx) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{21} }

func (m *ResponseCheckTx) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *ResponseCheckTx) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *ResponseCheckTx) GetLog() string {
	if m != nil {
		return m.Log
	}
	return ""
}

func (m *ResponseCheckTx) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

func (m *ResponseCheckTx) GetGasWanted() uint64 {
	if m != nil {
		return m.GasLimit
	}
	return 0
}

func (m *ResponseCheckTx) GetGasUsed() uint64 {
	if m != nil {
		return m.GasUsed
	}
	return 0
}

func (m *ResponseCheckTx) GetTags() []common.KVPair {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *ResponseCheckTx) GetFee() uint64 {
	if m != nil {
		return m.Fee
	}
	return 0
}

type ResponseDeliverTx struct {
	Code		uint32		`protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Data		string		`protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Log		string		`protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	Info		string		`protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	GasLimit	uint64		`protobuf:"varint,5,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit,omitempty"`
	GasUsed		uint64		`protobuf:"varint,6,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	Fee		uint64		`protobuf:"varint,7,opt,name=fee,json=fee,proto3" json:"fee,omitempty"`
	Tags		[]common.KVPair	`protobuf:"bytes,8,rep,name=tags" json:"tags,omitempty"`
	TxHash		common.HexBytes	`protobuf:"bytes,9,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	Height		int64		`protobuf:"varint,10,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *ResponseDeliverTx) Reset()			{ *m = ResponseDeliverTx{} }
func (m *ResponseDeliverTx) String() string		{ return proto.CompactTextString(m) }
func (*ResponseDeliverTx) ProtoMessage()		{}
func (*ResponseDeliverTx) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{22} }

func (m *ResponseDeliverTx) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *ResponseDeliverTx) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *ResponseDeliverTx) GetLog() string {
	if m != nil {
		return m.Log
	}
	return ""
}

func (m *ResponseDeliverTx) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

func (m *ResponseDeliverTx) GetGasLimit() uint64 {
	if m != nil {
		return m.GasLimit
	}
	return 0
}

func (m *ResponseDeliverTx) GetGasUsed() uint64 {
	if m != nil {
		return m.GasUsed
	}
	return 0
}

func (m *ResponseDeliverTx) GetTags() []common.KVPair {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *ResponseDeliverTx) GetFee() uint64 {
	if m != nil {
		return m.Fee
	}
	return 0
}

type ResponseEndBlock struct {
	ValidatorUpdates	[]Validator		`protobuf:"bytes,1,rep,name=validator_updates,json=validatorUpdates" json:"validator_updates"`
	ConsensusParamUpdates	*ConsensusParams	`protobuf:"bytes,2,opt,name=consensus_param_updates,json=consensusParamUpdates" json:"consensus_param_updates,omitempty"`
}

func (m *ResponseEndBlock) Reset()			{ *m = ResponseEndBlock{} }
func (m *ResponseEndBlock) String() string		{ return proto.CompactTextString(m) }
func (*ResponseEndBlock) ProtoMessage()			{}
func (*ResponseEndBlock) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{23} }

func (m *ResponseEndBlock) GetValidatorUpdates() []Validator {
	if m != nil {
		return m.ValidatorUpdates
	}
	return nil
}

func (m *ResponseEndBlock) GetConsensusParamUpdates() *ConsensusParams {
	if m != nil {
		return m.ConsensusParamUpdates
	}
	return nil
}

type ResponseCommit struct {
	AppState []byte `protobuf:"bytes,1,opt,name=appstate,proto3" json:"app_state,omitempty"`
}

func (m *ResponseCommit) Reset()			{ *m = ResponseCommit{} }
func (m *ResponseCommit) String() string		{ return proto.CompactTextString(m) }
func (*ResponseCommit) ProtoMessage()			{}
func (*ResponseCommit) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{24} }
func (m *ResponseCommit) GetReward() []common.KVPair {
	if m != nil {
		appState := ByteToAppState(m.AppState)
		return appState.Rewards
	}
	return nil
}
func (m *ResponseCommit) GetFee() uint64 {
	if m != nil {
		appState := ByteToAppState(m.AppState)
		return appState.Fee
	}
	return 0
}
func (m *ResponseCommit) GetLastAppHash() crypto.Hash {
	if m != nil {
		appState := ByteToAppState(m.AppState)
		return appState.AppHash
	}
	return nil
}
func (m *ResponseCommit) GetHashLists() []crypto.Hash {
	if m != nil {
		appState := ByteToAppState(m.AppState)
		return appState.TxsHashList
	}
	return nil
}

type ConsensusParams struct {
	BlockSize	*BlockSize	`protobuf:"bytes,1,opt,name=block_size,json=blockSize" json:"block_size,omitempty"`
	TxSize		*TxSize		`protobuf:"bytes,2,opt,name=tx_size,json=txSize" json:"tx_size,omitempty"`
	BlockGossip	*BlockGossip	`protobuf:"bytes,3,opt,name=block_gossip,json=blockGossip" json:"block_gossip,omitempty"`
}

func (m *ConsensusParams) Reset()			{ *m = ConsensusParams{} }
func (m *ConsensusParams) String() string		{ return proto.CompactTextString(m) }
func (*ConsensusParams) ProtoMessage()			{}
func (*ConsensusParams) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{25} }

func (m *ConsensusParams) GetBlockSize() *BlockSize {
	if m != nil {
		return m.BlockSize
	}
	return nil
}

func (m *ConsensusParams) GetTxSize() *TxSize {
	if m != nil {
		return m.TxSize
	}
	return nil
}

func (m *ConsensusParams) GetBlockGossip() *BlockGossip {
	if m != nil {
		return m.BlockGossip
	}
	return nil
}

type BlockSize struct {
	MaxBytes	int32	`protobuf:"varint,1,opt,name=max_bytes,json=maxBytes,proto3" json:"max_bytes,omitempty"`
	MaxTxs		int32	`protobuf:"varint,2,opt,name=max_txs,json=maxTxs,proto3" json:"max_txs,omitempty"`
	MaxGas		int64	`protobuf:"varint,3,opt,name=max_gas,json=maxGas,proto3" json:"max_gas,omitempty"`
}

func (m *BlockSize) Reset()			{ *m = BlockSize{} }
func (m *BlockSize) String() string		{ return proto.CompactTextString(m) }
func (*BlockSize) ProtoMessage()		{}
func (*BlockSize) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{26} }

func (m *BlockSize) GetMaxBytes() int32 {
	if m != nil {
		return m.MaxBytes
	}
	return 0
}

func (m *BlockSize) GetMaxTxs() int32 {
	if m != nil {
		return m.MaxTxs
	}
	return 0
}

func (m *BlockSize) GetMaxGas() int64 {
	if m != nil {
		return m.MaxGas
	}
	return 0
}

type TxSize struct {
	MaxBytes	int32	`protobuf:"varint,1,opt,name=max_bytes,json=maxBytes,proto3" json:"max_bytes,omitempty"`
	MaxGas		int64	`protobuf:"varint,2,opt,name=max_gas,json=maxGas,proto3" json:"max_gas,omitempty"`
}

func (m *TxSize) Reset()			{ *m = TxSize{} }
func (m *TxSize) String() string		{ return proto.CompactTextString(m) }
func (*TxSize) ProtoMessage()			{}
func (*TxSize) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{27} }

func (m *TxSize) GetMaxBytes() int32 {
	if m != nil {
		return m.MaxBytes
	}
	return 0
}

func (m *TxSize) GetMaxGas() int64 {
	if m != nil {
		return m.MaxGas
	}
	return 0
}

type BlockGossip struct {
	BlockPartSizeBytes int32 `protobuf:"varint,1,opt,name=block_part_size_bytes,json=blockPartSizeBytes,proto3" json:"block_part_size_bytes,omitempty"`
}

func (m *BlockGossip) Reset()				{ *m = BlockGossip{} }
func (m *BlockGossip) String() string			{ return proto.CompactTextString(m) }
func (*BlockGossip) ProtoMessage()			{}
func (*BlockGossip) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{28} }

func (m *BlockGossip) GetBlockPartSizeBytes() int32 {
	if m != nil {
		return m.BlockPartSizeBytes
	}
	return 0
}

type Allocation struct {
	Addr	string	`protobuf:"bytes,1,opt,name=addr,json=addr,proto3" json:"addr"`
	Fee	uint64	`protobuf:"varint,2,opt,name=fee,json=fee,proto3" json:"fee"`
}

type Header struct {
	ChainID		string		`protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Height		int64		`protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	Time		int64		`protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
	NumTxs		int32		`protobuf:"varint,4,opt,name=num_txs,json=numTxs,proto3" json:"num_txs,omitempty"`
	LastBlockID	BlockID		`protobuf:"bytes,5,opt,name=last_block_id,json=lastBlockId" json:"last_block_id"`
	LastCommitHash	[]byte		`protobuf:"bytes,6,opt,name=last_commit_hash,json=lastCommitHash,proto3" json:"last_commit_hash,omitempty"`
	DataHash	[]byte		`protobuf:"bytes,7,opt,name=data_hash,json=dataHash,proto3" json:"data_hash,omitempty"`
	ValidatorsHash	[]byte		`protobuf:"bytes,8,opt,name=validators_hash,json=validatorsHash,proto3" json:"validators_hash,omitempty"`
	LastAppHash	[]byte		`protobuf:"bytes,9,opt,name=last_app_hash,json=lastAppHash,proto3" json:"last_app_hash,omitempty"`
	LastFee		uint64		`protobuf:"varint,10,opt,name=last_fee,json=lastFee,proto3" json:"last_fee,omitempty"`
	LastAllocation	[]Allocation	`protobuf:"bytes,11,rep,name=last_allocation" json:"last_allocation,omitempty"`
	ProposerAddress	string		`protobuf:"bytes,12,opt,name=proposer_address,json=proposerAddress,proto3" json:"proposer_address,omitempty"`
	RewardAddress	string		`protobuf:"bytes,13,opt,name=reward_address,json=rewardAddress,proto3" json:"reward_address,omitempty"`
	RandomeOfBlock	[]byte		`protobuf:"bytes,14,opt,name=random_of_block,json=randomOfBlock,proto3" json:"random_of_block,omitempty"`
}

func (m *Header) Reset()			{ *m = Header{} }
func (m *Header) String() string		{ return proto.CompactTextString(m) }
func (*Header) ProtoMessage()			{}
func (*Header) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{29} }

func (m *Header) GetChainID() string {
	if m != nil {
		return m.ChainID
	}
	return ""
}

func (m *Header) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Header) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *Header) GetNumTxs() int32 {
	if m != nil {
		return m.NumTxs
	}
	return 0
}

func (m *Header) GetLastBlockID() BlockID {
	if m != nil {
		return m.LastBlockID
	}
	return BlockID{}
}

func (m *Header) GetLastCommitHash() []byte {
	if m != nil {
		return m.LastCommitHash
	}
	return nil
}

func (m *Header) GetDataHash() []byte {
	if m != nil {
		return m.DataHash
	}
	return nil
}

func (m *Header) GetValidatorsHash() []byte {
	if m != nil {
		return m.ValidatorsHash
	}
	return nil
}

func (m *Header) GetLastAppHash() []byte {
	if m != nil {
		return m.LastAppHash
	}
	return nil
}

type BlockID struct {
	Hash	[]byte		`protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Parts	PartSetHeader	`protobuf:"bytes,2,opt,name=parts" json:"parts"`
}

func (m *BlockID) Reset()			{ *m = BlockID{} }
func (m *BlockID) String() string		{ return proto.CompactTextString(m) }
func (*BlockID) ProtoMessage()			{}
func (*BlockID) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{30} }

func (m *BlockID) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *BlockID) GetParts() PartSetHeader {
	if m != nil {
		return m.Parts
	}
	return PartSetHeader{}
}

type PartSetHeader struct {
	Total	int32	`protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Hash	[]byte	`protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (m *PartSetHeader) Reset()				{ *m = PartSetHeader{} }
func (m *PartSetHeader) String() string			{ return proto.CompactTextString(m) }
func (*PartSetHeader) ProtoMessage()			{}
func (*PartSetHeader) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{31} }

func (m *PartSetHeader) GetTotal() int32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *PartSetHeader) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

type Validator struct {
	PubKey		[]byte	`protobuf:"bytes,1,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
	Power		uint64	`protobuf:"varint,2,opt,name=power,proto3" json:"power,omitempty"`
	RewardAddr	string	`protobuf:"bytes,3,opt,name=reward_addr,json=rewardAddr,proto3" json:"reward_addr"`
	Name		string	`protobuf:"bytes,4,opt,name=name,json=name,proto3" json:"name"`
}

func (m *Validator) Reset()			{ *m = Validator{} }
func (m *Validator) String() string		{ return proto.CompactTextString(m) }
func (*Validator) ProtoMessage()		{}
func (*Validator) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{32} }

func (m *Validator) GetPubKey() []byte {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *Validator) GetPower() uint64 {
	if m != nil {
		return m.Power
	}
	return 0
}

type Evidence struct {
	PubKey	string	`protobuf:"bytes,1,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
	Height	int64	`protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *Evidence) Reset()			{ *m = Evidence{} }
func (m *Evidence) String() string		{ return proto.CompactTextString(m) }
func (*Evidence) ProtoMessage()			{}
func (*Evidence) Descriptor() ([]byte, []int)	{ return fileDescriptorTypes, []int{33} }

func (m *Evidence) GetPubKey() string {
	if m != nil {
		return m.PubKey
	}
	return ""
}

func (m *Evidence) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func init() {
	proto.RegisterType((*Request)(nil), "types.Request")
	proto.RegisterType((*RequestEcho)(nil), "types.RequestEcho")
	proto.RegisterType((*RequestFlush)(nil), "types.RequestFlush")
	proto.RegisterType((*RequestInfo)(nil), "types.RequestInfo")
	proto.RegisterType((*RequestSetOption)(nil), "types.RequestSetOption")
	proto.RegisterType((*RequestInitChain)(nil), "types.RequestInitChain")
	proto.RegisterType((*RequestQuery)(nil), "types.RequestQuery")
	proto.RegisterType((*RequestBeginBlock)(nil), "types.RequestBeginBlock")
	proto.RegisterType((*RequestCheckTx)(nil), "types.RequestCheckTx")
	proto.RegisterType((*RequestDeliverTx)(nil), "types.RequestDeliverTx")
	proto.RegisterType((*RequestEndBlock)(nil), "types.RequestEndBlock")
	proto.RegisterType((*RequestCommit)(nil), "types.RequestCommit")
	proto.RegisterType((*Response)(nil), "types.Response")
	proto.RegisterType((*ResponseException)(nil), "types.ResponseException")
	proto.RegisterType((*ResponseEcho)(nil), "types.ResponseEcho")
	proto.RegisterType((*ResponseFlush)(nil), "types.ResponseFlush")
	proto.RegisterType((*ResponseInfo)(nil), "types.ResponseInfo")
	proto.RegisterType((*ResponseSetOption)(nil), "types.ResponseSetOption")
	proto.RegisterType((*ResponseInitChain)(nil), "types.ResponseInitChain")
	proto.RegisterType((*ResponseQuery)(nil), "types.ResponseQuery")
	proto.RegisterType((*ResponseBeginBlock)(nil), "types.ResponseBeginBlock")
	proto.RegisterType((*ResponseCheckTx)(nil), "types.ResponseCheckTx")
	proto.RegisterType((*ResponseDeliverTx)(nil), "types.ResponseDeliverTx")
	proto.RegisterType((*ResponseEndBlock)(nil), "types.ResponseEndBlock")
	proto.RegisterType((*ResponseCommit)(nil), "types.ResponseCommit")
	proto.RegisterType((*ConsensusParams)(nil), "types.ConsensusParams")
	proto.RegisterType((*BlockSize)(nil), "types.BlockSize")
	proto.RegisterType((*TxSize)(nil), "types.TxSize")
	proto.RegisterType((*BlockGossip)(nil), "types.BlockGossip")
	proto.RegisterType((*Header)(nil), "types.Header")
	proto.RegisterType((*BlockID)(nil), "types.BlockID")
	proto.RegisterType((*PartSetHeader)(nil), "types.PartSetHeader")
	proto.RegisterType((*Validator)(nil), "types.Validator")
	proto.RegisterType((*Evidence)(nil), "types.Evidence")
}

var _ context.Context
var _ grpc.ClientConn

const _ = grpc.SupportPackageIsVersion4

type ABCIApplicationClient interface {
	Echo(ctx context.Context, in *RequestEcho, opts ...grpc.CallOption) (*ResponseEcho, error)
	Flush(ctx context.Context, in *RequestFlush, opts ...grpc.CallOption) (*ResponseFlush, error)
	Info(ctx context.Context, in *RequestInfo, opts ...grpc.CallOption) (*ResponseInfo, error)
	SetOption(ctx context.Context, in *RequestSetOption, opts ...grpc.CallOption) (*ResponseSetOption, error)
	DeliverTx(ctx context.Context, in *RequestDeliverTx, opts ...grpc.CallOption) (*ResponseDeliverTx, error)
	CheckTx(ctx context.Context, in *RequestCheckTx, opts ...grpc.CallOption) (*ResponseCheckTx, error)
	Query(ctx context.Context, in *RequestQuery, opts ...grpc.CallOption) (*ResponseQuery, error)
	Commit(ctx context.Context, in *RequestCommit, opts ...grpc.CallOption) (*ResponseCommit, error)
	InitChain(ctx context.Context, in *RequestInitChain, opts ...grpc.CallOption) (*ResponseInitChain, error)
	BeginBlock(ctx context.Context, in *RequestBeginBlock, opts ...grpc.CallOption) (*ResponseBeginBlock, error)
	EndBlock(ctx context.Context, in *RequestEndBlock, opts ...grpc.CallOption) (*ResponseEndBlock, error)
}

type aBCIApplicationClient struct {
	cc *grpc.ClientConn
}

func NewABCIApplicationClient(cc *grpc.ClientConn) ABCIApplicationClient {
	return &aBCIApplicationClient{cc}
}

func (c *aBCIApplicationClient) Echo(ctx context.Context, in *RequestEcho, opts ...grpc.CallOption) (*ResponseEcho, error) {
	out := new(ResponseEcho)
	err := grpc.Invoke(ctx, "/types.ABCIApplication/Echo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aBCIApplicationClient) Flush(ctx context.Context, in *RequestFlush, opts ...grpc.CallOption) (*ResponseFlush, error) {
	out := new(ResponseFlush)
	err := grpc.Invoke(ctx, "/types.ABCIApplication/Flush", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aBCIApplicationClient) Info(ctx context.Context, in *RequestInfo, opts ...grpc.CallOption) (*ResponseInfo, error) {
	out := new(ResponseInfo)
	err := grpc.Invoke(ctx, "/types.ABCIApplication/Info", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aBCIApplicationClient) SetOption(ctx context.Context, in *RequestSetOption, opts ...grpc.CallOption) (*ResponseSetOption, error) {
	out := new(ResponseSetOption)
	err := grpc.Invoke(ctx, "/types.ABCIApplication/SetOption", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aBCIApplicationClient) DeliverTx(ctx context.Context, in *RequestDeliverTx, opts ...grpc.CallOption) (*ResponseDeliverTx, error) {
	out := new(ResponseDeliverTx)
	err := grpc.Invoke(ctx, "/types.ABCIApplication/DeliverTx", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aBCIApplicationClient) CheckTx(ctx context.Context, in *RequestCheckTx, opts ...grpc.CallOption) (*ResponseCheckTx, error) {
	out := new(ResponseCheckTx)
	err := grpc.Invoke(ctx, "/types.ABCIApplication/CheckTx", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aBCIApplicationClient) Query(ctx context.Context, in *RequestQuery, opts ...grpc.CallOption) (*ResponseQuery, error) {
	out := new(ResponseQuery)
	err := grpc.Invoke(ctx, "/types.ABCIApplication/Query", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aBCIApplicationClient) Commit(ctx context.Context, in *RequestCommit, opts ...grpc.CallOption) (*ResponseCommit, error) {
	out := new(ResponseCommit)
	err := grpc.Invoke(ctx, "/types.ABCIApplication/Commit", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aBCIApplicationClient) InitChain(ctx context.Context, in *RequestInitChain, opts ...grpc.CallOption) (*ResponseInitChain, error) {
	out := new(ResponseInitChain)
	err := grpc.Invoke(ctx, "/types.ABCIApplication/InitChain", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aBCIApplicationClient) BeginBlock(ctx context.Context, in *RequestBeginBlock, opts ...grpc.CallOption) (*ResponseBeginBlock, error) {
	out := new(ResponseBeginBlock)
	err := grpc.Invoke(ctx, "/types.ABCIApplication/BeginBlock", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aBCIApplicationClient) EndBlock(ctx context.Context, in *RequestEndBlock, opts ...grpc.CallOption) (*ResponseEndBlock, error) {
	out := new(ResponseEndBlock)
	err := grpc.Invoke(ctx, "/types.ABCIApplication/EndBlock", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type ABCIApplicationServer interface {
	Echo(context.Context, *RequestEcho) (*ResponseEcho, error)
	Flush(context.Context, *RequestFlush) (*ResponseFlush, error)
	Info(context.Context, *RequestInfo) (*ResponseInfo, error)
	SetOption(context.Context, *RequestSetOption) (*ResponseSetOption, error)
	DeliverTx(context.Context, *RequestDeliverTx) (*ResponseDeliverTx, error)
	CheckTx(context.Context, *RequestCheckTx) (*ResponseCheckTx, error)
	Query(context.Context, *RequestQuery) (*ResponseQuery, error)
	Commit(context.Context, *RequestCommit) (*ResponseCommit, error)
	InitChain(context.Context, *RequestInitChain) (*ResponseInitChain, error)
	BeginBlock(context.Context, *RequestBeginBlock) (*ResponseBeginBlock, error)
	EndBlock(context.Context, *RequestEndBlock) (*ResponseEndBlock, error)
}

func RegisterABCIApplicationServer(s *grpc.Server, srv ABCIApplicationServer) {
	s.RegisterService(&_ABCIApplication_serviceDesc, srv)
}

func _ABCIApplication_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestEcho)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIApplicationServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:		srv,
		FullMethod:	"/types.ABCIApplication/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIApplicationServer).Echo(ctx, req.(*RequestEcho))
	}
	return interceptor(ctx, in, info, handler)
}

func _ABCIApplication_Flush_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestFlush)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIApplicationServer).Flush(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:		srv,
		FullMethod:	"/types.ABCIApplication/Flush",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIApplicationServer).Flush(ctx, req.(*RequestFlush))
	}
	return interceptor(ctx, in, info, handler)
}

func _ABCIApplication_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIApplicationServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:		srv,
		FullMethod:	"/types.ABCIApplication/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIApplicationServer).Info(ctx, req.(*RequestInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _ABCIApplication_SetOption_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestSetOption)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIApplicationServer).SetOption(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:		srv,
		FullMethod:	"/types.ABCIApplication/SetOption",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIApplicationServer).SetOption(ctx, req.(*RequestSetOption))
	}
	return interceptor(ctx, in, info, handler)
}

func _ABCIApplication_DeliverTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestDeliverTx)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIApplicationServer).DeliverTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:		srv,
		FullMethod:	"/types.ABCIApplication/DeliverTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIApplicationServer).DeliverTx(ctx, req.(*RequestDeliverTx))
	}
	return interceptor(ctx, in, info, handler)
}

func _ABCIApplication_CheckTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestCheckTx)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIApplicationServer).CheckTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:		srv,
		FullMethod:	"/types.ABCIApplication/CheckTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIApplicationServer).CheckTx(ctx, req.(*RequestCheckTx))
	}
	return interceptor(ctx, in, info, handler)
}

func _ABCIApplication_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIApplicationServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:		srv,
		FullMethod:	"/types.ABCIApplication/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIApplicationServer).Query(ctx, req.(*RequestQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _ABCIApplication_Commit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestCommit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIApplicationServer).Commit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:		srv,
		FullMethod:	"/types.ABCIApplication/Commit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIApplicationServer).Commit(ctx, req.(*RequestCommit))
	}
	return interceptor(ctx, in, info, handler)
}

func _ABCIApplication_InitChain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestInitChain)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIApplicationServer).InitChain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:		srv,
		FullMethod:	"/types.ABCIApplication/InitChain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIApplicationServer).InitChain(ctx, req.(*RequestInitChain))
	}
	return interceptor(ctx, in, info, handler)
}

func _ABCIApplication_BeginBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestBeginBlock)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIApplicationServer).BeginBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:		srv,
		FullMethod:	"/types.ABCIApplication/BeginBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIApplicationServer).BeginBlock(ctx, req.(*RequestBeginBlock))
	}
	return interceptor(ctx, in, info, handler)
}

func _ABCIApplication_EndBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestEndBlock)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIApplicationServer).EndBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:		srv,
		FullMethod:	"/types.ABCIApplication/EndBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIApplicationServer).EndBlock(ctx, req.(*RequestEndBlock))
	}
	return interceptor(ctx, in, info, handler)
}

var _ABCIApplication_serviceDesc = grpc.ServiceDesc{
	ServiceName:	"types.ABCIApplication",
	HandlerType:	(*ABCIApplicationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName:	"Echo",
			Handler:	_ABCIApplication_Echo_Handler,
		},
		{
			MethodName:	"Flush",
			Handler:	_ABCIApplication_Flush_Handler,
		},
		{
			MethodName:	"Info",
			Handler:	_ABCIApplication_Info_Handler,
		},
		{
			MethodName:	"SetOption",
			Handler:	_ABCIApplication_SetOption_Handler,
		},
		{
			MethodName:	"DeliverTx",
			Handler:	_ABCIApplication_DeliverTx_Handler,
		},
		{
			MethodName:	"CheckTx",
			Handler:	_ABCIApplication_CheckTx_Handler,
		},
		{
			MethodName:	"Query",
			Handler:	_ABCIApplication_Query_Handler,
		},
		{
			MethodName:	"Commit",
			Handler:	_ABCIApplication_Commit_Handler,
		},
		{
			MethodName:	"InitChain",
			Handler:	_ABCIApplication_InitChain_Handler,
		},
		{
			MethodName:	"BeginBlock",
			Handler:	_ABCIApplication_BeginBlock_Handler,
		},
		{
			MethodName:	"EndBlock",
			Handler:	_ABCIApplication_EndBlock_Handler,
		},
	},
	Streams:	[]grpc.StreamDesc{},
	Metadata:	"types/types.proto",
}

func init()	{ proto.RegisterFile("types/types.proto", fileDescriptorTypes) }

var fileDescriptorTypes = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe4, 0x58, 0xcd, 0x6e, 0x1b, 0xc9,
	0x11, 0x16, 0xff, 0x39, 0x45, 0x89, 0x94, 0x5a, 0xb2, 0x4d, 0x73, 0x11, 0x58, 0x18, 0x04, 0x5e,
	0x3a, 0xf6, 0x8a, 0x89, 0x36, 0x36, 0x6c, 0x6f, 0xb0, 0x88, 0x29, 0x39, 0x26, 0xb1, 0x49, 0xd6,
	0x19, 0x7b, 0x1d, 0x20, 0x17, 0xa2, 0xc9, 0x69, 0x91, 0x03, 0x73, 0x7e, 0x76, 0xba, 0xa9, 0xa5,
	0x7c, 0xcb, 0x7d, 0xef, 0x39, 0xe7, 0x94, 0x27, 0xc8, 0x2b, 0x04, 0x09, 0xf2, 0x0e, 0x3a, 0xec,
	0x31, 0x2f, 0x91, 0xa0, 0xba, 0x7b, 0x7e, 0x35, 0xb3, 0x58, 0xe4, 0xba, 0x17, 0xb2, 0xab, 0xeb,
	0xab, 0xee, 0xae, 0xee, 0xea, 0xaf, 0x6a, 0x1a, 0x0e, 0xc4, 0x55, 0xc0, 0xf8, 0x48, 0xfe, 0x9e,
	0x04, 0xa1, 0x2f, 0x7c, 0xd2, 0x90, 0xc2, 0xe0, 0x93, 0xa5, 0x23, 0x56, 0x9b, 0xf9, 0xc9, 0xc2,
	0x77, 0x47, 0x4b, 0x7f, 0xe9, 0x8f, 0xa4, 0x76, 0xbe, 0xb9, 0x90, 0x92, 0x14, 0x64, 0x4b, 0x59,
	0x0d, 0x46, 0x29, 0xb8, 0x60, 0x9e, 0xcd, 0x42, 0xd7, 0xf1, 0xc4, 0x48, 0xb8, 0x6b, 0x67, 0xce,
	0x47, 0x0b, 0xdf, 0x75, 0x7d, 0x2f, 0x3d, 0x8d, 0xf9, 0x8f, 0x3a, 0xb4, 0x2c, 0xf6, 0xf5, 0x86,
	0x71, 0x41, 0x86, 0x50, 0x67, 0x8b, 0x95, 0xdf, 0xaf, 0x1e, 0x57, 0x86, 0x9d, 0x53, 0x72, 0xa2,
	0x70, 0x5a, 0xfb, 0x72, 0xb1, 0xf2, 0x27, 0x3b, 0x96, 0x44, 0x90, 0x87, 0xd0, 0xb8, 0x58, 0x6f,
	0xf8, 0xaa, 0x5f, 0x93, 0xd0, 0xc3, 0x2c, 0xf4, 0x37, 0xa8, 0x9a, 0xec, 0x58, 0x0a, 0x83, 0xc3,
	0x3a, 0xde, 0x85, 0xdf, 0xaf, 0x17, 0x0d, 0x3b, 0xf5, 0x2e, 0xe4, 0xb0, 0x88, 0x20, 0x4f, 0x01,
	0x38, 0x13, 0x33, 0x3f, 0x10, 0x8e, 0xef, 0xf5, 0x1b, 0x12, 0x7f, 0x27, 0x8b, 0x7f, 0xc3, 0xc4,
	0x97, 0x52, 0x3d, 0xd9, 0xb1, 0x0c, 0x1e, 0x09, 0x68, 0xe9, 0x78, 0x8e, 0x98, 0x2d, 0x56, 0xd4,
	0xf1, 0xfa, 0xcd, 0x22, 0xcb, 0xa9, 0xe7, 0x88, 0x33, 0x54, 0xa3, 0xa5, 0x13, 0x09, 0xe8, 0xca,
	0xd7, 0x1b, 0x16, 0x5e, 0xf5, 0x5b, 0x45, 0xae, 0xfc, 0x01, 0x55, 0xe8, 0x8a, 0xc4, 0x90, 0xcf,
	0xa0, 0x33, 0x67, 0x4b, 0xc7, 0x9b, 0xcd, 0xd7, 0xfe, 0xe2, 0x7d, 0xbf, 0x2d, 0x4d, 0xfa, 0x59,
	0x93, 0x31, 0x02, 0xc6, 0xa8, 0x9f, 0xec, 0x58, 0x30, 0x8f, 0x25, 0x72, 0x0a, 0xed, 0xc5, 0x8a,
	0x2d, 0xde, 0xcf, 0xc4, 0xb6, 0x6f, 0x48, 0xcb, 0x5b, 0x59, 0xcb, 0x33, 0xd4, 0xbe, 0xdd, 0x4e,
	0x76, 0xac, 0xd6, 0x42, 0x35, 0xd1, 0x2f, 0x9b, 0xad, 0x9d, 0x4b, 0x16, 0xa2, 0xd5, 0x61, 0x91,
	0x5f, 0xe7, 0x4a, 0x2f, 0xed, 0x0c, 0x3b, 0x12, 0xc8, 0x63, 0x30, 0x98, 0x67, 0xeb, 0x85, 0x76,
	0xa4, 0xe1, 0xed, 0xdc, 0x89, 0x7a, 0x76, 0xb4, 0xcc, 0x36, 0xd3, 0x6d, 0x72, 0x02, 0x4d, 0x8c,
	0x12, 0x47, 0xf4, 0x77, 0xa5, 0xcd, 0x51, 0x6e, 0x89, 0x52, 0x37, 0xd9, 0xb1, 0x34, 0x6a, 0xdc,
	0x82, 0xc6, 0x25, 0x5d, 0x6f, 0x98, 0xf9, 0x31, 0x74, 0x52, 0x91, 0x42, 0xfa, 0xd0, 0x72, 0x19,
	0xe7, 0x74, 0xc9, 0xfa, 0x95, 0xe3, 0xca, 0xd0, 0xb0, 0x22, 0xd1, 0xec, 0xc2, 0x6e, 0x3a, 0x4e,
	0x52, 0x86, 0x18, 0x0b, 0x68, 0x78, 0xc9, 0x42, 0x8e, 0x01, 0xa0, 0x0d, 0xb5, 0x68, 0x3e, 0x87,
	0xfd, 0x7c, 0x10, 0x90, 0x7d, 0xa8, 0xbd, 0x67, 0x57, 0x1a, 0x89, 0x4d, 0x72, 0xa4, 0x17, 0x24,
	0xa3, 0xd8, 0xb0, 0xf4, 0xea, 0xc2, 0xd8, 0x36, 0x0e, 0x03, 0xf2, 0x04, 0xe0, 0x92, 0xae, 0x1d,
	0x9b, 0x0a, 0x3f, 0xe4, 0xfd, 0xca, 0x71, 0x6d, 0xd8, 0x39, 0xdd, 0xd7, 0xee, 0xbe, 0x8b, 0x14,
	0xe3, 0xfa, 0x3f, 0xaf, 0xef, 0xed, 0x58, 0x29, 0x24, 0xb9, 0x0f, 0x3d, 0x1a, 0x04, 0x33, 0x2e,
	0xa8, 0x60, 0xb3, 0xf9, 0x95, 0x60, 0x5c, 0xce, 0xb5, 0x6b, 0xed, 0xd1, 0x20, 0x78, 0x83, 0xbd,
	0x63, 0xec, 0x34, 0xed, 0xd8, 0x51, 0x19, 0x45, 0x84, 0x40, 0xdd, 0xa6, 0x82, 0xca, 0xc5, 0xee,
	0x5a, 0xb2, 0x8d, 0x7d, 0x01, 0x15, 0x2b, 0xbd, 0x58, 0xd9, 0x26, 0xb7, 0xa1, 0xb9, 0x62, 0xce,
	0x72, 0x25, 0xe4, 0xed, 0xaa, 0x59, 0x5a, 0x42, 0xcf, 0x82, 0xd0, 0xbf, 0x64, 0xf2, 0x22, 0xb5,
	0x2d, 0x25, 0x98, 0xff, 0xae, 0xc0, 0xc1, 0x8d, 0xc8, 0xc3, 0x71, 0x57, 0x94, 0xaf, 0xa2, 0xb9,
	0xb0, 0x4d, 0x1e, 0xe2, 0xb8, 0xd4, 0x66, 0xa1, 0xbe, 0xe0, 0x7b, 0xda, 0xd7, 0x89, 0xec, 0xd4,
	0x8e, 0x6a, 0x08, 0x79, 0x08, 0x07, 0x74, 0xce, 0x99, 0x27, 0x66, 0xa9, 0x3d, 0xaa, 0x1d, 0xd7,
	0x86, 0x0d, 0x6b, 0x5f, 0x29, 0xde, 0x25, 0x3b, 0x32, 0x81, 0xa3, 0xf9, 0xd5, 0x07, 0xea, 0x09,
	0xc7, 0x63, 0x69, 0x7c, 0x5d, 0xee, 0x69, 0x4f, 0xcf, 0xf3, 0xf2, 0xd2, 0xb1, 0x99, 0xb7, 0x60,
	0x7a, 0xa6, 0xc3, 0xd8, 0x24, 0x19, 0xc9, 0x3c, 0x86, 0x6e, 0xf6, 0x32, 0x90, 0x2e, 0x54, 0xc5,
	0x56, 0xfb, 0x51, 0x15, 0x5b, 0xd3, 0x8c, 0x4f, 0x32, 0x0e, 0xfc, 0x1b, 0x98, 0x07, 0xd0, 0xcb,
	0xc5, 0x78, 0x6a, 0x53, 0x2b, 0xe9, 0x4d, 0x35, 0x7b, 0xb0, 0x97, 0x09, 0x6d, 0xf3, 0xdb, 0x06,
	0xb4, 0x2d, 0xc6, 0x03, 0xdf, 0xe3, 0x8c, 0x3c, 0x05, 0x83, 0x6d, 0x17, 0x4c, 0xf1, 0x51, 0x25,
	0x77, 0xdb, 0x15, 0xe6, 0x65, 0xa4, 0xc7, 0xeb, 0x17, 0x83, 0xc9, 0x83, 0x0c, 0x97, 0x1e, 0xe6,
	0x8d, 0xd2, 0x64, 0xfa, 0x28, 0x4b, 0xa6, 0x47, 0x39, 0x6c, 0x8e, 0x4d, 0x1f, 0x64, 0xd8, 0x34,
	0x3f, 0x70, 0x86, 0x4e, 0x9f, 0x15, 0xd0, 0x69, 0x7e, 0xf9, 0x25, 0x7c, 0xfa, 0xac, 0x80, 0x4f,
	0xfb, 0x37, 0xe6, 0x2a, 0x24, 0xd4, 0x47, 0x59, 0x42, 0xcd, 0xbb, 0x93, 0x63, 0xd4, 0x5f, 0x15,
	0x31, 0xea, 0xdd, 0x9c, 0x4d, 0x29, 0xa5, 0x7e, 0x7a, 0x83, 0x52, 0x6f, 0xe7, 0x4c, 0x0b, 0x38,
	0xf5, 0x59, 0x86, 0x53, 0xa1, 0xd0, 0xb7, 0x12, 0x52, 0x7d, 0x72, 0x93, 0x54, 0xef, 0xe4, 0x8f,
	0xb6, 0x88, 0x55, 0x47, 0x39, 0x56, 0xbd, 0x95, 0x5f, 0x65, 0x29, 0xad, 0x3e, 0xc0, 0xdb, 0x9d,
	0x8b, 0x34, 0x64, 0x02, 0x16, 0x86, 0x7e, 0xa8, 0x79, 0x4f, 0x09, 0xe6, 0x10, 0xf9, 0x26, 0x89,
	0xaf, 0xef, 0xa1, 0x60, 0x19, 0xf4, 0xa9, 0xe8, 0x32, 0xff, 0x52, 0x49, 0x6c, 0x25, 0x0b, 0xa7,
	0xb9, 0xca, 0xd0, 0x5c, 0x95, 0x62, 0xe6, 0x6a, 0x86, 0x99, 0xc9, 0xcf, 0xe0, 0x60, 0x4d, 0xb9,
	0x50, 0xfb, 0x32, 0xcb, 0x90, 0x57, 0x0f, 0x15, 0x6a, 0x43, 0x14, 0x8b, 0x7d, 0x02, 0x87, 0x29,
	0x2c, 0x12, 0xa9, 0x24, 0xaa, 0xba, 0xbc, 0xbc, 0xfb, 0x31, 0xfa, 0x45, 0x10, 0x4c, 0x28, 0x5f,
	0x99, 0xbf, 0x4b, 0xfc, 0x4f, 0x58, 0x9f, 0x40, 0x7d, 0xe1, 0xdb, 0xca, 0xad, 0x3d, 0x4b, 0xb6,
	0x31, 0x13, 0xac, 0xfd, 0xa5, 0x9c, 0xd5, 0xb0, 0xb0, 0x89, 0xa8, 0xf8, 0xa6, 0x18, 0xea, 0x4a,
	0x98, 0x87, 0xc9, 0x70, 0x71, 0xf8, 0x9a, 0x7f, 0xaf, 0x24, 0xfb, 0x11, 0x53, 0xf5, 0xff, 0x37,
	0x01, 0x1e, 0x8d, 0xe3, 0xd9, 0x6c, 0x2b, 0xaf, 0x5b, 0xcd, 0x52, 0x42, 0x94, 0xa6, 0x9a, 0xd2,
	0xc9, 0x6c, 0x9a, 0x6a, 0xc9, 0x3e, 0x25, 0x68, 0x8a, 0xf7, 0x2f, 0xe4, 0x3d, 0xd8, 0xb5, 0x94,
	0x90, 0xe2, 0x2e, 0x23, 0xc3, 0x5d, 0x47, 0x40, 0x6e, 0xde, 0x10, 0xf3, 0xbf, 0x15, 0x64, 0xbf,
	0x4c, 0xf4, 0x17, 0xfa, 0x13, 0x1d, 0x71, 0x35, 0x95, 0x8e, 0x7e, 0x98, 0x8f, 0x3f, 0x01, 0x58,
	0x52, 0x3e, 0xfb, 0x86, 0x7a, 0x82, 0xd9, 0xda, 0x51, 0x63, 0x49, 0xf9, 0x1f, 0x65, 0x07, 0xb9,
	0x0b, 0x6d, 0x54, 0x6f, 0x38, 0xb3, 0xa5, 0xc7, 0x35, 0xab, 0xb5, 0xa4, 0xfc, 0x2b, 0xce, 0x6c,
	0xf2, 0x1c, 0xea, 0x82, 0x2e, 0x79, 0xbf, 0x25, 0x13, 0x43, 0xf7, 0x44, 0x15, 0xa4, 0x27, 0x5f,
	0xbc, 0x7b, 0x4d, 0x9d, 0x70, 0x7c, 0x1b, 0xf3, 0xc2, 0x7f, 0xae, 0xef, 0x75, 0x11, 0xf3, 0xc8,
	0x77, 0x1d, 0xc1, 0xdc, 0x40, 0x5c, 0x59, 0xd2, 0x86, 0x0c, 0xa1, 0x76, 0xc1, 0x98, 0x66, 0x88,
	0xfd, 0xd8, 0x74, 0xfa, 0xe4, 0x97, 0xd2, 0x58, 0x25, 0x15, 0x84, 0x98, 0x7f, 0xae, 0x26, 0xa7,
	0x9c, 0x24, 0x89, 0x1f, 0xd7, 0x1e, 0xfc, 0xad, 0x82, 0x79, 0x32, 0x4b, 0x49, 0xe4, 0x0c, 0x0e,
	0xe2, 0xec, 0x3c, 0xdb, 0x04, 0x36, 0xc5, 0xda, 0xe5, 0xfb, 0x0b, 0x9f, 0xfd, 0xd8, 0xe0, 0x2b,
	0x85, 0x27, 0xbf, 0x87, 0x3b, 0x0b, 0x1c, 0xd5, 0xe3, 0x1b, 0x3e, 0x0b, 0x68, 0x48, 0xdd, 0x78,
	0xa8, 0x6a, 0x86, 0x82, 0xcf, 0x22, 0xd4, 0x6b, 0x04, 0x71, 0xeb, 0xd6, 0x22, 0xd3, 0xa1, 0xc7,
	0x33, 0x7f, 0x8a, 0x29, 0x3f, 0x4d, 0x83, 0x45, 0xa7, 0x62, 0xfe, 0xb5, 0x02, 0xbd, 0xdc, 0x80,
	0x64, 0x04, 0xa0, 0x58, 0x84, 0x3b, 0x1f, 0x98, 0x4e, 0xcf, 0x91, 0x1f, 0xd2, 0xe1, 0x37, 0xce,
	0x07, 0x66, 0x19, 0xf3, 0xa8, 0x49, 0xee, 0x43, 0x4b, 0x6c, 0x15, 0x3a, 0x5b, 0x02, 0xbd, 0xdd,
	0x4a, 0x68, 0x53, 0xc8, 0x7f, 0xf2, 0x18, 0x76, 0xd5, 0xc0, 0x4b, 0x9f, 0x73, 0x27, 0xd0, 0x89,
	0x99, 0xa4, 0x87, 0x7e, 0x25, 0x35, 0x56, 0x67, 0x9e, 0x08, 0xe6, 0x9f, 0xc0, 0x88, 0xa7, 0x25,
	0x1f, 0x81, 0xe1, 0xd2, 0xad, 0xae, 0x0f, 0x71, 0x6d, 0x0d, 0xab, 0xed, 0xd2, 0xad, 0x2c, 0x0d,
	0xc9, 0x1d, 0x68, 0xa1, 0x52, 0x6c, 0xd5, 0x9e, 0x35, 0xac, 0xa6, 0x4b, 0xb7, 0x6f, 0xb7, 0xb1,
	0x62, 0x49, 0x79, 0x54, 0xfc, 0xb9, 0x74, 0xfb, 0x8a, 0x72, 0xf3, 0x73, 0x68, 0xaa, 0x45, 0xfe,
	0xa0, 0x81, 0xd1, 0xbe, 0x9a, 0xb1, 0xff, 0x35, 0x74, 0x52, 0xeb, 0x26, 0xbf, 0x80, 0x5b, 0xca,
	0xc3, 0x80, 0x86, 0x42, 0xee, 0x48, 0x66, 0x40, 0x22, 0x95, 0xaf, 0x69, 0x28, 0x70, 0x4a, 0x55,
	0xce, 0xfe, 0xab, 0x0a, 0x4d, 0x55, 0x2a, 0x92, 0xfb, 0x98, 0x76, 0xa9, 0xe3, 0xcd, 0x1c, 0x5b,
	0x65, 0x88, 0x71, 0xe7, 0xbb, 0xeb, 0x7b, 0x2d, 0xc9, 0xa6, 0xd3, 0x73, 0xcc, 0xb4, 0xd8, 0xb0,
	0x53, 0xc4, 0x55, 0xcd, 0x54, 0xb2, 0x04, 0xea, 0xc2, 0x71, 0x99, 0x76, 0x51, 0xb6, 0x71, 0xe5,
	0xde, 0xc6, 0x95, 0x5b, 0x52, 0x57, 0x5b, 0xe2, 0x6d, 0x5c, 0xdc, 0x92, 0x57, 0xb0, 0x97, 0x4a,
	0x18, 0x8e, 0xad, 0x0b, 0x99, 0x6e, 0xfa, 0x34, 0xa6, 0xe7, 0xe3, 0x43, 0x0c, 0xd7, 0xef, 0xae,
	0xef, 0x75, 0x7e, 0x1b, 0xa5, 0x90, 0xe9, 0xb9, 0xd5, 0x89, 0xf3, 0xc9, 0xd4, 0x26, 0x43, 0x90,
	0xe9, 0x65, 0xa6, 0x52, 0xac, 0x4a, 0x3b, 0x8a, 0x91, 0xbb, 0xd8, 0xaf, 0x73, 0x30, 0x56, 0xca,
	0x1f, 0x81, 0x81, 0x41, 0xa7, 0x20, 0x8a, 0xa0, 0xdb, 0xd8, 0x21, 0x95, 0x1f, 0x43, 0x2f, 0x29,
	0x71, 0x15, 0x44, 0xb1, 0x75, 0x37, 0xe9, 0x96, 0xc0, 0xbb, 0xd0, 0x8e, 0xd3, 0x9b, 0x21, 0x11,
	0x2d, 0xaa, 0xb3, 0xda, 0x97, 0xd0, 0xd2, 0x4b, 0x2c, 0xac, 0xd4, 0x7f, 0x0e, 0x0d, 0x3c, 0x97,
	0xe8, 0x42, 0x45, 0x25, 0x94, 0x3c, 0x0f, 0x26, 0x32, 0xf5, 0xba, 0x02, 0x9a, 0xcf, 0x60, 0x2f,
	0xa3, 0xc5, 0x4c, 0x22, 0x7c, 0x41, 0xd7, 0xfa, 0x40, 0x95, 0x10, 0x4f, 0x56, 0x4d, 0x26, 0x33,
	0x9f, 0x83, 0x11, 0x5f, 0x7a, 0x3c, 0x85, 0x60, 0x33, 0x9f, 0x45, 0xdf, 0x54, 0xbb, 0x56, 0x33,
	0xd8, 0xcc, 0xbf, 0x50, 0xf9, 0x2a, 0xf0, 0xbf, 0xd1, 0xdf, 0x0e, 0x35, 0x4b, 0x09, 0xe6, 0x67,
	0xd0, 0x8e, 0xaa, 0xfa, 0x72, 0xd3, 0x92, 0x28, 0x38, 0xfd, 0xb6, 0x01, 0xbd, 0x17, 0xe3, 0xb3,
	0xe9, 0x8b, 0x20, 0x58, 0x3b, 0x0b, 0x2a, 0x33, 0xfb, 0x08, 0xea, 0xb2, 0x76, 0x29, 0x78, 0x7c,
	0x18, 0x14, 0x15, 0xd1, 0xe4, 0x14, 0x1a, 0xb2, 0x84, 0x21, 0x45, 0x6f, 0x10, 0x83, 0xc2, 0x5a,
	0x1a, 0x27, 0x51, 0x45, 0xce, 0xcd, 0xa7, 0x88, 0x41, 0x51, 0x41, 0x4d, 0x3e, 0x07, 0x23, 0x29,
	0x3e, 0xca, 0x1e, 0x24, 0x06, 0xa5, 0xa5, 0x35, 0xda, 0x27, 0x79, 0xa8, 0xec, 0xf3, 0x7d, 0x50,
	0x5a, 0x83, 0x92, 0xa7, 0xd0, 0x8a, 0x32, 0x79, 0xf1, 0x93, 0xc1, 0xa0, 0xa4, 0xec, 0xc5, 0xed,
	0x51, 0x15, 0x4d, 0xd1, 0xbb, 0xc6, 0xa0, 0xb0, 0x36, 0x27, 0x8f, 0xa1, 0xa9, 0x89, 0xb8, 0xf0,
	0xe3, 0x7f, 0x50, 0x5c, 0xbc, 0xa2, 0x93, 0xc9, 0xb7, 0x75, 0xd9, 0xdb, 0xcb, 0xa0, 0xf4, 0x23,
	0x82, 0xbc, 0x00, 0x48, 0x7d, 0xc0, 0x96, 0x3e, 0xaa, 0x0c, 0xca, 0x3f, 0x0e, 0x08, 0x86, 0x63,
	0xfc, 0xc1, 0x57, 0xfc, 0xd8, 0x31, 0x28, 0xab, 0xd7, 0xe7, 0x4d, 0xf9, 0x20, 0xf6, 0xe9, 0xff,
	0x02, 0x00, 0x00, 0xff, 0xff, 0x46, 0xbe, 0x48, 0x9c, 0x8c, 0x13, 0x00, 0x00,
}
