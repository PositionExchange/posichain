package consensus

import (
	"bytes"
	"encoding/binary"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/harmony-one/bls/ffi/go/bls"
	msg_pb "github.com/harmony-one/harmony/api/proto/message"
	bls_cosi "github.com/harmony-one/harmony/crypto/bls"
	"github.com/harmony-one/harmony/internal/utils"
	"github.com/harmony-one/harmony/p2p"
	"github.com/harmony-one/harmony/p2p/host"
)

// PbftPhase  PBFT phases: pre-prepare, prepare and commit
type PbftPhase int

// Enum for PbftPhase
const (
	Announce PbftPhase = iota
	Prepare
	Commit
)

// Mode determines whether a node is in normal or viewchanging mode
type Mode int

// Enum for node Mode
const (
	Normal Mode = iota
	ViewChanging
	Syncing
)

// PbftMode contains mode and viewID of viewchanging
type PbftMode struct {
	mode   Mode
	viewID uint32
	mux    sync.Mutex
}

// Mode return the current node mode
func (pm *PbftMode) Mode() Mode {
	return pm.mode
}

//String print mode string
func (mode Mode) String() string {
	if mode == Normal {
		return "Normal"
	} else if mode == ViewChanging {
		return "ViewChanging"
	} else if mode == Syncing {
		return "Sycning"
	}
	return "Unknown"
}

// String print phase string
func (phase PbftPhase) String() string {
	if phase == Announce {
		return "Announce"
	} else if phase == Prepare {
		return "Prepare"
	} else if phase == Commit {
		return "Commit"
	}
	return "Unknown"
}

// SetMode set the node mode as required
func (pm *PbftMode) SetMode(m Mode) {
	pm.mux.Lock()
	defer pm.mux.Unlock()
	pm.mode = m
}

// ViewID return the current viewchanging id
func (pm *PbftMode) ViewID() uint32 {
	return pm.viewID
}

// SetViewID sets the viewchanging id accordingly
func (pm *PbftMode) SetViewID(viewID uint32) {
	pm.mux.Lock()
	defer pm.mux.Unlock()
	pm.viewID = viewID
}

// GetViewID returns the current viewchange viewID
func (pm *PbftMode) GetViewID() uint32 {
	return pm.viewID
}

// switchPhase will switch PbftPhase to nextPhase if the desirePhase equals the nextPhase
func (consensus *Consensus) switchPhase(desirePhase PbftPhase, override bool) {
	utils.GetLogInstance().Debug("switchPhase: ", "desirePhase", desirePhase, "myPhase", consensus.phase, "override", override)
	if override {
		consensus.phase = desirePhase
		return
	}

	var nextPhase PbftPhase
	switch consensus.phase {
	case Announce:
		nextPhase = Prepare
	case Prepare:
		nextPhase = Commit
	case Commit:
		nextPhase = Announce
	}
	if nextPhase == desirePhase {
		consensus.phase = nextPhase
	}
}

// GetNextLeaderKey uniquely determine who is the leader for given viewID
func (consensus *Consensus) GetNextLeaderKey() *bls.PublicKey {
	idx := consensus.getIndexOfPubKey(consensus.LeaderPubKey)
	if idx == -1 {
		utils.GetLogInstance().Warn("GetNextLeaderKey: currentLeaderKey not found", "key", consensus.LeaderPubKey.SerializeToHexStr())
	}
	idx = (idx + 1) % len(consensus.PublicKeys)
	return consensus.PublicKeys[idx]
}

func (consensus *Consensus) getIndexOfPubKey(pubKey *bls.PublicKey) int {
	for k, v := range consensus.PublicKeys {
		if v.IsEqual(pubKey) {
			return k
		}
	}
	return -1
}

// ResetViewChangeState reset the state for viewchange
func (consensus *Consensus) ResetViewChangeState() {
	consensus.mode.SetMode(Normal)
	bhpBitmap, _ := bls_cosi.NewMask(consensus.PublicKeys, nil)
	nilBitmap, _ := bls_cosi.NewMask(consensus.PublicKeys, nil)
	viewIDBitmap, _ := bls_cosi.NewMask(consensus.PublicKeys, nil)
	consensus.bhpBitmap = bhpBitmap
	consensus.nilBitmap = nilBitmap
	consensus.viewIDBitmap = viewIDBitmap
	consensus.m1Payload = []byte{}

	consensus.bhpSigs = map[common.Address]*bls.Sign{}
	consensus.nilSigs = map[common.Address]*bls.Sign{}
	consensus.viewIDSigs = map[common.Address]*bls.Sign{}
}

