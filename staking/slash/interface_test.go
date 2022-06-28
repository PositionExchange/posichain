package slash

import (
	"math/big"

	"github.com/PositionExchange/posichain/core/types"
	"github.com/PositionExchange/posichain/internal/params"
	"github.com/PositionExchange/posichain/shard"
	staking "github.com/PositionExchange/posichain/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	errFakeChainUnexpectEpoch = errors.New("epoch not expected")
)

type fakeBlockChain struct {
	config         params.ChainConfig
	currentBlock   types.Block
	superCommittee shard.State
	snapshots      map[common.Address]staking.ValidatorWrapper
}

func (bc *fakeBlockChain) Config() *params.ChainConfig {
	return &bc.config
}

func (bc *fakeBlockChain) CurrentBlock() *types.Block {
	return &bc.currentBlock
}

func (bc *fakeBlockChain) ReadShardState(epoch *big.Int) (*shard.State, error) {
	if epoch.Cmp(big.NewInt(currentEpoch)) == 0 {
		return nil, errFakeChainUnexpectEpoch
	}
	return &bc.superCommittee, nil
}

func (bc *fakeBlockChain) ReadValidatorSnapshotAtEpoch(epoch *big.Int, addr common.Address) (*staking.ValidatorSnapshot, error) {
	vw, ok := bc.snapshots[addr]
	if !ok {
		return nil, errors.New("missing snapshot")
	}
	return &staking.ValidatorSnapshot{
		Validator: &vw,
		Epoch:     new(big.Int).Set(epoch),
	}, nil
}
