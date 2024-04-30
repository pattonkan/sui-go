package client

import (
	"context"

	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/types"
)

func (c *Client) GetAllBalances(ctx context.Context, owner *sui_types.SuiAddress) ([]types.Balance, error) {
	var resp []types.Balance
	return resp, c.CallContext(ctx, &resp, getAllBalances, owner)
}

// start with the first object when cursor is nil
func (c *Client) GetAllCoins(
	ctx context.Context,
	owner *sui_types.SuiAddress,
	cursor *sui_types.ObjectID,
	limit uint,
) (*types.CoinPage, error) {
	var resp types.CoinPage
	return &resp, c.CallContext(ctx, &resp, getAllCoins, owner, cursor, limit)
}

// GetBalance to use default sui coin(0x2::sui::SUI) when coinType is empty
func (c *Client) GetBalance(ctx context.Context, owner *sui_types.SuiAddress, coinType string) (*types.Balance, error) {
	resp := types.Balance{}
	if coinType == "" {
		return &resp, c.CallContext(ctx, &resp, getBalance, owner)
	} else {
		return &resp, c.CallContext(ctx, &resp, getBalance, owner, coinType)
	}
}

func (c *Client) GetCoinMetadata(ctx context.Context, coinType string) (*types.SuiCoinMetadata, error) {
	var resp types.SuiCoinMetadata
	return &resp, c.CallContext(ctx, &resp, getCoinMetadata, coinType)
}

// GetCoins to use default sui coin(0x2::sui::SUI) when coinType is nil
// start with the first object when cursor is nil
func (c *Client) GetCoins(
	ctx context.Context,
	owner *sui_types.SuiAddress,
	coinType *string,
	cursor *sui_types.ObjectID,
	limit uint,
) (*types.CoinPage, error) {
	var resp types.CoinPage
	return &resp, c.CallContext(ctx, &resp, getCoins, owner, coinType, cursor, limit)
}

func (c *Client) GetTotalSupply(ctx context.Context, coinType string) (*types.Supply, error) {
	var resp types.Supply
	return &resp, c.CallContext(ctx, &resp, getTotalSupply, coinType)
}
