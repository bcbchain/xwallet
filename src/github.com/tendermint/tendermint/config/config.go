package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	DefaultTendermintDir	= ".tendermint"
	defaultConfigDir	= "config"
	defaultDataDir		= "data"

	defaultConfigFileName	= "config.toml"
	defaultGenesisJSONName	= "genesis.json"

	defaultPrivValName	= "priv_validator.json"
	defaultNodeKeyName	= "node_key.json"
	defaultValsFileName	= "validators.json"
	defaultAddrBookName	= "addrbook.json"
	defaultAppProxy		= []string{"tcp://127.0.0.1:46658"}

	defaultConfigFilePath	= defaultConfigDir + "/" + defaultConfigFileName
	defaultGenesisJSONPath	= defaultConfigDir + "/" + defaultGenesisJSONName
	defaultPrivValPath	= defaultConfigDir + "/" + defaultPrivValName
	defaultNodeKeyPath	= defaultConfigDir + "/" + defaultNodeKeyName
	defaultValsPath		= defaultConfigDir + "/" + defaultValsFileName
	defaultAddrBookPath	= defaultConfigDir + "/" + defaultAddrBookName
)

type Config struct {
	BaseConfig	`mapstructure:",squash"`

	RPC		*RPCConfig		`mapstructure:"rpc"`
	P2P		*P2PConfig		`mapstructure:"p2p"`
	Mempool		*MempoolConfig		`mapstructure:"mempool"`
	Consensus	*ConsensusConfig	`mapstructure:"consensus"`
	TxIndex		*TxIndexConfig		`mapstructure:"tx_index"`
}

func DefaultConfig() *Config {
	return &Config{
		BaseConfig:	DefaultBaseConfig(),
		RPC:		DefaultRPCConfig(),
		P2P:		DefaultP2PConfig(),
		Mempool:	DefaultMempoolConfig(),
		Consensus:	DefaultConsensusConfig(),
		TxIndex:	DefaultTxIndexConfig(),
	}
}

func TestConfig() *Config {
	return &Config{
		BaseConfig:	TestBaseConfig(),
		RPC:		TestRPCConfig(),
		P2P:		TestP2PConfig(),
		Mempool:	TestMempoolConfig(),
		Consensus:	TestConsensusConfig(),
		TxIndex:	TestTxIndexConfig(),
	}
}

func (cfg *Config) SetRoot(root string) *Config {
	cfg.BaseConfig.RootDir = root
	cfg.RPC.RootDir = root
	cfg.P2P.RootDir = root
	cfg.Mempool.RootDir = root
	cfg.Consensus.RootDir = root
	return cfg
}

type BaseConfig struct {
	chainID	string

	RootDir	string	`mapstructure:"home"`

	Genesis	string	`mapstructure:"genesis_file"`

	PrivValidator	string	`mapstructure:"priv_validator_file"`

	NodeKey	string	`mapstructure:"node_key_file"`

	Validators	string	`mapstructure:"validators_file"`

	Moniker	string	`mapstructure:"moniker"`

	PrivValidatorListenAddr	string	`mapstructure:"priv_validator_laddr"`

	ProxyApp	[]string	`mapstructure:"proxy_app"`

	ABCI	string	`mapstructure:"abci"`

	Persist	string	`mapstructure:"persist"`

	LogLevel	string	`mapstructure:"log_level"`

	LogPath	string	`mapstructure:"log_path"`
	LogFile	string	`mapstructure:"log_file"`

	ProfListenAddress	string	`mapstructure:"prof_laddr"`

	FastSync	bool	`mapstructure:"fast_sync"`

	FilterPeers	bool	`mapstructure:"filter_peers"`

	DBBackend	string	`mapstructure:"db_backend"`

	DBPath	string	`mapstructure:"db_path"`
}

