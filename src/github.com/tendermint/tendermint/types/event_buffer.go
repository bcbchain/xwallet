package types

var _ TxEventPublisher = (*TxEventBuffer)(nil)

type TxEventBuffer struct {
	next		TxEventPublisher
	capacity	int
	events		[]EventDataTx
}

func NewTxEventBuffer(next TxEventPublisher, capacity int) *TxEventBuffer {
	return &TxEventBuffer{
		next:		next,
		capacity:	capacity,
		events:		make([]EventDataTx, 0, capacity),
	}
}

func (b TxEventBuffer) Len() int {
	return len(b.events)
}

func (b *TxEventBuffer) PublishEventTx(e EventDataTx) error {
	b.events = append(b.events, e)
	return nil
}

func (b *TxEventBuffer) Flush() error {
	for _, e := range b.events {
		err := b.next.PublishEventTx(e)
		if err != nil {
			return err
		}
	}

	b.events = b.events[:0]
	return nil
}
