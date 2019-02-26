package events

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddListenerForEventFireOnce(t *testing.T) {
	evsw := NewEventSwitch()
	err := evsw.Start()
	if err != nil {
		t.Errorf("Failed to start EventSwitch, error: %v", err)
	}
	messages := make(chan EventData)
	evsw.AddListenerForEvent("listener", "event",
		func(data EventData) {
			messages <- data
		})
	go evsw.FireEvent("event", "data")
	received := <-messages
	if received != "data" {
		t.Errorf("Message received does not match: %v", received)
	}
}

func TestAddListenerForEventFireMany(t *testing.T) {
	evsw := NewEventSwitch()
	err := evsw.Start()
	if err != nil {
		t.Errorf("Failed to start EventSwitch, error: %v", err)
	}
	doneSum := make(chan uint64)
	doneSending := make(chan uint64)
	numbers := make(chan uint64, 4)

	evsw.AddListenerForEvent("listener", "event",
		func(data EventData) {
			numbers <- data.(uint64)
		})

	go sumReceivedNumbers(numbers, doneSum)

	go fireEvents(evsw, "event", doneSending, uint64(1))
	checkSum := <-doneSending
	close(numbers)
	eventSum := <-doneSum
	if checkSum != eventSum {
		t.Errorf("Not all messages sent were received.\n")
	}
}

func TestAddListenerForDifferentEvents(t *testing.T) {
	evsw := NewEventSwitch()
	err := evsw.Start()
	if err != nil {
		t.Errorf("Failed to start EventSwitch, error: %v", err)
	}
	doneSum := make(chan uint64)
	doneSending1 := make(chan uint64)
	doneSending2 := make(chan uint64)
	doneSending3 := make(chan uint64)
	numbers := make(chan uint64, 4)

	evsw.AddListenerForEvent("listener", "event1",
		func(data EventData) {
			numbers <- data.(uint64)
		})
	evsw.AddListenerForEvent("listener", "event2",
		func(data EventData) {
			numbers <- data.(uint64)
		})
	evsw.AddListenerForEvent("listener", "event3",
		func(data EventData) {
			numbers <- data.(uint64)
		})

	go sumReceivedNumbers(numbers, doneSum)

	go fireEvents(evsw, "event1", doneSending1, uint64(1))
	go fireEvents(evsw, "event2", doneSending2, uint64(1))
	go fireEvents(evsw, "event3", doneSending3, uint64(1))
	var checkSum uint64 = 0
	checkSum += <-doneSending1
	checkSum += <-doneSending2
	checkSum += <-doneSending3
	close(numbers)
	eventSum := <-doneSum
	if checkSum != eventSum {
		t.Errorf("Not all messages sent were received.\n")
	}
}

func TestAddDifferentListenerForDifferentEvents(t *testing.T) {
	evsw := NewEventSwitch()
	err := evsw.Start()
	if err != nil {
		t.Errorf("Failed to start EventSwitch, error: %v", err)
	}
	doneSum1 := make(chan uint64)
	doneSum2 := make(chan uint64)
	doneSending1 := make(chan uint64)
	doneSending2 := make(chan uint64)
	doneSending3 := make(chan uint64)
	numbers1 := make(chan uint64, 4)
	numbers2 := make(chan uint64, 4)

	evsw.AddListenerForEvent("listener1", "event1",
		func(data EventData) {
			numbers1 <- data.(uint64)
		})
	evsw.AddListenerForEvent("listener1", "event2",
		func(data EventData) {
			numbers1 <- data.(uint64)
		})
	evsw.AddListenerForEvent("listener1", "event3",
		func(data EventData) {
			numbers1 <- data.(uint64)
		})
	evsw.AddListenerForEvent("listener2", "event2",
		func(data EventData) {
			numbers2 <- data.(uint64)
		})
	evsw.AddListenerForEvent("listener2", "event3",
		func(data EventData) {
			numbers2 <- data.(uint64)
		})

	go sumReceivedNumbers(numbers1, doneSum1)

	go sumReceivedNumbers(numbers2, doneSum2)

	go fireEvents(evsw, "event1", doneSending1, uint64(1))
	go fireEvents(evsw, "event2", doneSending2, uint64(1001))
	go fireEvents(evsw, "event3", doneSending3, uint64(2001))
	checkSumEvent1 := <-doneSending1
	checkSumEvent2 := <-doneSending2
	checkSumEvent3 := <-doneSending3
	checkSum1 := checkSumEvent1 + checkSumEvent2 + checkSumEvent3
	checkSum2 := checkSumEvent2 + checkSumEvent3
	close(numbers1)
	close(numbers2)
	eventSum1 := <-doneSum1
	eventSum2 := <-doneSum2
	if checkSum1 != eventSum1 ||
		checkSum2 != eventSum2 {
		t.Errorf("Not all messages sent were received for different listeners to different events.\n")
	}
}

