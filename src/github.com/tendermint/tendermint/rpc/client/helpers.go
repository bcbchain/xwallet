package client

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/types"
)

type Waiter func(delta int64) (abort error)

func DefaultWaitStrategy(delta int64) (abort error) {
	if delta > 10 {
		return errors.Errorf("Waiting for %d blocks... aborting", delta)
	} else if delta > 0 {

		delay := time.Duration(delta-1)*time.Second + 500*time.Millisecond
		time.Sleep(delay)
	}
	return nil
}

func WaitForHeight(c StatusClient, h int64, waiter Waiter) error {
	if waiter == nil {
		waiter = DefaultWaitStrategy
	}
	delta := int64(1)
	for delta > 0 {
		s, err := c.Status()
		if err != nil {
			return err
		}
		delta = h - s.SyncInfo.LatestBlockHeight

		if err := waiter(delta); err != nil {
			return err
		}
	}
	return nil
}

func WaitForOneEvent(c EventsClient, evtTyp string, timeout time.Duration) (types.TMEventData, error) {
	const subscriber = "helpers"
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	evts := make(chan interface{}, 1)

	query := types.QueryForEvent(evtTyp)
	err := c.Subscribe(ctx, subscriber, query, evts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe")
	}

	defer c.UnsubscribeAll(ctx, subscriber)

	select {
	case evt := <-evts:
		return evt.(types.TMEventData), nil
	case <-ctx.Done():
		return nil, errors.New("timed out waiting for event")
	}
}
