package rpc

import (
	rpcserver "common/rpc/lib/server"
)

var Routes = map[string]*rpcserver.RPCFunc{
	// bcbXWallet api
	"bcb_walletCreate":    rpcserver.NewRPCFunc(WalletCreate, "name,password"),
	"bcb_walletExport":    rpcserver.NewRPCFunc(WalletExport, "name,password,accessKey,plainText"),
	"bcb_walletImport":    rpcserver.NewRPCFunc(WalletImport, "name,privateKey,password,accessKey,plainText"),
	"bcb_walletList":      rpcserver.NewRPCFunc(WalletList, "pageNum"),
	"bcb_transfer":        rpcserver.NewRPCFunc(WalletTransfer, "name,accessKey,walletParams"),
	"bcb_transferOffline": rpcserver.NewRPCFunc(WalletTransferOffline, "name,accessKey,walletParams"),

	// block chain api
	"bcb_blockHeight":    rpcserver.NewRPCFunc(BlockHeight, ""),
	"bcb_block":          rpcserver.NewRPCFunc(Block, "height"),
	"bcb_transaction":    rpcserver.NewRPCFunc(Transaction, "txHash"),
	"bcb_balance":        rpcserver.NewRPCFunc(Balance, "address"),
	"bcb_balanceOfToken": rpcserver.NewRPCFunc(BalanceOfToken, "address,tokenAddress,tokenName"),
	"bcb_allBalance":     rpcserver.NewRPCFunc(AllBalance, "address"),
	"bcb_nonce":          rpcserver.NewRPCFunc(Nonce, "address"),
	"bcb_commitTx":       rpcserver.NewRPCFunc(CommitTx, "tx"),
	"bcb_version":        rpcserver.NewRPCFunc(Version, ""),
}
