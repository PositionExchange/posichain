package node

import (
	"bytes"
	"encoding/gob"
	"harmony-benchmark/blockchain"
	"harmony-benchmark/client"
	"harmony-benchmark/common"
	"harmony-benchmark/consensus"
	"harmony-benchmark/p2p"
	"net"
	"os"
	"time"
)

const (
	// The max number of transaction per a block.
	MaxNumberOfTransactionsPerBlock = 3000
)

// NodeHandler handles a new incoming connection.
func (node *Node) NodeHandler(conn net.Conn) {
	defer conn.Close()

	// Read p2p message payload
	content, err := p2p.ReadMessageContent(conn)

	if err != nil {
		node.log.Error("Read p2p data failed", "err", err, "node", node)
		return
	}
	consensusObj := node.Consensus

	msgCategory, err := common.GetMessageCategory(content)
	if err != nil {
		node.log.Error("Read node type failed", "err", err, "node", node)
		return
	}

	msgType, err := common.GetMessageType(content)
	if err != nil {
		node.log.Error("Read action type failed", "err", err, "node", node)
		return
	}

	msgPayload, err := common.GetMessagePayload(content)
	if err != nil {
		node.log.Error("Read message payload failed", "err", err, "node", node)
		return
	}

	switch msgCategory {
	case common.CONSENSUS:
		actionType := consensus.ConsensusMessageType(msgType)
		switch actionType {
		case consensus.CONSENSUS:
			if consensusObj.IsLeader {
				consensusObj.ProcessMessageLeader(msgPayload)
			} else {
				consensusObj.ProcessMessageValidator(msgPayload)
			}
		}
	case common.NODE:
		actionType := NodeMessageType(msgType)
		switch actionType {
		case TRANSACTION:
			node.transactionMessageHandler(msgPayload)
		case BLOCK:
			if node.Client != nil {
				blockMsgType := BlockMessageType(msgPayload[0])
				switch blockMsgType {
				case SYNC:
					decoder := gob.NewDecoder(bytes.NewReader(msgPayload[1:])) // skip the SYNC messge type
					blocks := new([]*blockchain.Block)
					decoder.Decode(blocks)
					if node.Client != nil && blocks != nil {
						node.Client.UpdateBlocks(*blocks)
					}
				}
			}
		case CONTROL:
			controlType := msgPayload[0]
			if ControlMessageType(controlType) == STOP {
				node.log.Debug("Stopping Node", "node", node, "numBlocks", len(node.blockchain.Blocks), "numTxsProcessed", node.countNumTransactionsInBlockchain())

				sizeInBytes := node.UtxoPool.GetSizeInByteOfUtxoMap()
				node.log.Debug("UtxoPool Report", "numEntries", len(node.UtxoPool.UtxoMap), "sizeInBytes", sizeInBytes)

				avgBlockSizeInBytes := 0
				txCount := 0
				avgTxSize := 0

				for _, block := range node.blockchain.Blocks {
					byteBuffer := bytes.NewBuffer([]byte{})
					encoder := gob.NewEncoder(byteBuffer)
					encoder.Encode(block)
					avgBlockSizeInBytes += len(byteBuffer.Bytes())

					txCount += len(block.Transactions)

					byteBuffer = bytes.NewBuffer([]byte{})
					encoder = gob.NewEncoder(byteBuffer)
					encoder.Encode(block.Transactions)
					avgTxSize += len(byteBuffer.Bytes())
				}
				avgBlockSizeInBytes = avgBlockSizeInBytes / len(node.blockchain.Blocks)
				avgTxSize = avgTxSize / txCount

				node.log.Debug("Blockchain Report", "numBlocks", len(node.blockchain.Blocks), "avgBlockSize", avgBlockSizeInBytes, "numTxs", txCount, "avgTxSzie", avgTxSize)

				os.Exit(0)
			}
		}
	case common.CLIENT:
		actionType := client.ClientMessageType(msgType)
		switch actionType {
		case client.TRANSACTION:
			if node.Client != nil {
				node.Client.TransactionMessageHandler(msgPayload)
			}
		}
	}
}

func (node *Node) transactionMessageHandler(msgPayload []byte) {
	txMessageType := TransactionMessageType(msgPayload[0])

	switch txMessageType {
	case SEND:
		txDecoder := gob.NewDecoder(bytes.NewReader(msgPayload[1:])) // skip the SEND messge type

		txList := new([]*blockchain.Transaction)
		err := txDecoder.Decode(&txList)
		if err != nil {
			node.log.Error("Failed deserializing transaction list", "node", node)
		}
		node.addPendingTransactions(*txList)
	case REQUEST:
		reader := bytes.NewBuffer(msgPayload[1:])
		var txIds map[[32]byte]bool
		buf := make([]byte, 32) // 32 byte hash Id
		for {
			_, err := reader.Read(buf)
			if err != nil {
				break
			}

			var txId [32]byte
			copy(txId[:], buf)
			txIds[txId] = true
		}

		var txToReturn []*blockchain.Transaction
		for _, tx := range node.pendingTransactions {
			if txIds[tx.ID] {
				txToReturn = append(txToReturn, tx)
			}
		}
		// TODO: return the transaction list to requester
	case UNLOCK:
		txAndProofDecoder := gob.NewDecoder(bytes.NewReader(msgPayload[1:])) // skip the UNLOCK messge type

		txAndProofs := new([]*blockchain.Transaction)
		err := txAndProofDecoder.Decode(&txAndProofs)
		if err != nil {
			node.log.Error("Failed deserializing transaction and proofs list", "node", node)
		}
		node.log.Debug("RECEIVED UNLOCK MESSAGE", "num", len(*txAndProofs))

		node.addPendingTransactions(*txAndProofs)
	}
}

