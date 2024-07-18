package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
)

type GetAllCoinsRequest struct {
	Owner  *sui_types.SuiAddress
	Cursor *sui_types.ObjectID // optional
	Limit  uint                // optional
}

func (s *ImplSuiAPI) GetAllBalances(ctx context.Context, owner *sui_types.SuiAddress) ([]*models.Balance, error) {
	var resp []*models.Balance
	return resp, s.http.CallContext(ctx, &resp, getAllBalances, owner)
}

// start with the first object when cursor is nil
func (s *ImplSuiAPI) GetAllCoins(ctx context.Context, req *GetAllCoinsRequest) (*models.CoinPage, error) {
	var resp models.CoinPage
	return &resp, s.http.CallContext(ctx, &resp, getAllCoins, req.Owner, req.Cursor, req.Limit)
}

type GetBalanceRequest struct {
	Owner    *sui_types.SuiAddress
	CoinType sui_types.ObjectType // optional
}

// GetBalance to use default sui coin(0x2::sui::SUI) when coinType is empty
func (s *ImplSuiAPI) GetBalance(ctx context.Context, req *GetBalanceRequest) (*models.Balance, error) {
	resp := models.Balance{}
	if req.CoinType == "" {
		return &resp, s.http.CallContext(ctx, &resp, getBalance, req.Owner)
	} else {
		return &resp, s.http.CallContext(ctx, &resp, getBalance, req.Owner, req.CoinType)
	}
}

func (s *ImplSuiAPI) GetCoinMetadata(ctx context.Context, coinType string) (*models.SuiCoinMetadata, error) {
	var resp models.SuiCoinMetadata
	return &resp, s.http.CallContext(ctx, &resp, getCoinMetadata, coinType)
}

type GetCoinsRequest struct {
	Owner    *sui_types.SuiAddress
	CoinType *sui_types.ObjectType // optional
	Cursor   *sui_types.ObjectID   // optional
	Limit    uint                  // optional
}

// GetCoins to use default sui coin(0x2::sui::SUI) when coinType is nil
// start with the first object when cursor is nil
func (s *ImplSuiAPI) GetCoins(ctx context.Context, req *GetCoinsRequest) (*models.CoinPage, error) {
	var resp models.CoinPage
	return &resp, s.http.CallContext(ctx, &resp, getCoins, req.Owner, req.CoinType, req.Cursor, req.Limit)
}

func (s *ImplSuiAPI) GetTotalSupply(ctx context.Context, coinType sui_types.ObjectType) (*models.Supply, error) {
	var resp models.Supply
	return &resp, s.http.CallContext(ctx, &resp, getTotalSupply, coinType)
}
