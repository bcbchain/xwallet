package flowrate

import (
	"math"
	"sync"
	"time"
)

type Monitor struct {
	mu	sync.Mutex
	active	bool
	start	time.Duration
	bytes	int64
	samples	int64

	rSample	float64
	rEMA	float64
	rPeak	float64
	rWindow	float64

	sBytes	int64
	sLast	time.Duration
	sRate	time.Duration

	tBytes	int64
	tLast	time.Duration
}

func New(sampleRate, windowSize time.Duration) *Monitor {
	if sampleRate = clockRound(sampleRate); sampleRate <= 0 {
		sampleRate = 5 * clockRate
	}
	if windowSize <= 0 {
		windowSize = 1 * time.Second
	}
	now := clock()
	return &Monitor{
		active:		true,
		start:		now,
		rWindow:	windowSize.Seconds(),
		sLast:		now,
		sRate:		sampleRate,
		tLast:		now,
	}
}

func (m *Monitor) Update(n int) int {
	m.mu.Lock()
	m.update(n)
	m.mu.Unlock()
	return n
}

func (m *Monitor) SetREMA(rEMA float64) {
	m.mu.Lock()
	m.rEMA = rEMA
	m.samples++
	m.mu.Unlock()
}

func (m *Monitor) IO(n int, err error) (int, error) {
	return m.Update(n), err
}

func (m *Monitor) Done() int64 {
	m.mu.Lock()
	if now := m.update(0); m.sBytes > 0 {
		m.reset(now)
	}
	m.active = false
	m.tLast = 0
	n := m.bytes
	m.mu.Unlock()
	return n
}

const timeRemLimit = 999*time.Hour + 59*time.Minute + 59*time.Second

type Status struct {
	Active		bool
	Start		time.Time
	Duration	time.Duration
	Idle		time.Duration
	Bytes		int64
	Samples		int64
	InstRate	int64
	CurRate		int64
	AvgRate		int64
	PeakRate	int64
	BytesRem	int64
	TimeRem		time.Duration
	Progress	Percent
}

func (m *Monitor) Status() Status {
	m.mu.Lock()
	now := m.update(0)
	s := Status{
		Active:		m.active,
		Start:		clockToTime(m.start),
		Duration:	m.sLast - m.start,
		Idle:		now - m.tLast,
		Bytes:		m.bytes,
		Samples:	m.samples,
		PeakRate:	round(m.rPeak),
		BytesRem:	m.tBytes - m.bytes,
		Progress:	percentOf(float64(m.bytes), float64(m.tBytes)),
	}
	if s.BytesRem < 0 {
		s.BytesRem = 0
	}
	if s.Duration > 0 {
		rAvg := float64(s.Bytes) / s.Duration.Seconds()
		s.AvgRate = round(rAvg)
		if s.Active {
			s.InstRate = round(m.rSample)
			s.CurRate = round(m.rEMA)
			if s.BytesRem > 0 {
				if tRate := 0.8*m.rEMA + 0.2*rAvg; tRate > 0 {
					ns := float64(s.BytesRem) / tRate * 1e9
					if ns > float64(timeRemLimit) {
						ns = float64(timeRemLimit)
					}
					s.TimeRem = clockRound(time.Duration(ns))
				}
			}
		}
	}
	m.mu.Unlock()
	return s
}

func (m *Monitor) Limit(want int, rate int64, block bool) (n int) {
	if want < 1 || rate < 1 {
		return want
	}
	m.mu.Lock()

	limit := round(float64(rate) * m.sRate.Seconds())
	if limit <= 0 {
		limit = 1
	}

	if now := m.update(0); block {
		for m.sBytes >= limit && m.active {
			now = m.waitNextSample(now)
		}
	}

	if limit -= m.sBytes; limit > int64(want) || !m.active {
		limit = int64(want)
	}
	m.mu.Unlock()

	if limit < 0 {
		limit = 0
	}
	return int(limit)
}

func (m *Monitor) SetTransferSize(bytes int64) {
	if bytes < 0 {
		bytes = 0
	}
	m.mu.Lock()
	m.tBytes = bytes
	m.mu.Unlock()
}

func (m *Monitor) update(n int) (now time.Duration) {
	if !m.active {
		return
	}
	if now = clock(); n > 0 {
		m.tLast = now
	}
	m.sBytes += int64(n)
	if sTime := now - m.sLast; sTime >= m.sRate {
		t := sTime.Seconds()
		if m.rSample = float64(m.sBytes) / t; m.rSample > m.rPeak {
			m.rPeak = m.rSample
		}

		if m.samples > 0 {
			w := math.Exp(-t / m.rWindow)
			m.rEMA = m.rSample + w*(m.rEMA-m.rSample)
		} else {
			m.rEMA = m.rSample
		}
		m.reset(now)
	}
	return
}

func (m *Monitor) reset(sampleTime time.Duration) {
	m.bytes += m.sBytes
	m.samples++
	m.sBytes = 0
	m.sLast = sampleTime
}

func (m *Monitor) waitNextSample(now time.Duration) time.Duration {
	const minWait = 5 * time.Millisecond
	current := m.sLast

	for m.sLast == current && m.active {
		d := current + m.sRate - now
		m.mu.Unlock()
		if d < minWait {
			d = minWait
		}
		time.Sleep(d)
		m.mu.Lock()
		now = m.update(0)
	}
	return now
}
