package shardingconfig

import (
	"math/big"

	"github.com/PositionExchange/posichain/numeric"

	"github.com/PositionExchange/posichain/internal/genesis"
)

// DockernetSchedule is the local docker testnet sharding
// configuration schedule.
var DockernetSchedule dockernetSchedule

type dockernetSchedule struct{}

const (
	dockernetBlocksPerEpoch = 10

	dockernetVdfDifficulty = 5000 // This takes about 10s to finish the vdf

	// Epoch versions
	dockernetV1Epoch = 110

	// DockerNetHTTPPattern is the http pattern for devnet.
	DockerNetHTTPPattern = "https://api.s%d.k.posichain.org"
	// DockerNetWSPattern is the websocket pattern for devnet.
	DockerNetWSPattern = "wss://ws.s%d.k.posichain.org"
)

func (ds dockernetSchedule) InstanceForEpoch(epoch *big.Int) Instance {
	switch {
	case epoch.Cmp(big.NewInt(dockernetV1Epoch)) >= 0:
		return dockernetV1
	default: // genesis
		return dockernetV0
	}
}

func (ds dockernetSchedule) BlocksPerEpoch() uint64 {
	return dockernetBlocksPerEpoch
}

func (ds dockernetSchedule) CalcEpochNumber(blockNum uint64) *big.Int {
	oldEpoch := blockNum / ds.BlocksPerEpoch()
	return big.NewInt(int64(oldEpoch))
}

func (ds dockernetSchedule) IsLastBlock(blockNum uint64) bool {
	return (blockNum+1)%ds.BlocksPerEpoch() == 0
}

func (ds dockernetSchedule) EpochLastBlock(epochNum uint64) uint64 {
	return ds.BlocksPerEpoch()*(epochNum+1) - 1
}

func (ds dockernetSchedule) VdfDifficulty() int {
	return dockernetVdfDifficulty
}

func (ds dockernetSchedule) GetNetworkID() NetworkID {
	return DockerNet
}

// GetShardingStructure is the sharding structure for dockernet.
func (ds dockernetSchedule) GetShardingStructure(numShard, shardID int) []map[string]interface{} {
	return genShardingStructure(numShard, shardID, DockerNetHTTPPattern, DockerNetWSPattern)
}

// IsSkippedEpoch returns if an epoch was skipped on shard due to staking epoch
func (ds dockernetSchedule) IsSkippedEpoch(shardID uint32, epoch *big.Int) bool {
	return false
}

var (
	dockernetReshardingEpoch = []*big.Int{
		big.NewInt(0),
		big.NewInt(dockernetV1Epoch),
	}
	dockernetV0 = MustNewInstance(2, 4, 2, 0, numeric.OneDec(), genesis.DockernetOperatedAccounts, genesis.DockernetFoundationalAccounts, emptyAllowlist, dockernetReshardingEpoch, DockernetSchedule.BlocksPerEpoch())
	dockernetV1 = MustNewInstance(2, 8, 2, 0, numeric.MustNewDecFromStr("0.7"), genesis.DockernetOperatedAccounts, genesis.DockernetFoundationalAccounts, emptyAllowlist, dockernetReshardingEpoch, DockernetSchedule.BlocksPerEpoch())
)
