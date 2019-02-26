package p2p

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	cmn "github.com/tendermint/tmlibs/common"
)

type NetAddress struct {
	ID	ID
	IP	net.IP
	Port	uint16
	str	string
}

func IDAddressString(id ID, hostPort string) string {
	return fmt.Sprintf("%s@%s", id, hostPort)
}

func NewNetAddress(id ID, addr net.Addr) *NetAddress {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		if flag.Lookup("test.v") == nil {
			cmn.PanicSanity(cmn.Fmt("Only TCPAddrs are supported. Got: %v", addr))
		} else {
			netAddr := NewNetAddressIPPort(net.IP("0.0.0.0"), 0)
			netAddr.ID = id
			return netAddr
		}
	}
	ip := tcpAddr.IP
	port := uint16(tcpAddr.Port)
	na := NewNetAddressIPPort(ip, port)
	na.ID = id
	return na
}

func NewNetAddressString(addr string) (*NetAddress, error) {
	spl := strings.Split(addr, "@")
	if len(spl) < 2 {
		return nil, fmt.Errorf("Address (%s) does not contain ID", addr)
	}
	return NewNetAddressStringWithOptionalID(addr)
}

func NewNetAddressStringWithOptionalID(addr string) (*NetAddress, error) {
	addrWithoutProtocol := removeProtocolIfDefined(addr)

	var id ID
	spl := strings.Split(addrWithoutProtocol, "@")
	if len(spl) == 2 {
		idStr := spl[0]
		id, addrWithoutProtocol = ID(idStr), spl[1]
	}

	host, portStr, err := net.SplitHostPort(addrWithoutProtocol)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(host)
	if ip == nil {
		if len(host) > 0 {
			ips, err := net.LookupIP(host)
			if err != nil {
				return nil, err
			}
			ip = ips[0]
		}
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, err
	}

	na := NewNetAddressIPPort(ip, uint16(port))
	na.ID = id
	return na, nil
}

func NewNetAddressStrings(addrs []string) ([]*NetAddress, []error) {
	netAddrs := make([]*NetAddress, 0)
	errs := make([]error, 0)
	for _, addr := range addrs {
		netAddr, err := NewNetAddressString(addr)
		if err != nil {
			errs = append(errs, fmt.Errorf("Error in address %s: %v", addr, err))
		} else {
			netAddrs = append(netAddrs, netAddr)
		}
	}
	return netAddrs, errs
}

func NewNetAddressIPPort(ip net.IP, port uint16) *NetAddress {
	return &NetAddress{
		IP:	ip,
		Port:	port,
	}
}

func (na *NetAddress) Equals(other interface{}) bool {
	if o, ok := other.(*NetAddress); ok {
		return na.String() == o.String()
	}
	return false
}

func (na *NetAddress) Same(other interface{}) bool {
	if o, ok := other.(*NetAddress); ok {
		if na.DialString() == o.DialString() {
			return true
		}
		if na.ID != "" && na.ID == o.ID {
			return true
		}
	}
	return false
}

func (na *NetAddress) String() string {
	if na.str == "" {
		addrStr := na.DialString()
		if na.ID != "" {
			addrStr = IDAddressString(na.ID, addrStr)
		}
		na.str = addrStr
	}
	return na.str
}

func (na *NetAddress) DialString() string {
	return net.JoinHostPort(
		na.IP.String(),
		strconv.FormatUint(uint64(na.Port), 10),
	)
}

func (na *NetAddress) Dial() (net.Conn, error) {
	conn, err := net.Dial("tcp", na.DialString())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (na *NetAddress) DialTimeout(timeout time.Duration) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", na.DialString(), timeout)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (na *NetAddress) Routable() bool {

	return na.Valid() && !(na.RFC1918() || na.RFC3927() || na.RFC4862() ||
		na.RFC4193() || na.RFC4843() || na.Local())
}

func (na *NetAddress) Valid() bool {
	return na.IP != nil && !(na.IP.IsUnspecified() || na.RFC3849() ||
		na.IP.Equal(net.IPv4bcast))
}

