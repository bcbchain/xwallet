package trust

import (
	"time"
)

type MetricTicker interface {
	GetChannel() <-chan time.Time

	Stop()
}

type TestTicker struct {
	C	chan time.Time
	stopped	bool
}

func NewTestTicker() *TestTicker {
	c := make(chan time.Time)
	return &TestTicker{
		C: c,
	}
}

func (t *TestTicker) GetChannel() <-chan time.Time {
	return t.C
}

func (t *TestTicker) Stop() {
	t.stopped = true
}

func (t *TestTicker) NextTick() {
	if t.stopped {
		return
	}
	t.C <- time.Now()
}

type Ticker struct {
	*time.Ticker
}

func NewTicker(d time.Duration) *Ticker {
	return &Ticker{time.NewTicker(d)}
}

func (t *Ticker) GetChannel() <-chan time.Time {
	return t.C
}