func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		Genesis:		defaultGenesisJSONPath,
		PrivValidator:		defaultPrivValPath,
		NodeKey:		defaultNodeKeyPath,
		Validators:		defaultValsPath,
		Moniker:		defaultMoniker,
		ProxyApp:		defaultAppProxy,
		ABCI:			"socket",
		LogLevel:		DefaultPackageLogLevels(),
		LogPath:		"log",
		LogFile:		"",
		ProfListenAddress:	"",
		FastSync:		true,
		FilterPeers:		false,
		DBBackend:		"leveldb",
		DBPath:			"data",
	}
}

func TestBaseConfig() BaseConfig {
	cfg := DefaultBaseConfig()
	cfg.chainID = "tendermint_test"
	cfg.ProxyApp = defaultAppProxy
	cfg.FastSync = false
	cfg.DBBackend = "memdb"
	return cfg
}

func (cfg BaseConfig) ChainID() string {
	return cfg.chainID
}

func (cfg BaseConfig) GenesisFile() string {
	return rootify(cfg.Genesis, cfg.RootDir)
}

func (cfg BaseConfig) ConfigFilePath() string {
	return rootify(defaultConfigFilePath, cfg.RootDir)
}

func (cfg BaseConfig) PrivValidatorFile() string {
	return rootify(cfg.PrivValidator, cfg.RootDir)
}

func (cfg BaseConfig) NodeKeyFile() string {
	return rootify(cfg.NodeKey, cfg.RootDir)
}

func (cfg BaseConfig) ValidatorsFile() string {
	return rootify(cfg.Validators, cfg.RootDir)
}

func (cfg BaseConfig) DBDir() string {
	return rootify(cfg.DBPath, cfg.RootDir)
}

func DefaultLogLevel() string {
	return "error"
}

func (cfg BaseConfig) LogDir() string {
	return rootify(cfg.LogPath, cfg.RootDir)
}

func DefaultPackageLogLevels() string {
	return fmt.Sprintf("main:info,state:info,*:%s", DefaultLogLevel())
}

type RPCConfig struct {
	RootDir	string	`mapstructure:"home"`

	ListenAddress	string	`mapstructure:"laddr"`

	CertFile	string	`mapstructure:"cert_file"`
	KeyFile		string	`mapstructure:"key_file"`

	GRPCListenAddress	string	`mapstructure:"grpc_laddr"`

	Unsafe	bool	`mapstructure:"unsafe"`
}

func DefaultRPCConfig() *RPCConfig {
	return &RPCConfig{
		ListenAddress:		"tcp://0.0.0.0:46657",
		GRPCListenAddress:	"",
		Unsafe:			false,
		CertFile:		"STAR.bcbchain.io.crt",
		KeyFile:		"STAR.bcbchain.io.key",
	}
}

func TestRPCConfig() *RPCConfig {
	cfg := DefaultRPCConfig()
	cfg.ListenAddress = "tcp://0.0.0.0:36657"
	cfg.GRPCListenAddress = "tcp://0.0.0.0:36658"
	cfg.Unsafe = true
	return cfg
}

type P2PConfig struct {
	RootDir	string	`mapstructure:"home"`

	ListenAddress	string	`mapstructure:"laddr"`

	AAddress	string	`mapstructure:"aaddr"`

	Seeds	string	`mapstructure:"seeds"`

	PersistentPeers	string	`mapstructure:"persistent_peers"`

	SkipUPNP	bool	`mapstructure:"skip_upnp"`

	AddrBook	string	`mapstructure:"addr_book_file"`

	AddrBookStrict	bool	`mapstructure:"addr_book_strict"`

	MaxNumPeers	int	`mapstructure:"max_num_peers"`

	FlushThrottleTimeout	int	`mapstructure:"flush_throttle_timeout"`

	MaxPacketMsgPayloadSize	int	`mapstructure:"max_packet_msg_payload_size"`

	SendRate	int64	`mapstructure:"send_rate"`

	RecvRate	int64	`mapstructure:"recv_rate"`

	PexReactor	bool	`mapstructure:"pex"`

	SeedMode	bool	`mapstructure:"seed_mode"`

	AuthEnc	bool	`mapstructure:"auth_enc"`

	PrivatePeerIDs	string	`mapstructure:"private_peer_ids"`
}

