package staking

import (
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	isValidatorKeyStr     = "Posichain/IsValidator/Key/v1"
	isValidatorStr        = "Posichain/IsValidator/Value/v1"
	collectRewardsStr     = "Posichain/CollectRewards"
	delegateStr           = "Posichain/Delegate"
	unDelegateStr         = "Posichain/UnDelegate"
	firstElectionEpochStr = "Posichain/FirstElectionEpoch/Key/v1"
)

// keys used to retrieve staking related information
var (
	IsValidatorKey        = crypto.Keccak256Hash([]byte(isValidatorKeyStr))
	IsValidator           = crypto.Keccak256Hash([]byte(isValidatorStr))
	CollectRewardsTopic   = crypto.Keccak256Hash([]byte(collectRewardsStr))
	DelegateTopic         = crypto.Keccak256Hash([]byte(delegateStr))
	UnDelegateTopic       = crypto.Keccak256Hash([]byte(unDelegateStr))
	FirstElectionEpochKey = crypto.Keccak256Hash([]byte(firstElectionEpochStr))
)
