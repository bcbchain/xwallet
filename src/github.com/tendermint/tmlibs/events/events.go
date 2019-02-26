package events

import (
	"sync"

	cmn "github.com/tendermint/tmlibs/common"
)

type EventData interface {
}

type Eventable interface {
	SetEventSwitch(evsw EventSwitch)
}

type Fireable interface {
	FireEvent(event string, data EventData)
}

type EventSwitch interface {
	cmn.Service
	Fireable

	AddListenerForEvent(listenerID, event string, cb EventCallback)
	RemoveListenerForEvent(event string, listenerID string)
	RemoveListener(listenerID string)
}

type eventSwitch struct {
	cmn.BaseService

	mtx		sync.RWMutex
	eventCells	map[string]*eventCell
	listeners	map[string]*eventListener
}

func NewEventSwitch() EventSwitch {
	evsw := &eventSwitch{}
	evsw.BaseService = *cmn.NewBaseService(nil, "EventSwitch", evsw)
	return evsw
}

func (evsw *eventSwitch) OnStart() error {
	evsw.BaseService.OnStart()
	evsw.eventCells = make(map[string]*eventCell)
	evsw.listeners = make(map[string]*eventListener)
	return nil
}

func (evsw *eventSwitch) OnStop() {
	evsw.mtx.Lock()
	defer evsw.mtx.Unlock()
	evsw.BaseService.OnStop()
	evsw.eventCells = nil
	evsw.listeners = nil
}

func (evsw *eventSwitch) AddListenerForEvent(listenerID, event string, cb EventCallback) {

	evsw.mtx.Lock()
	eventCell := evsw.eventCells[event]
	if eventCell == nil {
		eventCell = newEventCell()
		evsw.eventCells[event] = eventCell
	}
	listener := evsw.listeners[listenerID]
	if listener == nil {
		listener = newEventListener(listenerID)
		evsw.listeners[listenerID] = listener
	}
	evsw.mtx.Unlock()

	eventCell.AddListener(listenerID, cb)
	listener.AddEvent(event)
}

func (evsw *eventSwitch) RemoveListener(listenerID string) {

	evsw.mtx.RLock()
	listener := evsw.listeners[listenerID]
	evsw.mtx.RUnlock()
	if listener == nil {
		return
	}

	evsw.mtx.Lock()
	delete(evsw.listeners, listenerID)
	evsw.mtx.Unlock()

	listener.SetRemoved()
	for _, event := range listener.GetEvents() {
		evsw.RemoveListenerForEvent(event, listenerID)
	}
}

func (evsw *eventSwitch) RemoveListenerForEvent(event string, listenerID string) {

	evsw.mtx.Lock()
	eventCell := evsw.eventCells[event]
	evsw.mtx.Unlock()

	if eventCell == nil {
		return
	}

	numListeners := eventCell.RemoveListener(listenerID)

	if numListeners == 0 {

		evsw.mtx.Lock()
		eventCell.mtx.Lock()
		if len(eventCell.listeners) == 0 {
			delete(evsw.eventCells, event)
		}
		eventCell.mtx.Unlock()
		evsw.mtx.Unlock()
	}
}

func (evsw *eventSwitch) FireEvent(event string, data EventData) {

	evsw.mtx.RLock()
	eventCell := evsw.eventCells[event]
	evsw.mtx.RUnlock()

	if eventCell == nil {
		return
	}

	eventCell.FireEvent(data)
}

type eventCell struct {
	mtx		sync.RWMutex
	listeners	map[string]EventCallback
}

func newEventCell() *eventCell {
	return &eventCell{
		listeners: make(map[string]EventCallback),
	}
}

func (cell *eventCell) AddListener(listenerID string, cb EventCallback) {
	cell.mtx.Lock()
	cell.listeners[listenerID] = cb
	cell.mtx.Unlock()
}

func (cell *eventCell) RemoveListener(listenerID string) int {
	cell.mtx.Lock()
	delete(cell.listeners, listenerID)
	numListeners := len(cell.listeners)
	cell.mtx.Unlock()
	return numListeners
}

func (cell *eventCell) FireEvent(data EventData) {
	cell.mtx.RLock()
	for _, listener := range cell.listeners {
		listener(data)
	}
	cell.mtx.RUnlock()
}

type EventCallback func(data EventData)

type eventListener struct {
	id	string

	mtx	sync.RWMutex
	removed	bool
	events	[]string
}

func newEventListener(id string) *eventListener {
	return &eventListener{
		id:		id,
		removed:	false,
		events:		nil,
	}
}

func (evl *eventListener) AddEvent(event string) {
	evl.mtx.Lock()
	defer evl.mtx.Unlock()

	if evl.removed {
		return
	}
	evl.events = append(evl.events, event)
}

func (evl *eventListener) GetEvents() []string {
	evl.mtx.RLock()
	defer evl.mtx.RUnlock()

	events := make([]string, len(evl.events))
	copy(events, evl.events)
	return events
}

func (evl *eventListener) SetRemoved() {
	evl.mtx.Lock()
	defer evl.mtx.Unlock()
	evl.removed = true
}
