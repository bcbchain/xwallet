package tx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"

	atm "bcbchain.io/algorithm"
	"bcbchain.io/prototype"
	"bcbchain.io/rlp"
	"bcbchain.io/tx/contract/baccarat.v1.0"
	"bcbchain.io/tx/contract/blacklist.v1.0"
	"bcbchain.io/tx/contract/common"
	"bcbchain.io/tx/contract/crypto-gods.v1.0"
	"bcbchain.io/tx/contract/dc_yuebao.v.1.0"
	"bcbchain.io/tx/contract/dice2win"
	"bcbchain.io/tx/contract/dragonvstiger.v1.0"
	"bcbchain.io/tx/contract/everycolor.v1.0"
	"bcbchain.io/tx/contract/incentive.v1.0"
	"bcbchain.io/tx/contract/system.v1.0"
	"bcbchain.io/tx/contract/token-basic.v1.0"
	"bcbchain.io/tx/contract/token-byb.v1.0"
	"bcbchain.io/tx/contract/token-clt.v1.0"
	"bcbchain.io/tx/contract/token-issue.v1.0"
	"bcbchain.io/tx/contract/token-templet.v1.0"
	"bcbchain.io/tx/contract/token-united.v1.0"
	"bcbchain.io/tx/contract/transferagency.v1.0"

	"bcbchain.io/tx/contract/dice2win.v2.0"
	"bcbchain.io/tx/tx"
)

var chainId string

var methodIdToItems = map[uint32]interface{}{}

