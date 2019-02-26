package client_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/abci/types"

	"github.com/tendermint/tendermint/rpc/client"
	rpctest "github.com/tendermint/tendermint/rpc/test"
	"github.com/tendermint/tendermint/types"
)

func getHTTPClient() *client.HTTP {
	rpcAddr := rpctest.GetConfig().RPC.ListenAddress
	return client.NewHTTP(rpcAddr, "/websocket")
}

func getLocalClient() *client.Local {
	return client.NewLocal(node)
}

func GetClients() []client.Client {
	return []client.Client{
		getHTTPClient(),
		getLocalClient(),
	}
}

func TestStatus(t *testing.T) {
	for i, c := range GetClients() {
		moniker := rpctest.GetConfig().Moniker
		status, err := c.Status()
		require.Nil(t, err, "%d: %+v", i, err)
		assert.Equal(t, moniker, status.NodeInfo.Moniker)
	}
}

func TestInfo(t *testing.T) {
	for i, c := range GetClients() {

		info, err := c.ABCIInfo()
		require.Nil(t, err, "%d: %+v", i, err)

		assert.True(t, strings.Contains(info.Response.Data, "size"))
	}
}

func TestNetInfo(t *testing.T) {
	for i, c := range GetClients() {
		nc, ok := c.(client.NetworkClient)
		require.True(t, ok, "%d", i)
		netinfo, err := nc.NetInfo()
		require.Nil(t, err, "%d: %+v", i, err)
		assert.True(t, netinfo.Listening)
		assert.Equal(t, 0, len(netinfo.Peers))
	}
}

func TestDumpConsensusState(t *testing.T) {
	for i, c := range GetClients() {

		nc, ok := c.(client.NetworkClient)
		require.True(t, ok, "%d", i)
		cons, err := nc.DumpConsensusState()
		require.Nil(t, err, "%d: %+v", i, err)
		assert.NotEmpty(t, cons.RoundState)
		assert.Empty(t, cons.PeerRoundStates)
	}
}

func TestHealth(t *testing.T) {
	for i, c := range GetClients() {
		nc, ok := c.(client.NetworkClient)
		require.True(t, ok, "%d", i)
		_, err := nc.Health()
		require.Nil(t, err, "%d: %+v", i, err)
	}
}

func TestGenesisAndValidators(t *testing.T) {
	for i, c := range GetClients() {

		gen, err := c.Genesis()
		require.Nil(t, err, "%d: %+v", i, err)

		require.Equal(t, 1, len(gen.Genesis.Validators))
		gval := gen.Genesis.Validators[0]

		vals, err := c.Validators(nil)
		require.Nil(t, err, "%d: %+v", i, err)
		require.Equal(t, 1, len(vals.Validators))
		val := vals.Validators[0]

		assert.Equal(t, gval.Power, val.VotingPower)
		assert.Equal(t, gval.PubKey, val.PubKey)
	}
}

func TestABCIQuery(t *testing.T) {
	for i, c := range GetClients() {

		k, v, tx := MakeTxKV()
		bres, err := c.BroadcastTxCommit(tx)
		require.Nil(t, err, "%d: %+v", i, err)
		apph := bres.Height + 1

		client.WaitForHeight(c, apph, nil)
		res, err := c.ABCIQuery("/key", k)
		qres := res.Response
		if assert.Nil(t, err) && assert.True(t, qres.IsOK()) {
			assert.EqualValues(t, v, qres.Value)
		}
	}
}

func TestAppCalls(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	for i, c := range GetClients() {

		s, err := c.Status()
		require.Nil(err, "%d: %+v", i, err)

		sh := s.SyncInfo.LatestBlockHeight

		h := sh + 2
		_, err = c.Block(&h)
		assert.NotNil(err)

		k, v, tx := MakeTxKV()
		bres, err := c.BroadcastTxCommit(tx)
		require.Nil(err, "%d: %+v", i, err)
		require.True(bres.DeliverTx.IsOK())
		txh := bres.Height
		apph := txh + 1

		if err := client.WaitForHeight(c, apph, nil); err != nil {
			t.Error(err)
		}
		_qres, err := c.ABCIQueryWithOptions("/key", k, client.ABCIQueryOptions{Trusted: true})
		qres := _qres.Response
		if assert.Nil(err) && assert.True(qres.IsOK()) {

			assert.EqualValues(v, qres.Value)
		}

		ptx, err := c.Tx(bres.Hash, true)
		require.Nil(err, "%d: %+v", i, err)
		assert.EqualValues(txh, ptx.Height)
		assert.EqualValues(tx, ptx.Tx)

		block, err := c.Block(&apph)
		require.Nil(err, "%d: %+v", i, err)
		appHash := block.BlockMeta.Header.LastAppHash
		assert.True(len(appHash) > 0)
		assert.EqualValues(apph, block.BlockMeta.Header.Height)

		blockResults, err := c.BlockResults(&txh)
		require.Nil(err, "%d: %+v", i, err)
		assert.Equal(txh, blockResults.Height)
		if assert.Equal(1, len(blockResults.Results.DeliverTx)) {

			assert.EqualValues(0, blockResults.Results.DeliverTx[0].Code)
		}

		info, err := c.BlockchainInfo(apph, apph)
		require.Nil(err, "%d: %+v", i, err)
		assert.True(info.LastHeight >= apph)
		if assert.Equal(1, len(info.BlockMetas)) {
			lastMeta := info.BlockMetas[0]
			assert.EqualValues(apph, lastMeta.Header.Height)
			bMeta := block.BlockMeta
			assert.Equal(bMeta.Header.LastAppHash, lastMeta.Header.LastAppHash)
			assert.Equal(bMeta.BlockID, lastMeta.BlockID)
		}

		commit, err := c.Commit(&apph)
		require.Nil(err, "%d: %+v", i, err)
		cappHash := commit.Header.LastAppHash
		assert.Equal(appHash, cappHash)
		assert.NotNil(commit.Commit)

		h = apph - 1
		commit2, err := c.Commit(&h)
		require.Nil(err, "%d: %+v", i, err)
		assert.Equal(block.Block.LastCommit, commit2.Commit)

		_pres, err := c.ABCIQueryWithOptions("/key", k, client.ABCIQueryOptions{Trusted: false})
		pres := _pres.Response
		assert.Nil(err)
		assert.True(pres.IsOK())
	}
}

