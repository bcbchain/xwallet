package trust

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTrustMetricScores(t *testing.T) {
	tm := NewMetric()
	tm.Start()

	tm.GoodEvents(1)
	score := tm.TrustScore()
	assert.Equal(t, 100, score)

	tm.BadEvents(10)
	score = tm.TrustScore()
	assert.NotEqual(t, 100, score)
	tm.Stop()
}

func TestTrustMetricConfig(t *testing.T) {

	window := time.Minute * 60 * 24 * 7
	config := TrustMetricConfig{
		TrackingWindow:	window,
		IntervalLength:	2 * time.Minute,
	}

	tm := NewMetricWithConfig(config)
	tm.Start()

	assert.Equal(t, int(config.TrackingWindow/config.IntervalLength), tm.maxIntervals)

	dc := DefaultConfig()

	assert.Equal(t, dc.ProportionalWeight, tm.proportionalWeight)
	assert.Equal(t, dc.IntegralWeight, tm.integralWeight)
	tm.Stop()
	tm.Wait()

	config.ProportionalWeight = 0.3
	config.IntegralWeight = 0.7
	tm = NewMetricWithConfig(config)
	tm.Start()

	assert.Equal(t, config.ProportionalWeight, tm.proportionalWeight)
	assert.Equal(t, config.IntegralWeight, tm.integralWeight)
	tm.Stop()
	tm.Wait()
}

func _TestTrustMetricStopPause(t *testing.T) {

	tt := NewTestTicker()
	tm := NewMetric()
	tm.SetTicker(tt)
	tm.Start()

	tt.NextTick()
	tt.NextTick()
	tm.Pause()

	first := tm.Copy().numIntervals

	tt.NextTick()
	tt.NextTick()
	assert.Equal(t, first, tm.Copy().numIntervals)

	tm.GoodEvents(5)

	tt.NextTick()
	tt.NextTick()
	tm.Stop()
	tm.Wait()

	second := tm.Copy().numIntervals

	tm.NextTimeInterval()
	tm.NextTimeInterval()

	assert.Equal(t, second+2, tm.Copy().numIntervals)

	if first > second {
		t.Fatalf("numIntervals should always increase or stay the same over time")
	}
}
