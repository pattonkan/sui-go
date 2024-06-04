package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
)

func (s *ImplSuiAPI) GetCommitteeInfo(ctx context.Context, epoch *models.SafeSuiBigInt[sui_types.EpochId]) (*models.CommitteeInfo, error) {
	var resp models.CommitteeInfo
	return &resp, s.http.CallContext(ctx, &resp, getCommitteeInfo, epoch)
}

func (s *ImplSuiAPI) GetLatestSuiSystemState(ctx context.Context) (*models.SuiSystemStateSummary, error) {
	var resp models.SuiSystemStateSummary
	return &resp, s.http.CallContext(ctx, &resp, getLatestSuiSystemState)
}

func (s *ImplSuiAPI) GetReferenceGasPrice(ctx context.Context) (*models.SafeSuiBigInt[uint64], error) {
	var resp models.SafeSuiBigInt[uint64]
	return &resp, s.http.CallContext(ctx, &resp, getReferenceGasPrice)
}

func (s *ImplSuiAPI) GetStakes(ctx context.Context, owner *sui_types.SuiAddress) ([]*models.DelegatedStake, error) {
	var resp []*models.DelegatedStake
	return resp, s.http.CallContext(ctx, &resp, getStakes, owner)
}

func (s *ImplSuiAPI) GetStakesByIds(ctx context.Context, stakedSuiIds []sui_types.ObjectID) (
	[]*models.DelegatedStake,
	error,
) {
	var resp []*models.DelegatedStake
	return resp, s.http.CallContext(ctx, &resp, getStakesByIds, stakedSuiIds)
}

func (s *ImplSuiAPI) GetValidatorsApy(ctx context.Context) (*models.ValidatorsApy, error) {
	var resp models.ValidatorsApy
	return &resp, s.http.CallContext(ctx, &resp, getValidatorsApy)
}
