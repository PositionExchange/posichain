package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"harmony-benchmark/blockchain"
	"harmony-benchmark/consensus"
	"harmony-benchmark/log"
	"harmony-benchmark/node"
	"harmony-benchmark/p2p"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Get numTxs number of Fake transactions based on the existing UtxoPool.
// The transactions are generated by going through the existing utxos and
// randomly select a subset of them as input to new transactions. The output
// address of the new transaction are randomly selected from 1 - 1000.
// NOTE: the genesis block should contain 1000 coinbase transactions adding
//       value to each address in [1 - 1000]. See node.AddMoreFakeTransactions()
func getNewFakeTransactions(dataNode *node.Node, numTxs int) []*blockchain.Transaction {
	/*
	  UTXO map structure:
	     address - [
	                txId1 - [
	                        outputIndex1 - value1
	                        outputIndex2 - value2
	                       ]
	                txId2 - [
	                        outputIndex1 - value1
	                        outputIndex2 - value2
	                       ]
	               ]
	*/
	var outputs []*blockchain.Transaction
	count := 0
	countAll := 0

	for address, txMap := range dataNode.UtxoPool.UtxoMap {
		for txIdStr, utxoMap := range txMap {
			txId, err := hex.DecodeString(txIdStr)
			if err != nil {
				continue
			}
			for index, value := range utxoMap {
				countAll++
				if rand.Intn(100) < 30 { // 30% sample rate to select UTXO to use for new transactions
					// Spend the money of current UTXO to a random address in [1 - 1000]
					txin := blockchain.TXInput{txId, index, address, dataNode.Consensus.ShardID}
					txout := blockchain.TXOutput{value, strconv.Itoa(rand.Intn(10000))}
					tx := blockchain.Transaction{[32]byte{}, []blockchain.TXInput{txin}, []blockchain.TXOutput{txout}}
					tx.SetID()

					if count >= numTxs {
						continue
					}
					outputs = append(outputs, &tx)
					count++
				}
			}
		}
	}

	log.Debug("UTXO", "poolSize", countAll, "numTxsToSend", numTxs)
	return outputs
}

func getValidators(config string) []p2p.Peer {
	file, _ := os.Open(config)
	fscanner := bufio.NewScanner(file)
	var peerList []p2p.Peer
	for fscanner.Scan() {
		p := strings.Split(fscanner.Text(), " ")
		ip, port, status := p[0], p[1], p[2]
		if status == "leader" {
			continue
		}
		peer := p2p.Peer{Port: port, Ip: ip}
		peerList = append(peerList, peer)
	}
	return peerList
}

func getLeadersAndShardIds(config *[][]string) ([]p2p.Peer, []uint32) {
	var peerList []p2p.Peer
	var shardIds []uint32
	for _, node := range *config {
		ip, port, status, shardId := node[0], node[1], node[2], node[3]
		if status == "leader" {
			peerList = append(peerList, p2p.Peer{Ip: ip, Port: port})
			val, err := strconv.Atoi(shardId)
			if err == nil {
				shardIds = append(shardIds, uint32(val))
			} else {
				log.Error("[Generator] Error parsing the shard Id ", shardId)
			}
		}
	}
	return peerList, shardIds
}

func readConfigFile(configFile string) [][]string {
	file, _ := os.Open(configFile)
	fscanner := bufio.NewScanner(file)

	result := [][]string{}
	for fscanner.Scan() {
		p := strings.Split(fscanner.Text(), " ")
		result = append(result, p)
	}
	return result
}

func main() {
	configFile := flag.String("config_file", "local_config.txt", "file containing all ip addresses and config")
	numTxsPerBatch := flag.Int("num_txs_per_batch", 10000, "number of transactions to send per message")
	logFolder := flag.String("log_folder", "latest", "the folder collecting the logs of this execution")
	flag.Parse()
	config := readConfigFile(*configFile)
	leaders, shardIds := getLeadersAndShardIds(&config)

	// Setup a logger to stdout and log file.
	logFileName := fmt.Sprintf("./%v/tx-generator.log", *logFolder)
	h := log.MultiHandler(
		log.Must.FileHandler(logFileName, log.LogfmtFormat()),
		log.StdoutHandler)
	// In cases where you just want a stdout logger, use the following one instead.
	// h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)

	// Testing node to mirror the node data in consensus
	nodes := []node.Node{}
	for _, shardId := range shardIds {
		node := node.NewNode(&consensus.Consensus{ShardID: shardId})
		node.AddMoreFakeTransactions(10000)
		nodes = append(nodes, node)
	}

	time.Sleep(10 * time.Second) // wait for nodes to be ready

	start := time.Now()
	totalTime := 60.0
	for true {
		t := time.Now()
		if t.Sub(start).Seconds() >= totalTime {
			fmt.Println(int(t.Sub(start)), start, totalTime)
			break
		}

		t = time.Now()
		for i, leader := range leaders {
			txsToSend := getNewFakeTransactions(&nodes[i], *numTxsPerBatch)
			msg := node.ConstructTransactionListMessage(txsToSend)
			fmt.Printf("[Generator] Creating fake txs for leader took %s", time.Since(t))

			log.Debug("[Generator] Sending txs ...", "leader", leader, "numTxs", len(txsToSend))
			p2p.SendMessage(leader, msg)

			// Update local utxo pool to mirror the utxo pool of a real node
			nodes[i].UtxoPool.Update(txsToSend)
		}

		time.Sleep(500 * time.Millisecond) // Send a batch of transactions periodically
	}

	// Send a stop message to stop the nodes at the end
	msg := node.ConstructStopMessage()
	peers := append(getValidators(*configFile), leaders...)
	p2p.BroadcastMessage(peers, msg)
}
