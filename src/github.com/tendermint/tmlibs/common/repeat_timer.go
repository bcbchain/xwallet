package common

import (
	"sync"
	"time"
)

type TickerMaker func(dur time.Duration) Ticker

type Ticker interface {
	Chan() <-chan time.Time

	Stop()
}

var _ Ticker = (*defaultTicker)(nil)

type defaultTicker time.Ticker

func defaultTickerMaker(dur time.Duration) Ticker {
	ticker := time.NewTicker(dur)
	return (*defaultTicker)(ticker)
}

func (t *defaultTicker) Chan() <-chan time.Time {
	return t.C
}

func (t *defaultTicker) Stop() {
	((*time.Ticker)(t)).Stop()
}

func NewLogicalTickerMaker(source chan time.Time) TickerMaker {
	return func(dur time.Duration) Ticker {
		return newLogicalTicker(source, dur)
	}
}

type logicalTicker struct {
	source	<-chan time.Time
	ch	chan time.Time
	quit	chan struct{}
}

func newLogicalTicker(source <-chan time.Time, interval time.Duration) Ticker {
	lt := &logicalTicker{
		source:	source,
		ch:	make(chan time.Time),
		quit:	make(chan struct{}),
	}
	go lt.fireRoutine(interval)
	return lt
}

func (t *logicalTicker) fireRoutine(interval time.Duration) {
	source := t.source

	lasttime := time.Time{}
	select {
	case lasttime = <-source:
	case <-t.quit:
		return
	}

	for {
		select {
		case newtime := <-source:
			elapsed := newtime.Sub(lasttime)
			if interval <= elapsed {

				select {
				case t.ch <- newtime:
				case <-t.quit:
					return
				}

				lasttime = newtime
			}
		case <-t.quit:
			return
		}
	}
}

func (t *logicalTicker) Chan() <-chan time.Time {
	return t.ch
}

func (t *logicalTicker) Stop() {
	close(t.quit)
}

type RepeatTimer struct {
	name	string
	ch	chan time.Time
	tm	TickerMaker

	mtx	sync.Mutex
	dur	time.Duration
	ticker	Ticker
	quit	chan struct{}
}

func NewRepeatTimer(name string, dur time.Duration) *RepeatTimer {
	return NewRepeatTimerWithTickerMaker(name, dur, defaultTickerMaker)
}

func NewRepeatTimerWithTickerMaker(name string, dur time.Duration, tm TickerMaker) *RepeatTimer {
	var t = &RepeatTimer{
		name:	name,
		ch:	make(chan time.Time),
		tm:	tm,
		dur:	dur,
		ticker:	nil,
		quit:	nil,
	}
	t.reset()
	return t
}

func (t *RepeatTimer) fireRoutine(ch <-chan time.Time, quit <-chan struct{}) {
	for {
		select {
		case tick := <-ch:
			select {
			case t.ch <- tick:
			case <-quit:
				return
			}
		case <-quit:
			return
		}
	}
}

func (t *RepeatTimer) Chan() <-chan time.Time {
	return t.ch
}

func (t *RepeatTimer) Stop() {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.stop()
}

func (t *RepeatTimer) Reset() {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.reset()
}

func (t *RepeatTimer) reset() {
	if t.ticker != nil {
		t.stop()
	}
	t.ticker = t.tm(t.dur)
	t.quit = make(chan struct{})
	go t.fireRoutine(t.ticker.Chan(), t.quit)
}

func (t *RepeatTimer) stop() {
	if t.ticker == nil {

		panic("Tried to stop a stopped RepeatTimer")
	}
	t.ticker.Stop()
	t.ticker = nil

	close(t.quit)
}
