package rpcclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	metrics "github.com/rcrowley/go-metrics"

	"github.com/tendermint/go-amino"
	types "bcbchain.io/rpc/lib/types"
	cmn "github.com/tendermint/tmlibs/common"
)

const (
	defaultMaxReconnectAttempts	= 25
	defaultWriteWait		= 0
	defaultReadWait			= 0
	defaultPingPeriod		= 0
)

type WSClient struct {
	cmn.BaseService

	conn	*websocket.Conn
	cdc	*amino.Codec

	Address		string
	Endpoint	string
	Dialer		func(string, string) (net.Conn, error)

	PingPongLatencyTimer	metrics.Timer

	ResponsesCh	chan types.RPCResponse

	onReconnect	func()

	send		chan types.RPCRequest
	backlog		chan types.RPCRequest
	reconnectAfter	chan error
	readRoutineQuit	chan struct{}

	wg	sync.WaitGroup

	mtx		sync.RWMutex
	sentLastPingAt	time.Time
	reconnecting	bool

	maxReconnectAttempts	int

	writeWait	time.Duration

	readWait	time.Duration

	pingPeriod	time.Duration
}

func NewWSClient(remoteAddr, endpoint string, options ...func(*WSClient)) *WSClient {
	addr, dialer := makeHTTPDialer(remoteAddr)
	c := &WSClient{
		cdc:			amino.NewCodec(),
		Address:		addr,
		Dialer:			dialer,
		Endpoint:		endpoint,
		PingPongLatencyTimer:	metrics.NewTimer(),

		maxReconnectAttempts:	defaultMaxReconnectAttempts,
		readWait:		defaultReadWait,
		writeWait:		defaultWriteWait,
		pingPeriod:		defaultPingPeriod,
	}
	c.BaseService = *cmn.NewBaseService(nil, "WSClient", c)
	for _, option := range options {
		option(c)
	}
	return c
}

func MaxReconnectAttempts(max int) func(*WSClient) {
	return func(c *WSClient) {
		c.maxReconnectAttempts = max
	}
}

func ReadWait(readWait time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.readWait = readWait
	}
}

func WriteWait(writeWait time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.writeWait = writeWait
	}
}

func PingPeriod(pingPeriod time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.pingPeriod = pingPeriod
	}
}

func OnReconnect(cb func()) func(*WSClient) {
	return func(c *WSClient) {
		c.onReconnect = cb
	}
}

func (c *WSClient) String() string {
	return fmt.Sprintf("%s (%s)", c.Address, c.Endpoint)
}

func (c *WSClient) OnStart() error {
	err := c.dial()
	if err != nil {
		return err
	}

	c.ResponsesCh = make(chan types.RPCResponse)

	c.send = make(chan types.RPCRequest)

	c.reconnectAfter = make(chan error, 1)

	c.backlog = make(chan types.RPCRequest, 1)

	c.startReadWriteRoutines()
	go c.reconnectRoutine()

	return nil
}

func (c *WSClient) OnStop()	{}

func (c *WSClient) Stop() error {
	if err := c.BaseService.Stop(); err != nil {
		return err
	}

	c.wg.Wait()
	close(c.ResponsesCh)

	return nil
}

func (c *WSClient) IsReconnecting() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.reconnecting
}

func (c *WSClient) IsActive() bool {
	return c.IsRunning() && !c.IsReconnecting()
}

