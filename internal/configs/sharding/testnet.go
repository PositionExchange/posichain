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
	// ~4.5 hours per epoch (given 2s block time)
	testnetBlocksPerEpoch = 8192 // 2^13

	// This takes about 20s to finish the vdf
	testnetVdfDifficulty = 10000

	// Epoch versions
	testnetV1Epoch = 68
	testnetV2Epoch = 860 // Around January 10, 2023, 04:42:03 AM (Tuesday)

	// TestNetHTTPPattern is the http pattern for testnet.
	TestNetHTTPPattern = "https://api.s%d.t.posichain.org"
	// TestNetWSPattern is the websocket pattern for testnet.
	TestNetWSPattern = "wss://ws.s%d.t.posichain.org"
)

var testnetFeeCollector = mustAddress("0x0000000000000000000000000000000000000000")

func (ts testnetSchedule) InstanceForEpoch(epoch *big.Int) Instance {
	switch {
	case epoch.Cmp(big.NewInt(testnetV2Epoch)) >= 0:
		return testnetV2
	case epoch.Cmp(big.NewInt(testnetV1Epoch)) >= 0:
		return testnetV1
	default: // genesis
		return testnetV0
	}
}

func (ts testnetSchedule) BlocksPerEpoch() uint64 {
	return testnetBlocksPerEpoch
}

func (ts testnetSchedule) CalcEpochNumber(blockNum uint64) *big.Int {
	epoch := blockNum / ts.BlocksPerEpoch()
	return big.NewInt(int64(epoch))
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

var testnetReshardingEpoch = []*big.Int{
	big.NewInt(0),
	big.NewInt(testnetV1Epoch),
	big.NewInt(testnetV2Epoch),
}

var (
	testnetV0 = MustNewInstance(1, 5, 4, 0, numeric.OneDec(), genesis.TestnetOperatedAccounts, genesis.TestnetFoundationalAccounts, emptyAllowlist, emptyAddress, testnetReshardingEpoch, TestnetSchedule.BlocksPerEpoch())
	testnetV1 = MustNewInstance(1, 15, 4, 0, numeric.MustNewDecFromStr("0.7"), genesis.TestnetOperatedAccounts, genesis.TestnetFoundationalAccounts, emptyAllowlist, emptyAddress, testnetReshardingEpoch, TestnetSchedule.BlocksPerEpoch())
	testnetV2 = MustNewInstance(1, 15, 4, 0, numeric.MustNewDecFromStr("0.7"), genesis.TestnetOperatedAccounts, genesis.TestnetFoundationalAccounts, emptyAllowlist, testnetFeeCollector, testnetReshardingEpoch, TestnetSchedule.BlocksPerEpoch())
)
