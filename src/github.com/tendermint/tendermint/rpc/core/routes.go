package core

import (
	rpc "github.com/tendermint/tendermint/rpc/lib/server"
)

var Routes = map[string]*rpc.RPCFunc{

	"health":		rpc.NewRPCFunc(Health, ""),
	"status":		rpc.NewRPCFunc(Status, ""),
	"net_info":		rpc.NewRPCFunc(NetInfo, ""),
	"blockchain":		rpc.NewRPCFunc(BlockchainInfo, "minHeight,maxHeight"),
	"genesis":		rpc.NewRPCFunc(Genesis, ""),
	"block":		rpc.NewRPCFunc(Block, "height"),
	"block_results":	rpc.NewRPCFunc(BlockResults, "height"),
	"commit":		rpc.NewRPCFunc(Commit, "height"),
	"tx":			rpc.NewRPCFunc(Tx, "hash,prove"),

	"validators":		rpc.NewRPCFunc(Validators, "height"),
	"dump_consensus_state":	rpc.NewRPCFunc(DumpConsensusState, ""),
	"unconfirmed_txs":	rpc.NewRPCFunc(UnconfirmedTxs, ""),
	"num_unconfirmed_txs":	rpc.NewRPCFunc(NumUnconfirmedTxs, ""),

	"broadcast_tx_commit":	rpc.NewRPCFunc(BroadcastTxCommit, "tx"),
	"broadcast_tx_sync":	rpc.NewRPCFunc(BroadcastTxSync, "tx"),
	"broadcast_tx_async":	rpc.NewRPCFunc(BroadcastTxAsync, "tx"),

	"abci_query":	rpc.NewRPCFunc(ABCIQuery, "path,data,height,trusted"),
	"abci_info":	rpc.NewRPCFunc(ABCIInfo, ""),
}

func AddUnsafeRoutes() {

	Routes["dial_seeds"] = rpc.NewRPCFunc(UnsafeDialSeeds, "seeds")
	Routes["dial_peers"] = rpc.NewRPCFunc(UnsafeDialPeers, "peers,persistent")
	Routes["unsafe_flush_mempool"] = rpc.NewRPCFunc(UnsafeFlushMempool, "")

	Routes["unsafe_start_cpu_profiler"] = rpc.NewRPCFunc(UnsafeStartCPUProfiler, "filename")
	Routes["unsafe_stop_cpu_profiler"] = rpc.NewRPCFunc(UnsafeStopCPUProfiler, "")
	Routes["unsafe_write_heap_profile"] = rpc.NewRPCFunc(UnsafeWriteHeapProfile, "filename")
}
