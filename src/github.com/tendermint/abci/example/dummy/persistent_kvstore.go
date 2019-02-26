package dummy

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/tendermint/abci/example/code"
	"github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

const (
	ValidatorSetChangePrefix string = "val:"
)

var _ types.Application = (*PersistentDummyApplication)(nil)

type PersistentDummyApplication struct {
	app	*DummyApplication

	ValUpdates	[]types.Validator

	logger	log.Logger
}

func NewPersistentDummyApplication(dbDir string) *PersistentDummyApplication {
	name := "dummy"
	db, err := dbm.NewGoLevelDB(name, dbDir)
	if err != nil {
		panic(err)
	}

	state := loadState(db)

	return &PersistentDummyApplication{
		app:	&DummyApplication{state: state},
		logger:	log.NewNopLogger(),
	}
}

func (app *PersistentDummyApplication) SetLogger(l log.Logger) {
	app.logger = l
}

func (app *PersistentDummyApplication) Info(req types.RequestInfo) types.ResponseInfo {
	res := app.app.Info(req)
	res.LastBlockHeight = app.app.state.Height
	respAppState := types.AppState{
		BlockHeight:	app.app.state.Height,
		AppHash:	app.app.state.AppHash,
	}
	res.LastAppState = types.AppStateToByte(&respAppState)

	return res
}

func (app *PersistentDummyApplication) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return app.app.SetOption(req)
}

func (app *PersistentDummyApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {

	if isValidatorTx(tx) {

		return app.execValidatorTx(tx)
	}

	return app.app.DeliverTx(tx)
}

func (app *PersistentDummyApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	return app.app.CheckTx(tx)
}

func (app *PersistentDummyApplication) Commit() types.ResponseCommit {
	return app.app.Commit()
}

func (app *PersistentDummyApplication) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	return app.app.Query(reqQuery)
}

func (app *PersistentDummyApplication) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	for _, v := range req.Validators {
		r := app.updateValidator(v)
		if r.IsErr() {
			app.logger.Error("Error updating validators", "r", r)
		}
	}
	return types.ResponseInitChain{}
}

func (app *PersistentDummyApplication) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {

	app.ValUpdates = make([]types.Validator, 0)
	return types.ResponseBeginBlock{}
}

func (app *PersistentDummyApplication) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
}

func (app *PersistentDummyApplication) Validators() (validators []types.Validator) {
	itr := app.app.state.db.Iterator(nil, nil)
	for ; itr.Valid(); itr.Next() {
		if isValidatorTx(itr.Key()) {
			validator := new(types.Validator)
			err := types.ReadMessage(bytes.NewBuffer(itr.Value()), validator)
			if err != nil {
				panic(err)
			}
			validators = append(validators, *validator)
		}
	}
	return
}

func MakeValSetChangeTx(pubkey []byte, power int64) []byte {
	return []byte(cmn.Fmt("val:%X/%d", pubkey, power))
}

func isValidatorTx(tx []byte) bool {
	return strings.HasPrefix(string(tx), ValidatorSetChangePrefix)
}

func (app *PersistentDummyApplication) execValidatorTx(tx []byte) types.ResponseDeliverTx {
	tx = tx[len(ValidatorSetChangePrefix):]

	pubKeyAndPower := strings.Split(string(tx), "/")
	if len(pubKeyAndPower) < 3 {
		return types.ResponseDeliverTx{
			Code:	code.CodeTypeEncodingError,
			Log:	fmt.Sprintf("Expected 'pubkey/power'. Got %v", pubKeyAndPower)}
	}
	pubkeyS, powerS, rewardAddrS := pubKeyAndPower[0], pubKeyAndPower[1], pubKeyAndPower[2]

	pubkey, err := hex.DecodeString(pubkeyS)
	if err != nil {
		return types.ResponseDeliverTx{
			Code:	code.CodeTypeEncodingError,
			Log:	fmt.Sprintf("Pubkey (%s) is invalid hex", pubkeyS)}
	}

	if err != nil {
		return types.ResponseDeliverTx{
			Code:	code.CodeTypeEncodingError,
			Log:	fmt.Sprintf("RewardAddr (%s) is invalid hex", rewardAddrS)}
	}

	power, err := strconv.ParseUint(powerS, 10, 64)
	if err != nil {
		return types.ResponseDeliverTx{
			Code:	code.CodeTypeEncodingError,
			Log:	fmt.Sprintf("Power (%s) is not an int", powerS)}
	}

	return app.updateValidator(types.Validator{pubkey, power, rewardAddrS, ""})
}

func (app *PersistentDummyApplication) updateValidator(v types.Validator) types.ResponseDeliverTx {
	key := []byte("val:" + string(v.PubKey))
	if v.Power == 0 {

		if !app.app.state.db.Has(key) {
			return types.ResponseDeliverTx{
				Code:	code.CodeTypeUnauthorized,
				Log:	fmt.Sprintf("Cannot remove non-existent validator %X", key)}
		}
		app.app.state.db.Delete(key)
	} else {

		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(&v, value); err != nil {
			return types.ResponseDeliverTx{
				Code:	code.CodeTypeEncodingError,
				Log:	fmt.Sprintf("Error encoding validator: %v", err)}
		}
		app.app.state.db.Set(key, value.Bytes())
	}

	app.ValUpdates = append(app.ValUpdates, v)

	return types.ResponseDeliverTx{Code: code.CodeTypeOK}
}
