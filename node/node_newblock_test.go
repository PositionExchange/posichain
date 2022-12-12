package node

import (
	"github.com/PositionExchange/posichain/core"
	"github.com/PositionExchange/posichain/internal/chain"
	"github.com/stretchr/testify/require"
	"math/big"
	"strings"
	"testing"

	"github.com/PositionExchange/posichain/consensus"
	"github.com/PositionExchange/posichain/consensus/quorum"
	"github.com/PositionExchange/posichain/core/types"
	"github.com/PositionExchange/posichain/crypto/bls"
	nodeconfig "github.com/PositionExchange/posichain/internal/configs/node"
	shardingconfig "github.com/PositionExchange/posichain/internal/configs/sharding"
	"github.com/PositionExchange/posichain/internal/shardchain"
	"github.com/PositionExchange/posichain/internal/utils"
	"github.com/PositionExchange/posichain/multibls"
	"github.com/PositionExchange/posichain/p2p"
	"github.com/PositionExchange/posichain/shard"
	staking "github.com/PositionExchange/posichain/staking/types"
	"github.com/ethereum/go-ethereum/common"
)

func TestFinalizeNewBlockAsync(t *testing.T) {
	blsKey := bls.RandPrivateKey()
	pubKey := blsKey.GetPublicKey()
	leader := p2p.Peer{IP: "127.0.0.1", Port: "8882", ConsensusPubKey: pubKey}
	priKey, _, _ := utils.GenKeyP2P("127.0.0.1", "9902")
	host, err := p2p.NewHost(p2p.HostConfig{
		Self:   &leader,
		BLSKey: priKey,
	})
	if err != nil {
		t.Fatalf("newhost failure: %v", err)
	}
	var testDBFactory = &shardchain.MemDBFactory{}
	engine := chain.NewEngine()
	chainconfig := nodeconfig.GetShardConfig(shard.BeaconChainShardID).GetNetworkType().ChainConfig()
	collection := shardchain.NewCollection(
		nil, testDBFactory, &core.GenesisInitializer{NetworkType: nodeconfig.GetShardConfig(shard.BeaconChainShardID).GetNetworkType()}, engine, &chainconfig,
	)
	blockchain, err := collection.ShardChain(shard.BeaconChainShardID)
	require.NoError(t, err)

	decider := quorum.NewDecider(
		quorum.SuperMajorityVote, shard.BeaconChainShardID,
	)
	consensus, err := consensus.New(
		host, shard.BeaconChainShardID, multibls.GetPrivateKeys(blsKey), nil, decider, 3, false,
	)
	if err != nil {
		t.Fatalf("Cannot craeate consensus: %v", err)
	}

	shard.Schedule = shardingconfig.DevnetSchedule
	nodeconfig.SetNetworkType(nodeconfig.Devnet)
	node := New(host, consensus, engine, collection, nil, nil, nil, nil, nil)

	node.Worker.UpdateCurrent()

	txs := make(map[common.Address]types.Transactions)
	stks := staking.StakingTransactions{}
	node.Worker.CommitTransactions(
		txs, stks, common.Address{},
	)
	commitSigs := make(chan []byte)
	go func() {
		commitSigs <- []byte{}
	}()

	block, _ := node.Worker.FinalizeNewBlock(
		commitSigs, func() uint64 { return 0 }, common.Address{}, nil, nil,
	)

	// work around vrf verification as it's tested in another test.
	node.Blockchain().Config().VRFEpoch = big.NewInt(2)
	if err := VerifyNewBlock(nil, blockchain, nil)(block); err != nil {
		t.Error("New block is not verified successfully:", err)
	}

	node.Blockchain().InsertChain(types.Blocks{block}, false)

	node.Worker.UpdateCurrent()

	_, err = node.Worker.FinalizeNewBlock(
		commitSigs, func() uint64 { return 0 }, common.Address{}, nil, nil,
	)

	if !strings.Contains(err.Error(), "cannot finalize block") {
		t.Error("expect timeout on FinalizeNewBlock")
	}
}
