package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
)

func (s *ImplSuiAPI) DevInspectTransactionBlock(
	ctx context.Context,
	senderAddress *sui_types.SuiAddress,
	txByte sui_types.Base64Data,
	gasPrice *models.SafeSuiBigInt[uint64], // optional
	epoch *uint64, // optional
) (*models.DevInspectResults, error) {
	var resp models.DevInspectResults
	return &resp, s.http.CallContext(ctx, &resp, devInspectTransactionBlock, senderAddress, txByte, gasPrice, epoch)
}

func (s *ImplSuiAPI) DryRunTransaction(
	ctx context.Context,
	txBytes sui_types.Base64Data,
) (*models.DryRunTransactionBlockResponse, error) {
	var resp models.DryRunTransactionBlockResponse
	return &resp, s.http.CallContext(ctx, &resp, dryRunTransactionBlock, txBytes)
}

func (s *ImplSuiAPI) ExecuteTransactionBlock(
	ctx context.Context,
	txBytes sui_types.Base64Data,
	signatures []*sui_signer.Signature,
	options *models.SuiTransactionBlockResponseOptions,
	requestType models.ExecuteTransactionRequestType,
) (*models.SuiTransactionBlockResponse, error) {
	resp := models.SuiTransactionBlockResponse{}
	return &resp, s.http.CallContext(ctx, &resp, executeTransactionBlock, txBytes, signatures, options, requestType)
}
