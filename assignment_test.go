package harmony

import (
	"fmt"
	"testing"
)

func TestPreStakingEnabledCommittee(T *testing.T) {
	shardNum := 1
	shardSize := 5
	shardHarmonyNodes := 4
	for i := 0; i < shardNum; i++ {
		fmt.Printf("Shard %d\n", i)
		for j := 0; j < shardHarmonyNodes; j++ {
			index := i + j*shardNum // The initial account to use for genesis nodes
			fmt.Printf("hmy account index %d\n", index)
		}
		// add FN runner's key
		for j := shardHarmonyNodes; j < shardSize; j++ {
			index := i + (j-shardHarmonyNodes)*shardNum
			fmt.Printf("fn account index %d\n", index)
		}
	}
}
