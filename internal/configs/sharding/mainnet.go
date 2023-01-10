package shardingconfig

import (
	"math/big"

	"github.com/PositionExchange/posichain/internal/common"
	"github.com/PositionExchange/posichain/internal/genesis"
	"github.com/PositionExchange/posichain/numeric"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

const (
	blocksPerEpoch = 16384 // 2^14

	// This takes about 100s to finish the vdf
	mainnetVdfDifficulty = 50000

	// Epoch versions
	mainnetV1Epoch = 2
	mainnetV2Epoch = 105
	mainnetV3Epoch = 308 // Around November 26, 2022, 08:31 UTC+0 (Saturday)
	mainnetV4Epoch = 442 // Around January 15, 2023, 23:58:50 UTC+0 (Sunday)

	// MainNetHTTPPattern is the http pattern for mainnet.
	MainNetHTTPPattern = "https://api.s%d.posichain.org"
	// MainNetWSPattern is the websocket pattern for mainnet.
	MainNetWSPattern = "wss://ws.s%d.posichain.org"
)

var (
	// map of epochs skipped due to staking launch on mainnet
	skippedEpochs = map[uint32][]*big.Int{}

	emptyAddress = ethCommon.Address{}

	mainnetFeeCollector = mustAddress("0x0000000000000000000000000000000000000000")
)

func mustAddress(addrStr string) ethCommon.Address {
	addr, err := common.ParseAddr(addrStr)
	if err != nil {
		panic("invalid address")
	}
	return addr
}

// MainnetSchedule is the mainnet sharding configuration schedule.
var MainnetSchedule mainnetSchedule

type mainnetSchedule struct{}

func (ms mainnetSchedule) InstanceForEpoch(epoch *big.Int) Instance {
	switch {
	case epoch.Cmp(big.NewInt(mainnetV4Epoch)) >= 0:
		return mainnetV4
	case epoch.Cmp(big.NewInt(mainnetV3Epoch)) >= 0:
		return mainnetV3
	case epoch.Cmp(big.NewInt(mainnetV2Epoch)) >= 0:
		return mainnetV2
	case epoch.Cmp(big.NewInt(mainnetV1Epoch)) >= 0:
		return mainnetV1
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

var mainnetReshardingEpoch = []*big.Int{
	big.NewInt(0),
	big.NewInt(mainnetV1Epoch),
	big.NewInt(mainnetV2Epoch),
	big.NewInt(mainnetV3Epoch),
	big.NewInt(mainnetV4Epoch),
}

var (
	mainnetV0 = MustNewInstance(1, 21, 16, 0, numeric.OneDec(), genesis.MainnetOperatedAccounts, genesis.FoundationalNodeAccounts, emptyAllowlist, emptyAddress, mainnetReshardingEpoch, MainnetSchedule.BlocksPerEpoch())
	mainnetV1 = MustNewInstance(1, 21, 16, 0, numeric.MustNewDecFromStr("0.7"), genesis.MainnetOperatedAccounts, genesis.FoundationalNodeAccounts, emptyAllowlist, emptyAddress, mainnetReshardingEpoch, MainnetSchedule.BlocksPerEpoch())
	mainnetV2 = MustNewInstance(1, 41, 16, 0, numeric.MustNewDecFromStr("0.7"), genesis.MainnetOperatedAccounts, genesis.FoundationalNodeAccounts, emptyAllowlist, emptyAddress, mainnetReshardingEpoch, MainnetSchedule.BlocksPerEpoch())
	mainnetV3 = MustNewInstance(1, 61, 16, 0, numeric.MustNewDecFromStr("0.7"), genesis.MainnetOperatedAccounts, genesis.FoundationalNodeAccounts, emptyAllowlist, emptyAddress, mainnetReshardingEpoch, MainnetSchedule.BlocksPerEpoch())
	mainnetV4 = MustNewInstance(1, 61, 16, 0, numeric.MustNewDecFromStr("0.7"), genesis.MainnetOperatedAccounts, genesis.FoundationalNodeAccounts, emptyAllowlist, mainnetFeeCollector, mainnetReshardingEpoch, MainnetSchedule.BlocksPerEpoch())
)
