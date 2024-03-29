// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	blockfactory "github.com/PositionExchange/posichain/block/factory"
	"github.com/PositionExchange/posichain/internal/genesis"
	"github.com/PositionExchange/posichain/internal/params"
	"github.com/PositionExchange/posichain/staking/slash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/PositionExchange/posichain/core/rawdb"
	"github.com/PositionExchange/posichain/core/state"
	"github.com/PositionExchange/posichain/core/types"
	nodeconfig "github.com/PositionExchange/posichain/internal/configs/node"
	shardingconfig "github.com/PositionExchange/posichain/internal/configs/sharding"
	"github.com/PositionExchange/posichain/internal/utils"
	"github.com/PositionExchange/posichain/shard"
)

// no go:generate gencodec -type Genesis -field-override genesisSpecMarshaling -out gen_genesis.go
// no go:generate gencodec -type GenesisAccount -field-override genesisAccountMarshaling -out gen_genesis_account.go

var errGenesisNoConfig = errors.New("genesis has no chain configuration")

const (
	// GenesisEpoch is the number of the genesis epoch.
	GenesisEpoch = 0
)

// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
	Config         *params.ChainConfig  `json:"config"`
	Factory        blockfactory.Factory `json:"-"`
	Nonce          uint64               `json:"nonce"`
	ShardID        uint32               `json:"shardID"`
	Timestamp      uint64               `json:"timestamp"`
	ExtraData      []byte               `json:"extraData"`
	GasLimit       uint64               `json:"gasLimit"       gencodec:"required"`
	Mixhash        common.Hash          `json:"mixHash"`
	Coinbase       common.Address       `json:"coinbase"`
	Alloc          GenesisAlloc         `json:"alloc"          gencodec:"required"`
	ShardStateHash common.Hash          `json:"shardStateHash" gencodec:"required"`
	ShardState     shard.State          `json:"shardState"     gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`
}

// NewGenesisSpec creates a new genesis spec for the given network type and shard ID.
// Note that the shard state is NOT initialized.
func NewGenesisSpec(netType nodeconfig.NetworkType, shardID uint32) *Genesis {
	genesisAlloc := make(GenesisAlloc)
	chainConfig := netType.ChainConfig()
	gasLimit := params.GenesisGasLimit

	if shardID == 0 {
		if genesisAllocList := ResolveGenesisAlloc(netType); genesisAllocList != nil {
			for addr, account := range genesisAllocList {
				genesisAlloc[common.HexToAddress(addr)] = account
			}
		} else if netType != nodeconfig.Mainnet { // All non-mainnet chains get test accounts
			utils.Logger().Info().Msgf("No foundation accounts are provided, initialing test accounts for non-mainnet network")
			gasLimit = params.TestGenesisGasLimit
			foundationTestKey, _ := ecdsa.GenerateKey(
				crypto.S256(),
				strings.NewReader("Test contract key string stream that is fixed so that generated test key are deterministic every time"),
			)
			foundationTestPrivateKey := hex.EncodeToString(crypto.FromECDSA(foundationTestKey))
			foundationTestAddress := crypto.PubkeyToAddress(foundationTestKey.PublicKey)
			utils.Logger().Info().Msgf("Foundation account initiated, balances: %d, address: %s, private key: %s",
				GenesisTestFoundationBalance,
				foundationTestAddress.Hex(),
				foundationTestPrivateKey)
			genesisAlloc[foundationTestAddress] = GenesisAccount{
				Balance: toStorableBalance(GenesisTestFoundationBalance),
			}
		}
	}

	return &Genesis{
		Config:    &chainConfig,
		Factory:   blockfactory.NewFactory(&chainConfig),
		Alloc:     genesisAlloc,
		ShardID:   shardID,
		GasLimit:  gasLimit,
		Timestamp: 1654480800, // GMT: Monday, June 6, 2022 2:00:00 AM
		ExtraData: []byte(genesis.Signature),
	}
}

// GenesisAlloc specifies the initial state that is part of the genesis block.
type GenesisAlloc map[common.Address]GenesisAccount

// UnmarshalJSON is to deserialize the data into GenesisAlloc.
func (ga *GenesisAlloc) UnmarshalJSON(data []byte) error {
	m := make(map[common.UnprefixedAddress]GenesisAccount)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*ga = make(GenesisAlloc)
	for addr, a := range m {
		(*ga)[common.Address(addr)] = a
	}
	return nil
}

// GenesisAccount is an account in the state of the genesis block.
type GenesisAccount struct {
	Code       []byte                      `json:"code,omitempty"`
	Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance    *big.Int                    `json:"balance" gencodec:"required"`
	Nonce      uint64                      `json:"nonce,omitempty"`
	PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
}

// field type overrides for gencodec
type genesisSpecMarshaling struct {
	Nonce      math.HexOrDecimal64
	Timestamp  math.HexOrDecimal64
	ExtraData  hexutil.Bytes
	GasLimit   math.HexOrDecimal64
	GasUsed    math.HexOrDecimal64
	Number     math.HexOrDecimal64
	Difficulty *math.HexOrDecimal256
	Alloc      map[common.UnprefixedAddress]GenesisAccount
}

type genesisAccountMarshaling struct {
	Code       hexutil.Bytes
	Balance    *math.HexOrDecimal256
	Nonce      math.HexOrDecimal64
	Storage    map[storageJSON]storageJSON
	PrivateKey hexutil.Bytes
}

