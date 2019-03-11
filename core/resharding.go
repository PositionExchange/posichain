package core

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"math/rand"
	"sort"
	"strconv"

	"github.com/harmony-one/bls/ffi/go/bls"
	"github.com/harmony-one/harmony/internal/utils/contract"

	"github.com/harmony-one/harmony/core/types"
)

const (
	// InitialSeed is the initial random seed, a magic number to answer everything, remove later
	InitialSeed uint32 = 42
	// FirstEpoch is the number of the first epoch.
	// TODO(minhdoan): we should design the first epoch as 0. Please figure out how to change other logic to make it 0
	FirstEpoch = 1
)

// ShardingState is data structure hold the sharding state
type ShardingState struct {
	epoch      uint64 // current epoch
	rnd        uint64 // random seed for resharding
	numShards  int
	shardState types.ShardState
}

// sortedCommitteeBySize will sort shards by size
// Suppose there are N shards, the first N/2 larger shards are called active committees
// the rest N/2 smaller committees are called inactive committees
// actually they are all just normal shards
// TODO: sort the committee weighted by total staking instead of shard size
func (ss *ShardingState) sortCommitteeBySize() {
	sort.Slice(ss.shardState, func(i, j int) bool {
		return len(ss.shardState[i].NodeList) > len(ss.shardState[j].NodeList)
	})
}

// assignNewNodes add new nodes into the N/2 active committees evenly
func (ss *ShardingState) assignNewNodes(newNodeList []types.NodeID) {
	ss.sortCommitteeBySize()
	numActiveShards := ss.numShards / 2
	Shuffle(newNodeList)
	for i, nid := range newNodeList {
		id := i % numActiveShards
		ss.shardState[id].NodeList = append(ss.shardState[id].NodeList, nid)
	}
}

// cuckooResharding uses cuckoo rule to reshard X% of active committee(shards) into inactive committee(shards)
func (ss *ShardingState) cuckooResharding(percent float64) {
	numActiveShards := ss.numShards / 2
	kickedNodes := []types.NodeID{}
	for i := range ss.shardState {
		if i >= numActiveShards {
			break
		}
		numKicked := int(percent * float64(len(ss.shardState[i].NodeList)))
		if numKicked == 0 {
			numKicked++
		}
		length := len(ss.shardState[i].NodeList)
		tmp := ss.shardState[i].NodeList[length-numKicked:]
		kickedNodes = append(kickedNodes, tmp...)
		ss.shardState[i].NodeList = ss.shardState[i].NodeList[:length-numKicked]
	}

	Shuffle(kickedNodes)
	numInactiveShards := ss.numShards - numActiveShards
	for i, nid := range kickedNodes {
		id := numActiveShards + i%numInactiveShards
		ss.shardState[id].NodeList = append(ss.shardState[id].NodeList, nid)
	}
}

// assignLeaders will first add new nodes into shards, then use cuckoo rule to reshard to get new shard state
func (ss *ShardingState) assignLeaders() {
	for i := 0; i < ss.numShards; i++ {
		Shuffle(ss.shardState[i].NodeList)
		ss.shardState[i].Leader = ss.shardState[i].NodeList[0]
	}
}

// UpdateShardState will first add new nodes into shards, then use cuckoo rule to reshard to get new shard state
func (ss *ShardingState) UpdateShardState(newNodeList []types.NodeID, percent float64) {
	rand.Seed(int64(ss.rnd))
	ss.sortCommitteeBySize()
	// TODO: separate shuffling and leader assignment
	ss.assignLeaders()
	ss.assignNewNodes(newNodeList)
	ss.cuckooResharding(percent)
}

// Shuffle will shuffle the list with result uniquely determined by seed, assuming there is no repeat items in the list
func Shuffle(list []types.NodeID) {
	// Sort to make sure everyone will generate the same with the same rand seed.
	sort.Slice(list, func(i, j int) bool {
		return types.CompareNodeID(list[i], list[j]) == -1
	})
	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
}

// GetBlockNumberFromEpoch calculates the block number where epoch sharding information is stored
func GetBlockNumberFromEpoch(epoch uint64) uint64 {
	number := epoch * uint64(BlocksPerEpoch) // currently we use the first block in each epoch
	return number
}

// GetEpochFromBlockNumber calculates the epoch number the block belongs to
func GetEpochFromBlockNumber(blockNumber uint64) uint64 {
	return blockNumber / uint64(BlocksPerEpoch)
}

// GetShardingStateFromBlockChain will retrieve random seed and shard map from beacon chain for given a epoch
func GetShardingStateFromBlockChain(bc *BlockChain, epoch uint64) *ShardingState {
	number := GetBlockNumberFromEpoch(epoch)
	shardState := bc.GetShardStateByNumber(number)

	rndSeedBytes := bc.GetRandSeedByNumber(number)
	rndSeed := binary.BigEndian.Uint64(rndSeedBytes[:])

	return &ShardingState{epoch: epoch, rnd: rndSeed, shardState: shardState, numShards: len(shardState)}
}

// CalculateNewShardState get sharding state from previous epoch and calculate sharding state for new epoch
// TODO: currently, we just mock everything
func CalculateNewShardState(bc *BlockChain, epoch uint64) types.ShardState {
	if epoch == FirstEpoch {
		return getInitShardState(3, 10)
	}
	ss := GetShardingStateFromBlockChain(bc, epoch-1)
	newNodeList := fakeNewNodeList(int64(ss.rnd))
	percent := ss.calculateKickoutRate(newNodeList)
	ss.UpdateShardState(newNodeList, percent)
	return ss.shardState
}

// calculateKickoutRate calculates the cuckoo rule kick out rate in order to make committee balanced
func (ss *ShardingState) calculateKickoutRate(newNodeList []types.NodeID) float64 {
	numActiveCommittees := ss.numShards / 2
	newNodesPerShard := math.Ceil(float64(len(newNodeList)) / float64(numActiveCommittees))
	ss.sortCommitteeBySize()
	L := len(ss.shardState[0].NodeList)
	if L == 0 {
		return 0.0
	}
	rate := newNodesPerShard / float64(L)
	return math.Max(0.1, math.Min(rate, 1.0))
}

// getInitShardState returns the initial shard state at genesis.
func getInitShardState(numberOfShards, numNodesPerShard int) types.ShardState {
	shardState := types.ShardState{}
	for i := 0; i < numberOfShards; i++ {
		com := types.Committee{ShardID: uint32(i)}
		if i == 0 {
			for j := 0; j < numNodesPerShard; j++ {
				priKey := bls.SecretKey{}
				priKey.SetHexString(contract.InitialBeaconChainAccounts[i].Private)
				addrBytes := priKey.GetPublicKey().GetAddress()
				address := hex.EncodeToString(addrBytes[:])
				com.NodeList = append(com.NodeList, types.NodeID(address))
			}
		}
		shardState = append(shardState, com)
	}
	return shardState
}

// remove later after new nodes list generation ready
func fakeNewNodeList(seed int64) []types.NodeID {
	rand.Seed(seed)
	numNewNodes := rand.Intn(10)
	nodeList := []types.NodeID{}
	for i := 0; i < numNewNodes; i++ {
		nid := strconv.Itoa(int(rand.Int63()))
		nodeList = append(nodeList, types.NodeID(nid))
	}
	return nodeList
}
