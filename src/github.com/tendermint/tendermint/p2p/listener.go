package p2p

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/tendermint/tendermint/p2p/upnp"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/log"
)

type Listener interface {
	Connections() <-chan net.Conn
	InternalAddress() *NetAddress
	ExternalAddress() *NetAddress
	String() string
	Stop() error
}

type DefaultListener struct {
	cmn.BaseService

	listener	net.Listener
	intAddr		*NetAddress
	extAddr		*NetAddress
	anoAddr		string
	connections	chan net.Conn
}

const (
	numBufferedConnections	= 10
	defaultExternalPort	= 8770
	tryListenSeconds	= 5
)

func splitHostPort(addr string) (host string, port int) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}
	port, err = strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	return host, port
}

func NewDefaultListener(protocol, lAddr, aAddr string, skipUPNP bool, logger log.Logger) Listener {

	lAddrIP, lAddrPort := splitHostPort(lAddr)

	var listener net.Listener
	var err error
	for i := 0; i < tryListenSeconds; i++ {
		listener, err = net.Listen(protocol, lAddr)
		if err == nil {
			break
		} else if i < tryListenSeconds-1 {
			time.Sleep(time.Second * 1)
		}
	}
	if err != nil {
		panic(err)
	}

	listenerIP, listenerPort := splitHostPort(listener.Addr().String())
	logger.Info("Local listener", "ip", listenerIP, "port", listenerPort)

	var intAddr *NetAddress
	intAddr, err = NewNetAddressStringWithOptionalID(lAddr)
	if err != nil {
		panic(err)
	}

	var extAddr *NetAddress
	if !skipUPNP {

		if lAddrIP == "" || lAddrIP == "0.0.0.0" {
			extAddr = getUPNPExternalAddress(lAddrPort, listenerPort, logger)
		}
	}

	if extAddr == nil {
		extAddr = getNaiveExternalAddress(listenerPort, false, logger)
	}
	if extAddr == nil {
		panic("Could not determine external address!")
	}

	dl := &DefaultListener{
		listener:	listener,
		intAddr:	intAddr,
		extAddr:	extAddr,
		anoAddr:	aAddr,
		connections:	make(chan net.Conn, numBufferedConnections),
	}
	dl.BaseService = *cmn.NewBaseService(logger, "DefaultListener", dl)
	err = dl.Start()
	if err != nil {
		logger.Error("Error starting base service", "err", err)
	}
	return dl
}

func (l *DefaultListener) OnStart() error {
	if err := l.BaseService.OnStart(); err != nil {
		return err
	}
	go l.listenRoutine()
	return nil
}

func (l *DefaultListener) OnStop() {
	l.BaseService.OnStop()
	l.listener.Close()
}

func (l *DefaultListener) listenRoutine() {
	for {
		conn, err := l.listener.Accept()

		if !l.IsRunning() {
			break
		}

		if err != nil {
			panic(err)
		}

		l.connections <- conn
	}

	close(l.connections)
	for range l.connections {

	}
}

func (l *DefaultListener) Connections() <-chan net.Conn {
	return l.connections
}

func (l *DefaultListener) InternalAddress() *NetAddress {
	return l.intAddr
}

func (l *DefaultListener) ExternalAddress() *NetAddress {
	return l.extAddr
}

func (l *DefaultListener) NetListener() net.Listener {
	return l.listener
}

func (l *DefaultListener) String() string {
	if l.anoAddr != "" {
		return l.anoAddr
	} else {
		return fmt.Sprintf("%v", l.extAddr)
	}
}

func getUPNPExternalAddress(externalPort, internalPort int, logger log.Logger) *NetAddress {
	logger.Info("Getting UPNP external address")
	nat, err := upnp.Discover()
	if err != nil {
		logger.Info("Could not perform UPNP discover", "err", err)
		return nil
	}

	ext, err := nat.GetExternalAddress()
	if err != nil {
		logger.Info("Could not get UPNP external address", "err", err)
		return nil
	}

	if externalPort == 0 {
		externalPort = defaultExternalPort
	}

	externalPort, err = nat.AddPortMapping("tcp", externalPort, internalPort, "tendermint", 0)
	if err != nil {
		logger.Info("Could not add UPNP port mapping", "err", err)
		return nil
	}

	logger.Info("Got UPNP external address", "address", ext)
	return NewNetAddressIPPort(ext, uint16(externalPort))
}

func getNaiveExternalAddress(port int, settleForLocal bool, logger log.Logger) *NetAddress {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(cmn.Fmt("Could not fetch interface addresses: %v", err))
	}

	for _, a := range addrs {
		ipnet, ok := a.(*net.IPNet)
		if !ok {
			continue
		}
		v4 := ipnet.IP.To4()
		if v4 == nil || (!settleForLocal && v4[0] == 127) {
			continue
		}
		return NewNetAddressIPPort(ipnet.IP, uint16(port))
	}

	logger.Info("Node may not be connected to internet. Settling for local address")
	return getNaiveExternalAddress(port, true, logger)
}
