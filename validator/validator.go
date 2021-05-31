package validator

import (
	"context"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/bloxapp/ssv/storage/collections"
	"github.com/bloxapp/ssv/validator/storage"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/bloxapp/ssv/ibft"
	"github.com/bloxapp/ssv/network"
	"github.com/bloxapp/ssv/network/msgqueue"
	"github.com/bloxapp/ssv/slotqueue"
)

// Options to add in validator struct creation
type Options struct {
	Context                    context.Context
	Logger                     *zap.Logger
	Share                      *storage.Share
	SignatureCollectionTimeout time.Duration
	SlotQueue                  slotqueue.Queue
	Network                    network.Network
	Beacon                     beacon.Beacon
	ETHNetwork                 *core.Network
}

// Validator struct that manages all ibft wrappers
type Validator struct {
	ctx    context.Context
	logger *zap.Logger
	Share  *storage.Share
	//ibftStorage                collections.Iibft
	ethNetwork                 *core.Network
	beacon                     beacon.Beacon
	ibfts                      map[beacon.Role]ibft.IBFT
	msgQueue                   *msgqueue.MessageQueue
	network                    network.Network
	slotQueue                  slotqueue.Queue
	signatureCollectionTimeout time.Duration
}

// New Validator creation
func New(opt Options, ibftStorage collections.Iibft) *Validator {
	logger := opt.Logger.With(zap.String("pubKey", opt.Share.PublicKey.SerializeToHexStr())).
		With(zap.Uint64("node_id", opt.Share.NodeID))

	msgQueue := msgqueue.New()

	ibfts := make(map[beacon.Role]ibft.IBFT)
	ibfts[beacon.RoleAttester] = ibft.New(
		logger,
		ibftStorage,
		opt.Network,
		msgQueue,
		&proto.InstanceParams{
			ConsensusParams: proto.DefaultConsensusParams(),
			IbftCommittee:   opt.Share.Committee,
		},
		opt.Share,
	)
	go ibfts[beacon.RoleAttester].Init()

	return &Validator{
		ctx:                        opt.Context,
		logger:                     logger,
		msgQueue:                   msgQueue,
		Share:                      opt.Share,
		signatureCollectionTimeout: opt.SignatureCollectionTimeout,
		slotQueue:                  opt.SlotQueue,
		network:                    opt.Network,
		ibfts:                      ibfts,
		ethNetwork:                 opt.ETHNetwork,
		beacon:                     opt.Beacon,
	}
}

// Start validator
func (v *Validator) Start() error {
	if err := v.network.SubscribeToValidatorNetwork(v.Share.PublicKey); err != nil {
		return errors.Wrap(err, "failed to subscribe topic")
	}
	go v.startSlotQueueListener()
	go v.listenToNetworkMessages()
	return nil
}

// startSlotQueueListener starts slot queue listener
func (v *Validator) startSlotQueueListener() {
	v.logger.Info("start listening slot queue")

	for {
		slot, duty, ok, err := v.slotQueue.Next(v.Share.PublicKey.Serialize())
		if err != nil {
			v.logger.Error("failed to get next slot data", zap.Error(err))
			continue
		}

		if !ok {
			v.logger.Debug("no duties for slot scheduled")
			continue
		}
		go v.ExecuteDuty(v.ctx, slot, duty)
	}
}

func (v *Validator) listenToNetworkMessages() {
	sigChan := v.network.ReceivedSignatureChan()
	for sigMsg := range sigChan {
		v.msgQueue.AddMessage(&network.Message{
			Lambda:        sigMsg.Message.Lambda,
			SignedMessage: sigMsg,
			Type:          network.NetworkMsg_SignatureType,
		})
	}
}

// getSlotStartTime returns the start time for the given slot  TODO: redundant func (in ssvNode) need to fix
func (v *Validator) getSlotStartTime(slot uint64) time.Time {
	timeSinceGenesisStart := slot * uint64(v.ethNetwork.SlotDurationSec().Seconds())
	start := time.Unix(int64(v.ethNetwork.MinGenesisTime()+timeSinceGenesisStart), 0)
	return start
}
