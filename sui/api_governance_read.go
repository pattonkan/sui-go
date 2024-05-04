package sui

import (
	"context"

	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/types"
)

// TODO GetCommitteeInfo

func (s *ImplSuiAPI) GetLatestSuiSystemState(ctx context.Context) (*types.SuiSystemStateSummary, error) {
	var resp types.SuiSystemStateSummary
	return &resp, s.http.CallContext(ctx, &resp, getLatestSuiSystemState)
}

func (s *ImplSuiAPI) GetReferenceGasPrice(ctx context.Context) (*types.SafeSuiBigInt[uint64], error) {
	var resp types.SafeSuiBigInt[uint64]
	return &resp, s.http.CallContext(ctx, &resp, getReferenceGasPrice)
}

func (s *ImplSuiAPI) GetStakes(ctx context.Context, owner *sui_types.SuiAddress) ([]types.DelegatedStake, error) {
	var resp []types.DelegatedStake
	return resp, s.http.CallContext(ctx, &resp, getStakes, owner)
}

func (s *ImplSuiAPI) GetStakesByIds(ctx context.Context, stakedSuiIds []sui_types.ObjectID) ([]types.DelegatedStake, error) {
	var resp []types.DelegatedStake
	return resp, s.http.CallContext(ctx, &resp, getStakesByIds, stakedSuiIds)
}

func (s *ImplSuiAPI) GetValidatorsApy(ctx context.Context) (*types.ValidatorsApy, error) {
	var resp types.ValidatorsApy
	return &resp, s.http.CallContext(ctx, &resp, getValidatorsApy)
}
