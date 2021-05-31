package ibft

import (
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/bloxapp/ssv/storage"
	"github.com/bloxapp/ssv/storage/basedb"
	"github.com/bloxapp/ssv/storage/collections"
	validatorstorage "github.com/bloxapp/ssv/validator/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func testIBFTInstance(t *testing.T) *ibftImpl {
	return &ibftImpl{
		//instances: make([]*Instance, 0),
	}
}

func TestCanStartNewInstance(t *testing.T) {
	sks, nodes := GenerateNodes(4)

	tests := []struct {
		name          string
		opts          StartOptions
		storage       collections.Iibft
		initFinished  bool
		expectedError string
	}{
		{
			"valid next instance start",
			StartOptions{
				Identifier: []byte("lambda_10"),
				SeqNumber:  11,
				Duty:       nil,
				ValidatorShare: validatorstorage.Share{
					NodeID:    1,
					PublicKey: validatorPK(sks),
					ShareKey:  sks[1],
					Committee: nodes,
				},
			},
			populatedStorage(t, sks, 10),
			true,
			"",
		},
		{
			"valid first instance",
			StartOptions{
				Identifier: []byte("lambda_0"),
				SeqNumber:  0,
				Duty:       nil,
				ValidatorShare: validatorstorage.Share{
					NodeID:    1,
					PublicKey: validatorPK(sks),
					ShareKey:  sks[1],
					Committee: nodes,
				},
			},
			nil,
			true,
			"",
		},
		{
			"didn't finish initialization",
			StartOptions{
				Identifier: []byte("lambda_0"),
				SeqNumber:  0,
				Duty:       nil,
				ValidatorShare: validatorstorage.Share{
					NodeID:    1,
					PublicKey: validatorPK(sks),
					ShareKey:  sks[1],
					Committee: nodes,
				},
			},
			nil,
			false,
			"iBFT hasn't initialized yet",
		},
		{
			"sequence skips",
			StartOptions{
				Identifier: []byte("lambda_12"),
				SeqNumber:  12,
				Duty:       nil,
				ValidatorShare: validatorstorage.Share{
					NodeID:    1,
					PublicKey: validatorPK(sks),
					ShareKey:  sks[1],
					Committee: nodes,
				},
			},
			populatedStorage(t, sks, 10),
			true,
			"instance seq invalid",
		},
		{
			"past instance",
			StartOptions{
				Identifier: []byte("lambda_10"),
				SeqNumber:  10,
				Duty:       nil,
				ValidatorShare: validatorstorage.Share{
					NodeID:    1,
					PublicKey: validatorPK(sks),
					ShareKey:  sks[1],
					Committee: nodes,
				},
			},
			populatedStorage(t, sks, 10),
			true,
			"instance seq invalid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := testIBFTInstance(t)
			i.initFinished = test.initFinished
			if test.storage != nil {
				i.ibftStorage = test.storage
			} else {
				options := basedb.Options{
					Type:   "badger-memory",
					Logger: zap.L(),
					Path:   "",
				}
				// TODO do we must create new db instnce for each test?
				db, err := storage.GetStorageFactory(options)
				require.NoError(t, err)
				s := collections.NewIbft(db, options.Logger, "attestation")
				i.ibftStorage = &s
			}

			i.ValidatorShare = &test.opts.ValidatorShare
			i.params = &proto.InstanceParams{
				ConsensusParams: proto.DefaultConsensusParams(),
				IbftCommittee:   nodes,
			}
			//i.instances = test.prevInstances
			instanceOpts := i.instanceOptionsFromStartOptions(test.opts)
			//instanceOpts.SeqNumber = test.seqNumber
			err := i.canStartNewInstance(instanceOpts)

			if len(test.expectedError) > 0 {
				require.EqualError(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
