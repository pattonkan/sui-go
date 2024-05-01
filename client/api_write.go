package client

import (
	"context"

	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/types"
)

func (c *Client) DevInspectTransactionBlock(
	ctx context.Context,
	senderAddress *sui_types.SuiAddress,
	txByte suiBase64Data,
	gasPrice *types.SafeSuiBigInt[uint64],
	epoch *uint64,
) (*types.DevInspectResults, error) {
	var resp types.DevInspectResults
	return &resp, c.CallContext(ctx, &resp, devInspectTransactionBlock, senderAddress, txByte, gasPrice, epoch)
}

func (c *Client) DryRunTransaction(
	ctx context.Context,
	txBytes suiBase64Data,
) (*types.DryRunTransactionBlockResponse, error) {
	var resp types.DryRunTransactionBlockResponse
	return &resp, c.CallContext(ctx, &resp, dryRunTransactionBlock, txBytes)
}

func (c *Client) ExecuteTransactionBlock(
	ctx context.Context, txBytes suiBase64Data, signatures []any,
	options *types.SuiTransactionBlockResponseOptions, requestType types.ExecuteTransactionRequestType,
) (*types.SuiTransactionBlockResponse, error) {
	resp := types.SuiTransactionBlockResponse{}
	return &resp, c.CallContext(ctx, &resp, executeTransactionBlock, txBytes, signatures, options, requestType)
}
