package shardingconfig

import (
	"fmt"
	"math/big"

	"github.com/PositionExchange/posichain/numeric"

	"github.com/PositionExchange/posichain/internal/genesis"
)

// LocalnetSchedule is the local testnet sharding
// configuration schedule.
var LocalnetSchedule localnetSchedule

type localnetSchedule struct{}

const (
	localnetBlocksPerEpoch = 10

	localnetVdfDifficulty = 5000 // This takes about 10s to finish the vdf
)

func (ls localnetSchedule) InstanceForEpoch(epoch *big.Int) Instance {
	switch {
	default: // genesis
		return localnetV0
	}
}

func (ls localnetSchedule) BlocksPerEpoch() uint64 {
	return localnetBlocksPerEpoch
}

func (ls localnetSchedule) CalcEpochNumber(blockNum uint64) *big.Int {
	oldEpoch := blockNum / ls.BlocksPerEpoch()
	return big.NewInt(int64(oldEpoch))
}

func (ls localnetSchedule) IsLastBlock(blockNum uint64) bool {
	return (blockNum+1)%ls.BlocksPerEpoch() == 0
}

func (ls localnetSchedule) EpochLastBlock(epochNum uint64) uint64 {
	return ls.BlocksPerEpoch()*(epochNum+1) - 1
}

func (ls localnetSchedule) VdfDifficulty() int {
	return localnetVdfDifficulty
}

func (ls localnetSchedule) GetNetworkID() NetworkID {
	return LocalNet
}

// GetShardingStructure is the sharding structure for localnet.
func (ls localnetSchedule) GetShardingStructure(numShard, shardID int) []map[string]interface{} {
	var res []map[string]interface{}
	for i := 0; i < numShard; i++ {
		res = append(res, map[string]interface{}{
			"current": shardID == i,
			"shardID": i,
			"http":    fmt.Sprintf("http://127.0.0.1:%d", 9500+i),
			"ws":      fmt.Sprintf("ws://127.0.0.1:%d", 9800+i),
		})
	}
	return res
}

// IsSkippedEpoch returns if an epoch was skipped on shard due to staking epoch
func (ls localnetSchedule) IsSkippedEpoch(shardID uint32, epoch *big.Int) bool {
	return false
}

var (
	localnetReshardingEpoch = []*big.Int{big.NewInt(0)}
	localnetV0              = MustNewInstance(2, 9, 6, 0, numeric.MustNewDecFromStr("0.68"), genesis.LocalHarmonyAccountsV2, genesis.LocalFnAccountsV2, emptyAllowlist, localnetReshardingEpoch, LocalnetSchedule.BlocksPerEpoch())
)
