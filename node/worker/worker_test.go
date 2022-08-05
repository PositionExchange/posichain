package worker

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/PositionExchange/posichain/core/state"

	"github.com/ethereum/go-ethereum/core/rawdb"

	blockfactory "github.com/PositionExchange/posichain/block/factory"
	"github.com/PositionExchange/posichain/common/denominations"
	"github.com/PositionExchange/posichain/core"
	"github.com/PositionExchange/posichain/core/types"
	"github.com/PositionExchange/posichain/core/vm"
	chain2 "github.com/PositionExchange/posichain/internal/chain"
	"github.com/PositionExchange/posichain/internal/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// Test accounts
	testBankKey, _  = crypto.GenerateKey()
	testBankAddress = crypto.PubkeyToAddress(testBankKey.PublicKey)
	testBankFunds   = big.NewInt(8000000000000000000)

	chainConfig  = params.TestChainConfig
	blockFactory = blockfactory.ForTest
)

func TestNewWorker(t *testing.T) {
	// Setup a new blockchain with genesis block containing test token on test address
	var (
		database = rawdb.NewMemoryDatabase()
		gspec    = core.Genesis{
			Config:  chainConfig,
			Factory: blockFactory,
			Alloc:   core.GenesisAlloc{testBankAddress: {Balance: testBankFunds}},
			ShardID: 10,
		}
		engine = chain2.NewEngine()
	)

	genesis := gspec.MustCommit(database)
	_ = genesis
	chain, err := core.NewBlockChain(database, state.NewDatabase(database), nil, gspec.Config, engine, vm.Config{}, nil)

	if err != nil {
		t.Error(err)
	}
	// Create a new worker
	worker := New(params.TestChainConfig, chain, engine)

	if worker.GetCurrentState().GetBalance(crypto.PubkeyToAddress(testBankKey.PublicKey)).Cmp(testBankFunds) != 0 {
		t.Error("Worker state is not setup correctly")
	}
}

func TestCommitTransactions(t *testing.T) {
	// Setup a new blockchain with genesis block containing test token on test address
	var (
		database = rawdb.NewMemoryDatabase()
		gspec    = core.Genesis{
			Config:  chainConfig,
			Factory: blockFactory,
			Alloc:   core.GenesisAlloc{testBankAddress: {Balance: testBankFunds}},
			ShardID: 0,
		}
		engine = chain2.NewEngine()
	)

	gspec.MustCommit(database)
	chain, _ := core.NewBlockChain(database, state.NewDatabase(database), nil, gspec.Config, engine, vm.Config{}, nil)

	// Create a new worker
	worker := New(params.TestChainConfig, chain, engine)

	// Generate a test tx
	baseNonce := worker.GetCurrentState().GetNonce(crypto.PubkeyToAddress(testBankKey.PublicKey))
	randAmount := rand.Float32()
	tx, _ := types.SignTx(types.NewTransaction(baseNonce, testBankAddress, uint32(0), big.NewInt(int64(denominations.One*randAmount)), params.TxGas, nil, nil), types.HomesteadSigner{}, testBankKey)

	// Commit the tx to the worker
	txs := make(map[common.Address]types.Transactions)
	txs[testBankAddress] = types.Transactions{tx}
	err := worker.CommitTransactions(
		txs, nil, testBankAddress,
	)
	if err != nil {
		t.Error(err)
	}

	if len(worker.GetCurrentReceipts()) == 0 {
		t.Error("No receipt is created for new transactions")
	}

	if len(worker.current.txs) != 1 {
		t.Error("Transaction is not committed")
	}
}
