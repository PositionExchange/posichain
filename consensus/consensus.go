// Consensus package implements the Cosi PBFT consensus
package consensus // consensus

import (
	"fmt"
	"harmony-benchmark/blockchain"
	"harmony-benchmark/log"
	"harmony-benchmark/p2p"
	"regexp"
	"strconv"
	"sync"
)

// Consensus data containing all info related to one round of consensus process
type Consensus struct {
	state ConsensusState
	// Signatures collected from validators
	commits map[string]string
	// Signatures collected from validators
	responses map[string]string
	// List of validators
	validators []p2p.Peer
	// Leader
	leader p2p.Peer
	// private key of current node
	priKey string
	// Whether I am leader. False means I am validator
	IsLeader bool
	// Leader or validator Id - 2 byte
	nodeId uint16
	// Consensus Id (View Id) - 4 byte
	consensusId uint32
	// Blockhash - 32 byte
	blockHash [32]byte
	// BlockHeader to run consensus on
	blockHeader []byte
	// Shard Id which this node belongs to
	ShardID uint32

	// global consensus mutex
	mutex sync.Mutex

	// Validator specific fields
	// Blocks received but not done with consensus yet
	blocksReceived map[uint32]*BlockConsensusStatus

	// Signal channel for starting a new consensus process
	ReadySignal chan int
	// The verifier func passed from Node object
	BlockVerifier func(*blockchain.Block) bool
	// The post-consensus processing func passed from Node object
	// Called when consensus on a new block is done
	OnConsensusDone func(*blockchain.Block)

	Log log.Logger
}

// This used to keep track of the consensus status of multiple blocks received so far
// This is mainly used in the case that this node is lagging behind and needs to catch up.
// For example, the consensus moved to round N and this node received message(N).
// However, this node may still not finished with round N-1, so the newly received message(N)
// should be stored in this temporary structure. In case the round N-1 finishes, it can catch
// up to the latest state of round N by using this structure.
type BlockConsensusStatus struct {
	blockHeader []byte         // the block header of the block which the consensus is running on
	state       ConsensusState // the latest state of the consensus
}

// NewConsensus creates a new Consensus object
// TODO(minhdoan): Maybe convert it into just New
// FYI, see https://golang.org/doc/effective_go.html?#package-names
func NewConsensus(ip, port, ShardID string, peers []p2p.Peer, leader p2p.Peer) *Consensus {
	consensus := Consensus{}
	Peers := peers
	leaderPeer := leader
	selfPeer := p2p.Peer{Port: port, Ip: ip}

	if leaderPeer == selfPeer {
		consensus.IsLeader = true
	} else {
		consensus.IsLeader = false
	}

	consensus.commits = make(map[string]string)
	consensus.responses = make(map[string]string)

	consensus.leader = leaderPeer
	consensus.validators = Peers

	consensus.priKey = ip + ":" + port // use ip:port as unique key for now

	consensus.consensusId = 0 // or view Id in the original pbft paper

	myShardID, err := strconv.Atoi(ShardID)
	if err != nil {
		panic("Unparseable shard Id" + ShardID)
	}
	consensus.ShardID = uint32(myShardID)

	// For validators to keep track of all blocks received but not yet committed, so as to catch up to latest consensus if lagged behind.
	consensus.blocksReceived = make(map[uint32]*BlockConsensusStatus)

	// For now use socket address as 16 byte Id
	// TODO: populate with correct Id
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		consensus.Log.Crit("Regex Compilation Failed", "err", err, "consensus", consensus)
	}
	socketId := reg.ReplaceAllString(consensus.priKey, "")
	value, err := strconv.Atoi(socketId)
	consensus.nodeId = uint16(value)

	if consensus.IsLeader {
		consensus.ReadySignal = make(chan int)
		// send a signal to indicate it's ready to run consensus
		// this signal is consumed by node object to create a new block and in turn trigger a new consensus on it
		// this is a goroutine because go channel without buffer will block
		go func() {
			consensus.ReadySignal <- 1
		}()
	}

	consensus.Log = log.New()
	return &consensus
}

// Reset the state of the consensus
func (consensus *Consensus) ResetState() {
	consensus.state = FINISHED
	consensus.commits = make(map[string]string)
	consensus.responses = make(map[string]string)
}

// Returns a string representation of this consensus
func (consensus *Consensus) String() string {
	var duty string
	if consensus.IsLeader {
		duty = "LDR" // leader
	} else {
		duty = "VLD" // validator
	}
	return fmt.Sprintf("[duty:%s, priKey:%s, ShardID:%v, nodeId:%v, state:%s]",
		duty, consensus.priKey, consensus.ShardID, consensus.nodeId, consensus.state)
}