func (c *WSClient) Send(ctx context.Context, request types.RPCRequest) error {
	select {
	case c.send <- request:
		c.Logger.Info("sent a request", "req", request)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *WSClient) Call(ctx context.Context, method string, params map[string]interface{}) error {
	request, err := types.MapToRequest(c.cdc, "ws-client", method, params)
	if err != nil {
		return err
	}
	return c.Send(ctx, request)
}

func (c *WSClient) CallWithArrayParams(ctx context.Context, method string, params []interface{}) error {
	request, err := types.ArrayToRequest(c.cdc, "ws-client", method, params)
	if err != nil {
		return err
	}
	return c.Send(ctx, request)
}

func (c *WSClient) Codec() *amino.Codec {
	return c.cdc
}

func (c *WSClient) SetCodec(cdc *amino.Codec) {
	c.cdc = cdc
}

func (c *WSClient) dial() error {
	dialer := &websocket.Dialer{
		NetDial:	c.Dialer,
		Proxy:		http.ProxyFromEnvironment,
	}
	rHeader := http.Header{}
	conn, _, err := dialer.Dial("ws://"+c.Address+c.Endpoint, rHeader)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *WSClient) reconnect() error {
	attempt := 0

	c.mtx.Lock()
	c.reconnecting = true
	c.mtx.Unlock()
	defer func() {
		c.mtx.Lock()
		c.reconnecting = false
		c.mtx.Unlock()
	}()

	for {
		jitterSeconds := time.Duration(cmn.RandFloat64() * float64(time.Second))
		backoffDuration := jitterSeconds + ((1 << uint(attempt)) * time.Second)

		c.Logger.Info("reconnecting", "attempt", attempt+1, "backoff_duration", backoffDuration)
		time.Sleep(backoffDuration)

		err := c.dial()
		if err != nil {
			c.Logger.Error("failed to redial", "err", err)
		} else {
			c.Logger.Info("reconnected")
			if c.onReconnect != nil {
				go c.onReconnect()
			}
			return nil
		}

		attempt++

		if attempt > c.maxReconnectAttempts {
			return errors.Wrap(err, "reached maximum reconnect attempts")
		}
	}
}

func (c *WSClient) startReadWriteRoutines() {
	c.wg.Add(2)
	c.readRoutineQuit = make(chan struct{})
	go c.readRoutine()
	go c.writeRoutine()
}

func (c *WSClient) processBacklog() error {
	select {
	case request := <-c.backlog:
		if c.writeWait > 0 {
			if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
				c.Logger.Error("failed to set write deadline", "err", err)
			}
		}
		if err := c.conn.WriteJSON(request); err != nil {
			c.Logger.Error("failed to resend request", "err", err)
			c.reconnectAfter <- err

			c.backlog <- request
			return err
		}
		c.Logger.Info("resend a request", "req", request)
	default:
	}
	return nil
}

func (c *WSClient) reconnectRoutine() {
	for {
		select {
		case originalError := <-c.reconnectAfter:

			c.wg.Wait()
			if err := c.reconnect(); err != nil {
				c.Logger.Error("failed to reconnect", "err", err, "original_err", originalError)
				c.Stop()
				return
			}

		LOOP:
			for {
				select {
				case <-c.reconnectAfter:
				default:
					break LOOP
				}
			}
			err := c.processBacklog()
			if err == nil {
				c.startReadWriteRoutines()
			}

		case <-c.Quit():
			return
		}
	}
}

func (c *WSClient) writeRoutine() {
	var ticker *time.Ticker
	if c.pingPeriod > 0 {

		ticker = time.NewTicker(c.pingPeriod)
	} else {

		ticker = &time.Ticker{C: make(<-chan time.Time)}
	}

	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {

		}
		c.wg.Done()
	}()

	for {
		select {
		case request := <-c.send:
			if c.writeWait > 0 {
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.conn.WriteJSON(request); err != nil {
				c.Logger.Error("failed to send request", "err", err)
				c.reconnectAfter <- err

				c.backlog <- request
				return
			}
		case <-ticker.C:
			if c.writeWait > 0 {
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.Logger.Error("failed to write ping", "err", err)
				c.reconnectAfter <- err
				return
			}
			c.mtx.Lock()
			c.sentLastPingAt = time.Now()
			c.mtx.Unlock()
			c.Logger.Debug("sent ping")
		case <-c.readRoutineQuit:
			return
		case <-c.Quit():
			if err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				c.Logger.Error("failed to write message", "err", err)
			}
			return
		}
	}
}

func (c *WSClient) readRoutine() {
	defer func() {
		if err := c.conn.Close(); err != nil {

		}
		c.wg.Done()
	}()

	c.conn.SetPongHandler(func(string) error {

		c.mtx.RLock()
		t := c.sentLastPingAt
		c.mtx.RUnlock()
		c.PingPongLatencyTimer.UpdateSince(t)

		c.Logger.Debug("got pong")
		return nil
	})

	for {

		if c.readWait > 0 {
			if err := c.conn.SetReadDeadline(time.Now().Add(c.readWait)); err != nil {
				c.Logger.Error("failed to set read deadline", "err", err)
			}
		}
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				return
			}

			c.Logger.Error("failed to read response", "err", err)
			close(c.readRoutineQuit)
			c.reconnectAfter <- err
			return
		}

		var response types.RPCResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			c.Logger.Error("failed to parse response", "err", err, "data", string(data))
			continue
		}
		c.Logger.Info("got response", "resp", response.Result)

		select {
		case <-c.Quit():
		case c.ResponsesCh <- response:
		}
	}
}

func (c *WSClient) Subscribe(ctx context.Context, query string) error {
	params := map[string]interface{}{"query": query}
	return c.Call(ctx, "subscribe", params)
}

func (c *WSClient) Unsubscribe(ctx context.Context, query string) error {
	params := map[string]interface{}{"query": query}
	return c.Call(ctx, "unsubscribe", params)
}

func (c *WSClient) UnsubscribeAll(ctx context.Context) error {
	params := map[string]interface{}{}
	return c.Call(ctx, "unsubscribe_all", params)
}
