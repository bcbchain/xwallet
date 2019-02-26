package bcerrors

type BCError struct {
	ErrorCode	uint32
	ErrorDesc	string
}

func (bcerror *BCError) Error() string {
	if bcerror.ErrorDesc != "" {
		return bcerror.ErrorDesc
	} else {

		for _, error := range bcErrors {
			if error.ErrorCode == bcerror.ErrorCode {
				return error.ErrorDesc
			}
		}
	}
	return ""
}

const (
	ErrCodeOK = 200 + iota
)

const (
	ErrCodeDockerDupRegist	= 500 + iota
	ErrCodeDockerNotFindContract
	ErrCodeDockerNotFindDocker
)

const (
	ErrCodeNoAuthorization = 600 + iota
)

const (
	ErrCodeCheckTxTransData	= 1000 + iota
	ErrCodeCheckTxInvalidNonce
	ErrCodeCheckTxNoContract
	ErrCodeCheckTxNoteExceedLimit
)

const (
	ErrCodeDeliverTxTransData	= 2000 + iota
	ErrCodeDeliverTxInvalidNonce
	ErrCodeDeliverTxNoContract
	ErrCodeDeliverTxNoteExceedLimit
)

const (
	ErrCodeStubUnregisteredContract = 3000 + iota
)

const (
	ErrCodeInterContractsNoAuthorization	= 4000 + iota
	ErrCodeInterContractsNoGenesis
	ErrCodeInterContractsNoToken
	ErrCodeInterContractsInvalidAddr
	ErrCodeInterContractsUnsupportAddSupply
	ErrCodeInterContractsUnsupportBurn
	ErrCodeInterContractsInvalidGasLimit
	ErrCodeInterContractsInvalidGasPrice
	ErrCodeInterContractsInvalidFee
	ErrCodeInterContractsInvalidBalance
	ErrCodeInterContractsInvalidSupply
	ErrCodeInterContractsInsufficientBalance
	ErrCodeInterContractsNoValidators
	ErrCodeInterContractsInvalidRewarderAddr
	ErrCodeInterContractsInvalidPower
	ErrCodeInterContractsInvalidParameter
	ErrCodeInterContractsInvalidMethod
	ErrCodeInterContractsDupName
	ErrCodeInterContractsDupSymbol
	ErrCodeInterContractsNoStrategys
	ErrCodeInterContractsOutOfRange
	ErrCodeInterContractsUnsupportTransToSelf
	ErrCodeInterContractsTokenNotInit
	ErrCodeInterContractsUnfinished
	ErrCodeInterContractsNameTooLong
	ErrCodeInterContractsSymbolTooLong
	ErrCodeInterContractsInvalidPercent
	ErrCodeInterContractsEmptyName
	ErrCodeInterContractsLoseNameOfValidators
	ErrCodeInterContractsRuntimeError
	ErrCodeInterContractsInvalidTeam
	ErrCodeInterContractsSenderInBlackList
)

const (
	ErrCodeLowLevelError = 5000 + iota
)

const (
	ErrCodeInterContractsMinUserCode	= 8000 + iota
	ErrCodeInterContractsBeyondMaximumStockHolders
	ErrCodeInterContractsBybInitialized
	ErrCodeInterContractsBybOwnedByb
	ErrCodeInterContractsBybHolderNotFound
)

