package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/simple-rules/harmony-benchmark/attack"
	"github.com/simple-rules/harmony-benchmark/configr"
	"github.com/simple-rules/harmony-benchmark/consensus"
	"github.com/simple-rules/harmony-benchmark/log"
	"github.com/simple-rules/harmony-benchmark/node"
	"github.com/simple-rules/harmony-benchmark/p2p"

	"github.com/shirou/gopsutil/process"
	"github.com/simple-rules/harmony-benchmark/crypto"
	"github.com/simple-rules/harmony-benchmark/utils"
)

const (
	AttackProbability = 20
)

func getShardId(myIp, myPort string, config *[][]string) string {
	for _, node := range *config {
		ip, port, shardId := node[0], node[1], node[3]
		if ip == myIp && port == myPort {
			return shardId
		}
	}
	return "N/A"
}

func getLeader(myShardId string, config *[][]string) p2p.Peer {
	var leaderPeer p2p.Peer
	for _, node := range *config {
		ip, port, status, shardId := node[0], node[1], node[2], node[3]
		if status == "leader" && myShardId == shardId {
			leaderPeer.Ip = ip
			leaderPeer.Port = port

			// Get public key deterministically based on ip and port
			priKey := crypto.Ed25519Curve.Scalar().SetInt64(int64(utils.GetUniqueIdFromPeer(leaderPeer))) // TODO: figure out why using a random hash value doesn't work for private key (schnorr)
			leaderPeer.PubKey = crypto.GetPublicKeyFromScalar(crypto.Ed25519Curve, priKey)
		}
	}
	return leaderPeer
}

func getPeers(myIp, myPort, myShardId string, config *[][]string) []p2p.Peer {
	var peerList []p2p.Peer
	for _, node := range *config {
		ip, port, status, shardId := node[0], node[1], node[2], node[3]
		if status != "validator" || ip == myIp && port == myPort || myShardId != shardId {
			continue
		}
		// Get public key deterministically based on ip and port
		peer := p2p.Peer{Port: port, Ip: ip}
		priKey := crypto.Ed25519Curve.Scalar().SetInt64(int64(utils.GetUniqueIdFromPeer(peer)))
		peer.PubKey = crypto.GetPublicKeyFromScalar(crypto.Ed25519Curve, priKey)
		peerList = append(peerList, peer)
	}
	return peerList
}

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

func logMemUsage(consensus *consensus.Consensus) {
	p, _ := process.NewProcess(int32(os.Getpid()))
	for {
		info, _ := p.MemoryInfo()
		memMap, _ := p.MemoryMaps(false)
		log.Info("Mem Report", "info", info, "map", memMap)
		time.Sleep(10 * time.Second)
	}
}

// TODO: @ricl, start another process for reporting.
func logCPUUsage(consensus *consensus.Consensus) {
	p, _ := process.NewProcess(int32(os.Getpid()))
	for {
		percent, _ := p.CPUPercent()
		times, _ := p.Times()
		log.Info("CPU Report", "percent", percent, "times", times, "consensus", consensus)
		time.Sleep(10 * time.Second)
	}
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "IP of the node")
	port := flag.String("port", "9000", "port of the node.")
	configFile := flag.String("config_file", "config.txt", "file containing all ip addresses")
	logFolder := flag.String("log_folder", "latest", "the folder collecting the logs of this execution")
	attackedMode := flag.Int("attacked_mode", 0, "0 means not attacked, 1 means attacked, 2 means being open to be selected as attacked")
	flag.Parse()

	// Set up randomization seed.
	rand.Seed(int64(time.Now().Nanosecond()))

	// Attack determination.
	attack.GetInstance().SetAttackEnabled(attackDetermination(*attackedMode))

	config, _ := configr.ReadConfigFile(*configFile)
	shardID := getShardId(*ip, *port, &config)
	peers := getPeers(*ip, *port, shardID, &config)
	leader := getLeader(shardID, &config)

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
	// Logging for consensus.
	go logMemUsage(consensus)
	go logCPUUsage(consensus)

	// Set logger to attack model.
	attack.GetInstance().SetLogger(consensus.Log)
	// Current node.
	currentNode := node.New(consensus)
	// Create client peer.
	clientPeer := configr.GetClientPeer(&config)
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