func (na *NetAddress) Local() bool {
	return na.IP.IsLoopback() || zero4.Contains(na.IP)
}

func (na *NetAddress) ReachabilityTo(o *NetAddress) int {
	const (
		Unreachable	= 0
		Default		= iota
		Teredo
		Ipv6_weak
		Ipv4
		Ipv6_strong
	)
	if !na.Routable() {
		return Unreachable
	} else if na.RFC4380() {
		if !o.Routable() {
			return Default
		} else if o.RFC4380() {
			return Teredo
		} else if o.IP.To4() != nil {
			return Ipv4
		} else {
			return Ipv6_weak
		}
	} else if na.IP.To4() != nil {
		if o.Routable() && o.IP.To4() != nil {
			return Ipv4
		}
		return Default
	} else {
		var tunnelled bool

		if o.RFC3964() || o.RFC6052() || o.RFC6145() {
			tunnelled = true
		}
		if !o.Routable() {
			return Default
		} else if o.RFC4380() {
			return Teredo
		} else if o.IP.To4() != nil {
			return Ipv4
		} else if tunnelled {

			return Ipv6_weak
		}
		return Ipv6_strong
	}
}

var rfc1918_10 = net.IPNet{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)}
var rfc1918_192 = net.IPNet{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)}
var rfc1918_172 = net.IPNet{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)}
var rfc3849 = net.IPNet{IP: net.ParseIP("2001:0DB8::"), Mask: net.CIDRMask(32, 128)}
var rfc3927 = net.IPNet{IP: net.ParseIP("169.254.0.0"), Mask: net.CIDRMask(16, 32)}
var rfc3964 = net.IPNet{IP: net.ParseIP("2002::"), Mask: net.CIDRMask(16, 128)}
var rfc4193 = net.IPNet{IP: net.ParseIP("FC00::"), Mask: net.CIDRMask(7, 128)}
var rfc4380 = net.IPNet{IP: net.ParseIP("2001::"), Mask: net.CIDRMask(32, 128)}
var rfc4843 = net.IPNet{IP: net.ParseIP("2001:10::"), Mask: net.CIDRMask(28, 128)}
var rfc4862 = net.IPNet{IP: net.ParseIP("FE80::"), Mask: net.CIDRMask(64, 128)}
var rfc6052 = net.IPNet{IP: net.ParseIP("64:FF9B::"), Mask: net.CIDRMask(96, 128)}
var rfc6145 = net.IPNet{IP: net.ParseIP("::FFFF:0:0:0"), Mask: net.CIDRMask(96, 128)}
var zero4 = net.IPNet{IP: net.ParseIP("0.0.0.0"), Mask: net.CIDRMask(8, 32)}

func (na *NetAddress) RFC1918() bool {
	return rfc1918_10.Contains(na.IP) ||
		rfc1918_192.Contains(na.IP) ||
		rfc1918_172.Contains(na.IP)
}
func (na *NetAddress) RFC3849() bool	{ return rfc3849.Contains(na.IP) }
func (na *NetAddress) RFC3927() bool	{ return rfc3927.Contains(na.IP) }
func (na *NetAddress) RFC3964() bool	{ return rfc3964.Contains(na.IP) }
func (na *NetAddress) RFC4193() bool	{ return rfc4193.Contains(na.IP) }
func (na *NetAddress) RFC4380() bool	{ return rfc4380.Contains(na.IP) }
func (na *NetAddress) RFC4843() bool	{ return rfc4843.Contains(na.IP) }
func (na *NetAddress) RFC4862() bool	{ return rfc4862.Contains(na.IP) }
func (na *NetAddress) RFC6052() bool	{ return rfc6052.Contains(na.IP) }
func (na *NetAddress) RFC6145() bool	{ return rfc6145.Contains(na.IP) }

func removeProtocolIfDefined(addr string) string {
	if strings.Contains(addr, "://") {
		return strings.Split(addr, "://")[1]
	}
	return addr

}
