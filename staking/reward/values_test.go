package reward

import (
	"math/big"
	"testing"

	shardingconfig "github.com/PositionExchange/posichain/internal/configs/sharding"
)

func TestGetPreStakingRewardsFromBlockNumber(t *testing.T) {
	refMainnetRewards, _ := new(big.Int).SetString("0", 10)
	mainnetRewards := getTotalPreStakingNetworkRewards(shardingconfig.MainNet)
	if refMainnetRewards.Cmp(mainnetRewards) != 0 {
		t.Errorf("Expected mainnet rewards to be %v NOT %v", refMainnetRewards, mainnetRewards)
	}

	refTestnetRewards, _ := new(big.Int).SetString("0", 10)
	testnetRewards := getTotalPreStakingNetworkRewards(shardingconfig.TestNet)
	if refTestnetRewards.Cmp(testnetRewards) != 0 {
		t.Errorf("Expected testnet rewards to be %v NOT %v", refTestnetRewards, testnetRewards)
	}
}
