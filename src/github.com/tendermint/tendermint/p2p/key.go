package p2p

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
)

type ID string

const IDByteLength = 20

type NodeKey struct {
	PrivKey crypto.PrivKey `json:"priv_key"`
}

func (nodeKey *NodeKey) ID() ID {
	return PubKeyToID(nodeKey.PubKey())
}

func (nodeKey *NodeKey) PubKey() crypto.PubKey {
	return nodeKey.PrivKey.PubKey()
}

func PubKeyToID(pubKey crypto.PubKey) ID {
	return ID(pubKey.Address())
}

func LoadOrGenNodeKey(filePath string) (*NodeKey, error) {
	if cmn.FileExists(filePath) {
		nodeKey, err := LoadNodeKey(filePath)
		if err != nil {
			return nil, err
		}
		return nodeKey, nil
	}
	return genNodeKey(filePath)
}

func LoadNodeKey(filePath string) (*NodeKey, error) {
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	nodeKey := new(NodeKey)
	privKey := crypto.PrivKeyEd25519{}
	err = wire.UnmarshalJSON(jsonBytes, &privKey)
	if err != nil {
		return nil, fmt.Errorf("Error reading NodeKey from %v: %v", filePath, err)
	}
	nodeKey.PrivKey = privKey
	return nodeKey, nil
}

func genNodeKey(filePath string) (*NodeKey, error) {
	privKey := crypto.GenPrivKeyEd25519()
	nodeKey := &NodeKey{
		PrivKey: privKey,
	}

	jsonBytes, err := wire.MarshalJSON(privKey)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(filePath, jsonBytes, 0600)
	if err != nil {
		return nil, err
	}
	return nodeKey, nil
}

func MakePoWTarget(difficulty, targetBits uint) []byte {
	if targetBits%8 != 0 {
		panic(fmt.Sprintf("targetBits (%d) not a multiple of 8", targetBits))
	}
	if difficulty >= targetBits {
		panic(fmt.Sprintf("difficulty (%d) >= targetBits (%d)", difficulty, targetBits))
	}
	targetBytes := targetBits / 8
	zeroPrefixLen := (int(difficulty) / 8)
	prefix := bytes.Repeat([]byte{0}, zeroPrefixLen)
	mod := (difficulty % 8)
	if mod > 0 {
		nonZeroPrefix := byte(1<<(8-mod) - 1)
		prefix = append(prefix, nonZeroPrefix)
	}
	tailLen := int(targetBytes) - len(prefix)
	return append(prefix, bytes.Repeat([]byte{0xFF}, tailLen)...)
}
