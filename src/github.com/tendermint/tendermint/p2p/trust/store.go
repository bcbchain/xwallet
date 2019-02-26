package trust

import (
	"encoding/json"
	"sync"
	"time"

	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
)

const defaultStorePeriodicSaveInterval = 1 * time.Minute

var trustMetricKey = []byte("trustMetricStore")

type TrustMetricStore struct {
	cmn.BaseService

	peerMetrics	map[string]*TrustMetric

	mtx	sync.Mutex

	db	dbm.DB

	config	TrustMetricConfig
}

func NewTrustMetricStore(db dbm.DB, tmc TrustMetricConfig) *TrustMetricStore {
	tms := &TrustMetricStore{
		peerMetrics:	make(map[string]*TrustMetric),
		db:		db,
		config:		tmc,
	}

	tms.BaseService = *cmn.NewBaseService(nil, "TrustMetricStore", tms)
	return tms
}

func (tms *TrustMetricStore) OnStart() error {
	if err := tms.BaseService.OnStart(); err != nil {
		return err
	}

	tms.mtx.Lock()
	defer tms.mtx.Unlock()

	tms.loadFromDB()
	go tms.saveRoutine()
	return nil
}

func (tms *TrustMetricStore) OnStop() {
	tms.BaseService.OnStop()

	tms.mtx.Lock()
	defer tms.mtx.Unlock()

	for _, tm := range tms.peerMetrics {
		tm.Stop()
	}

	tms.saveToDB()
}

func (tms *TrustMetricStore) Size() int {
	tms.mtx.Lock()
	defer tms.mtx.Unlock()

	return tms.size()
}

func (tms *TrustMetricStore) AddPeerTrustMetric(key string, tm *TrustMetric) {
	tms.mtx.Lock()
	defer tms.mtx.Unlock()

	if key == "" || tm == nil {
		return
	}
	tms.peerMetrics[key] = tm
}

func (tms *TrustMetricStore) GetPeerTrustMetric(key string) *TrustMetric {
	tms.mtx.Lock()
	defer tms.mtx.Unlock()

	tm, ok := tms.peerMetrics[key]
	if !ok {

		tm = NewMetricWithConfig(tms.config)
		tm.Start()

		tms.peerMetrics[key] = tm
	}
	return tm
}

func (tms *TrustMetricStore) PeerDisconnected(key string) {
	tms.mtx.Lock()
	defer tms.mtx.Unlock()

	if tm, ok := tms.peerMetrics[key]; ok {
		tm.Pause()
	}
}

func (tms *TrustMetricStore) SaveToDB() {
	tms.mtx.Lock()
	defer tms.mtx.Unlock()

	tms.saveToDB()
}

func (tms *TrustMetricStore) size() int {
	return len(tms.peerMetrics)
}

func (tms *TrustMetricStore) loadFromDB() bool {

	bytes := tms.db.Get(trustMetricKey)
	if bytes == nil {
		return false
	}

	peers := make(map[string]MetricHistoryJSON)
	err := json.Unmarshal(bytes, &peers)
	if err != nil {
		cmn.PanicCrisis(cmn.Fmt("Could not unmarshal Trust Metric Store DB data: %v", err))
	}

	for key, p := range peers {
		tm := NewMetricWithConfig(tms.config)

		tm.Start()
		tm.Init(p)

		tms.peerMetrics[key] = tm
	}
	return true
}

func (tms *TrustMetricStore) saveToDB() {
	tms.Logger.Debug("Saving TrustHistory to DB", "size", tms.size())

	peers := make(map[string]MetricHistoryJSON)

	for key, tm := range tms.peerMetrics {

		peers[key] = tm.HistoryJSON()
	}

	bytes, err := json.Marshal(peers)
	if err != nil {
		tms.Logger.Error("Failed to encode the TrustHistory", "err", err)
		return
	}
	tms.db.SetSync(trustMetricKey, bytes)
}

func (tms *TrustMetricStore) saveRoutine() {
	t := time.NewTicker(defaultStorePeriodicSaveInterval)
	defer t.Stop()
loop:
	for {
		select {
		case <-t.C:
			tms.SaveToDB()
		case <-tms.Quit():
			break loop
		}
	}
}