func TestAddAndRemoveListener(t *testing.T) {
	evsw := NewEventSwitch()
	err := evsw.Start()
	if err != nil {
		t.Errorf("Failed to start EventSwitch, error: %v", err)
	}
	doneSum1 := make(chan uint64)
	doneSum2 := make(chan uint64)
	doneSending1 := make(chan uint64)
	doneSending2 := make(chan uint64)
	numbers1 := make(chan uint64, 4)
	numbers2 := make(chan uint64, 4)

	evsw.AddListenerForEvent("listener", "event1",
		func(data EventData) {
			numbers1 <- data.(uint64)
		})
	evsw.AddListenerForEvent("listener", "event2",
		func(data EventData) {
			numbers2 <- data.(uint64)
		})

	go sumReceivedNumbers(numbers1, doneSum1)

	go sumReceivedNumbers(numbers2, doneSum2)

	go fireEvents(evsw, "event1", doneSending1, uint64(1))
	checkSumEvent1 := <-doneSending1

	evsw.RemoveListener("listener")
	go fireEvents(evsw, "event2", doneSending2, uint64(1001))
	checkSumEvent2 := <-doneSending2
	close(numbers1)
	close(numbers2)
	eventSum1 := <-doneSum1
	eventSum2 := <-doneSum2
	if checkSumEvent1 != eventSum1 ||

		checkSumEvent2 == uint64(0) ||
		eventSum2 != uint64(0) {
		t.Errorf("Not all messages sent were received or unsubscription did not register.\n")
	}
}

func TestRemoveListener(t *testing.T) {
	evsw := NewEventSwitch()
	err := evsw.Start()
	if err != nil {
		t.Errorf("Failed to start EventSwitch, error: %v", err)
	}
	count := 10
	sum1, sum2 := 0, 0

	evsw.AddListenerForEvent("listener", "event1",
		func(data EventData) {
			sum1++
		})
	evsw.AddListenerForEvent("listener", "event2",
		func(data EventData) {
			sum2++
		})
	for i := 0; i < count; i++ {
		evsw.FireEvent("event1", true)
		evsw.FireEvent("event2", true)
	}
	assert.Equal(t, count, sum1)
	assert.Equal(t, count, sum2)

	evsw.RemoveListenerForEvent("event2", "listener")
	for i := 0; i < count; i++ {
		evsw.FireEvent("event1", true)
		evsw.FireEvent("event2", true)
	}
	assert.Equal(t, count*2, sum1)
	assert.Equal(t, count, sum2)

	evsw.RemoveListener("listener")
	for i := 0; i < count; i++ {
		evsw.FireEvent("event1", true)
		evsw.FireEvent("event2", true)
	}
	assert.Equal(t, count*2, sum1)
	assert.Equal(t, count, sum2)
}

func TestRemoveListenersAsync(t *testing.T) {
	evsw := NewEventSwitch()
	err := evsw.Start()
	if err != nil {
		t.Errorf("Failed to start EventSwitch, error: %v", err)
	}
	doneSum1 := make(chan uint64)
	doneSum2 := make(chan uint64)
	doneSending1 := make(chan uint64)
	doneSending2 := make(chan uint64)
	doneSending3 := make(chan uint64)
	numbers1 := make(chan uint64, 4)
	numbers2 := make(chan uint64, 4)

	evsw.AddListenerForEvent("listener1", "event1",
		func(data EventData) {
			numbers1 <- data.(uint64)
		})
	evsw.AddListenerForEvent("listener1", "event2",
		func(data EventData) {
			numbers1 <- data.(uint64)
		})
	evsw.AddListenerForEvent("listener1", "event3",
		func(data EventData) {
			numbers1 <- data.(uint64)
		})
	evsw.AddListenerForEvent("listener2", "event1",
		func(data EventData) {
			numbers2 <- data.(uint64)
		})
	evsw.AddListenerForEvent("listener2", "event2",
		func(data EventData) {
			numbers2 <- data.(uint64)
		})
	evsw.AddListenerForEvent("listener2", "event3",
		func(data EventData) {
			numbers2 <- data.(uint64)
		})

	go sumReceivedNumbers(numbers1, doneSum1)

	go sumReceivedNumbers(numbers2, doneSum2)
	addListenersStress := func() {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		for k := uint16(0); k < 400; k++ {
			listenerNumber := r1.Intn(100) + 3
			eventNumber := r1.Intn(3) + 1
			go evsw.AddListenerForEvent(fmt.Sprintf("listener%v", listenerNumber),
				fmt.Sprintf("event%v", eventNumber),
				func(_ EventData) {})
		}
	}
	removeListenersStress := func() {
		s2 := rand.NewSource(time.Now().UnixNano())
		r2 := rand.New(s2)
		for k := uint16(0); k < 80; k++ {
			listenerNumber := r2.Intn(100) + 3
			go evsw.RemoveListener(fmt.Sprintf("listener%v", listenerNumber))
		}
	}
	addListenersStress()

	go fireEvents(evsw, "event1", doneSending1, uint64(1))
	removeListenersStress()
	go fireEvents(evsw, "event2", doneSending2, uint64(1001))
	go fireEvents(evsw, "event3", doneSending3, uint64(2001))
	checkSumEvent1 := <-doneSending1
	checkSumEvent2 := <-doneSending2
	checkSumEvent3 := <-doneSending3
	checkSum := checkSumEvent1 + checkSumEvent2 + checkSumEvent3
	close(numbers1)
	close(numbers2)
	eventSum1 := <-doneSum1
	eventSum2 := <-doneSum2
	if checkSum != eventSum1 ||
		checkSum != eventSum2 {
		t.Errorf("Not all messages sent were received.\n")
	}
}

func sumReceivedNumbers(numbers, doneSum chan uint64) {
	var sum uint64 = 0
	for {
		j, more := <-numbers
		sum += j
		if !more {
			doneSum <- sum
			close(doneSum)
			return
		}
	}
}

func fireEvents(evsw EventSwitch, event string, doneChan chan uint64,
	offset uint64) {
	var sentSum uint64 = 0
	for i := offset; i <= offset+uint64(999); i++ {
		sentSum += i
		evsw.FireEvent(event, i)
	}
	doneChan <- sentSum
	close(doneChan)
}