func DefaultP2PConfig() *P2PConfig {
	return &P2PConfig{
		ListenAddress:			"tcp://0.0.0.0:46656",
		AddrBook:			defaultAddrBookPath,
		AddrBookStrict:			true,
		MaxNumPeers:			50,
		FlushThrottleTimeout:		100,
		MaxPacketMsgPayloadSize:	1024,
		SendRate:			512000,
		RecvRate:			512000,
		PexReactor:			true,
		SeedMode:			false,
		AuthEnc:			true,
	}
}

func TestP2PConfig() *P2PConfig {
	cfg := DefaultP2PConfig()
	cfg.ListenAddress = "tcp://0.0.0.0:36656"
	cfg.SkipUPNP = true
	cfg.FlushThrottleTimeout = 10
	return cfg
}

func (cfg *P2PConfig) AddrBookFile() string {
	return rootify(cfg.AddrBook, cfg.RootDir)
}

type MempoolConfig struct {
	RootDir		string	`mapstructure:"home"`
	Recheck		bool	`mapstructure:"recheck"`
	RecheckEmpty	bool	`mapstructure:"recheck_empty"`
	Broadcast	bool	`mapstructure:"broadcast"`
	WalPath		string	`mapstructure:"wal_dir"`
	CacheSize	int	`mapstructure:"cache_size"`
	CTxCacheTime	int64	`mapstructure:"ctx_cache_time"`
}

func DefaultMempoolConfig() *MempoolConfig {
	return &MempoolConfig{
		Recheck:	true,
		RecheckEmpty:	true,
		Broadcast:	true,
		WalPath:	defaultDataDir + "/" + "mempool.wal",
		CTxCacheTime:	600,
		CacheSize:	100000,
	}
}

func TestMempoolConfig() *MempoolConfig {
	cfg := DefaultMempoolConfig()
	cfg.CacheSize = 1000
	cfg.CTxCacheTime = 600
	return cfg
}

func (cfg *MempoolConfig) WalDir() string {
	return rootify(cfg.WalPath, cfg.RootDir)
}

type ConsensusConfig struct {
	RootDir	string	`mapstructure:"home"`
	WalPath	string	`mapstructure:"wal_file"`
	walFile	string

	TimeoutPropose		int	`mapstructure:"timeout_propose"`
	TimeoutProposeDelta	int	`mapstructure:"timeout_propose_delta"`
	TimeoutPrevote		int	`mapstructure:"timeout_prevote"`
	TimeoutPrevoteDelta	int	`mapstructure:"timeout_prevote_delta"`
	TimeoutPrecommit	int	`mapstructure:"timeout_precommit"`
	TimeoutPrecommitDelta	int	`mapstructure:"timeout_precommit_delta"`
	TimeoutCommit		int	`mapstructure:"timeout_commit"`

	SkipTimeoutCommit	bool	`mapstructure:"skip_timeout_commit"`

	MaxBlockSizeTxs		int	`mapstructure:"max_block_size_txs"`
	MaxBlockSizeBytes	int	`mapstructure:"max_block_size_bytes"`

	CreateEmptyBlocks		bool	`mapstructure:"create_empty_blocks"`
	CreateEmptyBlocksInterval	int	`mapstructure:"create_empty_blocks_interval"`

	PeerGossipSleepDuration		int	`mapstructure:"peer_gossip_sleep_duration"`
	PeerQueryMaj23SleepDuration	int	`mapstructure:"peer_query_maj23_sleep_duration"`
}

