package common

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParallel(t *testing.T) {

	var counter = new(int32)
	var tasks = make([]Task, 100*1000)
	for i := 0; i < len(tasks); i++ {
		tasks[i] = func(i int) (res interface{}, err error, abort bool) {
			atomic.AddInt32(counter, 1)
			return -1 * i, nil, false
		}
	}

	var trs, ok = Parallel(tasks...)
	assert.True(t, ok)

	assert.Equal(t, int(*counter), len(tasks), "Each task should have incremented the counter already")
	var failedTasks int
	for i := 0; i < len(tasks); i++ {
		taskResult, ok := trs.LatestResult(i)
		if !ok {
			assert.Fail(t, "Task #%v did not complete.", i)
			failedTasks++
		} else if taskResult.Error != nil {
			assert.Fail(t, "Task should not have errored but got %v", taskResult.Error)
			failedTasks++
		} else if !assert.Equal(t, -1*i, taskResult.Value.(int)) {
			assert.Fail(t, "Task should have returned %v but got %v", -1*i, taskResult.Value.(int))
			failedTasks++
		} else {

		}
	}
	assert.Equal(t, failedTasks, 0, "No task should have failed")
	assert.Nil(t, trs.FirstError(), "There should be no errors")
	assert.Equal(t, 0, trs.FirstValue(), "First value should be 0")
}

func TestParallelAbort(t *testing.T) {

	var flow1 = make(chan struct{}, 1)
	var flow2 = make(chan struct{}, 1)
	var flow3 = make(chan struct{}, 1)
	var flow4 = make(chan struct{}, 1)

	var tasks = []Task{
		func(i int) (res interface{}, err error, abort bool) {
			assert.Equal(t, i, 0)
			flow1 <- struct{}{}
			return 0, nil, false
		},
		func(i int) (res interface{}, err error, abort bool) {
			assert.Equal(t, i, 1)
			flow2 <- <-flow1
			return 1, errors.New("some error"), false
		},
		func(i int) (res interface{}, err error, abort bool) {
			assert.Equal(t, i, 2)
			flow3 <- <-flow2
			return 2, nil, true
		},
		func(i int) (res interface{}, err error, abort bool) {
			assert.Equal(t, i, 3)
			<-flow4
			return 3, nil, false
		},
	}

	var taskResultSet, ok = Parallel(tasks...)
	assert.False(t, ok, "ok should be false since we aborted task #2.")

	waitTimeout(t, taskResultSet.chz[3], "Task #3")

	flow4 <- <-flow3

	taskResultSet.Wait()

	checkResult(t, taskResultSet, 0, 0, nil, nil)
	checkResult(t, taskResultSet, 1, 1, errors.New("some error"), nil)
	checkResult(t, taskResultSet, 2, 2, nil, nil)
	checkResult(t, taskResultSet, 3, 3, nil, nil)
}

func TestParallelRecover(t *testing.T) {

	var tasks = []Task{
		func(i int) (res interface{}, err error, abort bool) {
			return 0, nil, false
		},
		func(i int) (res interface{}, err error, abort bool) {
			return 1, errors.New("some error"), false
		},
		func(i int) (res interface{}, err error, abort bool) {
			panic(2)
		},
	}

	var taskResultSet, ok = Parallel(tasks...)
	assert.False(t, ok, "ok should be false since we panic'd in task #2.")

	checkResult(t, taskResultSet, 0, 0, nil, nil)
	checkResult(t, taskResultSet, 1, 1, errors.New("some error"), nil)
	checkResult(t, taskResultSet, 2, nil, nil, 2)
}

func checkResult(t *testing.T, taskResultSet *TaskResultSet, index int, val interface{}, err error, pnk interface{}) {
	taskResult, ok := taskResultSet.LatestResult(index)
	taskName := fmt.Sprintf("Task #%v", index)
	assert.True(t, ok, "TaskResultCh unexpectedly closed for %v", taskName)
	assert.Equal(t, val, taskResult.Value, taskName)
	if err != nil {
		assert.Equal(t, err, taskResult.Error, taskName)
	} else if pnk != nil {
		assert.Equal(t, pnk, taskResult.Error.(Error).Cause(), taskName)
	} else {
		assert.Nil(t, taskResult.Error, taskName)
	}
}

func waitTimeout(t *testing.T, taskResultCh TaskResultCh, taskName string) {
	select {
	case _, ok := <-taskResultCh:
		if !ok {
			assert.Fail(t, "TaskResultCh unexpectedly closed (%v)", taskName)
		} else {
			assert.Fail(t, "TaskResultCh unexpectedly returned for %v", taskName)
		}
	case <-time.After(1 * time.Second):

	}
}