// WaitForConsensusReady ...
func (node *Node) WaitForConsensusReady(readySignal chan int) {
	node.log.Debug("Waiting for Consensus ready", "node", node)

	var newBlock *blockchain.Block
	timeoutCount := 0
	for { // keep waiting for Consensus ready
		retry := false
		select {
		case <-readySignal:
			time.Sleep(100 * time.Millisecond) // Delay a bit so validator is catched up.
		case <-time.After(8 * time.Second):
			retry = true
			node.Consensus.ResetState()
			timeoutCount++
			node.log.Debug("Consensus timeout, retry!", "count", timeoutCount, "node", node)
		}

		//node.log.Debug("Adding new block", "currentChainSize", len(node.blockchain.Blocks), "numTxs", len(node.blockchain.GetLatestBlock().Transactions), "PrevHash", node.blockchain.GetLatestBlock().PrevBlockHash, "Hash", node.blockchain.GetLatestBlock().Hash)
		if !retry {
			for {
				// Once we have more than 100 transactions pending we will try creating a new block
				if len(node.pendingTransactions) >= 100 {
					selectedTxs, crossShardTxAndProofs := node.getTransactionsForNewBlock(MaxNumberOfTransactionsPerBlock)

					if len(selectedTxs) == 0 {
						node.log.Debug("No valid transactions exist", "pendingTx", len(node.pendingTransactions))
					} else {
						node.log.Debug("Creating new block", "numTxs", len(selectedTxs), "pendingTxs", len(node.pendingTransactions), "currentChainSize", len(node.blockchain.Blocks))

						node.transactionInConsensus = selectedTxs
						node.log.Debug("CROSS SHARD TX", "num", len(crossShardTxAndProofs))
						node.CrossTxsInConsensus = crossShardTxAndProofs
						newBlock = blockchain.NewBlock(selectedTxs, node.blockchain.GetLatestBlock().Hash, node.Consensus.ShardID)
						break
					}
				}
				// If not enough transactions to run Consensus,
				// periodically check whether we have enough transactions to package into block.
				time.Sleep(1 * time.Second)
			}
		}

		// Send the new block to Consensus so it can be confirmed.
		if newBlock != nil {
			node.BlockChannel <- *newBlock
		}
	}
}

// This is called by consensus participants to verify the block they are running consensus on
func (node *Node) SendBackProofOfAcceptOrReject() {
	if node.ClientPeer != nil && len(node.CrossTxsToReturn) != 0 {
		node.crossTxToReturnMutex.Lock()
		proofs := []blockchain.CrossShardTxProof{}
		for _, txAndProof := range node.CrossTxsToReturn {
			proofs = append(proofs, *txAndProof.Proof)
		}
		node.CrossTxsToReturn = nil
		node.crossTxToReturnMutex.Unlock()

		node.log.Debug("SENDING PROOF TO CLIENT", "proofs", len(proofs))
		p2p.SendMessage(*node.ClientPeer, client.ConstructProofOfAcceptOrRejectMessage(proofs))
	}
}

// This is called by consensus leader to sync new blocks with other clients/nodes.
// NOTE: For now, just send to the client (basically not broadcasting)
func (node *Node) BroadcastNewBlock(newBlock *blockchain.Block) {
	if node.ClientPeer != nil {
		node.log.Debug("SENDING NEW BLOCK TO CLIENT")
		p2p.SendMessage(*node.ClientPeer, ConstructBlocksSyncMessage([]blockchain.Block{*newBlock}))
	}
}

// This is called by consensus participants to verify the block they are running consensus on
func (node *Node) VerifyNewBlock(newBlock *blockchain.Block) bool {
	return node.UtxoPool.VerifyTransactions(newBlock.Transactions)
}

// This is called by consensus participants, after consensus is done, to:
// 1. add the new block to blockchain
// 2. [leader] move cross shard tx and proof to the list where they wait to be sent to the client
func (node *Node) PostConsensusProcessing(newBlock *blockchain.Block) {
	node.AddNewBlock(newBlock)

	if node.Consensus.IsLeader {
		// Move crossTx-in-consensus into the list to be returned to client
		for _, crossTxAndProof := range node.CrossTxsInConsensus {
			crossTxAndProof.Proof.BlockHash = newBlock.Hash
			// TODO: fill in the signature proofs
		}
		if len(node.CrossTxsInConsensus) != 0 {
			node.addCrossTxsToReturn(node.CrossTxsInConsensus)
			node.CrossTxsInConsensus = []*blockchain.CrossShardTxAndProof{}
		}

		node.SendBackProofOfAcceptOrReject()
		node.BroadcastNewBlock(newBlock)
	}
}

func (node *Node) AddNewBlock(newBlock *blockchain.Block) {
	// Add it to blockchain
	node.blockchain.Blocks = append(node.blockchain.Blocks, newBlock)
	// Update UTXO pool
	node.UtxoPool.Update(newBlock.Transactions)
	// Clear transaction-in-Consensus list
	node.transactionInConsensus = []*blockchain.Transaction{}
}
