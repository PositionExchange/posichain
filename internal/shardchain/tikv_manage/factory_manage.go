package tikv_manage

import (
	"sync"

	"github.com/PositionExchange/posichain/internal/tikv/statedb_cache"
	"github.com/ethereum/go-ethereum/ethdb"
)

type TiKvFactory interface {
	NewChainDB(shardID uint32) (ethdb.Database, error)
	NewStateDB(shardID uint32) (ethdb.Database, error)
	NewCacheStateDB(shardID uint32) (*statedb_cache.StateDBCacheDatabase, error)
	CloseAllDB()
}

var tikvInit = sync.Once{}
var tikvFactory TiKvFactory

func SetDefaultTiKVFactory(factory TiKvFactory) {
	tikvInit.Do(func() {
		tikvFactory = factory
	})
}

func GetDefaultTiKVFactory() (factory TiKvFactory) {
	return tikvFactory
}
