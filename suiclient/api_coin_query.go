package suiclient

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
)

type GetAllCoinsRequest struct {
	Owner  *sui.Address
	Cursor *sui.ObjectId // optional
	Limit  uint          // optional
}

func (s *ClientImpl) GetAllBalances(ctx context.Context, owner *sui.Address) ([]*Balance, error) {
	var resp []*Balance
	return resp, s.http.CallContext(ctx, &resp, getAllBalances, owner)
}

// start with the first object when cursor is nil
func (s *ClientImpl) GetAllCoins(ctx context.Context, req *GetAllCoinsRequest) (*CoinPage, error) {
	var resp CoinPage
	return &resp, s.http.CallContext(ctx, &resp, getAllCoins, req.Owner, req.Cursor, req.Limit)
}

type GetBalanceRequest struct {
	Owner    *sui.Address
	CoinType sui.ObjectType // optional
}

// GetBalance to use default sui coin(0x2::sui::SUI) when coinType is empty
func (s *ClientImpl) GetBalance(ctx context.Context, req *GetBalanceRequest) (*Balance, error) {
	resp := Balance{}
	if req.CoinType == "" {
		return &resp, s.http.CallContext(ctx, &resp, getBalance, req.Owner)
	} else {
		return &resp, s.http.CallContext(ctx, &resp, getBalance, req.Owner, req.CoinType)
	}
}

func (s *ClientImpl) GetCoinMetadata(ctx context.Context, coinType string) (*CoinMetadata, error) {
	var resp CoinMetadata
	return &resp, s.http.CallContext(ctx, &resp, getCoinMetadata, coinType)
}

type GetCoinsRequest struct {
	Owner    *sui.Address
	CoinType *sui.ObjectType // optional
	Cursor   *sui.ObjectId   // optional
	Limit    uint            // optional
}

// GetCoins to use default sui coin(0x2::sui::SUI) when coinType is nil
// start with the first object when cursor is nil
func (s *ClientImpl) GetCoins(ctx context.Context, req *GetCoinsRequest) (*CoinPage, error) {
	var resp CoinPage
	return &resp, s.http.CallContext(ctx, &resp, getCoins, req.Owner, req.CoinType, req.Cursor, req.Limit)
}

func (s *ClientImpl) GetTotalSupply(ctx context.Context, coinType sui.ObjectType) (*Supply, error) {
	var resp Supply
	return &resp, s.http.CallContext(ctx, &resp, getTotalSupply, coinType)
}