func InitUnWrapper(config string) error {

	chainId = config

	methodIdToItems[ConvertPrototype2ID(prototype.CltSetOwner)] = token_clt_v1_0.DecodeSetOwnerItems
	methodIdToItems[ConvertPrototype2ID(prototype.CltTriggerDistribute)] = token_clt_v1_0.DecodeTriggerDistributeItems
	methodIdToItems[ConvertPrototype2ID(prototype.CltBuy)] = token_clt_v1_0.DecodeBuyItems
	methodIdToItems[ConvertPrototype2ID(prototype.CltReinvest)] = token_clt_v1_0.DecodeReinvestItems
	methodIdToItems[ConvertPrototype2ID(prototype.CltExit)] = token_clt_v1_0.DecodeExitItems
	methodIdToItems[ConvertPrototype2ID(prototype.CltWithdraw)] = token_clt_v1_0.DecodeWithdrawItems
	methodIdToItems[ConvertPrototype2ID(prototype.CltSell)] = token_clt_v1_0.DecodeSellItems
	methodIdToItems[ConvertPrototype2ID(prototype.CltTransfer)] = token_clt_v1_0.DecodeTransferItems
	methodIdToItems[ConvertPrototype2ID(prototype.CltSetStakingRequirement)] = token_clt_v1_0.DecodeSetStakingRequirementItems
	methodIdToItems[ConvertPrototype2ID(prototype.CltCalculateTokenReceived)] = token_clt_v1_0.DecodeCalculateTokenReceivedItems
	methodIdToItems[ConvertPrototype2ID(prototype.CltCalculateBCBReceived)] = token_clt_v1_0.DecodeCalculateBCBReceivedItems
	methodIdToItems[ConvertPrototype2ID(prototype.CltGetGlobalData)] = token_clt_v1_0.DecodeGetGlobalDataItems
	methodIdToItems[ConvertPrototype2ID(prototype.CltGetPlayerInfo)] = token_clt_v1_0.DecodeGetPlayerInfoItems

	methodIdToItems[ConvertPrototype2ID(prototype.BYBInit)] = token_byb_v1_0.DecodeInitItems
	methodIdToItems[ConvertPrototype2ID(prototype.BYBSetOwner)] = token_byb_v1_0.DecodeSetOwnerItems
	methodIdToItems[ConvertPrototype2ID(prototype.BYBSetGasPrice)] = token_byb_v1_0.DecodeSetGasPriceItems
	methodIdToItems[ConvertPrototype2ID(prototype.BYBAddSupply)] = token_byb_v1_0.DecodeAddSupplyItems
	methodIdToItems[ConvertPrototype2ID(prototype.BYBBurn)] = token_byb_v1_0.DecodeBurnItems
	methodIdToItems[ConvertPrototype2ID(prototype.BYBNewBlackHole)] = token_byb_v1_0.DecodeNewBlackHoleItems
	methodIdToItems[ConvertPrototype2ID(prototype.BYBNewStockHolder)] = token_byb_v1_0.DecodeNewStockHolderItems
	methodIdToItems[ConvertPrototype2ID(prototype.BYBDelStockHolder)] = token_byb_v1_0.DecodeDelStockHolderItems
	methodIdToItems[ConvertPrototype2ID(prototype.BYBChangeChromoOwnerShip)] = token_byb_v1_0.DecodeChangeChromoOwnerShipItems
	methodIdToItems[ConvertPrototype2ID(prototype.BYBTransfer)] = token_byb_v1_0.DecodeTransferItems
	methodIdToItems[ConvertPrototype2ID(prototype.BYBTransferByChromo)] = token_byb_v1_0.DecodeTransferByChromoItems

	methodIdToItems[ConvertPrototype2ID(prototype.SysNewValidator)] = system_v1_0.DecodeNewValidatorItems
	methodIdToItems[ConvertPrototype2ID(prototype.SysSetPower)] = system_v1_0.DecodeSetPowerItems
	methodIdToItems[ConvertPrototype2ID(prototype.SysSetRewardAddr)] = system_v1_0.DecodeSetRewardAddrItems
	methodIdToItems[ConvertPrototype2ID(prototype.SysForbidInternalContract)] = system_v1_0.DecodeForbidInternalContractItems
	methodIdToItems[ConvertPrototype2ID(prototype.SysDeployInternalContract)] = system_v1_0.DecodeDeployInternalContractItems
	methodIdToItems[ConvertPrototype2ID(prototype.SysSetRewardStrategy)] = system_v1_0.DecodeSetRewardStrategyItems

	methodIdToItems[ConvertPrototype2ID(prototype.TbTransfer)] = token_basic_v1_0.DecodeTransferItems
	methodIdToItems[ConvertPrototype2ID(prototype.TbSetGasBasePrice)] = token_basic_v1_0.DecodeSetGasBasePriceItems
	methodIdToItems[ConvertPrototype2ID(prototype.TbSetGasPrice)] = token_basic_v1_0.DecodeSetGasPriceItems

	methodIdToItems[ConvertPrototype2ID(prototype.TiNewToken)] = token_issue_v1_0.DecodeNewTokenItems

	methodIdToItems[ConvertPrototype2ID(prototype.TtTransfer)] = token_templet_v1_0.DecodeTransferItems
	methodIdToItems[ConvertPrototype2ID(prototype.TtBatchTransfer)] = token_templet_v1_0.DecodeBatchTransferItems
	methodIdToItems[ConvertPrototype2ID(prototype.TtAddSupply)] = token_templet_v1_0.DecodeAddSupplyItems
	methodIdToItems[ConvertPrototype2ID(prototype.TtBurn)] = token_templet_v1_0.DecodeBurnItems
	methodIdToItems[ConvertPrototype2ID(prototype.TtSetGasPrice)] = token_templet_v1_0.DecodeSetGasPriceItems
	methodIdToItems[ConvertPrototype2ID(prototype.TtSetOwner)] = token_templet_v1_0.DecodeSetOwnerItems

	methodIdToItems[ConvertPrototype2ID(prototype.CgsInit)] = crypto_gods_v1_0.DecodeInitItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsActivate)] = crypto_gods_v1_0.DecodeActivateItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsBuy)] = crypto_gods_v1_0.DecodeBuyItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsBuyXid)] = crypto_gods_v1_0.DecodeBuyXidItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsBuyXaddr)] = crypto_gods_v1_0.DecodeBuyXaddrItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsBuyXname)] = crypto_gods_v1_0.DecodeBuyXnameItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsSetOwner)] = crypto_gods_v1_0.DecodeSetOwnerItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsSetSetting)] = crypto_gods_v1_0.DecodeSetSettingItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsGetCurrentRoundInfo)] = crypto_gods_v1_0.DecodeGetCurrentRoundInfoItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsGetBuyPrice)] = crypto_gods_v1_0.DecodeGetBuyPriceItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsGetPrice)] = crypto_gods_v1_0.DecodeGetPriceItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsGetKeys)] = crypto_gods_v1_0.DecodeGetKeysItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsGetPlayerInfoByAddress)] = crypto_gods_v1_0.DecodeGetPlayerInfoByAddressItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsGetPlayerVault)] = crypto_gods_v1_0.DecodeGetPlayerVaultItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsGetTimeLeft)] = crypto_gods_v1_0.DecodeGetTimeLeftItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsWithDraw)] = crypto_gods_v1_0.DecodeWithDrawItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsReloadXid)] = crypto_gods_v1_0.DecodeReloadXidItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsReloadXaddr)] = crypto_gods_v1_0.DecodeReloadXaddrItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsReloadXname)] = crypto_gods_v1_0.DecodeReloadXnameItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsRegisterNameXid)] = crypto_gods_v1_0.DecodeRegisterNameXidItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsRegisterNameXaddr)] = crypto_gods_v1_0.DecodeRegisterNameXaddrItems
	methodIdToItems[ConvertPrototype2ID(prototype.CgsRegisterNameXname)] = crypto_gods_v1_0.DecodeRegisterNameXnameItems

	methodIdToItems[ConvertPrototype2ID(prototype.BlmAddAddress)] = blacklist_v1_0.DecodeAddAddressItems
	methodIdToItems[ConvertPrototype2ID(prototype.BlmDelAddress)] = blacklist_v1_0.DecodeDelAddressItems
	methodIdToItems[ConvertPrototype2ID(prototype.BlmSetOwner)] = blacklist_v1_0.DecodeSetOwnerItems

	methodIdToItems[ConvertPrototype2ID(prototype.UtNewToken)] = token_united_v1_0.DecodeNewTokenItems
	methodIdToItems[ConvertPrototype2ID(prototype.UtTransfer)] = token_united_v1_0.DecodeTransferItems
	methodIdToItems[ConvertPrototype2ID(prototype.UtSetOwner)] = token_united_v1_0.DecodeSetOwnerItems
	methodIdToItems[ConvertPrototype2ID(prototype.UtSetGasPrice)] = token_united_v1_0.DecodeSetGasPriceItems
	methodIdToItems[ConvertPrototype2ID(prototype.UtSetFee)] = token_united_v1_0.DecodeSetFeeItems
	methodIdToItems[ConvertPrototype2ID(prototype.UtSetGasPayer)] = token_united_v1_0.DecodeSetGasPayerItems
	methodIdToItems[ConvertPrototype2ID(prototype.UtAddSupply)] = token_united_v1_0.DecodeAddSupplyItems
	methodIdToItems[ConvertPrototype2ID(prototype.UtBurn)] = token_united_v1_0.DecodeBurnItems
	methodIdToItems[ConvertPrototype2ID(prototype.UtWithdraw)] = common.DecodeNoItems

	methodIdToItems[ConvertPrototype2ID(prototype.DWSetOwner)] = dice2win.DecodeDWSetOwnerItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWSetSecretSigner)] = dice2win.DecodeDWSetSecretSignerItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWSetSettings)] = dice2win.DecodeDWSetSettingsItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWWithdrawFunds)] = dice2win.DecodeDWWithdrawFundsItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWPlaceBet)] = dice2win.DecodeDWPlaceBetItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWSettleBet)] = dice2win.DecodeDWSettleBetItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWRefundBet)] = dice2win.DecodeDWRefundBetItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWClearStorage)] = dice2win.DecodeDWClearStorageItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWSetRecFeeInfo)] = dice2win.DecodeDWSetRecFeeInfoItems

	methodIdToItems[ConvertPrototype2ID(prototype.ECSetOwner)] = everycolor_v1_0.DecodeECSetOwnerItems
	methodIdToItems[ConvertPrototype2ID(prototype.ECSetSecretSigner)] = everycolor_v1_0.DecodeECSetSecretSignerItems
	methodIdToItems[ConvertPrototype2ID(prototype.ECSetSettings)] = everycolor_v1_0.DecodeECSetSettingsItems
	methodIdToItems[ConvertPrototype2ID(prototype.ECSetRecvFeeInfo)] = everycolor_v1_0.DecodeECSetRecFeeInfoItems
	methodIdToItems[ConvertPrototype2ID(prototype.ECWithdrawFunds)] = everycolor_v1_0.DecodeECWithdrawFundsItems
	methodIdToItems[ConvertPrototype2ID(prototype.ECPlaceBet)] = everycolor_v1_0.DecodeECPlaceBetItems
	methodIdToItems[ConvertPrototype2ID(prototype.ECSettleBets)] = everycolor_v1_0.DecodeECSettleBetsItems
	methodIdToItems[ConvertPrototype2ID(prototype.ECRefundBets)] = everycolor_v1_0.DecodeECRefundBetsItems
	methodIdToItems[ConvertPrototype2ID(prototype.ECWithdrawWin)] = everycolor_v1_0.DecodeECWithdrawWinItems

	methodIdToItems[ConvertPrototype2ID(prototype.DTSetOwner)] = dragonvstiger_v1_0.DecodeDTSetOwnerItems
	methodIdToItems[ConvertPrototype2ID(prototype.DTSetSecretSigner)] = dragonvstiger_v1_0.DecodeDTSetSecretSignerItems
	methodIdToItems[ConvertPrototype2ID(prototype.DTSetSettings)] = dragonvstiger_v1_0.DecodeDTSetSettingsItems
	methodIdToItems[ConvertPrototype2ID(prototype.DTWithdrawFunds)] = dragonvstiger_v1_0.DecodeDTWithdrawFundsItems
	methodIdToItems[ConvertPrototype2ID(prototype.DTPlaceBet)] = dragonvstiger_v1_0.DecodeDTPlaceBetItems
	methodIdToItems[ConvertPrototype2ID(prototype.DTSettleBets)] = dragonvstiger_v1_0.DecodeDTSettleBetsItems
	methodIdToItems[ConvertPrototype2ID(prototype.DTSetRecvFeeInfo)] = dragonvstiger_v1_0.DecodeDTSetRecvFeeInfoItems
	methodIdToItems[ConvertPrototype2ID(prototype.DTRefundBets)] = dragonvstiger_v1_0.DecodeDTRefundBetsItems
	methodIdToItems[ConvertPrototype2ID(prototype.DTWithdrawWin)] = dragonvstiger_v1_0.DecodeDTWithdrawWinItems

	methodIdToItems[ConvertPrototype2ID(prototype.YEBSetOwner)] = dc_yuebao_v_1_0.DecodeYEBSetOwnerItems
	methodIdToItems[ConvertPrototype2ID(prototype.YEBPromiseEarningRate)] = dc_yuebao_v_1_0.DecodeYEBPromiseEarningRateItems
	methodIdToItems[ConvertPrototype2ID(prototype.YEBDepositsReceived)] = dc_yuebao_v_1_0.DecodeYEBDepositsReceivedItems
	methodIdToItems[ConvertPrototype2ID(prototype.YEBDepositsWithdrawal)] = dc_yuebao_v_1_0.DecodeYEBDepositsWithdrawalItems
	methodIdToItems[ConvertPrototype2ID(prototype.YEBDeposit)] = dc_yuebao_v_1_0.DecodeYEBDepositItems
	methodIdToItems[ConvertPrototype2ID(prototype.YEBWithdraw)] = dc_yuebao_v_1_0.DecodeYEBWithdrawItems
	methodIdToItems[ConvertPrototype2ID(prototype.YEBWithdrawEarnings)] = dc_yuebao_v_1_0.DecodeYEBWithdrawEarningsItems
	methodIdToItems[ConvertPrototype2ID(prototype.YEBReinvest)] = dc_yuebao_v_1_0.DecodeYEBReinvestItems
	methodIdToItems[ConvertPrototype2ID(prototype.YEBSetOperators)] = dc_yuebao_v_1_0.DecodeYEBSetOperatorsItems

	methodIdToItems[ConvertPrototype2ID(prototype.BACSetOwner)] = baccarat_v1_0.DecodeBACSetOwnerItems
	methodIdToItems[ConvertPrototype2ID(prototype.BACSetSecretSigner)] = baccarat_v1_0.DecodeBACSetSecretSignerItems
	methodIdToItems[ConvertPrototype2ID(prototype.BACSetSettings)] = baccarat_v1_0.DecodeBACSetSettingsItems
	methodIdToItems[ConvertPrototype2ID(prototype.BACWithdrawFunds)] = baccarat_v1_0.DecodeBACWithdrawFundsItems
	methodIdToItems[ConvertPrototype2ID(prototype.BACPlaceBet)] = baccarat_v1_0.DecodeBACPlaceBetItems
	methodIdToItems[ConvertPrototype2ID(prototype.BACSettleBets)] = baccarat_v1_0.DecodeBACSettleBetsItems
	methodIdToItems[ConvertPrototype2ID(prototype.BACSetRecvFeeInfo)] = baccarat_v1_0.DecodeBACSetRecvFeeInfoItems
	methodIdToItems[ConvertPrototype2ID(prototype.BACRefundBets)] = baccarat_v1_0.DecodeBACRefundBetsItems
	methodIdToItems[ConvertPrototype2ID(prototype.BACWithdrawWin)] = baccarat_v1_0.DecodeBACWithdrawWinItems

	methodIdToItems[ConvertPrototype2ID(prototype.ICTSetManager)] = incentive_v_1_0.DecodeICTSetManagerItems
	methodIdToItems[ConvertPrototype2ID(prototype.ICTSetPayManager)] = incentive_v_1_0.DecodeICTSetPayManagerItems
	methodIdToItems[ConvertPrototype2ID(prototype.ICTSetEmployeeBonuses)] = incentive_v_1_0.DecodeICTSetEmployeeBonusesItems
	methodIdToItems[ConvertPrototype2ID(prototype.ICTPayBonuses)] = incentive_v_1_0.DecodeICTPayBonusesItems

	methodIdToItems[ConvertPrototype2ID(prototype.TACSetManager)] = transferagency_v1_0.DecodeTACSetManagerItems
	methodIdToItems[ConvertPrototype2ID(prototype.TACSetTokenFee)] = transferagency_v1_0.DecodeTACSetTokenFeeItems
	methodIdToItems[ConvertPrototype2ID(prototype.TACTransfer)] = transferagency_v1_0.DecodeTACTransferItems
	methodIdToItems[ConvertPrototype2ID(prototype.TACWithdrawFunds)] = transferagency_v1_0.DecodeTACWithdrawFundsItems

	methodIdToItems[ConvertPrototype2ID(prototype.DWSetOwner2_0)] = dice2win2_0.DecodeDWSetOwnerItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWSetSecretSigner2_0)] = dice2win2_0.DecodeDWSetSecretSignerItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWSetSettings2_0)] = dice2win2_0.DecodeDWSetSettingsItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWSetRecvFeeInfo2_0)] = dice2win2_0.DecodeDWSetRecFeeInfoItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWWithdrawFunds2_0)] = dice2win2_0.DecodeDWWithdrawFundsItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWPlaceBet2_0)] = dice2win2_0.DecodeDWPlaceBetItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWSettleBets2_0)] = dice2win2_0.DecodeDWSettleBetsItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWRefundBets2_0)] = dice2win2_0.DecodeDWRefundBetsItems
	methodIdToItems[ConvertPrototype2ID(prototype.DWWithdrawWin2_0)] = dice2win2_0.DecodeDWWithdrawWinItems

	return nil
}

