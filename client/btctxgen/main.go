/*
The btctxgen iterates the btc tx history block by block, transaction by transaction.

The btxtxiter provide a simple api called `NextTx` for us to move thru TXs one by one.

Same as txgen, iterate on each shard to generate simulated TXs (GenerateSimulatedTransactions):

 1. Get a new btc tx
 2. If it's a coinbase tx, create a corresponding coinbase tx in our blockchain
 3. Otherwise, create a normal TX, which might be cross-shard and might not, depending on whether all the TX inputs belong to the current shard.

Same as txgen, send single shard tx shard by shard, then broadcast cross shard tx.

TODO

Some todos for ricl
  * correct the logic to outputing to one of the input shard, rather than the current shard
*/
package main

import (
	"flag"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/simple-rules/harmony-benchmark/blockchain"
	"github.com/simple-rules/harmony-benchmark/client"
	"github.com/simple-rules/harmony-benchmark/client/btctxiter"
	"github.com/simple-rules/harmony-benchmark/client/config"
	"github.com/simple-rules/harmony-benchmark/consensus"
	"github.com/simple-rules/harmony-benchmark/log"
	"github.com/simple-rules/harmony-benchmark/node"
	"github.com/simple-rules/harmony-benchmark/p2p"
	proto_node "github.com/simple-rules/harmony-benchmark/proto/node"

	"github.com/piotrnar/gocoin/lib/btc"
)

type txGenSettings struct {
	crossShard        bool
	maxNumTxsPerBatch int
}

type TXRef struct {
	txID    [32]byte
	shardID uint32
}

var (
	utxoPoolMutex sync.Mutex
	setting       txGenSettings
	btcTXIter     btctxiter.BTCTXIterator
	utxoMapping   map[string]TXRef // btcTXID to { txID, shardID }
)

// Generates at most "maxNumTxs" number of simulated transactions based on the current UtxoPools of all shards.
// The transactions are generated by going through the existing utxos and
// randomly select a subset of them as the input for each new transaction. The output
// address of the new transaction are randomly selected from [0 - N), where N is the total number of fake addresses.
//
// When crossShard=true, besides the selected utxo input, select another valid utxo as input from the same address in a second shard.
// Similarly, generate another utxo output in that second shard.
//
// NOTE: the genesis block should contain N coinbase transactions which add
//       token (1000) to each address in [0 - N). See node.AddTestingAddresses()
//
// Params:
//     shardID                    - the shardID for current shard
//     dataNodes                  - nodes containing utxopools of all shards
// Returns:
//     all single-shard txs
//     all cross-shard txs
func generateSimulatedTransactions(shardID int, dataNodes []*node.Node) ([]*blockchain.Transaction, []*blockchain.Transaction) {
	/*
		  UTXO map structure:
		  {
			  address: {
				  txID: {
					  outputIndex: value
				  }
			  }
		  }
	*/

	utxoPoolMutex.Lock()
	txs := []*blockchain.Transaction{}
	crossTxs := []*blockchain.Transaction{}

	nodeShardID := dataNodes[shardID].Consensus.ShardID
	cnt := 0

LOOP:
	for true {
		btcTx := btcTXIter.NextTx()
		tx := blockchain.Transaction{}
		isCrossShardTx := false
		if btcTx.IsCoinBase() {
			tx.TxInput = []blockchain.TXInput{*blockchain.NewTXInput(blockchain.NewOutPoint(&blockchain.TxID{}, math.MaxUint32), [20]byte{}, nodeShardID)}
		} else {
			for _, btcTXI := range btcTx.TxIn {
				btcTXIDStr := btc.NewUint256(btcTXI.Input.Hash[:]).String()
				txRef := utxoMapping[btcTXIDStr]
				if txRef.shardID != nodeShardID {
					isCrossShardTx = true
				}
				tx.TxInput = append(tx.TxInput, *blockchain.NewTXInput(blockchain.NewOutPoint(&txRef.txID, btcTXI.Input.Vout), [20]byte{}, txRef.shardID))
			}
		}

		for _, btcTXO := range btcTx.TxOut {
			btcTXOAddr := btc.NewAddrFromPkScript(btcTXO.Pk_script, false)
			if btcTXOAddr == nil {
				log.Warn("TxOut: can't decode address")
			}
			txo := blockchain.TXOutput{Amount: int(btcTXO.Value), Address: btcTXOAddr.Hash160, ShardID: nodeShardID}
			tx.TxOutput = append(tx.TxOutput, txo)
		}
		tx.SetID()
		utxoMapping[btcTx.Hash.String()] = TXRef{tx.ID, nodeShardID}
		if isCrossShardTx {
			crossTxs = append(crossTxs, &tx)
		} else {
			txs = append(txs, &tx)
		}
		// log.Debug("[Generator] transformed btc tx", "block height", btcTXIter.GetBlockIndex(), "block tx count", btcTXIter.GetBlock().TxCount, "block tx cnt", len(btcTXIter.GetBlock().Txs), "txi", len(tx.TxInput), "txo", len(tx.TxOutput), "txCount", cnt)
		cnt++
		if cnt >= setting.maxNumTxsPerBatch {
			break LOOP
		}
	}

	utxoPoolMutex.Unlock()

	log.Debug("[Generator] generated transations", "single-shard", len(txs), "cross-shard", len(crossTxs))
	return txs, crossTxs
}

