package shardingconfig

import (
	"fmt"
	"math/big"
	"testing"
)

func TestMainnetInstanceForEpoch(t *testing.T) {
	tests := []struct {
		epoch    *big.Int
		instance Instance
	}{
		{
			big.NewInt(0),
			mainnetV0,
		},
	}

	for _, test := range tests {
		in := MainnetSchedule.InstanceForEpoch(test.epoch)
		if in.NumShards() != test.instance.NumShards() || in.NumNodesPerShard() != test.instance.NumNodesPerShard() {
			t.Errorf("can't get the right instane for epoch: %v\n", test.epoch)
		}
	}
}

func TestCalcEpochNumber(t *testing.T) {
	tests := []struct {
		block uint64
		epoch *big.Int
	}{
		{
			0,
			big.NewInt(0),
		},
		{
			1,
			big.NewInt(0),
		},
		{
			16383,
			big.NewInt(0),
		},
		{
			16384,
			big.NewInt(1),
		},
		{
			16385,
			big.NewInt(1),
		},
		{
			81919,
			big.NewInt(4),
		},
		{
			81700,
			big.NewInt(4),
		},
		{
			5849087,
			big.NewInt(356),
		},
		{
			5849088,
			big.NewInt(357),
		},
	}

	for i, test := range tests {
		ep := MainnetSchedule.CalcEpochNumber(test.block)
		if ep.Cmp(test.epoch) != 0 {
			t.Errorf("CalcEpochNumber error: index %v, got %v, expect %v\n", i, ep, test.epoch)
		}
	}
}

func TestIsLastBlock(t *testing.T) {
	tests := []struct {
		block  uint64
		result bool
	}{
		{
			0,
			false,
		},
		{
			1,
			false,
		},
		{
			16384,
			false,
		},
		{
			16383,
			true,
		},
		{
			32768,
			false,
		},
		{
			32767,
			true,
		},
		{
			49151,
			true,
		},
	}

	for i, test := range tests {
		ep := MainnetSchedule.IsLastBlock(test.block)
		if test.result != ep {
			t.Errorf("IsLastBlock error: index %v, got %v, expect %v\n", i, ep, test.result)
		}
	}
}

func TestEpochLastBlock(t *testing.T) {
	tests := []struct {
		epoch     uint64
		lastBlock uint64
	}{
		{
			0,
			16383,
		},
		{
			1,
			32767,
		},
		{
			2,
			49151,
		},
		{
			3,
			65535,
		},
		{
			358,
			5881855,
		},
	}

	for i, test := range tests {
		ep := MainnetSchedule.EpochLastBlock(test.epoch)
		if test.lastBlock != ep {
			t.Errorf("EpochLastBlock error: index %v, got %v, expect %v\n", i, ep, test.lastBlock)
		}
	}
}

func TestGetShardingStructure(t *testing.T) {
	shardID := 0
	numShard := 4
	res := genShardingStructure(numShard, shardID, "https://api.s%d.posichain.org", "ws://ws.s%d.posichain.org")
	if len(res) != 4 || !res[0]["current"].(bool) || res[1]["current"].(bool) || res[2]["current"].(bool) || res[3]["current"].(bool) {
		t.Error("Error when generating sharding structure")
	}
	for i := 0; i < numShard; i++ {
		if res[i]["current"].(bool) != (i == shardID) {
			t.Error("Error when generating sharding structure")
		}
		if res[i]["shardID"].(int) != i {
			t.Error("Error when generating sharding structure")
		}
		if res[i]["http"].(string) != fmt.Sprintf("https://api.s%d.posichain.org", i) {
			t.Error("Error when generating sharding structure")
		}
		if res[i]["ws"].(string) != fmt.Sprintf("ws://ws.s%d.posichain.org", i) {
			t.Error("Error when generating sharding structure")
		}
	}
}
