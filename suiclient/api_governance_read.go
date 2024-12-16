package suiclient

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
)

func (s *ClientImpl) GetCommitteeInfo(
	ctx context.Context,
	epoch *sui.BigInt, //optional
) (*CommitteeInfo, error) {
	var resp CommitteeInfo
	return &resp, s.http.CallContext(ctx, &resp, getCommitteeInfo, epoch)
}

func (s *ClientImpl) GetLatestSuiSystemState(ctx context.Context) (*SuiSystemStateSummary, error) {
	var resp SuiSystemStateSummary
	return &resp, s.http.CallContext(ctx, &resp, getLatestSuiSystemState)
}

func (s *ClientImpl) GetReferenceGasPrice(ctx context.Context) (*sui.BigInt, error) {
	var resp sui.BigInt
	return &resp, s.http.CallContext(ctx, &resp, getReferenceGasPrice)
}

func (s *ClientImpl) GetStakes(ctx context.Context, owner *sui.Address) ([]*DelegatedStake, error) {
	var resp []*DelegatedStake
	return resp, s.http.CallContext(ctx, &resp, getStakes, owner)
}

func (s *ClientImpl) GetStakesByIds(ctx context.Context, stakedSuiIds []sui.ObjectId) ([]*DelegatedStake, error) {
	var resp []*DelegatedStake
	return resp, s.http.CallContext(ctx, &resp, getStakesByIds, stakedSuiIds)
}

func (s *ClientImpl) GetValidatorsApy(ctx context.Context) (*ValidatorsApy, error) {
	var resp ValidatorsApy
	return &resp, s.http.CallContext(ctx, &resp, getValidatorsApy)
}
