package rpcserver

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/tendermint/go-amino"
	types "github.com/tendermint/tendermint/rpc/lib/types"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/log"
)

func RegisterRPCFuncs(mux *http.ServeMux, funcMap map[string]*RPCFunc, cdc *amino.Codec, logger log.Logger) {

	for funcName, rpcFunc := range funcMap {
		mux.HandleFunc("/"+funcName, makeHTTPHandler(rpcFunc, cdc, logger))
	}

	mux.HandleFunc("/", makeJSONRPCHandler(funcMap, cdc, logger))
}

type RPCFunc struct {
	f		reflect.Value
	args		[]reflect.Type
	returns		[]reflect.Type
	argNames	[]string
	ws		bool
}

func NewRPCFunc(f interface{}, args string) *RPCFunc {
	return newRPCFunc(f, args, false)
}

func NewWSRPCFunc(f interface{}, args string) *RPCFunc {
	return newRPCFunc(f, args, true)
}

func newRPCFunc(f interface{}, args string, ws bool) *RPCFunc {
	var argNames []string
	if args != "" {
		argNames = strings.Split(args, ",")
	}
	return &RPCFunc{
		f:		reflect.ValueOf(f),
		args:		funcArgTypes(f),
		returns:	funcReturnTypes(f),
		argNames:	argNames,
		ws:		ws,
	}
}

func funcArgTypes(f interface{}) []reflect.Type {
	t := reflect.TypeOf(f)
	n := t.NumIn()
	typez := make([]reflect.Type, n)
	for i := 0; i < n; i++ {
		typez[i] = t.In(i)
	}
	return typez
}

func funcReturnTypes(f interface{}) []reflect.Type {
	t := reflect.TypeOf(f)
	n := t.NumOut()
	typez := make([]reflect.Type, n)
	for i := 0; i < n; i++ {
		typez[i] = t.Out(i)
	}
	return typez
}

func makeJSONRPCHandler(funcMap map[string]*RPCFunc, cdc *amino.Codec, logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteRPCResponseHTTP(w, types.RPCInvalidRequestError("", errors.Wrap(err, "Error reading request body")))
			return
		}

		if len(b) == 0 {
			writeListOfEndpoints(w, r, funcMap)
			return
		}

		var request types.RPCRequest
		err = json.Unmarshal(b, &request)
		if err != nil {
			WriteRPCResponseHTTP(w, types.RPCParseError("", errors.Wrap(err, "Error unmarshalling request")))
			return
		}

		if request.ID == "" {
			logger.Debug("HTTPJSONRPC received a notification, skipping... (please send a non-empty ID if you want to call a method)")
			return
		}
		if len(r.URL.Path) > 1 {
			WriteRPCResponseHTTP(w, types.RPCInvalidRequestError(request.ID, errors.Errorf("Path %s is invalid", r.URL.Path)))
			return
		}
		rpcFunc := funcMap[request.Method]
		if rpcFunc == nil || rpcFunc.ws {
			WriteRPCResponseHTTP(w, types.RPCMethodNotFoundError(request.ID))
			return
		}
		var args []reflect.Value
		if len(request.Params) > 0 {
			args, err = jsonParamsToArgsRPC(rpcFunc, cdc, request.Params)
			if err != nil {
				WriteRPCResponseHTTP(w, types.RPCInvalidParamsError(request.ID, errors.Wrap(err, "Error converting json params to arguments")))
				return
			}
		}
		returns := rpcFunc.f.Call(args)
		logger.Info("HTTPJSONRPC", "method", request.Method, "args", args, "returns", returns)
		result, err := unreflectResult(returns)
		if err != nil {
			WriteRPCResponseHTTP(w, types.RPCInternalError(request.ID, err))
			return
		}
		WriteRPCResponseHTTP(w, types.NewRPCSuccessResponse(cdc, request.ID, result))
	}
}

