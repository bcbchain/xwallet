package main

import (
	"bcbXwallet/client"
	"bcbXwallet/common"
	"bcbXwallet/rpc"
	"fmt"
	"bcbchain.io/rpc/lib/server"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	cmn "github.com/tendermint/tmlibs/common"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unitest/bcbXwallettest/bcbxcmn"
)

const (
	usage = "bcbXwallet_rpc's url"
)

func main() {
	err := common.InitAll()
	if err != nil {
		panic(err)
	}

	if bcbxcmn.Version() == "" {
		err = rpc.InitDB()
		if err != nil {
			panic(err)
		}

		rpcLogger := common.GetLogger()

		coreCodec := amino.NewCodec()

		mux := http.NewServeMux()

		rpcserver.RegisterRPCFuncs(mux, rpc.Routes, coreCodec, rpcLogger)
		if common.GetConfig().UseHttps {
			crtPath, keyPath := common.OutCertFileIsExist()
			_, err = rpcserver.StartHTTPAndTLSServer(serverAddr(common.GetConfig().ServerAddr, false), mux, crtPath, keyPath, rpcLogger)
			if err != nil {
				cmn.Exit(err.Error())
			}
		} else {
			_, err = rpcserver.StartHTTPServer(serverAddr(common.GetConfig().ServerAddr, false), mux, rpcLogger)
			if err != nil {
				cmn.Exit(err.Error())
			}
		}
	}

	err = Execute()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func serverAddr(address string, isHttps bool) string {
	splitAddr := strings.Split(address, ":")

	if len(splitAddr) != 3 {
		fmt.Println("invalid serverAddr=" + address)
		return ""
	}

	port, err := strconv.Atoi(splitAddr[2])
	if err != nil {
		fmt.Println("invalid serverAddr=" + address)
		return ""
	}

	if isHttps {
		if common.GetConfig().UseHttps {
			return fmt.Sprintf("https://127.0.0.1:%d", port)
		} else {
			return fmt.Sprintf("http://127.0.0.1:%d", port)
		}
	} else {
		return address
	}
}

var (
	flagRpcUrl	string

	flagHeight	int64

	flagTxHash	string

	flagAddress		string
	flagTokenAddress	string
	flagTokenName		string

	flagTx	string

	flagName		string
	flagPassword		string
	flagAccessKey		string
	flagEncPrivateKey	string
	flagSmcAddress		string
	flagGasLimit		string
	flagNote		string
	flagNonce		string
	flagTo			string
	flagValue		string
	flagPlainText		string
	flagPageNum		uint64
)

var RootCmd = &cobra.Command{
	Use:	"bcbXwallet",
	Short:	"bcbXwallet",
	Long:	"bcb exchange wallet console",
}

func Execute() error {
	addFlags()
	addCommands()
	return RootCmd.Execute()
}

func addFlags() {
	addWalletCreateFlag()
	addWalletExportFlag()
	addWalletImportFlag()
	addWalletListFlag()
	addTransferFlag()
	addTransferOfflineFlag()

	addBlockHeightFlag()
	addBlockFlag()
	addTransactionFlag()
	addBalanceFlag()
	addBalanceOfTokenFlag()
	addAllBalanceFlag()
	addNonceFlag()
	addCommitTxFlag()
}

func addCommands() {
	RootCmd.AddCommand(walletCreateCmd)
	RootCmd.AddCommand(walletExportCmd)
	RootCmd.AddCommand(walletImportCmd)
	RootCmd.AddCommand(walletListCmd)
	RootCmd.AddCommand(transferCmd)
	RootCmd.AddCommand(transferOfflineCmd)

	RootCmd.AddCommand(blockHeightCmd)
	RootCmd.AddCommand(blockCmd)
	RootCmd.AddCommand(transactionCmd)
	RootCmd.AddCommand(balanceCmd)
	RootCmd.AddCommand(balanceOfTokenCmd)
	RootCmd.AddCommand(allBalanceCmd)
	RootCmd.AddCommand(nonceCmd)
	RootCmd.AddCommand(commitTxCmd)
}

var walletCreateCmd = &cobra.Command{
	Use:	"walletCreate",
	Short:	"create wallet",
	Long:	"create wallet",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.WalletCreate(flagName, flagPassword, flagRpcUrl)
	},
}