func initClient(clientNode *node.Node, clientPort string, leaders *[]p2p.Peer, nodes *[]*node.Node) {
	if clientPort == "" {
		return
	}

	clientNode.Client = client.NewClient(leaders)

	// This func is used to update the client's utxopool when new blocks are received from the leaders
	updateBlocksFunc := func(blocks []*blockchain.Block) {
		log.Debug("Received new block from leader", "len", len(blocks))
		for _, block := range blocks {
			for _, node := range *nodes {
				if node.Consensus.ShardID == block.ShardId {
					log.Debug("Adding block from leader", "shardId", block.ShardId)
					// Add it to blockchain
					utxoPoolMutex.Lock()
					node.AddNewBlock(block)
					utxoPoolMutex.Unlock()
				} else {
					continue
				}
			}
		}
	}
	clientNode.Client.UpdateBlocks = updateBlocksFunc

	// Start the client server to listen to leader's message
	go func() {
		clientNode.StartServer(clientPort)
	}()
}

func main() {
	configFile := flag.String("config_file", "local_config.txt", "file containing all ip addresses and config")
	maxNumTxsPerBatch := flag.Int("max_num_txs_per_batch", 100, "number of transactions to send per message")
	logFolder := flag.String("log_folder", "latest", "the folder collecting the logs of this execution")
	flag.Parse()

	// Read the configs
	configr := config.NewConfig()
	configr.ReadConfigFile(*configFile)
	leaders, shardIDs := configr.GetLeadersAndShardIds()

	// Do cross shard tx if there are more than one shard
	setting.crossShard = len(shardIDs) > 1
	setting.maxNumTxsPerBatch = *maxNumTxsPerBatch

	// TODO(Richard): refactor this chuck to a single method
	// Setup a logger to stdout and log file.
	logFileName := fmt.Sprintf("./%v/txgen.log", *logFolder)
	h := log.MultiHandler(
		log.StdoutHandler,
		log.Must.FileHandler(logFileName, log.LogfmtFormat()), // Log to file
		// log.Must.NetHandler("tcp", ":3000", log.JSONFormat()) // Log to remote
	)
	log.Root().SetHandler(h)

	btcTXIter.Init()
	utxoMapping = make(map[string]TXRef)

	// Nodes containing utxopools to mirror the shards' data in the network
	nodes := []*node.Node{}
	for _, shardID := range shardIDs {
		nodes = append(nodes, node.New(&consensus.Consensus{ShardID: shardID}, nil))
	}

	// Client/txgenerator server node setup
	clientPort := configr.GetClientPort()
	consensusObj := consensus.NewConsensus("0", clientPort, "0", nil, p2p.Peer{})
	clientNode := node.New(consensusObj, nil)

	initClient(clientNode, clientPort, &leaders, &nodes)

	// Transaction generation process
	time.Sleep(10 * time.Second) // wait for nodes to be ready
	start := time.Now()
	totalTime := 300.0 //run for 5 minutes

	for true {
		t := time.Now()
		if t.Sub(start).Seconds() >= totalTime {
			log.Debug("Generator timer ended.", "duration", (int(t.Sub(start))), "startTime", start, "totalTime", totalTime)
			break
		}

		allCrossTxs := []*blockchain.Transaction{}
		// Generate simulated transactions
		for i, leader := range leaders {
			txs, crossTxs := generateSimulatedTransactions(i, nodes)
			allCrossTxs = append(allCrossTxs, crossTxs...)

			log.Debug("[Generator] Sending single-shard txs ...", "leader", leader, "numTxs", len(txs), "numCrossTxs", len(crossTxs), "block height", btcTXIter.GetBlockIndex())
			msg := proto_node.ConstructTransactionListMessage(txs)
			p2p.SendMessage(leader, msg)
			// Note cross shard txs are later sent in batch
		}

		if len(allCrossTxs) > 0 {
			log.Debug("[Generator] Broadcasting cross-shard txs ...", "allCrossTxs", len(allCrossTxs))
			msg := proto_node.ConstructTransactionListMessage(allCrossTxs)
			p2p.BroadcastMessage(leaders, msg)

			// Put cross shard tx into a pending list waiting for proofs from leaders
			if clientPort != "" {
				clientNode.Client.PendingCrossTxsMutex.Lock()
				for _, tx := range allCrossTxs {
					clientNode.Client.PendingCrossTxs[tx.ID] = tx
				}
				clientNode.Client.PendingCrossTxsMutex.Unlock()
			}
		}

		time.Sleep(500 * time.Millisecond) // Send a batch of transactions periodically
	}

	// Send a stop message to stop the nodes at the end
	msg := proto_node.ConstructStopMessage()
	peers := append(configr.GetValidators(), leaders...)
	p2p.BroadcastMessage(peers, msg)
}
