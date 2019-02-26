package events

type EventCache struct {
	evsw	Fireable
	events	[]eventInfo
}

func NewEventCache(evsw Fireable) *EventCache {
	return &EventCache{
		evsw: evsw,
	}
}

type eventInfo struct {
	event	string
	data	EventData
}

func (evc *EventCache) FireEvent(event string, data EventData) {

	evc.events = append(evc.events, eventInfo{event, data})
}

func (evc *EventCache) Flush() {
	for _, ei := range evc.events {
		evc.evsw.FireEvent(ei.event, ei.data)
	}

	evc.events = nil
}