func createTimeout() map[TimeoutType]*utils.Timeout {
	timeouts := make(map[TimeoutType]*utils.Timeout)
	timeouts[timeoutConsensus] = utils.NewTimeout(phaseDuration)
	timeouts[timeoutViewChange] = utils.NewTimeout(viewChangeDuration)
	timeouts[timeoutBootstrap] = utils.NewTimeout(bootstrapDuration)
	return timeouts
}

// startViewChange send a  new view change
func (consensus *Consensus) startViewChange(viewID uint32) {
	if consensus.disableViewChange {
		return
	}
	consensus.consensusTimeout[timeoutConsensus].Stop()
	consensus.consensusTimeout[timeoutBootstrap].Stop()
	consensus.mode.SetMode(ViewChanging)
	consensus.mode.SetViewID(viewID)
	consensus.LeaderPubKey = consensus.GetNextLeaderKey()

	diff := viewID - consensus.viewID
	duration := time.Duration(int64(diff) * int64(viewChangeDuration))
	utils.GetLogInstance().Info("startViewChange", "viewID", viewID, "timeoutDuration", duration, "nextLeader", consensus.LeaderPubKey.SerializeToHexStr()[:10])

	msgToSend := consensus.constructViewChangeMessage()
	consensus.host.SendMessageToGroups([]p2p.GroupID{p2p.NewGroupIDByShardID(p2p.ShardID(consensus.ShardID))}, host.ConstructP2pMessage(byte(17), msgToSend))

	consensus.consensusTimeout[timeoutViewChange].SetDuration(duration)
	consensus.consensusTimeout[timeoutViewChange].Start()
	utils.GetLogger().Debug("start view change timeout", "viewID", consensus.viewID, "block", consensus.blockNum, "viewChangingID", consensus.mode.ViewID())
}

// new leader send new view message
func (consensus *Consensus) startNewView() {
	utils.GetLogInstance().Info("startNewView", "viewID", consensus.mode.GetViewID())
	consensus.mode.SetMode(Normal)
	consensus.switchPhase(Announce, false)

	msgToSend := consensus.constructNewViewMessage()
	consensus.host.SendMessageToGroups([]p2p.GroupID{p2p.NewGroupIDByShardID(p2p.ShardID(consensus.ShardID))}, host.ConstructP2pMessage(byte(17), msgToSend))
}

