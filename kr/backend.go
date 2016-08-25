// Copyright 2014 The go-krypton Authors
// This file is part of the go-krypton library.
//
// The go-krypton library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-krypton library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-krypton library. If not, see <http://www.gnu.org/licenses/>.

// Package kr implements the Krypton protocol.
package kr

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/krypton/krash"
	"github.com/krypton/go-krypton/accounts"
	"github.com/krypton/go-krypton/common"
	"github.com/krypton/go-krypton/common/compiler"
	"github.com/krypton/go-krypton/common/httpclient"
	"github.com/krypton/go-krypton/core"
	"github.com/krypton/go-krypton/core/state"
	"github.com/krypton/go-krypton/core/types"
	"github.com/krypton/go-krypton/core/vm"
	"github.com/krypton/go-krypton/crypto"
	"github.com/krypton/go-krypton/kr/downloader"
	"github.com/krypton/go-krypton/krdb"
	"github.com/krypton/go-krypton/event"
	"github.com/krypton/go-krypton/logger"
	"github.com/krypton/go-krypton/logger/glog"
	"github.com/krypton/go-krypton/miner"
	"github.com/krypton/go-krypton/p2p"
	"github.com/krypton/go-krypton/p2p/discover"
	"github.com/krypton/go-krypton/p2p/nat"
	"github.com/krypton/go-krypton/rlp"
	"github.com/krypton/go-krypton/whisper"
)

const (
	epochLength    = 30000
	krashRevision = 23

	autoDAGcheckInterval = 10 * time.Hour
	autoDAGepochHeight   = epochLength / 2
)

var (
	jsonlogger = logger.NewJsonLogger()

	datadirInUseErrnos = map[uint]bool{11: true, 32: true, 35: true}
	portInUseErrRE     = regexp.MustCompile("address already in use")

	defaultBootNodes = []*discover.Node{
		// KR/DEV Go Bootnodes
		discover.MustParseNode("enode://0fb920467bdc123eeee3c132519457967227ed80d255fc210b09bc070fb4612df90532604629f91c3f1912c298d5a20ef4c22824d1be33fef270b07448760362@192.99.32.187:17171"),
		discover.MustParseNode("enode://382b74a663581fe3559b52521415f9e7ee596f4873e0e88f8406e0e025d6ed1f2bed55f1f9739cab7d1e247cd5bcc4d944add0c282ea7232155da74b895f0e75@167.114.13.142:17172"),
		discover.MustParseNode("enode://809e1781a9da785ee6084ff38011f771c61ba9dce1e372ca8cbf2de9af119d9df5ded7728bc03f0b765b8ec9d14a3b296e6b248239d6cc8e4b2fbc3d3f9796a1@167.114.13.143:17173"),
	}

	defaultTestNetBootNodes = []*discover.Node{
		//discover.MustParseNode("enode://e4533109cc9bd7604e4ff6c095f7a1d807e15b38e9bfeb05d3b7c423ba86af0a9e89abbf40bd9dde4250fef114cd09270fa4e224cbeef8b7bf05a51e8260d6b8@94.242.229.4:40404"),
		//discover.MustParseNode("enode://8c336ee6f03e99613ad21274f269479bf4413fb294d697ef15ab897598afb931f56beb8e97af530aee20ce2bcba5776f4a312bc168545de4d43736992c814592@94.242.229.203:30303"),
	}

	staticNodes  = "static-nodes.json"  // Path within <datadir> to search for the static node list
	trustedNodes = "trusted-nodes.json" // Path within <datadir> to search for the trusted node list
)

