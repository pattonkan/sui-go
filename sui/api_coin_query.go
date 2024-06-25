package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
)

func (s *ImplSuiAPI) GetAllBalances(ctx context.Context, owner *sui_types.SuiAddress) ([]*models.Balance, error) {
	var resp []*models.Balance
	return resp, s.http.CallContext(ctx, &resp, getAllBalances, owner)
}

// start with the first object when cursor is nil
func (s *ImplSuiAPI) GetAllCoins(
	ctx context.Context,
	owner *sui_types.SuiAddress,
	cursor *sui_types.ObjectID, // optional
	limit uint, // optional
) (*models.CoinPage, error) {
	var resp models.CoinPage
	return &resp, s.http.CallContext(ctx, &resp, getAllCoins, owner, cursor, limit)
}

// GetBalance to use default sui coin(0x2::sui::SUI) when coinType is empty
func (s *ImplSuiAPI) GetBalance(
	ctx context.Context,
	owner *sui_types.SuiAddress,
	coinType sui_types.ObjectType, // optional
) (
	*models.Balance,
	error,
) {
	resp := models.Balance{}
	if coinType == "" {
		return &resp, s.http.CallContext(ctx, &resp, getBalance, owner)
	} else {
		return &resp, s.http.CallContext(ctx, &resp, getBalance, owner, coinType)
	}
}

func (s *ImplSuiAPI) GetCoinMetadata(ctx context.Context, coinType string) (*models.SuiCoinMetadata, error) {
	var resp models.SuiCoinMetadata
	return &resp, s.http.CallContext(ctx, &resp, getCoinMetadata, coinType)
}

// GetCoins to use default sui coin(0x2::sui::SUI) when coinType is nil
// start with the first object when cursor is nil
func (s *ImplSuiAPI) GetCoins(
	ctx context.Context,
	owner *sui_types.SuiAddress,
	coinType *sui_types.ObjectType, // optional
	cursor *sui_types.ObjectID, // optional
	limit uint, // optional
) (*models.CoinPage, error) {
	var resp models.CoinPage
	return &resp, s.http.CallContext(ctx, &resp, getCoins, owner, coinType, cursor, limit)
}

func (s *ImplSuiAPI) GetTotalSupply(ctx context.Context, coinType sui_types.ObjectType) (*models.Supply, error) {
	var resp models.Supply
	return &resp, s.http.CallContext(ctx, &resp, getTotalSupply, coinType)
}
