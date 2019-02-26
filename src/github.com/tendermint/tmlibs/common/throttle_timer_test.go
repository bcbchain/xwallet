package common

import (
	"sync"
	"testing"
	"time"

	asrt "github.com/stretchr/testify/assert"
)

type thCounter struct {
	input	chan struct{}
	mtx	sync.Mutex
	count	int
}

func (c *thCounter) Increment() {
	c.mtx.Lock()
	c.count++
	c.mtx.Unlock()
}

func (c *thCounter) Count() int {
	c.mtx.Lock()
	val := c.count
	c.mtx.Unlock()
	return val
}

func (c *thCounter) Read() {
	for range c.input {
		c.Increment()
	}
}

func TestThrottle(test *testing.T) {
	assert := asrt.New(test)

	ms := 50
	delay := time.Duration(ms) * time.Millisecond
	longwait := time.Duration(2) * delay
	t := NewThrottleTimer("foo", delay)

	c := &thCounter{input: t.Ch}
	assert.Equal(0, c.Count())
	go c.Read()

	time.Sleep(longwait)
	assert.Equal(0, c.Count())

	t.Set()
	time.Sleep(longwait)
	assert.Equal(1, c.Count())

	for i := 0; i < 5; i++ {
		t.Set()
	}
	time.Sleep(longwait)
	assert.Equal(2, c.Count())

	short := time.Duration(ms/5) * time.Millisecond
	for i := 0; i < 13; i++ {
		t.Set()
		time.Sleep(short)
	}
	time.Sleep(longwait)
	assert.Equal(5, c.Count())

	close(t.Ch)
}