type Config struct {
	DevMode bool
	TestNet bool

	Name         string
	NetworkId    int
	GenesisFile  string
	GenesisBlock *types.Block // used by block tests
	FastSync     bool
	Olympic      bool

	BlockChainVersion  int
	SkipBcVersionCheck bool // e.g. blockchain export
	DatabaseCache      int

	DataDir   string
	LogFile   string
	Verbosity int
	VmDebug   bool
	NatSpec   bool
	DocRoot   string
	AutoDAG   bool
	PowTest   bool
	ExtraData []byte

	MaxPeers        int
	MaxPendingPeers int
	Discovery       bool
	Port            string

	// Space-separated list of discovery node URLs
	BootNodes string

	// This key is used to identify the node on the network.
	// If nil, an ephemeral key is used.
	NodeKey *ecdsa.PrivateKey

	NAT  nat.Interface
	Shh  bool
	Dial bool

	Krbase      common.Address
	GasPrice       *big.Int
	MinerThreads   int
	AccountManager *accounts.Manager
	SolcPath       string

	GpoMinGasPrice          *big.Int
	GpoMaxGasPrice          *big.Int
	GpoFullBlockRatio       int
	GpobaseStepDown         int
	GpobaseStepUp           int
	GpobaseCorrectionFactor int

	// NewDB is used to create databases.
	// If nil, the default is to create leveldb databases on disk.
	NewDB func(path string) (krdb.Database, error)
}

func (cfg *Config) parseBootNodes() []*discover.Node {
	if cfg.BootNodes == "" {
		if cfg.TestNet {
			return defaultTestNetBootNodes
		}

		return defaultBootNodes
	}
	var ns []*discover.Node
	for _, url := range strings.Split(cfg.BootNodes, " ") {
		if url == "" {
			continue
		}
		n, err := discover.ParseNode(url)
		if err != nil {
			glog.V(logger.Error).Infof("Bootstrap URL %s: %v\n", url, err)
			continue
		}
		ns = append(ns, n)
	}
	return ns
}