func mapParamsToArgs(rpcFunc *RPCFunc, cdc *amino.Codec, params map[string]json.RawMessage, argsOffset int) ([]reflect.Value, error) {
	values := make([]reflect.Value, len(rpcFunc.argNames))
	for i, argName := range rpcFunc.argNames {
		argType := rpcFunc.args[i+argsOffset]

		if p, ok := params[argName]; ok && p != nil && len(p) > 0 {
			val := reflect.New(argType)
			err := cdc.UnmarshalJSON(p, val.Interface())
			if err != nil {
				return nil, err
			}
			values[i] = val.Elem()
		} else {
			values[i] = reflect.Zero(argType)
		}
	}

	return values, nil
}

func arrayParamsToArgs(rpcFunc *RPCFunc, cdc *amino.Codec, params []json.RawMessage, argsOffset int) ([]reflect.Value, error) {
	if len(rpcFunc.argNames) != len(params) {
		return nil, errors.Errorf("Expected %v parameters (%v), got %v (%v)",
			len(rpcFunc.argNames), rpcFunc.argNames, len(params), params)
	}

	values := make([]reflect.Value, len(params))
	for i, p := range params {
		argType := rpcFunc.args[i+argsOffset]
		val := reflect.New(argType)
		err := cdc.UnmarshalJSON(p, val.Interface())
		if err != nil {
			return nil, err
		}
		values[i] = val.Elem()
	}
	return values, nil
}

func jsonParamsToArgs(rpcFunc *RPCFunc, cdc *amino.Codec, raw []byte, argsOffset int) ([]reflect.Value, error) {

	var m map[string]json.RawMessage
	err := json.Unmarshal(raw, &m)
	if err == nil {
		return mapParamsToArgs(rpcFunc, cdc, m, argsOffset)
	}

	var a []json.RawMessage
	err = json.Unmarshal(raw, &a)
	if err == nil {
		return arrayParamsToArgs(rpcFunc, cdc, a, argsOffset)
	}

	return nil, errors.Errorf("Unknown type for JSON params: %v. Expected map or array", err)
}

func jsonParamsToArgsRPC(rpcFunc *RPCFunc, cdc *amino.Codec, params json.RawMessage) ([]reflect.Value, error) {
	return jsonParamsToArgs(rpcFunc, cdc, params, 0)
}

func jsonParamsToArgsWS(rpcFunc *RPCFunc, cdc *amino.Codec, params json.RawMessage, wsCtx types.WSRPCContext) ([]reflect.Value, error) {
	values, err := jsonParamsToArgs(rpcFunc, cdc, params, 1)
	if err != nil {
		return nil, err
	}
	return append([]reflect.Value{reflect.ValueOf(wsCtx)}, values...), nil
}

