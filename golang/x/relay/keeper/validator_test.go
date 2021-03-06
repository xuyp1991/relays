package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

func (s *KeeperSuite) TestGetConfs() {
	header := s.Fixtures.ValidatorTestCases.ValidateProof[0].Proof.ConfirmingHeader
	bestKnown := s.Fixtures.ValidatorTestCases.ValidateProof[0].BestKnown

	// errors if Best Known Digest is not found
	confs, err := s.Keeper.getConfs(s.Context, header)
	s.Equal(sdk.CodeType(types.BadHash256Digest), err.Code())
	s.Equal(uint32(0), confs)

	// errors if Best Known Digest header is not found
	s.Keeper.setBestKnownDigest(s.Context, bestKnown.Hash)

	confs, err = s.Keeper.getConfs(s.Context, header)
	s.Equal(sdk.CodeType(types.UnknownBlock), err.Code())
	s.Equal(uint32(0), confs)

	// success
	s.Keeper.ingestHeader(s.Context, bestKnown)

	confs, err = s.Keeper.getConfs(s.Context, header)
	s.SDKNil(err)
	s.Equal(uint32(4), confs)
}

func (s *KeeperSuite) TestValidateProof() {
	proofCases := s.Fixtures.ValidatorTestCases.ValidateProof
	proof := proofCases[0].Proof

	// errors if LCA is not found
	err := s.Keeper.validateProof(s.Context, proof)
	s.Equal(sdk.CodeType(types.BadHash256Digest), err.Code())

	// errors if link is not found
	s.Keeper.setLastReorgLCA(s.Context, proofCases[0].LCA)

	err = s.Keeper.validateProof(s.Context, proof)
	s.Equal(sdk.CodeType(types.NotAncestor), err.Code())

	for i := range proofCases {
		// Store lots of stuff
		s.Keeper.setLastReorgLCA(s.Context, proofCases[i].LCA)
		s.Keeper.ingestHeader(s.Context, proofCases[i].Proof.ConfirmingHeader)
		s.Keeper.setLink(s.Context, proofCases[i].Proof.ConfirmingHeader)

		if proofCases[i].Error != 0 {
			err := s.Keeper.validateProof(s.Context, proofCases[i].Proof)
			s.Equal(sdk.CodeType(proofCases[i].Error), err.Code())
		} else {
			err := s.Keeper.validateProof(s.Context, proofCases[i].Proof)
			s.Nil(err)
		}
	}
}

func (s *KeeperSuite) TestCheckRequestsFilled() {
	tc := s.Fixtures.ValidatorTestCases.CheckRequestsFilled
	validProof := s.Fixtures.ValidatorTestCases.ValidateProof[0]

	s.Keeper.setLastReorgLCA(s.Context, validProof.LCA)
	s.Keeper.ingestHeader(s.Context, validProof.Proof.ConfirmingHeader)
	s.Keeper.setLink(s.Context, validProof.Proof.ConfirmingHeader)
	s.Keeper.ingestHeader(s.Context, validProof.BestKnown)
	requestErr := s.Keeper.setRequest(s.Context, []byte{}, []byte{}, 0, 4, types.Local, nil)
	s.Nil(requestErr)

	// errors if getConfs fails
	_, err := s.Keeper.checkRequestsFilled(s.Context, tc[0].FilledRequests)
	s.Equal(sdk.CodeType(types.BadHash256Digest), err.Code())

	s.Keeper.setBestKnownDigest(s.Context, validProof.BestKnown.Hash)

	// errors if checkRequest errors
	// deactivate request
	activeErr := s.Keeper.setRequestState(s.Context, types.RequestID{}, false)
	s.SDKNil(activeErr)

	_, err = s.Keeper.checkRequestsFilled(s.Context, tc[0].FilledRequests)
	s.Equal(sdk.CodeType(types.ClosedRequest), err.Code())

	// reactivate request
	activeErr = s.Keeper.setRequestState(s.Context, types.RequestID{}, true)
	s.SDKNil(activeErr)

	for i := range tc {
		_, err := s.Keeper.checkRequestsFilled(s.Context, tc[i].FilledRequests)
		if tc[i].Error != 0 {
			s.Equal(sdk.CodeType(tc[i].Error), err.Code())
		} else {
			s.SDKNil(err)
		}
	}

	// errors if number of confirmations is less than the number of confirmations on the request
	requestErr = s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 5, types.Local, nil)
	s.Nil(requestErr)

	copiedRequest := tc[0].FilledRequests
	copiedRequest.Filled[0].ID = types.RequestID{0, 0, 0, 0, 0, 0, 0, 1}
	_, err = s.Keeper.checkRequestsFilled(s.Context, copiedRequest)
	s.Equal(sdk.CodeType(types.NotEnoughConfs), err.Code())
}
