package txindex

import (
	"context"

	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

const (
	subscriber = "IndexerService"
)

type IndexerService struct {
	cmn.BaseService

	idr		TxIndexer
	eventBus	*types.EventBus
}

func NewIndexerService(idr TxIndexer, eventBus *types.EventBus) *IndexerService {
	is := &IndexerService{idr: idr, eventBus: eventBus}
	is.BaseService = *cmn.NewBaseService(nil, "IndexerService", is)
	return is
}

func (is *IndexerService) OnStart() error {
	ch := make(chan interface{})
	if err := is.eventBus.Subscribe(context.Background(), subscriber, types.EventQueryTx, ch); err != nil {
		return err
	}
	go func() {
		for event := range ch {

			txResult := event.(types.EventDataTx).TxResult
			is.idr.Index(&txResult)
		}
	}()
	return nil
}

func (is *IndexerService) OnStop() {
	if is.eventBus.IsRunning() {
		_ = is.eventBus.UnsubscribeAll(context.Background(), subscriber)
	}
}
