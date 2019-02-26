package rpctest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tendermint/tmlibs/log"

	abci "github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"

	cfg "github.com/tendermint/tendermint/config"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	core_grpc "github.com/tendermint/tendermint/rpc/grpc"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
	pvm "github.com/tendermint/tendermint/types/priv_validator"
)

var globalConfig *cfg.Config

func waitForRPC() {
	laddr := GetConfig().RPC.ListenAddress
	client := rpcclient.NewJSONRPCClient(laddr)
	ctypes.RegisterAmino(client.Codec())
	result := new(ctypes.ResultStatus)
	for {
		_, err := client.Call("status", map[string]interface{}{}, result)
		if err == nil {
			return
		} else {
			fmt.Println("error", err)
			time.Sleep(time.Millisecond)
		}
	}
}

func waitForGRPC() {
	client := GetGRPCClient()
	for {
		_, err := client.Ping(context.Background(), &core_grpc.RequestPing{})
		if err == nil {
			return
		}
	}
}

func makePathname() string {

	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	sep := string(filepath.Separator)
	return strings.Replace(p, sep, "_", -1)
}

func randPort() int {
	return int(cmn.RandUint16()/2 + 10000)
}

func makeAddrs() (string, string, string) {
	start := randPort()
	return fmt.Sprintf("tcp://0.0.0.0:%d", start),
		fmt.Sprintf("tcp://0.0.0.0:%d", start+1),
		fmt.Sprintf("tcp://0.0.0.0:%d", start+2)
}

func GetConfig() *cfg.Config {
	if globalConfig == nil {
		pathname := makePathname()
		globalConfig = cfg.ResetTestRoot(pathname)

		tm, rpc, grpc := makeAddrs()
		globalConfig.P2P.ListenAddress = tm
		globalConfig.RPC.ListenAddress = rpc
		globalConfig.RPC.GRPCListenAddress = grpc
		globalConfig.TxIndex.IndexTags = "app.creator"
	}
	return globalConfig
}

func GetGRPCClient() core_grpc.BroadcastAPIClient {
	grpcAddr := globalConfig.RPC.GRPCListenAddress
	return core_grpc.StartGRPCClient(grpcAddr)
}

func StartTendermint(app abci.Application) *nm.Node {
	node := NewTendermint(app)
	err := node.Start()
	if err != nil {
		panic(err)
	}

	waitForRPC()
	waitForGRPC()

	fmt.Println("Tendermint running!")

	return node
}

func NewTendermint(app abci.Application) *nm.Node {

	config := GetConfig()
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger = log.NewFilter(logger, log.AllowError())
	pvFile := config.PrivValidatorFile()
	pv := pvm.LoadOrGenFilePV(pvFile)
	papp := proxy.NewLocalClientCreator(app)
	node, err := nm.NewNode(config, pv, papp,
		nm.DefaultGenesisDocProviderFunc(config),
		nm.DefaultDBProvider, logger)
	if err != nil {
		panic(err)
	}
	return node
}
