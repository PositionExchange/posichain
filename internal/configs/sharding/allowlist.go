package shardingconfig

import (
	"fmt"

	bls_cosi "github.com/PositionExchange/posichain/crypto/bls"
)

type Allowlist struct {
	MaxLimitPerShard int
	BLSPublicKeys    []bls_cosi.PublicKeyWrapper
}

func BLS(pubkeys []string) []bls_cosi.PublicKeyWrapper {
	blsPubkeys := make([]bls_cosi.PublicKeyWrapper, len(pubkeys))
	for i := range pubkeys {
		if key, err := bls_cosi.WrapperPublicKeyFromString(pubkeys[i]); err != nil {
			panic(fmt.Sprintf("invalid bls key: %d:%s error:%s", i, pubkeys[i], err.Error()))
		} else {
			blsPubkeys[i] = *key
		}
	}
	return blsPubkeys
}

// each time to update the allowlist, it requires a hardfork.
// keep same version of mainnet Instance
var mainnetAllowlist_TBD = Allowlist{
	MaxLimitPerShard: 0,
	BLSPublicKeys:    BLS([]string{}),
}

var emptyAllowlist = Allowlist{}
