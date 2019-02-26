package core

import (
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func Health() (*ctypes.ResultHealth, error) {
	return &ctypes.ResultHealth{}, nil
}