// storageJSON represents a 256 bit byte array, but allows less than 256 bits when
// unmarshaling from hex.
type storageJSON common.Hash

func (h *storageJSON) UnmarshalText(text []byte) error {
	text = bytes.TrimPrefix(text, []byte("0x"))
	if len(text) > 64 {
		return fmt.Errorf("too many hex characters in storage key/value %q", text)
	}
	offset := len(h) - len(text)/2 // pad on the left
	if _, err := hex.Decode(h[offset:], text); err != nil {
		return fmt.Errorf("invalid hex storage key/value %q", text)
	}
	return nil
}

func (h storageJSON) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New common.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database already contains an incompatible genesis block (have %x, new %x)", e.Stored[:8], e.New[:8])
}

func (g *Genesis) configOrDefault(ghash common.Hash) *params.ChainConfig {
	switch {
	case g != nil:
		return g.Config
	default:
		return params.AllProtocolChanges
	}
}

// ToBlock creates the genesis block and writes state of a genesis specification
// to the given database (or discards it if nil).
func (g *Genesis) ToBlock(db ethdb.Database) *types.Block {
	if db == nil {
		utils.Logger().Error().Msg("db should be initialized")
		os.Exit(1)
	}
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))
	for addr, account := range g.Alloc {
		statedb.AddBalance(addr, account.Balance)
		statedb.SetCode(addr, account.Code)
		statedb.SetNonce(addr, account.Nonce)
		for key, value := range account.Storage {
			statedb.SetState(addr, key, value)
		}
	}
	root := statedb.IntermediateRoot(false)
	shardStateBytes, err := shard.EncodeWrapper(g.ShardState, false)
	if err != nil {
		utils.Logger().Error().Msg("failed to rlp-serialize genesis shard state")
		os.Exit(1)
	}
	head := g.Factory.NewHeader(common.Big0).With().
		Number(new(big.Int).SetUint64(g.Number)).
		ShardID(g.ShardID).
		Time(new(big.Int).SetUint64(g.Timestamp)).
		ParentHash(g.ParentHash).
		Extra(g.ExtraData).
		GasLimit(g.GasLimit).
		GasUsed(g.GasUsed).
		MixDigest(g.Mixhash).
		Coinbase(g.Coinbase).
		Root(root).
		ShardStateHash(g.ShardStateHash).
		ShardState(shardStateBytes).
		Header()
	statedb.Commit(false)
	statedb.Database().TrieDB().Commit(root, true)

	return types.NewBlock(head, nil, nil, nil, nil, nil)
}

// Commit writes the block and state of a genesis specification to the database.
// The block is committed as the canonical head block.
func (g *Genesis) Commit(db ethdb.Database) (*types.Block, error) {
	block := g.ToBlock(db)
	if block.Number().Sign() != 0 {
		return nil, fmt.Errorf("can't commit genesis block with number > 0")
	}

	if err := rawdb.WriteBlock(db, block); err != nil {
		return nil, err
	}
	if err := rawdb.WriteReceipts(db, block.Hash(), block.NumberU64(), nil); err != nil {
		return nil, err
	}
	if err := rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64()); err != nil {
		return nil, err
	}
	if err := rawdb.WriteHeadBlockHash(db, block.Hash()); err != nil {
		return nil, err
	}
	if err := rawdb.WriteHeadHeaderHash(db, block.Hash()); err != nil {
		return nil, err
	}

	err := rawdb.WriteShardStateBytes(db, block.Header().Epoch(), block.Header().ShardState())

	if err != nil {
		utils.Logger().Error().Err(err).Msg("Failed to store genesis shard state")
	}

	config := g.Config
	if config == nil {
		config = params.AllProtocolChanges
	}
	rawdb.WriteChainConfig(db, block.Hash(), config)
	return block, nil
}

// MustCommit writes the genesis block and state to db, panicking on error.
// The block is committed as the canonical head block.
func (g *Genesis) MustCommit(db ethdb.Database) *types.Block {
	block, err := g.Commit(db)
	if err != nil {
		panic(err)
	}
	rawdb.WriteBlockRewardAccumulator(db, big.NewInt(0), 0)
	data, err := rlp.EncodeToBytes(slash.Records{})
	if err != nil {
		panic(err)
	}
	if err := rawdb.WritePendingSlashingCandidates(db, data); err != nil {
		panic(err)
	}
	return block
}

// GetGenesisSpec for a given shard
func GetGenesisSpec(shardID uint32) *Genesis {
	if shard.Schedule.GetNetworkID() == shardingconfig.MainNet {
		return NewGenesisSpec(nodeconfig.Mainnet, shardID)
	}
	if shard.Schedule.GetNetworkID() == shardingconfig.DevNet {
		return NewGenesisSpec(nodeconfig.Devnet, shardID)
	}
	if shard.Schedule.GetNetworkID() == shardingconfig.DockerNet {
		return NewGenesisSpec(nodeconfig.Dockernet, shardID)
	}
	if shard.Schedule.GetNetworkID() == shardingconfig.LocalNet {
		return NewGenesisSpec(nodeconfig.Localnet, shardID)
	}
	return NewGenesisSpec(nodeconfig.Testnet, shardID)
}

// GetInitialFunds for a given shard
func GetInitialFunds(shardID uint32) *big.Int {
	spec, total := GetGenesisSpec(shardID), big.NewInt(0)
	for _, account := range spec.Alloc {
		total = new(big.Int).Add(account.Balance, total)
	}
	return total
}
