package rpc

import (
	"bcbXwallet/common"
	"errors"
	"bcbchain.io/keys"
	"github.com/tendermint/go-crypto"
	"io/ioutil"
)

func WalletCreate(name, password string) (result *WalletCreateResult, err error) {
	logger := common.GetLogger()

	defer common.FuncRecover(logger, &err)
	logger.Trace("bcb_walletCreate", "name", name)

	if err = checkName(name); err != nil {
		return
	}

	if password == "" || !checkPassword(password) {
		return nil, errors.New("Password contains by all printable ASCII and length must be [8-20] ")
	}

	result = new(WalletCreateResult)
	result, err = walletCreate(name, password)
	if err != nil {
		logger.Error("Cannot create wallet", "error", err)
	}

	return
}

func WalletExport(name, password, accessKey string, plainText bool) (result *WalletExportResult, err error) {
	logger := common.GetLogger()

	defer common.FuncRecover(logger, &err)
	logger.Trace("bcb_walletExport", "name", name, "plainText", plainText)

	if err = checkName(name); err != nil {
		return
	}

	if password == "" || !checkPassword(password) {
		return nil, errors.New("Password contains by all printable ASCII and length must be [8-20] ")
	}

	if accessKey == "" {
		return nil, errors.New("The accessKey can not be empty ")
	}

	result, err = walletExport(name, password, accessKey, plainText)
	if err != nil {
		logger.Error("Cannot export wallet", "error", err)
	}

	return
}

func WalletImport(name, privateKey, password, accessKey string, plainText bool) (result *WalletImportResult, err error) {
	logger := common.GetLogger()

	defer common.FuncRecover(logger, &err)
	logger.Trace("bcb_walletImport", "name", name)

	if err = checkName(name); err != nil {
		return
	}

	if err = checkPrivateKey(privateKey, plainText); err != nil {
		return
	}

	if password == "" || !checkPassword(password) {
		return nil, errors.New("Password contains by all printable ASCII and length must be [8-20] ")
	}

	result, err = walletImport(name, privateKey, password, accessKey, plainText)
	if err != nil {
		logger.Error("Cannot import wallet", "error", err)
	}

	return
}

func WalletList(pageNum uint64) (result *WalletListResult, err error) {
	logger := common.GetLogger()

	defer common.FuncRecover(logger, &err)
	logger.Trace("bcb_walletList")

	result, err = walletList(pageNum)
	if err != nil {
		logger.Error("Cannot list wallet", "error", err)
	}

	return
}

func WalletTransfer(name, accessKey string, walletParams TransferParam) (result *TransferResult, err error) {
	logger := common.GetLogger()

	defer common.FuncRecover(logger, &err)
	logger.Trace("bcb_transfer", "name", name, "gasLimit", walletParams.GasLimit, "note", walletParams.Note, "to", walletParams.To, "Value", walletParams.Value)

	if err = checkName(name); err != nil {
		return
	}

	gasLimit, err := requireUint64(walletParams.GasLimit)
	if err != nil {
		return
	}

	if _, err = requireUint64(walletParams.Value); err != nil {
		return
	}

	if err = checkAddress(crypto.GetChainId(), walletParams.SmcAddress); err != nil {
		return
	}

	if err = checkAddress(crypto.GetChainId(), walletParams.To); err != nil {
		return
	}

	result, err = transfer(name, accessKey, gasLimit, walletParams)
	if err != nil {
		logger.Error("Cannot transfer", "error", err)
	}

	return
}

func WalletTransferOffline(name, accessKey string, walletParams TransferOfflineParam) (result *TransferOfflineResult, err error) {
	logger := common.GetLogger()

	defer common.FuncRecover(logger, &err)
	logger.Trace("bcb_transferOffline", "name", name, "gasLimit", walletParams.GasLimit, "note", walletParams.Note, "to", walletParams.To, "Value", walletParams.Value)

	if err = checkName(name); err != nil {
		return
	}

	gasLimit, err := requireUint64(walletParams.GasLimit)
	if err != nil {
		return
	}

	if _, err = requireUint64(walletParams.Value); err != nil {
		return
	}

	if err = checkAddress(crypto.GetChainId(), walletParams.SmcAddress); err != nil {
		return
	}

	if err = checkAddress(crypto.GetChainId(), walletParams.To); err != nil {
		return
	}

	result, err = walletTransferOffline(name, accessKey, gasLimit, walletParams)
	if err != nil {
		logger.Error("Cannot pack transfer transaction", "error", err)
	}

	return
}

