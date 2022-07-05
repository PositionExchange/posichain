package shardingconfig

import (
	"math/big"

	"github.com/PositionExchange/posichain/internal/genesis"
	"github.com/PositionExchange/posichain/numeric"
)

// TestnetSchedule is the long-running public testnet sharding
// configuration schedule.
var TestnetSchedule testnetSchedule

type testnetSchedule struct{}

const (
	// 4.5 hours per epoch (given 2s block time)
	testnetBlocksPerEpoch = 8192

	testnetVdfDifficulty = 10000 // This takes about 20s to finish the vdf

	// TestNetHTTPPattern is the http pattern for testnet.
	TestNetHTTPPattern = "https://api.s%d.t.posichain.org"
	// TestNetWSPattern is the websocket pattern for testnet.
	TestNetWSPattern = "wss://ws.s%d.t.posichain.org"
)

func (ts testnetSchedule) InstanceForEpoch(epoch *big.Int) Instance {
	switch {
	// Add more version here
	default: // genesis
		return testnetV0
	}
}

func (ts testnetSchedule) BlocksPerEpoch() uint64 {
	return testnetBlocksPerEpoch
}

func (ts testnetSchedule) CalcEpochNumber(blockNum uint64) *big.Int {
	oldEpoch := blockNum / ts.BlocksPerEpoch()
	return big.NewInt(int64(oldEpoch))
}

func (ts testnetSchedule) IsLastBlock(blockNum uint64) bool {
	return (blockNum+1)%ts.BlocksPerEpoch() == 0
}

func (ts testnetSchedule) EpochLastBlock(epochNum uint64) uint64 {
	return ts.BlocksPerEpoch()*(epochNum+1) - 1
}

func (ts testnetSchedule) VdfDifficulty() int {
	return testnetVdfDifficulty
}

func (ts testnetSchedule) GetNetworkID() NetworkID {
	return TestNet
}

// GetShardingStructure is the sharding structure for testnet.
func (ts testnetSchedule) GetShardingStructure(numShard, shardID int) []map[string]interface{} {
	return genShardingStructure(numShard, shardID, TestNetHTTPPattern, TestNetWSPattern)
}

// IsSkippedEpoch returns if an epoch was skipped on shard due to staking epoch
func (ts testnetSchedule) IsSkippedEpoch(shardID uint32, epoch *big.Int) bool {
	return false
}

var testnetReshardingEpoch = []*big.Int{big.NewInt(0)}

var testnetV0 = MustNewInstance(2, 5, 4, 0, numeric.OneDec(), genesis.TNHarmonyAccounts, genesis.TNFoundationalAccounts, emptyAllowlist, testnetReshardingEpoch, TestnetSchedule.BlocksPerEpoch())
