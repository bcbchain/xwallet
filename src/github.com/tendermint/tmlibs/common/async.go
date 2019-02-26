package common

import (
	"sync/atomic"
)

type Task func(i int) (val interface{}, err error, abort bool)

type TaskResult struct {
	Value	interface{}
	Error	error
}

type TaskResultCh <-chan TaskResult

type taskResultOK struct {
	TaskResult
	OK	bool
}

type TaskResultSet struct {
	chz	[]TaskResultCh
	results	[]taskResultOK
}

func newTaskResultSet(chz []TaskResultCh) *TaskResultSet {
	return &TaskResultSet{
		chz:		chz,
		results:	make([]taskResultOK, len(chz)),
	}
}

func (trs *TaskResultSet) Channels() []TaskResultCh {
	return trs.chz
}

func (trs *TaskResultSet) LatestResult(index int) (TaskResult, bool) {
	if len(trs.results) <= index {
		return TaskResult{}, false
	}
	resultOK := trs.results[index]
	return resultOK.TaskResult, resultOK.OK
}

func (trs *TaskResultSet) Reap() *TaskResultSet {
	for i := 0; i < len(trs.results); i++ {
		var trch = trs.chz[i]
		select {
		case result, ok := <-trch:
			if ok {

				trs.results[i] = taskResultOK{
					TaskResult:	result,
					OK:		true,
				}
			} else {

			}
		default:

		}
	}
	return trs
}

func (trs *TaskResultSet) Wait() *TaskResultSet {
	for i := 0; i < len(trs.results); i++ {
		var trch = trs.chz[i]
		select {
		case result, ok := <-trch:
			if ok {

				trs.results[i] = taskResultOK{
					TaskResult:	result,
					OK:		true,
				}
			} else {

			}
		}
	}
	return trs
}

func (trs *TaskResultSet) FirstValue() interface{} {
	for _, result := range trs.results {
		if result.Value != nil {
			return result.Value
		}
	}
	return nil
}

func (trs *TaskResultSet) FirstError() error {
	for _, result := range trs.results {
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func Parallel(tasks ...Task) (trs *TaskResultSet, ok bool) {
	var taskResultChz = make([]TaskResultCh, len(tasks))
	var taskDoneCh = make(chan bool, len(tasks))
	var numPanics = new(int32)
	ok = true

	for i, task := range tasks {
		var taskResultCh = make(chan TaskResult, 1)
		taskResultChz[i] = taskResultCh
		go func(i int, task Task, taskResultCh chan TaskResult) {

			defer func() {
				if pnk := recover(); pnk != nil {
					atomic.AddInt32(numPanics, 1)

					taskResultCh <- TaskResult{nil, ErrorWrap(pnk, "Panic in task")}

					close(taskResultCh)

					taskDoneCh <- false
				}
			}()

			var val, err, abort = task(i)

			taskResultCh <- TaskResult{val, err}

			close(taskResultCh)

			taskDoneCh <- abort
		}(i, task, taskResultCh)
	}

	for i := 0; i < len(tasks); i++ {
		abort := <-taskDoneCh
		if abort {
			ok = false
			break
		}
	}

	ok = ok && (atomic.LoadInt32(numPanics) == 0)

	return newTaskResultSet(taskResultChz).Reap(), ok
}
