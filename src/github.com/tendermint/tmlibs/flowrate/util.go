package flowrate

import (
	"math"
	"strconv"
	"time"
)

const clockRate = 20 * time.Millisecond

var czero = time.Now().Round(clockRate)

func clock() time.Duration {
	return time.Now().Round(clockRate).Sub(czero)
}

func clockToTime(c time.Duration) time.Time {
	return czero.Add(c)
}

func clockRound(d time.Duration) time.Duration {
	return (d + clockRate>>1) / clockRate * clockRate
}

func round(x float64) int64 {
	if _, frac := math.Modf(x); frac >= 0.5 {
		return int64(math.Ceil(x))
	}
	return int64(math.Floor(x))
}

type Percent uint32

func percentOf(x, total float64) Percent {
	if x < 0 || total <= 0 {
		return 0
	} else if p := round(x / total * 1e5); p <= math.MaxUint32 {
		return Percent(p)
	}
	return Percent(math.MaxUint32)
}

func (p Percent) Float() float64 {
	return float64(p) * 1e-3
}

func (p Percent) String() string {
	var buf [12]byte
	b := strconv.AppendUint(buf[:0], uint64(p)/1000, 10)
	n := len(b)
	b = strconv.AppendUint(b, 1000+uint64(p)%1000, 10)
	b[n] = '.'
	return string(append(b, '%'))
}
