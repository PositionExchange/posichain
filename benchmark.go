package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/simple-rules/harmony-benchmark/attack"
	"github.com/simple-rules/harmony-benchmark/consensus"
	"github.com/simple-rules/harmony-benchmark/db"
	"github.com/simple-rules/harmony-benchmark/log"
	"github.com/simple-rules/harmony-benchmark/node"
	"github.com/simple-rules/harmony-benchmark/utils"
)

const (
	AttackProbability = 20
)

func attackDetermination(attackedMode int) bool {
	switch attackedMode {
	case 0:
		return false
	case 1:
		return true
	case 2:
		return rand.Intn(100) < AttackProbability
	}
	return false
}

func startProfiler(shardID string, logFolder string) {
	err := utils.RunCmd("./bin/profiler", "-pid", strconv.Itoa(os.Getpid()), "-shard_id", shardID, "-log_folder", logFolder)
	if err != nil {
		log.Error("Failed to start profiler")
	}
}

func InitLDBDatabase(ip string, port string) (*db.LDBDatabase, error) {
	// TODO(minhdoan): Refactor this.
	dbFileName := "/tmp/harmony_" + ip + port + ".dat"
	var err = os.RemoveAll(dbFileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	return db.NewLDBDatabase(dbFileName, 0, 0)
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "IP of the node")
	port := flag.String("port", "9000", "port of the node.")
	configFile := flag.String("config_file", "config.txt", "file containing all ip addresses")
	logFolder := flag.String("log_folder", "latest", "the folder collecting the logs of this execution")
	attackedMode := flag.Int("attacked_mode", 0, "0 means not attacked, 1 means attacked, 2 means being open to be selected as attacked")
	dbSupported := flag.Int("db_supported", 0, "0 means not db_supported, 1 means db_supported")
	flag.Parse()

	// Set up randomization seed.
	rand.Seed(int64(time.Now().Nanosecond()))

	// Attack determination.
	attack.GetInstance().SetAttackEnabled(attackDetermination(*attackedMode))

	distributionConfig := utils.NewDistributionConfig()
	distributionConfig.ReadConfigFile(*configFile)
	shardID := distributionConfig.GetShardID(*ip, *port)
	peers := distributionConfig.GetPeers(*ip, *port, shardID)
	leader := distributionConfig.GetLeader(shardID)

	var role string
	if leader.Ip == *ip && leader.Port == *port {
		role = "leader"
	} else {
		role = "validator"
	}

	// Setup a logger to stdout and log file.
	logFileName := fmt.Sprintf("./%v/%s-%v-%v.log", *logFolder, role, *ip, *port)
	h := log.MultiHandler(
		log.StdoutHandler,
		log.Must.FileHandler(logFileName, log.JSONFormat()), // Log to file
		// log.Must.NetHandler("tcp", ":3000", log.JSONFormat()) // Log to remote
	)
	log.Root().SetHandler(h)

	// Consensus object.
	consensus := consensus.NewConsensus(*ip, *port, shardID, peers, leader)

	// Start Profiler for leader
	if role == "leader" {
		startProfiler(shardID, *logFolder)
	}

	// Set logger to attack model.
	attack.GetInstance().SetLogger(consensus.Log)
	// Current node.
	currentNode := node.New(consensus, true)
	// Initialize leveldb if dbSupported.
	if *dbSupported == 1 {
		currentNode.DB, _ = InitLDBDatabase(*ip, *port)
	}
	// Create client peer.
	clientPeer := distributionConfig.GetClientPeer()
	// If there is a client configured in the node list.
	if clientPeer != nil {
		currentNode.ClientPeer = clientPeer
	}

	// Assign closure functions to the consensus object
	consensus.BlockVerifier = currentNode.VerifyNewBlock
	consensus.OnConsensusDone = currentNode.PostConsensusProcessing

	// Temporary testing code, to be removed.
	currentNode.AddTestingAddresses(10000)

	if consensus.IsLeader {
		// Let consensus run
		go func() {
			consensus.WaitForNewBlock(currentNode.BlockChannel)
		}()
		// Node waiting for consensus readiness to create new block
		go func() {
			currentNode.WaitForConsensusReady(consensus.ReadySignal)
		}()
	}

	currentNode.StartServer(*port)
}