func ConvertPrototype2ID(prototype string) uint32 {
	var id uint32
	bytesBuffer := bytes.NewBuffer(atm.CalcMethodId(prototype))
	binary.Read(bytesBuffer, binary.BigEndian, &id)
	return id
}

type resTx struct {
	Code		string		`json:"code,omitempty"`
	FromAddr	string		`json:"fromAddr,omitempty"`
	Nonce		string		`json:"nonce,omitempty"`
	GasLimit	string		`json:"gasLimit,omitempty"`
	Note		string		`json:"note,omitempty"`
	ToContractAddr	string		`json:"toContractAddr,omitempty"`
	MethodId	string		`json:"methodId,omitempty"`
	Items		[]string	`json:"items,omitempty"`
}

func UnpackAndParseTx(strTx string) string {

	var transaction tx.Transaction
	fromAddr, err := transaction.TxParse(chainId, strTx)
	if err != nil {
		errInfo := string("{\"code\":-2, \"message\":\"Transaction.TxParse failed(") + err.Error() + ")\",\"data\":\"\"}"
		return errInfo
	}

	var methodInfo tx.MethodInfo
	err = rlp.DecodeBytes(transaction.Data, &methodInfo)
	if err != nil {
		errInfo := string("{\"code\":-2, \"message\":\"rlp.DecodeBytes failed(") + err.Error() + ")\",\"data\":\"\"}"
		return errInfo
	}

	var methodID string
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, methodInfo.MethodID)
	methodID = hex.EncodeToString(buf)

	resultItems, err := callFunc(methodInfo.MethodID, methodInfo.ParamData)

	if err != nil {
		return string("{\"code\":-2, \"message\":\"") + err.Error() + "\",\"data\":\"\"}"
	}

	var items []string
	items, _ = resultItems[0].Interface().([]string)

	tx := resTx{
		"0",
		fromAddr,
		strconv.FormatUint(transaction.Nonce, 10),
		strconv.FormatUint(transaction.GasLimit, 10),
		transaction.Note,
		transaction.To,
		methodID,
		items,
	}

	res, _ := json.Marshal(&tx)
	return string(res)
}

func ByteSliceToInt64(b []byte) int64 {
	buf := bytes.NewBuffer(b)
	var v int64
	binary.Read(buf, binary.BigEndian, &v)
	return v
}

func callFunc(id uint32, params ...interface{}) (result []reflect.Value, err error) {

	items, ok := methodIdToItems[id]
	if !ok {
		err = errors.New("The specified method is unsupported")
		return
	}

	f := reflect.ValueOf(items)
	if f.IsNil() {
		err = errors.New("The specified method is unsupported")
		return
	}
	if len(params) != f.Type().NumIn() {
		err = errors.New("Invalid number of params passed.")
		return
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result = f.Call(in)
	return
}
