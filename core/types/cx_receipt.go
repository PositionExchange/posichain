package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/harmony-one/harmony/block"
	"github.com/harmony-one/harmony/internal/ctxerror"
)

// CXReceipt represents a receipt for cross-shard transaction
type CXReceipt struct {
	TxHash    common.Hash // hash of the cross shard transaction in source shard
	From      common.Address
	To        *common.Address
	ShardID   uint32
	ToShardID uint32
	Amount    *big.Int
}

// CXReceipts is a list of CXReceipt
type CXReceipts []*CXReceipt

// Len returns the length of s.
func (cs CXReceipts) Len() int { return len(cs) }

// Swap swaps the i'th and the j'th element in s.
func (cs CXReceipts) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (cs CXReceipts) GetRlp(i int) []byte {
	if len(cs) == 0 {
		return []byte{}
	}
	enc, _ := rlp.EncodeToBytes(cs[i])
	return enc
}

// ToShardID returns the destination shardID of the cxReceipt
func (cs CXReceipts) ToShardID(i int) uint32 {
	if len(cs) == 0 {
		return 0
	}
	return cs[i].ToShardID
}

// MaxToShardID returns the maximum destination shardID of cxReceipts
func (cs CXReceipts) MaxToShardID() uint32 {
	maxShardID := uint32(0)
	if len(cs) == 0 {
		return maxShardID
	}
	for i := 0; i < len(cs); i++ {
		if maxShardID < cs[i].ToShardID {
			maxShardID = cs[i].ToShardID
		}
	}
	return maxShardID
}

// CXMerkleProof represents the merkle proof of a collection of ordered cross shard transactions
type CXMerkleProof struct {
	BlockNum      *big.Int      // blockNumber of source shard
	BlockHash     common.Hash   // blockHash of source shard
	ShardID       uint32        // shardID of source shard
	CXReceiptHash common.Hash   // root hash of the cross shard receipts in a given block
	ShardIDs      []uint32      // order list, records destination shardID
	CXShardHashes []common.Hash // ordered hash list, each hash corresponds to one destination shard's receipts root hash
}

// CXReceiptsProof carrys the cross shard receipts and merkle proof
type CXReceiptsProof struct {
	Receipts     CXReceipts
	MerkleProof  *CXMerkleProof
	Header       *block.Header
	CommitSig    []byte
	CommitBitmap []byte
}

// CXReceiptsProofs is a list of CXReceiptsProof
type CXReceiptsProofs []*CXReceiptsProof

// Len returns the length of s.
func (cs CXReceiptsProofs) Len() int { return len(cs) }

// Swap swaps the i'th and the j'th element in s.
func (cs CXReceiptsProofs) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (cs CXReceiptsProofs) GetRlp(i int) []byte {
	if len(cs) == 0 {
		return []byte{}
	}
	enc, _ := rlp.EncodeToBytes(cs[i])
	return enc
}

// ToShardID returns the destination shardID of the cxReceipt
// Not used
func (cs CXReceiptsProofs) ToShardID(i int) uint32 {
	return 0
}

// MaxToShardID returns the maximum destination shardID of cxReceipts
// Not used
func (cs CXReceiptsProofs) MaxToShardID() uint32 {
	return 0
}

// GetToShardID get the destination shardID, return error if there is more than one unique shardID
func (cxp *CXReceiptsProof) GetToShardID() (uint32, error) {
	var shardID uint32
	if cxp == nil || len(cxp.Receipts) == 0 {
		return uint32(0), ctxerror.New("[GetShardID] CXReceiptsProof or its receipts is NIL")
	}
	for i, cx := range cxp.Receipts {
		if i == 0 {
			shardID = cx.ToShardID
		} else if shardID == cx.ToShardID {
			continue
		} else {
			return shardID, ctxerror.New("[GetShardID] CXReceiptsProof contains distinct ToShardID")
		}
	}
	return shardID, nil
}
