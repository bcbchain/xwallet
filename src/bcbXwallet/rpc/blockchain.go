package rpc

import (
	"bcbXwallet/common"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"bcbchain.io/keys"
	"bcbchain.io/rlp"
	tx3 "bcbchain.io/tx/tx"
	"bcbchain.io/types"
	types2 "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tendermint/rpc/core/types"
	"math/big"
	"strings"
)

var (
	genesisTokenAddr = ""
)

func blockHeight() (blkHeight *BlockHeightResult, err error) {

	result := new(core_types.ResultABCIInfo)
	params := map[string]interface{}{}
	err = common.DoHttpRequestAndParseEx(common.GetConfig().NodeAddrSlice, "abci_info", params, result)
	if err != nil {
		return
	}

	blkHeight = new(BlockHeightResult)
	blkHeight.LastBlock = result.Response.LastBlockHeight

	return
}

func block(height int64) (blk *BlockResult, err error) {

	result := new(core_types.ResultBlock)
	params := map[string]interface{}{"height": height}
	err = common.DoHttpRequestAndParseEx(common.GetConfig().NodeAddrSlice, "block", params, result)
	if err != nil {
		return
	}

	blk = new(BlockResult)
	blk.BlockHeight = result.BlockMeta.Header.Height
	blk.BlockHash = hex.EncodeToString(result.BlockMeta.BlockID.Hash)
	blk.ParentHash = hex.EncodeToString(result.BlockMeta.Header.LastBlockID.Hash)
	blk.ChainID = result.BlockMeta.Header.ChainID
	blk.ValidatorHash = hex.EncodeToString(result.BlockMeta.Header.ValidatorsHash)
	blk.ConsensusHash = hex.EncodeToString(result.BlockMeta.Header.ConsensusHash)
	blk.BlockTime = result.BlockMeta.Header.Time.String()
	blk.BlockSize = result.BlockSize
	blk.ProposerAddress = result.BlockMeta.Header.ProposerAddress

	blk.Txs = make([]TxResult, 0)
	var blkResults *core_types.ResultBlockResults
	if blkResults, err = blockResults(height); err != nil {
		return
	}

	for _, blkResult := range blkResults.Results.DeliverTx {
		var tx *TxResult
		if tx, err = transaction(hex.EncodeToString(blkResult.TxHash), result); err != nil {
			return
		}
		blk.Txs = append(blk.Txs, *tx)
	}

	return
}

func transaction(txHash string, resultBlock *core_types.ResultBlock) (tx *TxResult, err error) {

	result := new(core_types.ResultTx)
	params := map[string]interface{}{"hash": txHash}
	err = common.DoHttpRequestAndParseEx(common.GetConfig().NodeAddrSlice, "tx", params, result)
	if err != nil {
		return
	}

	if resultBlock == nil {
		resultBlock = new(core_types.ResultBlock)
		params = map[string]interface{}{"height": result.Height}
		err = common.DoHttpRequestAndParseEx(common.GetConfig().NodeAddrSlice, "block", params, resultBlock)
		if err != nil {
			return
		}
	}

	var transaction tx3.Transaction
	var fromAddr string
	fromAddr, err = transaction.TxParse(crypto.GetChainId(), string(resultBlock.Block.Txs[result.Index]))
	if err != nil {
		return
	}

	tx = new(TxResult)
	tx.TxHash = txHash
	tx.TxTime = resultBlock.BlockMeta.Header.Time.String()
	tx.Code = result.DeliverResult.Code
	tx.Log = result.DeliverResult.Log
	tx.BlockHash = hex.EncodeToString(resultBlock.BlockMeta.BlockID.Hash)
	tx.BlockHeight = resultBlock.BlockMeta.Header.Height
	tx.From = fromAddr
	tx.Nonce = transaction.Nonce
	tx.GasLimit = transaction.GasLimit
	tx.Fee = result.DeliverResult.Fee
	tx.Note = transaction.Note
	tx.Messages = make([]Message, 0)

	var msg Message
	msg, err = message(transaction)
	if err != nil {
		return
	}
	tx.Messages = append(tx.Messages, msg)

	return
}

func balance(address keys.Address) (result *BalanceResult, err error) {

	return balanceOfToken(address, genesisToken(), "")
}

func balanceOfToken(address, tokenAddress keys.Address, tokenName string) (result *BalanceResult, err error) {

	var value []byte
	if tokenName != "" {
		var tmpAddress keys.Address
		if value, err = common.DoHttpQuery(common.GetConfig().NodeAddrSlice, keyOfTokenName(tokenName)); err != nil {
			return
		}
		if len(value) == 0 {
			return nil, errors.New("invalid tokenName")
		}

		if err = json.Unmarshal(value, &tmpAddress); err != nil {
			return
		}

		if tokenAddress != "" && tokenAddress != tmpAddress {
			return nil, errors.New("tokenAddress and tokenName not be same token")
		}
		tokenAddress = tmpAddress
	} else if tokenAddress == "" {
		return nil, errors.New("tokenAddress and tokenName cannot be empty with both")
	}

	if value, err = common.DoHttpQuery(common.GetConfig().NodeAddrSlice, keyOfAccountToken(address, tokenAddress)); err != nil {
		return
	}
	result = new(BalanceResult)
	if len(value) == 0 {
		result.Balance = "0"
	} else {
		var tokenBalance types.TokenBalance
		if err = json.Unmarshal(value, &tokenBalance); err != nil {
			return
		}
		result.Balance = tokenBalance.Balance.String()
	}

	return
}

