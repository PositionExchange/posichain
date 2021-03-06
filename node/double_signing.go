package node

import (
	"github.com/PositionExchange/posichain/internal/utils"
	"github.com/PositionExchange/posichain/staking/slash"
	"github.com/ethereum/go-ethereum/rlp"
)

// ProcessSlashCandidateMessage ..
func (node *Node) processSlashCandidateMessage(msgPayload []byte) {
	if !node.IsRunningBeaconChain() {
		return
	}
	candidates := slash.Records{}

	if err := rlp.DecodeBytes(msgPayload, &candidates); err != nil {
		utils.Logger().Error().
			Err(err).Msg("unable to decode slash candidates message")
		return
	}

	if err := node.Blockchain().AddPendingSlashingCandidates(
		candidates,
	); err != nil {
		utils.Logger().Error().
			Err(err).Msg("unable to add slash candidates to pending ")
	}
}
