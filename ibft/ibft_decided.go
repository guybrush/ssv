package ibft

import (
	"github.com/bloxapp/ssv/ibft/pipeline"
	"github.com/bloxapp/ssv/ibft/pipeline/auth"
	"github.com/bloxapp/ssv/ibft/proto"
	ibft_sync "github.com/bloxapp/ssv/ibft/sync"
	"github.com/bloxapp/ssv/network/msgqueue"
	"github.com/bloxapp/ssv/storage/kv"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// processDecidedQueueMessages is listen for all the ibft decided msg's and process them
func (i *ibftImpl) processDecidedQueueMessages() {
	go func() {
		for {
			if decidedMsg := i.msgQueue.PopMessage(msgqueue.DecidedIndexKey(i.GetIdentifier())); decidedMsg != nil {
				i.ProcessDecidedMessage(decidedMsg.SignedMessage)
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

func (i *ibftImpl) validateDecidedMsg(msg *proto.SignedMessage) error {
	p := pipeline.Combine(
		//decided.PrevInstanceDecided(prevInstanceStatus == proto.RoundState_Decided),
		auth.MsgTypeCheck(proto.RoundState_Commit),
		//auth.ValidateLambdas(msg.Message.Lambda, expectedPrevIdentifier),
		auth.AuthorizeMsg(i.ValidatorShare),
		auth.ValidateQuorum(i.ValidatorShare.ThresholdSize()),
	)
	return p.Run(msg)
}

// ProcessDecidedMessage is responsible for processing an incoming decided message.
// If the decided message is known or belong to the current executing instance, do nothing.
// Else perform a sync operation
/* From https://arxiv.org/pdf/2002.03613.pdf
We can omit this if we assume some mechanism external to the consensus algorithm that ensures
synchronization of decided values.
upon receiving a valid hROUND-CHANGE, λi, −, −, −i message from pj ∧ pi has decided
by calling Decide(λi,− , Qcommit) do
	send Qcommit to process pj
*/
func (i *ibftImpl) ProcessDecidedMessage(msg *proto.SignedMessage) {
	if err := i.validateDecidedMsg(msg); err != nil {
		i.logger.Error("received invalid decided message", zap.Error(err), zap.Uint64s("signer ids", msg.SignerIds))
		return
	}

	i.logger.Debug("received valid decided msg", zap.Uint64("seq number", msg.Message.SeqNumber), zap.Uint64s("signer ids", msg.SignerIds))

	// if we already have this in storage, pass
	known, err := i.decidedMsgKnown(msg)
	if err != nil {
		i.logger.Error("can't check if decided msg is known", zap.Error(err))
		return
	}
	if known {
		i.logger.Debug("decided is known, skipped")
		return
	}

	shouldSync, err := i.decidedRequiresSync(msg)
	if err != nil {
		i.logger.Error("can't check decided msg", zap.Error(err))
		return
	}
	if shouldSync {
		if i.currentInstance != nil {
			i.currentInstance.Stop()
		}
		// sync
		s := ibft_sync.NewHistorySync(i.logger, i.ValidatorShare.PublicKey.Serialize(), i.GetIdentifier(), i.network, i.ibftStorage, i.validateDecidedMsg)
		go func() {
			err := s.Start()
			if err != nil {
				i.logger.Error("history sync failed", zap.Error(err))
			}
		}()
	}
}

// HighestKnownDecided returns the highest known decided instance
func (i *ibftImpl) HighestKnownDecided() (*proto.SignedMessage, error) {
	highestKnown, err := i.ibftStorage.GetHighestDecidedInstance(i.GetIdentifier())
	if err != nil && err.Error() != kv.EntryNotFoundError {
		return nil, err
	}
	return highestKnown, nil
}

func (i *ibftImpl) decidedMsgKnown(msg *proto.SignedMessage) (bool, error) {
	found, err := i.ibftStorage.GetDecided(msg.Message.Lambda, msg.Message.SeqNumber)
	if err != nil && err.Error() != kv.EntryNotFoundError {
		return false, errors.Wrap(err, "could not get decided instance from storage")
	}
	return found != nil, nil
}

// decidedForCurrentInstance returns true if msg has same seq number is current instance
func (i *ibftImpl) decidedForCurrentInstance(msg *proto.SignedMessage) bool {
	return i.currentInstance != nil && i.currentInstance.State.SeqNumber == msg.Message.SeqNumber
}

// decidedRequiresSync returns true if:
// 		- highest known seq lower than msg seq
// 		- AND msg is not for current instance
func (i *ibftImpl) decidedRequiresSync(msg *proto.SignedMessage) (bool, error) {
	if i.decidedForCurrentInstance(msg) {
		return false, nil
	}

	if msg.Message.SeqNumber == 0 {
		return false, nil
	}

	highest, err := i.ibftStorage.GetHighestDecidedInstance(msg.Message.Lambda)
	if err != nil {
		if err.Error() == kv.EntryNotFoundError {
			return msg.Message.SeqNumber > 0, nil
		}
		return false, errors.Wrap(err, "could not get highest decided instance from storage")
	}
	return highest.Message.SeqNumber < msg.Message.SeqNumber, nil
}
