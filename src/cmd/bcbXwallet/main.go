package main

import (
	"bcXwallet/client"
	"bcXwallet/common"
	"bcXwallet/rpc"
	rpcserver "common/rpc/lib/server"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	cmn "github.com/tendermint/tmlibs/common"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	usage = "bcbXwallet_rpc's url"
)

func main() {
	err := common.InitAll()
	if err != nil {
		panic(err)
	}

	if version() == "" {
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

func serverAddr(address string, bRequest bool) string {
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

	if bRequest {
		if common.GetConfig().UseHttps {
			return fmt.Sprintf("https://127.0.0.1:%d", port)
		} else {
			return fmt.Sprintf("http://127.0.0.1:%d", port)
		}
	} else {
		return address
	}
}

func version() string {
	params := map[string]interface{}{}
	result := new(rpc.VersionResult)

	err := common.DoHttpRequestAndParseEx([]string{serverAddr(common.GetConfig().ServerAddr, true)}, "bcb_version", params, result)
	if err != nil {
		return ""
	}

	return result.Version
}

// flags
var (
	// global flags
	flagRpcUrl string

	// block flag
	flagHeight int64

	// transaction flag
	flagTxHash string

	// address flag
	flagAddress      string
	flagTokenAddress string
	flagTokenName    string

	// commitTx flag
	flagTx string

	// wallet flag
	flagName          string
	flagPassword      string
	flagAccessKey     string
	flagEncPrivateKey string
	flagSmcAddress    string
	flagGasLimit      string
	flagNote          string
	flagNonce         string
	flagTo            string
	flagValue         string
	flagPlainText     string
	flagPageNum       uint64
)

var RootCmd = &cobra.Command{
	Use:   "bcbXwallet",
	Short: "bcb exchange wallet console",
	Long:  "bcbXwallet client that it can perform the wallet operation, query chain information and so on.",
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
	Use:   "walletCreate",
	Short: "Create wallet",
	Long:  "Create a new wallet",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.WalletCreate(flagName, flagPassword, flagRpcUrl)
	},
}

func addWalletCreateFlag() {
	walletCreateCmd.PersistentFlags().StringVarP(&flagName, "name", "n", "", "wallet name")
	walletCreateCmd.PersistentFlags().StringVarP(&flagPassword, "password", "p", "", "wallet password")
	walletCreateCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var walletExportCmd = &cobra.Command{
	Use:   "walletExport",
	Short: "Export wallet",
	Long:  "Export the private key and walletAddr of wallet",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.WalletExport(flagName, flagPassword, flagAccessKey, flagRpcUrl, flagPlainText)
	},
}

func addWalletExportFlag() {
	walletExportCmd.PersistentFlags().StringVarP(&flagName, "name", "n", "", "wallet name")
	walletExportCmd.PersistentFlags().StringVarP(&flagPassword, "password", "p", "", "wallet password")
	walletExportCmd.PersistentFlags().StringVarP(&flagAccessKey, "accessKey", "a", "", "wallet accessKey")
	walletExportCmd.PersistentFlags().StringVarP(&flagPlainText, "plainText", "t", "", "export plain text(default false)")
	walletExportCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var walletImportCmd = &cobra.Command{
	Use:   "walletImport",
	Short: "Import wallet",
	Long:  "Import the private key to a new wallet",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.WalletImport(flagName, flagEncPrivateKey, flagPassword, flagAccessKey, flagRpcUrl, flagPlainText)
	},
}

func addWalletImportFlag() {
	walletImportCmd.PersistentFlags().StringVarP(&flagName, "name", "n", "", "wallet name")
	walletImportCmd.PersistentFlags().StringVarP(&flagEncPrivateKey, "privateKey", "k", "", "wallet privateKey")
	walletImportCmd.PersistentFlags().StringVarP(&flagPassword, "password", "p", "", "wallet password")
	walletImportCmd.PersistentFlags().StringVarP(&flagAccessKey, "accessKey", "a", "", "wallet accessKey")
	walletImportCmd.PersistentFlags().StringVarP(&flagPlainText, "plainText", "t", "", "import plain text(default false)")
	walletImportCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var walletListCmd = &cobra.Command{
	Use:   "walletList",
	Short: "Wallet list",
	Long:  "Query all wallet names and walletAddrs",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.WalletList(flagPageNum, flagRpcUrl)
	},
}

func addWalletListFlag() {
	walletListCmd.PersistentFlags().Uint64VarP(&flagPageNum, "pageNum", "p", 1, "page index, default first page")
	walletListCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer token",
	Long:  "Transfer token to someone with value",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Transfer(flagName, flagAccessKey, flagSmcAddress, flagGasLimit, flagNote, flagTo, flagValue, flagRpcUrl)
	},
}

