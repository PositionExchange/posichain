package shardingconfig

import (
	"math/big"

	"github.com/PositionExchange/posichain/numeric"

	"github.com/PositionExchange/posichain/internal/genesis"
)

const (
	blocksPerEpoch = 16384 // 2^14

	// This takes about 100s to finish the vdf
	mainnetVdfDifficulty = 50000

	// MainNetHTTPPattern is the http pattern for mainnet.
	MainNetHTTPPattern = "https://api.s%d.posichain.org"
	// MainNetWSPattern is the websocket pattern for mainnet.
	MainNetWSPattern = "wss://ws.s%d.posichain.org"
)

var (
	// map of epochs skipped due to staking launch on mainnet
	skippedEpochs = map[uint32][]*big.Int{}
)

// MainnetSchedule is the mainnet sharding configuration schedule.
var MainnetSchedule mainnetSchedule

type mainnetSchedule struct{}

func (ms mainnetSchedule) InstanceForEpoch(epoch *big.Int) Instance {
	switch {
	default: // genesis
		return mainnetV0
	}
}

func (ms mainnetSchedule) BlocksPerEpoch() uint64 {
	return blocksPerEpoch
}

func (ms mainnetSchedule) CalcEpochNumber(blockNum uint64) *big.Int {
	epoch := blockNum / ms.BlocksPerEpoch()
	return big.NewInt(int64(epoch))
}

func (ms mainnetSchedule) IsLastBlock(blockNum uint64) bool {
	return (blockNum+1)%ms.BlocksPerEpoch() == 0
}

func (ms mainnetSchedule) EpochLastBlock(epochNum uint64) uint64 {
	return ms.BlocksPerEpoch()*(epochNum+1) - 1
}

func (ms mainnetSchedule) VdfDifficulty() int {
	return mainnetVdfDifficulty
}

func (ms mainnetSchedule) GetNetworkID() NetworkID {
	return MainNet
}

// GetShardingStructure is the sharding structure for mainnet.
func (ms mainnetSchedule) GetShardingStructure(numShard, shardID int) []map[string]interface{} {
	return genShardingStructure(numShard, shardID, MainNetHTTPPattern, MainNetWSPattern)
}

// IsSkippedEpoch returns if an epoch was skipped on shard due to staking epoch
func (ms mainnetSchedule) IsSkippedEpoch(shardID uint32, epoch *big.Int) bool {
	if skipped, exists := skippedEpochs[shardID]; exists {
		for _, e := range skipped {
			if epoch.Cmp(e) == 0 {
				return true
			}
		}
	}
	return false
}

var mainnetReshardingEpoch = []*big.Int{big.NewInt(0)}

var (
	mainnetV0 = MustNewInstance(1, 30, 25, 0, numeric.OneDec(), genesis.MainnetOperatedAccounts, genesis.FoundationalNodeAccounts, emptyAllowlist, mainnetReshardingEpoch, MainnetSchedule.BlocksPerEpoch())
)
