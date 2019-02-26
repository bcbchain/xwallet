package trust

import (
	"math"
	"sync"
	"time"

	cmn "github.com/tendermint/tmlibs/common"
)

const (
	defaultDerivativeGamma1	= 0

	defaultDerivativeGamma2	= 1.0

	defaultHistoryDataWeight	= 0.8
)

type MetricHistoryJSON struct {
	NumIntervals	int		`json:"intervals"`
	History		[]float64	`json:"history"`
}

type TrustMetric struct {
	cmn.BaseService

	mtx	sync.Mutex

	proportionalWeight	float64

	integralWeight	float64

	numIntervals	int

	maxIntervals	int

	intervalLen	time.Duration

	history	[]float64

	historyWeights	[]float64

	historyWeightSum	float64

	historySize	int

	historyMaxSize	int

	historyValue	float64

	bad, good	float64

	paused	bool

	testTicker	MetricTicker
}

func NewMetric() *TrustMetric {
	return NewMetricWithConfig(DefaultConfig())
}

func NewMetricWithConfig(tmc TrustMetricConfig) *TrustMetric {
	tm := new(TrustMetric)
	config := customConfig(tmc)

	tm.proportionalWeight = config.ProportionalWeight
	tm.integralWeight = config.IntegralWeight
	tm.intervalLen = config.IntervalLength

	tm.maxIntervals = int(config.TrackingWindow / tm.intervalLen)

	tm.historyMaxSize = intervalToHistoryOffset(tm.maxIntervals) + 1

	tm.historyValue = 1.0

	tm.BaseService = *cmn.NewBaseService(nil, "TrustMetric", tm)
	return tm
}

func (tm *TrustMetric) OnStart() error {
	if err := tm.BaseService.OnStart(); err != nil {
		return err
	}
	go tm.processRequests()
	return nil
}

func (tm *TrustMetric) OnStop()	{}

func (tm *TrustMetric) HistoryJSON() MetricHistoryJSON {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()

	return MetricHistoryJSON{
		NumIntervals:	tm.numIntervals,
		History:	tm.history,
	}
}

func (tm *TrustMetric) Init(hist MetricHistoryJSON) {

	if hist.NumIntervals > tm.maxIntervals {
		hist.NumIntervals = tm.maxIntervals
	}
	tm.numIntervals = hist.NumIntervals

	if len(hist.History) > tm.historyMaxSize {

		last := len(hist.History) - tm.historyMaxSize
		hist.History = hist.History[last:]
	}
	tm.history = hist.History
	tm.historySize = len(tm.history)

	for i := 1; i <= tm.numIntervals; i++ {
		x := math.Pow(defaultHistoryDataWeight, float64(i))
		tm.historyWeights = append(tm.historyWeights, x)
	}

	for _, v := range tm.historyWeights {
		tm.historyWeightSum += v
	}

	tm.historyValue = tm.calcHistoryValue()
}

func (tm *TrustMetric) Pause() {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()

	tm.paused = true
}

func (tm *TrustMetric) BadEvents(num int) {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()

	tm.unpause()
	tm.bad += float64(num)
}

func (tm *TrustMetric) GoodEvents(num int) {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()

	tm.unpause()
	tm.good += float64(num)
}

func (tm *TrustMetric) TrustValue() float64 {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()

	return tm.calcTrustValue()
}

func (tm *TrustMetric) TrustScore() int {
	score := tm.TrustValue() * 100

	return int(math.Floor(score))
}

func (tm *TrustMetric) NextTimeInterval() {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()

	if tm.paused {

		return
	}

	newHist := tm.calcTrustValue()
	tm.history = append(tm.history, newHist)

	if tm.historySize < tm.historyMaxSize {
		tm.historySize++
	} else {

		last := len(tm.history) - tm.historyMaxSize
		tm.history = tm.history[last:]
	}

	if tm.numIntervals < tm.maxIntervals {
		tm.numIntervals++

		wk := math.Pow(defaultHistoryDataWeight, float64(tm.numIntervals))
		tm.historyWeights = append(tm.historyWeights, wk)
		tm.historyWeightSum += wk
	}

	tm.updateFadedMemory()

	tm.historyValue = tm.calcHistoryValue()
	tm.good = 0
	tm.bad = 0
}

func (tm *TrustMetric) SetTicker(ticker MetricTicker) {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()

	tm.testTicker = ticker
}

func (tm *TrustMetric) Copy() *TrustMetric {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()
	if tm == nil {
		return nil
	}

	return &TrustMetric{
		proportionalWeight:	tm.proportionalWeight,
		integralWeight:		tm.integralWeight,
		numIntervals:		tm.numIntervals,
		maxIntervals:		tm.maxIntervals,
		intervalLen:		tm.intervalLen,
		history:		tm.history,
		historyWeights:		tm.historyWeights,
		historyWeightSum:	tm.historyWeightSum,
		historySize:		tm.historySize,
		historyMaxSize:		tm.historyMaxSize,
		historyValue:		tm.historyValue,
		good:			tm.good,
		bad:			tm.bad,
		paused:			tm.paused,
	}

}

func (tm *TrustMetric) processRequests() {
	t := tm.testTicker
	if t == nil {

		t = NewTicker(tm.intervalLen)
	}
	defer t.Stop()

	tick := t.GetChannel()
loop:
	for {
		select {
		case <-tick:
			tm.NextTimeInterval()
		case <-tm.Quit():

			break loop
		}
	}
}

func (tm *TrustMetric) unpause() {

	if tm.paused {
		tm.good = 0
		tm.bad = 0

		tm.paused = false
	}
}

func (tm *TrustMetric) calcTrustValue() float64 {
	weightedP := tm.proportionalWeight * tm.proportionalValue()
	weightedI := tm.integralWeight * tm.historyValue
	weightedD := tm.weightedDerivative()

	tv := weightedP + weightedI + weightedD

	if tv < 0 {
		tv = 0
	}
	return tv
}

func (tm *TrustMetric) proportionalValue() float64 {
	value := 1.0

	total := tm.good + tm.bad
	if total > 0 {
		value = tm.good / total
	}
	return value
}

func (tm *TrustMetric) weightedDerivative() float64 {
	var weight float64 = defaultDerivativeGamma1

	d := tm.derivativeValue()
	if d < 0 {
		weight = defaultDerivativeGamma2
	}
	return weight * d
}

func (tm *TrustMetric) derivativeValue() float64 {
	return tm.proportionalValue() - tm.historyValue
}

func (tm *TrustMetric) calcHistoryValue() float64 {
	var hv float64

	for i := 0; i < tm.numIntervals; i++ {
		hv += tm.fadedMemoryValue(i) * tm.historyWeights[i]
	}

	return hv / tm.historyWeightSum
}

func (tm *TrustMetric) fadedMemoryValue(interval int) float64 {
	first := tm.historySize - 1

	if interval == 0 {

		return tm.history[first]
	}

	offset := intervalToHistoryOffset(interval)
	return tm.history[first-offset]
}

func (tm *TrustMetric) updateFadedMemory() {
	if tm.historySize < 2 {
		return
	}

	end := tm.historySize - 1

	for count := 1; count < tm.historySize; count++ {
		i := end - count

		x := math.Pow(2, float64(count))

		tm.history[i] = ((tm.history[i] * (x - 1)) + tm.history[i+1]) / x
	}
}

func intervalToHistoryOffset(interval int) int {

	return int(math.Floor(math.Log2(float64(interval))))
}
