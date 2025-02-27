package validator

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"go.uber.org/zap"
	"log"
)

func (v *Validator) verifyPartialSignature(signature []byte, root []byte, ibftID uint64, committiee map[uint64]*proto.Node) error {
	if val, found := committiee[ibftID]; found {
		pk := &bls.PublicKey{}
		if err := pk.Deserialize(val.Pk); err != nil {
			return errors.Wrap(err, "could not deserialized pk")
		}
		sig := &bls.Sign{}
		if err := sig.Deserialize(signature); err != nil {
			return errors.Wrap(err, "could not deserialized signature")
		}

		// verify
		if !sig.VerifyByte(pk, root) {
			return errors.Errorf("could not verify signature from iBFT member %d", ibftID)
		}
		return nil
	}
	return errors.Errorf("could not find iBFT member %d", ibftID)
}

// signDuty signs the duty after iBFT came to consensus
func (v *Validator) signDuty(ctx context.Context, decidedValue []byte, role beacon.Role, duty *ethpb.DutiesResponse_Duty, shareKey *bls.SecretKey, ) ([]byte, []byte, *proto.InputValue, error) {
	// sign input value
	var sig []byte
	var root []byte
	retValueStruct := &proto.InputValue{}
	var err error
	switch role {
	case beacon.RoleAttester:
		s := &proto.InputValue_Attestation{}
		if err := json.Unmarshal(decidedValue, s); err != nil {
			return nil, nil, nil, errors.Wrap(err, "failed to marshal attestation")
		}
		signedAttestation, r, e := v.beacon.SignAttestation(ctx, s.Attestation.Data, duty, shareKey)
		if e != nil{
			return nil, nil, nil, errors.Wrap(err, "failed to sign attestation")
		}
		retValueStruct.SignedData = s
		retValueStruct.GetAttestation().Signature = signedAttestation.Signature
		retValueStruct.GetAttestation().AggregationBits = signedAttestation.AggregationBits
		err = e
		sig = signedAttestation.GetSignature()
		root = r
	case beacon.RoleAggregator:
		s := &proto.InputValue_Aggregation{}
		if err := json.Unmarshal(decidedValue, s); err != nil {
			return nil, nil, nil, errors.Wrap(err, "failed to marshal aggregator")
		}
		signedAggregation, e := v.beacon.SignAggregation(ctx, s.Aggregation.Message, v.Share.ShareKey)
		if e != nil{
			return nil, nil, nil, errors.Wrap(err, "failed to sign attestation")
		}
		retValueStruct.SignedData = s
		retValueStruct.GetAggregation().Signature = signedAggregation.Signature
		retValueStruct.GetAggregation().Message = signedAggregation.Message
		err = e
		sig = signedAggregation.GetSignature()
	case beacon.RoleProposer:
		s := &proto.InputValue_Block{}
		if err := json.Unmarshal(decidedValue, s); err != nil {
			return nil, nil, nil, errors.Wrap(err, "failed to marshal aggregator")
		}

		signedProposal, e := v.beacon.SignProposal(ctx, nil, s.Block.Block, v.Share.ShareKey)
		if e != nil{
			return nil, nil, nil, errors.Wrap(err, "failed to sign attestation")
		}

		retValueStruct.SignedData = s
		retValueStruct.GetBlock().Signature = signedProposal.Signature
		retValueStruct.GetBlock().Block = signedProposal.Block
		err = e
		sig = signedProposal.GetSignature()
	default:
		return nil, nil, nil, errors.New("unsupported role, can't sign")
	}
	return sig, root, retValueStruct, err
}

// reconstructAndBroadcastSignature reconstructs the received signatures from other
// nodes and broadcasts the reconstructed signature to the beacon-chain
func (v *Validator) reconstructAndBroadcastSignature(ctx context.Context, logger *zap.Logger, signatures map[uint64][]byte, root []byte, inputValue *proto.InputValue, role beacon.Role, duty *ethpb.DutiesResponse_Duty) error {
	// Reconstruct signatures
	signature, err := threshold.ReconstructSignatures(signatures)
	if err != nil {
		return errors.Wrap(err, "failed to reconstruct signatures")
	}
	// verify reconstructed sig
	if res := signature.VerifyByte(v.Share.PublicKey, root); !res {
		return errors.New("could not reconstruct a valid signature")
	}

	logger.Info("signatures successfully reconstructed", zap.String("signature", base64.StdEncoding.EncodeToString(signature.Serialize())), zap.Int("signature count", len(signatures)))

	// Submit validation to beacon node
	switch role {
	case beacon.RoleAttester:
		logger.Info("submitting attestation")
		inputValue.GetAttestation().Signature = signature.Serialize()
		log.Printf("%s, %d\n", inputValue.GetAttestation(), duty.GetValidatorIndex())
		if err := v.beacon.SubmitAttestation(ctx, inputValue.GetAttestation(), duty.GetValidatorIndex(), v.Share.PublicKey); err != nil {
			return errors.Wrap(err, "failed to broadcast attestation")
		}
	case beacon.RoleAggregator:
		inputValue.GetAggregation().Signature = signature.Serialize()
		if err := v.beacon.SubmitAggregation(ctx, inputValue.GetAggregation()); err != nil {
			return errors.Wrap(err, "failed to broadcast aggregation")
		}
	case beacon.RoleProposer:
		inputValue.GetBlock().Signature = signature.Serialize()
		if err := v.beacon.SubmitProposal(ctx, inputValue.GetBlock()); err != nil {
			return errors.Wrap(err, "failed to broadcast proposal")
		}
	default:
		return errors.New("role is undefined, can't reconstruct signature")
	}
	return nil
}