func (consensus *Consensus) onViewChange(msg *msg_pb.Message) {
	senderKey, validatorAddress, err := consensus.verifyViewChangeSenderKey(msg)
	if err != nil {
		utils.GetLogInstance().Debug("onViewChange verifySenderKey failed", "error", err)
		return
	}

	recvMsg, err := ParseViewChangeMessage(msg)
	if err != nil {
		utils.GetLogInstance().Warn("onViewChange unable to parse viewchange message")
		return
	}
	newLeaderKey := recvMsg.LeaderPubkey
	if !consensus.PubKey.IsEqual(newLeaderKey) {
		return
	}

	utils.GetLogInstance().Warn("onViewChange received", "viewChangeID", recvMsg.ViewID, "myCurrentID", consensus.viewID, "ValidatorAddress", consensus.SelfAddress)

	if consensus.blockNum > recvMsg.BlockNum {
		return
	}

	if consensus.blockNum < recvMsg.BlockNum {
		utils.GetLogger().Warn("new leader has lower blocknum", "have", consensus.blockNum, "got", recvMsg.BlockNum)
		return
	}

	if consensus.mode.Mode() == ViewChanging && consensus.mode.GetViewID() > recvMsg.ViewID {
		return
	}
	if err = verifyMessageSig(senderKey, msg); err != nil {
		utils.GetLogInstance().Debug("onViewChange Failed to verify sender's signature", "error", err)
		return
	}

	consensus.vcLock.Lock()
	defer consensus.vcLock.Unlock()

	// add self m1 or m2 type message signature and bitmap
	_, ok1 := consensus.nilSigs[consensus.SelfAddress]
	_, ok2 := consensus.bhpSigs[consensus.SelfAddress]
	if !(ok1 || ok2) {
		// add own signature for newview message
		preparedMsgs := consensus.pbftLog.GetMessagesByTypeSeq(msg_pb.MessageType_PREPARED, recvMsg.BlockNum)
		preparedMsg := consensus.pbftLog.FindMessageByMaxViewID(preparedMsgs)
		if preparedMsg == nil {
			sign := consensus.priKey.SignHash(NIL)
			consensus.nilSigs[consensus.SelfAddress] = sign
			consensus.nilBitmap.SetKey(consensus.PubKey, true)
		} else {
			msgToSign := append(preparedMsg.BlockHash[:], preparedMsg.Payload...)
			consensus.bhpSigs[consensus.SelfAddress] = consensus.priKey.SignHash(msgToSign)
			consensus.bhpBitmap.SetKey(consensus.PubKey, true)
		}
	}
	// add self m3 type message signature and bitmap
	_, ok3 := consensus.viewIDSigs[consensus.SelfAddress]
	if !ok3 {
		viewIDHash := make([]byte, 4)
		binary.LittleEndian.PutUint32(viewIDHash, recvMsg.ViewID)
		sign := consensus.priKey.SignHash(viewIDHash)
		consensus.viewIDSigs[consensus.SelfAddress] = sign
		consensus.viewIDBitmap.SetKey(consensus.PubKey, true)
	}

	if len(consensus.viewIDSigs) >= consensus.Quorum() {
		return
	}

	// m2 type message
	if len(recvMsg.Payload) == 0 {
		_, ok := consensus.nilSigs[validatorAddress]
		if ok {
			utils.GetLogInstance().Debug("onViewChange already received m2 message from the validator", "validatorAddress", validatorAddress)
			return
		}

		if !recvMsg.ViewchangeSig.VerifyHash(senderKey, NIL) {
			utils.GetLogInstance().Warn("onViewChange failed to verify signature for m2 type viewchange message")
			return
		}
		consensus.nilSigs[validatorAddress] = recvMsg.ViewchangeSig
		consensus.nilBitmap.SetKey(recvMsg.SenderPubkey, true) // Set the bitmap indicating that this validator signed.
	} else { // m1 type message
		_, ok := consensus.bhpSigs[validatorAddress]
		if ok {
			utils.GetLogInstance().Debug("onViewChange already received m1 message from the validator", "validatorAddress", validatorAddress)
			return
		}
		if !recvMsg.ViewchangeSig.VerifyHash(recvMsg.SenderPubkey, recvMsg.Payload) {
			utils.GetLogInstance().Warn("onViewChange failed to verify signature for m1 type viewchange message")
			return
		}

		// first time receive m1 type message, need verify validity of prepared message
		if len(consensus.m1Payload) == 0 || !bytes.Equal(consensus.m1Payload, recvMsg.Payload) {
			if len(recvMsg.Payload) <= 32 {
				utils.GetLogger().Debug("m1 recvMsg payload not enough length", "len", len(recvMsg.Payload))
				return
			}
			blockHash := recvMsg.Payload[:32]
			aggSig, mask, err := consensus.readSignatureBitmapPayload(recvMsg.Payload, 32)
			if err != nil {
				utils.GetLogger().Error("m1 recvMsg payload read error", "error", err)
				return
			}
			// check has 2f+1 signature in m1 type message
			if count := utils.CountOneBits(mask.Bitmap); count < consensus.Quorum() {
				utils.GetLogger().Debug("not have enough signature", "need", consensus.Quorum(), "have", count)
				return
			}

			// Verify the multi-sig for prepare phase
			if !aggSig.VerifyHash(mask.AggregatePublic, blockHash[:]) {
				utils.GetLogInstance().Warn("onViewChange failed to verify multi signature for m1 prepared payload", "blockHash", blockHash)
				return
			}
			if len(consensus.m1Payload) == 0 {
				consensus.m1Payload = append(recvMsg.Payload[:0:0], recvMsg.Payload...)
			}
		}
		consensus.bhpSigs[validatorAddress] = recvMsg.ViewchangeSig
		consensus.bhpBitmap.SetKey(recvMsg.SenderPubkey, true) // Set the bitmap indicating that this validator signed.
	}

	// check and add viewID (m3 type) message signature
	_, ok := consensus.viewIDSigs[validatorAddress]
	if ok {
		utils.GetLogInstance().Debug("onViewChange already received m3 viewID message from the validator", "validatorAddress", validatorAddress)
		return
	}
	viewIDHash := make([]byte, 4)
	binary.LittleEndian.PutUint32(viewIDHash, recvMsg.ViewID)
	if !recvMsg.ViewidSig.VerifyHash(recvMsg.SenderPubkey, viewIDHash) {
		utils.GetLogInstance().Warn("onViewChange failed to verify viewID signature", "viewID", recvMsg.ViewID)
		return
	}
	consensus.viewIDSigs[validatorAddress] = recvMsg.ViewidSig
	consensus.viewIDBitmap.SetKey(recvMsg.SenderPubkey, true) // Set the bitmap indicating that this validator signed.

	if len(consensus.viewIDSigs) >= consensus.Quorum() {
		consensus.mode.SetMode(Normal)
		consensus.LeaderPubKey = consensus.PubKey
		consensus.ResetState()
		if len(consensus.m1Payload) == 0 {
			go func() {
				consensus.ReadySignal <- struct{}{}
			}()
		} else {
			consensus.phase = Commit
			copy(consensus.blockHash[:], consensus.m1Payload[:32])
			aggSig, mask, err := consensus.readSignatureBitmapPayload(recvMsg.Payload, 32)
			if err != nil {
				utils.GetLogger().Error("readSignatureBitmapPayload fail", "error", err)
				return
			}
			consensus.aggregatedPrepareSig = aggSig
			consensus.prepareBitmap = mask

			// Leader sign commit message
			blockNumHash := make([]byte, 8)
			binary.LittleEndian.PutUint64(blockNumHash, consensus.blockNum)
			commitPayload := append(blockNumHash, consensus.blockHash[:]...)
			consensus.commitSigs[consensus.SelfAddress] = consensus.priKey.SignHash(commitPayload)
		}

		consensus.mode.SetViewID(recvMsg.ViewID)
		msgToSend := consensus.constructNewViewMessage()

		utils.GetLogInstance().Warn("onViewChange", "sent newview message", len(msgToSend))
		consensus.host.SendMessageToGroups([]p2p.GroupID{p2p.NewGroupIDByShardID(p2p.ShardID(consensus.ShardID))}, host.ConstructP2pMessage(byte(17), msgToSend))

		consensus.viewID = recvMsg.ViewID
		consensus.ResetViewChangeState()
		consensus.consensusTimeout[timeoutViewChange].Stop()
		consensus.consensusTimeout[timeoutConsensus].Start()
		utils.GetLogger().Debug("new leader start consensus timeout and stop view change timeout", "viewID", consensus.viewID, "block", consensus.blockNum, "viewChangingID", consensus.mode.ViewID())
		utils.GetLogger().Debug("I am the new leader", "myKey", consensus.PubKey.SerializeToHexStr()[:20], "viewID", consensus.viewID, "block", consensus.blockNum)
	}
	utils.GetLogInstance().Debug("onViewChange", "numSigs", len(consensus.viewIDSigs), "needed", consensus.Quorum())

}

