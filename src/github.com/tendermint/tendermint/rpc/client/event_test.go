package client_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

var waitForEventTimeout = 5 * time.Second

func MakeTxKV() ([]byte, []byte, []byte) {
	k := []byte(cmn.RandStr(8))
	v := []byte(cmn.RandStr(8))
	return k, v, append(k, append([]byte("="), v...)...)
}

func TestHeaderEvents(t *testing.T) {
	for i, c := range GetClients() {
		i, c := i, c
		t.Run(reflect.TypeOf(c).String(), func(t *testing.T) {

			if !c.IsRunning() {

				err := c.Start()
				require.Nil(t, err, "%d: %+v", i, err)
				defer c.Stop()
			}

			evtTyp := types.EventNewBlockHeader
			evt, err := client.WaitForOneEvent(c, evtTyp, waitForEventTimeout)
			require.Nil(t, err, "%d: %+v", i, err)
			_, ok := evt.(types.EventDataNewBlockHeader)
			require.True(t, ok, "%d: %#v", i, evt)

		})
	}
}

func TestBlockEvents(t *testing.T) {
	for i, c := range GetClients() {
		i, c := i, c
		t.Run(reflect.TypeOf(c).String(), func(t *testing.T) {

			if !c.IsRunning() {

				err := c.Start()
				require.Nil(t, err, "%d: %+v", i, err)
				defer c.Stop()
			}

			var firstBlockHeight int64
			for j := 0; j < 3; j++ {
				evtTyp := types.EventNewBlock
				evt, err := client.WaitForOneEvent(c, evtTyp, waitForEventTimeout)
				require.Nil(t, err, "%d: %+v", j, err)
				blockEvent, ok := evt.(types.EventDataNewBlock)
				require.True(t, ok, "%d: %#v", j, evt)

				block := blockEvent.Block
				if j == 0 {
					firstBlockHeight = block.Header.Height
					continue
				}

				require.Equal(t, block.Header.Height, firstBlockHeight+int64(j))
			}
		})
	}
}

func TestTxEventsSentWithBroadcastTxAsync(t *testing.T) {
	for i, c := range GetClients() {
		i, c := i, c
		t.Run(reflect.TypeOf(c).String(), func(t *testing.T) {

			if !c.IsRunning() {

				err := c.Start()
				require.Nil(t, err, "%d: %+v", i, err)
				defer c.Stop()
			}

			_, _, tx := MakeTxKV()
			evtTyp := types.EventTx

			txres, err := c.BroadcastTxAsync(tx)
			require.Nil(t, err, "%+v", err)
			require.Equal(t, txres.Code, abci.CodeTypeOK)

			evt, err := client.WaitForOneEvent(c, evtTyp, waitForEventTimeout)
			require.Nil(t, err, "%d: %+v", i, err)

			txe, ok := evt.(types.EventDataTx)
			require.True(t, ok, "%d: %#v", i, evt)

			require.EqualValues(t, tx, txe.Tx)
			require.True(t, txe.Result.IsOK())
		})
	}
}

func TestTxEventsSentWithBroadcastTxSync(t *testing.T) {
	for i, c := range GetClients() {
		i, c := i, c
		t.Run(reflect.TypeOf(c).String(), func(t *testing.T) {

			if !c.IsRunning() {

				err := c.Start()
				require.Nil(t, err, "%d: %+v", i, err)
				defer c.Stop()
			}

			_, _, tx := MakeTxKV()
			evtTyp := types.EventTx

			txres, err := c.BroadcastTxSync(tx)
			require.Nil(t, err, "%+v", err)
			require.Equal(t, txres.Code, abci.CodeTypeOK)

			evt, err := client.WaitForOneEvent(c, evtTyp, waitForEventTimeout)
			require.Nil(t, err, "%d: %+v", i, err)

			txe, ok := evt.(types.EventDataTx)
			require.True(t, ok, "%d: %#v", i, evt)

			require.EqualValues(t, tx, txe.Tx)
			require.True(t, txe.Result.IsOK())
		})
	}
}
