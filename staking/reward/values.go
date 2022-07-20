package reward

import (
	"fmt"
	"math/big"

	"github.com/PositionExchange/posichain/common/denominations"
	"github.com/PositionExchange/posichain/consensus/engine"
	shardingconfig "github.com/PositionExchange/posichain/internal/configs/sharding"
	"github.com/PositionExchange/posichain/internal/params"
	"github.com/PositionExchange/posichain/numeric"
	"github.com/PositionExchange/posichain/shard"
)

var (
	// PreStakedBlocks is the block reward, to be split evenly among block signers in pre-staking era.
	// 0.001 POSI per block.
	PreStakedBlocks = new(big.Int).Mul(big.NewInt(0.001*denominations.Nano), big.NewInt(denominations.Nano))

	// StakedBlocks is the DEFAULT flat-rate block reward for epos staking launch.
	// 0.001 POSI per block.
	StakedBlocks = numeric.NewDecFromBigInt(new(big.Int).Mul(
		big.NewInt(0.001*denominations.Nano), big.NewInt(denominations.Nano),
	))

	// Posichain tokenomics model
	// https://docs.google.com/spreadsheets/d/1aWqPpw2XO2guSJeXVC5n3p71Ug1yKMmpnwFOWYwq5e4/edit#gid=1008442802

	TwoSecStakedBlocks20222023 = numeric.NewDecFromBigInt(big.NewInt(0.1167975157 * denominations.One))
	TwoSecStakedBlocks20242025 = numeric.NewDecFromBigInt(big.NewInt(0.1016296444 * denominations.One))
	TwoSecStakedBlocks20262027 = numeric.NewDecFromBigInt(big.NewInt(0.0508148222 * denominations.One))
	TwoSecStakedBlocks20282029 = numeric.NewDecFromBigInt(big.NewInt(0.0254074111 * denominations.One))
	TwoSecStakedBlocks20302031 = numeric.NewDecFromBigInt(big.NewInt(0.0127037055 * denominations.One))
	TwoSecStakedBlocks20322033 = numeric.NewDecFromBigInt(big.NewInt(0.0063518528 * denominations.One))
	TwoSecStakedBlocks20342035 = numeric.NewDecFromBigInt(big.NewInt(0.0031759264 * denominations.One))

	// TotalInitialTokens is the total amount of tokens (in POSI) at block 0 of the network.
	// This should be set/change on the node's init according to the core.GenesisSpec.
	TotalInitialTokens = numeric.Dec{Int: big.NewInt(0)}

	// None ..
	None = big.NewInt(0)

	// ErrInvalidBeaconChain if given chain is not beacon chain
	ErrInvalidBeaconChain = fmt.Errorf("given chain is not beaconchain")
)

// getPreStakingRewardsFromBlockNumber returns the number of tokens injected into the network
// in the pre-staking era (epoch < staking epoch) in ATTO.
//
// If the block number is > than the last block of an epoch, the last block of the epoch is
// used for the calculation by default.
//
// WARNING: This assumes beacon chain is at most the same block height as another shard in the
// transition from pre-staking to staking era/epoch.
func getPreStakingRewardsFromBlockNumber(id shardingconfig.NetworkID, blockNum *big.Int) *big.Int {
	if blockNum.Cmp(big.NewInt(2)) == -1 {
		// block 0 & 1 does not contain block rewards
		return big.NewInt(0)
	}

	lastBlockInEpoch := blockNum

	switch id {
	case shardingconfig.MainNet:
		lastBlockInEpoch = new(big.Int).SetUint64(shardingconfig.MainnetSchedule.EpochLastBlock(
			params.MainnetChainConfig.StakingEpoch.Uint64() - 1,
		))
	case shardingconfig.TestNet:
		lastBlockInEpoch = new(big.Int).SetUint64(shardingconfig.TestnetSchedule.EpochLastBlock(
			params.TestnetChainConfig.StakingEpoch.Uint64() - 1,
		))
	case shardingconfig.DevNet:
		lastBlockInEpoch = new(big.Int).SetUint64(shardingconfig.DevnetSchedule.EpochLastBlock(
			params.DevnetChainConfig.StakingEpoch.Uint64() - 1,
		))
	case shardingconfig.LocalNet:
		lastBlockInEpoch = new(big.Int).SetUint64(shardingconfig.LocalnetSchedule.EpochLastBlock(
			params.LocalnetChainConfig.StakingEpoch.Uint64() - 1,
		))
	}

	if blockNum.Cmp(lastBlockInEpoch) == 1 {
		blockNum = lastBlockInEpoch
	}

	return new(big.Int).Mul(PreStakedBlocks, new(big.Int).Sub(blockNum, big.NewInt(1)))
}

