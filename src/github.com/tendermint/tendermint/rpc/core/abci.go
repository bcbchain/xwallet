package core

import (
	abci "github.com/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/version"
	cmn "github.com/tendermint/tmlibs/common"
)

func ABCIQuery(path string, data cmn.HexBytes, height int64, trusted bool) (*ctypes.ResultABCIQuery, error) {
	resQuery, err := proxyAppQuery.QuerySync(abci.RequestQuery{
		Path:	path,
		Data:	data,
		Height:	height,
		Prove:	!trusted,
	})
	if err != nil {
		return nil, err
	}
	logger.Info("ABCIQuery", "path", path, "data", data, "result", resQuery)
	return &ctypes.ResultABCIQuery{*resQuery}, nil
}

func ABCIInfo() (*ctypes.ResultABCIInfo, error) {
	resInfo, err := proxyAppQuery.InfoSync(abci.RequestInfo{version.Version})
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultABCIInfo{*resInfo}, nil
}
