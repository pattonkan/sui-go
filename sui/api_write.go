package sui

import (
	"context"

	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/types"
)

func (s *ImplSuiAPI) DevInspectTransactionBlock(
	ctx context.Context,
	senderAddress *sui_types.SuiAddress,
	txByte suiBase64Data,
	gasPrice *types.SafeSuiBigInt[uint64],
	epoch *uint64,
) (*types.DevInspectResults, error) {
	var resp types.DevInspectResults
	return &resp, s.http.CallContext(ctx, &resp, devInspectTransactionBlock, senderAddress, txByte, gasPrice, epoch)
}

func (s *ImplSuiAPI) DryRunTransaction(
	ctx context.Context,
	txBytes suiBase64Data,
) (*types.DryRunTransactionBlockResponse, error) {
	var resp types.DryRunTransactionBlockResponse
	return &resp, s.http.CallContext(ctx, &resp, dryRunTransactionBlock, txBytes)
}

func (s *ImplSuiAPI) ExecuteTransactionBlock(
	ctx context.Context, txBytes suiBase64Data, signatures []any,
	options *types.SuiTransactionBlockResponseOptions, requestType types.ExecuteTransactionRequestType,
) (*types.SuiTransactionBlockResponse, error) {
	resp := types.SuiTransactionBlockResponse{}
	return &resp, s.http.CallContext(ctx, &resp, executeTransactionBlock, txBytes, signatures, options, requestType)
}