func DefaultConsensusConfig() *ConsensusConfig {
	return &ConsensusConfig{
		WalPath:			defaultDataDir + "/" + "cs.wal" + "/" + "wal",
		TimeoutPropose:			3000,
		TimeoutProposeDelta:		500,
		TimeoutPrevote:			1000,
		TimeoutPrevoteDelta:		500,
		TimeoutPrecommit:		1000,
		TimeoutPrecommitDelta:		500,
		TimeoutCommit:			1000,
		SkipTimeoutCommit:		false,
		MaxBlockSizeTxs:		10000,
		MaxBlockSizeBytes:		1,
		CreateEmptyBlocks:		true,
		CreateEmptyBlocksInterval:	120,
		PeerGossipSleepDuration:	100,
		PeerQueryMaj23SleepDuration:	2000,
	}
}

func TestConsensusConfig() *ConsensusConfig {
	cfg := DefaultConsensusConfig()
	cfg.TimeoutPropose = 100
	cfg.TimeoutProposeDelta = 1
	cfg.TimeoutPrevote = 10
	cfg.TimeoutPrevoteDelta = 1
	cfg.TimeoutPrecommit = 10
	cfg.TimeoutPrecommitDelta = 1
	cfg.TimeoutCommit = 10
	cfg.SkipTimeoutCommit = true
	cfg.PeerGossipSleepDuration = 5
	cfg.PeerQueryMaj23SleepDuration = 250
	return cfg
}

func (cfg *ConsensusConfig) WaitForTxs() bool {
	return !cfg.CreateEmptyBlocks || cfg.CreateEmptyBlocksInterval > 0
}

func (cfg *ConsensusConfig) EmptyBlocksInterval() time.Duration {
	return time.Duration(cfg.CreateEmptyBlocksInterval) * time.Second
}

func (cfg *ConsensusConfig) Propose(round int) time.Duration {
	return time.Duration(cfg.TimeoutPropose+cfg.TimeoutProposeDelta*round) * time.Millisecond
}

func (cfg *ConsensusConfig) Prevote(round int) time.Duration {
	return time.Duration(cfg.TimeoutPrevote+cfg.TimeoutPrevoteDelta*round) * time.Millisecond
}

func (cfg *ConsensusConfig) Precommit(round int) time.Duration {
	return time.Duration(cfg.TimeoutPrecommit+cfg.TimeoutPrecommitDelta*round) * time.Millisecond
}

func (cfg *ConsensusConfig) Commit(t time.Time) time.Time {
	return t.Add(time.Duration(cfg.TimeoutCommit) * time.Millisecond)
}

func (cfg *ConsensusConfig) PeerGossipSleep() time.Duration {
	return time.Duration(cfg.PeerGossipSleepDuration) * time.Millisecond
}

func (cfg *ConsensusConfig) PeerQueryMaj23Sleep() time.Duration {
	return time.Duration(cfg.PeerQueryMaj23SleepDuration) * time.Millisecond
}

func (cfg *ConsensusConfig) WalFile() string {
	if cfg.walFile != "" {
		return cfg.walFile
	}
	return rootify(cfg.WalPath, cfg.RootDir)
}

func (cfg *ConsensusConfig) SetWalFile(walFile string) {
	cfg.walFile = walFile
}

type TxIndexConfig struct {
	Indexer	string	`mapstructure:"indexer"`

	IndexTags	string	`mapstructure:"index_tags"`

	IndexAllTags	bool	`mapstructure:"index_all_tags"`
}

func DefaultTxIndexConfig() *TxIndexConfig {
	return &TxIndexConfig{
		Indexer:	"kv",
		IndexTags:	"",
		IndexAllTags:	false,
	}
}

func TestTxIndexConfig() *TxIndexConfig {
	return DefaultTxIndexConfig()
}

func rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return root + "/" + path
}

var defaultMoniker = getDefaultMoniker()

func getDefaultMoniker() string {
	moniker, err := os.Hostname()
	if err != nil {
		moniker = "anonymous"
	}
	return moniker
}