func addWalletCreateFlag() {
	walletCreateCmd.PersistentFlags().StringVarP(&flagName, "name", "", "", "wallet name")
	walletCreateCmd.PersistentFlags().StringVarP(&flagPassword, "password", "", "", "wallet password")
	walletCreateCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var walletExportCmd = &cobra.Command{
	Use:	"walletExport",
	Short:	"export wallet",
	Long:	"export wallet",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.WalletExport(flagName, flagPassword, flagAccessKey, flagRpcUrl, flagPlainText)
	},
}

func addWalletExportFlag() {
	walletExportCmd.PersistentFlags().StringVarP(&flagName, "name", "", "", "wallet name")
	walletExportCmd.PersistentFlags().StringVarP(&flagPassword, "password", "", "", "wallet password")
	walletExportCmd.PersistentFlags().StringVarP(&flagAccessKey, "accessKey", "", "", "wallet accessKey")
	walletExportCmd.PersistentFlags().StringVarP(&flagPlainText, "plainText", "", "", "export plain text(default false)")
	walletExportCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var walletImportCmd = &cobra.Command{
	Use:	"walletImport",
	Short:	"import wallet",
	Long:	"import wallet",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.WalletImport(flagName, flagEncPrivateKey, flagPassword, flagAccessKey, flagRpcUrl, flagPlainText)
	},
}

func addWalletImportFlag() {
	walletImportCmd.PersistentFlags().StringVarP(&flagName, "name", "", "", "wallet name")
	walletImportCmd.PersistentFlags().StringVarP(&flagEncPrivateKey, "privateKey", "", "", "wallet privateKey")
	walletImportCmd.PersistentFlags().StringVarP(&flagPassword, "password", "", "", "wallet password")
	walletImportCmd.PersistentFlags().StringVarP(&flagAccessKey, "accessKey", "", "", "wallet accessKey")
	walletImportCmd.PersistentFlags().StringVarP(&flagPlainText, "plainText", "", "", "import plain text(default false)")
	walletImportCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var walletListCmd = &cobra.Command{
	Use:	"walletList",
	Short:	"list wallet",
	Long:	"list wallet",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.WalletList(flagPageNum, flagRpcUrl)
	},
}

func addWalletListFlag() {
	walletListCmd.PersistentFlags().Uint64VarP(&flagPageNum, "pageNum", "", 1, "page index, default first page")
	walletListCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var transferCmd = &cobra.Command{
	Use:	"transfer",
	Short:	"transfer token",
	Long:	"transfer token",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Transfer(flagName, flagAccessKey, flagSmcAddress, flagGasLimit, flagNote, flagTo, flagValue, flagRpcUrl)
	},
}

func addTransferFlag() {
	transferCmd.PersistentFlags().StringVarP(&flagName, "name", "", "", "wallet name")
	transferCmd.PersistentFlags().StringVarP(&flagAccessKey, "accessKey", "", "", "wallet accessKey")
	transferCmd.PersistentFlags().StringVarP(&flagSmcAddress, "smcAddress", "", "", "smart contract address")
	transferCmd.PersistentFlags().StringVarP(&flagGasLimit, "gasLimit", "", "5000", "gas limit ")
	transferCmd.PersistentFlags().StringVarP(&flagNote, "note", "", "", "note")
	transferCmd.PersistentFlags().StringVarP(&flagTo, "to", "", "", "to address")
	transferCmd.PersistentFlags().StringVarP(&flagValue, "value", "", "", "transfer value")
	transferCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var transferOfflineCmd = &cobra.Command{
	Use:	"transferOffline",
	Short:	"offline pack and sign transfer transaction",
	Long:	"offline pack and sign transfer transaction",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.TransferOffline(flagName, flagAccessKey, flagSmcAddress, flagGasLimit, flagNote, flagTo, flagValue, flagNonce, flagRpcUrl)
	},
}