func TestBroadcastTxSync(t *testing.T) {
	require := require.New(t)

	mempool := node.MempoolReactor().Mempool
	initMempoolSize := mempool.Size()

	for i, c := range GetClients() {
		_, _, tx := MakeTxKV()
		bres, err := c.BroadcastTxSync(tx)
		require.Nil(err, "%d: %+v", i, err)
		require.Equal(bres.Code, abci.CodeTypeOK)

		require.Equal(initMempoolSize+1, mempool.Size())

		txs := mempool.Reap(1)
		require.EqualValues(tx, txs[0])
		mempool.Flush()
	}
}

func TestBroadcastTxCommit(t *testing.T) {
	require := require.New(t)

	mempool := node.MempoolReactor().Mempool
	for i, c := range GetClients() {
		_, _, tx := MakeTxKV()
		bres, err := c.BroadcastTxCommit(tx)
		require.Nil(err, "%d: %+v", i, err)
		require.True(bres.CheckTx.IsOK())
		require.True(bres.DeliverTx.IsOK())

		require.Equal(0, mempool.Size())
	}
}

func TestTx(t *testing.T) {

	c := getHTTPClient()
	_, _, tx := MakeTxKV()
	bres, err := c.BroadcastTxCommit(tx)
	require.Nil(t, err, "%+v", err)

	txHeight := bres.Height
	txHash := bres.Hash

	anotherTxHash := types.Tx("a different tx").Hash()

	cases := []struct {
		valid	bool
		hash	[]byte
		prove	bool
	}{

		{true, txHash, false},
		{true, txHash, true},
		{false, anotherTxHash, false},
		{false, anotherTxHash, true},
		{false, nil, false},
		{false, nil, true},
	}

	for i, c := range GetClients() {
		for j, tc := range cases {
			t.Logf("client %d, case %d", i, j)

			ptx, err := c.Tx(tc.hash, tc.prove)

			if !tc.valid {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err, "%+v", err)
				assert.EqualValues(t, txHeight, ptx.Height)
				assert.EqualValues(t, tx, ptx.Tx)
				assert.Zero(t, ptx.Index)
				assert.True(t, ptx.TxResult.IsOK())
				assert.EqualValues(t, txHash, ptx.Hash)

				proof := ptx.Proof
				if tc.prove && assert.EqualValues(t, tx, proof.Data) {
					assert.True(t, proof.Proof.Verify(proof.Index, proof.Total, txHash, proof.RootHash))
				}
			}
		}
	}
}

func TestTxSearch(t *testing.T) {

	c := getHTTPClient()
	_, _, tx := MakeTxKV()
	bres, err := c.BroadcastTxCommit(tx)
	require.Nil(t, err, "%+v", err)

	txHeight := bres.Height
	txHash := bres.Hash

	anotherTxHash := types.Tx("a different tx").Hash()

	for i, c := range GetClients() {
		t.Logf("client %d", i)

		results, err := c.TxSearch(fmt.Sprintf("tx.hash='%v'", txHash), true)
		require.Nil(t, err, "%+v", err)
		require.Len(t, results, 1)

		ptx := results[0]
		assert.EqualValues(t, txHeight, ptx.Height)
		assert.EqualValues(t, tx, ptx.Tx)
		assert.Zero(t, ptx.Index)
		assert.True(t, ptx.TxResult.IsOK())
		assert.EqualValues(t, txHash, ptx.Hash)

		proof := ptx.Proof
		if assert.EqualValues(t, tx, proof.Data) {
			assert.True(t, proof.Proof.Verify(proof.Index, proof.Total, txHash, proof.RootHash))
		}

		results, err = c.TxSearch(fmt.Sprintf("tx.hash='%X'", anotherTxHash), false)
		require.Nil(t, err, "%+v", err)
		require.Len(t, results, 0)

		results, err = c.TxSearch("app.creator='jae'", false)
		require.Nil(t, err, "%+v", err)
		if len(results) == 0 {
			t.Fatal("expected a lot of transactions")
		}
	}
}
