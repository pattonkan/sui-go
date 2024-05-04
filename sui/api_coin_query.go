package sui

import (
	"context"

	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/types"
)

func (s *ImplSuiAPI) GetAllBalances(ctx context.Context, owner *sui_types.SuiAddress) ([]types.Balance, error) {
	var resp []types.Balance
	return resp, s.http.CallContext(ctx, &resp, getAllBalances, owner)
}

// start with the first object when cursor is nil
func (s *ImplSuiAPI) GetAllCoins(
	ctx context.Context,
	owner *sui_types.SuiAddress,
	cursor *sui_types.ObjectID,
	limit uint,
) (*types.CoinPage, error) {
	var resp types.CoinPage
	return &resp, s.http.CallContext(ctx, &resp, getAllCoins, owner, cursor, limit)
}

// GetBalance to use default sui coin(0x2::sui::SUI) when coinType is empty
func (s *ImplSuiAPI) GetBalance(ctx context.Context, owner *sui_types.SuiAddress, coinType string) (*types.Balance, error) {
	resp := types.Balance{}
	if coinType == "" {
		return &resp, s.http.CallContext(ctx, &resp, getBalance, owner)
	} else {
		return &resp, s.http.CallContext(ctx, &resp, getBalance, owner, coinType)
	}
}

func (s *ImplSuiAPI) GetCoinMetadata(ctx context.Context, coinType string) (*types.SuiCoinMetadata, error) {
	var resp types.SuiCoinMetadata
	return &resp, s.http.CallContext(ctx, &resp, getCoinMetadata, coinType)
}

// GetCoins to use default sui coin(0x2::sui::SUI) when coinType is nil
// start with the first object when cursor is nil
func (s *ImplSuiAPI) GetCoins(
	ctx context.Context,
	owner *sui_types.SuiAddress,
	coinType *string,
	cursor *sui_types.ObjectID,
	limit uint,
) (*types.CoinPage, error) {
	var resp types.CoinPage
	return &resp, s.http.CallContext(ctx, &resp, getCoins, owner, coinType, cursor, limit)
}

func (s *ImplSuiAPI) GetTotalSupply(ctx context.Context, coinType string) (*types.Supply, error) {
	var resp types.Supply
	return &resp, s.http.CallContext(ctx, &resp, getTotalSupply, coinType)
}