func addTransferOfflineFlag() {
	transferOfflineCmd.PersistentFlags().StringVarP(&flagName, "name", "", "", "wallet name")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagAccessKey, "accessKey", "", "", "wallet accessKey")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagSmcAddress, "smcAddress", "", "", "smart contract address")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagGasLimit, "gasLimit", "", "5000", "gas limit ")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagNonce, "nonce", "", "", "nonce")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagNote, "note", "", "", "note")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagTo, "to", "", "", "to address")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagValue, "value", "", "", "transfer value")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var blockHeightCmd = &cobra.Command{
	Use:	"blockHeight",
	Short:	"get current block height",
	Long:	"get current block height",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.BlockHeight(flagRpcUrl)
	},
}

func addBlockHeightFlag() {
	blockHeightCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var blockCmd = &cobra.Command{
	Use:	"block",
	Short:	"get block info with height",
	Long:	"get block info with height",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Block(flagHeight, flagRpcUrl)
	},
}

func addBlockFlag() {
	blockCmd.PersistentFlags().Int64VarP(&flagHeight, "height", "", 0, "block height")
	blockCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var transactionCmd = &cobra.Command{
	Use:	"transaction",
	Short:	"get transaction info with txHash",
	Long:	"get transaction info with txHash",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Transaction(flagTxHash, flagRpcUrl)
	},
}

func addTransactionFlag() {
	transactionCmd.PersistentFlags().StringVarP(&flagTxHash, "txHash", "", "", "transaction's hash")
	transactionCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var balanceCmd = &cobra.Command{
	Use:	"balance",
	Short:	"get balance of BCB token for specific address",
	Long:	"get balance of BCB token for specific address",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Balance(flagAddress, flagRpcUrl)
	},
}

func addBalanceFlag() {
	balanceCmd.PersistentFlags().StringVarP(&flagAddress, "address", "a", "", "account's address")
	balanceCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var balanceOfTokenCmd = &cobra.Command{
	Use:	"balanceOfToken",
	Short:	"get balance of specific token for specific address",
	Long:	"get balance of specific token for specific address",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.BalanceOfToken(flagAddress, flagTokenAddress, flagTokenName, flagRpcUrl)
	},
}

func addBalanceOfTokenFlag() {
	balanceOfTokenCmd.PersistentFlags().StringVarP(&flagAddress, "address", "", "", "account's address")
	balanceOfTokenCmd.PersistentFlags().StringVarP(&flagTokenAddress, "tokenAddress", "", "", "token's address")
	balanceOfTokenCmd.PersistentFlags().StringVarP(&flagTokenName, "tokenName", "", "", "token's address")
	balanceOfTokenCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var allBalanceCmd = &cobra.Command{
	Use:	"allBalance",
	Short:	"get balance of all tokens for specific address",
	Long:	"get balance of all tokens for specific address",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.AllBalance(flagAddress, flagRpcUrl)
	},
}

func addAllBalanceFlag() {
	allBalanceCmd.PersistentFlags().StringVarP(&flagAddress, "address", "", "", "account's address")
	allBalanceCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var nonceCmd = &cobra.Command{
	Use:	"nonce",
	Short:	"get the next usable nonce for specific address",
	Long:	"get the next usable nonce for specific address",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Nonce(flagAddress, flagRpcUrl)
	},
}

func addNonceFlag() {
	nonceCmd.PersistentFlags().StringVarP(&flagAddress, "address", "", "", "account's address")
	nonceCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var commitTxCmd = &cobra.Command{
	Use:	"commitTx",
	Short:	"commit transaction",
	Long:	"commit transaction",
	Args:	cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.CommitTx(flagTx, flagRpcUrl)
	},
}

func addCommitTxFlag() {
	commitTxCmd.PersistentFlags().StringVarP(&flagTx, "tx", "", "", "packed and signed transaction's data")
	commitTxCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "", serverAddr(common.GetConfig().ServerAddr, true), usage)
}