func addTransferFlag() {
	transferCmd.PersistentFlags().StringVarP(&flagName, "name", "n", "", "wallet name")
	transferCmd.PersistentFlags().StringVarP(&flagAccessKey, "accessKey", "a", "", "wallet accessKey")
	transferCmd.PersistentFlags().StringVarP(&flagSmcAddress, "smcAddress", "s", "", "smart contract address")
	transferCmd.PersistentFlags().StringVarP(&flagGasLimit, "gasLimit", "g", "5000", "gas limit ")
	transferCmd.PersistentFlags().StringVarP(&flagNote, "note", "o", "", "note")
	transferCmd.PersistentFlags().StringVarP(&flagTo, "to", "t", "", "to address")
	transferCmd.PersistentFlags().StringVarP(&flagValue, "value", "v", "", "transfer value")
	transferCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var transferOfflineCmd = &cobra.Command{
	Use:   "transferOffline",
	Short: "Offline transaction",
	Long:  "Offline pack and sign transfer transaction",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.TransferOffline(flagName, flagAccessKey, flagSmcAddress, flagGasLimit, flagNote, flagTo, flagValue, flagNonce, flagRpcUrl)
	},
}

func addTransferOfflineFlag() {
	transferOfflineCmd.PersistentFlags().StringVarP(&flagName, "name", "n", "", "wallet name")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagAccessKey, "accessKey", "a", "", "wallet accessKey")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagSmcAddress, "smcAddress", "s", "", "smart contract address")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagGasLimit, "gasLimit", "g", "5000", "gas limit ")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagNonce, "nonce", "c", "", "nonce")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagNote, "note", "o", "", "note")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagTo, "to", "t", "", "to address")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagValue, "value", "v", "", "transfer value")
	transferOfflineCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var blockHeightCmd = &cobra.Command{
	Use:   "blockHeight",
	Short: "Get current block height",
	Long:  "Get BlockChain current block height",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.BlockHeight(flagRpcUrl)
	},
}

func addBlockHeightFlag() {
	blockHeightCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "Get block information",
	Long:  "Get block information with height, must great than zero",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Block(flagHeight, flagRpcUrl)
	},
}

func addBlockFlag() {
	blockCmd.PersistentFlags().Int64VarP(&flagHeight, "height", "t", 0, "block height")
	blockCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var transactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "Get transaction information",
	Long:  "Get transaction information with txHash and cannot be empty",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Transaction(flagTxHash, flagRpcUrl)
	},
}

func addTransactionFlag() {
	transactionCmd.PersistentFlags().StringVarP(&flagTxHash, "txHash", "t", "", "transaction's hash")
	transactionCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Get balance information",
	Long:  "Get balance of BCB token for specific address",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Balance(flagAddress, flagRpcUrl)
	},
}

func addBalanceFlag() {
	balanceCmd.PersistentFlags().StringVarP(&flagAddress, "address", "a", "", "account's address")
	balanceCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var balanceOfTokenCmd = &cobra.Command{
	Use:   "balanceOfToken",
	Short: "Get balance information of address",
	Long:  "Get balance of specific token for specific address",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.BalanceOfToken(flagAddress, flagTokenAddress, flagTokenName, flagRpcUrl)
	},
}

func addBalanceOfTokenFlag() {
	balanceOfTokenCmd.PersistentFlags().StringVarP(&flagAddress, "address", "a", "", "account's address")
	balanceOfTokenCmd.PersistentFlags().StringVarP(&flagTokenAddress, "tokenAddress", "t", "", "token's address")
	balanceOfTokenCmd.PersistentFlags().StringVarP(&flagTokenName, "tokenName", "n", "", "token's address")
	balanceOfTokenCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var allBalanceCmd = &cobra.Command{
	Use:   "allBalance",
	Short: "Get all balance information",
	Long:  "Get balance of all tokens for specific address",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.AllBalance(flagAddress, flagRpcUrl)
	},
}

func addAllBalanceFlag() {
	allBalanceCmd.PersistentFlags().StringVarP(&flagAddress, "address", "a", "", "account's address")
	allBalanceCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var nonceCmd = &cobra.Command{
	Use:   "nonce",
	Short: "Get account nonce",
	Long:  "Get the next usable nonce for specific address",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Nonce(flagAddress, flagRpcUrl)
	},
}

func addNonceFlag() {
	nonceCmd.PersistentFlags().StringVarP(&flagAddress, "address", "a", "", "account's address")
	nonceCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}

var commitTxCmd = &cobra.Command{
	Use:   "commitTx",
	Short: "Commit transaction",
	Long:  "Commit transaction with tx's data",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.CommitTx(flagTx, flagRpcUrl)
	},
}

func addCommitTxFlag() {
	commitTxCmd.PersistentFlags().StringVarP(&flagTx, "tx", "t", "", "packed and signed transaction's data")
	commitTxCmd.PersistentFlags().StringVarP(&flagRpcUrl, "url", "u", serverAddr(common.GetConfig().ServerAddr, true), usage)
}
