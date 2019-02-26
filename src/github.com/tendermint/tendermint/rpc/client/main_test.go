package client_test

import (
	"os"
	"testing"

	"github.com/tendermint/abci/example/kvstore"
	nm "github.com/tendermint/tendermint/node"
	rpctest "github.com/tendermint/tendermint/rpc/test"
)

var node *nm.Node

func TestMain(m *testing.M) {

	app := kvstore.NewKVStoreApplication()
	node = rpctest.StartTendermint(app)
	code := m.Run()

	node.Stop()
	node.Wait()
	os.Exit(code)
}
