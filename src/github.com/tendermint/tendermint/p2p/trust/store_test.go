package trust

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

func TestTrustMetricStoreSaveLoad(t *testing.T) {
	dir, err := ioutil.TempDir("", "trust_test")
	if err != nil {
		panic(err)
	}
	defer os.Remove(dir)

	historyDB := dbm.NewDB("trusthistory", "goleveldb", dir)

	store := NewTrustMetricStore(historyDB, DefaultConfig())
	store.SetLogger(log.TestingLogger())
	store.saveToDB()

	store = NewTrustMetricStore(historyDB, DefaultConfig())
	store.SetLogger(log.TestingLogger())
	store.Start()

	assert.Zero(t, store.Size())

	var tt []*TestTicker
	for i := 0; i < 100; i++ {

		tt = append(tt, NewTestTicker())
	}

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("peer_%d", i)
		tm := NewMetric()

		tm.SetTicker(tt[i])
		tm.Start()
		store.AddPeerTrustMetric(key, tm)

		tm.BadEvents(10)
		tm.GoodEvents(1)
	}

	assert.Equal(t, 100, store.Size())

	for i := 0; i < 100; i++ {
		tt[i].NextTick()
		tt[i].NextTick()
	}

	store.Stop()

	store = NewTrustMetricStore(historyDB, DefaultConfig())
	store.SetLogger(log.TestingLogger())
	store.Start()

	assert.Equal(t, 100, store.Size())
	for _, tm := range store.peerMetrics {
		assert.NotEqual(t, 1.0, tm.TrustValue())
	}

	store.Stop()
}

func TestTrustMetricStoreConfig(t *testing.T) {
	historyDB := dbm.NewDB("", "memdb", "")

	config := TrustMetricConfig{
		ProportionalWeight:	0.5,
		IntegralWeight:		0.5,
	}

	store := NewTrustMetricStore(historyDB, config)
	store.SetLogger(log.TestingLogger())
	store.Start()

	tm := store.GetPeerTrustMetric("TestKey")

	assert.Equal(t, 0.5, tm.proportionalWeight)
	assert.Equal(t, 0.5, tm.integralWeight)
	store.Stop()
}

func TestTrustMetricStoreLookup(t *testing.T) {
	historyDB := dbm.NewDB("", "memdb", "")

	store := NewTrustMetricStore(historyDB, DefaultConfig())
	store.SetLogger(log.TestingLogger())
	store.Start()

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("peer_%d", i)
		store.GetPeerTrustMetric(key)

		ktm := store.peerMetrics[key]
		assert.NotNil(t, ktm, "Expected to find TrustMetric %s but wasn't there.", key)
	}

	store.Stop()
}

func TestTrustMetricStorePeerScore(t *testing.T) {
	historyDB := dbm.NewDB("", "memdb", "")

	store := NewTrustMetricStore(historyDB, DefaultConfig())
	store.SetLogger(log.TestingLogger())
	store.Start()

	key := "TestKey"
	tm := store.GetPeerTrustMetric(key)

	first := tm.TrustScore()
	assert.Equal(t, 100, first)

	tm.BadEvents(1)
	first = tm.TrustScore()
	assert.NotEqual(t, 100, first)
	tm.BadEvents(10)
	second := tm.TrustScore()

	if second > first {
		t.Errorf("A greater number of bad events should lower the trust score")
	}
	store.PeerDisconnected(key)

	tm = store.GetPeerTrustMetric(key)
	assert.NotEqual(t, 100, tm.TrustScore())
	store.Stop()
}
