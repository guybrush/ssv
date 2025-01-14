package prysmgrpc

import (
	"context"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// GetProposalData implements Beacon interface
func (b *prysmGRPC) GetProposalData(ctx context.Context, slot uint64, shareKey *bls.SecretKey) (*ethpb.BeaconBlock, error) {
	randaoReveal, err := b.signRandaoReveal(ctx, slot, shareKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign randao reveal")
	}

	block, err := b.validatorClient.GetBlock(ctx, &ethpb.BlockRequest{
		Slot:         slot,
		RandaoReveal: randaoReveal,
		Graffiti:     b.graffiti,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get block")
	}

	return block, nil
}

// SignProposal implements Beacon interface
func (b *prysmGRPC) SignProposal(ctx context.Context, domain *ethpb.DomainResponse, block *ethpb.BeaconBlock, shareKey *bls.SecretKey) (*ethpb.SignedBeaconBlock, error) {
	// TODO: Check this
	/*if err := b.preBlockSignValidations(ctx, block); err != nil {
		return nil, errors.Wrapf(err, "failed block safety check for slot %d", block.Slot)
	}*/

	sig, err := b.signBlock(ctx, domain, block, shareKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign block")
	}

	// TODO: Check this
	/*if err := b.postBlockSignUpdate(ctx, block, domain); err != nil {
		return nil, errors.Wrapf(err, "failed post block signing validations for slot %d", blk.Block.Slot)
	}*/

	return &ethpb.SignedBeaconBlock{
		Block:     block,
		Signature: sig,
	}, nil
}

// SubmitProposal implements Beacon interface
func (b *prysmGRPC) SubmitProposal(ctx context.Context, block *ethpb.SignedBeaconBlock) error {
	if _, err := b.validatorClient.ProposeBlock(ctx, block); err != nil {
		return errors.Wrap(err, "failed to propose block")
	}

	return nil
}

// signRandaoReveal signs randao reveal with randao domain and private key.
func (b *prysmGRPC) signRandaoReveal(ctx context.Context, slot uint64, shareKey *bls.SecretKey) ([]byte, error) {
	domain, err := b.domainData(ctx, slot, params.BeaconConfig().DomainRandao[:])
	if err != nil {
		return nil, errors.Wrap(err, "failed to get domain data")
	}

	if domain == nil {
		return nil, errors.New("domain data is empty")
	}

	root, err := helpers.ComputeSigningRoot(b.network.EstimatedEpochAtSlot(slot), domain.SignatureDomain)
	if err != nil {
		return nil, err
	}

	return shareKey.SignByte(root[:]).Serialize(), nil
}

func (b *prysmGRPC) signBlock(ctx context.Context, domain *ethpb.DomainResponse, block *ethpb.BeaconBlock, shareKey *bls.SecretKey) ([]byte, error) {
	var err error
	if domain == nil{
		domain, err = b.domainData(ctx, block.GetSlot(), params.BeaconConfig().DomainBeaconProposer[:])
		if err != nil {
			return nil, errors.Wrap(err, "failed to get domain data")
		}
	}

	// TODO: A patch to randao signature!!
	root, err := helpers.ComputeSigningRoot(block, domain.SignatureDomain)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compute signing root")
	}

	return shareKey.SignByte(root[:]).Serialize(), nil
}
