package shardingconfig

import (
	"fmt"
	"testing"
)

func TestGenerateAccountForNodeInstanceConfig(T *testing.T) {
	shardNum := 1
	shardSize := 30
	shardOperatedNodes := 25
	for i := 0; i < shardNum; i++ {
		fmt.Printf("Shard %d\n", i)
		for j := 0; j < shardOperatedNodes; j++ {
			index := i + j*shardNum
			fmt.Printf("operated account index %d\n", index)
		}
		// add FN runner's key
		for j := shardOperatedNodes; j < shardSize; j++ {
			index := i + (j-shardOperatedNodes)*shardNum
			fmt.Printf("fn account index %d\n", index)
		}
	}
}
