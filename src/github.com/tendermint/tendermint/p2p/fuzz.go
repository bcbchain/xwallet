package p2p

import (
	"net"
	"sync"
	"time"

	cmn "github.com/tendermint/tmlibs/common"
)

const (
	FuzzModeDrop	= iota

	FuzzModeDelay
)

type FuzzedConnection struct {
	conn	net.Conn

	mtx	sync.Mutex
	start	<-chan time.Time
	active	bool

	config	*FuzzConnConfig
}

type FuzzConnConfig struct {
	Mode		int
	MaxDelay	time.Duration
	ProbDropRW	float64
	ProbDropConn	float64
	ProbSleep	float64
}

func DefaultFuzzConnConfig() *FuzzConnConfig {
	return &FuzzConnConfig{
		Mode:		FuzzModeDrop,
		MaxDelay:	3 * time.Second,
		ProbDropRW:	0.2,
		ProbDropConn:	0.00,
		ProbSleep:	0.00,
	}
}

func FuzzConn(conn net.Conn) net.Conn {
	return FuzzConnFromConfig(conn, DefaultFuzzConnConfig())
}

func FuzzConnFromConfig(conn net.Conn, config *FuzzConnConfig) net.Conn {
	return &FuzzedConnection{
		conn:	conn,
		start:	make(<-chan time.Time),
		active:	true,
		config:	config,
	}
}

func FuzzConnAfter(conn net.Conn, d time.Duration) net.Conn {
	return FuzzConnAfterFromConfig(conn, d, DefaultFuzzConnConfig())
}

func FuzzConnAfterFromConfig(conn net.Conn, d time.Duration, config *FuzzConnConfig) net.Conn {
	return &FuzzedConnection{
		conn:	conn,
		start:	time.After(d),
		active:	false,
		config:	config,
	}
}

func (fc *FuzzedConnection) Config() *FuzzConnConfig {
	return fc.config
}

func (fc *FuzzedConnection) Read(data []byte) (n int, err error) {
	if fc.fuzz() {
		return 0, nil
	}
	return fc.conn.Read(data)
}

func (fc *FuzzedConnection) Write(data []byte) (n int, err error) {
	if fc.fuzz() {
		return 0, nil
	}
	return fc.conn.Write(data)
}

func (fc *FuzzedConnection) Close() error	{ return fc.conn.Close() }

func (fc *FuzzedConnection) LocalAddr() net.Addr	{ return fc.conn.LocalAddr() }

func (fc *FuzzedConnection) RemoteAddr() net.Addr	{ return fc.conn.RemoteAddr() }

func (fc *FuzzedConnection) SetDeadline(t time.Time) error	{ return fc.conn.SetDeadline(t) }

func (fc *FuzzedConnection) SetReadDeadline(t time.Time) error {
	return fc.conn.SetReadDeadline(t)
}

func (fc *FuzzedConnection) SetWriteDeadline(t time.Time) error {
	return fc.conn.SetWriteDeadline(t)
}

func (fc *FuzzedConnection) randomDuration() time.Duration {
	maxDelayMillis := int(fc.config.MaxDelay.Nanoseconds() / 1000)
	return time.Millisecond * time.Duration(cmn.RandInt()%maxDelayMillis)
}

func (fc *FuzzedConnection) fuzz() bool {
	if !fc.shouldFuzz() {
		return false
	}

	switch fc.config.Mode {
	case FuzzModeDrop:

		r := cmn.RandFloat64()
		if r <= fc.config.ProbDropRW {
			return true
		} else if r < fc.config.ProbDropRW+fc.config.ProbDropConn {

			fc.Close()
			return true
		} else if r < fc.config.ProbDropRW+fc.config.ProbDropConn+fc.config.ProbSleep {
			time.Sleep(fc.randomDuration())
		}
	case FuzzModeDelay:

		time.Sleep(fc.randomDuration())
	}
	return false
}

func (fc *FuzzedConnection) shouldFuzz() bool {
	if fc.active {
		return true
	}

	fc.mtx.Lock()
	defer fc.mtx.Unlock()

	select {
	case <-fc.start:
		fc.active = true
		return true
	default:
		return false
	}
}
