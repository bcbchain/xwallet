package common

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/assert"
)

func TestDefaultTicker(t *testing.T) {
	ticker := defaultTickerMaker(time.Millisecond * 10)
	<-ticker.Chan()
	ticker.Stop()
}

func TestRepeatTimer(t *testing.T) {

	ch := make(chan time.Time, 100)
	mtx := new(sync.Mutex)

	tick := func(startMs, endMs, incrMs time.Duration) {
		mtx.Lock()
		go func() {
			for tMs := startMs; tMs < endMs; tMs += incrMs {
				lt := time.Time{}
				lt = lt.Add(tMs * time.Millisecond)
				ch <- lt
			}
			mtx.Unlock()
		}()
	}

	tock := func(t *testing.T, rt *RepeatTimer, timesMs []int64) {

		for _, timeMs := range timesMs {
			tyme := <-rt.Chan()
			sinceMs := tyme.Sub(time.Time{}) / time.Millisecond
			assert.Equal(t, timeMs, int64(sinceMs))
		}

		time.Sleep(time.Millisecond * 100)
		done := true
		select {
		case <-rt.Chan():
			done = false
		default:
		}
		assert.True(t, done)
	}

	tm := NewLogicalTickerMaker(ch)
	rt := NewRepeatTimerWithTickerMaker("bar", time.Second, tm)

	tick(0, 1000, 10)
	tock(t, rt, []int64{})
	tick(1000, 2000, 10)
	tock(t, rt, []int64{1000})
	tick(2005, 5000, 10)
	tock(t, rt, []int64{2005, 3005, 4005})
	tick(5001, 5999, 1)

	tock(t, rt, []int64{5005})
	tick(6000, 7005, 1)
	tock(t, rt, []int64{6005})
	tick(7033, 8032, 1)
	tock(t, rt, []int64{7033})

	rt.Reset()
	tock(t, rt, []int64{})
	tick(8040, 8041, 1)
	tock(t, rt, []int64{})
	tick(9555, 9556, 1)
	tock(t, rt, []int64{9555})

	rt.Stop()
	tock(t, rt, []int64{})

	assert.Panics(t, func() { rt.Stop() })
}

func TestRepeatTimerReset(t *testing.T) {

	defer leaktest.Check(t)()

	timer := NewRepeatTimer("test", 20*time.Millisecond)
	defer timer.Stop()

	select {
	case <-timer.Chan():
		t.Fatal("did not expect to receive tick")
	default:
	}

	timer.Reset()

	select {
	case <-timer.Chan():

	case <-time.After(40 * time.Millisecond):
		t.Fatal("expected to receive tick after reset")
	}

	for i := 0; i < 100; i++ {
		time.Sleep(time.Duration(rand.Intn(40)) * time.Millisecond)
		timer.Reset()
	}
}