// parseNodes parses a list of discovery node URLs loaded from a .json file.
func (cfg *Config) parseNodes(file string) []*discover.Node {
	// Short circuit if no node config is present
	path := filepath.Join(cfg.DataDir, file)
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	// Load the nodes from the config file
	blob, err := ioutil.ReadFile(path)
	if err != nil {
		glog.V(logger.Error).Infof("Failed to access nodes: %v", err)
		return nil
	}
	nodelist := []string{}
	if err := json.Unmarshal(blob, &nodelist); err != nil {
		glog.V(logger.Error).Infof("Failed to load nodes: %v", err)
		return nil
	}
	// Interpret the list as a discovery node array
	var nodes []*discover.Node
	for _, url := range nodelist {
		if url == "" {
			continue
		}
		node, err := discover.ParseNode(url)
		if err != nil {
			glog.V(logger.Error).Infof("Node URL %s: %v\n", url, err)
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func (cfg *Config) nodeKey() (*ecdsa.PrivateKey, error) {
	// use explicit key from command line args if set
	if cfg.NodeKey != nil {
		return cfg.NodeKey, nil
	}
	// use persistent key if present
	keyfile := filepath.Join(cfg.DataDir, "nodekey")
	key, err := crypto.LoadECDSA(keyfile)
	if err == nil {
		return key, nil
	}
	// no persistent key, generate and store a new one
	if key, err = crypto.GenerateKey(); err != nil {
		return nil, fmt.Errorf("could not generate server key: %v", err)
	}
	if err := crypto.SaveECDSA(keyfile, key); err != nil {
		glog.V(logger.Error).Infoln("could not persist nodekey: ", err)
	}
	return key, nil
}

type Krypton struct {
	// Channel for shutting down the krypton
	shutdownChan chan bool

	// DB interfaces
	chainDb krdb.Database // Block chain database
	dappDb  krdb.Database // Dapp database

	// Handlers
	txPool          *core.TxPool
	blockchain      *core.BlockChain
	accountManager  *accounts.Manager
	whisper         *whisper.Whisper
	pow             *krash.Krash
	protocolManager *ProtocolManager
	SolcPath        string
	solc            *compiler.Solidity

	GpoMinGasPrice          *big.Int
	GpoMaxGasPrice          *big.Int
	GpoFullBlockRatio       int
	GpobaseStepDown         int
	GpobaseStepUp           int
	GpobaseCorrectionFactor int

	httpclient *httpclient.HTTPClient

	net      *p2p.Server
	eventMux *event.TypeMux
	miner    *miner.Miner

	// logger logger.LogSystem

	Mining        bool
	MinerThreads  int
	NatSpec       bool
	DataDir       string
	AutoDAG       bool
	PowTest       bool
	autodagquit   chan bool
	krbase     common.Address
	clientVersion string
	netVersionId  int
	shhVersionId  int
}

func New(config *Config) (*Krypton, error) {
	logger.New(config.DataDir, config.LogFile, config.Verbosity)

	// Let the database take 3/4 of the max open files (TODO figure out a way to get the actual limit of the open files)
	const dbCount = 3
	krdb.OpenFileLimit = 128 / (dbCount + 1)

	newdb := config.NewDB
	if newdb == nil {
		newdb = func(path string) (krdb.Database, error) { return krdb.NewLDBDatabase(path, config.DatabaseCache) }
	}

	// Open the chain database and perform any upgrades needed
	chainDb, err := newdb(filepath.Join(config.DataDir, "chaindata"))
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok && datadirInUseErrnos[uint(errno)] {
			err = fmt.Errorf("%v (check if another instance of gkr is already running with the same data directory '%s')", err, config.DataDir)
		}
		return nil, fmt.Errorf("blockchain db err: %v", err)
	}
	if db, ok := chainDb.(*krdb.LDBDatabase); ok {
		db.Meter("kr/db/chaindata/")
	}
	if err := upgradeChainDatabase(chainDb); err != nil {
		return nil, err
	}
	if err := addMipmapBloomBins(chainDb); err != nil {
		return nil, err
	}

	dappDb, err := newdb(filepath.Join(config.DataDir, "dapp"))
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok && datadirInUseErrnos[uint(errno)] {
			err = fmt.Errorf("%v (check if another instance of gkr is already running with the same data directory '%s')", err, config.DataDir)
		}
		return nil, fmt.Errorf("dapp db err: %v", err)
	}
	if db, ok := dappDb.(*krdb.LDBDatabase); ok {
		db.Meter("kr/db/dapp/")
	}

	nodeDb := filepath.Join(config.DataDir, "nodes")
	glog.V(logger.Info).Infof("Protocol Versions: %v, Network Id: %v", ProtocolVersions, config.NetworkId)

	if len(config.GenesisFile) > 0 {
		fr, err := os.Open(config.GenesisFile)
		if err != nil {
			return nil, err
		}

		block, err := core.WriteGenesisBlock(chainDb, fr)
		if err != nil {
			return nil, err
		}
		glog.V(logger.Info).Infof("Successfully wrote genesis block. New genesis hash = %x\n", block.Hash())
	}

	// different modes
	switch {
	case config.Olympic:
		glog.V(logger.Error).Infoln("Starting Olympic network")
		fallthrough
	case config.DevMode:
		_, err := core.WriteOlympicGenesisBlock(chainDb, 42)
		if err != nil {
			return nil, err
		}
	case config.TestNet:
		state.StartingNonce = 1048576 // (2**20)
		_, err := core.WriteTestNetGenesisBlock(chainDb, 0x6d6f7264656e)
		if err != nil {
			return nil, err
		}
	}
	// This is for testing only.
	if config.GenesisBlock != nil {
		core.WriteTd(chainDb, config.GenesisBlock.Hash(), config.GenesisBlock.Difficulty())
		core.WriteBlock(chainDb, config.GenesisBlock)
		core.WriteCanonicalHash(chainDb, config.GenesisBlock.Hash(), config.GenesisBlock.NumberU64())
		core.WriteHeadBlockHash(chainDb, config.GenesisBlock.Hash())
	}

	if !config.SkipBcVersionCheck {
		b, _ := chainDb.Get([]byte("BlockchainVersion"))
		bcVersion := int(common.NewValue(b).Uint())
		if bcVersion != config.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run gkr upgradedb.\n", bcVersion, config.BlockChainVersion)
		}
		saveBlockchainVersion(chainDb, config.BlockChainVersion)
	}
	glog.V(logger.Info).Infof("Blockchain DB Version: %d", config.BlockChainVersion)

	kr := &Krypton{
		shutdownChan:            make(chan bool),
		chainDb:                 chainDb,
		dappDb:                  dappDb,
		eventMux:                &event.TypeMux{},
		accountManager:          config.AccountManager,
		DataDir:                 config.DataDir,
		krbase:               config.Krbase,
		clientVersion:           config.Name, // TODO should separate from Name
		netVersionId:            config.NetworkId,
		NatSpec:                 config.NatSpec,
		MinerThreads:            config.MinerThreads,
		SolcPath:                config.SolcPath,
		AutoDAG:                 config.AutoDAG,
		PowTest:                 config.PowTest,
		GpoMinGasPrice:          config.GpoMinGasPrice,
		GpoMaxGasPrice:          config.GpoMaxGasPrice,
		GpoFullBlockRatio:       config.GpoFullBlockRatio,
		GpobaseStepDown:         config.GpobaseStepDown,
		GpobaseStepUp:           config.GpobaseStepUp,
		GpobaseCorrectionFactor: config.GpobaseCorrectionFactor,
		httpclient:              httpclient.New(config.DocRoot),
	}

	if config.PowTest {
		glog.V(logger.Info).Infof("krash used in test mode")
		kr.pow, err = krash.NewForTesting()
		if err != nil {
			return nil, err
		}
	} else {
		kr.pow = krash.New()
	}
	//genesis := core.GenesisBlock(uint64(config.GenesisNonce), stateDb)
	kr.blockchain, err = core.NewBlockChain(chainDb, kr.pow, kr.EventMux())
	if err != nil {
		if err == core.ErrNoGenesis {
			return nil, fmt.Errorf(`Genesis block not found. Please supply a genesis block with the "--genesis /path/to/file" argument`)
		}
		return nil, err
	}
	newPool := core.NewTxPool(kr.EventMux(), kr.blockchain.State, kr.blockchain.GasLimit)
	kr.txPool = newPool

	if kr.protocolManager, err = NewProtocolManager(config.FastSync, config.NetworkId, kr.eventMux, kr.txPool, kr.pow, kr.blockchain, chainDb); err != nil {
		return nil, err
	}
	kr.miner = miner.New(kr, kr.EventMux(), kr.pow)
	kr.miner.SetGasPrice(config.GasPrice)
	kr.miner.SetExtra(config.ExtraData)

	if config.Shh {
		kr.whisper = whisper.New()
		kr.shhVersionId = int(kr.whisper.Version())
	}

	netprv, err := config.nodeKey()
	if err != nil {
		return nil, err
	}
	protocols := append([]p2p.Protocol{}, kr.protocolManager.SubProtocols...)
	if config.Shh {
		protocols = append(protocols, kr.whisper.Protocol())
	}
	kr.net = &p2p.Server{
		PrivateKey:      netprv,
		Name:            config.Name,
		MaxPeers:        config.MaxPeers,
		MaxPendingPeers: config.MaxPendingPeers,
		Discovery:       config.Discovery,
		Protocols:       protocols,
		NAT:             config.NAT,
		NoDial:          !config.Dial,
		BootstrapNodes:  config.parseBootNodes(),
		StaticNodes:     config.parseNodes(staticNodes),
		TrustedNodes:    config.parseNodes(trustedNodes),
		NodeDatabase:    nodeDb,
	}
	if len(config.Port) > 0 {
		kr.net.ListenAddr = ":" + config.Port
	}

	vm.Debug = config.VmDebug

	return kr, nil
}

