package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
)

func (s *ImplSuiAPI) DevInspectTransactionBlock(
	ctx context.Context,
	senderAddress *sui_types.SuiAddress,
	txByte serialization.Base64Data,
	gasPrice *models.SafeSuiBigInt[uint64],
	epoch *uint64,
) (*models.DevInspectResults, error) {
	var resp models.DevInspectResults
	return &resp, s.http.CallContext(ctx, &resp, devInspectTransactionBlock, senderAddress, txByte, gasPrice, epoch)
}

func (s *ImplSuiAPI) DryRunTransaction(
	ctx context.Context,
	txBytes serialization.Base64Data,
) (*models.DryRunTransactionBlockResponse, error) {
	var resp models.DryRunTransactionBlockResponse
	return &resp, s.http.CallContext(ctx, &resp, dryRunTransactionBlock, txBytes)
}

func (s *ImplSuiAPI) ExecuteTransactionBlock(
	ctx context.Context, txBytes serialization.Base64Data, signatures []any,
	options *models.SuiTransactionBlockResponseOptions, requestType models.ExecuteTransactionRequestType,
) (*models.SuiTransactionBlockResponse, error) {
	resp := models.SuiTransactionBlockResponse{}
	return &resp, s.http.CallContext(ctx, &resp, executeTransactionBlock, txBytes, signatures, options, requestType)
}
