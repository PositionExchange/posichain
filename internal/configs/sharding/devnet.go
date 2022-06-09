package shardingconfig

import (
	"math/big"

	"github.com/harmony-one/harmony/internal/genesis"
	"github.com/harmony-one/harmony/internal/params"
	"github.com/harmony-one/harmony/numeric"
)

// DevnetSchedule is the long-running public devnet sharding
// configuration schedule.
var DevnetSchedule devnetSchedule

type devnetSchedule struct{}

const (
	devnetBlocksPerEpoch   = 5
	devnetBlocksPerEpochV2 = 10

	// This takes about 20s to finish the vdf
	devnetVdfDifficulty = 10000

	// DevNetHTTPPattern is the http pattern for devnet.
	DevNetHTTPPattern = "http://s%d.z.d.posichain.com:9500"
	// DevNetWSPattern is the websocket pattern for devnet.
	DevNetWSPattern = "wss://ws.s%d.z.d.posichain.com:9800"
)

func (ts devnetSchedule) InstanceForEpoch(epoch *big.Int) Instance {
	return devnetV0
}

func (ts devnetSchedule) BlocksPerEpochOld() uint64 {
	return devnetBlocksPerEpoch
}

func (ts devnetSchedule) BlocksPerEpoch() uint64 {
	return devnetBlocksPerEpochV2
}

func (ts devnetSchedule) CalcEpochNumber(blockNum uint64) *big.Int {

	firstBlock2s := params.DevnetChainConfig.TwoSecondsEpoch.Uint64() * ts.BlocksPerEpochOld()
	switch {
	case blockNum >= firstBlock2s:
		return big.NewInt(int64((blockNum-firstBlock2s)/ts.BlocksPerEpoch() + params.DevnetChainConfig.TwoSecondsEpoch.Uint64()))
	default: // genesis
		oldEpoch := blockNum / ts.BlocksPerEpochOld()
		return big.NewInt(int64(oldEpoch))
	}

}

func (ts devnetSchedule) IsLastBlock(blockNum uint64) bool {
	firstBlock2s := params.DevnetChainConfig.TwoSecondsEpoch.Uint64() * ts.BlocksPerEpochOld()

	switch {
	case blockNum >= firstBlock2s:
		return (blockNum-firstBlock2s)%ts.BlocksPerEpoch() == ts.BlocksPerEpoch()-1
	default: // genesis
		return (blockNum+1)%ts.BlocksPerEpochOld() == 0
	}
}

func (ts devnetSchedule) EpochLastBlock(epochNum uint64) uint64 {
	firstBlock2s := params.DevnetChainConfig.TwoSecondsEpoch.Uint64() * ts.BlocksPerEpochOld()

	switch {
	case params.DevnetChainConfig.IsTwoSeconds(big.NewInt(int64(epochNum))):
		return firstBlock2s - 1 + ts.BlocksPerEpoch()*(epochNum-params.DevnetChainConfig.TwoSecondsEpoch.Uint64()+1)
	default: // genesis
		return ts.BlocksPerEpochOld()*(epochNum+1) - 1
	}

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

var devnetReshardingEpoch = []*big.Int{
	big.NewInt(0),
	params.DevnetChainConfig.StakingEpoch,
	params.DevnetChainConfig.TwoSecondsEpoch,
}

var devnetV0 = MustNewInstance(2, 4, 2, 0, numeric.OneDec(), genesis.TNHarmonyAccounts, genesis.TNFoundationalAccounts, emptyAllowlist, devnetReshardingEpoch, DevnetSchedule.BlocksPerEpochOld())
