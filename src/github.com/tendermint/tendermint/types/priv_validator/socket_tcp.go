package privval

import (
	"net"
	"time"
)

type timeoutError interface {
	Timeout() bool
}

var _ net.Listener = (*tcpTimeoutListener)(nil)

type tcpTimeoutListener struct {
	*net.TCPListener

	acceptDeadline	time.Duration
	connDeadline	time.Duration
	period		time.Duration
}

func newTCPTimeoutListener(
	ln net.Listener,
	acceptDeadline, connDeadline time.Duration,
	period time.Duration,
) tcpTimeoutListener {
	return tcpTimeoutListener{
		TCPListener:	ln.(*net.TCPListener),
		acceptDeadline:	acceptDeadline,
		connDeadline:	connDeadline,
		period:		period,
	}
}

func (ln tcpTimeoutListener) Accept() (net.Conn, error) {
	err := ln.SetDeadline(time.Now().Add(ln.acceptDeadline))
	if err != nil {
		return nil, err
	}

	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}

	if err := tc.SetDeadline(time.Now().Add(ln.connDeadline)); err != nil {
		return nil, err
	}

	if err := tc.SetKeepAlive(true); err != nil {
		return nil, err
	}

	if err := tc.SetKeepAlivePeriod(ln.period); err != nil {
		return nil, err
	}

	return tc, nil
}