// Network retrieves the underlying P2P network server. This should eventually
// be moved out into a protocol independent package, but for now use an accessor.
func (s *Krypton) Network() *p2p.Server {
	return s.net
}

func (s *Krypton) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Krypton) Krbase() (eb common.Address, err error) {
	eb = s.krbase
	if (eb == common.Address{}) {
		addr, e := s.AccountManager().AddressByIndex(0)
		if e != nil {
			err = fmt.Errorf("krbase address must be explicitly specified")
		}
		eb = common.HexToAddress(addr)
	}
	return
}

// set in js console via admin interface or wrapper from cli flags
func (self *Krypton) SetKrbase(krbase common.Address) {
	self.krbase = krbase
	self.miner.SetKrbase(krbase)
}

func (s *Krypton) StopMining()         { s.miner.Stop() }
func (s *Krypton) IsMining() bool      { return s.miner.Mining() }
func (s *Krypton) Miner() *miner.Miner { return s.miner }

// func (s *Krypton) Logger() logger.LogSystem             { return s.logger }
func (s *Krypton) Name() string                       { return s.net.Name }
func (s *Krypton) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *Krypton) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *Krypton) TxPool() *core.TxPool               { return s.txPool }
func (s *Krypton) Whisper() *whisper.Whisper          { return s.whisper }
func (s *Krypton) EventMux() *event.TypeMux           { return s.eventMux }
func (s *Krypton) ChainDb() krdb.Database            { return s.chainDb }
func (s *Krypton) DappDb() krdb.Database             { return s.dappDb }
func (s *Krypton) IsListening() bool                  { return true } // Always listening
func (s *Krypton) PeerCount() int                     { return s.net.PeerCount() }
func (s *Krypton) Peers() []*p2p.Peer                 { return s.net.Peers() }
func (s *Krypton) MaxPeers() int                      { return s.net.MaxPeers }
func (s *Krypton) ClientVersion() string              { return s.clientVersion }
func (s *Krypton) KrVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *Krypton) NetVersion() int                    { return s.netVersionId }
func (s *Krypton) ShhVersion() int                    { return s.shhVersionId }
func (s *Krypton) Downloader() *downloader.Downloader { return s.protocolManager.downloader }

