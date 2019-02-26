package p2p

import (
	"fmt"
	"strings"

	cmn "github.com/tendermint/tmlibs/common"
)

const (
	maxNodeInfoSize	= 10240
	maxNumChannels	= 16
)

func MaxNodeInfoSize() int {
	return maxNodeInfoSize
}

type NodeInfo struct {
	ID		ID	`json:"id"`
	ListenAddr	string	`json:"listen_addr"`

	Network		string		`json:"network"`
	Channels	cmn.HexBytes	`json:"channels"`
	Version		string		`json:"version"`
	BaseVersion	string		`json:"base_version"`

	Moniker	string		`json:"moniker"`
	Other	[]string	`json:"other"`
}

func (info NodeInfo) Validate() error {
	if len(info.Channels) > maxNumChannels {
		return fmt.Errorf("info.Channels is too long (%v). Max is %v", len(info.Channels), maxNumChannels)
	}

	channels := make(map[byte]struct{})
	for _, ch := range info.Channels {
		_, ok := channels[ch]
		if ok {
			return fmt.Errorf("info.Channels contains duplicate channel id %v", ch)
		}
		channels[ch] = struct{}{}
	}
	return nil
}

func (info NodeInfo) CompatibleWith(other NodeInfo) error {
	iMajor, iMinor, _, iErr := splitVersion(info.Version)
	oMajor, oMinor, _, oErr := splitVersion(other.Version)

	if iErr != nil {
		return iErr
	}

	if oErr != nil {
		return oErr
	}

	if iMajor != oMajor {
		return fmt.Errorf("Peer is on a different major version. Got %v, expected %v", oMajor, iMajor)
	}

	if iMinor != oMinor {
		return fmt.Errorf("Peer is on a different minor version. Got %v, expected %v", oMinor, iMinor)
	}

	if info.Network != other.Network {
		return fmt.Errorf("Peer is on a different network. Got %v, expected %v", other.Network, info.Network)
	}

	if len(info.Channels) == 0 {
		return nil
	}

	found := false
OUTER_LOOP:
	for _, ch1 := range info.Channels {
		for _, ch2 := range other.Channels {
			if ch1 == ch2 {
				found = true
				break OUTER_LOOP
			}
		}
	}
	if !found {
		return fmt.Errorf("Peer has no common channels. Our channels: %v ; Peer channels: %v", info.Channels, other.Channels)
	}
	return nil
}

func (info NodeInfo) NetAddress() *NetAddress {
	netAddr, err := NewNetAddressString(IDAddressString(info.ID, info.ListenAddr))
	if err != nil {
		panic(err)
	}
	return netAddr
}

func (info NodeInfo) String() string {
	return fmt.Sprintf("NodeInfo{id: %v, moniker: %v, network: %v [listen %v], version: %v (%v)}",
		info.ID, info.Moniker, info.Network, info.ListenAddr, info.Version, info.Other)
}

func splitVersion(version string) (string, string, string, error) {
	spl := strings.Split(version, ".")
	if len(spl) != 4 {
		return "", "", "", fmt.Errorf("Invalid version format %v", version)
	}
	return spl[0], spl[1], spl[2], nil
}
