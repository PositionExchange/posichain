package shardingconfig

import (
	"math/big"

	"github.com/PositionExchange/posichain/internal/genesis"
	"github.com/PositionExchange/posichain/numeric"
)

// DevnetSchedule is the long-running public devnet sharding
// configuration schedule.
var DevnetSchedule devnetSchedule

type devnetSchedule struct{}

const (
	devnetBlocksPerEpoch = 5

	// This takes about 20s to finish the vdf
	devnetVdfDifficulty = 10000

	// DevNetHTTPPattern is the http pattern for devnet.
	DevNetHTTPPattern = "http://s%d.d.posichain.org"
	// DevNetWSPattern is the websocket pattern for devnet.
	DevNetWSPattern = "wss://ws.s%d.d.posichain.org"
)

func (ts devnetSchedule) InstanceForEpoch(epoch *big.Int) Instance {
	return devnetV0
}

func (ts devnetSchedule) BlocksPerEpoch() uint64 {
	return devnetBlocksPerEpoch
}

func (ts devnetSchedule) CalcEpochNumber(blockNum uint64) *big.Int {
	oldEpoch := blockNum / ts.BlocksPerEpoch()
	return big.NewInt(int64(oldEpoch))
}

func (ts devnetSchedule) IsLastBlock(blockNum uint64) bool {
	return (blockNum+1)%ts.BlocksPerEpoch() == 0
}

func (ts devnetSchedule) EpochLastBlock(epochNum uint64) uint64 {
	return ts.BlocksPerEpoch()*(epochNum+1) - 1
}

func (ts devnetSchedule) VdfDifficulty() int {
	return devnetVdfDifficulty
}

func (ts devnetSchedule) GetNetworkID() NetworkID {
	return DevNet
}

// GetShardingStructure is the sharding structure for devnet.
func (ts devnetSchedule) GetShardingStructure(numShard, shardID int) []map[string]interface{} {
	return genShardingStructure(numShard, shardID, DevNetHTTPPattern, DevNetWSPattern)
}

// IsSkippedEpoch returns if an epoch was skipped on shard due to staking epoch
func (ts devnetSchedule) IsSkippedEpoch(shardID uint32, epoch *big.Int) bool {
	return false
}

var devnetReshardingEpoch = []*big.Int{big.NewInt(0)}
var devnetV0 = MustNewInstance(2, 4, 2, 0, numeric.OneDec(), genesis.HarmonyAccounts, genesis.FoundationalNodeAccounts, emptyAllowlist, devnetReshardingEpoch, DevnetSchedule.BlocksPerEpoch())
