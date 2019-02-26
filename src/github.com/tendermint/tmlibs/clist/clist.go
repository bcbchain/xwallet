package clist

import (
	"sync"
)

type CElement struct {
	mtx		sync.RWMutex
	prev		*CElement
	prevWg		*sync.WaitGroup
	prevWaitCh	chan struct{}
	next		*CElement
	nextWg		*sync.WaitGroup
	nextWaitCh	chan struct{}
	removed		bool

	Value	interface{}
}

func (e *CElement) NextWait() *CElement {
	for {
		e.mtx.RLock()
		next := e.next
		nextWg := e.nextWg
		removed := e.removed
		e.mtx.RUnlock()

		if next != nil || removed {
			return next
		}

		nextWg.Wait()

	}
}

func (e *CElement) PrevWait() *CElement {
	for {
		e.mtx.RLock()
		prev := e.prev
		prevWg := e.prevWg
		removed := e.removed
		e.mtx.RUnlock()

		if prev != nil || removed {
			return prev
		}

		prevWg.Wait()
	}
}

func (e *CElement) PrevWaitChan() <-chan struct{} {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	return e.prevWaitCh
}

func (e *CElement) NextWaitChan() <-chan struct{} {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	return e.nextWaitCh
}

func (e *CElement) Next() *CElement {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	return e.next
}

func (e *CElement) Prev() *CElement {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	return e.prev
}

func (e *CElement) Removed() bool {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	return e.removed
}

func (e *CElement) DetachNext() {
	if !e.Removed() {
		panic("DetachNext() must be called after Remove(e)")
	}
	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.next = nil
}

func (e *CElement) DetachPrev() {
	if !e.Removed() {
		panic("DetachPrev() must be called after Remove(e)")
	}
	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.prev = nil
}

func (e *CElement) SetNext(newNext *CElement) {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	oldNext := e.next
	e.next = newNext
	if oldNext != nil && newNext == nil {

		e.nextWg = waitGroup1()
		e.nextWaitCh = make(chan struct{})
	}
	if oldNext == nil && newNext != nil {
		e.nextWg.Done()
		close(e.nextWaitCh)
	}
}

func (e *CElement) SetPrev(newPrev *CElement) {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	oldPrev := e.prev
	e.prev = newPrev
	if oldPrev != nil && newPrev == nil {
		e.prevWg = waitGroup1()
		e.prevWaitCh = make(chan struct{})
	}
	if oldPrev == nil && newPrev != nil {
		e.prevWg.Done()
		close(e.prevWaitCh)
	}
}

func (e *CElement) SetRemoved() {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.removed = true

	if e.prev == nil {
		e.prevWg.Done()
		close(e.prevWaitCh)
	}
	if e.next == nil {
		e.nextWg.Done()
		close(e.nextWaitCh)
	}
}

type CList struct {
	mtx	sync.RWMutex
	wg	*sync.WaitGroup
	waitCh	chan struct{}
	head	*CElement
	tail	*CElement
	len	int
}

func (l *CList) Init() *CList {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	l.wg = waitGroup1()
	l.waitCh = make(chan struct{})
	l.head = nil
	l.tail = nil
	l.len = 0
	return l
}

func New() *CList	{ return new(CList).Init() }

func (l *CList) Len() int {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	return l.len
}

func (l *CList) Front() *CElement {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	return l.head
}

func (l *CList) FrontWait() *CElement {

	for {
		l.mtx.RLock()
		head := l.head
		wg := l.wg
		l.mtx.RUnlock()

		if head != nil {
			return head
		}
		wg.Wait()

	}
}

func (l *CList) Back() *CElement {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	return l.tail
}

func (l *CList) BackWait() *CElement {
	for {
		l.mtx.RLock()
		tail := l.tail
		wg := l.wg
		l.mtx.RUnlock()

		if tail != nil {
			return tail
		}
		wg.Wait()

	}
}

func (l *CList) WaitChan() <-chan struct{} {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	return l.waitCh
}

func (l *CList) PushBack(v interface{}) *CElement {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	e := &CElement{
		prev:		nil,
		prevWg:		waitGroup1(),
		prevWaitCh:	make(chan struct{}),
		next:		nil,
		nextWg:		waitGroup1(),
		nextWaitCh:	make(chan struct{}),
		removed:	false,
		Value:		v,
	}

	if l.len == 0 {
		l.wg.Done()
		close(l.waitCh)
	}
	l.len++

	if l.tail == nil {
		l.head = e
		l.tail = e
	} else {
		e.SetPrev(l.tail)
		l.tail.SetNext(e)
		l.tail = e
	}

	return e
}

func (l *CList) Remove(e *CElement) interface{} {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	prev := e.Prev()
	next := e.Next()

	if l.head == nil || l.tail == nil {
		panic("Remove(e) on empty CList")
	}
	if prev == nil && l.head != e {
		panic("Remove(e) with false head")
	}
	if next == nil && l.tail != e {
		panic("Remove(e) with false tail")
	}

	if l.len == 1 {
		l.wg = waitGroup1()
		l.waitCh = make(chan struct{})
	}

	l.len--

	if prev == nil {
		l.head = next
	} else {
		prev.SetNext(next)
	}
	if next == nil {
		l.tail = prev
	} else {
		next.SetPrev(prev)
	}

	e.SetRemoved()

	return e.Value
}

func waitGroup1() (wg *sync.WaitGroup) {
	wg = &sync.WaitGroup{}
	wg.Add(1)
	return
}
