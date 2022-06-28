package slash

import (
	"math/big"

	"github.com/PositionExchange/posichain/core/types"
	"github.com/PositionExchange/posichain/internal/params"
	"github.com/PositionExchange/posichain/shard"
)

// CommitteeReader ..
type CommitteeReader interface {
	Config() *params.ChainConfig
	ReadShardState(epoch *big.Int) (*shard.State, error)
	CurrentBlock() *types.Block
}