// Start the krypton
func (s *Krypton) Start() error {
	jsonlogger.LogJson(&logger.LogStarting{
		ClientString:    s.net.Name,
		ProtocolVersion: s.KrVersion(),
	})
	err := s.net.Start()
	if err != nil {
		if portInUseErrRE.MatchString(err.Error()) {
			err = fmt.Errorf("%v (possibly another instance of gkr is using the same port)", err)
		}
		return err
	}

	if s.AutoDAG {
		s.StartAutoDAG()
	}

	s.protocolManager.Start()

	if s.whisper != nil {
		s.whisper.Start()
	}

	glog.V(logger.Info).Infoln("Server started")
	return nil
}

func (s *Krypton) StartForTest() {
	jsonlogger.LogJson(&logger.LogStarting{
		ClientString:    s.net.Name,
		ProtocolVersion: s.KrVersion(),
	})
}

// AddPeer connects to the given node and maintains the connection until the
// server is shut down. If the connection fails for any reason, the server will
// attempt to reconnect the peer.
func (self *Krypton) AddPeer(nodeURL string) error {
	n, err := discover.ParseNode(nodeURL)
	if err != nil {
		return fmt.Errorf("invalid node URL: %v", err)
	}
	self.net.AddPeer(n)
	return nil
}

func (s *Krypton) Stop() {
	s.net.Stop()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	s.txPool.Stop()
	s.eventMux.Stop()
	if s.whisper != nil {
		s.whisper.Stop()
	}
	s.StopAutoDAG()

	s.chainDb.Close()
	s.dappDb.Close()
	close(s.shutdownChan)
}

// This function will wait for a shutdown and resumes main thread execution
func (s *Krypton) WaitForShutdown() {
	<-s.shutdownChan
}