func makeHTTPHandler(rpcFunc *RPCFunc, cdc *amino.Codec, logger log.Logger) func(http.ResponseWriter, *http.Request) {

	if rpcFunc.ws {
		return func(w http.ResponseWriter, r *http.Request) {
			WriteRPCResponseHTTP(w, types.RPCMethodNotFoundError(""))
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("HTTP HANDLER", "req", r)
		args, err := httpParamsToArgs(rpcFunc, cdc, r)
		if err != nil {
			WriteRPCResponseHTTP(w, types.RPCInvalidParamsError("", errors.Wrap(err, "Error converting http params to arguments")))
			return
		}
		returns := rpcFunc.f.Call(args)
		logger.Trace("HTTPRestRPC", "method", r.URL.Path, "args", args, "returns", returns)
		result, err := unreflectResult(returns)
		if err != nil {
			WriteRPCResponseHTTP(w, types.RPCInternalError("", err))
			return
		}
		WriteRPCResponseHTTP(w, types.NewRPCSuccessResponse(cdc, "", result))
	}
}

func httpParamsToArgs(rpcFunc *RPCFunc, cdc *amino.Codec, r *http.Request) ([]reflect.Value, error) {
	values := make([]reflect.Value, len(rpcFunc.args))

	for i, name := range rpcFunc.argNames {
		argType := rpcFunc.args[i]

		values[i] = reflect.Zero(argType)

		arg := GetParam(r, name)

		if "" == arg {
			continue
		}

		v, err, ok := nonJSONToArg(cdc, argType, arg)
		if err != nil {
			return nil, err
		}
		if ok {
			values[i] = v
			continue
		}

		values[i], err = _jsonStringToArg(cdc, argType, arg)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func _jsonStringToArg(cdc *amino.Codec, ty reflect.Type, arg string) (reflect.Value, error) {
	v := reflect.New(ty)
	err := cdc.UnmarshalJSON([]byte(arg), v.Interface())
	if err != nil {
		return v, err
	}
	v = v.Elem()
	return v, nil
}

func nonJSONToArg(cdc *amino.Codec, ty reflect.Type, arg string) (reflect.Value, error, bool) {
	isQuotedString := strings.HasPrefix(arg, `"`) && strings.HasSuffix(arg, `"`)
	isHexString := strings.HasPrefix(strings.ToLower(arg), "0x")
	expectingString := ty.Kind() == reflect.String
	expectingByteSlice := ty.Kind() == reflect.Slice && ty.Elem().Kind() == reflect.Uint8

	if isHexString {
		if !expectingString && !expectingByteSlice {
			err := errors.Errorf("Got a hex string arg, but expected '%s'",
				ty.Kind().String())
			return reflect.ValueOf(nil), err, false
		}

		var value []byte
		value, err := hex.DecodeString(arg[2:])
		if err != nil {
			return reflect.ValueOf(nil), err, false
		}
		if ty.Kind() == reflect.String {
			return reflect.ValueOf(string(value)), nil, true
		}
		return reflect.ValueOf([]byte(value)), nil, true
	}

	if isQuotedString && expectingByteSlice {
		v := reflect.New(reflect.TypeOf(""))
		err := cdc.UnmarshalJSON([]byte(arg), v.Interface())
		if err != nil {
			return reflect.ValueOf(nil), err, false
		}
		v = v.Elem()
		return reflect.ValueOf([]byte(v.String())), nil, true
	}

	return reflect.ValueOf(nil), nil, false
}

const (
	defaultWSWriteChanCapacity	= 1000
	defaultWSWriteWait		= 10 * time.Second
	defaultWSReadWait		= 30 * time.Second
	defaultWSPingPeriod		= (defaultWSReadWait * 9) / 10
)

type wsConnection struct {
	cmn.BaseService

	remoteAddr	string
	baseConn	*websocket.Conn
	writeChan	chan types.RPCResponse

	funcMap	map[string]*RPCFunc
	cdc	*amino.Codec

	writeChanCapacity	int

	writeWait	time.Duration

	readWait	time.Duration

	pingPeriod	time.Duration

	eventSub	types.EventSubscriber
}

func NewWSConnection(baseConn *websocket.Conn, funcMap map[string]*RPCFunc, cdc *amino.Codec, options ...func(*wsConnection)) *wsConnection {
	wsc := &wsConnection{
		remoteAddr:		baseConn.RemoteAddr().String(),
		baseConn:		baseConn,
		funcMap:		funcMap,
		cdc:			cdc,
		writeWait:		defaultWSWriteWait,
		writeChanCapacity:	defaultWSWriteChanCapacity,
		readWait:		defaultWSReadWait,
		pingPeriod:		defaultWSPingPeriod,
	}
	for _, option := range options {
		option(wsc)
	}
	wsc.BaseService = *cmn.NewBaseService(nil, "wsConnection", wsc)
	return wsc
}

func EventSubscriber(eventSub types.EventSubscriber) func(*wsConnection) {
	return func(wsc *wsConnection) {
		wsc.eventSub = eventSub
	}
}

func WriteWait(writeWait time.Duration) func(*wsConnection) {
	return func(wsc *wsConnection) {
		wsc.writeWait = writeWait
	}
}

func WriteChanCapacity(cap int) func(*wsConnection) {
	return func(wsc *wsConnection) {
		wsc.writeChanCapacity = cap
	}
}

func ReadWait(readWait time.Duration) func(*wsConnection) {
	return func(wsc *wsConnection) {
		wsc.readWait = readWait
	}
}

func PingPeriod(pingPeriod time.Duration) func(*wsConnection) {
	return func(wsc *wsConnection) {
		wsc.pingPeriod = pingPeriod
	}
}

func (wsc *wsConnection) OnStart() error {
	wsc.writeChan = make(chan types.RPCResponse, wsc.writeChanCapacity)

	go wsc.readRoutine()

	wsc.writeRoutine()

	return nil
}

func (wsc *wsConnection) OnStop() {

	if wsc.eventSub != nil {
		wsc.eventSub.UnsubscribeAll(context.TODO(), wsc.remoteAddr)
	}
}

func (wsc *wsConnection) GetRemoteAddr() string {
	return wsc.remoteAddr
}

func (wsc *wsConnection) GetEventSubscriber() types.EventSubscriber {
	return wsc.eventSub
}

func (wsc *wsConnection) WriteRPCResponse(resp types.RPCResponse) {
	select {
	case <-wsc.Quit():
		return
	case wsc.writeChan <- resp:
	}
}

func (wsc *wsConnection) TryWriteRPCResponse(resp types.RPCResponse) bool {
	select {
	case <-wsc.Quit():
		return false
	case wsc.writeChan <- resp:
		return true
	default:
		return false
	}
}

func (wsc *wsConnection) Codec() *amino.Codec {
	return wsc.cdc
}

func (wsc *wsConnection) readRoutine() {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("WSJSONRPC: %v", r)
			}
			wsc.Logger.Error("Panic in WSJSONRPC handler", "err", err, "stack", string(debug.Stack()))
			wsc.WriteRPCResponse(types.RPCInternalError("unknown", err))
			go wsc.readRoutine()
		} else {
			wsc.baseConn.Close()
		}
	}()

	wsc.baseConn.SetPongHandler(func(m string) error {
		return wsc.baseConn.SetReadDeadline(time.Now().Add(wsc.readWait))
	})

	for {
		select {
		case <-wsc.Quit():
			return
		default:

			if err := wsc.baseConn.SetReadDeadline(time.Now().Add(wsc.readWait)); err != nil {
				wsc.Logger.Error("failed to set read deadline", "err", err)
			}
			var in []byte
			_, in, err := wsc.baseConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					wsc.Logger.Info("Client closed the connection")
				} else {
					wsc.Logger.Error("Failed to read request", "err", err)
				}
				wsc.Stop()
				return
			}

			var request types.RPCRequest
			err = json.Unmarshal(in, &request)
			if err != nil {
				wsc.WriteRPCResponse(types.RPCParseError("", errors.Wrap(err, "Error unmarshaling request")))
				continue
			}

			if request.ID == "" {
				wsc.Logger.Debug("WSJSONRPC received a notification, skipping... (please send a non-empty ID if you want to call a method)")
				continue
			}

			rpcFunc := wsc.funcMap[request.Method]
			if rpcFunc == nil {
				wsc.WriteRPCResponse(types.RPCMethodNotFoundError(request.ID))
				continue
			}
			var args []reflect.Value
			if rpcFunc.ws {
				wsCtx := types.WSRPCContext{Request: request, WSRPCConnection: wsc}
				if len(request.Params) > 0 {
					args, err = jsonParamsToArgsWS(rpcFunc, wsc.cdc, request.Params, wsCtx)
				}
			} else {
				if len(request.Params) > 0 {
					args, err = jsonParamsToArgsRPC(rpcFunc, wsc.cdc, request.Params)
				}
			}
			if err != nil {
				wsc.WriteRPCResponse(types.RPCInternalError(request.ID, errors.Wrap(err, "Error converting json params to arguments")))
				continue
			}
			returns := rpcFunc.f.Call(args)

			wsc.Logger.Info("WSJSONRPC", "method", request.Method)

			result, err := unreflectResult(returns)
			if err != nil {
				wsc.WriteRPCResponse(types.RPCInternalError(request.ID, err))
				continue
			} else {
				wsc.WriteRPCResponse(types.NewRPCSuccessResponse(wsc.cdc, request.ID, result))
				continue
			}

		}
	}
}