func BlockHeight() (result *BlockHeightResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_blockHeight")

	result, err = blockHeight()
	if err != nil {
		common.GetLogger().Error("Cannot get current block height", "error", err)
	}

	return
}

func Block(height int64) (result *BlockResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_block", "height", height)

	if height == 0 {
		var blkHeight *BlockHeightResult
		if blkHeight, err = blockHeight(); err != nil {
			common.GetLogger().Error("Cannot get current block height", "error", err)
			return
		}
		height = blkHeight.LastBlock
	} else if height < 0 {
		return nil, errors.New("Height cannot be negative ")
	}

	result, err = block(height)
	if err != nil {
		common.GetLogger().Error("Cannot get block data", "height", height, "error", err)
	}

	return
}

func Transaction(txHash string) (result *TxResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_transaction", "txHash", txHash)

	if txHash == "" {
		return nil, errors.New("TxHash cannot be empty ")
	}

	result, err = transaction(txHash, nil)
	if err != nil {
		common.GetLogger().Error("Cannot get transaction data", "error", err)
	}

	return
}

func Balance(address keys.Address) (result *BalanceResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_balance", "address", address)

	if address == "" {
		return nil, errors.New("Address cannot be empty ")
	}

	if err = checkAddress(crypto.GetChainId(), address); err != nil {
		return
	}

	result, err = balance(address)
	if err != nil {
		common.GetLogger().Error("Cannot get balance", "error", err)
	}

	return
}

func BalanceOfToken(address, tokenAddress keys.Address, tokenName string) (result *BalanceResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_balanceOfToken", "address", address, "tokenAddress", tokenAddress, "tokenName", tokenName)

	if address == "" {
		return nil, errors.New("Address cannot be empty ")
	}

	if err = checkAddress(crypto.GetChainId(), address); err != nil {
		return
	}

	if tokenAddress != "" {
		if err = checkAddress(crypto.GetChainId(), tokenAddress); err != nil {
			return
		}
	} else if tokenName == "" {
		return nil, errors.New("TokenAddress and TokenName cannot empty with both ")
	}

	result, err = balanceOfToken(address, tokenAddress, tokenName)
	if err != nil {
		common.GetLogger().Error("Cannot get balance of token", "error", err)
	}

	return
}

func AllBalance(address keys.Address) (result *[]AllBalanceItemResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_allBalance", "address", address)

	if err = checkAddress(crypto.GetChainId(), address); err != nil {
		return
	}

	result, err = allBalance(address)
	if err != nil {
		common.GetLogger().Error("Cannot get all balance", "error", err)
	}

	return
}

func Nonce(address keys.Address) (result *NonceResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_nonce", "address", address)

	if err = checkAddress(crypto.GetChainId(), address); err != nil {
		return
	}

	result, err = nonce(address)
	if err != nil {
		common.GetLogger().Error("Cannot get nonce", "error", err)
	}

	return
}

func CommitTx(tx string) (result *CommitTxResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_commitTx", "tx", tx)

	if tx == "" {
		return nil, errors.New("Tx cannot be empty ")
	}

	result, err = commitTx(tx)
	if err != nil {
		common.GetLogger().Error("Cannot commit tx", "error", err)
	}

	return
}

func Version() (result *VersionResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_version")

	var version []byte
	version, err = ioutil.ReadFile("./.config/version")
	if err != nil {
		common.GetLogger().Error("Read version file error", "error", err)
		return
	}
	result = new(VersionResult)
	result.Version = string(version)

	return
}
