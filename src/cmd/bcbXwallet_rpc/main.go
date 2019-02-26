package main

import (
	"bcbXwallet/common"
	"bcbXwallet/rpc"
	"bcbchain.io/rpc/lib/server"
	"github.com/tendermint/go-amino"
	cmn "github.com/tendermint/tmlibs/common"
	"net/http"
	"os"
)

func main() {
	err := common.InitAll()
	if err != nil {
		panic(err)
	}

	err = rpc.InitDB()
	if err != nil {
		common.GetLogger().Error("open db failed", "error", err.Error())
		panic(err)
	}

	rpcLogger := common.GetLogger()

	coreCodec := amino.NewCodec()

	mux := http.NewServeMux()

	rpcserver.RegisterRPCFuncs(mux, rpc.Routes, coreCodec, rpcLogger)

	if common.GetConfig().UseHttps {
		crtPath, keyPath := common.OutCertFileIsExist()
		_, err = rpcserver.StartHTTPAndTLSServer(common.GetConfig().ServerAddr, mux, crtPath, keyPath, rpcLogger)
		if err != nil {
			cmn.Exit(err.Error())
		}
	} else {
		_, err = rpcserver.StartHTTPServer(common.GetConfig().ServerAddr, mux, rpcLogger)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	cmn.TrapSignal(func(signal os.Signal) {
	})

}
