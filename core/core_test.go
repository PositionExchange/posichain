package core

import (
	"math/big"
	"testing"

	blockfactory "github.com/PositionExchange/posichain/block/factory"
	"github.com/PositionExchange/posichain/core/types"
	shardingconfig "github.com/PositionExchange/posichain/internal/configs/sharding"
	"github.com/PositionExchange/posichain/shard"
)

func TestIsEpochBlock(t *testing.T) {
	blockNumbered := func(n int64) *types.Block {
		return types.NewBlock(
			blockfactory.NewTestHeader().With().Number(big.NewInt(n)).Header(),
			nil, nil, nil, nil, nil,
		)
	}
	tests := []struct {
		schedule shardingconfig.Schedule
		block    *types.Block
		expected bool
	}{
		{
			shardingconfig.MainnetSchedule,
			blockNumbered(1002),
			false,
		},
		{
			shardingconfig.MainnetSchedule,
			blockNumbered(0),
			true,
		},
		{
			shardingconfig.MainnetSchedule,
			blockNumbered(16384),
			true,
		},
		{
			shardingconfig.TestnetSchedule,
			blockNumbered(75),
			false,
		},
		{
			shardingconfig.TestnetSchedule,
			blockNumbered(8192),
			true,
		},
		{
			shardingconfig.TestnetSchedule,
			blockNumbered(16384),
			true,
		},
		{
			shardingconfig.DevnetSchedule,
			blockNumbered(0),
			true,
		},
		{
			shardingconfig.DevnetSchedule,
			blockNumbered(2),
			false,
		},
		{
			shardingconfig.DevnetSchedule,
			blockNumbered(450),
			true,
		},
		{
			shardingconfig.DevnetSchedule,
			blockNumbered(900),
			true,
		},
	}
	for i, test := range tests {
		shard.Schedule = test.schedule
		r := IsEpochBlock(test.block)
		if r != test.expected {
			t.Errorf("index: %v, expected: %v, got: %v\n", i, test.expected, r)
		}
	}
}