// StartAutoDAG() spawns a go routine that checks the DAG every autoDAGcheckInterval
// by default that is 10 times per epoch
// in epoch n, if we past autoDAGepochHeight within-epoch blocks,
// it calls krash.MakeDAG  to pregenerate the DAG for the next epoch n+1
// if it does not exist yet as well as remove the DAG for epoch n-1
// the loop quits if autodagquit channel is closed, it can safely restart and
// stop any number of times.
// For any more sophisticated pattern of DAG generation, use CLI subcommand
// makedag
func (self *Krypton) StartAutoDAG() {
	if self.autodagquit != nil {
		return // already started
	}
	go func() {
		glog.V(logger.Info).Infof("Automatic pregeneration of krash DAG ON (krash dir: %s)", krash.DefaultDir)
		var nextEpoch uint64
		timer := time.After(0)
		self.autodagquit = make(chan bool)
		for {
			select {
			case <-timer:
				glog.V(logger.Info).Infof("checking DAG (krash dir: %s)", krash.DefaultDir)
				currentBlock := self.BlockChain().CurrentBlock().NumberU64()
				thisEpoch := currentBlock / epochLength
				if nextEpoch <= thisEpoch {
					if currentBlock%epochLength > autoDAGepochHeight {
						if thisEpoch > 0 {
							previousDag, previousDagFull := dagFiles(thisEpoch - 1)
							os.Remove(filepath.Join(krash.DefaultDir, previousDag))
							os.Remove(filepath.Join(krash.DefaultDir, previousDagFull))
							glog.V(logger.Info).Infof("removed DAG for epoch %d (%s)", thisEpoch-1, previousDag)
						}
						nextEpoch = thisEpoch + 1
						dag, _ := dagFiles(nextEpoch)
						if _, err := os.Stat(dag); os.IsNotExist(err) {
							glog.V(logger.Info).Infof("Pregenerating DAG for epoch %d (%s)", nextEpoch, dag)
							err := krash.MakeDAG(nextEpoch*epochLength, "") // "" -> krash.DefaultDir
							if err != nil {
								glog.V(logger.Error).Infof("Error generating DAG for epoch %d (%s)", nextEpoch, dag)
								return
							}
						} else {
							glog.V(logger.Error).Infof("DAG for epoch %d (%s)", nextEpoch, dag)
						}
					}
				}
				timer = time.After(autoDAGcheckInterval)
			case <-self.autodagquit:
				return
			}
		}
	}()
}

// stopAutoDAG stops automatic DAG pregeneration by quitting the loop
func (self *Krypton) StopAutoDAG() {
	if self.autodagquit != nil {
		close(self.autodagquit)
		self.autodagquit = nil
	}
	glog.V(logger.Info).Infof("Automatic pregeneration of krash DAG OFF (krash dir: %s)", krash.DefaultDir)
}

// HTTPClient returns the light http client used for fetching offchain docs
// (natspec, source for verification)
func (self *Krypton) HTTPClient() *httpclient.HTTPClient {
	return self.httpclient
}

func (self *Krypton) Solc() (*compiler.Solidity, error) {
	var err error
	if self.solc == nil {
		self.solc, err = compiler.New(self.SolcPath)
	}
	return self.solc, err
}

// set in js console via admin interface or wrapper from cli flags
func (self *Krypton) SetSolc(solcPath string) (*compiler.Solidity, error) {
	self.SolcPath = solcPath
	self.solc = nil
	return self.Solc()
}

// dagFiles(epoch) returns the two alternative DAG filenames (not a path)
// 1) <revision>-<hex(seedhash[8])> 2) full-R<revision>-<hex(seedhash[8])>
func dagFiles(epoch uint64) (string, string) {
	seedHash, _ := krash.GetSeedHash(epoch * epochLength)
	dag := fmt.Sprintf("full-R%d-%x", krashRevision, seedHash[:8])
	return dag, "full-R" + dag
}

