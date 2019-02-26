package dummy

import (
	"github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
)

func RandVal(i int) types.Validator {
	pubkey := cmn.RandBytes(33)
	power := cmn.RandUint16() + 1
	return types.Validator{pubkey, uint64(power), "", ""}
}

func RandVals(cnt int) []types.Validator {
	res := make([]types.Validator, cnt)
	for i := 0; i < cnt; i++ {
		res[i] = RandVal(i)
	}
	return res
}

func InitDummy(app *PersistentDummyApplication) {
	app.InitChain(types.RequestInitChain{
		Validators:	RandVals(1),
		AppStateBytes:	[]byte("[]"),
	})
}
