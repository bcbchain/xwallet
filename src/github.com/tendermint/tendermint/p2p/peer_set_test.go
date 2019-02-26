package p2p

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	crypto "github.com/tendermint/go-crypto"
	cmn "github.com/tendermint/tmlibs/common"
)

func randPeer() *peer {
	nodeKey := NodeKey{PrivKey: crypto.GenPrivKeyEd25519()}
	return &peer{
		nodeInfo: NodeInfo{
			ID:		nodeKey.ID(),
			ListenAddr:	cmn.Fmt("%v.%v.%v.%v:46656", rand.Int()%256, rand.Int()%256, rand.Int()%256, rand.Int()%256),
		},
	}
}

func TestPeerSetAddRemoveOne(t *testing.T) {
	t.Parallel()
	peerSet := NewPeerSet()

	var peerList []Peer
	for i := 0; i < 5; i++ {
		p := randPeer()
		if err := peerSet.Add(p); err != nil {
			t.Error(err)
		}
		peerList = append(peerList, p)
	}

	n := len(peerList)

	for i, peerAtFront := range peerList {
		peerSet.Remove(peerAtFront)
		wantSize := n - i - 1
		for j := 0; j < 2; j++ {
			assert.Equal(t, false, peerSet.Has(peerAtFront.ID()), "#%d Run #%d: failed to remove peer", i, j)
			assert.Equal(t, wantSize, peerSet.Size(), "#%d Run #%d: failed to remove peer and decrement size", i, j)

			peerSet.Remove(peerAtFront)
		}
	}

	for _, peer := range peerList {
		if err := peerSet.Add(peer); err != nil {
			t.Error(err)
		}
	}

	for i := n - 1; i >= 0; i-- {
		peerAtEnd := peerList[i]
		peerSet.Remove(peerAtEnd)
		assert.Equal(t, false, peerSet.Has(peerAtEnd.ID()), "#%d: failed to remove item at end", i)
		assert.Equal(t, i, peerSet.Size(), "#%d: differing sizes after peerSet.Remove(atEndPeer)", i)
	}
}

func TestPeerSetAddRemoveMany(t *testing.T) {
	t.Parallel()
	peerSet := NewPeerSet()

	peers := []Peer{}
	N := 100
	for i := 0; i < N; i++ {
		peer := randPeer()
		if err := peerSet.Add(peer); err != nil {
			t.Errorf("Failed to add new peer")
		}
		if peerSet.Size() != i+1 {
			t.Errorf("Failed to add new peer and increment size")
		}
		peers = append(peers, peer)
	}

	for i, peer := range peers {
		peerSet.Remove(peer)
		if peerSet.Has(peer.ID()) {
			t.Errorf("Failed to remove peer")
		}
		if peerSet.Size() != len(peers)-i-1 {
			t.Errorf("Failed to remove peer and decrement size")
		}
	}
}

func TestPeerSetAddDuplicate(t *testing.T) {
	t.Parallel()
	peerSet := NewPeerSet()
	peer := randPeer()

	n := 20
	errsChan := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			errsChan <- peerSet.Add(peer)
		}()
	}

	errsTally := make(map[error]int)
	for i := 0; i < n; i++ {
		err := <-errsChan
		errsTally[err]++
	}

	wantErrCount, gotErrCount := n-1, errsTally[ErrSwitchDuplicatePeer]
	assert.Equal(t, wantErrCount, gotErrCount, "invalid ErrSwitchDuplicatePeer count")

	wantNilErrCount, gotNilErrCount := 1, errsTally[nil]
	assert.Equal(t, wantNilErrCount, gotNilErrCount, "invalid nil errCount")
}

func TestPeerSetGet(t *testing.T) {
	t.Parallel()
	peerSet := NewPeerSet()
	peer := randPeer()
	assert.Nil(t, peerSet.Get(peer.ID()), "expecting a nil lookup, before .Add")

	if err := peerSet.Add(peer); err != nil {
		t.Fatalf("Failed to add new peer: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			got, want := peerSet.Get(peer.ID()), peer
			assert.Equal(t, got, want, "#%d: got=%v want=%v", i, got, want)
		}(i)
	}
	wg.Wait()
}