func (wsc *wsConnection) writeRoutine() {
	pingTicker := time.NewTicker(wsc.pingPeriod)
	defer func() {
		pingTicker.Stop()
		if err := wsc.baseConn.Close(); err != nil {
			wsc.Logger.Error("Error closing connection", "err", err)
		}
	}()

	pongs := make(chan string, 1)
	wsc.baseConn.SetPingHandler(func(m string) error {
		select {
		case pongs <- m:
		default:
		}
		return nil
	})

	for {
		select {
		case m := <-pongs:
			err := wsc.writeMessageWithDeadline(websocket.PongMessage, []byte(m))
			if err != nil {
				wsc.Logger.Info("Failed to write pong (client may disconnect)", "err", err)
			}
		case <-pingTicker.C:
			err := wsc.writeMessageWithDeadline(websocket.PingMessage, []byte{})
			if err != nil {
				wsc.Logger.Error("Failed to write ping", "err", err)
				wsc.Stop()
				return
			}
		case msg := <-wsc.writeChan:
			jsonBytes, err := json.MarshalIndent(msg, "", "  ")
			if err != nil {
				wsc.Logger.Error("Failed to marshal RPCResponse to JSON", "err", err)
			} else {
				if err = wsc.writeMessageWithDeadline(websocket.TextMessage, jsonBytes); err != nil {
					wsc.Logger.Error("Failed to write response", "err", err)
					wsc.Stop()
					return
				}
			}
		case <-wsc.Quit():
			return
		}
	}
}