func allBalance(address keys.Address) (items *[]AllBalanceItemResult, err error) {

	tokens := make([]string, 0)
	if err = common.DoHttpQueryAndParse(common.GetConfig().NodeAddrSlice, keyOfAccount(address), &tokens); err != nil {
		return
	}

	balanceItems := make([]AllBalanceItemResult, 0)
	for _, token := range tokens {
		splitToken := strings.Split(token, "/")
		if splitToken[4] != "token" || len(splitToken) != 6 {
			continue
		}
		tokenBalance := new(types.TokenBalance)
		if err = common.DoHttpQueryAndParse(common.GetConfig().NodeAddrSlice, token, tokenBalance); err != nil {
			return
		}

		var name string
		if name, err = tokenName(tokenBalance.Address); err != nil {
			return
		}

		balanceItems = append(balanceItems,
			AllBalanceItemResult{
				TokenAddress:	tokenBalance.Address,
				TokenName:	name,
				Balance:	tokenBalance.Balance.String()})
	}

	return &balanceItems, err
}

func nonce(acctAddress keys.Address) (result *NonceResult, err error) {

	type account struct {
		Nonce uint64 `json:"nonce"`
	}

	a := new(account)
	value, err := common.DoHttpQuery(common.GetConfig().NodeAddrSlice, keyOfAccountNonce(acctAddress))
	if err != nil {
		return
	}

	result = new(NonceResult)
	if len(value) == 0 {
		result.Nonce = 1
	} else {
		err = json.Unmarshal(value, a)
		if err != nil {
			return
		}

		result.Nonce = a.Nonce + 1
	}

	return
}

func commitTx(tx string) (commit *CommitTxResult, err error) {

	var result *types.ResultBroadcastTxCommit
	result, err = common.DoHttpRequestAndParse(common.GetConfig().NodeAddrSlice, tx)
	if err != nil {
		return
	}

	commit = new(CommitTxResult)
	if result.CheckTx.Code != types2.CodeTypeOK {
		commit.Code = result.CheckTx.Code
		commit.Log = result.CheckTx.Log
	} else {
		commit.Code = result.DeliverTx.Code
		commit.Log = result.DeliverTx.Log
	}
	commit.Fee = result.DeliverTx.Fee
	commit.TxHash = hex.EncodeToString(result.Hash)
	commit.Height = result.Height

	return
}

func blockResults(height int64) (blkResults *core_types.ResultBlockResults, err error) {

	blkResults = new(core_types.ResultBlockResults)
	params := map[string]interface{}{"height": height}
	err = common.DoHttpRequestAndParseEx(common.GetConfig().NodeAddrSlice, "block_results", params, blkResults)
	if err != nil {
		return
	}

	return
}

func message(transation tx3.Transaction) (msg Message, err error) {

	var methodInfo tx3.MethodInfo
	if err = rlp.DecodeBytes(transation.Data, &methodInfo); err != nil {
		return
	}
	methodID := fmt.Sprintf("%x", methodInfo.MethodID)

	msg.SmcAddress = transation.To
	if msg.SmcName, msg.Method, err = contractNameAndMethod(transation.To, methodID); err != nil {
		return
	}

	if methodID == transferMethodID {
		var itemsBytes = make([][]byte, 0)
		if err = rlp.DecodeBytes(methodInfo.ParamData, &itemsBytes); err != nil {
			return
		}
		msg.To = string(itemsBytes[0])
		msg.Value = new(big.Int).SetBytes(itemsBytes[1][:]).String()
	}

	return
}

func contractNameAndMethod(contractAddress keys.Address, methodID string) (contractName string, method string, err error) {

	contract := new(types.Contract)
	if err = common.DoHttpQueryAndParse(common.GetConfig().NodeAddrSlice, keyOfContract(contractAddress), contract); err != nil {
		return
	}

	for _, methodItem := range contract.Methods {
		if methodItem.MethodId == methodID {
			method = methodItem.Prototype
			break
		}
	}

	return contract.Name, method, nil
}

func tokenName(tokenAddress keys.Address) (name string, err error) {

	token := new(types.IssueToken)

	if err = common.DoHttpQueryAndParse(common.GetConfig().NodeAddrSlice, keyOfToken(tokenAddress), token); err != nil {
		return
	}

	return token.Name, err
}

func genesisToken() string {
	if genesisTokenAddr == "" {
		token := new(types.IssueToken)

		if err := common.DoHttpQueryAndParse(common.GetConfig().NodeAddrSlice, keyOfGenesisToken(), token); err != nil {
			return ""
		} else {
			genesisTokenAddr = token.Address
		}
	}

	return genesisTokenAddr
}