func saveBlockchainVersion(db krdb.Database, bcVersion int) {
	d, _ := db.Get([]byte("BlockchainVersion"))
	blockchainVersion := common.NewValue(d).Uint()

	if blockchainVersion == 0 {
		db.Put([]byte("BlockchainVersion"), common.NewValue(bcVersion).Bytes())
	}
}

// upgradeChainDatabase ensures that the chain database stores block split into
// separate header and body entries.
func upgradeChainDatabase(db krdb.Database) error {
	// Short circuit if the head block is stored already as separate header and body
	data, err := db.Get([]byte("LastBlock"))
	if err != nil {
		return nil
	}
	head := common.BytesToHash(data)

	if block := core.GetBlockByHashOld(db, head); block == nil {
		return nil
	}
	// At least some of the database is still the old format, upgrade (skip the head block!)
	glog.V(logger.Info).Info("Old database detected, upgrading...")

	if db, ok := db.(*krdb.LDBDatabase); ok {
		blockPrefix := []byte("block-hash-")
		for it := db.NewIterator(); it.Next(); {
			// Skip anything other than a combined block
			if !bytes.HasPrefix(it.Key(), blockPrefix) {
				continue
			}
			// Skip the head block (merge last to signal upgrade completion)
			if bytes.HasSuffix(it.Key(), head.Bytes()) {
				continue
			}
			// Load the block, split and serialize (order!)
			block := core.GetBlockByHashOld(db, common.BytesToHash(bytes.TrimPrefix(it.Key(), blockPrefix)))

			if err := core.WriteTd(db, block.Hash(), block.DeprecatedTd()); err != nil {
				return err
			}
			if err := core.WriteBody(db, block.Hash(), &types.Body{block.Transactions(), block.Uncles()}); err != nil {
				return err
			}
			if err := core.WriteHeader(db, block.Header()); err != nil {
				return err
			}
			if err := db.Delete(it.Key()); err != nil {
				return err
			}
		}
		// Lastly, upgrade the head block, disabling the upgrade mechanism
		current := core.GetBlockByHashOld(db, head)

		if err := core.WriteTd(db, current.Hash(), current.DeprecatedTd()); err != nil {
			return err
		}
		if err := core.WriteBody(db, current.Hash(), &types.Body{current.Transactions(), current.Uncles()}); err != nil {
			return err
		}
		if err := core.WriteHeader(db, current.Header()); err != nil {
			return err
		}
	}
	return nil
}

func addMipmapBloomBins(db krdb.Database) (err error) {
	const mipmapVersion uint = 2

	// check if the version is set. We ignore data for now since there's
	// only one version so we can easily ignore it for now
	var data []byte
	data, _ = db.Get([]byte("setting-mipmap-version"))
	if len(data) > 0 {
		var version uint
		if err := rlp.DecodeBytes(data, &version); err == nil && version == mipmapVersion {
			return nil
		}
	}

	defer func() {
		if err == nil {
			var val []byte
			val, err = rlp.EncodeToBytes(mipmapVersion)
			if err == nil {
				err = db.Put([]byte("setting-mipmap-version"), val)
			}
			return
		}
	}()
	latestBlock := core.GetBlock(db, core.GetHeadBlockHash(db))
	if latestBlock == nil { // clean database
		return
	}

	tstart := time.Now()
	glog.V(logger.Info).Infoln("upgrading db log bloom bins")
	for i := uint64(0); i <= latestBlock.NumberU64(); i++ {
		hash := core.GetCanonicalHash(db, i)
		if (hash == common.Hash{}) {
			return fmt.Errorf("chain db corrupted. Could not find block %d.", i)
		}
		core.WriteMipmapBloom(db, i, core.GetBlockReceipts(db, hash))
	}
	glog.V(logger.Info).Infoln("upgrade completed in", time.Since(tstart))
	return nil
}