func (wsc *wsConnection) writeMessageWithDeadline(msgType int, msg []byte) error {
	if err := wsc.baseConn.SetWriteDeadline(time.Now().Add(wsc.writeWait)); err != nil {
		return err
	}
	return wsc.baseConn.WriteMessage(msgType, msg)
}

type WebsocketManager struct {
	websocket.Upgrader
	funcMap		map[string]*RPCFunc
	cdc		*amino.Codec
	logger		log.Logger
	wsConnOptions	[]func(*wsConnection)
}

func NewWebsocketManager(funcMap map[string]*RPCFunc, cdc *amino.Codec, wsConnOptions ...func(*wsConnection)) *WebsocketManager {
	return &WebsocketManager{
		funcMap:	funcMap,
		cdc:		cdc,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {

				return true
			},
		},
		logger:		log.NewNopLogger(),
		wsConnOptions:	wsConnOptions,
	}
}

func (wm *WebsocketManager) SetLogger(l log.Logger) {
	wm.logger = l
}

func (wm *WebsocketManager) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wm.Upgrade(w, r, nil)
	if err != nil {

		wm.logger.Error("Failed to upgrade to websocket connection", "err", err)
		return
	}

	con := NewWSConnection(wsConn, wm.funcMap, wm.cdc, wm.wsConnOptions...)
	con.SetLogger(wm.logger.With("remote", wsConn.RemoteAddr()))
	wm.logger.Info("New websocket connection", "remote", con.remoteAddr)
	err = con.Start()
	if err != nil {
		wm.logger.Error("Error starting connection", "err", err)
	}
}

func unreflectResult(returns []reflect.Value) (interface{}, error) {
	errV := returns[1]
	if errV.Interface() != nil {
		return nil, errors.Errorf("%v", errV.Interface())
	}
	rv := returns[0]

	rvp := reflect.New(rv.Type())
	rvp.Elem().Set(rv)
	return rvp.Interface(), nil
}

func writeListOfEndpoints(w http.ResponseWriter, r *http.Request, funcMap map[string]*RPCFunc) {
	noArgNames := []string{}
	argNames := []string{}
	for name, funcData := range funcMap {
		if len(funcData.args) == 0 {
			noArgNames = append(noArgNames, name)
		} else {
			argNames = append(argNames, name)
		}
	}
	sort.Strings(noArgNames)
	sort.Strings(argNames)
	buf := new(bytes.Buffer)
	buf.WriteString("<html><body>")
	buf.WriteString("<br>Available endpoints:<br>")

	for _, name := range noArgNames {
		link := fmt.Sprintf("//%s/%s", r.Host, name)
		buf.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a></br>", link, link))
	}

	buf.WriteString("<br>Endpoints that require arguments:<br>")
	for _, name := range argNames {
		link := fmt.Sprintf("//%s/%s?", r.Host, name)
		funcData := funcMap[name]
		for i, argName := range funcData.argNames {
			link += argName + "=_"
			if i < len(funcData.argNames)-1 {
				link += "&"
			}
		}
		buf.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a></br>", link, link))
	}
	buf.WriteString("</body></html>")
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write(buf.Bytes())
}