// TODO: move to consensus_leader.go later
func (consensus *Consensus) onNewView(msg *msg_pb.Message) {
	utils.GetLogInstance().Debug("onNewView received new view message")
	senderKey, _, err := consensus.verifyViewChangeSenderKey(msg)
	if err != nil {
		utils.GetLogInstance().Warn("onNewView verifySenderKey failed", "error", err)
		return
	}
	recvMsg, err := consensus.ParseNewViewMessage(msg)
	if err != nil {
		utils.GetLogInstance().Warn("onViewChange unable to parse viewchange message")
		return
	}

	if err = verifyMessageSig(senderKey, msg); err != nil {
		utils.GetLogInstance().Error("onNewView failed to verify new leader's signature", "error", err)
		return
	}
	consensus.vcLock.Lock()
	defer consensus.vcLock.Unlock()

	if recvMsg.M3AggSig == nil || recvMsg.M3Bitmap == nil {
		utils.GetLogInstance().Error("onNewView M3AggSig or M3Bitmap is nil")
		return
	}
	m3Sig := recvMsg.M3AggSig
	m3Mask := recvMsg.M3Bitmap
	viewIDHash := make([]byte, 4)
	binary.LittleEndian.PutUint32(viewIDHash, recvMsg.ViewID)
	// check total number of sigs >= 2f+1
	if count := utils.CountOneBits(m3Mask.Bitmap); count < consensus.Quorum() {
		utils.GetLogger().Debug("not have enough signature", "need", consensus.Quorum(), "have", count)
		return
	}

	if !m3Sig.VerifyHash(m3Mask.AggregatePublic, viewIDHash) {
		utils.GetLogInstance().Warn("onNewView unable to verify aggregated signature of m3 payload", "m3Sig", m3Sig.SerializeToHexStr()[:10], "m3Mask", m3Mask.Bitmap, "viewID", recvMsg.ViewID)
		return
	}

	m2Mask := recvMsg.M2Bitmap
	if recvMsg.M2AggSig != nil {
		m2Sig := recvMsg.M2AggSig
		if !m2Sig.VerifyHash(m2Mask.AggregatePublic, NIL) {
			utils.GetLogInstance().Warn("onNewView unable to verify aggregated signature of m2 payload")
			return
		}
	}

	if m3Mask == nil || m2Mask == nil {
		utils.GetLogInstance().Error("onNewView m3Mask or m2Mask is nil")
		return
	}

	// check when M3 sigs > M2 sigs, then M1 (recvMsg.Payload) should not be empty
	if utils.CountOneBits(m3Mask.Bitmap) > utils.CountOneBits(m2Mask.Bitmap) {
		if len(recvMsg.Payload) <= 32 {
			utils.GetLogger().Debug("we should have m1 message payload non-empty and valid")
			return
		}
	}

	// check validity of m1 type payload
	if len(recvMsg.Payload) > 32 {
		blockHash := recvMsg.Payload[:32]
		aggSig, mask, err := consensus.readSignatureBitmapPayload(recvMsg.Payload, 32)
		if err != nil {
			utils.GetLogger().Error("unable to read signature/bitmap", "error", err)
			return
		}
		if !aggSig.VerifyHash(mask.AggregatePublic, blockHash) {
			utils.GetLogInstance().Warn("onNewView failed to verify signature for prepared message")
			return
		}
		copy(consensus.blockHash[:], blockHash)
		consensus.aggregatedPrepareSig = aggSig
		consensus.prepareBitmap = mask

		// create prepared message from newview
		preparedMsg := PbftMessage{MessageType: msg_pb.MessageType_PREPARED, ViewID: recvMsg.ViewID, BlockNum: recvMsg.BlockNum}
		preparedMsg.BlockHash = common.Hash{}
		copy(preparedMsg.BlockHash[:], blockHash[:])
		preparedMsg.Payload = make([]byte, len(recvMsg.Payload)-32)
		copy(preparedMsg.Payload[:], recvMsg.Payload[32:])
		preparedMsg.SenderPubkey = senderKey
		consensus.pbftLog.AddMessage(&preparedMsg)
	}

	// newView message verified success, override my state
	consensus.viewID = recvMsg.ViewID
	consensus.mode.SetViewID(recvMsg.ViewID)
	consensus.LeaderPubKey = senderKey
	consensus.ResetViewChangeState()

	// change view and leaderKey to keep in sync with network
	if consensus.blockNum != recvMsg.BlockNum {
		utils.GetLogger().Debug("new leader changed", "newLeaderKey", consensus.LeaderPubKey.SerializeToHexStr()[:20], "viewID", consensus.viewID, "myBlock", consensus.blockNum, "newViewBlockNum", recvMsg.BlockNum)
		return
	}

	if len(recvMsg.Payload) > 32 {
		// Construct and send the commit message
		blockNumHash := make([]byte, 8)
		binary.LittleEndian.PutUint64(blockNumHash, consensus.blockNum)
		commitPayload := append(blockNumHash, consensus.blockHash[:]...)
		msgToSend := consensus.constructCommitMessage(commitPayload)

		utils.GetLogInstance().Info("onNewView === commit", "sent commit message", len(msgToSend), "viewID", consensus.viewID)
		consensus.host.SendMessageToGroups([]p2p.GroupID{p2p.NewGroupIDByShardID(p2p.ShardID(consensus.ShardID))}, host.ConstructP2pMessage(byte(17), msgToSend))
		consensus.switchPhase(Commit, true)
	} else {
		consensus.ResetState()
		utils.GetLogInstance().Info("onNewView === announce")
	}
	utils.GetLogger().Debug("new leader changed", "newLeaderKey", consensus.LeaderPubKey.SerializeToHexStr()[:20], "viewID", consensus.viewID, "block", consensus.blockNum)
	utils.GetLogger().Debug("validator start consensus timeout and stop view change timeout", "viewID", consensus.viewID, "block", consensus.blockNum)
	consensus.consensusTimeout[timeoutConsensus].Start()
	consensus.consensusTimeout[timeoutViewChange].Stop()
}