var bcErrors = []BCError{
	{ErrCodeOK, ""},

	{ErrCodeDockerDupRegist, "Failed to register contract due to it has been registered into docker"},
	{ErrCodeDockerNotFindContract, "Did not find contract from docker"},
	{ErrCodeDockerNotFindDocker, "Did not find the docker"},

	{ErrCodeNoAuthorization, "No authorization"},

	{ErrCodeCheckTxTransData, "Transaction data is invalid"},
	{ErrCodeCheckTxInvalidNonce, "Nonce is invalid"},
	{ErrCodeCheckTxNoContract, "Did not find this contract"},
	{ErrCodeCheckTxNoteExceedLimit, "Invalid note, it must be stay within 256 characters limit"},

	{ErrCodeDeliverTxTransData, "Transaction data is invalid"},
	{ErrCodeDeliverTxInvalidNonce, "Nonce is invalid"},
	{ErrCodeDeliverTxNoContract, "Did not find this contract"},
	{ErrCodeDeliverTxNoteExceedLimit, "Invalid note, it must be stay within 256 characters limit"},

	{ErrCodeStubUnregisteredContract, "Contract did not be registered into docker"},

	{ErrCodeInterContractsNoAuthorization, "No authorization to execute contract"},
	{ErrCodeInterContractsNoGenesis, "No genesis"},
	{ErrCodeInterContractsNoToken, "Did not find this token"},
	{ErrCodeInterContractsInvalidAddr, "The specified contract address is invalid"},
	{ErrCodeInterContractsUnsupportAddSupply, "The token does not support to add supply"},
	{ErrCodeInterContractsUnsupportBurn, "The token does not support to burn"},
	{ErrCodeInterContractsInvalidGasLimit, "The specified gaslimit is too less to execute"},
	{ErrCodeInterContractsInvalidGasPrice, "The token's gasprice is invalid"},
	{ErrCodeInterContractsInvalidFee, "The fee of operation is invalid"},
	{ErrCodeInterContractsInvalidBalance, "The account's balance is invalid"},
	{ErrCodeInterContractsInvalidSupply, "The token supply is incorrect"},
	{ErrCodeInterContractsInsufficientBalance, "The accounts' balance is insufficient"},
	{ErrCodeInterContractsNoValidators, "Did not get validators list"},
	{ErrCodeInterContractsInvalidRewarderAddr, "The proposer's reward address is mismatch"},
	{ErrCodeInterContractsInvalidPower, "Validator's power is invalid"},
	{ErrCodeInterContractsInvalidParameter, "Invalid parameter"},
	{ErrCodeInterContractsInvalidMethod, "The specified method is unsupported"},
	{ErrCodeInterContractsDupName, "The specified token name is duplicated with existing one"},
	{ErrCodeInterContractsDupSymbol, "The specified token symbol is duplicated with existing one"},
	{ErrCodeInterContractsNoStrategys, "Did not get reward strategy"},
	{ErrCodeInterContractsOutOfRange, "It exceeds the limit"},
	{ErrCodeInterContractsUnsupportTransToSelf, "Do not support to transfer to yourself"},
	{ErrCodeInterContractsTokenNotInit, "Invalid operation for token what is not completed issued yet"},
	{ErrCodeInterContractsUnfinished, "Unfinished contract, coming soon"},
	{ErrCodeInterContractsNameTooLong, "Name is too long"},
	{ErrCodeInterContractsSymbolTooLong, "Symbol is too long"},
	{ErrCodeInterContractsInvalidPercent, "Percentage of reward is wrong"},
	{ErrCodeInterContractsEmptyName, "Name cannot be empty"},
	{ErrCodeInterContractsLoseNameOfValidators, "There is no validators in reward strategy list"},
	{ErrCodeInterContractsRuntimeError, "Running time error"},

	{ErrCodeInterContractsInvalidTeam, "The team is invalid"},

	{ErrCodeLowLevelError, ""},
	{ErrCodeInterContractsMinUserCode, ""},
	{ErrCodeInterContractsBeyondMaximumStockHolders, "The number of stockHolder exceeds the limit"},
	{ErrCodeInterContractsBybInitialized, "The token byb has been initialized"},
	{ErrCodeInterContractsBybOwnedByb, "The stockholder is still owning byb token"},
	{ErrCodeInterContractsBybHolderNotFound, "The stockholder was not found"},

	{ErrCodeInterContractsSenderInBlackList, "This sender is in black list"},
}