// WARNING: the data collected here are calculated from a consumer of the Rosetta API.
// If data becomes mission critical, implement a cross-link based approach.
//
// Data Source: https://github.com/harmony-one/jupyter
//
// TODO (dm): use first crosslink of all shards to compute rewards on network instead of relying on constants.
var (
	totalPreStakingNetworkRewardsInAtto = map[shardingconfig.NetworkID][]*big.Int{
		shardingconfig.MainNet: {
			// Below are all of the last blocks of pre-staking era for mainnet.
			//getPreStakingRewardsFromBlockNumber(shardingconfig.MainNet, big.NewInt(999999)),
			//getPreStakingRewardsFromBlockNumber(shardingconfig.MainNet, big.NewInt(999999)),
		},
		shardingconfig.TestNet: {
			// Below are all of the placeholders 'last blocks' of pre-staking era for testnet.
			//getPreStakingRewardsFromBlockNumber(shardingconfig.TestNet, big.NewInt(999999)),
			//getPreStakingRewardsFromBlockNumber(shardingconfig.TestNet, big.NewInt(999999)),
		},
		shardingconfig.DevNet: {
			// Below are all of the placeholders 'last blocks' of pre-staking era for testnet.
			//getPreStakingRewardsFromBlockNumber(shardingconfig.DevNet, big.NewInt(999999)),
			//getPreStakingRewardsFromBlockNumber(shardingconfig.DevNet, big.NewInt(999999)),
		},
		shardingconfig.LocalNet: {
			// Below are all of the placeholders 'last blocks' of pre-staking era for localnet.
			//getPreStakingRewardsFromBlockNumber(shardingconfig.LocalNet, big.NewInt(999999)),
			//getPreStakingRewardsFromBlockNumber(shardingconfig.LocalNet, big.NewInt(999999)),
		},
	}
)

// getTotalPreStakingNetworkRewards in ATTO for given NetworkID
func getTotalPreStakingNetworkRewards(id shardingconfig.NetworkID) *big.Int {
	totalRewards := big.NewInt(0)
	if allRewards, ok := totalPreStakingNetworkRewardsInAtto[id]; ok {
		for _, reward := range allRewards {
			totalRewards = new(big.Int).Add(reward, totalRewards)
		}
	}
	return totalRewards
}

// GetTotalTokens in the network for all shards in POSI.
// This can only be computed with beaconchain if in staking era.
// If not in staking era, returns the rewards given out by the start of staking era.
func GetTotalTokens(chain engine.ChainReader) (numeric.Dec, error) {
	currHeader := chain.CurrentHeader()
	if !chain.Config().IsStaking(currHeader.Epoch()) {
		return GetTotalPreStakingTokens(), nil
	}
	if chain.ShardID() != shard.BeaconChainShardID {
		return numeric.Dec{}, ErrInvalidBeaconChain
	}

	stakingRewards, err := chain.ReadBlockRewardAccumulator(currHeader.Number().Uint64())
	if err != nil {
		return numeric.Dec{}, err
	}
	return GetTotalPreStakingTokens().Add(numeric.NewDecFromBigIntWithPrec(stakingRewards, 18)), nil
}

// GetTotalPreStakingTokens returns the total amount of tokens (in POSI) in the
// network at the the last block of the pre-staking era (epoch < staking epoch).
func GetTotalPreStakingTokens() numeric.Dec {
	preStakingRewards := numeric.NewDecFromBigIntWithPrec(
		getTotalPreStakingNetworkRewards(shard.Schedule.GetNetworkID()), 18,
	)
	return TotalInitialTokens.Add(preStakingRewards)
}

// SetTotalInitialTokens with the given initial tokens (from genesis in ATTO).
func SetTotalInitialTokens(initTokensAsAtto *big.Int) {
	TotalInitialTokens = numeric.NewDecFromBigIntWithPrec(initTokensAsAtto, 18)
}
